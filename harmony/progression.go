package harmony

import (
	"fmt"
	"iter"
	"slices"
)

// A Step is one chord of a [Progression], placed by its distance from
// the progression's origin.
//
// A distance and not a root. A progression has no key of its own, and
// giving one to its first chord would make every transposition of the
// same progression a different value.
type Step struct {
	// Offset is the distance from the origin to this chord's root.
	// The first step of a progression is at zero by construction.
	Offset Semitones

	// Tetrad is what the chord is. Extensions colour a chord without
	// changing which chord it is, so they are not part of the shape a
	// progression describes.
	Tetrad ChordPattern
}

// A Progression is a sequence of chords held relative to its own
// origin.
//
// # Why relative
//
// A player without perfect pitch hears intervals, not keys. Asked for
// a two five one, they play one, and whether they started on D or on
// F sharp is not what was being asked. So the progression is the
// shape, and the key is a fact about one performance of it.
//
// This is the same reasoning that keeps [ScalePattern] free of a
// tonic, and it has the same consequence: a progression has no
// Transpose method, because moving a relative object by a distance
// means nothing. Anchoring one on a root is [Progression.At].
//
// # Why the whole sequence, not each chord
//
// A progression played a tone higher has every chord moved by the same
// tone. Normalising each chord on its own would lose that: it would
// accept a sequence whose chords were individually right and whose
// relations were wrong, which is precisely the mistake a learner makes
// and the one worth catching.
type Progression struct {
	steps []Step

	// origin is the root the progression was built from, when it was
	// built from absolute roots at all.
	//
	// A progression is a shape and carries no key, so this is not part
	// of its identity: two performances of the same shape in different
	// keys compare equal. It is kept only so that [Progression.Compare]
	// can report how far apart they were, which the caller needs in
	// order to decide whether to care.
	origin PitchClass

	// anchored records whether origin means anything. A shape declared
	// from steps has no origin, and a shift against it is zero rather
	// than a number computed from a field nobody set.
	anchored bool
}

// NewProgression builds a progression from chord roots given as
// absolute classes, rebasing them on the first.
//
// The roots are read as an ascending line where each chord is taken to
// be within a tritone of the last, which is how a listener hears a
// progression: the nearest motion, not the largest. A caller wanting
// wider leaps builds the steps directly with [NewProgressionSteps].
func NewProgression(tetrads []ChordPattern, roots []PitchClass) (Progression, error) {
	if len(tetrads) != len(roots) {
		return Progression{}, fmt.Errorf(
			"gohar: %d tetrads for %d roots", len(tetrads), len(roots))
	}
	if len(roots) == 0 {
		return Progression{}, fmt.Errorf("gohar: a progression has at least one chord")
	}

	steps := make([]Step, len(roots))
	offset := Semitones(0)
	for i := range roots {
		if !roots[i].IsValid() {
			return Progression{}, fmt.Errorf("gohar: root %d is out of range", roots[i])
		}
		if !tetrads[i].HasOffset(0) {
			return Progression{}, fmt.Errorf(
				"gohar: chord %d does not contain its own root", i)
		}
		if i > 0 {
			offset += nearestMotion(roots[i-1], roots[i])
		}
		steps[i] = Step{Offset: offset, Tetrad: tetrads[i]}
	}
	return Progression{steps: steps, origin: roots[0], anchored: true}, nil
}

// nearestMotion returns the shortest distance from one root to the
// next, preferring the ascending one when they tie.
//
// A listener hears the nearest motion. Told that a progression goes
// from C to G, nobody hears a descent of a fifth followed by nothing;
// they hear the fourth up. Taking the raw difference between classes
// would make the same progression a different shape depending on which
// octave the player happened to voice it in.
func nearestMotion(from, to PitchClass) Semitones {
	d := from.Up(to)
	if d > 6 {
		d -= 12
	}
	return d
}

// NewProgressionSteps builds a progression from explicit steps.
//
// Fails when the first step is not at offset zero, which would mean
// the progression carries an origin it does not have.
func NewProgressionSteps(steps ...Step) (Progression, error) {
	if len(steps) == 0 {
		return Progression{}, fmt.Errorf("gohar: a progression has at least one chord")
	}
	if steps[0].Offset != 0 {
		return Progression{}, fmt.Errorf(
			"gohar: the first step sits at offset %d, but it is the origin",
			steps[0].Offset)
	}
	return Progression{steps: slices.Clone(steps)}, nil
}

// Len returns the number of chords in p.
func (p Progression) Len() int {
	return len(p.steps)
}

// Steps iterates over the chords of p in order.
func (p Progression) Steps() iter.Seq2[int, Step] {
	return func(yield func(int, Step) bool) {
		for i, s := range p.steps {
			if !yield(i, s) {
				return
			}
		}
	}
}

