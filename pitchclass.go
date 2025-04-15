package gohar

import (
	"fmt"
	"iter"
)

// A PitchClass can be thought of as an "absolute" or "ideal" note, that we can manifest
// as a distinct music note by giving it an octave and a duration.
// This is the kind of construct we manipulate intuitively when we reason about whether a
// note belongs to a chord or scale: we know E belongs to the C major scale, regardless
// of the octave.
//
// This type is encoded as a single byte.
//
//	bit number : 76543210
//	values:      AAAAXBBB
//	A: Alterations (4 bits)
//	B: Base Pitch Class (3 bits)
//	X: unused
type PitchClass uint8

const (
	PitchClassC PitchClass = 0 | PitchClassNatural
	PitchClassD PitchClass = 1 | PitchClassNatural
	PitchClassE PitchClass = 2 | PitchClassNatural
	PitchClassF PitchClass = 3 | PitchClassNatural
	PitchClassG PitchClass = 4 | PitchClassNatural
	PitchClassA PitchClass = 5 | PitchClassNatural
	PitchClassB PitchClass = 6 | PitchClassNatural

	PitchClassNatural     PitchClass = (8 + 0) << 4
	PitchClassSharp       PitchClass = (8 + 1) << 4
	PitchClassFlat        PitchClass = (8 - 1) << 4
	PitchClassDoubleSharp PitchClass = (8 + 2) << 4
	PitchClassDoubleFlat  PitchClass = (8 - 2) << 4
)

// NewPitchClassFromChar creates a new PitchClass.
// ErrInvalidPitchClass is returned if b isn't in the range [A-G].
// ErrInvalidAlteration is returned if alt isn't in the range [-2,2].
func NewPitchClassFromChar(char byte, alt Pitch) (PitchClass, error) {
	if char < 'A' || char > 'G' {
		return 0, wrapErrorf(ErrInvalidPitchClass,
			"expected range 'A' <= b <= 'G', got %c", char,
		)
	}
	if alt < -2 || alt > 2 {
		return 0, wrapErrorf(ErrInvalidAlteration,
			"expected range -2 <= alt <= +2, got %d", alt,
		)
	}
	if char < 'C' {
		char += 7
	}
	pc := PitchClass(char - 'C')
	pc |= PitchClass(alt+8) << 4
	return pc, nil
}

// DefaultPitchClass returns the default PitchClass associated with
// given pitch. On pitches that do not map to a base pitch class ("black keys"),
// it picks the base pitch class immediately above and flattens it.
func DefaultPitchClass(p Pitch) PitchClass {
	return defaultPitchClass[p.Normalize()]
}

var defaultPitchClass = [12]PitchClass{
	PitchClassC,
	PitchClassD.Flat(),
	PitchClassD,
	PitchClassE.Flat(),
	PitchClassE,
	PitchClassF,
	PitchClassG.Flat(),
	PitchClassG,
	PitchClassA.Flat(),
	PitchClassA,
	PitchClassB.Flat(),
	PitchClassB,
}

// String returns a string representation of the PitchClass.
// "<invalid>" is returned if the PitchClass is invalid.
func (p PitchClass) String() string {
	if !p.IsValid() {
		return "<invalid>"
	}
	return fmt.Sprintf("%c%s", p.BaseName(), altToString(p.Alt()))
}

func altToString(alt Pitch) string {
	switch alt {
	case -2:
		return AltDoubleFlat
	case -1:
		return AltFlat
	case 0:
		return ""
	case 1:
		return AltSharp
	case 2:
		return AltDoubleSharp
	default:
		return fmt.Sprintf("(%+d)", alt)
	}
}

// BaseName returns the base pitch class.
func (p PitchClass) BaseName() byte {
	return "CDEFGAB"[int(p&0x0f)]
}

func (p PitchClass) Base() int8 {
	return int8(p & 0xf)
}

// Alt returns the alteration of this PitchClass.
func (p PitchClass) Alt() Pitch {
	return Pitch((p&0xf0)>>4) - 8
}

// IsValid returns true if this pitch class has:
// - a base in the range [A-G],
// - an alt in the range [-2;+2].
func (p PitchClass) IsValid() bool {
	alt := p.Alt()
	return p.Base() < 7 && -2 <= alt && alt <= 2
}

// IsEnharmonic returns true if both PitchClasses map to the same
// pitches.
func (p PitchClass) IsEnharmonic(o PitchClass) bool {
	return p.Pitch(0) == o.Pitch(0)
}

// Sharp adds one sharp to the PitchClass.
func (p PitchClass) Sharp() PitchClass {
	return p.WithAlt(PitchClassSharp)
}

// Flat adds one flat to the PitchClass.
func (p PitchClass) Flat() PitchClass {
	return p.WithAlt(PitchClassFlat)
}

// DoubleSharp adds a double sharp to the PitchClass.
func (p PitchClass) DoubleSharp() PitchClass {
	return p.WithAlt(PitchClassDoubleSharp)
}

// DoubleFlat adds a double flat to the PitchClass.
func (p PitchClass) DoubleFlat() PitchClass {
	return p.WithAlt(PitchClassDoubleFlat)
}

// WithAlt applies given alteration to the PitchClass.
func (p PitchClass) WithAlt(alt PitchClass) PitchClass {
	return p&0x0f | alt&0xf0
}

// Pitch returns the pitch of the PitchClass at given octave.
func (p PitchClass) Pitch(oct int8) Pitch {
	return (asPitch[p.Base()] + p.Alt()).AtOctave(oct)
}

var asPitch = [7]Pitch{0, 2, 4, 5, 7, 9, 11}

// Pitches iterates over all pitches of the PitchClass within the
// specified (inclusive) Pitch range.
func (p PitchClass) Pitches(from, to Pitch) iter.Seq[Pitch] {
	if from > to {
		return func(func(Pitch) bool) {}
	}
	start := p.Pitch(from.GetOctave())
	for start < from {
		start += PitchDiffOctave
	}
	return func(yield func(Pitch) bool) {
		for pitch := start; pitch <= to; pitch += PitchDiffOctave {
			if !yield(pitch) {
				return
			}
		}
	}
}

// Transpose transposes a pitch class by given interval.
func (p PitchClass) Transpose(i Interval) PitchClass {
	return pitchClassWithPitch(
		wrapBasePitchClass(int8(p.Base())+i.ScaleDiff, 7),
		p.Pitch(0)+i.PitchDiff,
	)
}
func wrapBasePitchClass(value, max int8) PitchClass {
	value %= max
	for value < 0 {
		value += max
	}
	return PitchClass(value) | PitchClassNatural
}

func pitchClassWithPitch(base PitchClass, target Pitch) PitchClass {
	alt := target.Normalize() - base.Pitch(0)
	for alt < -6 {
		alt += 12
	}
	for alt > 6 {
		alt -= 12
	}
	return base.WithAlt(PitchClass(8+alt) << 4)
}
