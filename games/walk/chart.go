package main

import (
	"image/color"
	"math"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The chart on screen: the bars of the grid as a Real Book page lays
// them out, four a row (see "L'affichage" in docs/walk.md). The game
// draws three rows of it while it plays, the review all of them.

// chart writes the chords of `t` bar by bar, as the grid spells them:
// one tonality is not enough to spell a grid that modulates, the C of
// Tune Up's Cmaj7 would be a B♯ in D. A bar where no chord starts is
// left empty, and drawn as the repeat sign.
//
// Respelling by the tonalities the analysis hears, zone by zone, as
// charts/cmd/analyse does, is for later (see docs/chantiers.md).
func chart(t tune) [][]cell {
	bars := make([][]cell, len(t.grid.Bars))
	for i, ch := range t.grid.Chords {
		bar := t.grid.Bar(ch.Start)
		if bar < 0 {
			continue
		}
		beat := int((ch.Start - t.grid.Bars[bar]) / analysis.TicksPerBeat)
		bars[bar] = append(bars[bar], cell{name: symbolOf(t.written[i]), beat: beat})
	}
	return bars
}

// drawChart draws the bars, the one being played shaded, and a cursor
// at `pos`, in beats; minus infinity when nothing plays.
//
// Three rows show at a time. A longer grid turns its pages a row at a
// time, as a reader of a Real Book does: the row being played second,
// the one before it above, the next one below, to read ahead.
func (g *game) drawChart(c screen.Canvas, pos float64) {
	bars := len(g.bars)
	playing, chorus := -1, 0
	if pos >= 0 {
		playing, chorus = barOf(int(pos), g.chorusLen)
	}
	rows := (bars + barsPerRow - 1) / barsPerRow
	top := 0
	if rows > chartRows && playing >= 0 {
		top = min(max(playing/barsPerRow-1, 0), rows-chartRows)
	}
	for i, cells := range g.bars {
		row := i / barsPerRow
		if row < top || row >= top+chartRows {
			continue
		}
		x := float32(chartX + i%barsPerRow*barW)
		y := float32(chartY + (row-top)*rowH)
		var fill color.Color
		if i == playing {
			fill = pale
		}
		g.drawBar(c, cells, x, y, barH, fill, false, i%barsPerRow == barsPerRow-1 || i == bars-1)
		if i == playing {
			cx := x + float32(math.Mod(pos, perBar)/perBar)*barW
			c.Line(cx, y, cx, y+barH, 1.5, ink)
		}
		for beat := range perBar {
			if k, ok := g.markAt(i*perBar+beat, chorus); ok {
				drawMark(c, k, x+markDX+float32(beat)*barW/perBar, y+markDY)
			}
		}
	}
}

// drawBar draws the bar `cells` from (`x`, `y`), `h` high: its `fill`,
// if not nil, its barline, the closing one too when it ends its row
// (`last`), and its chords, or the repeat sign for a bar where none
// starts. Two chords sharing the bar are written small, and all of them
// when the rows are shrunk (`small`).
func (g *game) drawBar(c screen.Canvas, cells []cell, x, y, h float32, fill color.Color, small, last bool) {
	if fill != nil {
		c.Rect(x, y, barW, h, fill)
	}
	c.Line(x, y, x, y+h, 1, ink)
	if last {
		c.Line(x+barW, y, x+barW, y+h, 1, ink)
	}
	dy := float64(y) + min(chordDY, float64(h)/2-10) // about the middle of a shrunk bar
	if len(cells) == 0 {
		c.Centred("%", g.fonts.chord, float64(x)+barW/2, dy, ink)
	}
	font, raised := g.fonts.chord, g.fonts.chordRaised
	if small || len(cells) > 1 {
		font, raised = g.fonts.chordSmall, g.fonts.chordSmallRaised
	}
	for _, cl := range cells {
		cx := float64(x) + chordDX + float64(cl.beat)*barW/perBar
		drawSymbol(c, cl.name, font, raised, cx, dy)
	}
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

// drawSymbol draws the symbol `s` from (`x`, `y`): its line in `font`,
// its raised run in the smaller `raised`, a little above, as an
// exponent.
func drawSymbol(c screen.Canvas, s symbol, font, raised *screen.Font, x, y float64) {
	c.Text(s.line, font, x, y, ink)
	w, _ := c.Measure(s.line, font)
	x += w
	if s.raised != "" {
		c.Text(s.raised, raised, x, y-raisedDY, ink)
		w, _ = c.Measure(s.raised, raised)
		x += w
	}
	c.Text(s.bass, font, x, y, ink)
}
