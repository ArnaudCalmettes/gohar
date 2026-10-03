package synth

import (
	"testing"
	"time"
)

var t0 = time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)

// firstSound returns the index of the first frame that is not silent,
// or -1.
func firstSound(e *Engine, buf []byte, now time.Time) int {
	e.render(e, buf, now)
	for f := range len(buf) / BytesPerFrame {
		for _, b := range buf[f*BytesPerFrame : f*BytesPerFrame+4] {
			if b != 0 {
				return f
			}
		}
	}
	return -1
}

// A square starts at full level on its first sample, so the first
// frame that sounds is the frame the note was dated on.
func square() *Engine { return &Engine{Timbre: Square} }

// 10 ms after the first Read is frame 480 at 48 kHz.
func TestScheduledNoteLandsOnItsFrame(t *testing.T) {
	e := square()
	buf := make([]byte, 1024*BytesPerFrame)
	e.ScheduleOn(60, 1, t0.Add(10*time.Millisecond))
	if f := firstSound(e, buf, t0); f != 480 {
		t.Errorf("a note dated 10 ms ahead sounds on frame %d, want 480", f)
	}
}

// A date beyond the buffer waits for the next one, and lands on its
// frame there, the clock having moved on by exactly one buffer.
func TestScheduledNoteWaitsForItsBuffer(t *testing.T) {
	e := square()
	buf := make([]byte, 480*BytesPerFrame) // 10 ms
	e.ScheduleOn(60, 1, t0.Add(15*time.Millisecond))
	if f := firstSound(e, buf, t0); f != -1 {
		t.Fatalf("the first 10 ms sound from frame %d, want silence", f)
	}
	if f := firstSound(e, buf, t0.Add(10*time.Millisecond)); f != 240 {
		t.Errorf("the note sounds on frame %d of the second buffer, want 240", f)
	}
}

// Jitter in the driver's calls does not move a date: the clock only
// follows a sixty fourth of each gap.
func TestJitterDoesNotMoveADate(t *testing.T) {
	e := square()
	buf := make([]byte, 480*BytesPerFrame)
	e.ScheduleOn(60, 1, t0.Add(25*time.Millisecond))
	firstSound(e, buf, t0)
	firstSound(e, buf, t0.Add(13*time.Millisecond)) // 3 ms late
	if f := firstSound(e, buf, t0.Add(20*time.Millisecond)); f < 238 || f > 242 {
		t.Errorf("the note sounds on frame %d of the third buffer, want 240 give or take 2", f)
	}
}

// A date already past plays at once, on the first frame.
func TestLateNotePlaysAtOnce(t *testing.T) {
	e := square()
	buf := make([]byte, 256*BytesPerFrame)
	firstSound(e, buf, t0)
	e.ScheduleOn(60, 1, t0)
	if f := firstSound(e, buf, t0.Add(5*time.Millisecond)); f != 0 {
		t.Errorf("a late note sounds on frame %d, want 0", f)
	}
}

// A release and a press of the same key on the same beat keep their
// order: the key sounds after the beat.
func TestSameDateKeepsOrder(t *testing.T) {
	e := square()
	buf := make([]byte, 1024*BytesPerFrame)
	beat := t0.Add(10 * time.Millisecond)
	e.ScheduleOn(60, 1, t0.Add(time.Millisecond))
	e.ScheduleOff(60, beat)
	e.ScheduleOn(60, 1, beat)
	e.render(e, buf, t0)
	if e.voices[0].state == voiceRelease || e.voices[0].state == voiceOff {
		t.Errorf("the key was released after being struck again: state %d", e.voices[0].state)
	}
}

// Dated commands are kept in order whatever order they arrive in, and
// the latest ones are dropped when there is no room left.
func TestInsertSortsAndDrops(t *testing.T) {
	var q queue
	for _, ms := range []int{30, 10, 20, 10} {
		q.insert(command{key: ms, at: t0.Add(time.Duration(ms) * time.Millisecond), dated: true})
	}
	var got []int
	for i := range q.nDated {
		got = append(got, q.dated[i].key)
	}
	if len(got) != 4 || got[0] != 10 || got[1] != 10 || got[2] != 20 || got[3] != 30 {
		t.Errorf("dates in order %v, want 10 10 20 30", got)
	}

	q = queue{}
	for i := range maxDated {
		q.insert(command{at: t0.Add(time.Duration(i) * time.Millisecond), dated: true})
	}
	q.insert(command{key: 1, at: t0.Add(-time.Second), dated: true})
	if q.nDated != maxDated || q.dated[0].key != 1 {
		t.Errorf("a full queue keeps %d dates, the earliest first: key %d", q.nDated, q.dated[0].key)
	}
}

// Scheduled notes are early on purpose: they stay out of Delays.
func TestScheduledNotesAreNotDelays(t *testing.T) {
	e := square()
	buf := make([]byte, 256*BytesPerFrame)
	e.ScheduleOn(60, 1, t0.Add(-time.Second))
	e.render(e, buf, t0)
	if h := e.Histogram(); h.Total() != 0 {
		t.Errorf("a scheduled note counted as a delay: %v", h.Counts)
	}
}

func TestScheduleDoesNotAllocate(t *testing.T) {
	e := square()
	buf := make([]byte, 256*BytesPerFrame)
	allocs := testing.AllocsPerRun(50, func() {
		now := time.Now()
		e.ScheduleOn(60, 1, now.Add(2*time.Millisecond))
		e.ScheduleOff(60, now.Add(4*time.Millisecond))
		_, _ = e.Read(buf)
	})
	if allocs != 0 {
		t.Errorf("%v allocations per run", allocs)
	}
}
