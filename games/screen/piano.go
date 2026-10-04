package screen

import "image/color"

// MIDIKeys is how many keys MIDI numbers, 0 to 127.
const MIDIKeys = 128

const (
	// glowDecay is what a released key keeps of its light each tick:
	// at 60 ticks a second, it is mostly gone in a quarter of a second,
	// long enough for the eye to follow a scale, short enough for the
	// next note to stand out.
	glowDecay = 0.85

	// glowOff is the light under which a key is dark again, and pressed
	// the light over which it is drawn down.
	glowOff = 0.01
	pressed = 0.95

	// ghost dims a key folded in from outside the range: a pedal two
	// octaves below, or a MIDI keyboard set to another octave. It still
	// shows which class sounds, without pretending to be that key.
	ghost = 0.5

	// A black key, against a white one.
	blackWidth  = 0.6
	blackLength = 0.62

	// A key at rest shows its front, a lip a shade darker than its top,
	// and a key held down hides it under the keys around it: what a
	// keyboard seen from above and in front shows. The pressed key also
	// takes a shadow where it goes under the fallboard.
	whiteLip   = 4
	blackLip   = 3
	lipShade   = 0.78
	shadowRows = 3
	shadowFrom = 0.55

	labelGap = 4 // a key's label, above its bottom edge
	octave   = 12
)

var (
	whiteKey  = color.RGBA{0xe8, 0xe6, 0xe0, 0xff}
	blackKey  = color.RGBA{0x2a, 0x2a, 0x32, 0xff}
	keyBorder = color.RGBA{0x1c, 0x1c, 0x22, 0xff}

	darkLabel  = color.RGBA{0x1c, 0x1c, 0x22, 0xff}
	lightLabel = color.RGBA{0xe8, 0xe6, 0xe0, 0xff}
)

// A Piano is the keyboard on screen, from `Low` to `High`, in the box
// at (`X`, `Y`), `W` by `H`, and its light from one frame to the next.
// A key outside the range is folded into it by octaves, dimmed.
type Piano struct {
	Low, High  int
	X, Y, W, H float32

	glow  [MIDIKeys]float32
	color [MIDIKeys]color.RGBA
}

// Lit is a set of keys down and the colour they light.
type Lit struct {
	Down  *[MIDIKeys]bool
	Color color.RGBA
}

// Update lights the keys down in `lights`, the first set first when two
// hold the same key, and lets the rest fade. Once a tick.
func (p *Piano) Update(lights ...Lit) {
	for k := range p.glow {
		lit := false
		for _, l := range lights {
			if l.Down != nil && l.Down[k] {
				p.glow[k], p.color[k], lit = 1, l.Color, true
				break
			}
		}
		if !lit {
			p.glow[k] *= glowDecay
			if p.glow[k] < glowOff {
				p.glow[k] = 0
			}
		}
	}
}

// fold brings a key into the displayed range by octaves.
func (p *Piano) fold(k int) int {
	for k < p.Low {
		k += octave
	}
	for k > p.High {
		k -= octave
	}
	return k
}

// levels folds every lit key into the range, keeping the brightest
// light per displayed key, and its colour.
func (p *Piano) levels() (out [MIDIKeys]float32, cols [MIDIKeys]color.RGBA) {
	for k, g := range p.glow {
		if g == 0 {
			continue
		}
		d := p.fold(k)
		if d != k {
			g *= ghost
		}
		if g > out[d] {
			out[d], cols[d] = g, p.color[k]
		}
	}
	return out, cols
}

// IsBlack tells a black key.
func IsBlack(k int) bool {
	switch k % octave {
	case 1, 3, 6, 8, 10:
		return true
	}
	return false
}

func (p *Piano) whiteCount() int {
	n := 0
	for k := p.Low; k <= p.High; k++ {
		if !IsBlack(k) {
			n++
		}
	}
	return n
}

// KeyRect places key `k`. A black key sits across the boundary of the
// white key below it.
func (p *Piano) KeyRect(k int) (x, y, w, h float32) {
	whiteW := p.W / float32(p.whiteCount())
	whites := 0
	for i := p.Low; i < k; i++ {
		if !IsBlack(i) {
			whites++
		}
	}
	if !IsBlack(k) {
		return p.X + float32(whites)*whiteW, p.Y, whiteW, p.H
	}
	bw := whiteW * blackWidth
	return p.X + float32(whites)*whiteW - bw/2, p.Y, bw, p.H * blackLength
}

// A Dress is how a game paints a key besides its light: another colour
// for its top, a name at its bottom. The zero Dress is a plain key.
type Dress struct {
	Base  color.RGBA // none when its alpha is 0
	Label string
	Font  *Font
}

// Draw paints the keyboard, whites then blacks over them. `dress`, when
// not nil, says how to paint each key; whatever sounds glows over it,
// so that playing stays readable.
func (p *Piano) Draw(c Canvas, dress func(k int) Dress) {
	levels, cols := p.levels()
	for _, black := range []bool{false, true} {
		for k := p.Low; k <= p.High; k++ {
			if IsBlack(k) != black {
				continue
			}
			x, y, w, h := p.KeyRect(k)
			base, label, lip := whiteKey, darkLabel, float32(whiteLip)
			if black {
				base, label, lip = blackKey, lightLabel, blackLip
			}
			var d Dress
			if dress != nil {
				d = dress(k)
			}
			if d.Base.A != 0 {
				base = d.Base
			}
			top := blend(base, cols[k], levels[k])
			c.Rect(x, y, w, h, keyBorder)
			if levels[k] > pressed {
				// Down: the lip goes under, the top runs to the edge,
				// darkened at the back where it sinks.
				c.Rect(x+1, y, w-2, h-1, top)
				for i := range shadowRows {
					f := shadowFrom + (1-shadowFrom)*float32(i)/shadowRows
					c.Rect(x+1, y+float32(i), w-2, 1, shade(top, f))
				}
			} else {
				c.Rect(x+1, y, w-2, h-1-lip, top)
				c.Rect(x+1, y+h-1-lip, w-2, lip, shade(top, lipShade))
			}
			if d.Label != "" && d.Font != nil {
				tw, th := c.Measure(d.Label, d.Font)
				c.Text(d.Label, d.Font, float64(x)+(float64(w)-tw)/2, float64(y+h)-th-labelGap, label)
			}
		}
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
