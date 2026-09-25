package harmony

import (
	"fmt"
	"slices"
)

// A Target is what a slot asks the player to reach.
//
// # Why a target and not a chord
//
// Asking for a chord and asking for a dominant are not the same
// question, and the difference is not a matter of strictness. A
// tritone substitution is a dominant of C only in relation to C: the
// same D flat seventh is a dominant of G flat in another context. So
// the role cannot be read off the chord, and a table from tetrad to
// function is necessary but never sufficient.
//
// Once a target exists, a single chord is the one step case of a move
// toward it, which is why there is no separate type for it.
type Target struct {
	// Root is the pitch class the approach is expected to reach.
	Root PitchClass

	// Tetrad is the chord expected on arrival, when the slot names
	// one. Zero when the slot only asks for a destination and leaves
	// the colour to the player.
	Tetrad ChordPattern

	// Key is the tonal context the approach is judged in.
	//
	// It decides which extensions belong: a minor ninth over a
	// dominant belongs to a mode of the harmonic minor, while the same
	// note over a major tonic belongs to nothing. Zero means no
	// context, and then extensions are recorded without being judged.
	Key Tonality
}

// An Approach is what a player produced inside one slot: the chords
// played, in order, on their way to a target.
//
// Usually one chord. Sometimes more: a suspension resolving onto its own
// third, or a secondary dominant inserted before the arrival. The type
// does not care how many: a move of one step is the ordinary case, not
// a special one.
type Approach struct {
	Chords []Chord
}

// A Resolution is what an approach was found to do.
//
// # Facts, not a verdict
//
// Every field below records something observed. None of them scores the
// approach, and there is deliberately no total.
//
// This is the same choice made for chord identification and for the
// same reason: a score always produces a best answer, including for a
// path that makes no sense. What is not recognised should be said,
// not rated. A game builds its damage out of these fields, and a
// teaching interface can name exactly what it saw.
//
// What is missing from this type is any judgement of how far a path may
// wander. Whether a two five may stand in for a five, or how long a
// chain of secondary dominants may run before it stops being an
// approach, is not derivable from anything in this package. Those rules
// belong to a musician; when they are written down they will sit above
// this type, not inside it.
type Resolution struct {
	// Arrived reports whether the last chord of the approach is built
	// on the target's root, and matches its tetrad when the target
	// names one.
	Arrived bool

	// Dominant reports whether some chord of the approach carried the
	// tritone that resolves onto the target.
	//
	// Derived rather than tabulated. What makes a dominant pull is the
	// tritone between its third and its seventh, and that same tritone
	// belongs to two chords a tritone apart. A G seventh and a D flat
	// seventh share it, so the substitution is not a special case to
	// be listed: it is the same object seen from another root.
	Dominant bool

	// DominantRoot is the root of the chord that carried it, and
	// Substitute reports whether that root was the tritone
	// substitution rather than the expected fifth above the target.
	DominantRoot PitchClass
	Substitute   bool

	// SuspensionAt is the index of a suspension that a later chord on
	// the same root resolved onto its third, or minus one when there
	// was none.
	//
	// Two chords, one degree. An evaluation that counted chords would
	// see a spurious extra one here, which is why an approach is a
	// sequence and not a chord.
	//
	// A suspension is a subdominant, never a dominant. The seventh sus
	// four has no third, so it has no tritone: Bill Evans made it a
	// chord in its own right by removing exactly the interval that
	// makes a dominant pull, which turns the cadence plagal. That is
	// why [Target.Resolve] does not report it under Dominant, and why a
	// game should not reward it as one. Where it sits is the caller's
	// to read: on the target root it is the plagal colour that makes
	// the ionian audible, its fourth being that mode's characteristic
	// degree; on another root it is an approach in its own right.
	SuspensionAt int

	// Extensions holds the classes sounded beyond the tetrads of the
	// approach, and InKey those of them that belong to the mode of
	// their degree in the target's key.
	//
	// The two are separate on purpose. A colour note that belongs is
	// not the same event as one that does not, and a game that
	// rewarded the first while ignoring the difference would teach
	// players to pile extensions on rather than to choose them.
	Extensions PitchSet
	InKey      PitchSet

	// Motion is the total voice movement across the approach, in
	// semitones, counting each voice's shortest path.
	//
	// Low means the chords were connected. This is what tells a
	// well chosen inversion from a mechanical root position, and it
	// is the one thing here that has nothing to do with which chords
	// were played and everything to do with how.
	Motion Semitones

	// Steps is the number of chords the approach took.
	Steps int
}

