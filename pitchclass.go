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
// This type is maps to a 16-bit structure that can be efficiently passed around by
// immutable copies.
type PitchClass struct {
	// base is the base pitch class, represented as a byte between 'A' and 'G'.
	base byte
	// Alt is the alteration that is applied to the Base class.
	// Its value should be kept between -2 (double flat) and +2 (double sharp).
	alt Pitch
}

var (
	PitchClassA = PitchClass{'A', 0}
	PitchClassB = PitchClass{'B', 0}
	PitchClassC = PitchClass{'C', 0}
	PitchClassD = PitchClass{'D', 0}
	PitchClassE = PitchClass{'E', 0}
	PitchClassF = PitchClass{'F', 0}
	PitchClassG = PitchClass{'G', 0}
)

// NewPitchClass creates a new PitchClass.
// ErrInvalidPitchClass is returned if b isn't in the range [A-G].
// ErrInvalidAlteration is returned if alt isn't in the range [-2,2].
func NewPitchClass(b byte, alt Pitch) (PitchClass, error) {
	if b < 'A' || b > 'G' {
		return PitchClass{}, wrapErrorf(ErrInvalidPitchClass,
			"expected range 'A' <= b <= 'G', got %c", b,
		)
	}
	if alt < -2 || alt > 2 {
		return PitchClass{}, wrapErrorf(ErrInvalidAlteration,
			"expected range -2 <= alt <= +2, got %d", alt,
		)
	}
	return PitchClass{b, alt}, nil
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
	return fmt.Sprintf("%c%s", p.base, altToString(p.alt))
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

// Base returns the base pitch class.
func (p PitchClass) Base() byte {
	return p.base
}

// Alt returns the alteration of this PitchClass.
func (p PitchClass) Alt() Pitch {
	return p.alt
}

// IsValid returns true if this pitch class has:
// - a base in the range [A-G],
// - an alt in the range [-2;+2].
func (p PitchClass) IsValid() bool {
	return p.base >= 'A' && p.base <= 'G' && p.alt >= -2 && p.alt <= 2
}

// IsEnharmonic returns true if both PitchClasses map to the same
// pitches.
func (p PitchClass) IsEnharmonic(o PitchClass) bool {
	return p.Pitch(0) == o.Pitch(0)
}

// Sharp adds one sharp to the PitchClass.
func (p PitchClass) Sharp() PitchClass {
	return PitchClass{p.base, p.alt + 1}
}

// Flat adds one flat to the PitchClass.
func (p PitchClass) Flat() PitchClass {
	return PitchClass{p.base, p.alt - 1}
}

// DoubleSharp adds a double sharp to the PitchClass.
func (p PitchClass) DoubleSharp() PitchClass {
	return PitchClass{p.base, p.alt + 2}
}

// DoubleFlat adds a double flat to the PitchClass.
func (p PitchClass) DoubleFlat() PitchClass {
	return PitchClass{p.base, p.alt - 2}
}

// Pitch returns the pitch of the PitchClass at given octave.
func (p PitchClass) Pitch(oct int8) Pitch {
	return (asPitch[int(p.base-'A')] + p.alt).AtOctave(oct)
}

var asPitch = [7]Pitch{9, 11, 0, 2, 4, 5, 7}

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
		for pitch := start; pitch <= to; pitch += 12 {
			if !yield(pitch) {
				return
			}
		}
	}
}

// Transpose transposes a pitch class by given interval.
func (p PitchClass) Transpose(i Interval) PitchClass {
	return pitchClassWithPitch(
		moveBaseNote(p.base, int(i.ScaleDiff)),
		p.Pitch(0)+i.PitchDiff,
	)
}

func pitchClassWithPitch(b byte, target Pitch) PitchClass {
	pc := PitchClass{b, 0}
	pc.alt = target.Normalize() - pc.Pitch(0)
	if pc.alt < -6 {
		pc.alt += 12
	}
	if pc.alt > 6 {
		pc.alt -= 12
	}
	return pc
}

func (p PitchClass) Serialize() uint16 {
	return uint16(p.base) | uint16(p.alt+64)<<8
}

func DeserializePitchClass(repr uint16) PitchClass {
	return PitchClass{
		base: byte(repr & 0xff),
		alt:  Pitch((repr & 0xff00) >> 8),
	}
}
