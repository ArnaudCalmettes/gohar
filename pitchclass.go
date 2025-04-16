package gohar

import (
	"fmt"
	"iter"
)

// A PitchClass can be thought of as an "absolute" or "ideal" note; it represents
// all octaves of a note. In other words, it represents the class of a note's pitches.
//
// This is the kind of construct we manipulate intuitively when we say things like:
// - Eb belongs to the Cm11 chord.
// - F# belongs to the C lydian mode.
//
// This type is encoded as a single byte.
//
//	bit number : 7654 3 210
//	values:      AAAA X BBB
//	A: Alteration (4 bits) encoded as 8 + A where -2 (bb) <= A <= +2 (##).
//	B: Base Pitch Class as a number between 0 (C) and 6 (B).
//	X: unused
//
// Examples:
//
//	C (natural)
//	PitchClass: 0b10000000 (= 128)
//	fields:       AAAA BBB
//	A: 0b1000 (= 8 + 0) -> no alteration
//	B: 0b000 (= 0) -> C
//
//	F#
//	PitchClass: 0b10010011 (= 115)
//	fields:       AAAA BBB
//	A: 0b1001 (= 8 + 1) -> 1 sharp
//	B: 0b011 (= 3) -> F
type PitchClass uint8

const (
	PitchClassC PitchClass = 0 | PitchClassNatural
	PitchClassD PitchClass = 1 | PitchClassNatural
	PitchClassE PitchClass = 2 | PitchClassNatural
	PitchClassF PitchClass = 3 | PitchClassNatural
	PitchClassG PitchClass = 4 | PitchClassNatural
	PitchClassA PitchClass = 5 | PitchClassNatural
	PitchClassB PitchClass = 6 | PitchClassNatural

	PitchClassNatural PitchClass = (8 + 0) << 4
)

// NewPitchClassFromChar creates a new PitchClass.
// ErrInvalidPitchClass is returned if char isn't in the range ['A','G'].
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
	return defaultPitchClasses[p.Normalize()]
}

var defaultPitchClasses = [12]PitchClass{
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

// BaseName returns the base pitch class as an ASCII byte in the ['A','G'] range.
func (p PitchClass) BaseName() byte {
	return "CDEFGAB"[int(p&0x0f)]
}

// Base returns the base as a number between 0 (C) and 6 (B)
func (p PitchClass) Base() int8 {
	return int8(p & 0xf)
}

// Alt returns the alteration of this PitchClass.
func (p PitchClass) Alt() Pitch {
	return Pitch((p&0xf0)>>4) - 8
}

// IsValid returns true if this pitch class has:
// - a base in the range [1,6],
// - an alt in the range [-2,+2].
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
	return p.WithAlt(+1)
}

// Flat adds one flat to the PitchClass.
func (p PitchClass) Flat() PitchClass {
	return p.WithAlt(-1)
}

// DoubleSharp adds a double sharp to the PitchClass.
func (p PitchClass) DoubleSharp() PitchClass {
	return p.WithAlt(+2)
}

// DoubleFlat adds a double flat to the PitchClass.
func (p PitchClass) DoubleFlat() PitchClass {
	return p.WithAlt(-2)
}

// WithAlt applies given alteration to the PitchClass.
func (p PitchClass) WithAlt(alt Pitch) PitchClass {
	return PitchClass(alt+8)<<4 | (p & 0x0f)
}

// Pitch returns the pitch of the PitchClass at given octave.
func (p PitchClass) Pitch(oct int8) Pitch {
	return (asPitch[p.Base()] + p.Alt()).AtOctave(oct)
}

var asPitch = [7]Pitch{0, 2, 4, 5, 7, 9, 11}

// Pitches iterates over all pitches of this class in ascending order
// within the specified range.
func (p PitchClass) Pitches(from, to Pitch) iter.Seq[Pitch] {
	return func(yield func(Pitch) bool) {
		if from > to {
			return
		}
		start := p.Pitch(from.GetOctave())
		for start < from {
			start += PitchDiffOctave
		}
		for pitch := start; pitch <= to; pitch += PitchDiffOctave {
			if !yield(pitch) {
				return
			}
		}
	}
}

// Transpose transposes a pitch class by given interval.
func (p PitchClass) Transpose(i Interval) PitchClass {
	return p.MoveBase(i.ScaleDiff).ClipToPitch(p.Pitch(0) + i.PitchDiff)
}

// MoveBase moves the base pitch class by given (positive or negative) steps.
func (p PitchClass) MoveBase(step int8) PitchClass {
	b := (int8(p.Base()) + step) % 7
	if b < 0 {
		b += 7
	}
	return PitchClass(b) | p&0xf0
}

// ClipToPitch changes p's alteration so that the target pitch belongs
// to the resulting PitchClass.
func (p PitchClass) ClipToPitch(target Pitch) PitchClass {
	alt := target.Normalize() - p.WithAlt(0).Pitch(0)
	switch {
	case alt < -6:
		alt += 12
	case alt > 6:
		alt -= 12
	}
	return p.WithAlt(alt)
}

func PitchesWithClasses(from, to Pitch, pcs []PitchClass) iter.Seq2[Pitch, PitchClass] {
	if len(pcs) == 0 {
		return func(yield func(Pitch, PitchClass) bool) {}
	}
	return func(yield func(Pitch, PitchClass) bool) {
		pitch := from - 1
		oct := pitch.GetOctave()
		for pitch < to {
			for _, pc := range pcs {
				p := pc.Pitch(oct)
				for p <= pitch {
					oct++
					p += 12
				}
				if p > to {
					return
				}
				if !yield(p, pc) {
					return
				}
				pitch = p
			}
		}
	}
}
