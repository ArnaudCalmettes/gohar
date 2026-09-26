package harmony

// A Pitch is an absolute sound height, identified with its MIDI note
// number. Middle C is 60.
//
// A Pitch is a position, never a distance. The distance between two
// pitches is a [Semitones]. Keeping the two apart is deliberate: a
// single type for both invites shifting a bitmask by an absolute
// height, which panics for negative values and silently overflows past
// 31. With two types the mistake does not compile.
//
// Nothing in this package accepts a Pitch where a distance is meant.
// Set operations take [Semitones] or [PitchClass] only.
type Pitch int8

// The valid MIDI range. Pitch is signed and negative values are
// representable, so the range is an invariant to check, not one the
// type enforces.
const (
	PitchMin Pitch = 0
	PitchMax Pitch = 127
)

// MiddleC anchors the MIDI numbering.
//
// This is the one place where the core names a pitch. MIDI note
// numbering is a fact of the protocol rather than a naming convention
// of this library, so the exception stops here: there is no PitchD,
// no PitchEFlat, and no table of note names anywhere in this package.
const MiddleC Pitch = 60

// IsValid reports whether `p` falls within the MIDI range.
func (p Pitch) IsValid() bool {
	return p >= PitchMin && p <= PitchMax
}

// Class returns the pitch class of `p`, discarding its octave.
//
// Defined for every Pitch, including those outside the MIDI range.
func (p Pitch) Class() PitchClass {
	return PitchClass(((int(p) % 12) + 12) % 12)
}

// Octave returns the octave number of `p`, with [MiddleC] in octave 4.
//
// This follows scientific pitch notation, where MIDI 60 is C4. Some
// hardware numbers the same pitch C3. The library commits to C4 and
// leaves the offset to whatever adapter talks to the device.
func (p Pitch) Octave() int8 {
	oct := int(p) / 12
	if int(p)%12 < 0 {
		oct--
	}
	return int8(oct - 1)
}

// Transpose moves `p` by `s` semitones.
//
// The result may fall outside the MIDI range. Transpose does not clamp
// and does not report an error: an out of range intermediate value is
// legitimate during a computation. Check with [Pitch.IsValid] at the
// point where the value reaches a device or a display.
func (p Pitch) Transpose(s Semitones) Pitch {
	return p + Pitch(s)
}

// Sub returns the signed distance from `other` to `p`.
//
// Positive when `p` is above `other`. The result can exceed one octave.
func (p Pitch) Sub(other Pitch) Semitones {
	return Semitones(p - other)
}

// Semitones is a signed distance between two heights, counted in
// semitones. It is never an absolute position.
type Semitones int8

const (
	// Unison is the null distance.
	Unison Semitones = 0
	// Semitone is the smallest distance the library represents.
	Semitone Semitones = 1
	// SemitonesPerOctave is the size of an octave.
	//
	// Named for its unit rather than simply Octave, so that it does not
	// read like a counterpart of [Pitch.Octave], which returns an
	// octave number and not a distance.
	SemitonesPerOctave Semitones = 12
)

// The core deliberately declares no constants for the usual interval
// distances.
//
// Every such name is a spelling. Three semitones is a minor third in
// one context and an augmented second in another; six semitones is an
// augmented fourth or a diminished fifth. Choosing between those names
// requires a tonal context, which is exactly what the core does not
// have. Naming them here would smuggle spelling back into the layer
// that was built to be free of it.
//
// Write Semitones(7) in the core. The readable constants belong in the
// naming package, where a context is available to justify them.

// Abs returns the absolute value of `s`.
func (s Semitones) Abs() Semitones {
	if s < 0 {
		return -s
	}
	return s
}

// Fold reduces `s` to the pitch class it reaches from class 0.
//
// Defined for negative distances: Fold(-1) is 11, not an error.
func (s Semitones) Fold() PitchClass {
	return PitchClass(((int(s) % 12) + 12) % 12)
}

// Transposable constrains the types that a distance moves, yielding the
// same type. [Pitch], [PitchClass] and [PitchSet] all satisfy it.
//
// It exists so that [ScalePattern] and [ChordPattern] expose one
// generic method each instead of one method per target type. The
// previous iteration of the library carried six near identical methods
// on ScalePattern for want of this. Generic methods require Go 1.27.
type Transposable[T any] interface {
	Transpose(Semitones) T
}
