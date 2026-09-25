// Package analysis turns what a player holds into what it means.
//
// # Identification, not scoring
//
// An earlier draft of this package scored every pattern against every
// root with tunable weights and returned a ranking. That was the wrong
// shape. Harmony has construction rules, and those rules determine
// most of what a weighting would otherwise guess at: whether a D in a
// C chord is a second or a ninth follows from whether the chord holds
// a third, and no amount of tuning improves on knowing that.
//
// So the pattern is normalised first, by [harmony.ChordPattern.Normalize],
// and then matched against a table of tetrads by equality. A shape
// that matches nothing is not a chord, and saying so is worth more
// than a plausible ranking of wrong answers.
//
// What genuinely remains uncertain is narrow: which held note is the
// root, and what to do when a symmetric chord admits several. Those
// are resolved by ordering rules, stated in [Recognizer.Identify], not
// by numbers.
//
// # Two layers
//
// The recognition layer is pure: a snapshot goes in, readings come
// out, with no state, no clock and no goroutine. The engine above it
// holds the sliding window, the stickiness and the inferred tonality.
// The pure layer is checked by table; the engine is tuned by ear.
package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Snapshot is what the pure layer sees: the pitches sounding at one
// instant.
//
// Pitches rather than classes. The bass decides between readings that
// the classes alone cannot: the same four notes are C major with an
// added sixth over C, and A minor seventh over A.
//
// The slice belongs to the caller. Nothing here retains or mutates it,
// though [Recognizer.Identify] sorts a copy.
type Snapshot struct {
	Pitches []harmony.Pitch
}

// Classes folds the snapshot into the set of classes it sounds.
func (s Snapshot) Classes() harmony.PitchSet {
	var out harmony.PitchSet
	for _, p := range s.Pitches {
		out = out.With(p.Class())
	}
	return out
}

// Bass returns the lowest sounding pitch, and whether anything sounds.
func (s Snapshot) Bass() (harmony.Pitch, bool) {
	if len(s.Pitches) == 0 {
		return 0, false
	}
	bass := s.Pitches[0]
	for _, p := range s.Pitches[1:] {
		if p < bass {
			bass = p
		}
	}
	return bass, true
}

// IsEmpty reports whether nothing sounds.
func (s Snapshot) IsEmpty() bool {
	return len(s.Pitches) == 0
}

// A Reading is one identification of a snapshot.
//
// There is no score. A reading is either what the notes are, under some
// root, or it does not exist. Where several readings exist they are
// genuinely several, and ordering them is a matter of stated rules
// rather than of degree.
type Reading struct {
	// Root is the pitch class the reading is built on, always one of
	// the sounding classes.
	Root harmony.PitchClass

	// Pattern is the normalised pattern, spread over two octaves with
	// every member on the rung harmony assigns it.
	Pattern harmony.ChordPattern

	// Tetrad is the first octave of Pattern, the part that fixes what
	// the chord is. This is what matched the table.
	Tetrad harmony.ChordPattern

	// Extensions is what sits above the tetrad and colours it. Empty
	// for a plain triad or seventh chord.
	Extensions harmony.ChordPattern

	// BassIsRoot reports whether the lowest sounding pitch is Root.
	// The single most useful fact for ordering readings, and about the
	// best one can do with no tonal context at all.
	BassIsRoot bool

	// Degree is the degree of [Context.Tonality] that Root sits on, or
	// zero when there is no tonality or the root falls outside it.
	// Degrees count from one, so zero is unambiguous.
	//
	// # Why this is the field that matters
	//
	// A player without perfect pitch does not hear that they played a
	// G seventh. They hear that they played the dominant. Telling them
	// the absolute root is telling them the one thing they cannot use,
	// and telling them the degree is the whole of the feedback.
	//
	// Which makes the tonality a precondition rather than a comfort:
	// with no tonic there is no degree, and with no degree there is
	// nothing to say to a player who transposed the exercise.
	Degree harmony.Degree
}

// IsZero reports whether r identifies nothing.
func (r Reading) IsZero() bool {
	return r.Pattern == 0
}

// A Context is what the pure layer knows beyond the notes.
//
// It holds no history and no timing, which is what keeps this layer
// testable by table.
type Context struct {
	// Tonality is the key the engine believes it is hearing, or the
	// zero value when it has inferred none.
	//
	// Used to order readings that the bass rule leaves tied, never to
	// admit or reject one. A player who modulates without warning is
	// still identified.
	Tonality harmony.Tonality

	// Current is the reading the caller is already showing, if any.
	//
	// When it is still among the valid readings it is kept, which is
	// what stops a display from flickering between the four equally
	// correct roots of a diminished seventh chord. Stickiness rather than a
	// score margin: with no scores to compare, the rule is simply that
	// a still-valid answer does not get replaced by an equal one.
	Current Reading
}
