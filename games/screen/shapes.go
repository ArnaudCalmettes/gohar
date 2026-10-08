package screen

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// The shapes the buttons are drawn with: a frame with round corners,
// a polygon, a ring. Ebitengine fills a path with triangles of a
// plain image, tinted with the colour.

// whiteSubImage is the plain image the triangles are drawn from: a
// pixel inside a white image, its edge clear of the filter.
var whiteSubImage = func() *ebiten.Image {
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	return img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}()

// draw draws the triangles of a path, in `col`.
func (c Canvas) draw(vs []ebiten.Vertex, is []uint16, col color.Color) {
	r, g, b, a := col.RGBA()
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	c.Dst.DrawTriangles(vs, is, whiteSubImage, &ebiten.DrawTrianglesOptions{AntiAlias: true})
}

// RoundRect strokes a frame from (`x`, `y`), `w` by `h`, its corners
// rounded with radius `r`, its line `width` wide.
func (c Canvas) RoundRect(x, y, w, h, r, width float32, col color.Color) {
	s := float32(c.Scale)
	x, y, w, h, r = x*s, y*s, w*s, h*s, r*s
	var p vector.Path
	p.MoveTo(x+r, y)
	p.ArcTo(x+w, y, x+w, y+h, r)
	p.ArcTo(x+w, y+h, x, y+h, r)
	p.ArcTo(x, y+h, x, y, r)
	p.ArcTo(x, y, x+w, y, r)
	p.Close()
	vs, is := p.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: width * s, LineJoin: vector.LineJoinRound})
	c.draw(vs, is, col)
}

// Polygon fills the polygon of `points`, x and y in turn.
func (c Canvas) Polygon(points []float32, col color.Color) {
	s := float32(c.Scale)
	var p vector.Path
	for i := 0; i+1 < len(points); i += 2 {
		if i == 0 {
			p.MoveTo(points[i]*s, points[i+1]*s)
			continue
		}
		p.LineTo(points[i]*s, points[i+1]*s)
	}
	p.Close()
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	c.draw(vs, is, col)
}

// Ring strokes a circle centred on (`x`, `y`), of radius `r`.
func (c Canvas) Ring(x, y, r, width float32, col color.Color) {
	s := float32(c.Scale)
	vector.StrokeCircle(c.Dst, x*s, y*s, r*s, width*s, col, true)
}
