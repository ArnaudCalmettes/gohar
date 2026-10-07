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
		r := NewRunner(&stage{}, steps)
		r.Start()
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
			case *Guess:
				play(r, st.Want)
			case *Hang:
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
