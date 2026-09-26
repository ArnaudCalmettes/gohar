package harmony

import (
	"fmt"
	"iter"
	"math/bits"
)

// A PitchSet is a set of pitch classes, held as a twelve bit mask. Bit
// n is set when class n belongs to the set.
//
// This is the type the analysis engine works on. What reaches it from a
// keyboard is a handful of held keys with no octave, no order and no
// spelling, which is exactly a PitchSet.
//
// # Invariant
//
// Bits 12 through 15 are always zero. Every operation below preserves
// this, and every operation added later must preserve it too.
//
// The invariant is what makes == the correct comparison for this type.
// Two sets with the same members but different high bits would compare
// unequal while being logically identical, and the bug would surface as
// a chord that is recognised on one code path and not on another. There
// is no Equal method: use ==, and keep the invariant true.
type PitchSet uint16

// EmptyPitchSet contains no class. It is the zero value, so a declared
// PitchSet is usable without construction.
const EmptyPitchSet PitchSet = 0

// ChromaticPitchSet contains all twelve classes.
const ChromaticPitchSet PitchSet = 0x0fff

// NewPitchSet builds a set from the given classes.
//
// Duplicates are harmless. Classes outside [0, 12) are rejected rather
// than folded: folding would silently accept a caller who confused a
// distance with a class, which is the confusion the type system is
// meant to catch here.
func NewPitchSet(classes ...PitchClass) (PitchSet, error) {
	var s PitchSet
	for _, c := range classes {
		if !c.IsValid() {
			return 0, fmt.Errorf("gohar: pitch class %d is out of range", c)
		}
		s |= 1 << c
	}
	return s, nil
}

// Contains reports whether `c` belongs to `s`.
func (s PitchSet) Contains(c PitchClass) bool {
	return s&(1<<c) != 0
}

// With returns `s` with `c` added. Adding a member already present returns
// `s` unchanged.
func (s PitchSet) With(c PitchClass) PitchSet {
	return s | 1<<c
}

// Without returns `s` with `c` removed. Removing an absent member returns `s`
// unchanged.
func (s PitchSet) Without(c PitchClass) PitchSet {
	return s &^ (1 << c)
}

// Union returns the classes belonging to `s` or `other`.
func (s PitchSet) Union(other PitchSet) PitchSet {
	return s | other
}

// Intersect returns the classes belonging to both `s` and `other`.
func (s PitchSet) Intersect(other PitchSet) PitchSet {
	return s & other
}

// Difference returns the classes of `s` absent from `other`.
func (s PitchSet) Difference(other PitchSet) PitchSet {
	return s &^ other
}

// IsSubsetOf reports whether every class of `s` belongs to `other`.
func (s PitchSet) IsSubsetOf(other PitchSet) bool {
	return s&other == s
}

// Len returns the number of classes in `s`.
func (s PitchSet) Len() int {
	return bits.OnesCount16(uint16(s))
}

// IsEmpty reports whether `s` contains no class.
func (s PitchSet) IsEmpty() bool {
	return s == 0
}

// Transpose moves every class of `s` by `n` semitones.
//
// This is a rotation over twelve bits, not a shift: a class leaving one
// end re-enters at the other. A plain shift would drop members and the
// loss would be silent. Negative distances rotate the other way.
//
// Transpose makes PitchSet satisfy [Transposable].
func (s PitchSet) Transpose(n Semitones) PitchSet {
	shift := ((int(n) % 12) + 12) % 12
	if shift == 0 {
		return s
	}
	return ((s << shift) | (s >> (12 - shift))) & ChromaticPitchSet
}

// Classes iterates over the members of `s` in ascending numeric order.
//
// Numeric order is not musical order. The set holds no root, so there
// is no note to start from; a caller who wants an ordering relative to
// a root transposes first.
func (s PitchSet) Classes() iter.Seq[PitchClass] {
	return func(yield func(PitchClass) bool) {
		for c := PitchClass(0); c < PitchClassCount; c++ {
			if s.Contains(c) && !yield(c) {
				return
			}
		}
	}
}

// Rotations iterates over the twelve transpositions of `s`, yielding the
// distance applied along with the result. The first pair is (0, `s`).
//
// This is the brute force half of root finding: rotate the played set
// until it lines up with a pattern anchored at class 0, and the
// distance that made it line up is the root. Twelve iterations of a
// register sized operation, so the loop is cheap. See [PitchSet.Canonical]
// for the version that trades the loop for a lookup.
func (s PitchSet) Rotations() iter.Seq2[Semitones, PitchSet] {
	return func(yield func(Semitones, PitchSet) bool) {
		for n := Semitones(0); n < 12; n++ {
			if !yield(n, s.Transpose(n)) {
				return
			}
		}
	}
}

// Canonical returns the rotation of `s` with the smallest numeric value,
// along with the distance that produces it from `s`.
//
// All twelve transpositions of a set share one canonical form, so it
// serves as a key: a table indexed by canonical form turns pattern
// matching into a single lookup rather than twelve comparisons.
//
// The returned distance is not the musical root. It is the offset to
// the canonical representative, which is an arbitrary member chosen by
// numeric value. A caller resolving a root composes this offset with
// the root that the matched pattern declares. Getting this backwards
// yields chords that are consistently wrong by a fixed interval, which
// is a symptom worth recognising early.
func (s PitchSet) Canonical() (PitchSet, Semitones) {
	if s.IsEmpty() {
		return s, 0
	}
	best, at := s, Semitones(0)
	for n := Semitones(1); n < 12; n++ {
		if rotated := s.Transpose(n); rotated < best {
			best, at = rotated, n
		}
	}
	return best, at
}

// String returns the twelve bits of `s`, for debugging only. It never
// returns note names.
func (s PitchSet) String() string {
	return fmt.Sprintf("%012b", uint16(s))
}
