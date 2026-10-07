package lessons

import "math/rand/v2"

// The scripts of the lessons written so far, by ID. Each call builds
// fresh steps: Ask and Repeat keep where they are.
var scripts = map[string]func() []Step{
	"1.1": firstSteps,
}

// Script returns the steps of lesson `id`, or nil while it is not
// written.
func Script(id string) []Step {
	if s, ok := scripts[id]; ok {
		return s()
	}
	return nil
}

// Phrases returns the phrases every script and every activity says,
// for the game to check that it can say them all.
func Phrases() []string {
	var out []string
	add := func(steps []Step) {
		for _, st := range steps {
			if p := st.said(); p != "" {
				out = append(out, p)
			}
		}
	}
	for _, s := range scripts {
		add(s())
	}
	for _, a := range activities {
		add(a(rand.New(rand.NewPCG(1, 2)))) // any order says the same phrases
	}
	return out
}

// MIDI numbers of the walker's phrases: around middle C, under the
// right hand.
const (
	c4 = 60 + iota
	_
	d4
	_
	e4
	f4
	_
	g4
)

// firstSteps is lesson 1.1, "Premiers pas" (see
// docs/debutants/chapitre-1.md): the groups of black keys, C and F
// against each other, C D E F G sung, then the motifs from C.
func firstSteps() []Step {
	return []Step{
		Say{Phrase: "l1.1.hello"},
		Show{Phrase: "l1.1.pairs", PCs: Pair},
		&Ask{Phrase: "l1.1.press.pair", Want: Pair, Teasing: true},
		Show{Phrase: "l1.1.trios", PCs: Trio},
		&Ask{Phrase: "l1.1.press.trio", Want: Trio, Teasing: true},

		// C and F together, each by its group, in both notations.
		Show{Phrase: "l1.1.do", PCs: Do.Hint()},
		&Ask{Phrase: "l1.1.find.do", Want: Do},
		Show{Phrase: "l1.1.fa", PCs: Fa.Hint()},
		&Ask{Phrase: "l1.1.find.fa", Want: Fa},
		&Ask{Phrase: "l1.1.find.c", Want: Do},
		&Ask{Phrase: "l1.1.find.fa", Want: Fa},
		&Ask{Phrase: "l1.1.find.do", Want: Do},
		&Ask{Phrase: "l1.1.find.f", Want: Fa},

		// Five notes from C, under one hand, sung.
		Say{Phrase: "l1.1.five"},
		&Repeat{Phrase: "l1.1.sing", Keys: []int{c4, d4, e4, f4, g4}, Want: []Target{Do, Re, Mi, Fa, Sol}},

		// The motifs: C under the left hand, the right one moving up.
		Say{Phrase: "l1.1.motifs"},
		&Repeat{Phrase: "l1.1.motif.re", Keys: []int{c4, d4, c4}, Want: []Target{Do, Re, Do}},
		&Repeat{Phrase: "l1.1.motif.mi", Keys: []int{c4, e4, c4}, Want: []Target{Do, Mi, Do}},
		&Repeat{Phrase: "l1.1.motif.fa", Keys: []int{c4, f4, c4}, Want: []Target{Do, Fa, Do}},
		&Repeat{Phrase: "l1.1.motif.sol", Keys: []int{c4, g4, c4}, Want: []Target{Do, Sol, Do}},
		Say{Phrase: "l1.1.bravo"},
	}
}
