package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/chart"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// pen is what the chart is written with, in the game and the review.
func (a *app) pen() chart.Pen {
	return chart.Pen{
		Chord: a.fonts.chord, ChordRaised: a.fonts.chordRaised,
		Small: a.fonts.chordSmall, SmallRaised: a.fonts.chordSmallRaised,
		Ink: ink, Paper: paper,
	}
}

// The scrolling of the chart: how long it glides from one row to the
// next, and the room left of the bars kept in the page, for the
// rehearsal marks.
const (
	scrollTime = 350 * time.Millisecond
	markRoom   = 40
)

// A scroll is the row at the top of the page, gliding: from `from` to
// `to`, from `at` on, over scrollTime, slowing down as it arrives.
type scroll struct {
	from, to float64
	at       time.Time
}

// row returns the row at the top of the page at `now`, fractional while
// it glides.
func (s scroll) row(now time.Time) float64 {
	t := min(float64(now.Sub(s.at))/float64(scrollTime), 1)
	t = 1 - (1-t)*(1-t)*(1-t) // ease out
	return s.from + (s.to-s.from)*t
}

// aim makes `to` the row the page glides to, from where it is at `now`.
func (s *scroll) aim(to float64, now time.Time) {
	if to != s.to {
		s.from, s.to, s.at = s.row(now), to, now
	}
}

// jump puts the page on row `to` at once: a new grid.
func (s *scroll) jump(to float64) { *s = scroll{from: to, to: to} }

// drawChart draws the page of the chart being played: the bars, the
// one being played shaded, and a cursor at `pos`, in beats; minus
// infinity when nothing plays.
//
// Three rows show at a time. A longer grid turns its pages a row at a
// time, as a reader of a Real Book does: the row being played second,
// the one before it above, the next one below, to read ahead. The page
// glides from one row to the next, cut at its edges; at rest, it stays
// where the player scrolled it (see scrollRest). A chevron above or
// below says rows are hidden there. The
// rehearsal marks of the grid, as written, stand left of the bars they
// open.
func (g *game) drawChart(c screen.Canvas, pos float64) {
	now := time.Now()
	bars := len(g.bars)
	playing, chorus := -1, 0
	if pos >= 0 {
		playing, chorus = mark.BarOf(int(pos), g.chorusLen)
	}
	rows := (bars + barsPerRow - 1) / barsPerRow
	last := max(rows-chartRows, 0) // the lowest top row
	top := g.restRow
	if playing >= 0 {
		top = min(max(playing/barsPerRow-1, 0), last)
	}
	g.scroll.aim(float64(top), now)
	at := g.scroll.row(now)

	// The page: the rows showing, whole or cut at its edges.
	pageH := (chartRows-1)*rowH + barH
	s := c.Scale
	clip := image.Rect(int((chartX-markRoom)*s), int(chartY*s), int(screenWidth*s), int(float64(chartY+pageH)*s))
	page := screen.Canvas{Dst: c.Dst.SubImage(clip).(*ebiten.Image), Scale: s}
	pen := g.pen()
	for i, cells := range g.bars {
		row := float64(i / barsPerRow)
		if row < math.Floor(at) || row > at+chartRows {
			continue
		}
		x := float32(chartX + i%barsPerRow*barW)
		y := float32(chartY + (row-at)*rowH)
		var fill color.Color
		if i == playing {
			fill = pale
		}
		pen.DrawBar(page, cells, x, y, barW, barH, fill, g.small, i%barsPerRow == barsPerRow-1 || i == bars-1)
		if label, ok := g.rehearsals[i]; ok {
			pen.DrawRehearsal(page, label, x, y)
		}
		if i == playing {
			cx := x + float32(math.Mod(pos, perBar)/perBar)*barW
			page.Line(cx, y, cx, y+barH, 1.5, ink)
		}
		for beat := range perBar {
			if k, ok := g.markAt(i*perBar+beat, chorus); ok {
				chart.DrawMark(page, k, x+markDX+float32(beat)*barW/perBar, y+markDY)
			}
		}
	}

	// The chevrons, once the page has arrived.
	mid := float32(chartX + barsPerRow*barW/2)
	if at > 0.01 {
		drawChevron(c, mid, chartY-chevronGap, -1)
	}
	if at < float64(last)-0.01 {
		drawChevron(c, mid, float32(chartY+pageH)+chevronGap, 1)
	}
}

// chevronGap is how far a chevron stands from the page; chevronW, how
// wide it is.
const (
	chevronGap = 6
	chevronW   = 10
)

// drawChevron draws a chevron centred on (`x`, `y`), pointing up for
// `dir` -1, down for 1: rows hidden that way.
func drawChevron(c screen.Canvas, x, y float32, dir float32) {
	tip := y + dir*chevronW/4
	back := y - dir*chevronW/4
	c.Line(x-chevronW/2, back, x, tip, 1.5, faint)
	c.Line(x, tip, x+chevronW/2, back, 1.5, faint)
}

// scrollRest follows the wheel, and a finger dragged on the chart, while
// nothing plays: the page moves a row at a time, the player reading the
// grid before he plays it.
func (g *game) scrollRest() {
	rows := (len(g.bars) + barsPerRow - 1) / barsPerRow
	last := max(rows-chartRows, 0)
	if _, dy := ebiten.Wheel(); dy != 0 {
		step := 1
		if dy > 0 {
			step = -1 // the wheel up shows the rows above
		}
		g.restRow = min(max(g.restRow+step, 0), last)
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		lx, ly := float64(x)/g.scale, float64(y)/g.scale
		if lx >= chartX && ly >= chartY && ly < chartY+chartRows*rowH {
			g.drag = &drag{id: id, y: ly, row: g.restRow}
		}
	}
	if g.drag == nil {
		return
	}
	if inpututil.IsTouchJustReleased(g.drag.id) {
		g.drag = nil
		return
	}
	_, y := ebiten.TouchPosition(g.drag.id)
	moved := int(math.Round((g.drag.y - float64(y)/g.scale) / rowH))
	g.restRow = min(max(g.drag.row+moved, 0), last)
}

// A drag is a finger on the chart at rest: which one, where it touched
// it, and the row at the top then.
type drag struct {
	id  ebiten.TouchID
	y   float64
	row int
}
