package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// The menus, the title's and the options': the walker greyed on the
// left, on the jam, a heading and a list in the middle (see "L'écran
// titre" in docs/walk.md).
const (
	headingY = 70  // the top of the heading
	menuY    = 170 // the top of the first item
	menuGap  = 36  // from one item to the next
	menuPad  = 6   // around the item chosen
)

// A menu is the list both menus share: the item chosen, moved by the
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
func (a *app) drawMenu(dst *ebiten.Image, wk *walker, g gait, heading string, head *screen.Font, labels []string, chosen int, keys string) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: a.scale}
	now := time.Now()
	wk.draw(c, walkerX, walkerY, walkerScale, g, a.jam.m.Beats(now), now, faint)

	const mid = screenWidth / 2
	c.Centred(heading, head, mid, headingY, ink)
	for i, label := range labels {
		y := float64(menuY + i*menuGap)
		var col color.Color = faint
		if i == chosen {
			w, h := c.Measure(label, a.fonts.chord)
			c.Rect(float32(mid-w/2-menuPad), float32(y-menuPad), float32(w+2*menuPad), float32(h+2*menuPad), pale)
			col = ink
		}
		c.Centred(label, a.fonts.chord, mid, y, col)
	}
	c.Text(keys, a.fonts.ui, margin, statusY, faint)
}
