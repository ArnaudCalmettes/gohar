package main

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The verdict and the advice of a review (see "Le bilan" in
// docs/walk.md): how the notes sat on the beat, in one phrase, and one
// priority to work on, never more.

// A formula is a sequence of chords a musician hears as one thing, as
// the analysis finds it in the grid: the review counts the beats it
// missed under the formula their chord is part of.
type formula int

const (
	noFormula formula = iota
	anatole           // I VI II V
	threeSix          // III VI II V I
	aeolian           // ♭VI ♭VII I
	twoFive           // II V, and its I when it resolves
	formulas
)

// formulasOf finds the formulas of `c`, by the start of their chords.
// The cells of the analysis come first: the II V of an anatole is part
// of the anatole. A V alone, without its II, is no formula.
func formulasOf(c analysis.Changes) map[analysis.Ticks]formula {
	out := map[analysis.Ticks]formula{}
	mark := func(i int, f formula) {
		if i < 0 || i >= len(c.Chords) {
			return
		}
		if _, ok := out[c.Chords[i].Start]; !ok {
			out[c.Chords[i].Start] = f
		}
	}
	for _, cell := range analysis.Cells(c) {
		f := map[analysis.CellKind]formula{
			analysis.Anatole:            anatole,
			analysis.ThreeSixTwoFiveOne: threeSix,
			analysis.AeolianCadence:     aeolian,
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
			mark(i, twoFive)
		}
	}
	return out
}

// A feel is how the notes of a run sat on the beat, in one word.
type feel int

const (
	feelNone     feel = iota // no note to measure
	feelSteady               // close to the beat and steady
	feelRushing              // ahead of it, on average
	feelDragging             // behind it
	feelUnsteady             // close on average, but scattered
)

// The thresholds of the verdict: starting values, to set by playing,
// as the walker's. The latency is already taken off the notes.
const (
	steadyMean   = 20 * time.Millisecond // the mean offset, either way
	steadySpread = 30 * time.Millisecond // the spread around it
)

// feel says how the notes of the run sat on the beat.
func (s *summary) feel() feel {
	mean, spread, ok := s.offset()
	switch {
	case !ok:
		return feelNone
	case mean < -steadyMean:
		return feelRushing
	case mean > steadyMean:
		return feelDragging
	case spread > steadySpread:
		return feelUnsteady
	}
	return feelSteady
}

// An advice is the one priority a review gives.
type advice struct {
	kind adviceKind
	what formula   // what to work, when a formula
	sit  situation // or a situation, when no formula is to consolidate
	bpm  float64   // the tempo to try, when slower or faster
}

type adviceKind int

const (
	noAdvice    adviceKind = iota
	workFormula            // work a formula at this tempo
	workSit                // work a situation at this tempo
	slower                 // the time and the changes slip: slow down
	stay                   // the changes hold, the time does not: stay
	faster                 // everything worked: speed up
)

// advise gives the one priority of a run played at `bpm`: the weakest
// formula when only the changes slip, a slower tempo when the time
// slips too, the same tempo when only the time does, a faster one when
// everything worked.
func advise(s *summary, bpm float64) advice {
	steady := s.feel() == feelSteady || s.feel() == feelNone
	f, fok := s.weakestFormula()
	sit, sok := s.weakestSituation()
	switch {
	case (fok || sok) && !steady:
		if t := slowerThan(bpm); t < bpm {
			return advice{kind: slower, bpm: t}
		}
	case fok:
		return advice{kind: workFormula, what: f}
	case sok:
		return advice{kind: workSit, sit: sit}
	case !steady:
		return advice{kind: stay}
	}
	if fok || sok {
		return advice{} // already at the slowest
	}
	if t := fasterThan(bpm); t > bpm {
		return advice{kind: faster, bpm: t}
	}
	return advice{}
}

// fasterThan is the next tempo to try: 10 more under 140, 5 from
// there, never past the fastest.
func fasterThan(bpm float64) float64 {
	step := 10.0
	if bpm >= 140 {
		step = 5
	}
	return min(bpm+step, maxBPM)
}

// slowerThan is the same step down, never under the slowest.
func slowerThan(bpm float64) float64 {
	step := 10.0
	if bpm > 140 {
		step = 5
	}
	return max(bpm-step, minBPM)
}

// weakestFormula is the formula to consolidate that missed the most
// beats, false when none is to consolidate.
func (s *summary) weakestFormula() (formula, bool) {
	best, missed := noFormula, 0
	for f := anatole; f < formulas; f++ {
		t := s.formulas[f]
		if m := t.expected - t.landed; t.expected > 0 && !t.worked() && m > missed {
			best, missed = f, m
		}
	}
	return best, missed > 0
}

// weakestSituation is the same among the situations.
func (s *summary) weakestSituation() (situation, bool) {
	var best situation
	missed := 0
	for sit, t := range s.situations {
		if m := t.expected - t.landed; t.expected > 0 && !t.worked() && m > missed {
			best, missed = situation(sit), m
		}
	}
	return best, missed > 0
}
