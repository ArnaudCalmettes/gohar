package harmony

import "strconv"

// A PitchClass is a height modulo the octave: 0 through 11.
//
// A PitchClass carries no spelling. C sharp and D flat are the same
// value, 1, and compare equal. This is the point: a player at the
// keyboard produces heights, not spellings, and enharmonic equality is
// the equality the analysis engine needs.
//
// The written form is a naming.SpelledNote, built from a PitchClass
// and a tonal context.
type PitchClass uint8

// PitchClassCount is the number of distinct pitch classes.
const PitchClassCount = 12

// The core declares no named pitch class constants.
//
// PitchClassC would be a name, and the core carries none. Write
// PitchClass(0). To go the other way, from text to value, use the
// parser in the naming package, which has the tables and the context
// the core lacks.

// IsValid reports whether c is below [PitchClassCount].
//
// PitchClass is an unsigned byte, so out of range values are
// representable. The type does not enforce the invariant, for the same
// reason the library has no error returning constructor here: this is a
// value type used inside tight loops, and an allocation or a branch per
// construction is not worth the guarantee. Operations that can produce
// an out of range value are documented as such; the ones below cannot.
func (c PitchClass) IsValid() bool {
	return c < PitchClassCount
}

// Transpose moves c by s semitones, wrapping around the octave.
//
// Total and closed: every input yields a valid PitchClass, negative
// distances included. Transpose(Semitones(-1)) on class 0 is 11.
func (c PitchClass) Transpose(s Semitones) PitchClass {
	return PitchClass(((int(c)+int(s))%12 + 12) % 12)
}

// Up returns the ascending distance from c to other, in the range
// [0, 12).
//
// Not symmetric. Up from 0 to 11 is 11, while Up from 11 to 0 is 1.
// There is no signed variant, because the sign of a distance between
// two classes has no meaning without an octave to place them in.
func (c PitchClass) Up(other PitchClass) Semitones {
	return Semitones(((int(other)-int(c))%12 + 12) % 12)
}

// Pitch returns the [Pitch] of class c in the given octave, with octave
// 4 holding [MiddleC].
//
// The result may fall outside the MIDI range for extreme octaves. Check
// with [Pitch.IsValid] where it matters.
func (c PitchClass) Pitch(octave int8) Pitch {
	return Pitch(int(c) + 12*(int(octave)+1))
}

// Set returns the singleton [PitchSet] containing c.
func (c PitchClass) Set() PitchSet {
	return PitchSet(1) << c
}

// String returns the numeric value of c, for debugging only.
//
// It never returns a note name. Anything shown to a user goes through
// the naming package, which needs a tonal context to decide between
// spellings. A String method that guessed one would be a naming
// decision taken in the wrong layer, and it would be silently wrong
// most of the time.
func (c PitchClass) String() string {
	return strconv.Itoa(int(c))
}
