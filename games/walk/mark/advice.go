package mark

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The verdict and the advice of a review (see "Le bilan" in
// docs/walk.md): how the notes sat on the beat, in one phrase, and one
// priority to work on, never more.

// A Formula is a sequence of chords a musician hears as one thing, as
// the analysis finds it in the grid: the review counts the beats it
// missed under the formula their chord is part of.
type Formula int

const (
	NoFormula Formula = iota
	Anatole           // I VI II V
	ThreeSix          // III VI II V I
	Aeolian           // ♭VI ♭VII I
	TwoFive           // II V, and its I when it resolves
	Formulas
)

// FormulasOf finds the formulas of `c`, by the start of their chords.
// The cells of the analysis come first: the II V of an anatole is part
// of the anatole. A V alone, without its II, is no formula; a II V is
// one, whether or not its V resolves.
func FormulasOf(c analysis.Changes) map[analysis.Ticks]Formula {
	out := map[analysis.Ticks]Formula{}
	mark := func(i int, f Formula) {
		if i < 0 || i >= len(c.Chords) {
			return
		}
		if _, ok := out[c.Chords[i].Start]; !ok {
			out[c.Chords[i].Start] = f
		}
	}
	for _, cell := range analysis.Cells(c) {
		f := map[analysis.CellKind]Formula{
			analysis.Anatole:            Anatole,
			analysis.ThreeSixTwoFiveOne: ThreeSix,
			analysis.AeolianCadence:     Aeolian,
		}[cell.Kind]
		for i := cell.From; i <= cell.To; i++ {
			mark(i, f)
		}
	}
	for _, b := range analysis.Blocks(c, analysis.Approaches(c)) {
		if b.Two < 0 || b.Kind == harmony.PlagalApproach {
			continue
		}
		for _, i := range []int{b.Two, b.Sus, b.Five, b.Target} {
			mark(i, TwoFive)
		}
	}
	// A II V whose V prepares nothing is no block, but a musician hears
	// it all the same: the first Dm7 G7 of Dm7 G7 | Dm7 G7 in Satin Doll,
	// whose G7 goes back to its two.
	for i := 0; i+1 < len(c.Chords); i++ {
		two, five := c.Chords[i], c.Chords[i+1]
		if !two.Silent && !five.Silent && isTwo(two.Chord.Pattern) &&
			five.Chord.Pattern.HasAll(harmony.IntMajorThird, harmony.IntMinorSeventh) &&
			five.Chord.Root == two.Chord.Root.Transpose(5) {
			mark(i, TwoFive)
			mark(i+1, TwoFive)
		}
	}
	return out
}

// isTwo tells the qualities of a two: a minor seventh, or a half
// diminished seventh, as harmony reads them.
func isTwo(p harmony.ChordPattern) bool {
	t := p.Tetrad()
	return t == harmony.ChordMinorSeventh || t == harmony.ChordMinorSeventhNo5 || t == harmony.ChordHalfDiminished
}

// A Feel is how the notes of a run sat on the beat, in one word.
type Feel int

const (
	FeelNone     Feel = iota // no note to measure
	FeelSteady               // close to the beat and steady
	FeelRushing              // ahead of it, on average
	FeelDragging             // behind it
	FeelUnsteady             // close on average, but scattered
)

// The thresholds of the verdict: starting values, to set by playing,
// as the walker's. The latency is already taken off the notes.
const (
	steadyMean   = 20 * time.Millisecond // the mean offset, either way
	steadySpread = 30 * time.Millisecond // the spread around it
)

// Feel says how the notes of the run sat on the beat.
func (s *Summary) Feel() Feel {
	mean, spread, ok := s.Offset()
	switch {
	case !ok:
		return FeelNone
	case mean < -steadyMean:
		return FeelRushing
	case mean > steadyMean:
		return FeelDragging
	case spread > steadySpread:
		return FeelUnsteady
	}
	return FeelSteady
}

// An Advice is the one priority a review gives.
type Advice struct {
	Kind      AdviceKind
	What      Formula   // what to work, when a formula
	Situation Situation // or a situation, when no formula is to consolidate
	BPM       float64   // the tempo to try, when slower or faster
}

// An AdviceKind is what an Advice asks for.
type AdviceKind int

const (
	NoAdvice      AdviceKind = iota
	WorkFormula              // work a formula at this tempo
	WorkSituation            // work a situation at this tempo
	Slower                   // the time and the changes slip: slow down
	Stay                     // the changes hold, the time does not: stay
	Faster                   // everything worked: speed up
)

// Advise gives the one priority of a run played at `bpm`, the
// options' tempo bounded by `slowest` and `fastest`: the weakest
// formula when only the changes slip, a slower tempo when the time
// slips too, the same tempo when only the time does, a faster one when
// everything worked.
func Advise(s *Summary, bpm, slowest, fastest float64) Advice {
	steady := s.Feel() == FeelSteady || s.Feel() == FeelNone
	f, fok := s.weakestFormula()
	sit, sok := s.weakestSituation()
	switch {
	case (fok || sok) && !steady:
		if t := slowerThan(bpm, slowest); t < bpm {
			return Advice{Kind: Slower, BPM: t}
		}
	case fok:
		return Advice{Kind: WorkFormula, What: f}
	case sok:
		return Advice{Kind: WorkSituation, Situation: sit}
	case !steady:
		return Advice{Kind: Stay}
	}
	if fok || sok {
		return Advice{} // already at the slowest
	}
	if t := fasterThan(bpm, fastest); t > bpm {
		return Advice{Kind: Faster, BPM: t}
	}
	return Advice{}
}

// fasterThan is the next tempo to try: 10 more under 140, 5 from
// there, never past the fastest.
func fasterThan(bpm, fastest float64) float64 {
	step := 10.0
	if bpm >= 140 {
		step = 5
	}
	return min(bpm+step, fastest)
}

// slowerThan is the same step down, never under the slowest.
func slowerThan(bpm, slowest float64) float64 {
	step := 10.0
	if bpm > 140 {
		step = 5
	}
	return max(bpm-step, slowest)
}

// weakestFormula is the formula to consolidate that missed the most
// beats, false when none is to consolidate.
func (s *Summary) weakestFormula() (Formula, bool) {
	best, missed := NoFormula, 0
	for f := Anatole; f < Formulas; f++ {
		t := s.ByFormula[f]
		if m := t.Expected - t.Landed; t.Expected > 0 && !t.Worked() && m > missed {
			best, missed = f, m
		}
	}
	return best, missed > 0
}

// weakestSituation is the same among the situations.
func (s *Summary) weakestSituation() (Situation, bool) {
	var best Situation
	missed := 0
	for sit, t := range s.BySituation {
		if m := t.Expected - t.Landed; t.Expected > 0 && !t.Worked() && m > missed {
			best, missed = Situation(sit), m
		}
	}
	return best, missed > 0
}
