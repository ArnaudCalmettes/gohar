package main

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

// The game draws in logical coordinates, screenWidth by screenHeight,
// and a canvas renders them at the window's real resolution.
//
// # Why not let Ebitengine scale
//
// A layout of 640 by 360 in a 1280 by 720 window is scaled up as an
// image, and every edge shows its pixels. Rendering at the real size
// and scaling the coordinates instead keeps shapes and text sharp at any
// window size and on a high density screen, and the drawing code still
// thinks in one fixed grid.
type canvas struct {
	dst   *ebiten.Image
	scale float64
}

func (c canvas) rect(x, y, w, h float32, col color.Color) {
	s := float32(c.scale)
	vector.DrawFilledRect(c.dst, x*s, y*s, w*s, h*s, col, true)
}

func (c canvas) text(s string, f *font, x, y float64, col color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x*c.scale, y*c.scale)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(c.dst, s, f.at(c.scale), op)
}

// measure returns the size of s in logical units.
func (c canvas) measure(s string, f *font) (w, h float64) {
	w, h = text.Measure(s, f.at(c.scale), 0)
	return w / c.scale, h / c.scale
}

// A font is one size of the game's type: Go Regular, falling back on
// the signs font for the glyphs it lacks. It is set at the window's
// resolution each time it is used, so that glyphs are rasterised at
// their real size rather than scaled up.
type font struct {
	size    float64
	regular *text.GoTextFace
	signs   *text.GoTextFace
	face    text.Face
}

func newFont(size float64) (*font, error) {
	regular, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}
	signs, err := text.NewGoTextFaceSource(bytes.NewReader(signsTTF))
	if err != nil {
		return nil, err
	}
	f := &font{
		size:    size,
		regular: &text.GoTextFace{Source: regular, Size: size},
		signs:   &text.GoTextFace{Source: signs, Size: size},
	}
	// The multi face holds the two pointers, so resizing them in at
	// resizes it too.
	f.face, err = text.NewMultiFace(f.regular, f.signs)
	return f, err
}

func (f *font) at(scale float64) text.Face {
	f.regular.Size = f.size * scale
	f.signs.Size = f.size * scale
	return f.face
}
