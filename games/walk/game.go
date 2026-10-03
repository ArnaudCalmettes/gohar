package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// lookahead is how far ahead the band is scheduled. Ebiten's loop
// ticks every 16 ms, too coarse to play a beat on time: notes are
// queued a window ahead and the synth applies them to the sample (see
// "Le son" in docs/walk.md).
const lookahead = 100 * time.Millisecond

// The count-in: two bars of snaps on 2 and 4, the configuration
// validated by ear.
const (
	perBar  = 4
	countIn = 2 * perBar
)

// The chart, laid out as a Real Book page: four bars a row.
const (
	barsPerRow = 4
	chartX     = 20
	chartY     = 70
	barW       = 150
	rowH       = 56 // room for the keyboard below
)

// A cell is one chord written in a bar, at the beat it starts on.
type cell struct {
	name string
	beat int // in the bar, from 0
}

// game is the shell of Walk With Me: the chart, the count-in, the band,
// the player's hands, and the marks of the first palier. Two phases:
// with the tempo, the band plays and the marker marks; without, the
// chart waits for the player (practice).
type game struct {
	title      string
	bpm        float64
	choruses   int
	grid       analysis.Changes
	namer      *naming.Namer
	bars       [][]cell // the chart, bar by bar; empty for a chord held from before
	chorusLen  int      // beats in a chorus
	band       *band
	events     chan keyboard.Event
	midi       string // the keyboard's name, empty without one
	practicing bool   // the phase without tempo

	running  bool
	m        Metronome
	beats    []Beat
	next     int // the next beat to schedule
	marker   *Marker
	marks    map[int]BeatKind // the arrivals marked, by beat
	practice *practice
	lastMark string // what the last note of the bass was
	last     keyboard.Event
	down     [128]bool // the keys held, as the events tell
	piano    piano
	walker   walker
	line     []int     // the reference line, in demo: a key per beat
	heard    []Note    // its notes not yet sounded, for the marker
	rec      *recorder // nil without -record

	fonts fonts
	scale float64
}

func newGame(title string, grid analysis.Changes, bpm float64, choruses int, bd *band, midi string, practicing bool) (*game, error) {
	fs, err := newFonts()
	if err != nil {
		return nil, err
	}
	namer := namerFor(grid)
	return &game{
		title:      title,
		bpm:        bpm,
		choruses:   choruses,
		grid:       grid,
		namer:      namer,
		bars:       chart(grid, namer),
		chorusLen:  len(Expect(grid, NewMetronome(time.Time{}, bpm, perBar), 1)),
		band:       bd,
		events:     make(chan keyboard.Event, 64),
		midi:       midi,
		practicing: practicing,
		piano:      piano{split: bd.split},
		fonts:      fs,
	}, nil
}

// onKey is the MIDI callback: it sounds the key at once and hands the
// event to the game loop, never waiting. A full channel drops the
// event for the display only: the sound has already gone.
func (g *game) onKey(e keyboard.Event) {
	g.band.key(e)
	select {
	case g.events <- e:
	default:
	}
}

func (g *game) start(now time.Time) {
	g.running = true
	g.lastMark = ""
	g.walker = walker{} // each run starts from a walk
	if g.practicing {
		// No time: any metronome numbers the beats.
		g.practice = newPractice(Expect(g.grid, NewMetronome(now, g.bpm, perBar), 1))
		return
	}
	beat := time.Duration(float64(time.Minute) / g.bpm)
	g.m = NewMetronome(now.Add(lookahead+countIn*beat), g.bpm, perBar)
	g.beats = Expect(g.grid, g.m, g.choruses)
	g.next = -countIn
	g.marker = NewMarker(FirstPalier, g.m, g.beats)
	g.marks = map[int]BeatKind{}
	g.line, g.heard = nil, nil
	if g.rec != nil {
		g.rec.run(now, g.bpm, g.band.demo)
	}
	if g.band.demo {
		seed := uint64(now.UnixNano())
		g.line = Walk(g.beats, rand.New(rand.NewPCG(seed, seed>>32|1)))
	}
}

func (g *game) stop(now time.Time) {
	g.band.stop(now)
	g.running = false
	if g.rec != nil {
		g.rec.flush()
	}
}

