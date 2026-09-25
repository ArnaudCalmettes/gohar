package harmony

import (
	"fmt"
	"iter"
	"slices"
)

// A Phrase is a melodic shape: a sequence of offsets from its first
// note.
//
// # Why it is not a Progression
//
// The two look alike and one rule separates them. A progression rebases
// each chord by the nearest motion, folding a fifth up into a fourth
// down, because a listener hears the shortest path between two roots
// and the octave a player voices them in says nothing.
//
// A melody is the opposite. Rising a major seventh and falling a minor
// second land on the same pitch class and are not the same phrase: one
// leaps, the other steps. So the offsets here are signed and never
// folded, and they can exceed an octave.
//
// # What it deliberately does not carry
//
// No rhythm. Duration is a real part of a melody and none of what this
// library does needs it: recognising a shape, comparing two
// performances and telling a player where they diverged all work on the
// sequence of pitches alone. Adding durations would double the
// comparison logic to serve nothing yet.
//
// # Why every input device fits
//
// A phrase is relative, so a keyboard supplies absolute pitches to be
// rebased while a pad or a game controller supplies the offsets
// directly. A device that cannot produce a pitch at all can still
// produce a contour, which [Phrase.Contour] reduces a phrase to.
type Phrase struct {
	offsets []Semitones

	// origin is the pitch the phrase was played from, when it was
	// played rather than declared.
	//
	// Not part of its identity: the same shape from two starting notes
	// compares equal, which is the whole premise. Kept so that
	// [Phrase.Compare] can report how far apart the two performances
	// were, and let the caller decide whether to care.
	origin   Pitch
	anchored bool
}

// A Direction is one step of a contour: up, down, or neither.
type Direction int8

const (
	Down  Direction = -1
	Level Direction = 0
	Up    Direction = 1
)

// NewPhrase builds a phrase from the pitches of one performance.
//
// The first pitch is the origin and sits at offset zero. Every other
// offset is its signed distance from that note, unfolded, so a phrase
// that climbs two octaves has an offset of twenty four.
func NewPhrase(pitches ...Pitch) (Phrase, error) {
	if len(pitches) == 0 {
		return Phrase{}, fmt.Errorf("gohar: a phrase has at least one note")
	}

	offsets := make([]Semitones, len(pitches))
	for i, p := range pitches {
		offsets[i] = p.Sub(pitches[0])
	}
	return Phrase{offsets: offsets, origin: pitches[0], anchored: true}, nil
}

// NewPhraseOffsets declares a shape directly, with no performance
// behind it.
//
// The first offset is the origin, so it is zero by construction and
// anything else is refused.
func NewPhraseOffsets(offsets ...Semitones) (Phrase, error) {
	if len(offsets) == 0 {
		return Phrase{}, fmt.Errorf("gohar: a phrase has at least one note")
	}
	if offsets[0] != 0 {
		return Phrase{}, fmt.Errorf(
			"gohar: the first offset is %d, but it is the origin", offsets[0])
	}
	return Phrase{offsets: slices.Clone(offsets)}, nil
}

// Len returns the number of notes in p.
func (p Phrase) Len() int {
	return len(p.offsets)
}

// Offsets iterates over the notes of p, yielding each index with its
// distance from the first note.
func (p Phrase) Offsets() iter.Seq2[int, Semitones] {
	return func(yield func(int, Semitones) bool) {
		for i, n := range p.offsets {
			if !yield(i, n) {
				return
			}
		}
	}
}

// At voices p from a starting pitch.
//
// Pitches may fall outside the MIDI range for an extreme start; check
// with [Pitch.IsValid] where it matters.
func (p Phrase) At(root Pitch) iter.Seq[Pitch] {
	return func(yield func(Pitch) bool) {
		for _, n := range p.offsets {
			if !yield(root.Transpose(n)) {
				return
			}
		}
	}
}

// Contour reduces p to its directions, one per step, so it has one
// fewer element than the phrase has notes.
//
// This is the shape a device with no pitch can still express, and the
// least a listener can be asked to hear. A phrase and its contour are
// not equivalent: many phrases share a contour, which is the point when
// the exercise is to hear direction before distance.
func (p Phrase) Contour() iter.Seq[Direction] {
	return func(yield func(Direction) bool) {
		for i := 1; i < len(p.offsets); i++ {
			var d Direction
			switch {
			case p.offsets[i] > p.offsets[i-1]:
				d = Up
			case p.offsets[i] < p.offsets[i-1]:
				d = Down
			}
			if !yield(d) {
				return
			}
		}
	}
}

// Compare reads a performance against p.
//
// Same reports whether the shapes match, transposition ignored.
// FirstDivergence locates where they parted, or is -1 when they did
// not, because telling a player which note went wrong is worth more
// than telling them how many did.
//
// Shift is the distance between the two starting notes, in semitones
// and unfolded. A phrase sung an octave higher is transposed by twelve,
// not by nothing: unlike a chord root, a melodic register is audible
// and worth reporting.
func (p Phrase) Compare(played Phrase) Match {
	m := Match{FirstDivergence: -1}

	if p.anchored && played.anchored {
		m.Shift = played.origin.Sub(p.origin)
		m.Transposed = m.Shift != Unison
	}

	shorter := min(len(p.offsets), len(played.offsets))
	for i := range shorter {
		if p.offsets[i] != played.offsets[i] {
			m.FirstDivergence = i
			return m
		}
	}
	if len(p.offsets) != len(played.offsets) {
		m.FirstDivergence = shorter
		return m
	}

	m.Same = true
	return m
}

// CompareContour reads a performance against p by direction alone.
//
// The forgiving comparison, for an exercise that asks a player to hear
// where the line goes before asking by how much. FirstDivergence
// indexes the step, so a value of zero means the very first move went
// the wrong way.
func (p Phrase) CompareContour(played Phrase) Match {
	m := Match{FirstDivergence: -1}

	mine := slices.Collect(p.Contour())
	theirs := slices.Collect(played.Contour())

	shorter := min(len(mine), len(theirs))
	for i := range shorter {
		if mine[i] != theirs[i] {
			m.FirstDivergence = i
			return m
		}
	}
	if len(mine) != len(theirs) {
		m.FirstDivergence = shorter
		return m
	}

	m.Same = true
	return m
}
