package harmony

import (
	"fmt"
	"iter"
	"math/bits"
)

// A ChordPattern is the structure of a chord, independent of any root.
// Bit n is set when the chord contains a note n semitones above its
// root.
//
// # Why twenty four bits
//
// A chord stacks past the octave, and the position in the stack carries
// meaning that a fold into twelve classes destroys. A ninth sits at 14
// and a second at 2; a minor ninth at 13 and a minor second at 1; a
// raised eleventh at 18 and a tritone at 6. Those pairs sound different
// and function differently, so they occupy different bits.
//
// This is not spelling. Spelling is choosing a name for a fixed
// distance; this is a distance that genuinely differs.
//
// # Invariants
//
// Bits 24 through 31 are always zero, and bit 0 is always set: a chord
// contains its own root.
//
// # Why there is no Transpose
//
// Like [ScalePattern], a chord pattern is relative, so moving it by a
// distance means nothing. Anchoring it on a root is [ChordPattern.From]
// or [Chord].
type ChordPattern uint32

// The tetrads, each given with the offsets it encodes so that the
// hexadecimal stays checkable by eye.
//
// # Why the omit-five variants are here
//
// The fifth of a tetrad is not always played, and in jazz a plain
// perfect fifth on a dominant chord is avoided outright, since it
// settles a chord whose job is to be unstable. So a rootless or
// fifthless voicing is not a chord with a note missing, it is a chord
// as commonly played, and it deserves its own pattern rather than being
// read as an incomplete one.
//
// This is the difference between identifying and guessing. A
// normalised pattern either equals one of these or is not a chord.
//
// Chord qualities are named here for the same reason scale patterns
// are: one arrangement of distances, one name, no second reading to
// choose between.
const (
	ChordMajorTriad            ChordPattern = 0x000091 // 0, 4, 7      triade majeure
	ChordMinorTriad            ChordPattern = 0x000089 // 0, 3, 7      triade mineure
	ChordDiminishedTriad       ChordPattern = 0x000049 // 0, 3, 6      triade diminuée
	ChordAugmentedTriad        ChordPattern = 0x000111 // 0, 4, 8      triade augmentée
	ChordSus2                  ChordPattern = 0x000085 // 0, 2, 7      sus2
	ChordSus4                  ChordPattern = 0x0000a1 // 0, 5, 7      sus4
	ChordMajorSixth            ChordPattern = 0x000291 // 0, 4, 7, 9   6
	ChordMinorSixth            ChordPattern = 0x000289 // 0, 3, 7, 9   mineur 6
	ChordMajorSeventh          ChordPattern = 0x000891 // 0, 4, 7, 11  majeur 7
	ChordMajorSeventhNo5       ChordPattern = 0x000811 // 0, 4, 11     majeur 7 sans quinte
	ChordMajorSeventhSharp5    ChordPattern = 0x000911 // 0, 4, 8, 11  majeur 7 dièse 5
	ChordDominantSeventh       ChordPattern = 0x000491 // 0, 4, 7, 10  7
	ChordDominantSeventhNo5    ChordPattern = 0x000411 // 0, 4, 10     7 sans quinte
	ChordDominantSeventhFlat5  ChordPattern = 0x000451 // 0, 4, 6, 10  7 bémol 5
	ChordDominantSeventhSharp5 ChordPattern = 0x000511 // 0, 4, 8, 10  7 dièse 5
	ChordMinorSeventh          ChordPattern = 0x000489 // 0, 3, 7, 10  mineur 7
	ChordMinorSeventhNo5       ChordPattern = 0x000409 // 0, 3, 10     mineur 7 sans quinte
	ChordHalfDiminished        ChordPattern = 0x000449 // 0, 3, 6, 10  mineur 7 bémol 5
	ChordMinorMajorSeventh     ChordPattern = 0x000889 // 0, 3, 7, 11  mineur majeur 7
	ChordMinorMajorSeventhNo5  ChordPattern = 0x000809 // 0, 3, 11     mineur majeur 7 sans quinte
	ChordDiminishedSeventh     ChordPattern = 0x000249 // 0, 3, 6, 9   diminué 7
	ChordDominantSeventhSus2   ChordPattern = 0x000485 // 0, 2, 7, 10  7sus2
	ChordDominantSeventhSus4   ChordPattern = 0x0004a1 // 0, 5, 7, 10  7sus4
)

