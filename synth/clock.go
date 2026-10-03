package synth

import (
	"math"
	"time"
)

const (
	// smoothing is how slowly the clock follows the driver: each Read
	// moves the origin by 1/smoothing of the gap it sees. The driver
	// calls in bursts, a few buffers at once to refill its ring, so a
	// single reading says little; sixty four of them average the bursts
	// out and still follow the slow drift between the sound card's
	// crystal and the system's.
	smoothing = 64

	// resync is the gap past which the clock stops averaging and jumps:
	// an underrun, a suspended laptop. Smoothing out half a second would
	// play every dated note off for minutes.
	resync = 100 * time.Millisecond
)

// A clock ties the frames an instrument writes to the wall clock, so
// that a note dated at an instant lands on a frame.
//
// The tie is an estimate. What the engine sees is when Read is called,
// not when its samples reach the speaker, which happens later by the
// queues below it (see Device.Floor). That later is constant for a
// given setting, so it shifts every dated note alike: the notes keep
// their spacing to the sample, and the shift is the output latency the
// calibration measures.
//
// Touched by the audio goroutine alone.
type clock struct {
	origin  time.Time // the estimated instant of frame 0
	frames  int64     // frames written so far
	started bool
}

// tick is called at the top of each Read, before any frame is written:
// `now` is the estimated instant of frame `frames`.
func (c *clock) tick(now time.Time) {
	if !c.started {
		c.origin, c.started = now, true
		return
	}
	gap := now.Sub(c.at(c.frames))
	if gap > resync || gap < -resync {
		c.origin = c.origin.Add(gap)
		return
	}
	c.origin = c.origin.Add(gap / smoothing)
}

// at returns the instant of frame `f`.
func (c *clock) at(f int64) time.Time {
	return c.origin.Add(time.Duration(f * int64(time.Second) / SampleRate))
}

// frameOf returns the frame closest to `t`.
func (c *clock) frameOf(t time.Time) int64 {
	return int64(math.Round(t.Sub(c.origin).Seconds() * SampleRate))
}
