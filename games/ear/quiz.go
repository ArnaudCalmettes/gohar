package main

import (
	"math/rand/v2"
	"slices"
	"time"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// The rules of the activity, kept apart from Ebitengine so that they
// can be tested without a window. Everything here is decided in
// OREILLE.md; the comments only say where.

const (
	gameID   dex.GameID = "ear"
	activity            = "natural-modes"

	// system is the only one this first activity plays.
	system = harmony.NaturalMajor

	// retryGap is how many questions come between a mistake and its
	// retry: two, so that the retry lands third.
	retryGap = 2

	// farEnough is the least number of notes a distractor must differ
	// by, at the same tonic. Lydian against ionian differs by one and
	// is kept for a later level.
	farEnough = 2
)

// A Question is one mode on one tonic, and the three names offered.
type Question struct {
	Tonic   harmony.PitchClass
	Mode    harmony.Degree
	Choices [3]harmony.Degree

	// Retry marks a question that came back after a mistake. It does
	// not change the facts it produces: the dex is not told.
	Retry bool
}

// A Series is ten questions, plus at most one retry per mistake.
type Series struct {
	rng       *rand.Rand
	questions []Question
	current   int
	facts     []dex.Fact
}

// NewSeries draws a series. The seven modes all come once before any
// comes twice, so that ten questions cover the whole system.
func NewSeries(rng *rand.Rand, length int) *Series {
	s := &Series{rng: rng}
	var order []harmony.Degree
	for len(order) < length {
		perm := rng.Perm(7)
		for _, i := range perm {
			order = append(order, harmony.Degree(i+1))
		}
	}
	for _, mode := range order[:length] {
		s.questions = append(s.questions, s.question(mode, harmony.PitchClass(rng.IntN(12))))
	}
	return s
}

func (s *Series) question(mode harmony.Degree, tonic harmony.PitchClass) Question {
	var far []harmony.Degree
	for d := harmony.Degree(1); d <= 7; d++ {
		if d != mode && distance(mode, d) >= farEnough {
			far = append(far, d)
		}
	}
	s.rng.Shuffle(len(far), func(i, j int) { far[i], far[j] = far[j], far[i] })

	q := Question{Tonic: tonic, Mode: mode, Choices: [3]harmony.Degree{mode, far[0], far[1]}}
	s.rng.Shuffle(3, func(i, j int) { q.Choices[i], q.Choices[j] = q.Choices[j], q.Choices[i] })
	return q
}

// distance counts the notes two modes of the system differ by on the
// same tonic.
func distance(a, b harmony.Degree) int {
	pa, _ := system.Mode(a)
	pb, _ := system.Mode(b)
	return pa.At(0).Difference(pb.At(0)).Len()
}

// Current returns the question being asked, and false once the series
// is over.
func (s *Series) Current() (Question, bool) {
	if s.current >= len(s.questions) {
		return Question{}, false
	}
	return s.questions[s.current], true
}

// Position returns the number of the current question and how many
// there are now, retries included.
func (s *Series) Position() (int, int) {
	return s.current + 1, len(s.questions)
}

// Answer records the player's choice and reports whether it was right.
//
// A right answer, first time or retry, is a FactNamed and a FactHeard:
// the game names the mode as it confirms. A wrong one is nothing at
// all; it only schedules the retry, once, on another tonic.
func (s *Series) Answer(choice harmony.Degree) bool {
	q, ok := s.Current()
	if !ok {
		return false
	}

	if choice == q.Mode {
		n := dex.ModeOf(system, q.Mode)
		s.facts = append(s.facts,
			dex.Fact{Kind: dex.FactNamed, Notion: n, Tonic: q.Tonic},
			dex.Fact{Kind: dex.FactHeard, Notion: n, Tonic: q.Tonic},
		)
		return true
	}

	if !q.Retry {
		tonic := (q.Tonic + 1 + harmony.PitchClass(s.rng.IntN(11))) % 12
		retry := s.question(q.Mode, tonic)
		retry.Retry = true
		at := min(s.current+1+retryGap, len(s.questions))
		s.questions = slices.Insert(s.questions, at, retry)
	}
	return false
}

// Next moves on, and reports false once the series is over.
func (s *Series) Next() bool {
	s.current++
	return s.current < len(s.questions)
}

// Report is the one report of the activity. Facts from a series that
// was left halfway are still facts, which is why this does not check
// that the series ended.
func (s *Series) Report(at time.Time) dex.Report {
	return dex.Report{Game: gameID, Activity: activity, At: at, Facts: s.facts}
}