// NewChordPattern builds a pattern from offsets above the root.
//
// Offset 0 is implied and may be omitted. Offsets outside [0, 24) are
// rejected. Duplicates are harmless.
func NewChordPattern(offsets ...Semitones) (ChordPattern, error) {
	p := ChordPattern(1)
	for _, n := range offsets {
		if n < 0 || n >= 24 {
			return 0, fmt.Errorf("gohar: chord offset %d is out of range", n)
		}
		p |= 1 << n
	}
	return p, nil
}

// HasOffset reports whether `p` includes a note `n` semitones above its
// root.
//
// Position sensitive: a normalised pattern holding a ninth does not
// have offset 2. Callers asking the pitch class question want
// [ChordPattern.Fold] first, and callers asking the harmonic question
// want [ChordPattern.Has] with an interval.
func (p ChordPattern) HasOffset(n Semitones) bool {
	if n < 0 || n >= 24 {
		return false
	}
	return p&(1<<n) != 0
}

// With returns `p` with an offset added.
func (p ChordPattern) With(n Semitones) ChordPattern {
	if n < 0 || n >= 24 {
		return p
	}
	return p | 1<<n
}

// Without returns `p` with an offset removed. Removing offset 0 is
// refused, since it would break the invariant.
func (p ChordPattern) Without(n Semitones) ChordPattern {
	if n <= 0 || n >= 24 {
		return p
	}
	return p &^ (1 << n)
}

// Len returns the number of notes in `p`.
func (p ChordPattern) Len() int {
	return bits.OnesCount32(uint32(p))
}

// Offsets iterates over the members of `p` in ascending order, starting
// at 0.
//
// Derived from the bits alone. No member is inserted, moved or patched
// afterwards to accommodate an irregular chord: an earlier version of
// this library rewrote a fixed index of its output to handle diminished
// sevenths and raised ninths, which held only as long as no pattern
// gained a member before that index. Irregular cases are decided by
// which bits are set, never by correcting a result already built.
func (p ChordPattern) Offsets() iter.Seq[Semitones] {
	return func(yield func(Semitones) bool) {
		for n := Semitones(0); n < 24; n++ {
			if p.HasOffset(n) && !yield(n) {
				return
			}
		}
	}
}

// Has reports whether `p` contains the given interval.
//
// Enharmonic by necessity: the test is on the semitone count, since a
// bit mask holds no degree. A pattern with a raised ninth answers true
// to a minor tenth, and that is not a defect. One hears the raised
// ninth in an altered chord, and calls it a tenth only once the chord
// has been identified, so the two have to match the same bit until the
// identification happens.
func (p ChordPattern) Has(i Interval) bool {
	return p.HasOffset(i.Semitones)
}

// HasAny and HasAll report whether `p` contains any, or all, of the
// given intervals.
func (p ChordPattern) HasAny(intervals ...Interval) bool {
	for _, i := range intervals {
		if p.Has(i) {
			return true
		}
	}
	return false
}

// HasAll reports whether `p` contains every one of the given intervals.
func (p ChordPattern) HasAll(intervals ...Interval) bool {
	if len(intervals) == 0 {
		return false
	}
	for _, i := range intervals {
		if !p.Has(i) {
			return false
		}
	}
	return true
}

// Contains reports whether every member of `other` belongs to `p`.
func (p ChordPattern) Contains(other ChordPattern) bool {
	return p&other == other
}

// Collapse folds `p` into its first octave, keeping it a pattern.
//
// The counterpart of [ChordPattern.Fold], which leaves pattern space
// for a [PitchSet]. Collapse stays here, because it is the first step
// of [ChordPattern.Normalize] and normalising needs a pattern to work
// on.
func (p ChordPattern) Collapse() ChordPattern {
	return (p | p>>12) & 0x000fff
}

