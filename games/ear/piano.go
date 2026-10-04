package main

import (
	"image/color"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// The keyboard on screen: three octaves, C3 to C6, lit by whatever is
// sounding, from any source (see screen.Piano).
//
// # What it must never do
//
// Give the answer away. During a question it shows only what sounds,
// note by note, which is what the ear hears anyway. The whole mode is
// marked on the keys only once it is revealed.

func newPiano() screen.Piano {
	return screen.Piano{Low: 48, High: 84, X: 40, Y: 262, W: 560, H: 72} // C3 to C6
}

// lit is the light of a key that sounds.
var lit = color.RGBA{0xe0, 0xa0, 0x40, 0xff}

// A tint is how a revealed key is painted: light on a white key, dark
// on a black one, so that the two rows stay apart once coloured.
type tint struct{ light, dark color.RGBA }

var (
	tonicTint  = tint{color.RGBA{0x6c, 0xb4, 0xe8, 0xff}, color.RGBA{0x2a, 0x6a, 0xa8, 0xff}}
	rightTint  = tint{color.RGBA{0x8f, 0xd9, 0xb6, 0xff}, color.RGBA{0x2a, 0x7a, 0x5a, 0xff}}
	wrongTint  = tint{color.RGBA{0xe8, 0x4a, 0x4a, 0xff}, color.RGBA{0xb0, 0x22, 0x22, 0xff}}
	sharedTint = tint{color.RGBA{0xc4, 0xc2, 0xbc, 0xff}, color.RGBA{0x5a, 0x5a, 0x64, 0xff}}
)

// A reveal says what to mark on the keys once the answer is out: the
// classes of the right mode, and when the player chose another, the
// classes of that one. Where they differ is what the ear had to catch.
type reveal struct {
	on      bool
	tonic   harmony.PitchClass
	right   harmony.PitchSet
	chosen  harmony.PitchSet
	mistake bool

	// labels holds the name of each marked class, spelled in the mode
	// it belongs to: mi♭ and fa♯ in G harmonic minor, never ré♯ or
	// sol♭. Empty for the classes left unmarked.
	labels [12]string
}

func (r reveal) tint(k int) (tint, bool) {
	if !r.on {
		return tint{}, false
	}
	c := harmony.PitchClass(k % 12)
	inRight, inChosen := r.right.Contains(c), r.mistake && r.chosen.Contains(c)
	switch {
	case c == r.tonic:
		return tonicTint, true
	case inRight && r.mistake && !inChosen:
		return rightTint, true
	case inChosen && !inRight:
		return wrongTint, true
	case inRight && r.mistake:
		return sharedTint, true
	case inRight:
		return rightTint, true
	}
	return tint{}, false
}

// dress paints the keys the reveal marks: its tint whole, light on a
// white key, dark on a black one, and its name in `f` at the bottom.
func (r reveal) dress(f *screen.Font) func(k int) screen.Dress {
	return func(k int) screen.Dress {
		t, ok := r.tint(k)
		if !ok {
			return screen.Dress{}
		}
		base := t.light
		if screen.IsBlack(k) {
			base = t.dark
		}
		return screen.Dress{Base: base, Label: r.labels[k%12], Font: f}
	}
}
