package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/games/walk/band"
	"github.com/ArnaudCalmettes/gohar/games/walk/bass"
	"github.com/ArnaudCalmettes/gohar/games/walk/chart"
	"github.com/ArnaudCalmettes/gohar/games/walk/figure"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// The count-in: two bars of hi-hat, "1, 3, 1, 2, 3, 4", the time to
// move the hands from the space bar to the keyboard.
const (
	perBar  = mark.BeatsPerBar
	countIn = band.CountIn
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

	// Within a bar: a mark from the start of its beat.
	markDX, markDY = 14, barH - 10

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

// The walker's bubble: how long it stays, and where, above his head.
const (
	bubbleHold = 1500 * time.Millisecond
	bubbleLift = 130 // from his feet to the line under the phrase
)

// game is the game scene of Walk With Me: the chart, the count-in, the
// band, the player's hands, and the marks of the first palier. Two
// phases: with the tempo, the band plays and the marker marks; without,
// the chart waits for the player (practice). At rest, the arrows choose
// the grid and the tempo; Escape goes back to the title.
type game struct {
	*app

	// The grid chosen, and the tempo: the options' at first, changed
	// here for the session.
	title string
	grid  analysis.Changes
	bpm   float64

	namer      *naming.Namer                     // for a note over no chord
	written    map[analysis.Ticks]chordpro.Chord // the chords as the grid writes them, by start
	bars       [][]chart.Cell                    // the chart, bar by bar; empty for a chord held from before
	formulas   map[analysis.Ticks]mark.Formula   // the formulas its chords make, for the review
	small      bool                              // a bar holds two chords: all are written small
	chorusLen  int                               // beats in a chorus
	rules      mark.Rules                        // what the marker and the practice expect
	practicing bool                              // the phase without tempo

	running  bool
	m        tempo.Metronome
	beats    []mark.Beat
	next     int // the next beat to schedule
	marker   *mark.Marker
	marks    map[int]mark.BeatKind // the arrivals marked, by beat
	practice *mark.Practice
	lastMark string // what the last note of the bass was
	last     keyboard.Event
	down     [screen.MIDIKeys]bool // the keys held, as the events tell
	piano    screen.Piano
	walker   figure.Walker
	line     []int       // the reference line, in demo: a key per beat
	heard    []mark.Note // its notes not yet sounded, for the marker
	ending   int         // the key the demo ends on, after the last beat; 0 for none

	summary  *mark.Summary // the marks of the run, for its review
	coach    *figure.Coach // what the walker says, the player's runs only
	bubble   string        // the phrase in his bubble, empty for none
	bubbleAt time.Time     // when he said it
}

func newGame(a *app) *game {
	rules := mark.FirstPalier
	rules.Split = a.band.Split // -split moves the marker's zone with the sound's
	g := &game{
		app:        a,
		bpm:        a.bpm,
		rules:      rules,
		practicing: a.untimed,
		piano:      newPiano(),
	}
	g.load(a.current)
	return g
}

// load makes the grid `i` of the game the one played, and remembers it
// for the session.
func (g *game) load(i int) {
	g.current = i
	t := g.tunes[i]
	g.title, g.grid = t.Title, t.Grid
	g.namer = namerFor(t.Grid)
	g.written = map[analysis.Ticks]chordpro.Chord{}
	for i, ch := range t.Grid.Chords {
		g.written[ch.Start] = t.Written[i]
	}
	g.bars = chart.Bars(t)
	g.formulas = mark.FormulasOf(t.Grid)
	g.small = chart.Small(g.bars)
	g.chorusLen = len(mark.Expect(t.Grid, tempo.NewMetronome(time.Time{}, g.bpm, perBar), 1))
	g.marks, g.practice = nil, nil
}

// Enter stops the music of the menus: the game counts in on its own,
// at its own tempo. The space bar starts it.
func (g *game) Enter() {
	if g.jam != nil {
		g.jam.Stop(time.Now())
		g.jam = nil
	}
	g.band.BassHand.Store(true) // the left hand walks: a double bass under the split
}

// Leave stops the run, if any: the bass released, the recording
// flushed. Elsewhere the keyboard is a piano from end to end: over the
// jam of the menus, in the lessons.
func (g *game) Leave() {
	if g.running {
		g.stop(time.Now())
	}
	g.band.BassHand.Store(false)
}

func (g *game) start(now time.Time) {
	g.running = true
	g.lastMark = ""
	g.walker = figure.Walker{} // each run starts from a walk
	seed := uint64(now.UnixNano())
	rng := rand.New(rand.NewPCG(seed, seed>>32|1))
	g.coach, g.bubble = figure.NewCoach(rng, coachLines), ""
	if g.practicing {
		// No time: any metronome numbers the beats.
		g.practice = mark.NewPractice(g.rules, mark.Expect(g.grid, tempo.NewMetronome(now, g.bpm, perBar), 1))
		return
	}
	beat := time.Duration(float64(time.Minute) / g.bpm)
	g.m = tempo.NewMetronome(now.Add(band.Lookahead+countIn*beat), g.bpm, perBar)
	g.beats = mark.Expect(g.grid, g.m, g.choruses)
	g.next = -countIn
	g.marker = mark.NewMarker(g.rules, g.m, g.beats)
	g.marks = map[int]mark.BeatKind{}
	g.summary = mark.NewSummary(len(g.bars), g.chorusLen, g.formulas)
	g.line, g.heard, g.ending = nil, nil, 0
	if g.rec != nil {
		who := g.lang.T(msgPlayer)
		if g.band.Demo {
			who = g.lang.T(msgDemo)
		}
		g.rec.run(now, g.bpm, who)
	}
	if g.band.Demo {
		g.line = bass.Walk(g.beats, rng)
		g.ending = bass.Ending(g.beats, g.line, g.grid.End)
	}
}

func (g *game) stop(now time.Time) {
	g.band.Stop(now)
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
	case inpututil.IsKeyJustPressed(ebiten.KeyD) && !g.running && !g.practicing:
		g.band.Demo = !g.band.Demo
	case g.running:
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp):
		g.load((g.current + len(g.tunes) - 1) % len(g.tunes))
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown):
		g.load((g.current + 1) % len(g.tunes))
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
		g.bpm = max(g.bpm-bpmStep, minBPM)
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight):
		g.bpm = min(g.bpm+bpmStep, maxBPM)
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
		g.summary.Beat(bm.Beat, g.beats[bm.Beat], bm.Kind)
		g.walker.Mark(bm.Kind, now)
		if !g.band.Demo {
			g.speak(g.coach.Arrival(bm.Beat, bm.Kind == mark.Landed, g.walker.Gait(true) == figure.Snapping), now)
		}
	}
	_, end := g.m.Due(now, now.Add(band.Lookahead))
	for ; g.next < end && g.next < len(g.beats); g.next++ {
		key, at := 0, g.m.At(g.next)
		if g.line != nil && g.next >= 0 {
			key = g.line[g.next]
			g.heard = append(g.heard, mark.Note{Key: key, At: at})
		}
		c := band.Cue{Pos: g.m.Position(g.next), Key: key, Snap: g.walker.Gait(true) == figure.Snapping}
		a := band.Full
		if g.next < 0 {
			a = band.HatCountIn
		}
		g.band.Play(a.Strokes(c), g.next, g.m)
	}
	// The demo's last note, on the beat after the last one, held a bar:
	// the run stops once it has rung.
	last := len(g.beats)
	if g.ending != 0 {
		if g.next == last && end > last {
			g.band.Play(band.Ending.Strokes(band.Cue{Key: g.ending}), last, g.m)
			g.next++
		}
		last += perBar
	}
	if now.After(g.m.At(last)) { // played through: the player's run gets its review
		g.stop(now)
		if !g.band.Demo {
			return scene.Push(newReview(g))
		}
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
	if !g.running || e.Key >= g.band.Split {
		return
	}
	if g.practicing {
		chord, at := g.practice.Waiting().Chord, g.practice.Arrival()
		pitch := g.practice.Play(e.Key)
		g.lastMark = g.lang.T(msgMark, "Note", g.noteName(e.Key, chord), "Words", g.lang.T(pitchPhrase[pitch]))
		if g.practice.Arrival() != at {
			g.walker.Mark(mark.Landed, time.Now())
		}
		return
	}
	// The calibrated latency, taken off: the instant the player heard
	// the beat he answered, rather than the one the game received.
	m, ok := g.marker.Play(mark.Note{Key: e.Key, At: e.At.Add(-g.latency)})
	if !ok {
		return
	}
	g.lastMark = g.lang.T(msgMark, "Note", g.noteName(e.Key, g.beats[m.Beat].Chord), "Words", g.markWords(m))
	g.record(m)
	g.summary.Note(m)
	if !g.band.Demo {
		g.speak(g.coach.Note(m.Off, m.Timing, m.Beat), time.Now())
	}
}

