package screen

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// GoRegular is the type of everything a game says rather than shows:
// menus, status lines.
var GoRegular = goregular.TTF

// A Font is one size of one type, falling back on further types for
// the glyphs it lacks: Go Regular, say, then a font of musical signs.
// It is set at the window's resolution each time it is used, so that
// glyphs are rasterised at their real size rather than scaled up.
type Font struct {
	size  float64
	faces []*text.GoTextFace
	face  text.Face
}

// NewFont makes a font of `size` logical units from the font files
// `files`, the first one first, the others for what it lacks.
func NewFont(size float64, files ...[]byte) (*Font, error) {
	f := &Font{size: size}
	var faces []text.Face
	for _, file := range files {
		src, err := text.NewGoTextFaceSource(bytes.NewReader(file))
		if err != nil {
			return nil, err
		}
		face := &text.GoTextFace{Source: src, Size: size}
		f.faces = append(f.faces, face)
		faces = append(faces, face)
	}
	if len(faces) == 1 {
		f.face = faces[0]
		return f, nil
	}
	// The multi face holds the pointers, so resizing them in at resizes
	// it too.
	var err error
	f.face, err = text.NewMultiFace(faces...)
	return f, err
}

func (f *Font) at(scale float64) text.Face {
	for _, face := range f.faces {
		face.Size = f.size * scale
	}
	return f.face
}