// At anchors p on a root and yields the chords it becomes.
//
// The conversion out of relative space, as [ScalePattern.At] is for
// scales.
func (p Progression) At(root PitchClass) iter.Seq[Chord] {
	return func(yield func(Chord) bool) {
		for _, s := range p.steps {
			chord := Chord{Root: root.Transpose(s.Offset), Pattern: s.Tetrad}
			if !yield(chord) {
				return
			}
		}
	}
}

// Degrees returns the degree each chord of p is built on, once p is
// anchored on the tonic of t.
//
// This is what a player is told. Naming the chords of an exercise by
// their absolute roots is useless to someone who transposed it; naming
// them by degree is the whole point, and it is why a tonality is a
// precondition of teaching rather than a convenience.
//
// A chord whose root falls outside t has no degree, and its entry is
// zero. Degrees are counted from one, so zero is unambiguous.
func (p Progression) Degrees(root PitchClass, t Tonality) []Degree {
	out := make([]Degree, len(p.steps))
	for i, s := range p.steps {
		if d, ok := t.DegreeOf(root.Transpose(s.Offset)); ok {
			out[i] = d
		}
	}
	return out
}

// A Match is the result of comparing two progressions.
//
// # Why the shift is always returned
//
// Because ignoring it has to be a decision, not a side effect. A game
// that accepts any key still knows which one was played, and can say
// so; a game that wants the requested key has the number it needs to
// say how far off. Collapsing the comparison to a boolean would take
// that choice away from the caller, who is the only one who knows what
// the exercise was asking.
type Match struct {
	// Shift is the distance from the expected progression to the one
	// played. Zero when they start on the same root.
	Shift Semitones

	// Transposed reports whether Shift is anything but zero. A
	// convenience, so that a caller reads the intent rather than
	// comparing to zero.
	Transposed bool

	// Same reports whether the two progressions have the same shape,
	// which is the question transposition does not affect.
	Same bool

	// FirstDivergence is the index of the earliest step that differs,
	// or minus one when none does.
	//
	// The index matters more than the count. A learner who fumbles the
	// third chord of a two five one did the first two right, and being
	// told where it went wrong is worth more than being told how many
	// were wrong.
	FirstDivergence int
}

// Compare matches p against played, ignoring transposition but
// reporting it.
//
// # How the shift is found
//
// One shift for the whole sequence, taken from the first step, then
// checked against the rest. Resolving a shift per chord would accept a
// sequence whose chords were each some transposition of the expected
// one while the relations between them were wrong.
//
// # What counts as the same shape
//
// The tetrads, and the offsets between roots. Extensions do not: a
// player who voices a ninth on the dominant played the progression
// that was asked for.
//
// Progressions of different lengths never match. FirstDivergence is
// then the first index past the shorter one.
func (p Progression) Compare(played Progression) Match {
	m := Match{FirstDivergence: -1}

	if p.anchored && played.anchored {
		m.Shift = nearestMotion(p.origin, played.origin)
		m.Transposed = m.Shift != Unison
	}

	shorter := min(len(p.steps), len(played.steps))
	for i := range shorter {
		if p.steps[i] != played.steps[i] {
			m.FirstDivergence = i
			return m
		}
	}
	if len(p.steps) != len(played.steps) {
		m.FirstDivergence = shorter
		return m
	}

	m.Same = true
	return m
}

// Common progressions, given as the shapes a learner practises rather
// than as chords in any key.
var (
	// TwoFiveOneMajor is the two five one in major: minor seventh,
	// dominant seventh, major seventh, each a fourth above the last.
	TwoFiveOneMajor = Progression{steps: []Step{
		{0, ChordMinorSeventh},
		{5, ChordDominantSeventh},
		{10, ChordMajorSeventh},
	}}

	// TwoFiveOneMinor is its minor counterpart, built on the half
	// diminished second degree.
	//
	// The resolution is a minor major seventh and not a minor seventh.
	// A minor seventh sounds subdominant, which is what it is two bars
	// earlier; landing on one would leave the progression unresolved.
	// See [TwoFiveOneMinorSixth] for the other tonic minor.
	TwoFiveOneMinor = Progression{steps: []Step{
		{0, ChordHalfDiminished},
		{5, ChordDominantSeventh},
		{10, ChordMinorMajorSeventh},
	}}

	// TwoFiveOneMinorSixth resolves onto a minor sixth instead.
	//
	// The two tonic minors, not two spellings of one. Both are what a
	// minor two five one lands on, and which one a player practises is
	// a choice about sound rather than a detail of notation.
	TwoFiveOneMinorSixth = Progression{steps: []Step{
		{0, ChordHalfDiminished},
		{5, ChordDominantSeventh},
		{10, ChordMinorSixth},
	}}
)
