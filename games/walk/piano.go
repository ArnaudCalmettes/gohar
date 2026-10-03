package main

import "image/color"

// The keyboard on screen, lit by what the player holds and by what the
// band plays in demo: the one of ear (games/ear/piano.go), without its
// reveal, over a wider range, from the double bass's low E up to the
// right hand's comping.
//
// A copy for now, trimmed to what walk needs. The day a third game
// wants it, it moves to a package of its own.

const (
	midiKeys = 128 // 0 to 127

	pianoLow  = 24 // C1, under the low E of the double bass
	pianoHigh = 84 // C6

	pianoX      = margin
	pianoY      = 250
	pianoWidth  = screenWidth - 2*margin
	pianoHeight = 70

	// A black key, against a white one.
	blackWidth  = 0.6
	blackLength = 0.62

	// glowDecay is what a released key keeps of its light each tick:
	// mostly gone in a quarter of a second, as in ear.
	glowDecay = 0.85

	// glowOff is the light under which a key is dark again, and pressed
	// the light over which it is drawn down.
	glowOff = 0.01
	pressed = 0.95

	// ghost dims a key folded in from outside the range.
	ghost = 0.5

	whiteLip   = 4
	blackLip   = 3
	lipShade   = 0.78
	shadowRows = 3
	shadowFrom = 0.55

	splitTick = 5 // the mark of the split, above the keys
)

var (
	whiteKey  = color.RGBA{0xe8, 0xe6, 0xe0, 0xff}
	blackKey  = color.RGBA{0x2a, 0x2a, 0x32, 0xff}
	keyBorder = color.RGBA{0x1c, 0x1c, 0x22, 0xff}

	// The player warm, the demo cool: what the hands play and what the
	// band plays for them read apart at a glance.
	playerLit = color.RGBA{0xe0, 0xa0, 0x40, 0xff}
	demoLit   = color.RGBA{0x6c, 0xb4, 0xe8, 0xff}
)

// A piano is the keyboard's state from one frame to the next: how lit
// each MIDI key still is.
type piano struct {
	glow   [midiKeys]float32
	byDemo [midiKeys]bool // the light comes from the demo, not the player
	split  int            // the lowest key of the right hand
}

// update lights what the player holds and what the demo plays, the
// player first when both play the same key, and lets the rest fade.
func (p *piano) update(player, demo *[midiKeys]bool) {
	for k := range p.glow {
		switch {
		case player[k]:
			p.glow[k], p.byDemo[k] = 1, false
		case demo[k]:
			p.glow[k], p.byDemo[k] = 1, true
		default:
			p.glow[k] *= glowDecay
			if p.glow[k] < glowOff {
				p.glow[k] = 0
			}
		}
	}
}

// fold brings a key into the displayed range by octaves.
func fold(k int) int {
	for k < pianoLow {
		k += octave
	}
	for k > pianoHigh {
		k -= octave
	}
	return k
}

// levels folds every lit key into the range, keeping the brightest
// light per displayed key, and whether the demo lit it.
func (p *piano) levels() (out [midiKeys]float32, demo [midiKeys]bool) {
	for k, g := range p.glow {
		if g == 0 {
			continue
		}
		d := fold(k)
		if d != k {
			g *= ghost
		}
		if g > out[d] {
			out[d], demo[d] = g, p.byDemo[k]
		}
	}
	return out, demo
}

func isBlack(k int) bool {
	switch k % octave {
	case 1, 3, 6, 8, 10:
		return true
	}
	return false
}

func whiteCount() int {
	n := 0
	for k := pianoLow; k <= pianoHigh; k++ {
		if !isBlack(k) {
			n++
		}
	}
	return n
}

// keyRect places key `k`. A black key sits across the boundary of the
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
	bw := whiteW * blackWidth
	return pianoX + float32(whites)*whiteW - bw/2, pianoY, bw, pianoHeight * blackLength
}

// draw paints the keyboard, whites then blacks over them, and marks
// the split with a thin line above the keys.
func (p *piano) draw(c canvas) {
	levels, demo := p.levels()
	for _, black := range []bool{false, true} {
		for k := pianoLow; k <= pianoHigh; k++ {
			if isBlack(k) != black {
				continue
			}
			x, y, w, h := keyRect(k)
			base := whiteKey
			if black {
				base = blackKey
			}
			glow := playerLit
			if demo[k] {
				glow = demoLit
			}
			top := blend(base, glow, levels[k])
			lip := float32(whiteLip)
			if black {
				lip = blackLip
			}
			c.rect(x, y, w, h, keyBorder)
			if levels[k] > pressed {
				// Down: the lip goes under, the top runs to the edge,
				// darkened at the back where it sinks.
				c.rect(x+1, y, w-2, h-1, top)
				for i := range shadowRows {
					f := shadowFrom + (1-shadowFrom)*float32(i)/shadowRows
					c.rect(x+1, y+float32(i), w-2, 1, shade(top, f))
				}
			} else {
				c.rect(x+1, y, w-2, h-1-lip, top)
				c.rect(x+1, y+h-1-lip, w-2, lip, shade(top, lipShade))
			}
		}
	}
	if p.split > pianoLow && p.split <= pianoHigh {
		x, y, _, _ := keyRect(p.split)
		c.line(x, y-1-splitTick, x, y-1, 1, faint)
	}
}

// shade darkens a colour to `f` of its brightness.
func shade(a color.RGBA, f float32) color.RGBA {
	return blend(color.RGBA{A: 0xff}, a, f)
}

// blend moves from `a` toward `b` by `t`, between 0 and 1.
func blend(a, b color.RGBA, t float32) color.RGBA {
	mix := func(x, y uint8) uint8 { return uint8(float32(x) + (float32(y)-float32(x))*t) }
	return color.RGBA{mix(a.R, b.R), mix(a.G, b.G), mix(a.B, b.B), 0xff}
}
