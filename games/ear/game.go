package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/ArnaudCalmettes/gohar/synth"
)

// signsTTF holds the glyphs Go Regular lacks: ♭ ♮ ♯ 𝄫 𝄪. A two
// kilobyte cut of Noto Music, embedded rather than read from the system
// so that the browser build has it too. See fonts/README.md.
//
//go:embed fonts/signs.ttf
var signsTTF []byte

const (
	screenWidth  = 640
	screenHeight = 360
	seriesLength = 10
)

type state uint8

const (
	stateMenu state = iota
	stateAsking
	stateRevealed
	stateEnd
)

var (
	background = color.RGBA{0x1c, 0x1c, 0x22, 0xff}
	foreground = color.RGBA{0xe8, 0xe6, 0xe0, 0xff}
	dim        = color.RGBA{0x8a, 0x88, 0x84, 0xff}
	button     = color.RGBA{0x33, 0x33, 0x3d, 0xff}
	rightColor = color.RGBA{0x3d, 0x7a, 0x4f, 0xff}
	wrongColor = color.RGBA{0x8a, 0x3b, 0x3b, 0xff}
)

// activities are what the menu offers, in the order of oreille.md: the
// tetrachords before the modes they build. All open from the start, for
// now; levels come with the settings.
var activities = []Activity{
	tetrachords{shapes: naturalTetrachords},
	modes{system: harmony.NaturalMajor},
	modes{system: harmony.NaturalMajor, all: true},
}

// keysDown is what is sounding right now, from any source. Written on
// the sources' goroutines, copied once per tick by Update.
//
// A fixed array rather than a map: the copy is then a plain assignment,
// and nothing allocates on the path of a key press.
type keysDown struct {
	mu   sync.Mutex
	down [128]bool
}

func (k *keysDown) set(e keyboard.Event) {
	if e.Key < 0 || e.Key >= len(k.down) {
		return
	}
	k.mu.Lock()
	k.down[e.Key] = e.Down
	k.mu.Unlock()
}

func (k *keysDown) snapshot(dst *[128]bool) {
	k.mu.Lock()
	*dst = k.down
	k.mu.Unlock()
}

type game struct {
	lang    language
	other   language // the one L switches to
	engine  synth.Instrument
	dex     *dex.Dex
	dexPath string
	rng     *rand.Rand
	midi    string // the MIDI input's name, empty without one

	face  *font
	small *font
	scale float64 // physical pixels per logical unit, set by Layout
	keys  keysDown
	down  [128]bool // this tick's copy of keys
	piano piano

	state    state
	activity Activity
	series   *Series
	seq      *keyboard.Sequence
	chosen   int // the index of the choice made
	correct  bool
	changes  []dex.Change // what the end screen says, in any language
	debug    bool
}

func newGame(lang, other language, engine synth.Instrument, d *dex.Dex, dexPath string, rng *rand.Rand, midi string) (*game, error) {
	face, err := newFont(16)
	if err != nil {
		return nil, err
	}
	small, err := newFont(10)
	if err != nil {
		return nil, err
	}
	return &game{
		lang: lang, other: other, engine: engine, dex: d, dexPath: dexPath, rng: rng, midi: midi,
		face: face, small: small, scale: 1,
	}, nil
}

// toggleNotation switches both languages between signs and words, so
// that L keeps the choice N made.
func (g *game) toggleNotation() {
	n := naming.Words
	if g.lang.namer.Notation() == naming.Words {
		n = naming.Signs
	}
	g.lang = g.lang.withNotation(n)
	g.other = g.other.withNotation(n)
}

// onKey is the one path from any source to the sound and the screen.
// It returns at once, as keyboard.Source asks.
func (g *game) onKey(e keyboard.Event) {
	if e.Down {
		g.engine.NoteOn(e.Key, e.Velocity, e.At)
	} else {
		g.engine.NoteOff(e.Key, e.At)
	}
	g.keys.set(e)
}

// play replaces whatever sequence is running. Close releases what the
// old one held, so a skipped question never leaves its pedal behind.
func (g *game) play(notes []keyboard.Note) {
	g.stop()
	g.seq = keyboard.NewSequence("ear", notes)
	if err := g.seq.Listen(g.onKey); err != nil {
		log.Println(err)
	}
}

func (g *game) stop() {
	if g.seq != nil {
		g.seq.Close()
		g.seq = nil
	}
}

func (g *game) startSeries() {
	g.series = NewSeries(g.activity, g.rng, seriesLength)
	g.ask()
}

func (g *game) ask() {
	q, _ := g.series.Current()
	g.state = stateAsking
	g.play(g.activity.Sound(q))
}

func (g *game) answer(i int) {
	q, _ := g.series.Current()
	g.chosen = i
	g.correct = g.series.Answer(i)
	g.state = stateRevealed
	if !g.correct {
		g.play(g.activity.Correction(q, i))
	}
}

func (g *game) next() {
	if g.series.Next() {
		g.ask()
		return
	}
	g.finish()
}

