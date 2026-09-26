package main

import (
	"math/rand/v2"
	"slices"
	"time"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// What every activity shares, kept apart from Ebitengine so that it can
// be tested without a window. The rules are in docs/oreille.md; the
// comments only say where they land.

const gameID dex.GameID = "ear"

// retryGap is how many questions come between a mistake and its retry:
// two, so that the retry lands third.
const retryGap = 2

// A Question is one thing to hear and the answers offered for it.
//
// The answers are notions of the dex, so that the game names them in
// any language and notation without the activity knowing either. What
// sounds is the right one, Choices[Answer]: the activity keeps no
// other secret.
//
// # Answering by playing
//
// Only by choosing, for now. Answering by playing the note on the MIDI
// keyboard will come with the degrees: the first note played after the
// question is then the answer, and a right one is also a production.
type Question struct {
	Tonic   harmony.PitchClass
	Choices []dex.Notion
	Answer  int

	// Retry marks a question that came back after a mistake. It does
	// not change the facts it produces: the dex is not told.
	Retry bool
}

// Right returns the notion that sounds.
func (q Question) Right() dex.Notion {
	return q.Choices[q.Answer]
}

// A display says what the keyboard may show once the answer is out:
// the right shape on the tonic, and the one chosen by mistake if any.
// Patterns rather than sets, so that the notes can be spelled in the
// shape they belong to.
type display struct {
	tonic   harmony.PitchClass
	right   harmony.ScalePattern
	chosen  harmony.ScalePattern
	mistake bool
}

// An Activity is what differs from one exercise to the next: what to
// sound, what to offer, what to learn from a right answer.
//
// What does not differ lives in Series: the retry after a mistake, the
// facts sent only on a right answer, the one report at the end. Those
// are the principles of the game, not of an exercise.
type Activity interface {
	// ID is what the report names the activity.
	ID() string

	// Plan draws the questions of a series. The whole series at once,
	// so that an activity can make sure it covers its ground before
	// repeating itself.
	Plan(rng *rand.Rand, n int) []Question

	// Retry is the question that comes back after a mistake: the same
	// notion, asked again differently. What "differently" means is the
	// activity's call: another tonic for a mode, the same key for a
	// degree.
	Retry(q Question, rng *rand.Rand) Question

	// Prompt is the sentence that asks the question.
	Prompt(l language, q Question) string

	// Sound is what the player hears when the question is asked.
	Sound(q Question) []keyboard.Note

	// Correction is what the player hears after choosing `chosen` by
	// mistake: the right answer, then the wrong one, to compare.
	Correction(q Question, chosen int) []keyboard.Note

	// Show is what the keyboard shows once the answer is out.
	Show(q Question, chosen int) display

	// Facts are what a right answer teaches the dex.
	Facts(q Question) []dex.Fact
}

// A Series is a planned run of questions from one activity, plus at
// most one retry per mistake.
type Series struct {
	activity  Activity
	rng       *rand.Rand
	questions []Question
	current   int
	facts     []dex.Fact
}

func NewSeries(a Activity, rng *rand.Rand, length int) *Series {
	return &Series{activity: a, rng: rng, questions: a.Plan(rng, length)}
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

// Answer records the choice `i` and reports whether it was right.
//
// A right answer, first time or retry, sends the activity's facts. A
// wrong one is nothing at all: it only schedules the retry, once.
func (s *Series) Answer(i int) bool {
	q, ok := s.Current()
	if !ok {
		return false
	}
	if i == q.Answer {
		s.facts = append(s.facts, s.activity.Facts(q)...)
		return true
	}
	if !q.Retry {
		retry := s.activity.Retry(q, s.rng)
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

// Report is the one report of the series. Facts from a series left
// halfway are still facts, which is why this does not check that the
// series ended.
func (s *Series) Report(at time.Time) dex.Report {
	return dex.Report{Game: gameID, Activity: s.activity.ID(), At: at, Facts: s.facts}
}

// otherTonic draws a tonic different from `t`.
func otherTonic(t harmony.PitchClass, rng *rand.Rand) harmony.PitchClass {
	return (t + 1 + harmony.PitchClass(rng.IntN(11))) % 12
}
