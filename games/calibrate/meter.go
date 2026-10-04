// Package calibrate is the latency calibration, a scene shared by the
// games (see "La calibration" in docs/architecture.md).
//
// The method is Rhythm Paradise Groove's: a bar beaten "ta, ta, ta,
// TI", the player pressing a key on the TI, the fourth beat. The three
// beats before it set the pulse in the player's head, and the press on
// the fourth is all the more precise.
package calibrate

import (
	"math"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// When the measure is steady: the last `window` taps, spread by less
// than `steady`. A first guess of ours, to set with players.
const (
	window = 8
	steady = 20 * time.Millisecond // their standard deviation
)

// A Meter measures how late the player's presses fall after the fourth
// beat of each bar, as the game hears them: the sound's way out to the
// ear and the key's way in to the game, added up, plus the player's own
// habit of pressing early or late.
type Meter struct {
	m    tempo.Metronome
	offs []time.Duration // the last taps, at most `window`
}

func NewMeter(m tempo.Metronome) *Meter {
	return &Meter{m: m}
}

// Tap takes note of a press at `at`, against the fourth beat nearest to
// it: negative when early, positive when late. A press half a beat or
// more away from any fourth beat is no answer to it, and is not kept.
func (k *Meter) Tap(at time.Time) (off time.Duration, ok bool) {
	per := k.m.PerBar()
	last := per - 1 // the fourth beat, counted from 0
	bar := math.Round((k.m.Beats(at) - float64(last)) / float64(per))
	n := int(bar)*per + last
	off = at.Sub(k.m.At(n))
	if n < 0 || off.Abs() >= k.m.Beat()/2 {
		return off, false
	}
	k.offs = append(k.offs, off)
	if len(k.offs) > window {
		k.offs = k.offs[1:]
	}
	return off, true
}

// Taps returns how many taps the measure holds, up to `window`.
func (k *Meter) Taps() int { return len(k.offs) }

// Mean returns the mean of the taps kept, and whether it is steady:
// `window` of them, spread by less than `steady`.
func (k *Meter) Mean() (mean time.Duration, ok bool) {
	if len(k.offs) == 0 {
		return 0, false
	}
	var sum float64
	for _, o := range k.offs {
		sum += float64(o)
	}
	mu := sum / float64(len(k.offs))
	var sq float64
	for _, o := range k.offs {
		sq += (float64(o) - mu) * (float64(o) - mu)
	}
	sd := math.Sqrt(sq / float64(len(k.offs)))
	return time.Duration(math.Round(mu)), len(k.offs) == window && sd < float64(steady)
}
