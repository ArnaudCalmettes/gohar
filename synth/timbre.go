package synth

import (
	"fmt"
	"math"
)

// A Timbre is the waveform an Engine plays.
//
// The sine is the zero value and stays the default: pure, and the
// plainest thing to judge an interval by. The others are the voices of
// the consoles of the eighties, a pulse at three widths, a triangle and
// a noise.
type Timbre uint8

const (
	Sine Timbre = iota

	// Pulse12 and Pulse25 are thin and nasal; Square is the hollow one,
	// half up and half down.
	Pulse12
	Pulse25
	Square

	// Triangle is soft, and in the authentic rendering climbs in
	// sixteen steps, which is where its grain comes from.
	Triangle

	// Noise is pitched by the key: the higher, the faster it changes.
	// Meant for percussion rather than for scales.
	Noise
)

var timbreNames = map[string]Timbre{
	"sine": Sine, "pulse12": Pulse12, "pulse25": Pulse25,
	"square": Square, "triangle": Triangle, "noise": Noise,
}

// ParseTimbre reads a timbre by name: sine, pulse12, pulse25, square,
// triangle or noise.
func ParseTimbre(name string) (Timbre, error) {
	t, ok := timbreNames[name]
	if !ok {
		return 0, fmt.Errorf("synth: unknown timbre %q", name)
	}
	return t, nil
}

// level is how loud a timbre is written, so that switching from one to
// another does not jump in volume. A square carries far more energy
// than a sine of the same peak.
func (t Timbre) level() float64 {
	switch t {
	case Pulse12, Pulse25, Square:
		return 0.5
	case Triangle:
		return 0.9
	case Noise:
		return 0.4
	}
	return 1
}

// duty is the part of the period a pulse spends high.
func (t Timbre) duty() float64 {
	switch t {
	case Pulse12:
		return 0.125
	case Pulse25:
		return 0.25
	}
	return 0.5
}

// oscillator is the waveform side of a voice: where it is in its
// period, how fast it moves, and the noise register.
type oscillator struct {
	phase float64 // in periods, between 0 and 1
	inc   float64 // periods per sample
	lfsr  uint16
}

// next returns the oscillator's sample and moves it on by one.
//
// # Authentic or smooth
//
// A square computed naively jumps between two samples, and the jump
// folds every harmonic above half the sample rate back into the audible
// range. That aliasing is part of how the consoles sounded, so the
// authentic rendering keeps it, along with the stepped triangle. The
// smooth rendering rounds each jump with a PolyBLEP, a two sample
// correction that costs a few multiplications, and draws the triangle
// as a line.
func (o *oscillator) next(t Timbre, smooth bool) float64 {
	p := o.phase
	var v float64

	switch t {
	case Sine:
		v = math.Sin(2 * math.Pi * p)

	case Pulse12, Pulse25, Square:
		d := t.duty()
		v = -1
		if p < d {
			v = 1
		}
		if smooth {
			v += polyBLEP(p, o.inc) - polyBLEP(frac(p-d+1), o.inc)
		}
		// A narrow pulse sits mostly low. Removing its mean keeps it
		// centred, which is what the consoles' output filter did and
		// what keeps chords from eating the headroom.
		v -= 2*d - 1

	case Triangle:
		if smooth {
			v = 4*math.Abs(p-0.5) - 1
		} else {
			step := int(p * 32)
			level := 15 - step
			if step >= 16 {
				level = step - 16
			}
			v = float64(level)/7.5 - 1
		}

	case Noise:
		if o.lfsr == 0 {
			o.lfsr = 1
		}
		v = -1
		if o.lfsr&1 == 1 {
			v = 1
		}
	}

	o.phase += o.inc
	if t == Noise {
		// The register is clocked sixteen times per period of the key,
		// fast enough to sound as a hiss and still follow the key.
		o.phase += 15 * o.inc
	}
	if o.phase >= 1 {
		o.phase -= math.Floor(o.phase)
		if t == Noise {
			bit := (o.lfsr ^ o.lfsr>>1) & 1
			o.lfsr = o.lfsr>>1 | bit<<14
		}
	}
	return v * t.level()
}

// polyBLEP is the correction for a jump of two at phase zero, `t` being
// the phase and `dt` the increment per sample. Zero away from the jump.
func polyBLEP(t, dt float64) float64 {
	switch {
	case t < dt:
		t /= dt
		return t + t - t*t - 1
	case t > 1-dt:
		t = (t - 1) / dt
		return t*t + t + t + 1
	}
	return 0
}

func frac(x float64) float64 {
	return x - math.Floor(x)
}
