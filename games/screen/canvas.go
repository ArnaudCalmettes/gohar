// Package screen holds what the games draw with: a canvas in logical
// coordinates, their fonts, and the keyboard on screen (see "Les
// scènes" in docs/architecture.md).
package screen

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// A Canvas draws in logical coordinates and renders them at the
// window's real resolution, `Scale` times larger.
//
// # Why not let Ebitengine scale
//
// A layout of 640 by 360 in a 1280 by 720 window is scaled up as an
// image, and every edge shows its pixels. Rendering at the real size
// and scaling the coordinates instead keeps shapes and text sharp at any
// window size and on a high density screen, and the drawing code still
// thinks in one fixed grid.
type Canvas struct {
	Dst   *ebiten.Image
	Scale float64
}

func (c Canvas) Rect(x, y, w, h float32, col color.Color) {
	s := float32(c.Scale)
	vector.DrawFilledRect(c.Dst, x*s, y*s, w*s, h*s, col, true)
}

func (c Canvas) Circle(x, y, r float32, col color.Color) {
	s := float32(c.Scale)
	vector.DrawFilledCircle(c.Dst, x*s, y*s, r*s, col, true)
}

func (c Canvas) Line(x0, y0, x1, y1, width float32, col color.Color) {
	s := float32(c.Scale)
	vector.StrokeLine(c.Dst, x0*s, y0*s, x1*s, y1*s, width*s, col, true)
}

// Text draws `s` with its top left corner on (`x`, `y`).
func (c Canvas) Text(s string, f *Font, x, y float64, col color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x*c.Scale, y*c.Scale)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(c.Dst, s, f.at(c.Scale), op)
}

// Centred draws `s` centred on `x`.
func (c Canvas) Centred(s string, f *Font, x, y float64, col color.Color) {
	w, _ := c.Measure(s, f)
	c.Text(s, f, x-w/2, y, col)
}

// Measure returns the size of `s` in logical units.
func (c Canvas) Measure(s string, f *Font) (w, h float64) {
	w, h = text.Measure(s, f.at(c.Scale), 0)
	return w / c.Scale, h / c.Scale
}
