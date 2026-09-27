package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// Ticks measure time in a sequence of changes: TicksPerBeat to a beat.
type Ticks int

// TicksPerBeat leaves room below the beat. Chords read from a chart
// start on a beat, but chords a player sounds do not.
const TicksPerBeat Ticks = 2520

// A Change is one chord of a sequence, with its bass and how long it
// sounds.
//
// The bass is part of it. It tells a perfect cadence from an imperfect
// one, a first degree in first inversion from a third degree, and a
// chromatic passing chord from a leap: an analysis that dropped it
// could not see half of what the book teaches.
type Change struct {
	Chord harmony.Chord
	Bass  harmony.PitchClass // the root, unless the chord is inverted

	// Silent is a change where no chord sounds: "N.C." on a chart, a
	// chord no one could read, the rest before a pickup. Its Chord is
	// meaningless.
	Silent bool

	Start  Ticks // from the start of the sequence
	Length Ticks
}

// Inverted reports whether the bass is not the root.
func (c Change) Inverted() bool {
	return !c.Silent && c.Bass != c.Chord.Root
}

// Changes are the chords of a tune or of a performance, in the order
// they sound: what the analysis of a sequence reads.
//
// # Two sources, one type
//
// A chart gives them all at once, and they loop: the app plays chorus
// after chorus, so the last chord leads to the first, and a turnaround
// is read as the preparation of the start. A player at the keyboard
// gives them one by one, and the sequence only grows: what follows the
// last chord is not known yet. Loops tells the two apart, and every
// analysis runs the same on both.
//
// # The coda
//
// A coda concludes a tune and is played on the last chorus only: the
// chorus loops back from the change before it, and after the coda
// nothing follows.
type Changes struct {
	Chords []Change
	Loops  bool
	Coda   int     // the first change of the coda, 0 when there is none
	Bars   []Ticks // where each bar starts, when the source has bars
}

// Next returns the index of the change that follows change i, or -1
// when none does: after the coda, or after the last change of a
// sequence that does not loop.
func (c Changes) Next(i int) int {
	last := len(c.Chords) - 1
	switch {
	case i < 0 || i > last:
		return -1
	case c.Coda > 0 && i == c.Coda-1:
		return 0
	case i < last:
		return i + 1
	case c.Loops && c.Coda == 0:
		return 0
	}
	return -1
}

// Prev returns the index of the change that comes before change i, or
// -1 when none does: before the first change of a sequence that does
// not loop. The first change of a looping sequence comes after the last
// change of the chorus, the one before the coda when there is one.
func (c Changes) Prev(i int) int {
	switch {
	case i < 0 || i >= len(c.Chords):
		return -1
	case i > 0:
		return i - 1
	case !c.Loops:
		return -1
	case c.Coda > 0:
		return c.Coda - 1
	}
	return len(c.Chords) - 1
}

// Bar returns the index of the bar a tick falls in, or -1 when the
// sequence has no bars.
func (c Changes) Bar(at Ticks) int {
	bar := -1
	for i, start := range c.Bars {
		if start > at {
			break
		}
		bar = i
	}
	return bar
}