func (g *game) Update() error {
	now := time.Now()
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return ebiten.Termination
	case inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.running:
		g.stop(now)
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		g.start(now)
	case inpututil.IsKeyJustPressed(ebiten.KeyT) && !g.running:
		g.practicing = !g.practicing
	}

	for drained := false; !drained; {
		select {
		case e := <-g.events:
			if e.Down {
				g.last = e
				g.mark(e)
			}
			if e.Key >= 0 && e.Key < len(g.down) {
				g.down[e.Key] = e.Down
			}
		default:
			drained = true
		}
	}

	var demo [128]bool
	if k := g.demoKey(now); k > 0 {
		demo[k] = true
	}
	g.piano.update(&g.down, &demo)

	if !g.running || g.practicing {
		return nil
	}
	// The demo's notes reach the marker when they sound, as the
	// player's would: a generated line must land every arrival.
	for len(g.heard) > 0 && !g.heard[0].At.After(now) {
		if m, ok := g.marker.Play(g.heard[0]); ok {
			g.record(m)
		}
		g.heard = g.heard[1:]
	}
	for _, bm := range g.marker.Close(now) {
		g.marks[bm.Beat] = bm.Kind
		g.walker.mark(bm.Kind, now)
	}
	_, end := g.m.Due(now, now.Add(lookahead))
	for ; g.next < end && g.next < len(g.beats); g.next++ {
		key, at := 0, g.m.At(g.next)
		if g.line != nil && g.next >= 0 {
			key = g.line[g.next]
			g.heard = append(g.heard, Note{Key: key, At: at})
		}
		g.band.beat(g.m.Position(g.next), key, at)
	}
	if now.After(g.m.At(len(g.beats))) {
		g.stop(now)
	}
	return nil
}

// demoKey is the key the reference bass sounds at `now`, 0 for none.
// It lights for the first four fifths of the beat only, so that a root
// repeated from one beat to the next pulses on the keyboard as it does
// in the ear, rather than staying lit.
func (g *game) demoKey(now time.Time) int {
	if g.line == nil || !g.running || g.practicing {
		return 0
	}
	x := g.m.Beats(now)
	n := int(math.Floor(x))
	if n < 0 || n >= len(g.beats) || x-float64(n) > 0.8 {
		return 0
	}
	return g.line[n]
}

// mark hands a key pressed in the bass zone to the marker, or to the
// practice, and keeps what it was for the status line.
func (g *game) mark(e keyboard.Event) {
	if !g.running || e.Key >= g.band.split {
		return
	}
	if g.practicing {
		chord, at := g.practice.waiting().Chord, g.practice.at
		g.lastMark = g.noteName(e.Key, chord) + " : " + pitchWords[g.practice.play(e.Key)]
		if g.practice.at != at {
			g.walker.mark(Landed, time.Now())
		}
		return
	}
	m, ok := g.marker.Play(Note{Key: e.Key, At: e.At})
	if !ok {
		return
	}
	g.lastMark = g.noteName(e.Key, g.beats[m.Beat].Chord) + " : " + markWords(m)
	g.record(m)
}

// markWords says what a note was: its timing, and its pitch when a beat
// claims it.
func markWords(m NoteMark) string {
	if m.Timing == Between {
		return timingWords[m.Timing]
	}
	return timingWords[m.Timing] + ", " + pitchWords[m.Pitch]
}

// record writes `m` down, when the game records: the bar counted within
// its chorus, the chord as the chart writes it.
func (g *game) record(m NoteMark) {
	if g.rec == nil {
		return
	}
	ch := g.beats[m.Beat].Chord
	p := g.m.Position(m.Beat)
	bars := g.chorusLen / perBar
	chorus := (p.Bar-1)/bars + 1
	p.Bar = (p.Bar-1)%bars + 1
	g.rec.note(chorus, p, m.Off, symbol(g.namer, ch), octaveName(g.noteName(m.Key, ch), m.Key), markWords(m))
}

