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
// all measure their delays and keep their dates the same way.
type Instrument interface {
	io.Reader

	// NoteOn starts a key now, or restarts it if it was already
	// sounding. `at` is when the press happened, recorded for Delays.
	NoteOn(key int, velocity float64, at time.Time)

	// NoteOff releases a key now. A key not sounding is not an error.
	NoteOff(key int, at time.Time)

	// ScheduleOn starts a key at instant `at`, to the sample. Meant for
	// what the program plays, a metronome or a bass line, queued ahead
	// of time: a date already past plays at once, late.
	ScheduleOn(key int, velocity float64, at time.Time)

	// ScheduleOff releases a key at instant `at`, to the sample.
	ScheduleOff(key int, at time.Time)

	// Delays returns the last and the worst delay between a key event
	// and the moment it reached the samples. Scheduled notes do not
	// count: they are early on purpose.
	Delays() (last, worst time.Duration)

	// Histogram returns how those delays were spread.
	Histogram() Histogram
}

type command struct {
	on       bool
	key      int
	velocity float64

	// at is when the key was pressed, or with `dated`, when to play it.
	at    time.Time
	dated bool
}

// maxDated is how many scheduled commands may wait for their frame. A
// game schedules a window ahead, a tenth of a second or so: even a
// walking bass under a comping piano and a metronome stays far below.
// Past it, the latest dates are dropped.
const maxDated = 256

// A queue is the handover between the goroutine that presses keys and
// the audio goroutine, the measure of the delay across it, and the
// keeper of the dated commands until their frame comes.
//
// The mutex is held only to push a command or to copy the pending ones
// out, never while samples are computed, so the audio goroutine waits
// at most for a copy of sixty four small structs.
//
// Embedded in an instrument, it provides NoteOn, NoteOff, ScheduleOn,
// ScheduleOff, Delays and Histogram; the instrument implements renderer
// and its Read calls render.
type queue struct {
	mu        sync.Mutex
	commands  [maxCommands]command
	nCommands int
	last      time.Duration
	worst     time.Duration
	histogram Histogram

	// Audio side, touched by render alone.
	taken  [maxCommands]command
	dated  [maxDated]command // sorted by date, earliest first
	nDated int
	clock  clock
}

// A renderer is the instrument as render sees it: what a command does
// to its voices, and how to fill frames with what is sounding.
type renderer interface {
	apply(c command)

	// fill writes whole frames into `buf`, moving the voices on.
	fill(buf []byte)
}

func (q *queue) NoteOn(key int, velocity float64, at time.Time) {
	q.push(command{on: true, key: key, velocity: velocity, at: at})
}

func (q *queue) NoteOff(key int, at time.Time) {
	q.push(command{key: key, at: at})
}

func (q *queue) ScheduleOn(key int, velocity float64, at time.Time) {
	q.push(command{on: true, key: key, velocity: velocity, at: at, dated: true})
}

func (q *queue) ScheduleOff(key int, at time.Time) {
	q.push(command{key: key, at: at, dated: true})
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
// key event waited until `now`, and empties the queue.
func (q *queue) take(now time.Time, dst *[maxCommands]command) int {
	q.mu.Lock()
	n := q.nCommands
	for i := range n {
		c := q.commands[i]
		if !c.dated {
			d := now.Sub(c.at)
			q.last = d
			if d > q.worst {
				q.worst = d
			}
			q.histogram.add(d)
		}
		dst[i] = c
	}
	q.nCommands = 0
	q.mu.Unlock()
	return n
}

// render is the body of an instrument's Read: it applies the key events
// at the top of the buffer, and each dated command on its own frame,
// filling the frames in between.
//
// A dated command already past is applied on the first frame. Two
// commands dated alike keep the order they were pushed in, so that a
// release and a press of the same key on the same beat do not swap.
func (q *queue) render(r renderer, buf []byte, now time.Time) int {
	q.clock.tick(now)
	n := q.take(now, &q.taken)
	for i := range n {
		if q.taken[i].dated {
			q.insert(q.taken[i])
		} else {
			r.apply(q.taken[i])
		}
	}

	count := int64(len(buf) / BytesPerFrame)
	start := q.clock.frames
	var done int64
	for q.nDated > 0 {
		f := q.clock.frameOf(q.dated[0].at) - start
		if f >= count {
			break
		}
		if f > done {
			r.fill(buf[done*BytesPerFrame : f*BytesPerFrame])
			done = f
		}
		r.apply(q.dated[0])
		q.nDated--
		copy(q.dated[:q.nDated], q.dated[1:q.nDated+1])
	}
	if done < count {
		r.fill(buf[done*BytesPerFrame : count*BytesPerFrame])
	}
	q.clock.frames += count
	return int(count) * BytesPerFrame
}

// insert places a dated command after every command dated earlier or
// alike. Linear, from the end: dates mostly arrive in order, so the
// loop rarely runs.
func (q *queue) insert(c command) {
	i := q.nDated
	for i > 0 && q.dated[i-1].at.After(c.at) {
		i--
	}
	if i == len(q.dated) {
		return // later than everything, and no room left
	}
	end := min(q.nDated, len(q.dated)-1)
	copy(q.dated[i+1:end+1], q.dated[i:end])
	q.dated[i] = c
	q.nDated = end + 1
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
