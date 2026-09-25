package main

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

func newTestSeries() *Series {
	return NewSeries(rand.New(rand.NewPCG(1, 2)), 10)
}

func TestSeriesCoversTheSystem(t *testing.T) {
	s := newTestSeries()
	seen := map[harmony.Degree]bool{}
	for _, q := range s.questions {
		seen[q.Mode] = true
	}
	if len(s.questions) != 10 || len(seen) != 7 {
		t.Errorf("%d questions over %d modes, want 10 over 7", len(s.questions), len(seen))
	}
}

func TestChoicesAreFarAndDistinct(t *testing.T) {
	s := newTestSeries()
	for i, q := range s.questions {
		seen := map[harmony.Degree]bool{}
		hasAnswer := false
		for _, c := range q.Choices {
			if seen[c] {
				t.Errorf("question %d offers %d twice", i, c)
			}
			seen[c] = true
			if c == q.Mode {
				hasAnswer = true
				continue
			}
			if distance(q.Mode, c) < farEnough {
				t.Errorf("question %d: %d is too close to %d", i, c, q.Mode)
			}
		}
		if !hasAnswer {
			t.Errorf("question %d does not offer its answer", i)
		}
	}
}

// Lydian and ionian differ by one note, which is what keeps them apart
// at this level.
func TestDistance(t *testing.T) {
	if d := distance(1, 4); d != 1 {
		t.Errorf("ionian to lydian: %d, want 1", d)
	}
	if d := distance(4, 7); d != 5 {
		t.Errorf("lydian to locrian: %d, want 5, the tonic and the tritone shared", d)
	}
}

func wrongChoice(q Question) harmony.Degree {
	for _, c := range q.Choices {
		if c != q.Mode {
			return c
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
	if !retry.Retry || retry.Mode != q.Mode || retry.Tonic == q.Tonic {
		t.Errorf("retry %+v for %+v: same mode, another tonic expected", retry, q)
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
		s.Answer(q.Mode)
		s.Next()
	}

	r := s.Report(time.Now())
	if r.Game != gameID || len(r.Facts) != 20 {
		t.Fatalf("report %s with %d facts, want %s with 20", r.Game, len(r.Facts), gameID)
	}
	for _, f := range r.Facts {
		if f.Kind != dex.FactNamed && f.Kind != dex.FactHeard {
			t.Errorf("unexpected fact %v", f.Kind)
		}
	}
}
