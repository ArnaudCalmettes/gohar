package main

import (
	"image/color"
	"math"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/chart"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// pen is what the chart is written with, in the game and the review.
func (a *app) pen() chart.Pen {
	return chart.Pen{
		Chord: a.fonts.chord, ChordRaised: a.fonts.chordRaised,
		Small: a.fonts.chordSmall, SmallRaised: a.fonts.chordSmallRaised,
		Ink: ink,
	}
}

// drawChart draws the page of the chart being played: the bars, the
// one being played shaded, and a cursor at `pos`, in beats; minus
// infinity when nothing plays.
//
// Three rows show at a time. A longer grid turns its pages a row at a
// time, as a reader of a Real Book does: the row being played second,
// the one before it above, the next one below, to read ahead.
func (g *game) drawChart(c screen.Canvas, pos float64) {
	bars := len(g.bars)
	playing, chorus := -1, 0
	if pos >= 0 {
		playing, chorus = mark.BarOf(int(pos), g.chorusLen)
	}
	rows := (bars + barsPerRow - 1) / barsPerRow
	top := 0
	if rows > chartRows && playing >= 0 {
		top = min(max(playing/barsPerRow-1, 0), rows-chartRows)
	}
	pen := g.pen()
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
		pen.DrawBar(c, cells, x, y, barW, barH, fill, g.small, i%barsPerRow == barsPerRow-1 || i == bars-1)
		if i == playing {
			cx := x + float32(math.Mod(pos, perBar)/perBar)*barW
			c.Line(cx, y, cx, y+barH, 1.5, ink)
		}
		for beat := range perBar {
			if k, ok := g.markAt(i*perBar+beat, chorus); ok {
				chart.DrawMark(c, k, x+markDX+float32(beat)*barW/perBar, y+markDY)
			}
		}
	}
}