// Normalize spreads `p` over two octaves, putting every member on the
// rung harmony assigns it.
//
// # What this replaces
//
// Whether a D in a C chord is a second or a ninth is not a matter of
// taste, and not something to weigh against alternatives: it is
// determined by whether the chord holds a third. Normalize is where
// that determination happens, once, before anything is compared. A
// pattern that has been normalised can be matched against a table of
// tetrads by equality, and a shape that matches nothing is not a chord.
// There is no closest match to fall back on.
//
// # The rules it encodes
//
// These are constraints on how chords are built, not heuristics:
//
//   - No chord holds a minor second. It is always a minor ninth.
//   - A chord holding a third holds neither second nor fourth. They
//     become a ninth and an eleventh.
//   - A chord holding both thirds has the minor one as a raised ninth,
//     which is the altered chord: its diminished fourth is heard as
//     the major third of the tetrad, pushing the minor third up.
//   - A chord holding a fourth and no third has its second as a ninth.
//   - A chord holding a perfect fifth has any diminished fifth as a
//     raised eleventh.
//   - A chord holding a diminished triad has any major seventh as a
//     major fourteenth.
//   - A chord holding a seventh holds no sixth. The sixth becomes a
//     thirteenth, minor or major.
//
// Applied in that order, since the later tests read a pattern the
// earlier ones have already moved.
//
// # How a member moves
//
// A member is raised by swapping its bit with the bit an octave above,
// and only when exactly one of the two is set: see moveUp.
func (p ChordPattern) Normalize() ChordPattern {
	p = p.Collapse()
	p = p.moveUp(IntMinorSecond)

	if p.HasAny(IntMajorThird, IntMinorThird) {
		p = p.moveUp(IntMajorSecond)
		p = p.moveUp(IntPerfectFourth)
		if p.HasAll(IntMajorThird, IntMinorThird) {
			p = p.moveUp(IntMinorThird)
		}
	}
	if p.Has(IntPerfectFourth) {
		p = p.moveUp(IntMajorSecond)
	}
	if p.Has(IntPerfectFifth) {
		p = p.moveUp(IntDiminishedFifth)
	}
	if p.Contains(ChordDiminishedTriad) {
		p = p.moveUp(IntMajorSeventh)
	}
	if p.Has(IntMinorSeventh) {
		if p.hasFifth() {
			p = p.moveUp(IntAugmentedFifth)
		}
		p = p.moveUp(IntMajorSixth)
	}
	if p.Has(IntMajorSeventh) {
		p = p.moveUp(IntDiminishedFifth)
		p = p.moveUp(IntMajorSixth)
		if p.hasFifth() {
			p = p.moveUp(IntMinorSixth)
		}
	}
	return p
}

// Intervals iterates over the members of `p` as intervals, degrees
// included.
//
// Meaningful on a normalised pattern only. On a raw one the degree
// assignment is the very question Normalize answers, and this would be
// guessing at it.
func (p ChordPattern) Intervals() iter.Seq[Interval] {
	return func(yield func(Interval) bool) {
		diminished := p.Contains(ChordDiminishedTriad)
		for n := range p.Offsets() {
			i, ok := intervalAtOffset(n, diminished)
			if !ok {
				continue
			}
			if !yield(i) {
				return
			}
		}
	}
}

// Tetrad returns the first octave of `p`, which is what determines the
// chord's nature.
//
// Harmony holds that the first four sounds fix what a chord is, and
// everything above colours it. So identification compares tetrads, and
// the upper structure is read afterwards as extensions.
func (p ChordPattern) Tetrad() ChordPattern {
	return p & 0x000fff
}

// Fold collapses `p` into the set of pitch classes it sounds.
//
// Lossy and one way. A ninth and a second fold to the same class, so
// the distinction that justified twenty four bits does not survive.
// There is no unfold.
//
// This is nonetheless the bridge to recognition. What a keyboard
// produces is a [PitchSet], so matching a played set against a pattern
// happens on folded forms. The ambiguity Fold introduces is lifted on
// the way back, not here: recognition rebuilds a pattern from the held
// notes and [ChordPattern.Normalize] decides whether a D is a second or
// a ninth.
func (p ChordPattern) Fold() PitchSet {
	return PitchSet(p.Collapse())
}

// IsSubsetOf reports whether every member of `p` belongs to `other`,
// compared at full twenty four bit resolution.
func (p ChordPattern) IsSubsetOf(other ChordPattern) bool {
	return p&other == p
}

// From iterates over the members of p, starting at root and moving up.
//
// On a [Pitch] the iteration crosses octaves and preserves the stack:
// a ninth lands an octave and a tone above the root, where it belongs.
//
// On a [PitchClass] the octave is lost but the order is not. A ninth
// still comes last, folded onto the class of a second, so the sequence
// wraps rather than climbs. That is the same behaviour as
// [ScalePattern.From], where a scale read from its own tonic wraps
// past the octave, and it is not what [ChordPattern.Fold] gives: a set
// has no order at all.
//
// Meaningful on a normalised pattern, for the same reason as
// [ChordPattern.Intervals].
func (p ChordPattern) From[T Transposable[T]](root T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := range p.Offsets() {
			if !yield(root.Transpose(n)) {
				return
			}
		}
	}
}

// String returns the twenty four bits of `p`, for debugging only.
func (p ChordPattern) String() string {
	return fmt.Sprintf("%024b", uint32(p))
}

