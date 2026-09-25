package keyboard

import (
	"cmp"
	"errors"
	"slices"
	"time"
)

// A Note is one key held for a while, placed on a sequence's timeline.
//
// Notes rather than raw events because that is how a game writes what
// it wants heard: a pedal held under a scale is one long note and seven
// short ones, and nobody should have to pair the key ups by hand.
type Note struct {
	Key      int
	Velocity float64
	Start    time.Duration
	Length   time.Duration
}

// A Sequence replays notes as if a finger played them.
//
// # Why a source and not a scheduler in the engine
//
// The engine applies what it receives at the top of its next buffer,
// and NoteOn's time is a measurement, not a date to wait for. Driving a
// melody from the game loop would add a frame of jitter to every note,
// sixteen milliseconds, which a trained ear hears in a scale.
//
// A sequence runs on its own goroutine and wakes on timers, so it
// reaches the engine exactly the way a MIDI keyboard does, with the
// same granularity: the device buffer. What sounds in tune and on time
// under a finger sounds the same here, and the game reads both through
// the same interface.
//
// It is the replayed sequence this package was written for, not a test
// double.
//
// One playback per sequence. Playing it again is a new sequence, which
// costs one slice.
type Sequence struct {
	name   string
	events []scheduled

	stop chan struct{}
	done chan struct{}
}

type scheduled struct {
	after time.Duration
	event Event
}

// NewSequence lays notes on a timeline. The notes may come in any
// order; a note of zero length is dropped, since it would be a key up
// racing its own key down.
func NewSequence(name string, notes []Note) *Sequence {
	s := &Sequence{name: name, done: make(chan struct{})}
	for _, n := range notes {
		if n.Length <= 0 {
			continue
		}
		s.events = append(s.events,
			scheduled{n.Start, Event{Key: n.Key, Velocity: n.Velocity, Down: true}},
			scheduled{n.Start + n.Length, Event{Key: n.Key}},
		)
	}
	// Stable, and ups before downs at the same instant: a note ending
	// where the same key starts again must let go first, or the release
	// would cut the new one.
	slices.SortStableFunc(s.events, func(a, b scheduled) int {
		if a.after != b.after {
			return cmp.Compare(a.after, b.after)
		}
		switch {
		case !a.event.Down && b.event.Down:
			return -1
		case a.event.Down && !b.event.Down:
			return 1
		}
		return 0
	})
	return s
}

// Name returns what the sequence was called.
func (s *Sequence) Name() string { return s.name }

// Listen starts the playback. Events are stamped when they are sent,
// like every other source does on arrival.
func (s *Sequence) Listen(recv func(Event)) error {
	if s.stop != nil {
		return errors.New("keyboard: sequence already playing")
	}
	s.stop = make(chan struct{})
	go s.play(recv)
	return nil
}

// Done is closed once the last event has been sent, or once Close has
// released what was held.
func (s *Sequence) Done() <-chan struct{} { return s.done }

// Close stops the playback and returns once it has.
//
// Keys still down are released on the way out: a question skipped in
// the middle of its scale must not leave a pedal sounding forever. No
// event is sent after Close returns.
func (s *Sequence) Close() error {
	if s.stop == nil {
		return nil
	}
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
	<-s.done
	return nil
}

func (s *Sequence) play(recv func(Event)) {
	defer close(s.done)

	start := time.Now()
	held := map[int]bool{}

	timer := time.NewTimer(0)
	defer timer.Stop()
	<-timer.C

	for _, e := range s.events {
		// Every wait is measured from the start, never from the
		// previous event, so that timer lateness does not add up along
		// the sequence.
		if wait := time.Until(start.Add(e.after)); wait > 0 {
			timer.Reset(wait)
			select {
			case <-timer.C:
			case <-s.stop:
				s.release(recv, held)
				return
			}
		}

		ev := e.event
		ev.At = time.Now()
		recv(ev)
		if ev.Down {
			held[ev.Key] = true
		} else {
			delete(held, ev.Key)
		}
	}
}

func (s *Sequence) release(recv func(Event), held map[int]bool) {
	now := time.Now()
	for key := range held {
		recv(Event{Key: key, At: now})
	}
}