// finish sends the one report of the series and saves the collection.
func (g *game) finish() {
	g.stop()
	g.changes = g.dex.Apply(g.series.Report(time.Now()))
	if err := saveDex(g.dexPath, g.dex); err != nil {
		log.Println("saving the dex:", err)
	}
	g.state = stateEnd
}

// buttonRect places answer `i` of `n`, side by side across the width:
// three names of modes, four of tetrachords, seven degrees.
func buttonRect(i, n int) (x, y, w, h float32) {
	const left, width, gutter = 40, 560, 10
	w = (width - float32(n-1)*gutter) / float32(n)
	return left + float32(i)*(w+gutter), 180, w, 48
}

// clickedButton returns the answer among `n` clicked this frame, or -1.
func (g *game) clickedButton(n int) int {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return -1
	}
	px, py := ebiten.CursorPosition()
	cx, cy := float32(float64(px)/g.scale), float32(float64(py)/g.scale)
	for i := range n {
		x, y, w, h := buttonRect(i, n)
		if cx >= x && cx < x+w && cy >= y && cy < y+h {
			return i
		}
	}
	return -1
}

// answerKeys are the digits, one per choice: seven at most, the degrees
// of a scale.
var answerKeys = []ebiten.Key{
	ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4,
	ebiten.Key5, ebiten.Key6, ebiten.Key7,
}

// pressedAnswerKey returns the answer among `n` typed this frame, or -1.
func pressedAnswerKey(n int) int {
	for i, k := range answerKeys[:min(n, len(answerKeys))] {
		if inpututil.IsKeyJustPressed(k) {
			return i
		}
	}
	return -1
}

func (g *game) Update() error {
	// Escape alone: Ebiten names keys by their place on a US keyboard,
	// and the place of Q is A on an AZERTY one.
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.stop()
		return ebiten.Termination
	}
	g.keys.snapshot(&g.down)
	g.piano.update(&g.down)

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.debug = !g.debug
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.lang, g.other = g.other, g.lang
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.toggleNotation()
	}

	proceed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	switch g.state {
	case stateMenu:
		i := g.clickedButton(len(activities))
		if i < 0 {
			i = pressedAnswerKey(len(activities))
		}
		if i >= 0 {
			g.activity = activities[i]
			g.startSeries()
		}

	case stateEnd:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = stateMenu
		} else if proceed {
			g.startSeries()
		}

	case stateAsking:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			q, _ := g.series.Current()
			g.play(g.activity.Sound(q))
		}
		q, _ := g.series.Current()
		i := g.clickedButton(len(q.Choices))
		if i < 0 {
			i = pressedAnswerKey(len(q.Choices))
		}
		if i >= 0 {
			g.answer(i)
		}

	case stateRevealed:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			q, _ := g.series.Current()
			if g.correct {
				g.play(g.activity.Sound(q))
			} else {
				g.play(g.activity.Correction(q, g.chosen))
			}
		}
		// A click on the buttons is the answer just given, not a
		// request to move on.
		q, _ := g.series.Current()
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			(inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.clickedButton(len(q.Choices)) < 0) {
			g.next()
		}
	}
	return nil
}

func (g *game) canvas(screen *ebiten.Image) canvas {
	return canvas{dst: screen, scale: g.scale}
}

func (g *game) print(screen *ebiten.Image, s string, x, y float64, c color.Color) {
	g.canvas(screen).text(s, g.face, x, y, c)
}

func (g *game) printWith(screen *ebiten.Image, f *font, s string, x, y float64, c color.Color) {
	g.canvas(screen).text(s, f, x, y, c)
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(background)
	w := g.lang.words

	switch g.state {
	case stateMenu:
		g.print(screen, w.menu, 40, 90, foreground)
		for i, a := range activities {
			x, y, bw, bh := buttonRect(i, len(activities))
			g.canvas(screen).rect(x, y, bw, bh, button)
			g.print(screen, fmt.Sprintf("%d  %s", i+1, g.lang.activity(a)), float64(x)+12, float64(y)+14, foreground)
		}
		g.print(screen, w.keys, 40, 236, dim)

	case stateAsking, stateRevealed:
		q, _ := g.series.Current()
		n, total := g.series.Position()
		head := fmt.Sprintf(w.question, n, total)
		if q.Retry {
			head += " " + w.retry
		}
		g.print(screen, head, 40, 40, dim)
		g.print(screen, g.activity.Prompt(g.lang, q), 40, 90, foreground)

		for i, choice := range q.Choices {
			x, y, bw, bh := buttonRect(i, len(q.Choices))
			fill := button
			if g.state == stateRevealed {
				switch i {
				case q.Answer:
					fill = rightColor
				case g.chosen:
					fill = wrongColor
				}
			}
			g.canvas(screen).rect(x, y, bw, bh, fill)
			g.buttonText(screen, i+1, g.lang.label(choice), x, y, bw)
		}

		if g.state == stateAsking {
			g.print(screen, w.answerKeys+"   "+w.replay, 40, 236, dim)
		} else {
			verdict := w.right
			if !g.correct {
				verdict = fmt.Sprintf(w.wrong, g.lang.notion(q.Right()))
			}
			g.print(screen, verdict, 40, 130, foreground)
			g.print(screen, w.next+"   "+w.replay, 40, 236, dim)
		}

	case stateEnd:
		g.print(screen, w.end, 40, 40, foreground)
		var lines []string
		for _, c := range g.changes {
			if line := g.lang.change(c); line != "" {
				lines = append(lines, line)
			}
		}
		if len(lines) == 0 {
			lines = append(lines, w.nothingNew)
		}
		// Seven discoveries at most fit above the keyboard.
		for i, line := range lines {
			g.print(screen, line, 40, 80+float64(i)*22, foreground)
		}
		g.print(screen, w.again, 40, 236, dim)
	}

	g.drawFooter(screen)
}

