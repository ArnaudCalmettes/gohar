package mark

import (
	"math"
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The summary of a run (see "Le bilan" in docs/walk.md): what worked
// and what is to consolidate, by situation and bar by bar, and how the
// notes sat on the beat. It takes note, as the marker does, and never
// grades the player: the counts are said as they are.

// A Situation is where in the grid a beat expects its note.
type Situation int

const (
	BarStart Situation = iota // a chord arrives on beat 1
	MidBar                    // a chord arrives within the bar: two share it
	Held                      // beat 1 of a chord that carries on
	Situations
)

// situationOf says where `b` expects its note, false for a beat that
// expects none.
func situationOf(b Beat) (Situation, bool) {
	switch {
	case b.Arrives && b.Position.Beat == 1:
		return BarStart, true
	case b.Arrives:
		return MidBar, true
	case b.Holds:
		return Held, true
	}
	return 0, false
}

// A Tally counts the beats that expected a note, and those landed.
type Tally struct {
	Landed, Expected int
}

// workedShare is the share of landed beats from which a situation or a bar
// has worked; below it, it is to consolidate. A starting value, to set
// by playing.
const workedShare = 0.75

// Worked says the tally has worked; false too when nothing was
// expected.
func (t Tally) Worked() bool {
	return t.Expected > 0 && float64(t.Landed) >= workedShare*float64(t.Expected)
}

// A Summary gathers the marks of one run.
type Summary struct {
	BySituation [Situations]Tally
	ByFormula   [Formulas]Tally            // by formula, the beats out of any under noFormula
	formulaOf   map[analysis.Ticks]Formula // of each chord, by its start
	ByBar       []Tally                    // by bar of the chart, every chorus together
	chorusLen   int                        // in beats

	notes      int     // the notes a beat claimed
	sum, sumSq float64 // of their offsets, in milliseconds
}

// NewSummary sums up a run over a chart of `bars` bars, in choruses of
// `chorusLen` beats, whose chords make the formulas `f` (see FormulasOf).
func NewSummary(bars, chorusLen int, f map[analysis.Ticks]Formula) *Summary {
	return &Summary{ByBar: make([]Tally, bars), chorusLen: chorusLen, formulaOf: f}
}

// Beat takes the mark `k` of beat `b`, the beat `n` of the run.
func (s *Summary) Beat(n int, b Beat, k BeatKind) {
	sit, ok := situationOf(b)
	if !ok {
		return
	}
	i, _ := BarOf(n, s.chorusLen)
	bar := &s.ByBar[i]
	f := &s.ByFormula[s.formulaOf[b.Chord.Start]]
	s.BySituation[sit].Expected++
	bar.Expected++
	f.Expected++
	if k == Landed {
		s.BySituation[sit].Landed++
		bar.Landed++
		f.Landed++
	}
}

// Note takes the mark of a note. A note between two beats has no
// offset to measure.
func (s *Summary) Note(m NoteMark) {
	if m.Timing == Between {
		return
	}
	ms := float64(m.Off) / float64(time.Millisecond)
	s.notes++
	s.sum += ms
	s.sumSq += ms * ms
}

// Offset returns the mean offset of the notes, negative when early, and
// their spread around it, as a standard deviation; false without a note.
func (s *Summary) Offset() (mean, spread time.Duration, ok bool) {
	if s.notes == 0 {
		return 0, 0, false
	}
	n := float64(s.notes)
	m := s.sum / n
	v := max(s.sumSq/n-m*m, 0)
	return time.Duration(m * float64(time.Millisecond)), time.Duration(math.Sqrt(v) * float64(time.Millisecond)), true
}
