package main

import (
	"image/color"

	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// The keyboard on screen (see screen.Piano), lit by what the player
// holds and by what the band plays in demo, over a wider range than
// ear's: from the double bass's low E up to the right hand's comping.

func newPiano() screen.Piano {
	return screen.Piano{
		Low: 24, High: 84, // C1, under the low E of the double bass, to C6
		X: margin, Y: 250, W: screenWidth - 2*margin, H: 70,
	}
}

var (
	// The player warm, the demo cool: what the hands play and what the
	// band plays for them read apart at a glance.
	playerLit = color.RGBA{0xe0, 0xa0, 0x40, 0xff}
	demoLit   = color.RGBA{0x6c, 0xb4, 0xe8, 0xff}
)

// splitTick is the mark of the split, above the keys.
const splitTick = 5

// drawSplit marks the lowest key of the right hand, `split`, with a thin
// line above the keys.
func drawSplit(c screen.Canvas, p *screen.Piano, split int) {
	if split > p.Low && split <= p.High {
		x, y, _, _ := p.KeyRect(split)
		c.Line(x, y-1-splitTick, x, y-1, 1, faint)
	}
}
