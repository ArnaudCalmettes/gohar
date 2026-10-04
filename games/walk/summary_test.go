package main

import (
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// A run of two choruses of the blues where every arrival lands but the
// D7 of bar 11, on beat 3, missed twice, and the G7 of bar 12 missed
// once: the changes within the bar are to consolidate, bar 11 too; bar
// 12, landed three times out of four, has worked.
func TestSummary(t *testing.T) {
	grid, err := readGrid(jazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	beats := Expect(grid, tempo.NewMetronome(time.Time{}, 120, perBar), 2)
	chorus := len(beats) / 2
	missed := map[int]bool{
		10*perBar + 2:          true, // bar 11, beat 3: D7
		chorus + 10*perBar + 2: true,
		chorus + 11*perBar:     true, // bar 12, beat 1: G7, second chorus
	}
	s := newSummary(12)
	for n, b := range beats {
		k := Landed
		if missed[n] {
			k = Missed
		}
		s.beat(n, b, k)
	}
	if !s.situations[barStart].worked() || !s.situations[held].worked() {
		t.Errorf("the changes on beat 1 and the chords held should have worked: %+v", s.situations)
	}
	if got := s.situations[midBar]; got != (tally{landed: 2, expected: 4}) || got.worked() {
		t.Errorf("the changes within the bar (D7, C7 in bars 11 and 12): %+v, want 2 landed of 4, to consolidate", got)
	}
	if s.bars[10].worked() {
		t.Errorf("bar 11 (F6 D7): %+v, want to consolidate", s.bars[10])
	}
	if !s.bars[11].worked() {
		t.Errorf("bar 12 (G7 C7): %+v, want worked", s.bars[11])
	}
	if s.bars[3].expected != 2 {
		t.Errorf("bar 4, F7 held: %d beats expected, want 2", s.bars[3].expected)
	}
}

// The notes ahead of the beat by 10 and 30 ms: 20 ms early on average,
// give or take 10. A note between two beats does not count.
func TestSummaryOffset(t *testing.T) {
	s := newSummary(12)
	if _, _, ok := s.offset(); ok {
		t.Error("an offset without a note")
	}
	for _, m := range []NoteMark{
		{Off: -10 * time.Millisecond, Timing: OnTime},
		{Off: -30 * time.Millisecond, Timing: Early},
		{Off: 400 * time.Millisecond, Timing: Between},
	} {
		s.note(m)
	}
	mean, spread, _ := s.offset()
	if mean != -20*time.Millisecond || spread != 10*time.Millisecond {
		t.Errorf("mean %v, spread %v; want -20ms, 10ms", mean, spread)
	}
}
