package lessons

import (
	"math/rand/v2"
	"testing"
)

// play answers `want` as a student who gets it right: its keys down
// around middle C, then up.
func play(r *Runner, want Target) {
	var ks []int
	for _, p := range want.Hint() {
		ks = append(ks, 60+p)
	}
	for _, k := range ks {
		r.NoteOn(k)
	}
	for _, k := range ks {
		r.NoteOff(k)
	}
}

// Every script and every activity runs to its end for a student who reads every bubble and
// plays every answer right, and says only phrases it names.
func TestScripts(t *testing.T) {
	all := map[string][]Step{}
	for id := range scripts {
		all[id] = Script(id)
	}
	for id := range activities {
		all["activity "+id] = Activity(id, rand.New(rand.NewPCG(1, 2)))
	}
	for id, steps := range all {
		s := &stage{}
		r := NewRunner(s, steps)
		r.Start()
		beat := 0 // the pulse, through the steps in rhythm
		for i, st := range steps {
			if r.Done() {
				t.Fatalf("%s: done at step %d of %d", id, i, len(steps))
			}
			switch st := st.(type) {
			case Say, Show:
				r.Read()
			case Play:
				r.PhraseEnded()
			case *Ask:
				play(r, st.Want)
			case *Repeat:
				r.PhraseEnded()
				for _, w := range st.Want {
					play(r, w)
				}
			case *PlayLine:
				if st.Keys != nil {
					r.PhraseEnded()
				}
				for _, w := range st.Want {
					play(r, w)
				}
			case Write, Band:
				r.Read()
			case *Guess:
				play(r, st.Want)
			case *Groove:
				beat = groove(r, s, st, beat)
			case *Hang:
				r.PhraseEnded()
				for _, w := range st.Want {
					play(r, w)
				}
				play(r, st.Resolve)
			}
			r.Resume() // the keys up, a beat gone
		}
		if !r.Done() {
			t.Errorf("%s: not done after its last step", id)
		}
	}
	if len(Phrases()) == 0 {
		t.Error("no phrases")
	}
}

// groove plays `st` through as a student who gets every beat right,
// from beat `n` of the pulse on, and returns the beat after it ends.
func groove(r *Runner, s *stage, st *Groove, n int) int {
	on := func() bool { return !r.Done() && r.steps[r.at] == Step(st) }
	for ; on(); n++ {
		r.Beat(n)
		bar, ok := st.wanted[n]
		if st.first < 0 && n%beatsPerBar == 0 {
			bar, ok = 0, true // the 1 that starts the series
		}
		if ok && on() {
			s.at, s.off = n, 0
			if st.Want != nil {
				play(r, st.Want[bar])
				continue
			}
			r.NoteOn(60)
			r.NoteOff(60)
		}
	}
	return n
}
