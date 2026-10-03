package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
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

// game is the shell of Walk With Me, first delivery: the chart, the
// count-in, the band and the player's hands. No marks yet.
type game struct {
	title    string
	bpm      float64
	choruses int
	grid     analysis.Changes
	bars     [][]cell // the chart, bar by bar; empty for a chord held from before
	band     *band
	events   chan keyboard.Event
	midi     string // the keyboard's name, empty without one

	running bool
	m       Metronome
	beats   []Beat
	next    int // the next beat to schedule
	last    keyboard.Event
	down    [128]bool // the keys held, as the events tell
	piano   piano

	fonts fonts
	scale float64
}

func newGame(title string, grid analysis.Changes, bpm float64, choruses int, bd *band, midi string) (*game, error) {
	fs, err := newFonts()
	if err != nil {
		return nil, err
	}
	return &game{
		title:    title,
		bpm:      bpm,
		choruses: choruses,
		grid:     grid,
		bars:     chart(grid),
		band:     bd,
		events:   make(chan keyboard.Event, 64),
		midi:     midi,
		piano:    piano{split: bd.split},
		fonts:    fs,
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
	beat := time.Duration(float64(time.Minute) / g.bpm)
	g.m = NewMetronome(now.Add(lookahead+countIn*beat), g.bpm, perBar)
	g.beats = Expect(g.grid, g.m, g.choruses)
	g.next = -countIn
	g.running = true
}

func (g *game) stop(now time.Time) {
	g.band.stop(now)
	g.running = false
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
	}

	for drained := false; !drained; {
		select {
		case e := <-g.events:
			if e.Down {
				g.last = e
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

	if !g.running {
		return nil
	}
	_, end := g.m.Due(now, now.Add(lookahead))
	for ; g.next < end && g.next < len(g.beats); g.next++ {
		var b *Beat
		if g.next >= 0 {
			b = &g.beats[g.next]
		}
		g.band.beat(g.m.Position(g.next), b, g.m.At(g.next))
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
	if !g.band.demo || !g.running {
		return 0
	}
	x := g.m.Beats(now)
	n := int(math.Floor(x))
	if n < 0 || n >= len(g.beats) || x-float64(n) > 0.8 {
		return 0
	}
	return bassKey(g.beats[n].Chord.Bass)
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(paper)
	c := canvas{dst: screen, scale: g.scale}
	c.text(g.title, g.fonts.ui, chartX, 24, ink)
	c.text(fmt.Sprintf("%.0f à la noire", g.bpm), g.fonts.ui, chartX, 40, faint)

	pos := math.Inf(-1)
	if g.running {
		pos = g.m.Beats(time.Now())
	}
	g.drawChart(c, pos)

	if pos < 0 && pos > -countIn {
		// The count-in, in big: "1, 2, 3, 4".
		p := g.m.Position(int(math.Floor(pos)))
		c.centred(fmt.Sprint(p.Beat), g.fonts.count, screenWidth/2, chartY+rowH/2, ink)
	}

	g.piano.draw(c)

	status := "Espace : jouer ou arrêter   Échap : quitter"
	if g.midi == "" {
		status = "sans clavier MIDI   " + status
	} else if g.last.Down {
		status = fmt.Sprintf("%s, dernière touche %d   %s", g.midi, g.last.Key, status)
	}
	c.text(status, g.fonts.ui, chartX, screenHeight-28, faint)
}

// drawChart draws the twelve bars, the one being played shaded, and a
// cursor at `pos`, in beats; minus infinity when nothing plays.
func (g *game) drawChart(c canvas, pos float64) {
	bars := len(g.bars)
	playing := -1
	if pos >= 0 {
		playing = int(pos) / perBar % bars
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
		if i%barsPerRow == barsPerRow-1 || i == bars-1 {
			end := x + barW
			c.line(end, y, end, y+rowH-10, 1, ink)
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	w, h := float64(outsideWidth)*s, float64(outsideHeight)*s
	g.scale = min(w/screenWidth, h/screenHeight)
	return int(screenWidth * g.scale), int(screenHeight * g.scale)
}

// chart writes the chords of `grid` bar by bar, spelled in the
// tonality the analysis hears: the tonality is found, never read from
// the key the chart declares. A bar where no chord starts is left
// empty, and drawn as the repeat sign.
func chart(grid analysis.Changes) [][]cell {
	namer, err := naming.NewNamer(naming.English)
	if err != nil {
		panic(err) // English is a locale of naming: it cannot fail
	}
	if tune := analysis.Hear(grid, nil).Tune; len(tune) > 0 {
		namer = namer.WithTonality(tune[0])
	}

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
