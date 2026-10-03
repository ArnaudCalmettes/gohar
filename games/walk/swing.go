package main

import (
	"math"
	"time"
)

// A Pulse maps positions, counted in beats, to instants and back. The
// Metronome is the straight one; Swing bends another.
type Pulse interface {
	AtBeats(x float64) time.Time
	Beats(t time.Time) float64
}

// Swing delays the "and" of every beat, the second eighth, by moving it
// to `ratio` of the beat. 0.5 keeps the eighths even, 2/3 is the
// triplet swing of the method books, 0.75 the tight dotted swing of
// early New Orleans jazz (see "Le son" in docs/walk.md).
//
// Swing belongs to a style, not to the pulse, hence a decorator: the
// beats themselves never move, only what falls between them. Within a
// beat, the first half stretches over [0, ratio] and the second
// squeezes into [ratio, 1]; both ways are linear, so the map inverts
// exactly, and the marker can tell a swung "and" from a late one.
type Swing struct {
	Pulse
	ratio float64
}

var _ Pulse = Swing{}

// NewSwing swings `p` at `ratio`, strictly between 0 and 1.
func NewSwing(p Pulse, ratio float64) Swing {
	if !(ratio > 0 && ratio < 1) {
		panic("walk: a swing ratio lies strictly between 0 and 1")
	}
	return Swing{Pulse: p, ratio: ratio}
}

// AtBeats returns the instant of straight position `x`, swung.
func (s Swing) AtBeats(x float64) time.Time {
	return s.Pulse.AtBeats(s.bend(x))
}

// Beats returns the straight position of `t`: a swung "and" reads as
// 0.5 past its beat.
func (s Swing) Beats(t time.Time) float64 {
	return s.unbend(s.Pulse.Beats(t))
}

func (s Swing) bend(x float64) float64 {
	n := math.Floor(x)
	f := x - n
	if f < 0.5 {
		return n + f*2*s.ratio
	}
	return n + s.ratio + (f-0.5)*2*(1-s.ratio)
}

func (s Swing) unbend(y float64) float64 {
	n := math.Floor(y)
	f := y - n
	if f < s.ratio {
		return n + f/(2*s.ratio)
	}
	return n + 0.5 + (f-s.ratio)/(2*(1-s.ratio))
}

// DueEvery returns the positions, in steps of 1/`sub` beat, whose
// instants fall in [`from`, `to`): step k is position k/`sub`. With
// `sub` at 2, the eighths, swung or not. Beats only ever grows with
// time, so a step lies in the window exactly when its position lies
// between the window's ends, and windows laid end to end hand out each
// step once.
func DueEvery(p Pulse, sub int, from, to time.Time) (first, end int) {
	return int(math.Ceil(p.Beats(from) * float64(sub))), int(math.Ceil(p.Beats(to) * float64(sub)))
}
