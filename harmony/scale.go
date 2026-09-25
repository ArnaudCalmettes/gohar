package harmony

import (
	"fmt"
	"iter"
	"math/bits"
)

// A Degree is the ordinal rank of a member within a scale pattern,
// counted from 1.
//
// Its upper bound is the size of the pattern, not seven. A pentatonic
// pattern has degrees 1 through 5. Seven is only the common case.
//
// A Degree is not a chord figure. The 9 in a ninth chord is an
// Extension, declared in chord.go, and the two must not be mixed: they
// count different things over different ranges.
type Degree uint8

// A ScalePattern is the structure of a scale, independent of any
// tonic. Bit n is set when the scale contains a note n semitones above
// its tonic.
//
// # Invariants
//
// Bits 12 through 15 are always zero, and bit 0 is always set: a scale
// contains its own tonic. Both hold for every value the constructors
// and operations below can produce.
//
// # Why there is no Transpose
//
// A pattern is relative, so moving it by a distance means nothing. The
// operation that does mean something is [ScalePattern.Mode], which
// re-anchors the pattern on one of its own members. Anchoring it on an
// absolute class is [ScalePattern.At], which yields a [PitchSet] and
// leaves pattern space entirely.
type ScalePattern uint16

// The usual patterns.
//
// Unlike interval distances, these names carry no spelling. A major
// scale is one specific arrangement of semitones in every context, with
// no second reading to choose between, so naming it here smuggles
// nothing into the core.
//
// # How to read the literals
//
// The underscores split the octave at the tritone, into its two
// tetrachords rather than into groups of four bits. The group on the
// right runs from the tonic to the fourth, the one on the left from the
// tritone to the seventh. Read right to left within a group, as bits
// are numbered. Only the twelve low bits ever carry anything, so the
// four above them are left unwritten.
//
// Grouped this way the scales show how they are built. The three minors
// share their lower tetrachord, and the melodic minor differs from the
// major in that tetrachord alone: it is the major with a lowered third,
// which is exactly what its name as a mode, ionian flat 3, says.
//
// Five of these are the mother scales the modes are drawn from, and a
// [System] designates each of them. The natural minor is not one: it is
// the aeolian, the sixth mode of the major, kept here because it is
// asked for often enough to deserve a name.
const (
	ScaleMajor               ScalePattern = 0b101010_110101
	ScaleNaturalMinor        ScalePattern = 0b010110_101101
	ScaleHarmonicMinor       ScalePattern = 0b100110_101101
	ScaleMelodicMinor        ScalePattern = 0b101010_101101
	ScaleHarmonicMajor       ScalePattern = 0b100110_110101
	ScaleDoubleHarmonicMajor ScalePattern = 0b100110_110011
)

// NewScalePattern builds the pattern that s describes when read from
// tonic.
//
// Fails when tonic does not belong to s, since bit 0 would come out
// clear and break the invariant. That failure is worth surfacing: it
// means the caller proposed a tonic the set does not contain, which in
// the analysis engine is a candidate to reject rather than a value to
// repair.
func NewScalePattern(s PitchSet, tonic PitchClass) (ScalePattern, error) {
	if !tonic.IsValid() {
		return 0, fmt.Errorf("gohar: tonic %d is out of range", tonic)
	}
	if !s.Contains(tonic) {
		return 0, fmt.Errorf("gohar: tonic %d does not belong to the set", tonic)
	}
	return ScalePattern(s.Transpose(-Semitones(tonic))), nil
}

// At anchors p on tonic and returns the resulting set of classes.
//
// This is the conversion out of pattern space. [NewScalePattern] is the
// conversion back in.
func (p ScalePattern) At(tonic PitchClass) PitchSet {
	return PitchSet(p).Transpose(Semitones(tonic))
}

// Len returns the number of notes in p.
func (p ScalePattern) Len() int {
	return bits.OnesCount16(uint16(p))
}

// IsHeptatonic reports whether p has exactly seven notes.
//
// This is the precondition of [Tonality] and of spelling by degree: one
// letter per degree only works when there are seven of them.
func (p ScalePattern) IsHeptatonic() bool {
	return p.Len() == 7
}

// Contains reports whether p includes a note n semitones above its
// tonic.
func (p ScalePattern) Contains(n Semitones) bool {
	if n < 0 || n >= 12 {
		return false
	}
	return p&(1<<n) != 0
}

// Offsets iterates over the members of p, yielding each degree with its
// distance above the tonic. The first pair is always (1, 0).
func (p ScalePattern) Offsets() iter.Seq2[Degree, Semitones] {
	return func(yield func(Degree, Semitones) bool) {
		d := Degree(1)
		for n := Semitones(0); n < 12; n++ {
			if !p.Contains(n) {
				continue
			}
			if !yield(d, n) {
				return
			}
			d++
		}
	}
}