// speak puts the phrase `id` in the walker's bubble, at `now`; an empty
// one leaves the bubble as it is.
func (g *game) speak(id string, now time.Time) {
	if id != "" {
		g.bubble, g.bubbleAt = g.lang.T(id), now
	}
}

// markWords says what a note was: its timing, and its pitch when a beat
// claims it.
func (g *game) markWords(m mark.NoteMark) string {
	timing := g.lang.T(timingPhrase[m.Timing])
	if m.Timing == mark.Between {
		return timing
	}
	return timing + ", " + g.lang.T(pitchPhrase[m.Pitch])
}

// record writes `m` down, when the game records: the bar counted within
// its chorus, the chord as the chart writes it.
func (g *game) record(m mark.NoteMark) {
	if g.rec == nil {
		return
	}
	ch := g.beats[m.Beat].Chord
	p := g.m.Position(m.Beat)
	bar, chorus := mark.BarOf(m.Beat, g.chorusLen)
	p.Bar = bar + 1
	g.rec.note(chorus+1, p, m.Off, chart.SymbolOf(g.written[ch.Start]).String(), octaveName(g.noteName(m.Key, ch), m.Key), g.markWords(m))
}

// noteName spells `key` over the chord it was played on, as a lead
// sheet writes the degrees of a chord, its root as the grid writes it:
// the E♭ of F7, never a D♯. In ASCII: the status line is in Go Regular,
// which has no musical signs.
func (g *game) noteName(key int, ch analysis.Change) string {
	pc := harmony.PitchClass(key % octave)
	n := g.namer.Note(pc)
	if !ch.Silent {
		root := g.namer.Note(ch.Chord.Root)
		if w, ok := g.written[ch.Start]; ok && !w.NoChord {
			root = w.Root
		}
		n = naming.SpellOver(root, ch.Chord.Pattern, pc)
	}
	return naming.English.Name(n, naming.ASCII)
}

