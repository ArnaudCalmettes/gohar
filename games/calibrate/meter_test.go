package calibrate

import (
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// At 120 to the quarter note, a beat lasts 500 ms: the fourth beat of
// bar `bar`, counted from 1, falls at 1.5 s, 3.5 s...
var start = time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)

func fourth(m tempo.Metronome, bar int) time.Time {
	return m.At((bar-1)*m.PerBar() + m.PerBar() - 1)
}

func TestTapAgainstTheFourthBeat(t *testing.T) {
	m := tempo.NewMetronome(start, 120, 4)
	for _, c := range []struct {
		name string
		at   time.Time
		off  time.Duration
		ok   bool
	}{
		{"on the fourth beat of bar 1", fourth(m, 1), 0, true},
		{"40 ms late on bar 3", fourth(m, 3).Add(40 * time.Millisecond), 40 * time.Millisecond, true},
		{"30 ms early on bar 2", fourth(m, 2).Add(-30 * time.Millisecond), -30 * time.Millisecond, true},
		{"on beat 2 of bar 2, no answer", m.At(5), -time.Second, false},
		{"before the first bar", start.Add(-2 * time.Second), 500 * time.Millisecond, false},
	} {
		off, ok := NewMeter(m).Tap(c.at)
		if ok != c.ok || (ok && off != c.off) {
			t.Errorf("%s: got %v, %v; want %v, %v", c.name, off, ok, c.off, c.ok)
		}
	}
}

// A player 60 ms late, give or take 10 ms: steady after eight taps,
// never before.
func TestSteadyAfterAWindow(t *testing.T) {
	m := tempo.NewMetronome(start, 120, 4)
	k := NewMeter(m)
	spread := []time.Duration{50, 70, 60, 55, 65, 60, 50, 70}
	for i, s := range spread {
		if _, ok := k.Mean(); ok {
			t.Fatalf("steady after %d taps", i)
		}
		k.Tap(fourth(m, i+1).Add(s * time.Millisecond))
	}
	mean, ok := k.Mean()
	if !ok || mean != 60*time.Millisecond {
		t.Errorf("got %v, %v; want 60ms, steady", mean, ok)
	}
}

// A player all over the place never steadies: the scene keeps beating.
func TestUnsteady(t *testing.T) {
	m := tempo.NewMetronome(start, 120, 4)
	k := NewMeter(m)
	for i := range 2 * window {
		s := 0 * time.Millisecond
		if i%2 == 0 {
			s = 150 * time.Millisecond
		}
		k.Tap(fourth(m, i+1).Add(s))
	}
	if _, ok := k.Mean(); ok {
		t.Error("steady, with taps 150 ms apart")
	}
}