// buttonText writes a choice on its button. When the number and the
// name do not fit side by side, the number goes on top and the name
// below it in the small face: seven modes on one line leave about
// seventy units a button, and « 7  mixolydien » needs more.
func (g *game) buttonText(screen *ebiten.Image, number int, name string, x, y, w float32) {
	const pad = 12
	c := g.canvas(screen)
	line := fmt.Sprintf("%d  %s", number, name)
	if tw, _ := c.measure(line, g.face); tw <= float64(w)-2*pad {
		c.text(line, g.face, float64(x)+pad, float64(y)+14, foreground)
		return
	}
	c.text(fmt.Sprint(number), g.face, float64(x)+pad/2, float64(y)+4, foreground)
	c.text(name, g.small, float64(x)+pad/2, float64(y)+28, foreground)
}

// drawFooter shows the keyboard, and on H the delay figures of the
// test under load.
func (g *game) drawFooter(screen *ebiten.Image) {
	g.piano.draw(g.canvas(screen), g.reveal(), g.small)
	g.corner(screen)

	if g.debug {
		h := g.engine.Histogram()
		g.printWith(screen, g.small, fmt.Sprintf("n %d  p50 %v  p99 %v  max %v  fps %.0f",
			h.Total(), h.Quantile(0.5), h.Quantile(0.99), h.Quantile(1), ebiten.ActualFPS()),
			40, 342, dim)
	}
}

// reveal says what the keyboard may show. Nothing but what sounds
// until the answer is out: marking the mode during the question would
// be showing the answer.
func (g *game) reveal() reveal {
	if g.state != stateRevealed {
		return reveal{}
	}
	q, _ := g.series.Current()
	d := g.activity.Show(q, g.chosen)
	r := reveal{
		on:      true,
		tonic:   d.tonic,
		right:   d.right.At(d.tonic),
		chosen:  d.chosen.At(d.tonic),
		mistake: d.mistake,
	}

	// Signs on the keys whatever the notation: « si bémol » does not fit
	// on a key, and what the keys show is read, not spoken.
	signs := g.lang.namer.WithNotation(naming.Signs)
	if r.mistake {
		g.label(&r, signs, d.tonic, d.chosen)
	}
	// The right shape is spelled last, so that it names the shared
	// classes: that is the one the player is asked to learn.
	g.label(&r, signs, d.tonic, d.right)
	return r
}

// label spells the classes of one shape on its tonic into r: in its own
// tonality when it has seven notes, by the default convention when it
// has fewer, a tetrachord having no key to be spelled in.
func (g *game) label(r *reveal, n *naming.Namer, tonic harmony.PitchClass, p harmony.ScalePattern) {
	if t, err := harmony.NewTonality(tonic, p); err == nil {
		n = n.WithTonality(t)
	}
	for c := range p.At(tonic).Classes() {
		r.labels[c] = n.Name(c)
	}
}

// corner shows the MIDI input in small, right aligned, cut to fit.
//
// The client name alone: ALSA names a port "client:port 24:0", and the
// part after the colon repeats the client and adds numbers nobody
// reads. The full name is what -port matches and what keys -list shows.
func (g *game) corner(screen *ebiten.Image) {
	label := g.lang.words.noMIDI
	if g.midi != "" {
		label, _, _ = strings.Cut(g.midi, ":")
	}
	const room = 240
	c := g.canvas(screen)
	if w, _ := c.measure(label, g.small); w > room {
		r := []rune(label)
		for len(r) > 0 {
			r = r[:len(r)-1]
			if w, _ := c.measure(string(r)+"…", g.small); w <= room {
				break
			}
		}
		label = string(r) + "…"
	}
	w, _ := c.measure(label, g.small)
	c.text(label, g.small, screenWidth-12-w, 10, dim)
}

// Layout renders at the window's physical resolution, keeping the
// logical grid's proportions: see canvas.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	w, h := float64(outsideWidth)*s, float64(outsideHeight)*s
	g.scale = min(w/screenWidth, h/screenHeight)
	return int(screenWidth * g.scale), int(screenHeight * g.scale)
}
