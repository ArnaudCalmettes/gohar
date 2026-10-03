package main

import (
	"math"
	"testing"
	"time"
)

// At 120 a beat lasts 500 ms, so the "and" of bar 1, beat 1 falls at
// 500 ms times the ratio.
func TestSwingPlacesTheAnd(t *testing.T) {
	ms := time.Millisecond
	for _, tc := range []struct {
		name  string
		ratio float64
		want  time.Duration
	}{
		{"even eighths", 0.5, 250 * ms},
		{"triplet swing", 2.0 / 3, 333 * ms},
		{"New Orleans dotted swing", 0.75, 375 * ms},
	} {
		s := NewSwing(at120(), tc.ratio)
		got := s.AtBeats(0.5).Sub(t0).Round(ms)
		if got != tc.want {
			t.Errorf("%s: the \"and\" of beat 1 at %v, want %v", tc.name, got, tc.want)
		}
		if beat3 := s.AtBeats(2).Sub(t0); beat3 != time.Second {
			t.Errorf("%s: beat 3 moved to %v", tc.name, beat3)
		}
	}
}

// What the marker needs: a swung "and" reads back as an "and", and a
// position read back from its instant is the position, to the
// nanosecond an instant is rounded to (2e-9 beats at 120).
func TestSwingInverts(t *testing.T) {
	s := NewSwing(at120(), 2.0/3)
	for _, x := range []float64{-4, -3.5, 0, 0.25, 0.5, 0.9, 1, 2.5, 47.5} {
		if got := s.Beats(s.AtBeats(x)); math.Abs(got-x) > 1e-8 {
			t.Errorf("position %v reads back as %v", x, got)
		}
	}
	// Played at 333 ms, the triplet "and": exactly halfway, swung.
	if got := s.Beats(t0.Add(time.Second / 3)); math.Abs(got-0.5) > 1e-6 {
		t.Errorf("a note at 333 ms reads as %v beats, want the \"and\", 0.5", got)
	}
}

// The ride pattern of a swing tune, "ding, ding-a ding, ding-a": 1, 2,
// the "and" of 2, 3, 4, the "and" of 4. Windows laid end to end hand
// out each eighth once, swung like the rest.
func TestDueEveryEighth(t *testing.T) {
	s := NewSwing(at120(), 2.0/3)
	window := 70 * time.Millisecond
	seen := map[int]int{}
	for from := t0.Add(-time.Second); from.Before(t0.Add(3 * time.Second)); from = from.Add(window) {
		first, end := DueEvery(s, 2, from, from.Add(window))
		for k := first; k < end; k++ {
			seen[k]++
		}
	}
	for k := -4; k < 12; k++ {
		if seen[k] != 1 {
			t.Errorf("eighth %d (position %v) handed out %d times", k, float64(k)/2, seen[k])
		}
	}
}

func TestSwingRefusesAFlatRatio(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a ratio of 1 accepted")
		}
	}()
	NewSwing(at120(), 1)
}