// A Chord is a pattern anchored on a root.
//
// The zero value is not usable: its pattern is empty, which violates
// the bit 0 invariant.
type Chord struct {
	Root    PitchClass
	Pattern ChordPattern
}

// NewChord pairs a root with a pattern.
func NewChord(root PitchClass, p ChordPattern) (Chord, error) {
	if !root.IsValid() {
		return Chord{}, fmt.Errorf("gohar: root %d is out of range", root)
	}
	if !p.HasOffset(0) {
		return Chord{}, fmt.Errorf("gohar: chord pattern does not contain its own root")
	}
	return Chord{Root: root, Pattern: p}, nil
}

// Set returns the classes `c` sounds, folded.
func (c Chord) Set() PitchSet {
	return c.Pattern.Fold().Transpose(Semitones(c.Root))
}

// Pitches voices `c` once, `from` the lowest root at or above `from`, and
// stops at `to`.
//
// One stack, not every chord tone in the range. This differs from
// [Scale.Pitches] on purpose, and the difference is what the two are
// for: a scale in a register is every one of its notes there, while a
// chord in a register is one voicing of it. Yielding the upper notes
// of the octave below as well would hand back something no player is
// holding.
//
// The bounds are therefore a register rather than a filter: `from` picks
// which voicing, `to` cuts it short if the register is too narrow `to`
// hold the whole chord.
func (c Chord) Pitches(from, to Pitch) iter.Seq[Pitch] {
	return func(yield func(Pitch) bool) {
		for octave := int8(-1); octave < 10; octave++ {
			root := c.Root.Pitch(octave)
			if root < from {
				continue
			}
			for n := range c.Pattern.Offsets() {
				p := root.Transpose(n)
				if p > to || !yield(p) {
					return
				}
			}
			return
		}
	}
}

// moveUp raises a member by an octave, and only when exactly one of the
// two positions is occupied.
//
// The check guards the swap, not the music. With both bits set, the XOR
// would clear both and the class would vanish from the pattern. A valid
// pattern never holds a class twice, so this only happens to one forged
// by hand, which is left as it is rather than silently damaged.
func (p ChordPattern) moveUp(i Interval) ChordPattern {
	if !p.Has(i) {
		return p
	}
	mask := ChordPattern(1)<<i.Semitones | ChordPattern(1)<<(i.Semitones+12)
	if bits.OnesCount32(uint32(p&mask)) != 1 {
		return p
	}
	return p ^ mask
}

// hasFifth reports whether the chord already holds a fifth other than
// an augmented one.
//
// Offset eight is both an augmented fifth and a minor sixth, and no bit
// mask distinguishes them. The chord does: a note at eight can only be
// a sixth if a fifth is already sounding. Otherwise it is the fifth,
// and pushing it up to a lowered thirteenth would leave the chord
// without one.
//
// This diverges from an earlier published version of this algorithm,
// which raised offset eight unconditionally under a seventh. That made
// the major seventh sharp five collapse onto the fifthless major
// seventh, so a lydian sharp five chord could never be identified as
// itself despite the nomenclature giving it a sharp five.
func (p ChordPattern) hasFifth() bool {
	return p.Has(IntPerfectFifth) || p.Has(IntDiminishedFifth)
}

// intervalAtOffset reads a normalised offset as the interval harmony
// assigns it.
//
// Offset nine is the one place the reading needs context: it is a
// diminished seventh in a diminished chord and a major sixth
// everywhere else, and both spellings are the same bit.
func intervalAtOffset(n Semitones, diminished bool) (Interval, bool) {
	switch n {
	case 0:
		return IntUnison, true
	case 2:
		return IntMajorSecond, true
	case 3:
		return IntMinorThird, true
	case 4:
		return IntMajorThird, true
	case 5:
		return IntPerfectFourth, true
	case 6:
		return IntDiminishedFifth, true
	case 7:
		return IntPerfectFifth, true
	case 8:
		return IntAugmentedFifth, true
	case 9:
		if diminished {
			return IntDiminishedSeventh, true
		}
		return IntMajorSixth, true
	case 10:
		return IntMinorSeventh, true
	case 11:
		return IntMajorSeventh, true
	case 13:
		return IntMinorNinth, true
	case 14:
		return IntMajorNinth, true
	case 15:
		return IntAugmentedNinth, true
	case 17:
		return IntPerfectEleventh, true
	case 18:
		return IntAugmentedEleventh, true
	case 20:
		return IntMinorThirteenth, true
	case 21:
		return IntMajorThirteenth, true
	case 23:
		return IntMajorFourteenth, true
	}
	return Interval{}, false
}
