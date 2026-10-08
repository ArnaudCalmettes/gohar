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

	// tapPad widens the box of an item, left and right, for a finger.
	tapPad = 12
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

// A hit is the box of an item of the list last drawn, in logical
// units: where a click or a tap chooses it.
type hit struct{ x, y, w, h float64 }

// has tells whether the point (`x`, `y`) is in the box.
func (h hit) has(x, y float64) bool {
	return x >= h.x && x < h.x+h.w && y >= h.y && y < h.y+h.h
}

// pressed returns where the screen was just clicked or touched, in
// logical units, if it was.
func (a *app) pressed() (x, y float64, ok bool) {
	var px, py int
	switch ids := inpututil.AppendJustPressedTouchIDs(nil); {
	case len(ids) > 0:
		px, py = ebiten.TouchPosition(ids[0])
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		px, py = ebiten.CursorPosition()
	default:
		return 0, 0, false
	}
	return float64(px) / a.scale, float64(py) / a.scale, true // the screen is scaled (see layout)
}

// A tap is a click or a touch on an item of the list last drawn: its
// index, and `side`, -1 left of its middle, 1 right of it, for the
// options that a tap changes rather than opens.
type tap struct {
	item, side int
}

// tapped returns the item of the list last drawn just clicked or
// touched, if any. The boxes are those of the last frame: Update runs
// before Draw, and the list stays where it was.
func (a *app) tapped() (tap, bool) {
	lx, ly, ok := a.pressed()
	if !ok {
		return tap{}, false
	}
	for i, h := range a.hits {
		if h.has(lx, ly) {
			side := 1
			if lx < h.x+h.w/2 {
				side = -1
			}
			return tap{item: i, side: side}, true
		}
	}
	return tap{}, false
}

// drawMenu draws a menu screen on the jam's beats: the walker `wk` in
// gait `g`, the heading in `head`, the `labels`, and the keys in the
// status line.
func (a *app) drawMenu(dst *ebiten.Image, wk *figure.Walker, g figure.Gait, heading string, head *screen.Font, labels []string, chosen int, keys string) {
	a.drawList(dst, wk, g, heading, head, entries(labels, nil), chosen, keys, layout{top: menuY, gap: menuGap})
}

// A layout places a list: its heading at `heading` (headingY when 0),
// its first item at `top`, `gap` from one item to the next, and `apart`
// more before the last one, Back, so that it does not read as one of
// the list. `left` aligns the items but Back on their left, the block of
// them centred: their checks fall in a column. The items under another,
// a tree's branches, come `branchGap` apart, indented, in a smaller
// hand, a little further from the next item, for their highlight to
// stay clear of it.
type layout struct {
	heading, top, gap, apart, branchGap float64
	left                                bool
}

// An entry is an item of a list as drawList draws it: its label, a check
// when done, and whether it is a branch, under the item before it.
type entry struct {
	label  string
	done   bool
	branch bool
}

// branchIndent is how far a branch stands right of its item.
const branchIndent = 28

// entries makes plain entries of `labels`, a check on those `done`
// marks; `done` may be nil, or shorter than `labels`.
func entries(labels []string, done []bool) []entry {
	out := make([]entry, len(labels))
	for i, l := range labels {
		out[i] = entry{label: l, done: i < len(done) && done[i]}
	}
	return out
}

// drawList is drawMenu with its list of `items` laid out by `l`, a green
// check left of each one done; the last item is Back. It keeps the box
// of each item, for tapped.
func (a *app) drawList(dst *ebiten.Image, wk *figure.Walker, g figure.Gait, heading string, head *screen.Font, items []entry, chosen int, keys string, l layout) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: a.scale}
	if wk != nil { // none in the settings of a grid, over the game
		now := time.Now()
		wk.Draw(c, walkerX, walkerY, walkerScale, g, a.jam.Metronome().Beats(now), now, faint)
	}

	const mid = screenWidth / 2
	hy := l.heading
	if hy == 0 {
		hy = headingY
	}
	c.Centred(heading, head, mid, hy, ink)

	// The left edge of the items, Back and the branches aside, when
	// aligned on it.
	widest := 0.0
	for _, it := range items[:len(items)-1] {
		if !it.branch {
			w, _ := c.Measure(it.label, a.fonts.chord)
			widest = max(widest, w)
		}
	}
	a.hits = a.hits[:0]
	y := l.top
	for i, it := range items {
		back := i == len(items)-1
		switch {
		case i == 0:
		case back:
			y += l.gap + l.apart
		case it.branch && !items[i-1].branch:
			y += l.gap // under its item
		case it.branch:
			y += l.branchGap
		case items[i-1].branch:
			y += l.gap + menuPad
		default:
			y += l.gap
		}
		font := a.fonts.chord
		if it.branch {
			font = a.fonts.chordSmall
		}
		w, h := c.Measure(it.label, font)
		x := mid - w/2
		if l.left && !back {
			x = mid - widest/2
			if it.branch {
				x += branchIndent
			}
		}
		// The box: the highlight, a little wider for a finger.
		a.hits = append(a.hits, hit{x: x - menuPad - tapPad, y: y - menuPad, w: w + 2*(menuPad+tapPad), h: h + 2*menuPad})
		var col color.Color = faint
		if i == chosen {
			c.Rect(float32(x-menuPad), float32(y-menuPad), float32(w+2*menuPad), float32(h+2*menuPad), pale)
			col = ink
		}
		c.Text(it.label, font, x, y, col)
		if it.done {
			drawCheck(c, x-menuPad-checkGap-checkW, y+h/2)
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