// markAt returns the mark of beat `n` of the chorus playing, if any.
func (g *game) markAt(n int, chorus int) (mark.BeatKind, bool) {
	if g.practicing {
		if g.practice != nil && g.practice.Landed(n) {
			return mark.Landed, true
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
		pos = float64(g.practice.Waiting().N)
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
	drawSplit(c, &g.piano, g.band.Split)
	g.drawWalker(c)
	g.drawBubble(c)

	status := g.keys()
	switch {
	case g.midi == "":
		status = g.lang.T(msgNoMIDI) + "   " + status
	case g.lastMark != "":
		status = g.lastMark + "   " + status
	}
	c.Text(status, g.fonts.ui, margin, statusY, faint)
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
	if g.band.Demo && !g.practicing {
		col = faint // a silhouette: the band plays, not a player
	}
	g.walker.Draw(c, walkerX, walkerY, walkerScale, g.walker.Gait(onBeat), beats, now, col)
}

// drawBubble draws what the walker says above his head, for a while: one
// line, centred over him (see figure.Bubble).
func (g *game) drawBubble(c screen.Canvas) {
	if g.bubble == "" || time.Since(g.bubbleAt) > bubbleHold {
		return
	}
	w, h := c.Measure(g.bubble, g.fonts.bubble)
	x := max(float64(margin), walkerX-w/2) // never off the screen
	under := float64(walkerY - bubbleLift)
	b := figure.Bubble{Font: g.fonts.bubble, Ink: ink, Left: x, Top: under - figure.BubbleGap - h, Width: w, LineH: h}
	b.Draw(c, []string{g.bubble}, walkerX, walkerY, walkerScale)
}

// keys says the keys of the game as they stand, each with what it does
// now: Space starts or stops; at rest, T and D name the phase and the
// demo they switch to, and the arrows choose the grid and the tempo. D
// only means something with the tempo: without, the chart waits for
// the player.
func (g *game) keys() string {
	if g.running {
		return g.lang.T(msgKeyStop) + "   " + g.lang.T(msgKeyMenu)
	}
	keys := []string{g.lang.T(msgKeyPlay)}
	if g.practicing {
		keys = append(keys, g.lang.T(msgKeyTempo))
	} else {
		keys = append(keys, g.lang.T(msgKeyFree))
		demo := msgKeyDemoOn
		if g.band.Demo {
			demo = msgKeyDemoOff
		}
		keys = append(keys, g.lang.T(demo))
	}
	keys = append(keys, g.lang.T(msgKeyGrid), g.lang.T(msgKeyMenu))
	return strings.Join(keys, "   ")
}

// drawMode draws, top right, the phase the space bar starts or is
// playing: a framed label, filled while it plays, so that the mode
// reads before the first note.
func (g *game) drawMode(c screen.Canvas) {
	label := g.lang.T(msgModeTempo) // Italian, as on a score, in any language
	switch {
	case g.practicing:
		label = g.lang.T(msgModeFree)
	case g.band.Demo:
		label = g.lang.T(msgModeDemo)
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

// namerFor spells in the tonality the analysis hears in `grid`: the
// tonality is found, never read from the key the chart declares. It
// only spells a note over no chord: the chords keep the spelling of the
// grid, which follows its modulations (see chart).
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
