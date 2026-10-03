package main

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

// The game draws in logical coordinates, screenWidth by screenHeight,
// and the canvas renders them at the window's real resolution, as in
// ear: shapes and text stay sharp at any size.
const (
	screenWidth  = 640
	screenHeight = 360
)

// museJazzText is the hand of the charts, as on a Real Book. See
// fonts/README.md.
//
//go:embed fonts/MuseJazzText.otf
var museJazzText []byte

var (
	paper = color.White
	ink   = color.Black
	pale  = color.RGBA{0xe8, 0xe8, 0xe8, 0xff} // the bar being played
	faint = color.RGBA{0x88, 0x88, 0x88, 0xff} // what is said rather than shown

	// The marks of the arrivals: landed, missed, doubled.
	landedInk  = color.RGBA{0x2a, 0x9d, 0x5a, 0xff}
	missedInk  = color.RGBA{0xc0, 0x39, 0x2b, 0xff}
	doubledInk = color.RGBA{0xe0, 0x8a, 0x1e, 0xff}
)

type canvas struct {
	dst   *ebiten.Image
	scale float64
}

func (c canvas) rect(x, y, w, h float32, col color.Color) {
	s := float32(c.scale)
	vector.DrawFilledRect(c.dst, x*s, y*s, w*s, h*s, col, true)
}

func (c canvas) circle(x, y, r float32, col color.Color) {
	s := float32(c.scale)
	vector.DrawFilledCircle(c.dst, x*s, y*s, r*s, col, true)
}

func (c canvas) line(x0, y0, x1, y1, width float32, col color.Color) {
	s := float32(c.scale)
	vector.StrokeLine(c.dst, x0*s, y0*s, x1*s, y1*s, width*s, col, true)
}

func (c canvas) text(s string, f *font, x, y float64, col color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x*c.scale, y*c.scale)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(c.dst, s, f.at(c.scale), op)
}

// centred draws `s` centred on `x`.
func (c canvas) centred(s string, f *font, x, y float64, col color.Color) {
	w, _ := c.measure(s, f)
	c.text(s, f, x-w/2, y, col)
}

// measure returns the size of `s` in logical units.
func (c canvas) measure(s string, f *font) (w, h float64) {
	w, h = text.Measure(s, f.at(c.scale), 0)
	return w / c.scale, h / c.scale
}

// A font is one face at one logical size, rasterised at the window's
// resolution each time it is used.
type font struct {
	size float64
	face *text.GoTextFace
}

func newFont(ttf []byte, size float64) (*font, error) {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(ttf))
	if err != nil {
		return nil, err
	}
	return &font{size: size, face: &text.GoTextFace{Source: src, Size: size}}, nil
}

func (f *font) at(scale float64) text.Face {
	f.face.Size = f.size * scale
	return f.face
}

// fonts are the faces the game uses: the chart's hand for the chords
// and the countdown, Go Regular for the rest.
type fonts struct {
	chord, count, ui *font
}

// Their sizes, in logical units.
const (
	chordSize = 24
	countSize = 64
	uiSize    = 11
)

func newFonts() (fonts, error) {
	var fs fonts
	var err error
	if fs.chord, err = newFont(museJazzText, chordSize); err != nil {
		return fs, err
	}
	if fs.count, err = newFont(museJazzText, countSize); err != nil {
		return fs, err
	}
	fs.ui, err = newFont(goregular.TTF, uiSize)
	return fs, err
}
