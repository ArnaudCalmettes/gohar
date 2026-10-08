package lessons

import "testing"

// strike presses and lets go of `name` at beat `at` of the pulse, `off`
// beats from it.
func strike(r *Runner, s *stage, name string, at int, off float64) {
	s.at, s.off = at, off
	r.NoteOn(key(name))
	r.NoteOff(key(name))
}

// beats has the pulse go from beat `from` to beat `to`, both included.
func beats(r *Runner, from, to int) {
	for n := from; n <= to; n++ {
		r.Beat(n)
	}
}

// On the 1 only, two bars, the line of the blues in C: nothing counts
// until the player plays a 1, here the second one; then C, a little
// late, and F on the next 1, a little early, and the step is over at
// the next downbeat.
func TestGrooveOnTheOne(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Groove{
		Phrase: "on the 1", Again: "again", On: []int{1}, Bars: 2,
		Want: []Target{Do, Fa}, Chords: []string{"C7", "F7"},
	}})
	r.Start()
	s.take()
	beats(r, 0, 2)
	strike(r, s, "C2", 2, 0) // on a 2: nothing starts
	beats(r, 3, 4)
	strike(r, s, "C2", 4, 0.1)
	beats(r, 5, 8)
	strike(r, s, "F3", 8, -0.2)
	beats(r, 9, 12)
	expect(t, s, "line C7 | F7, bar 0", "beat 4 landed", "line C7 | F7, bar 1", "beat 8 landed")
	if !r.Done() {
		t.Error("not over after two bars landed")
	}
}

// A miss never stops the pulse: a G on a 1 starts nothing, and shows
// why; once started on C, a note on beat 2 and the 1 of the second bar
// left empty are marked as they come, and the series waits for a 1
// again.
func TestGrooveMisses(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Groove{Phrase: "on the 1", Again: "again", On: []int{1}, Bars: 2, Want: []Target{Do, Fa}}})
	r.Start()
	s.take()
	beats(r, 0, 4)
	strike(r, s, "G2", 4, 0)
	beats(r, 5, 8)
	strike(r, s, "C2", 8, 0)
	strike(r, s, "C2", 9, 0)
	beats(r, 9, 16)
	expect(t, s, "beat 4 missed", "beat 8 landed", "beat 9 missed", "beat 12 missed", "says again")
	if r.Done() {
		t.Fatal("over despite the misses")
	}
	// Started again on the 1 of beat 20: landed, it is over.
	beats(r, 17, 20)
	strike(r, s, "C2", 20, 0)
	beats(r, 21, 24)
	strike(r, s, "F2", 24, 0.3)
	beats(r, 25, 28)
	if !r.Done() {
		t.Errorf("not over after the series landed again: %v", s.take())
	}
}

// Too far from the beat, a note misses it, however right.
func TestGrooveOffTheBeat(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Groove{Phrase: "every beat", Again: "again", On: []int{1, 2, 3, 4}, Bars: 1}})
	r.Start()
	s.take()
	beats(r, 0, 4)
	strike(r, s, "C4", 4, 0)
	strike(r, s, "C4", 5, 0.45)
	expect(t, s, "beat 4 landed", "beat 5 missed")
}
