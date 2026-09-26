package synth

import (
	"io"
	"sync"
	"time"
)

// An Instrument turns key presses into samples: what a game plays
// through, whatever makes the sound.
//
// Read belongs to the audio driver's goroutine and the rest to whoever
// holds the keys, a MIDI callback or a game loop. Every instrument here
// hands over between the two through a queue, the same one, so they
// all measure their delays the same way.
type Instrument interface {
	io.Reader

	// NoteOn starts a key, or restarts it if it was already sounding.
	// `at` is when the press happened, recorded for Delays.
	NoteOn(key int, velocity float64, at time.Time)

	// NoteOff releases a key. A key not sounding is not an error.
	NoteOff(key int, at time.Time)

	// Delays returns the last and the worst delay between a key event
	// and the moment it reached the samples.
	Delays() (last, worst time.Duration)

	// Histogram returns how those delays were spread.
	Histogram() Histogram
}

type command struct {
	on       bool
	key      int
	velocity float64
	at       time.Time
}

// A queue is the handover between the goroutine that presses keys and
// the audio goroutine, and the measure of the delay across it.
//
// The mutex is held only to push a command or to copy the pending ones
// out, never while samples are computed, so the audio goroutine waits
// at most for a copy of sixty four small structs.
//
// Embedded in an instrument, it provides NoteOn, NoteOff, Delays and
// Histogram; the instrument writes Read, taking the commands with take
// at the top of it.
type queue struct {
	mu        sync.Mutex
	commands  [maxCommands]command
	nCommands int
	last      time.Duration
	worst     time.Duration
	histogram Histogram
}

func (q *queue) NoteOn(key int, velocity float64, at time.Time) {
	q.push(command{on: true, key: key, velocity: velocity, at: at})
}

func (q *queue) NoteOff(key int, at time.Time) {
	q.push(command{key: key, at: at})
}

func (q *queue) push(c command) {
	q.mu.Lock()
	if q.nCommands < len(q.commands) {
		q.commands[q.nCommands] = c
		q.nCommands++
	}
	q.mu.Unlock()
}

// take copies the pending commands into `dst`, records how long each
// waited until `now`, and empties the queue. Called by Read alone.
func (q *queue) take(now time.Time, dst *[maxCommands]command) int {
	q.mu.Lock()
	n := q.nCommands
	for i := range n {
		d := now.Sub(q.commands[i].at)
		q.last = d
		if d > q.worst {
			q.worst = d
		}
		q.histogram.add(d)
		dst[i] = q.commands[i]
	}
	q.nCommands = 0
	q.mu.Unlock()
	return n
}

// Delays returns the last and the worst delay between a key event and
// the moment it reached the samples.
//
// This is the part of the latency that belongs to us. What the driver,
// the bus and the card add afterwards is not visible from here, and the
// only honest measurement of the whole remains the ear.
func (q *queue) Delays() (last, worst time.Duration) {
	q.mu.Lock()
	last, worst = q.last, q.worst
	q.mu.Unlock()
	return
}

// Histogram returns how the delays between a key event and the samples
// were spread since the instrument started.
//
// A copy, taken under a lock the audio goroutine holds a few
// microseconds per buffer: cheap enough to read every frame.
func (q *queue) Histogram() Histogram {
	q.mu.Lock()
	h := q.histogram
	q.mu.Unlock()
	return h
}
