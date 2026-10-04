// Package tempo is the map of musical time the games share: which
// instant each beat falls on, and which beat an instant belongs to,
// straight or swung. Walk with me schedules its band on it and marks
// the player against it; the calibration of the latency counts on it.
package tempo

import (
	"math"
	"time"
)

// A Metronome is the map of musical time: which instant each beat
// falls on, and which beat an instant belongs to. The scrolling, the
// marker and the sound all read the same map, so they cannot drift
// apart.
//
// Beats are counted from the first beat of bar 1, which is beat 0. The
// count-in has negative beats and sits in bar 0, the way a musician
// counts "one, two, three, four" before the tune.
//
// Every instant is computed from the start, never by adding beat after
// beat: a thousand beats later, there is no accumulated rounding.
type Metronome struct {
	start  time.Time
	beat   time.Duration
	perBar int
}

// NewMetronome returns the map of a tune whose bar 1 starts at
// `start`, at `bpm` beats per minute and `perBar` beats per bar.
func NewMetronome(start time.Time, bpm float64, perBar int) Metronome {
	return Metronome{
		start:  start,
		beat:   time.Duration(float64(time.Minute) / bpm),
		perBar: perBar,
	}
}

// Beat returns the length of one beat.
func (m Metronome) Beat() time.Duration { return m.beat }

// PerBar returns the beats in a bar.
func (m Metronome) PerBar() int { return m.perBar }

// At returns the instant beat `n` falls on.
func (m Metronome) At(n int) time.Time {
	return m.start.Add(time.Duration(n) * m.beat)
}

// Beats returns where `t` falls, in beats: 2.5 is halfway between the
// third and the fourth beat of bar 1. What the scrolling reads.
func (m Metronome) Beats(t time.Time) float64 {
	return float64(t.Sub(m.start)) / float64(m.beat)
}

// AtBeats returns the instant of position `x`, in beats: the inverse of
// Beats. 2.5 is the "and" of beat 3 in bar 1, straight; Swing moves it.
func (m Metronome) AtBeats(x float64) time.Time {
	return m.start.Add(time.Duration(math.Round(x * float64(m.beat))))
}

// Nearest returns the beat closest to `t`, and how far `t` is from it:
// negative when early, positive when late. What the marker reads, with
// `t` already corrected for the latency.
func (m Metronome) Nearest(t time.Time) (n int, off time.Duration) {
	d := t.Sub(m.start)
	n = floorDiv(d+m.beat/2, m.beat)
	return n, d - time.Duration(n)*m.beat
}

// Due returns the beats whose instants fall in [`from`, `to`), as the
// half-open range [first, end). The sound schedules them ahead, a
// window at a time: a beat at the very edge belongs to one window
// only, so nothing is played twice or skipped.
func (m Metronome) Due(from, to time.Time) (first, end int) {
	return ceilDiv(from.Sub(m.start), m.beat), ceilDiv(to.Sub(m.start), m.beat)
}

// A Position names a beat the way a musician does, bar and beat both
// counted from 1. The count-in is bar 0.
type Position struct {
	Bar, Beat int
}

// Position names beat `n`.
func (m Metronome) Position(n int) Position {
	bar := floorDiv(time.Duration(n), time.Duration(m.perBar))
	return Position{Bar: bar + 1, Beat: n - bar*m.perBar + 1}
}

// Strong tells whether a beat is strong: 1 and 3 in four, the first
// one in any other meter. Siron and Siskind both place the roots there
// (see "Ce qu'on attend, temps par temps" in docs/walk.md).
//
// A stylistic frame decides this, not the pulse: it lives here only
// until a game has a second style to tell apart.
func (m Metronome) Strong(p Position) bool {
	if m.perBar == 4 {
		return p.Beat == 1 || p.Beat == 3
	}
	return p.Beat == 1
}

// floorDiv divides rounding toward minus infinity, so that the
// count-in, before the start, rounds the same way as the tune.
func floorDiv(a, b time.Duration) int {
	q := a / b
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return int(q)
}

func ceilDiv(a, b time.Duration) int {
	return -floorDiv(-a, b)
}