// NewApproach collects the chords of one slot.
func NewApproach(chords ...Chord) (Approach, error) {
	if len(chords) == 0 {
		return Approach{}, fmt.Errorf("gohar: an approach has at least one chord")
	}
	return Approach{Chords: slices.Clone(chords)}, nil
}


// Resolve reads an approach against a target and reports what it did.
//
// Nothing here refuses. An approach that reaches nothing comes back
// with Arrived false and whatever else was observed, because telling a
// player that they never landed is more useful than telling them their
// attempt scored badly.
func (t Target) Resolve(a Approach) Resolution {
	r := Resolution{Steps: len(a.Chords), SuspensionAt: -1}
	if len(a.Chords) == 0 {
		return r
	}

	last := a.Chords[len(a.Chords)-1]
	r.Arrived = last.Root == t.Root &&
		(t.Tetrad == 0 || last.Pattern.Tetrad() == t.Tetrad)

	for i, c := range a.Chords {
		if root, ok := t.dominantOf(c); ok && !r.Dominant {
			r.Dominant = true
			r.DominantRoot = root
			r.Substitute = root != t.Root.Transpose(7)
		}
		if r.SuspensionAt < 0 && suspensionResolves(a.Chords, i) {
			r.SuspensionAt = i
		}
		r.Extensions = r.Extensions.Union(extensionClasses(c))
	}

	r.InKey = t.extensionsInKey(r.Extensions)
	r.Motion = voiceMotion(a.Chords)
	return r
}

// dominantOf reports whether a chord carries the tritone that resolves
// onto the target, and on which root it sits.
//
// The test is on the tritone rather than on the chord's name, which is
// what makes the substitution fall out instead of needing a table. The
// tritone of a dominant on the fifth above the target and the tritone
// of the one a tritone away are the same two classes.
func (t Target) dominantOf(c Chord) (PitchClass, bool) {
	if !c.Pattern.HasAll(IntMajorThird, IntMinorSeventh) {
		return 0, false
	}

	third := c.Root.Transpose(IntMajorThird.Semitones)
	seventh := c.Root.Transpose(IntMinorSeventh.Semitones)

	// The leading tone rises a semitone to the target and the seventh
	// falls a semitone to its third. Either orientation of the pair is
	// the same tritone, which is exactly why two roots share it.
	if third.Transpose(1) == t.Root || seventh.Transpose(-1) == t.Root.Transpose(4) {
		return c.Root, true
	}
	if seventh.Transpose(1) == t.Root || third.Transpose(-1) == t.Root.Transpose(4) {
		return c.Root, true
	}
	return 0, false
}

// suspensionResolves reports whether the chord at index i is a
// suspension that a later chord on the same root resolves.
//
// A suspension has no third and holds a fourth or a second in its
// place; it resolves when a following chord on that same root has the
// third the suspension was standing in for.
func suspensionResolves(chords []Chord, i int) bool {
	c := chords[i]
	if c.Pattern.HasAny(IntMajorThird, IntMinorThird) {
		return false
	}
	if !c.Pattern.HasAny(IntPerfectFourth, IntMajorSecond) {
		return false
	}

	for _, next := range chords[i+1:] {
		if next.Root != c.Root {
			continue
		}
		if next.Pattern.HasAny(IntMajorThird, IntMinorThird) {
			return true
		}
	}
	return false
}

// extensionClasses returns the classes a chord sounds beyond its
// tetrad.
func extensionClasses(c Chord) PitchSet {
	upper := c.Pattern &^ c.Pattern.Tetrad()
	return upper.Fold().Transpose(Semitones(c.Root))
}

// extensionsInKey keeps the extensions that belong to the mode of their
// own degree in the target's key.
//
// With no key there is nothing to judge against, and the result is
// empty rather than everything: an unjudged colour is not a validated
// one.
func (t Target) extensionsInKey(extensions PitchSet) PitchSet {
	if t.Key.IsZero() {
		return EmptyPitchSet
	}

	out := EmptyPitchSet
	for c := range extensions.Classes() {
		if t.Key.Contains(c) {
			out = out.With(c)
		}
	}
	return out
}

// voiceMotion sums how far the chords moved from one to the next,
// taking each class by its shortest path.
//
// A rough measure: it works on classes, so it cannot tell a voice that
// leapt an octave from one that stayed put. Refining it needs the
// pitches actually held, which live in the analysis engine rather than
// here.
func voiceMotion(chords []Chord) Semitones {
	var total Semitones
	for i := 1; i < len(chords); i++ {
		total += setDistance(chords[i-1].Set(), chords[i].Set())
	}
	return total
}

// setDistance measures how far one set of classes sits from another, by
// counting the classes that changed.
func setDistance(from, to PitchSet) Semitones {
	return Semitones(from.Difference(to).Len() + to.Difference(from).Len())
}
