package naming

import (
	"fmt"
	"slices"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Namer turns the numbers the core computes into words.
//
// # Why it is a value and not a set of functions
//
// Naming a chord can happen on every frame, so no derivation may take
// place at call time. The catalogue is walked once at construction and
// indexed by canonical form, which turns mode recognition into a single
// map lookup instead of thirty five comparisons across twelve
// rotations.
//
// Holding the locale and the tonal context in the same value is what
// replaces the mutable globals the previous iteration used. A game
// builds one Namer at startup and updates its context when the analysis
// engine revises what key it thinks it hears.
type Namer struct {
	locale   Locale
	tonic    SpelledNote
	tonality harmony.Tonality

	// byPattern maps a scale pattern to the catalogue entry that owns
	// it, built once at construction.
	//
	// A plain map on the pattern rather than an index on canonical
	// forms: the thirty five patterns are distinct, so the pattern
	// alone is a key. Indexing on canonical forms would let a lookup
	// skip the twelve rotations, but there are no rotations to skip
	// here, since a mode is already anchored on its own tonic.
	byPattern map[harmony.ScalePattern]Mode

	// spelled holds the degree spellings of the current tonality,
	// derived once when the context is set.
	spelled []SpelledNote
}

// NewNamer builds a namer for a locale, with no tonal context.
//
// With no context, spelling falls back to the C natural convention. The
// namer does not claim to be in C major, and [Namer.Tonality] returns
// the zero value so that an interface can tell the difference.
func NewNamer(l Locale) (*Namer, error) {
	if l.Letters[0] == "" {
		return nil, fmt.Errorf("naming: locale has no letter names")
	}

	byPattern := make(map[harmony.ScalePattern]Mode, len(catalogue))
	for _, m := range catalogue {
		p := m.Pattern()
		if previous, clash := byPattern[p]; clash {
			return nil, fmt.Errorf(
				"naming: system %d degree %d and system %d degree %d share a pattern",
				previous.System, previous.Degree, m.System, m.Degree)
		}
		byPattern[p] = m
	}

	return &Namer{locale: l, byPattern: byPattern}, nil
}

// WithTonality returns a namer spelling in the given tonal context.
//
// The tonic spelling is chosen by the default convention. Use
// [Namer.WithSpelledTonality] when the caller knows which of two
// enharmonic spellings the tonic should take, which matters as soon as
// a key signature is involved: the same class is F sharp in one key and
// G flat in another, and only the caller knows which.
func (n *Namer) WithTonality(t harmony.Tonality) *Namer {
	return n.WithSpelledTonality(t, defaultSpelling(t.Tonic()))
}

// WithSpelledTonality returns a namer spelling in the given context,
// with the tonic written as given.
//
// Every other spelling in the scale follows from this one. Nothing else
// is chosen.
func (n *Namer) WithSpelledTonality(t harmony.Tonality, tonic SpelledNote) *Namer {
	out := *n
	out.tonality = t
	out.tonic = tonic
	out.spelled = spellDegrees(t, tonic)
	return &out
}

// Tonality returns the context n spells in, or the zero value when it
// has none.
//
// The zero value means no context, never C major. An interface showing
// a key name must check this before displaying anything, or it will
// tell the player they are in a key nobody inferred.
func (n *Namer) Tonality() harmony.Tonality {
	return n.tonality
}

// Note spells a pitch class in the current context.
//
// Constant time: the seven degree spellings are derived when the
// context is set, and a class outside the context falls back to the
// nearest degree carrying an accidental.
func (n *Namer) Note(c harmony.PitchClass) SpelledNote {
	for _, note := range n.spelled {
		if note.Class() == c {
			return note
		}
	}
	if note, ok := nearestSpelling(n.spelled, c); ok {
		return note
	}
	return defaultSpelling(c)
}

// Name spells a pitch class and renders it in the locale.
func (n *Namer) Name(c harmony.PitchClass) string {
	return n.locale.Name(n.Note(c))
}

// Scale spells every degree of the current context, in degree order.
//
// Returns nil when n has no context: there is no scale to spell, and an
// empty result says so more honestly than seven notes of C major would.
func (n *Namer) Scale() []SpelledNote {
	if n.tonality.IsZero() {
		return nil
	}
	return slices.Clone(n.spelled)
}

// Mode identifies the mode a pattern is, anchored on a tonic.
//
// One map lookup on the canonical form, then a comparison against the
// entries reaching it. The thirty five patterns are distinct, so at
// most one matches and the result needs no tonal context to
// disambiguate.
func (n *Namer) Mode(t harmony.Tonality) (Mode, bool) {
	if t.IsZero() {
		return Mode{}, false
	}
	m, ok := n.byPattern[t.Pattern()]
	return m, ok
}

// ModeName identifies a mode and renders its name in the locale.
func (n *Namer) ModeName(t harmony.Tonality) (string, bool) {
	m, ok := n.Mode(t)
	if !ok {
		return "", false
	}
	return n.locale.ModeName(m), true
}