// noteName spells `key` over the chord it was played on, as a lead
// sheet writes the degrees of a chord: the E♭ of F7, never a D♯. In
// ASCII: the status line is in Go Regular, which has no musical signs.
func (g *game) noteName(key int, ch analysis.Change) string {
	pc := harmony.PitchClass(key % 12)
	n := g.namer.Note(pc)
	if !ch.Silent {
		n = naming.SpellOver(g.namer.Note(ch.Chord.Root), ch.Chord.Pattern, pc)
	}
	return naming.English.Name(n, naming.ASCII)
}

var (
	timingWords = map[Timing]string{OnTime: "sur le temps", Early: "en avance", Late: "en retard", Between: "entre deux temps"}
	pitchWords  = map[Pitch]string{Root: "fondamentale", ChordTone: "note de l'accord", Outside: "hors de l'accord"}
)

// markAt returns the mark of beat `n` of the chorus playing, if any.
func (g *game) markAt(n int, chorus int) (BeatKind, bool) {
	if g.practicing {
		if g.practice != nil && g.practice.landed[n] {
			return Landed, true
		}
		return 0, false
	}
	k, ok := g.marks[chorus*g.chorusLen+n]
	return k, ok
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(paper)
	c := canvas{dst: screen, scale: g.scale}
	c.text(g.title, g.fonts.ui, chartX, 24, ink)
	tempo := fmt.Sprintf("%.0f à la noire", g.bpm)
	if g.practicing {
		tempo = "sans tempo : la grille attend la fondamentale"
	}
	c.text(tempo, g.fonts.ui, chartX, 40, faint)
	g.drawMode(c)

	pos := math.Inf(-1)
	switch {
	case g.running && g.practicing:
		pos = float64(g.practice.waiting().N)
	case g.running:
		pos = g.m.Beats(time.Now())
	}
	g.drawChart(c, pos)

	if !g.practicing && pos < 0 && pos > -countIn {
		// The count-in, in big: "1, 2, 3, 4".
		p := g.m.Position(int(math.Floor(pos)))
		c.centred(fmt.Sprint(p.Beat), g.fonts.count, screenWidth/2, chartY+rowH/2, ink)
	}

	g.piano.draw(c)
	g.drawWalker(c)

	status := "Espace : jouer ou arrêter   T : avec ou sans tempo   Échap : quitter"
	switch {
	case g.midi == "":
		status = "sans clavier MIDI   " + status
	case g.lastMark != "":
		status = g.lastMark + "   " + status
	}
	c.text(status, g.fonts.ui, chartX, screenHeight-28, faint)
}

// drawChart draws the twelve bars, the one being played shaded, and a
// cursor at `pos`, in beats; minus infinity when nothing plays.
func (g *game) drawChart(c canvas, pos float64) {
	bars := len(g.bars)
	playing, chorus := -1, 0
	if pos >= 0 {
		playing = int(pos) / perBar % bars
		chorus = int(pos) / g.chorusLen
	}
	for i, cells := range g.bars {
		x := float32(chartX + i%barsPerRow*barW)
		y := float32(chartY + i/barsPerRow*rowH)
		if i == playing {
			c.rect(x, y, barW, rowH-10, pale)
			frac := math.Mod(pos, perBar) / perBar
			cx := x + float32(frac)*barW
			c.line(cx, y, cx, y+rowH-10, 1.5, ink)
		}
		c.line(x, y, x, y+rowH-10, 1, ink)
		if len(cells) == 0 {
			c.centred("%", g.fonts.chord, float64(x)+barW/2, float64(y)+14, ink)
		}
		for _, cl := range cells {
			cx := float64(x) + 8 + float64(cl.beat)*barW/perBar
			c.text(cl.name, g.fonts.chord, cx, float64(y)+14, ink)
		}
		for beat := range perBar {
			if k, ok := g.markAt(i*perBar+beat, chorus); ok {
				drawMark(c, k, x+14+float32(beat)*barW/perBar, y+rowH-20)
			}
		}
		if i%barsPerRow == barsPerRow-1 || i == bars-1 {
			end := x + barW
			c.line(end, y, end, y+rowH-10, 1, ink)
		}
	}
}

