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
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// lookahead is how far ahead the band is scheduled. Ebiten's loop
// ticks every 16 ms, too coarse to play a beat on time: notes are
// queued a window ahead and the synth applies them to the sample (see
// "Le son" in docs/walk.md).
const lookahead = 100 * time.Millisecond

// swingRatio places the ride's "and", the triplet swing of the method
// books (see tempo.Swing).
const swingRatio = 2.0 / 3

// The count-in: two bars of hi-hat, "1, 3, 1, 2, 3, 4", the time to
// move the hands from the space bar to the keyboard.
const (
	perBar  = 4
	countIn = 2 * perBar
)

// The layout: the chart on the right two thirds of the screen, laid
// out as a Real Book page, four bars a row; the walker on the third
// left of it, standing level with its last row; the keyboard below
// (see piano.go).
const (
	margin  = 20 // around the screen
	titleY  = 24 // the title, then the tempo, in the margin
	tempoY  = 40
	modeY   = 18 // the mode, top right
	statusY = screenHeight - 28

	barsPerRow = 4
	chartRows  = 3 // the twelve bars of the blues
	barW       = 105
	rowH       = 56 // room for the keyboard below
	barH       = rowH - 10
	chartX     = screenWidth - margin - barsPerRow*barW
	chartY     = 70

	// Within a bar: a chord and a mark from the start of their beat.
	chordDX, chordDY = 8, 14
	markDX, markDY   = 14, barH - 10

	walkerX     = chartX / 2
	walkerY     = chartY + (chartRows-1)*rowH + barH
	walkerScale = 3
)

const (
	// demoHold is the part of a beat the demo's key stays lit: a root
	// repeated from one beat to the next pulses, as in the ear.
	demoHold = 0.8

	// swayPeriod is the walker's own beat at rest, in milliseconds.
	swayPeriod = 1500
)

// A cell is one chord written in a bar, at the beat it starts on.
type cell struct {
	name string
	beat int // in the bar, from 0
}

// game is the game scene of Walk With Me: the chart, the count-in, the
// band, the player's hands, and the marks of the first palier. Two
// phases: with the tempo, the band plays and the marker marks; without,
// the chart waits for the player (practice). Escape goes back to the
// title.
type game struct {
	*app

	namer      *naming.Namer
	bars       [][]cell // the chart, bar by bar; empty for a chord held from before
	chorusLen  int      // beats in a chorus
	rules      Rules    // what the marker and the practice expect
	practicing bool     // the phase without tempo

	running  bool
	m        tempo.Metronome
	swing    tempo.Swing // the ride's eighths
	beats    []Beat
	next     int // the next beat to schedule
	marker   *Marker
	marks    map[int]BeatKind // the arrivals marked, by beat
	practice *practice
	lastMark string // what the last note of the bass was
	last     keyboard.Event
	down     [screen.MIDIKeys]bool // the keys held, as the events tell
	piano    screen.Piano
	walker   walker
	line     []int  // the reference line, in demo: a key per beat
	heard    []Note // its notes not yet sounded, for the marker
	ending   int    // the key the demo ends on, after the last beat; 0 for none
}

func newGame(a *app) *game {
	namer := namerFor(a.grid)
	rules := FirstPalier
	rules.Split = a.band.split // -split moves the marker's zone with the sound's
	return &game{
		app:        a,
		namer:      namer,
		bars:       chart(a.grid, namer),
		chorusLen:  len(Expect(a.grid, tempo.NewMetronome(time.Time{}, a.bpm, perBar), 1)),
		rules:      rules,
		practicing: a.untimed,
		piano:      newPiano(),
	}
}

// Enter has nothing to start: the space bar does.
func (g *game) Enter() {}

// Leave stops the run, if any: the bass released, the recording
// flushed.
func (g *game) Leave() {
	if g.running {
		g.stop(time.Now())
	}
}

func (g *game) start(now time.Time) {
	g.running = true
	g.lastMark = ""
	g.walker = walker{} // each run starts from a walk
	if g.practicing {
		// No time: any metronome numbers the beats.
		g.practice = newPractice(g.rules, Expect(g.grid, tempo.NewMetronome(now, g.bpm, perBar), 1))
		return
	}
	beat := time.Duration(float64(time.Minute) / g.bpm)
	g.m = tempo.NewMetronome(now.Add(lookahead+countIn*beat), g.bpm, perBar)
	g.swing = tempo.NewSwing(g.m, swingRatio)
	g.beats = Expect(g.grid, g.m, g.choruses)
	g.next = -countIn
	g.marker = NewMarker(g.rules, g.m, g.beats)
	g.marks = map[int]BeatKind{}
	g.line, g.heard, g.ending = nil, nil, 0
	if g.rec != nil {
		who := g.lang.T(msgPlayer)
		if g.band.demo {
			who = g.lang.T(msgDemo)
		}
		g.rec.run(now, g.bpm, who)
	}
	if g.band.demo {
		seed := uint64(now.UnixNano())
		g.line = Walk(g.beats, rand.New(rand.NewPCG(seed, seed>>32|1)))
		g.ending = Ending(g.beats, g.line, g.grid.End)
	}
}

