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
	stateStart state = iota
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

	state   state
	series  *Series
	seq     *keyboard.Sequence
	chosen  harmony.Degree
	correct bool
	changes []dex.Change // what the end screen says, in any language
	debug   bool
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
	g.series = NewSeries(g.rng, seriesLength)
	g.ask()
}

func (g *game) ask() {
	q, _ := g.series.Current()
	g.state = stateAsking
	g.play(questionNotes(q))
}

func (g *game) answer(i int) {
	q, _ := g.series.Current()
	g.chosen = q.Choices[i]
	g.correct = g.series.Answer(g.chosen)
	g.state = stateRevealed
	if !g.correct {
		g.play(comparisonNotes(q, g.chosen))
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

// buttonRect places answer `i`.
func buttonRect(i int) (x, y, w, h float32) {
	return 40 + float32(i)*190, 180, 180, 48
}

// clickedButton returns the answer clicked this frame, or -1.
func (g *game) clickedButton() int {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return -1
	}
	px, py := ebiten.CursorPosition()
	cx, cy := float32(float64(px)/g.scale), float32(float64(py)/g.scale)
	for i := range 3 {
		x, y, w, h := buttonRect(i)
		if cx >= x && cx < x+w && cy >= y && cy < y+h {
			return i
		}
	}
	return -1
}

func pressedAnswerKey() int {
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if inpututil.IsKeyJustPressed(k) {
			return i
		}
	}
	return -1
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
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
	case stateStart, stateEnd:
		if proceed {
			g.startSeries()
		}

	case stateAsking:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			q, _ := g.series.Current()
			g.play(questionNotes(q))
		}
		i := g.clickedButton()
		if i < 0 {
			i = pressedAnswerKey()
		}
		if i >= 0 {
			g.answer(i)
		}

	case stateRevealed:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			q, _ := g.series.Current()
			if g.correct {
				g.play(questionNotes(q))
			} else {
				g.play(comparisonNotes(q, g.chosen))
			}
		}
		// A click on the buttons is the answer just given, not a
		// request to move on.
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			(inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.clickedButton() < 0) {
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
	case stateStart:
		g.print(screen, w.start, 40, 160, foreground)
		g.print(screen, w.keys, 40, 200, dim)

	case stateAsking, stateRevealed:
		q, _ := g.series.Current()
		n, total := g.series.Position()
		head := fmt.Sprintf(w.question, n, total)
		if q.Retry {
			head += " " + w.retry
		}
		g.print(screen, head, 40, 40, dim)
		g.print(screen, fmt.Sprintf(w.which, g.lang.note(q.Tonic)), 40, 90, foreground)

		for i, d := range q.Choices {
			x, y, bw, bh := buttonRect(i)
			fill := button
			if g.state == stateRevealed {
				switch {
				case d == q.Mode:
					fill = rightColor
				case d == g.chosen:
					fill = wrongColor
				}
			}
			g.canvas(screen).rect(x, y, bw, bh, fill)
			g.print(screen, fmt.Sprintf("%d  %s", i+1, g.lang.mode(d)), float64(x)+12, float64(y)+14, foreground)
		}

		if g.state == stateAsking {
			g.print(screen, w.answerKeys+"   "+w.replay, 40, 236, dim)
		} else {
			verdict := w.right
			if !g.correct {
				verdict = fmt.Sprintf(w.wrong, g.lang.mode(q.Mode))
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
	right, _ := system.Mode(q.Mode)
	chosen, _ := system.Mode(g.chosen)
	r := reveal{
		on:      true,
		tonic:   q.Tonic,
		right:   right.At(q.Tonic),
		chosen:  chosen.At(q.Tonic),
		mistake: !g.correct,
	}

	// Signs on the keys whatever the notation: « si bémol » does not fit
	// on a key, and what the keys show is read, not spoken.
	signs := g.lang.namer.WithNotation(naming.Signs)
	if r.mistake {
		g.label(&r, signs, q.Tonic, chosen)
	}
	// The right mode is spelled last, so that it names the shared
	// classes: that is the scale the player is asked to learn.
	g.label(&r, signs, q.Tonic, right)
	return r
}

// label spells the classes of one mode on its tonic into r.
func (g *game) label(r *reveal, n *naming.Namer, tonic harmony.PitchClass, p harmony.ScalePattern) {
	t, err := harmony.NewTonality(tonic, p)
	if err != nil {
		return
	}
	spelled := n.WithTonality(t)
	for c := range p.At(tonic).Classes() {
		r.labels[c] = spelled.Name(c)
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
