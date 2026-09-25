package keyboard_test

import (
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
)

// Plain testing on purpose: the games module does not depend on
// testify, and one file of tests is not worth adding it.

type recorder struct {
	mu     sync.Mutex
	events []keyboard.Event
}

func (r *recorder) recv(e keyboard.Event) {
	r.mu.Lock()
	r.events = append(r.events, e)
	r.mu.Unlock()
}

func (r *recorder) get() []keyboard.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]keyboard.Event(nil), r.events...)
}

const ms = time.Millisecond

// A pedal under two notes: the shape of a question.
func pedalAndTwoNotes() []keyboard.Note {
	return []keyboard.Note{
		{Key: 48, Velocity: 0.5, Start: 0, Length: 60 * ms},
		{Key: 62, Velocity: 0.8, Start: 20 * ms, Length: 20 * ms},
		{Key: 64, Velocity: 0.8, Start: 40 * ms, Length: 20 * ms},
	}
}

func TestSequencePlaysInOrderAndOnTime(t *testing.T) {
	s := keyboard.NewSequence("question", pedalAndTwoNotes())
	var r recorder
	if err := s.Listen(r.recv); err != nil {
		t.Fatal(err)
	}

	select {
	case <-s.Done():
	case <-time.After(time.Second):
		t.Fatal("the sequence never ended")
	}

	got := r.get()
	type key struct {
		key  int
		down bool
	}
	var order []key
	for _, e := range got {
		order = append(order, key{e.Key, e.Down})
	}
	// At 40 ms the 62 lets go before the 64 starts: ups come first.
	want := []key{
		{48, true}, {62, true}, {62, false}, {64, true},
		{48, false}, {64, false},
	}
	if !slices.Equal(order, want) {
		t.Fatalf("order %v, want %v", order, want)
	}

	t.Run("never early, and lateness does not add up", func(t *testing.T) {
		start := got[0].At
		wants := []time.Duration{0, 20 * ms, 40 * ms, 40 * ms, 60 * ms, 60 * ms}
		for i, e := range got {
			late := e.At.Sub(start) - wants[i]
			if late < -ms {
				t.Errorf("event %d early by %v", i, -late)
			}
			// Generous on purpose: this checks drift, not the scheduler
			// of a loaded CI machine.
			if late > 15*ms {
				t.Errorf("event %d late by %v", i, late)
			}
		}
	})

	t.Run("one playback per sequence", func(t *testing.T) {
		if s.Listen(r.recv) == nil {
			t.Error("a second Listen was accepted")
		}
	})
}

func TestSequenceCloseReleasesWhatIsHeld(t *testing.T) {
	s := keyboard.NewSequence("skipped", pedalAndTwoNotes())
	var r recorder
	if err := s.Listen(r.recv); err != nil {
		t.Fatal(err)
	}

	time.Sleep(30 * ms)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	after := len(r.get())
	time.Sleep(50 * ms)

	got := r.get()
	if len(got) != after {
		t.Errorf("%d events after Close returned", len(got)-after)
	}

	down := map[int]bool{}
	for _, e := range got {
		down[e.Key] = e.Down
	}
	for k, d := range down {
		if d {
			t.Errorf("key %d left sounding", k)
		}
	}
}

func TestSequenceDropsEmptyNotes(t *testing.T) {
	s := keyboard.NewSequence("empty", []keyboard.Note{{Key: 60, Velocity: 1}})
	var r recorder
	if err := s.Listen(r.recv); err != nil {
		t.Fatal(err)
	}
	<-s.Done()
	if got := r.get(); len(got) != 0 {
		t.Errorf("an empty note sent %v", got)
	}
}

func TestCloseWithoutListening(t *testing.T) {
	if err := keyboard.NewSequence("idle", nil).Close(); err != nil {
		t.Error(err)
	}
}
