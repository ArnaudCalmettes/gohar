package main

import (
	"math/rand/v2"
	"testing"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

var natural = modes{system: harmony.NaturalMajor}

func planModes() []Question {
	return natural.Plan(rand.New(rand.NewPCG(1, 2)), 10)
}

func TestModesCoverTheSystem(t *testing.T) {
	questions := planModes()
	seen := map[dex.Notion]bool{}
	for _, q := range questions {
		seen[q.Right()] = true
	}
	if len(questions) != 10 || len(seen) != 7 {
		t.Errorf("%d questions over %d modes, want 10 over 7", len(questions), len(seen))
	}
}

func TestModeChoicesAreFarAndDistinct(t *testing.T) {
	for i, q := range planModes() {
		seen := map[dex.Notion]bool{}
		for j, c := range q.Choices {
			if seen[c] {
				t.Errorf("question %d offers %v twice", i, c)
			}
			seen[c] = true
			if j != q.Answer && natural.distance(q.Right().Degree, c.Degree) < farEnough {
				t.Errorf("question %d: %v is too close to %v", i, c, q.Right())
			}
		}
		if len(q.Choices) != 3 {
			t.Errorf("question %d offers %d choices, want 3", i, len(q.Choices))
		}
	}
}

// Lydian and ionian differ by one note, which is what keeps them apart
// at this level.
func TestDistance(t *testing.T) {
	if d := natural.distance(1, 4); d != 1 {
		t.Errorf("ionian to lydian: %d, want 1", d)
	}
	if d := natural.distance(4, 7); d != 5 {
		t.Errorf("lydian to locrian: %d, want 5, the tonic and the tritone shared", d)
	}
}

// A mode comes back on another tonic: the ear has to find the colour,
// not remember the notes.
func TestModeRetryChangesTonic(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 4))
	q := planModes()[0]
	r := natural.Retry(q, rng)
	if r.Right() != q.Right() || r.Tonic == q.Tonic {
		t.Errorf("retry %+v of %+v: same mode, another tonic expected", r, q)
	}
}

// All seven, each in the place of its degree: the answer's index is
// its degree minus one.
func TestAllModesInTheOrderOfTheDegrees(t *testing.T) {
	all := modes{system: harmony.NaturalMajor, all: true}
	for i, q := range all.Plan(rand.New(rand.NewPCG(1, 2)), 10) {
		if len(q.Choices) != 7 {
			t.Fatalf("question %d offers %d choices, want 7", i, len(q.Choices))
		}
		for j, c := range q.Choices {
			if c.Degree != harmony.Degree(j+1) {
				t.Errorf("question %d: choice %d is degree %d", i, j+1, c.Degree)
			}
		}
		if int(q.Right().Degree) != q.Answer+1 {
			t.Errorf("question %d: answer %d is not degree %d", i, q.Answer+1, q.Right().Degree)
		}
	}
}
