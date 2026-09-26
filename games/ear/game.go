package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/ArnaudCalmettes/gohar/synth"
)

// signsTTF holds the glyphs Go Regular lacks: ♭ ♮ ♯ and the ellipsis.
// A four kilobyte cut of DejaVu Sans, embedded rather than read from
// the system so that the browser build has it too. See fonts/README.md.
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
// the sources' goroutines, read in Draw.
type keysDown struct {
	mu   sync.Mutex
	down map[int]bool
}

func (k *keysDown) set(e keyboard.Event) {
	k.mu.Lock()
	if e.Down {
		k.down[e.Key] = true
	} else {
		delete(k.down, e.Key)
	}
	k.mu.Unlock()
}

func (k *keysDown) sorted() []int {
	k.mu.Lock()
	defer k.mu.Unlock()
	keys := make([]int, 0, len(k.down))
	for key := range k.down {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

type game struct {
	lang    language
	other   language // the one L switches to
	engine  *synth.Engine
	dex     *dex.Dex
	dexPath string
	rng     *rand.Rand
	midi    string // the MIDI input's name, empty without one

	face  text.Face
	small text.Face
	keys  keysDown

	state   state
	series  *Series
	seq     *keyboard.Sequence
	chosen  harmony.Degree
	correct bool
	changes []dex.Change // what the end screen says, in any language
	debug   bool
}

func newGame(lang, other language, engine *synth.Engine, d *dex.Dex, dexPath string, rng *rand.Rand, midi string) (*game, error) {
	face, err := newFace(16)
	if err != nil {
		return nil, err
	}
	small, err := newFace(10)
	if err != nil {
		return nil, err
	}
	return &game{
		lang: lang, other: other, engine: engine, dex: d, dexPath: dexPath, rng: rng, midi: midi,
		face: face, small: small,
		keys: keysDown{down: map[int]bool{}},
	}, nil
}

// newFace is Go Regular, falling back on the signs font for what it
// lacks. A multi face picks the first face holding each glyph.
func newFace(size float64) (text.Face, error) {
	regular, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}
	signs, err := text.NewGoTextFaceSource(bytes.NewReader(signsTTF))
	if err != nil {
		return nil, err
	}
	return text.NewMultiFace(
		&text.GoTextFace{Source: regular, Size: size},
		&text.GoTextFace{Source: signs, Size: size},
	)
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
func clickedButton() int {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return -1
	}
	cx, cy := ebiten.CursorPosition()
	for i := range 3 {
		x, y, w, h := buttonRect(i)
		if float32(cx) >= x && float32(cx) < x+w && float32(cy) >= y && float32(cy) < y+h {
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
		i := clickedButton()
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
			(inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && clickedButton() < 0) {
			g.next()
		}
	}
	return nil
}

func (g *game) print(screen *ebiten.Image, s string, x, y float64, c color.Color) {
	g.printWith(screen, g.face, s, x, y, c)
}

func (g *game) printWith(screen *ebiten.Image, face text.Face, s string, x, y float64, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(screen, s, face, op)
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
			vector.DrawFilledRect(screen, x, y, bw, bh, fill, false)
			g.print(screen, fmt.Sprintf("%d  %s", i+1, g.lang.mode(d)), float64(x)+12, float64(y)+14, foreground)
		}

		if g.state == stateAsking {
			g.print(screen, w.answerKeys+"   "+w.replay, 40, 260, dim)
		} else {
			verdict := w.right
			if !g.correct {
				verdict = fmt.Sprintf(w.wrong, g.lang.mode(q.Mode))
			}
			g.print(screen, verdict, 40, 130, foreground)
			g.print(screen, w.next+"   "+w.replay, 40, 260, dim)
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
		for i, line := range lines {
			g.print(screen, line, 40, 90+float64(i)*26, foreground)
		}
		g.print(screen, w.again, 40, 300, dim)
	}

	g.drawFooter(screen)
}

// drawFooter shows what is sounding, and on H the delay figures of the
// test under load.
func (g *game) drawFooter(screen *ebiten.Image) {
	var names []string
	for _, key := range g.keys.sorted() {
		names = append(names, g.lang.note(harmony.PitchClass(key%12)))
	}
	g.print(screen, strings.Join(names, " "), 40, 316, dim)
	g.corner(screen)

	if g.debug {
		h := g.engine.Histogram()
		g.printWith(screen, g.small, fmt.Sprintf("n %d  p50 %v  p99 %v  max %v  fps %.0f",
			h.Total(), h.Quantile(0.5), h.Quantile(0.99), h.Quantile(1), ebiten.ActualFPS()),
			40, 340, dim)
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
	if w, _ := text.Measure(label, g.small, 0); w > room {
		r := []rune(label)
		for len(r) > 0 {
			r = r[:len(r)-1]
			if w, _ := text.Measure(string(r)+"…", g.small, 0); w <= room {
				break
			}
		}
		label = string(r) + "…"
	}
	w, _ := text.Measure(label, g.small, 0)
	g.printWith(screen, g.small, label, screenWidth-12-w, 10, dim)
}

func (g *game) Layout(int, int) (int, int) { return screenWidth, screenHeight }