func (g *game) stop(now time.Time) {
	g.band.stop(now)
	g.running = false
	if g.rec != nil {
		g.rec.flush()
	}
}

func (g *game) Update() scene.Transition {
	now := time.Now()
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return scene.Replace(newTitle(g.app))
	case inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.running:
		g.stop(now)
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		g.start(now)
	case inpututil.IsKeyJustPressed(ebiten.KeyT) && !g.running:
		g.practicing = !g.practicing
	}

	g.drain(func(e keyboard.Event) {
		if e.Down {
			g.last = e
			g.mark(e)
		}
		if e.Key >= 0 && e.Key < len(g.down) {
			g.down[e.Key] = e.Down
		}
	})

	var demo [screen.MIDIKeys]bool
	if k := g.demoKey(now); k > 0 {
		demo[k] = true
	}
	g.piano.Update(screen.Lit{Down: &g.down, Color: playerLit}, screen.Lit{Down: &demo, Color: demoLit})

	if !g.running || g.practicing {
		return scene.Stay
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
		snap := g.walker.gait(true) == snapping
		g.band.beat(g.m.Position(g.next), key, at, g.swing.AtBeats(float64(g.next)+0.5), snap)
	}
	// The demo's last note, on the beat after the last one, held a bar:
	// the run stops once it has rung.
	last := len(g.beats)
	if g.ending != 0 {
		if g.next == last && end > last {
			g.band.end(g.ending, g.m.At(last))
			g.next++
		}
		last += perBar
	}
	if now.After(g.m.At(last)) {
		g.stop(now)
	}
	return scene.Stay
}

// demoKey is the key the reference bass sounds at `now`, 0 for none,
// lit for the first `demoHold` of the beat only.
func (g *game) demoKey(now time.Time) int {
	if g.line == nil || !g.running || g.practicing {
		return 0
	}
	x := g.m.Beats(now)
	n := int(math.Floor(x))
	switch {
	case n >= len(g.beats) && g.ending != 0:
		return g.ending // held to the end
	case n < 0 || n >= len(g.beats) || x-float64(n) > demoHold:
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
		pitch := g.practice.play(e.Key)
		g.lastMark = g.lang.T(msgMark, "Note", g.noteName(e.Key, chord), "Words", g.lang.T(pitchPhrase[pitch]))
		if g.practice.at != at {
			g.walker.mark(Landed, time.Now())
		}
		return
	}
	m, ok := g.marker.Play(Note{Key: e.Key, At: e.At})
	if !ok {
		return
	}
	g.lastMark = g.lang.T(msgMark, "Note", g.noteName(e.Key, g.beats[m.Beat].Chord), "Words", g.markWords(m))
	g.record(m)
}

// markWords says what a note was: its timing, and its pitch when a beat
// claims it.
func (g *game) markWords(m NoteMark) string {
	timing := g.lang.T(timingPhrase[m.Timing])
	if m.Timing == Between {
		return timing
	}
	return timing + ", " + g.lang.T(pitchPhrase[m.Pitch])
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
	g.rec.note(chorus, p, m.Off, symbol(g.namer, ch), octaveName(g.noteName(m.Key, ch), m.Key), g.markWords(m))
}

// noteName spells `key` over the chord it was played on, as a lead
// sheet writes the degrees of a chord: the E♭ of F7, never a D♯. In
// ASCII: the status line is in Go Regular, which has no musical signs.
func (g *game) noteName(key int, ch analysis.Change) string {
	pc := harmony.PitchClass(key % octave)
	n := g.namer.Note(pc)
	if !ch.Silent {
		n = naming.SpellOver(g.namer.Note(ch.Chord.Root), ch.Chord.Pattern, pc)
	}
	return naming.English.Name(n, naming.ASCII)
}

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