// drawWalker draws the stick figure in the header, between the title
// and the mode: on the beats while the band plays, swaying on a slow
// clock of his own otherwise.
func (g *game) drawWalker(c canvas) {
	now := time.Now()
	beats := float64(now.UnixMilli()%1_000_000) / 1500 // a slow sway
	tempo := g.running && !g.practicing
	if tempo {
		beats = g.m.Beats(now)
		tempo = beats >= 0 // the count-in: still looking for the tempo
	}
	var col color.Color = ink
	if g.band.demo && !g.practicing {
		col = faint // a silhouette: the band plays, not a player
	}
	g.walker.draw(c, screenWidth/2, 60, g.walker.gait(tempo), beats, now, col)
}

// drawMode draws, top right, the phase the space bar starts or is
// playing: a framed label, filled while it plays, so that the mode
// reads before the first note.
func (g *game) drawMode(c canvas) {
	label := "A TEMPO" // Italian, as on a score
	switch {
	case g.practicing:
		label = "LIBRE"
	case g.band.demo:
		label = "A TEMPO, DÉMO"
	}
	if !g.running {
		label += " : Espace pour lancer, T pour changer"
	}
	w, h := c.measure(label, g.fonts.ui)
	const pad = 4
	x, y := float32(screenWidth-chartX-w-2*pad), float32(18)
	bw, bh := float32(w+2*pad), float32(h+2*pad)
	fill, text := paper, ink
	if g.running {
		fill, text = ink, paper
	}
	c.rect(x-1, y-1, bw+2, bh+2, ink)
	c.rect(x, y, bw, bh, fill)
	c.text(label, g.fonts.ui, float64(x)+pad, float64(y)+pad, text)
}

// drawMark draws the mark of an arrival centred on `x`, `y`: a green
// dot when landed, a red cross when missed, two orange strokes when
// doubled.
func drawMark(c canvas, k BeatKind, x, y float32) {
	switch k {
	case Landed:
		c.circle(x, y, 3.5, landedInk)
	case Missed:
		c.line(x-3.5, y-3.5, x+3.5, y+3.5, 1.5, missedInk)
		c.line(x-3.5, y+3.5, x+3.5, y-3.5, 1.5, missedInk)
	case Doubled:
		c.line(x-2, y-4, x-2, y+4, 1.5, doubledInk)
		c.line(x+2, y-4, x+2, y+4, 1.5, doubledInk)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	w, h := float64(outsideWidth)*s, float64(outsideHeight)*s
	g.scale = min(w/screenWidth, h/screenHeight)
	return int(screenWidth * g.scale), int(screenHeight * g.scale)
}

// namerFor spells in the tonality the analysis hears in `grid`: the
// tonality is found, never read from the key the chart declares.
func namerFor(grid analysis.Changes) *naming.Namer {
	namer, err := naming.NewNamer(naming.English)
	if err != nil {
		panic(err) // English is a locale of naming: it cannot fail
	}
	if tune := analysis.Hear(grid, nil).Tune; len(tune) > 0 {
		namer = namer.WithTonality(tune[0])
	}
	return namer
}

// chart writes the chords of `grid` bar by bar, spelled by `namer`. A
// bar where no chord starts is left empty, and drawn as the repeat
// sign.
func chart(grid analysis.Changes, namer *naming.Namer) [][]cell {
	bars := make([][]cell, len(grid.Bars))
	for _, ch := range grid.Chords {
		bar := grid.Bar(ch.Start)
		if bar < 0 {
			continue
		}
		beat := int((ch.Start - grid.Bars[bar]) / analysis.TicksPerBeat)
		bars[bar] = append(bars[bar], cell{name: symbol(namer, ch), beat: beat})
	}
	return bars
}

// symbol writes a change as a chart does: B♭7, Gm7, F6, D7/F♯, N.C.
func symbol(n *naming.Namer, ch analysis.Change) string {
	if ch.Silent {
		return "N.C."
	}
	q, ok := naming.ChordStyle{}.Symbol(ch.Chord.Pattern)
	if !ok {
		q = "?"
	}
	s := n.Name(ch.Chord.Root) + q
	if ch.Inverted() {
		s += "/" + n.Name(ch.Bass)
	}
	return s
}
