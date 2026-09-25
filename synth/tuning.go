package synth

import "math"

// A Tuning turns a MIDI key number into a frequency.
//
// # Why this lives here and not in harmony
//
// A harmony.Pitch exists precisely to say nothing about frequency. It
// is a position in a system of intervals, and the same position is a
// different number of hertz at a different diapason, in a different
// temperament, on a different instrument. Giving it a pitch standard
// would undo the abstraction it was built for.
//
// So the conversion belongs to whoever makes sound, and it takes a
// plain key number rather than a harmony.Pitch, which also keeps this
// package free of any dependency on the theory layer.
type Tuning struct {
	// A4 is the frequency of key 69 in hertz. Zero means 440.
	A4 float64
}

// Frequency returns the pitch of a MIDI key number in hertz.
//
// Equal temperament, which is what a keyboard gives us and what every
// game here assumes. A tuning that bends the steps would need a table
// rather than an exponent, and nothing calls for one yet.
func (t Tuning) Frequency(key int) float64 {
	a4 := t.A4
	if a4 == 0 {
		a4 = 440
	}
	return a4 * math.Exp2(float64(key-69)/12)
}
