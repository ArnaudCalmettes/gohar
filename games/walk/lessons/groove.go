package lessons

import (
	"math"

	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// onBeat is how far from the beat a note still lands it, in beats: the
// loose window of the first palier, a third of a beat.
var onBeat = mark.FirstPalier.Loose

// beatsPerBar is the meter of the chapter: four in a bar, always.
const beatsPerBar = mark.BeatsPerBar

// Groove is "Jouer en rythme": over the pulse, the player plays on the
// beats `On` of each bar (1 for the downbeat), `Bars` bars in a row, a
// key of his choice, or `Want`, the note of each bar, in any octave,
// when given; `Chords`, when given, shows the grid those notes come
// from, a bar a string ("C6 A7" for two chords), the row of the bar
// playing, the bar shaded.
//
// The series starts when the player starts it: on a 1 of his choice,
// played on time, the right note when one is wanted. Before, his notes
// count for nothing, but a wrong one on a 1 is marked, so that he sees
// why nothing starts. The pulse never stops (see "Ne jamais casser la
// pulse" in docs/architecture.md): in the series, a note off the beat,
// a wrong one, a beat left empty, is marked on the count, discreetly,
// and the series goes on. At its end, the step is over if nothing was
// missed; otherwise the walker says `Again`, and the series waits for
// a 1 again.
type Groove struct {
	still
	Phrase string
	Again  string
	On     []int
	Bars   int
	Want   []Target
	Chords []string

	first  int          // the series' first beat; -1 while it waits for a 1
	wanted map[int]int  // the beats to play, and their bar
	marked map[int]bool // the beats marked, landed or missed
	landed map[int]bool // those landed
	missed bool
}

func (st *Groove) said() string { return st.Phrase }

func (st *Groove) alsoSays() []string { return []string{st.Again} }

func (st *Groove) begin(s Stage) {
	s.Say(st.Phrase)
	if st.Want != nil {
		s.Bass(true) // the roots, on the double bass
	}
	if st.Chords != nil {
		s.Chords(nil)
	}
	st.wait(s)
}

// wait has the series wait for its first 1.
func (st *Groove) wait(s Stage) {
	st.first = -1
	if st.Chords != nil {
		s.Line(st.Chords, -1)
	}
}

// start lays the series out from beat `first`, a downbeat.
func (st *Groove) start(s Stage, first int) {
	st.first, st.missed = first, false
	st.wanted, st.marked, st.landed = map[int]int{}, map[int]bool{}, map[int]bool{}
	for bar := range st.Bars {
		for _, b := range st.On {
			st.wanted[first+bar*beatsPerBar+b-1] = bar
		}
	}
	if st.Chords != nil {
		s.Line(st.Chords, 0)
	}
}

func (st *Groove) beat(s Stage, n int) bool {
	if st.first < 0 {
		return false
	}
	// A beat to play gone by unplayed, its window closed.
	for b := range st.wanted {
		if b < n && !st.marked[b] {
			st.mark(s, b, false)
		}
	}
	end := st.first + st.Bars*beatsPerBar
	if bar := (n - st.first) / beatsPerBar; n > st.first && n < end && n%beatsPerBar == 0 && st.Chords != nil {
		s.Line(st.Chords, bar)
	}
	if n < end {
		return false
	}
	if !st.missed {
		return true
	}
	s.Say(st.Again)
	st.wait(s)
	return false
}

func (st *Groove) pressed(s Stage, held []int) bool {
	n, off := s.Timing()
	right := st.Want == nil || st.Want[0].Judge(held) == Hit
	if st.first < 0 {
		switch {
		case n%beatsPerBar != 0 || math.Abs(off) > onBeat:
			// Off the 1: nothing starts, nothing counts.
		case right:
			st.start(s, n)
			st.mark(s, n, true)
		default:
			s.MarkBeat(n, false) // a wrong note on the 1: why nothing starts
		}
		return false
	}
	switch {
	case n < st.first || n >= st.first+st.Bars*beatsPerBar:
		return false // after the series: free
	case st.marked[n] && !st.landed[n]:
		return false // missed already
	case st.marked[n]:
		st.mark(s, n, false) // a second note on a beat landed
		return false
	}
	bar, wanted := st.wanted[n]
	ok := wanted && math.Abs(off) <= onBeat
	if ok && st.Want != nil {
		ok = st.Want[bar].Judge(held) == Hit
	}
	st.mark(s, n, ok)
	return false
}

// mark marks beat `n`, landed or not, on the stage and in the series.
func (st *Groove) mark(s Stage, n int, landed bool) {
	st.marked[n], st.landed[n] = true, landed
	if !landed {
		st.missed = true
	}
	s.MarkBeat(n, landed)
}
