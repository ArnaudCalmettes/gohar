package main

import (
	_ "embed"
	"image/color"

	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// The game draws in logical coordinates, screenWidth by screenHeight
// (see screen.Canvas).
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

// fonts are the faces the game uses: the chart's hand for the chords,
// the countdown and the walker's bubbles, Go Regular for the rest.
type fonts struct {
	chord, count, bubble, ui *screen.Font
}

// Their sizes, in logical units.
const (
	chordSize  = 24
	countSize  = 64
	bubbleSize = 14
	uiSize     = 11
)

func newFonts() (fonts, error) {
	var fs fonts
	var err error
	if fs.chord, err = screen.NewFont(chordSize, museJazzText); err != nil {
		return fs, err
	}
	if fs.count, err = screen.NewFont(countSize, museJazzText); err != nil {
		return fs, err
	}
	if fs.bubble, err = screen.NewFont(bubbleSize, museJazzText, screen.GoRegular); err != nil { // Go Regular for any letter the hand lacks
		return fs, err
	}
	fs.ui, err = screen.NewFont(uiSize, screen.GoRegular)
	return fs, err
}