// Offset returns the distance from the tonic to degree d, and whether d
// exists in p.
func (p ScalePattern) Offset(d Degree) (Semitones, bool) {
	for degree, n := range p.Offsets() {
		if degree == d {
			return n, true
		}
	}
	return 0, false
}

// Mode returns the pattern read from degree d of p.
//
// Mode(2) of [ScaleMajor] is the dorian pattern. Mode(1) returns p
// unchanged. The result is a rotation that re-anchors bit 0 on the
// member at degree d, so the invariants hold by construction.
//
// Reports whether d exists in p.
func (p ScalePattern) Mode(d Degree) (ScalePattern, bool) {
	n, ok := p.Offset(d)
	if !ok {
		return 0, false
	}
	return ScalePattern(PitchSet(p).Transpose(-n)), true
}

// From iterates over the members of p, starting at root and moving up.
//
// One method covers every target: From on a [PitchClass] yields classes
// folded into one octave, From on a [Pitch] yields ascending pitches
// that cross octave boundaries. The previous iteration of this library
// carried six near identical methods here for want of generic methods,
// which arrived in Go 1.27.
//
// The iteration follows degree order, so the caller can pair it with
// [ScalePattern.Offsets] to recover which degree each element holds.
// That pairing is what spelling by degree needs.
func (p ScalePattern) From[T Transposable[T]](root T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := Semitones(0); n < 12; n++ {
			if p.Contains(n) && !yield(root.Transpose(n)) {
				return
			}
		}
	}
}

// Mirror returns the pattern read backwards from its tonic: every
// offset n becomes twelve minus n.
//
// This is the retrograde of the inverse that Bernard Maury's teaching
// uses to reharmonise, and it keeps the tonic in place. The ionian
// mirrors onto the phrygian, the lydian onto the locrian, and the
// dorian onto itself.
//
// # What this is not
//
// A naming rule. Mode names follow first how a mode sounds and what it
// does in a tonal context, and a mirror that sends a mode onto itself
// means only that its inverse sounds like it. The names happen to
// commute with the mirror across most of the catalogue and stop doing
// so at the palindromes, which is where the ear takes over from the
// symmetry. Nothing in this library derives a name from a mirror.
//
// It stays a structural fact worth checking, though: the thirty five
// modes form a set that the mirror maps onto itself.
func (p ScalePattern) Mirror() ScalePattern {
	var out ScalePattern
	for n := Semitones(0); n < 12; n++ {
		if p.Contains(n) {
			out |= 1 << ((12 - n) % 12)
		}
	}
	return out
}

// String returns the twelve bits of p, for debugging only.
func (p ScalePattern) String() string {
	return fmt.Sprintf("%012b", uint16(p))
}

// A Scale is a pattern anchored on a tonic.
//
// The zero value is not a usable scale: its pattern is empty, which
// violates the bit 0 invariant. Build one with [NewScale] or by taking
// a declared pattern. A Scale is not a [Tonality]: any pattern makes a
// Scale, only a heptatonic one makes a Tonality.
type Scale struct {
	Tonic   PitchClass
	Pattern ScalePattern
}

// NewScale pairs a tonic with a pattern.
func NewScale(tonic PitchClass, p ScalePattern) (Scale, error) {
	if !tonic.IsValid() {
		return Scale{}, fmt.Errorf("gohar: tonic %d is out of range", tonic)
	}
	if !p.Contains(0) {
		return Scale{}, fmt.Errorf("gohar: scale pattern does not contain its own tonic")
	}
	return Scale{Tonic: tonic, Pattern: p}, nil
}

// Set returns the classes of s.
func (s Scale) Set() PitchSet {
	return s.Pattern.At(s.Tonic)
}

// Classes iterates over the classes of s from its tonic upward.
func (s Scale) Classes() iter.Seq[PitchClass] {
	return s.Pattern.From(s.Tonic)
}

// Pitches iterates over the pitches of s between from and to
// inclusive, ascending.
//
// Bounds are pitches rather than octave numbers so that a keyboard with
// an odd range needs no arithmetic at the call site.
func (s Scale) Pitches(from, to Pitch) iter.Seq[Pitch] {
	return func(yield func(Pitch) bool) {
		set := s.Set()
		for p := from; p <= to; p++ {
			if set.Contains(p.Class()) && !yield(p) {
				return
			}
		}
	}
}

// Degree returns the class at degree d, and whether d exists in s.
func (s Scale) Degree(d Degree) (PitchClass, bool) {
	n, ok := s.Pattern.Offset(d)
	if !ok {
		return 0, false
	}
	return s.Tonic.Transpose(n), true
}
