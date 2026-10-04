package main

import (
	"math"
	"time"
)

// The summary of a run (see "Le bilan" in docs/walk.md): what worked
// and what is to consolidate, by situation and bar by bar, and how the
// notes sat on the beat. It takes note, as the marker does, and never
// grades the player: the counts are said as they are.

// A situation is where in the grid a beat expects its note.
type situation int

const (
	barStart situation = iota // a chord arrives on beat 1
	midBar                    // a chord arrives within the bar: two share it
	held                      // beat 1 of a chord that carries on
	situations
)

// situationOf says where `b` expects its note, false for a beat that
// expects none.
func situationOf(b Beat) (situation, bool) {
	switch {
	case b.Arrives && b.Position.Beat == 1:
		return barStart, true
	case b.Arrives:
		return midBar, true
	case b.Holds:
		return held, true
	}
	return 0, false
}

// A tally counts the beats that expected a note, and those landed.
type tally struct {
	landed, expected int
}

// worked is the share of landed beats from which a situation or a bar
// has worked; below it, it is to consolidate. A starting value, to set
// by playing.
const worked = 0.75

// worked says the tally has worked; false too when nothing was
// expected.
func (t tally) worked() bool {
	return t.expected > 0 && float64(t.landed) >= worked*float64(t.expected)
}

// A summary gathers the marks of one run.
type summary struct {
	situations [situations]tally
	bars       []tally // by bar of the chart, every chorus together
	chorusLen  int     // in beats

	notes      int     // the notes a beat claimed
	sum, sumSq float64 // of their offsets, in milliseconds
}

// newSummary sums up a run over a chart of `bars` bars, in choruses of
// `chorusLen` beats.
func newSummary(bars, chorusLen int) *summary {
	return &summary{bars: make([]tally, bars), chorusLen: chorusLen}
}

// beat takes the mark `k` of beat `b`, the beat `n` of the run.
func (s *summary) beat(n int, b Beat, k BeatKind) {
	sit, ok := situationOf(b)
	if !ok {
		return
	}
	i, _ := barOf(n, s.chorusLen)
	bar := &s.bars[i]
	s.situations[sit].expected++
	bar.expected++
	if k == Landed {
		s.situations[sit].landed++
		bar.landed++
	}
}

// note takes the mark of a note. A note between two beats has no
// offset to measure.
func (s *summary) note(m NoteMark) {
	if m.Timing == Between {
		return
	}
	ms := float64(m.Off) / float64(time.Millisecond)
	s.notes++
	s.sum += ms
	s.sumSq += ms * ms
}

// offset returns the mean offset of the notes, negative when early, and
// their spread around it, as a standard deviation; false without a note.
func (s *summary) offset() (mean, spread time.Duration, ok bool) {
	if s.notes == 0 {
		return 0, 0, false
	}
	n := float64(s.notes)
	m := s.sum / n
	v := max(s.sumSq/n-m*m, 0)
	return time.Duration(m * float64(time.Millisecond)), time.Duration(math.Sqrt(v) * float64(time.Millisecond)), true
}
