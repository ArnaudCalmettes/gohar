package main

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// The series rules hold whatever the activity; the modes stand in for
// all of them here.

func newTestSeries() *Series {
	return NewSeries(modes{system: harmony.NaturalMajor}, rand.New(rand.NewPCG(1, 2)), 10)
}

func wrongChoice(q Question) int {
	for i := range q.Choices {
		if i != q.Answer {
			return i
		}
	}
	panic("no wrong choice")
}

func TestAMistakeIsNothingButARetry(t *testing.T) {
	s := newTestSeries()
	q, _ := s.Current()

	if s.Answer(wrongChoice(q)) {
		t.Fatal("a wrong choice was accepted")
	}
	if len(s.facts) != 0 {
		t.Errorf("a mistake produced facts: %v", s.facts)
	}

	_, total := s.Position()
	if total != 11 {
		t.Fatalf("%d questions after a mistake, want 11", total)
	}
	retry := s.questions[1+retryGap]
	if !retry.Retry || retry.Right() != q.Right() {
		t.Errorf("retry %+v for %+v: the same notion expected", retry, q)
	}

	t.Run("a retry failed again does not come back", func(t *testing.T) {
		for s.Next() {
			if c, _ := s.Current(); c.Retry {
				s.Answer(wrongChoice(c))
				break
			}
		}
		if _, total := s.Position(); total != 11 {
			t.Errorf("%d questions, want 11", total)
		}
	})
}

func TestTheCorrectionIsWhatCounts(t *testing.T) {
	s := newTestSeries()
	for {
		q, ok := s.Current()
		if !ok {
			break
		}
		s.Answer(q.Answer)
		s.Next()
	}

	r := s.Report(time.Now())
	if r.Game != gameID || r.Activity != "modes" || len(r.Facts) != 20 {
		t.Fatalf("report %s/%s with %d facts, want %s/modes with 20",
			r.Game, r.Activity, len(r.Facts), gameID)
	}
	for _, f := range r.Facts {
		if f.Kind != dex.FactNamed && f.Kind != dex.FactHeard {
			t.Errorf("unexpected fact %v", f.Kind)
		}
	}
}
