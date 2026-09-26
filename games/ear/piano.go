package main

import (
	"image/color"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// The keyboard on screen: three octaves, C3 to C6, lit by whatever is
// sounding, from any source.
//
// # What it must never do
//
// Give the answer away. During a question it shows only what sounds,
// note by note, which is what the ear hears anyway. The whole mode is
// marked on the keys only once it is revealed.

const (
	pianoLow  = 48 // C3
	pianoHigh = 84 // C6

	pianoX      = 40
	pianoY      = 262
	pianoWidth  = 560
	pianoHeight = 72

	// glowDecay is what a released key keeps of its light each tick:
	// at 60 ticks a second, it is mostly gone in a quarter of a second,
	// long enough for the eye to follow a scale, short enough for the
	// next note to stand out.
	glowDecay = 0.85

	// ghost dims a key folded in from outside the range: the pedal two
	// octaves below, or a MIDI keyboard set to another octave. It still
	// shows which class sounds, without pretending to be that key.
	ghost = 0.5
)

var (
	whiteKey  = color.RGBA{0xe8, 0xe6, 0xe0, 0xff}
	blackKey  = color.RGBA{0x2a, 0x2a, 0x32, 0xff}
	keyBorder = color.RGBA{0x1c, 0x1c, 0x22, 0xff}
	lit       = color.RGBA{0xe0, 0xa0, 0x40, 0xff}

	darkLabel  = color.RGBA{0x1c, 0x1c, 0x22, 0xff}
	lightLabel = color.RGBA{0xe8, 0xe6, 0xe0, 0xff}
)

// A tint is how a revealed key is painted: light on a white key, dark
// on a black one, so that the two rows stay apart once coloured.
type tint struct{ light, dark color.RGBA }

var (
	tonicTint  = tint{color.RGBA{0x6c, 0xb4, 0xe8, 0xff}, color.RGBA{0x2a, 0x6a, 0xa8, 0xff}}
	rightTint  = tint{color.RGBA{0x8f, 0xd9, 0xb6, 0xff}, color.RGBA{0x2a, 0x7a, 0x5a, 0xff}}
	wrongTint  = tint{color.RGBA{0xe8, 0x9a, 0x9a, 0xff}, color.RGBA{0x8a, 0x3b, 0x3b, 0xff}}
	sharedTint = tint{color.RGBA{0xc4, 0xc2, 0xbc, 0xff}, color.RGBA{0x5a, 0x5a, 0x64, 0xff}}
)

// A piano is the keyboard's state from one frame to the next: how lit
// each MIDI key still is.
type piano struct {
	glow [128]float32
}

// update lights what is down and lets the rest fade.
func (p *piano) update(down *[128]bool) {
	for k := range p.glow {
		if down[k] {
			p.glow[k] = 1
		} else {
			p.glow[k] *= glowDecay
			if p.glow[k] < 0.01 {
				p.glow[k] = 0
			}
		}
	}
}

// fold brings a key into the displayed range by octaves.
func fold(k int) int {
	for k < pianoLow {
		k += 12
	}
	for k > pianoHigh {
		k -= 12
	}
	return k
}

// levels folds every lit key into the range, keeping the brightest
// light per displayed key.
func (p *piano) levels() [128]float32 {
	var out [128]float32
	for k, g := range p.glow {
		if g == 0 {
			continue
		}
		d := fold(k)
		if d != k {
			g *= ghost
		}
		out[d] = max(out[d], g)
	}
	return out
}

func isBlack(k int) bool {
	switch k % 12 {
	case 1, 3, 6, 8, 10:
		return true
	}
	return false
}

// whiteCount is the number of white keys in the range.
func whiteCount() int {
	n := 0
	for k := pianoLow; k <= pianoHigh; k++ {
		if !isBlack(k) {
			n++
		}
	}
	return n
}

// keyRect places key k. A black key sits across the boundary of the
// white key below it.
func keyRect(k int) (x, y, w, h float32) {
	whiteW := float32(pianoWidth) / float32(whiteCount())
	whites := 0
	for i := pianoLow; i < k; i++ {
		if !isBlack(i) {
			whites++
		}
	}
	if !isBlack(k) {
		return pianoX + float32(whites)*whiteW, pianoY, whiteW, pianoHeight
	}
	bw := whiteW * 0.6
	return pianoX + float32(whites)*whiteW - bw/2, pianoY, bw, pianoHeight * 0.62
}

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

// draw paints the keyboard: whites, then blacks over them. A revealed
// key takes its tint whole and its name at the bottom; whatever sounds
// glows over it, so that playing stays readable after the reveal.
func (p *piano) draw(c canvas, r reveal, f *font) {
	levels := p.levels()
	for _, black := range []bool{false, true} {
		for k := pianoLow; k <= pianoHigh; k++ {
			if isBlack(k) != black {
				continue
			}
			x, y, w, h := keyRect(k)
			base, label := whiteKey, darkLabel
			if black {
				base, label = blackKey, lightLabel
			}
			if t, ok := r.tint(k); ok {
				base, label = t.light, darkLabel
				if black {
					base, label = t.dark, lightLabel
				}
			}
			// Pressed keys sink a little. A placeholder: it reads as a
			// press, but looks cheap, and a better rendering is to find.
			if levels[k] > 0.95 {
				y += 2
			}
			c.rect(x, y, w, h, keyBorder)
			c.rect(x+1, y, w-2, h-1, blend(base, lit, levels[k]))

			if name := r.labels[k%12]; r.on && name != "" {
				tw, th := c.measure(name, f)
				c.text(name, f, float64(x)+(float64(w)-tw)/2, float64(y+h)-th-4, label)
			}
		}
	}
}

// blend moves from a toward b by t, between 0 and 1.
func blend(a, b color.RGBA, t float32) color.RGBA {
	mix := func(x, y uint8) uint8 { return uint8(float32(x) + (float32(y)-float32(x))*t) }
	return color.RGBA{mix(a.R, b.R), mix(a.G, b.G), mix(a.B, b.B), 0xff}
}