func (g *game) Draw(dst *ebiten.Image) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: g.scale}
	c.Text(g.title, g.fonts.ui, margin, titleY, ink)
	pace := g.lang.T(msgTempo, "BPM", fmt.Sprintf("%.0f", g.bpm))
	if g.practicing {
		pace = g.lang.T(msgFreeTempo)
	}
	c.Text(pace, g.fonts.ui, margin, tempoY, faint)
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
		// The count-in, in big: every beat of both bars, "1, 2, 3, 4",
		// even where the hi-hat only plays 1 and 3.
		p := g.m.Position(int(math.Floor(pos)))
		c.Centred(fmt.Sprint(p.Beat), g.fonts.count, chartX+barsPerRow*barW/2, chartY+rowH/2, ink)
	}

	g.piano.Draw(c, nil)
	drawSplit(c, &g.piano, g.band.split)
	g.drawWalker(c)

	status := g.lang.T(msgKeys)
	switch {
	case g.midi == "":
		status = g.lang.T(msgNoMIDI) + "   " + status
	case g.lastMark != "":
		status = g.lastMark + "   " + status
	}
	c.Text(status, g.fonts.ui, margin, statusY, faint)
}

// drawChart draws the twelve bars, the one being played shaded, and a
// cursor at `pos`, in beats; minus infinity when nothing plays.
func (g *game) drawChart(c screen.Canvas, pos float64) {
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
			c.Rect(x, y, barW, barH, pale)
			frac := math.Mod(pos, perBar) / perBar
			cx := x + float32(frac)*barW
			c.Line(cx, y, cx, y+barH, 1.5, ink)
		}
		c.Line(x, y, x, y+barH, 1, ink)
		if len(cells) == 0 {
			c.Centred("%", g.fonts.chord, float64(x)+barW/2, float64(y)+chordDY, ink)
		}
		for _, cl := range cells {
			cx := float64(x) + chordDX + float64(cl.beat)*barW/perBar
			c.Text(cl.name, g.fonts.chord, cx, float64(y)+chordDY, ink)
		}
		for beat := range perBar {
			if k, ok := g.markAt(i*perBar+beat, chorus); ok {
				drawMark(c, k, x+markDX+float32(beat)*barW/perBar, y+markDY)
			}
		}
		if i%barsPerRow == barsPerRow-1 || i == bars-1 {
			end := x + barW
			c.Line(end, y, end, y+barH, 1, ink)
		}
	}
}

// drawWalker draws the stick figure left of the chart: on the beats
// while the band plays, swaying on a slow clock of his own otherwise.
func (g *game) drawWalker(c screen.Canvas) {
	now := time.Now()
	beats := float64(now.UnixMilli()%1_000_000) / swayPeriod
	onBeat := g.running && !g.practicing
	if onBeat {
		beats = g.m.Beats(now)
		// Before the first beat, the count-in, he looks for the tempo;
		// after the last one, on the demo's last note, he stops.
		onBeat = beats >= 0 && beats < float64(len(g.beats))
	}
	var col color.Color = ink
	if g.band.demo && !g.practicing {
		col = faint // a silhouette: the band plays, not a player
	}
	g.walker.draw(c, walkerX, walkerY, walkerScale, g.walker.gait(onBeat), beats, now, col)
}

// drawMode draws, top right, the phase the space bar starts or is
// playing: a framed label, filled while it plays, so that the mode
// reads before the first note.
func (g *game) drawMode(c screen.Canvas) {
	label := g.lang.T(msgModeTempo) // Italian, as on a score, in any language
	switch {
	case g.practicing:
		label = g.lang.T(msgModeFree)
	case g.band.demo:
		label = g.lang.T(msgModeDemo)
	}
	if !g.running {
		label = g.lang.T(msgModeHint, "Mode", label)
	}
	w, h := c.Measure(label, g.fonts.ui)
	const pad = 4
	x, y := float32(screenWidth-margin-w-2*pad), float32(modeY)
	bw, bh := float32(w+2*pad), float32(h+2*pad)
	fill, text := paper, ink
	if g.running {
		fill, text = ink, paper
	}
	c.Rect(x-1, y-1, bw+2, bh+2, ink)
	c.Rect(x, y, bw, bh, fill)
	c.Text(label, g.fonts.ui, float64(x)+pad, float64(y)+pad, text)
}

// drawMark draws the mark of an arrival centred on `x`, `y`: a green
// dot when landed, a red cross when missed, two orange strokes when
// doubled.
func drawMark(c screen.Canvas, k BeatKind, x, y float32) {
	const (
		r      = 3.5 // the dot, the arms of the cross
		apart  = 2   // the two strokes, from the centre
		half   = 4   // their half length
		stroke = 1.5
	)
	switch k {
	case Landed:
		c.Circle(x, y, r, landedInk)
	case Missed:
		c.Line(x-r, y-r, x+r, y+r, stroke, missedInk)
		c.Line(x-r, y+r, x+r, y-r, stroke, missedInk)
	case Doubled:
		c.Line(x-apart, y-half, x-apart, y+half, stroke, doubledInk)
		c.Line(x+apart, y-half, x+apart, y+half, stroke, doubledInk)
	}
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
