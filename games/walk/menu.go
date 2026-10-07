package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/chart"
	"github.com/ArnaudCalmettes/gohar/games/walk/figure"
)

// The menus, the title's, the options' and the course's: the walker
// greyed on the left, on the jam, a heading and a list in the middle
// (see "L'écran titre" in docs/walk.md).
const (
	headingY = 70  // the top of the heading
	menuY    = 170 // the top of the first item
	menuGap  = 36  // from one item to the next
	menuPad  = 6   // around the item chosen

	// The check of an item done, left of its label.
	checkW   = 14
	checkGap = 10
)

// A menu is the list the menus share: the item chosen, moved by the
// arrows.
type menu struct {
	chosen, items int
}

// move follows the up and down arrows, round the list.
func (m *menu) move() {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp):
		m.chosen = (m.chosen + m.items - 1) % m.items
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown):
		m.chosen = (m.chosen + 1) % m.items
	}
}

// confirmed tells whether the item chosen was just confirmed.
func confirmed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// drawMenu draws a menu screen on the jam's beats: the walker `wk` in
// gait `g`, the heading in `head`, the `labels`, and the keys in the
// status line.
func (a *app) drawMenu(dst *ebiten.Image, wk *figure.Walker, g figure.Gait, heading string, head *screen.Font, labels []string, chosen int, keys string) {
	a.drawList(dst, wk, g, heading, head, labels, nil, chosen, keys, layout{top: menuY, gap: menuGap})
}

// A layout places a list: its first item at `top`, `gap` from one item
// to the next, and `apart` more before the last one, Back, so that it
// does not read as one of the list.
type layout struct {
	top, gap, apart float64
}

// drawList is drawMenu with its list laid out by `l`, and a green check
// left of each item `done` marks. `done` may be nil, or shorter than
// `labels`.
func (a *app) drawList(dst *ebiten.Image, wk *figure.Walker, g figure.Gait, heading string, head *screen.Font, labels []string, done []bool, chosen int, keys string, l layout) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: a.scale}
	now := time.Now()
	wk.Draw(c, walkerX, walkerY, walkerScale, g, a.jam.Metronome().Beats(now), now, faint)

	const mid = screenWidth / 2
	c.Centred(heading, head, mid, headingY, ink)
	for i, label := range labels {
		y := l.top + float64(i)*l.gap
		if i == len(labels)-1 {
			y += l.apart
		}
		w, h := c.Measure(label, a.fonts.chord)
		var col color.Color = faint
		if i == chosen {
			c.Rect(float32(mid-w/2-menuPad), float32(y-menuPad), float32(w+2*menuPad), float32(h+2*menuPad), pale)
			col = ink
		}
		c.Centred(label, a.fonts.chord, mid, y, col)
		if i < len(done) && done[i] {
			drawCheck(c, mid-w/2-menuPad-checkGap-checkW, y+h/2)
		}
	}
	c.Text(keys, a.fonts.ui, margin, statusY, faint)
}

// drawCheck draws a check mark from `x`, centred on `y`.
func drawCheck(c screen.Canvas, x, y float64) {
	const width = 2.5
	x0, y0 := float32(x), float32(y)
	c.Line(x0, y0, x0+checkW*0.35, y0+checkW*0.35, width, chart.LandedInk)
	c.Line(x0+checkW*0.35, y0+checkW*0.35, x0+checkW, y0-checkW*0.45, width, chart.LandedInk)
}
