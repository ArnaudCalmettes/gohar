package tempo

import (
	"testing"
	"time"
)

var t0 = time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)

// At 120, a beat lasts half a second: easy to check by hand.
func at120() Metronome { return NewMetronome(t0, 120, 4) }

func TestMetronomeAt(t *testing.T) {
	m := at120()
	for _, tc := range []struct {
		name string
		n    int
		want time.Duration
	}{
		{"bar 1, beat 1", 0, 0},
		{"bar 1, beat 3", 2, time.Second},
		{"bar 2, beat 1", 4, 2 * time.Second},
		{"count-in, beat 1", -4, -2 * time.Second},
	} {
		if got := m.At(tc.n).Sub(t0); got != tc.want {
			t.Errorf("%s: at %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestMetronomePosition(t *testing.T) {
	m := at120()
	for _, tc := range []struct {
		n      int
		want   Position
		strong bool
	}{
		{0, Position{1, 1}, true},
		{1, Position{1, 2}, false},
		{2, Position{1, 3}, true},
		{3, Position{1, 4}, false},
		{4, Position{2, 1}, true},
		{47, Position{12, 4}, false}, // last beat of a blues
		{-4, Position{0, 1}, true},   // count-in
		{-1, Position{0, 4}, false},
	} {
		got := m.Position(tc.n)
		if got != tc.want {
			t.Errorf("beat %d: %+v, want %+v", tc.n, got, tc.want)
		}
		if m.Strong(got) != tc.strong {
			t.Errorf("beat %d (bar %d, beat %d): strong %v", tc.n, got.Bar, got.Beat, !tc.strong)
		}
	}
}

func TestMetronomeWaltzHasOneStrongBeat(t *testing.T) {
	m := NewMetronome(t0, 120, 3)
	if m.Strong(Position{1, 3}) {
		t.Error("beat 3 of a waltz is weak")
	}
	if got := m.Position(3); got != (Position{2, 1}) {
		t.Errorf("beat 3 of a waltz: %+v, want bar 2, beat 1", got)
	}
}

func TestMetronomeNearest(t *testing.T) {
	m := at120()
	ms := time.Millisecond
	for _, tc := range []struct {
		name  string
		after time.Duration
		n     int
		off   time.Duration
	}{
		{"on bar 1, beat 1", 0, 0, 0},
		{"20 ms late on beat 2", 520 * ms, 1, 20 * ms},
		{"30 ms early on beat 3", 970 * ms, 2, -30 * ms},
		{"10 ms early on bar 1, beat 1", -10 * ms, 0, -10 * ms},
		{"40 ms late on the last count-in beat", -460 * ms, -1, 40 * ms},
	} {
		n, off := m.Nearest(t0.Add(tc.after))
		if n != tc.n || off != tc.off {
			t.Errorf("%s: beat %d off %v, want beat %d off %v", tc.name, n, off, tc.n, tc.off)
		}
	}
}

func TestMetronomeBeats(t *testing.T) {
	m := at120()
	if got := m.Beats(t0.Add(1250 * time.Millisecond)); got != 2.5 {
		t.Errorf("halfway between beats 3 and 4 of bar 1: %v, want 2.5", got)
	}
}

// Windows laid end to end must hand out every beat exactly once, a
// beat on the edge included.
func TestMetronomeDueWindows(t *testing.T) {
	m := at120()
	window := 300 * time.Millisecond
	seen := map[int]int{}
	for from := t0.Add(-2 * time.Second); from.Before(t0.Add(4 * time.Second)); from = from.Add(window) {
		first, end := m.Due(from, from.Add(window))
		for n := first; n < end; n++ {
			seen[n]++
		}
	}
	for n := -4; n < 8; n++ {
		if seen[n] != 1 {
			t.Errorf("beat %d handed out %d times", n, seen[n])
		}
	}
}
