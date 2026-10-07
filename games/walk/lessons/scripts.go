package lessons

import (
	"math/rand/v2"
	"slices"
)

// The scripts of the lessons written so far, by ID. Each call builds
// fresh steps: Ask and Repeat keep where they are.
var scripts = map[string]func() []Step{
	"1.1": firstSteps,
	"1.2": cMajorScale,
	"1.3": sharps,
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
			if o, ok := st.(interface{ alsoSays() []string }); ok {
				out = append(out, o.alsoSays()...)
			}
		}
	}
	for _, s := range scripts {
		add(s())
	}
	for _, a := range activities {
		add(a(rand.New(rand.NewPCG(1, 2)))) // any order says the same phrases
	}
	for _, n := range slices.Concat(sevenNotes, fiveSharps) { // drawn by a script, not all said each time
		out = append(out, findLetter[n], findSolfege[n])
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
	_
	a4
	_
	b4
	c5

	cs4 = c4 + 1
	fs4 = f4 + 1
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

		// Five notes from C, under one hand, sung, up and down.
		Say{Phrase: "l1.1.five"},
		&Repeat{Phrase: "l1.1.sing", Keys: []int{c4, d4, e4, f4, g4}, Want: []Target{Do, Re, Mi, Fa, Sol}},
		&Repeat{Phrase: "l1.1.sing.down", Keys: []int{g4, f4, e4, d4, c4}, Want: []Target{Sol, Fa, Mi, Re, Do}},

		// The motifs: C under the left hand, the right one moving up.
		Say{Phrase: "l1.1.motifs"},
		&Repeat{Phrase: "l1.1.motif.re", Keys: []int{c4, d4, c4}, Want: []Target{Do, Re, Do}},
		&Repeat{Phrase: "l1.1.motif.mi", Keys: []int{c4, e4, c4}, Want: []Target{Do, Mi, Do}},
		&Repeat{Phrase: "l1.1.motif.fa", Keys: []int{c4, f4, c4}, Want: []Target{Do, Fa, Do}},
		&Repeat{Phrase: "l1.1.motif.sol", Keys: []int{c4, g4, c4}, Want: []Target{Do, Sol, Do}},
		Say{Phrase: "l1.1.bravo"},
	}
}

// cMajorScale is lesson 1.2, "La gamme de do majeur" (see
// docs/debutants/chapitre-1.md): C D E F G gone over once more, the
// whole scale sung, the letters with a B left hanging, E and B by the
// black keys, the motifs down from C, then the seven notes to find, in
// a new order each time.
func cMajorScale() []Step {
	up := []int{c4, d4, e4, f4, g4, a4, b4, c5}
	scale := []Target{Do, Re, Mi, Fa, Sol, La, Si, Do}
	steps := []Step{
		// What lesson 1.1 left, once more before going on.
		&Repeat{Phrase: "l1.2.review", Want: scale[:5]},
		Say{Phrase: "l1.2.hello"},

		// The whole scale, sung up and down, in no hurry: the fingering
		// comes later, the thumb passing under is beyond a beginner.
		&Repeat{Phrase: "l1.2.scale.up", Keys: up, Want: scale},
		&Repeat{Phrase: "l1.2.scale.down", Keys: reversed(up), Want: reversed(scale)},

		// The letters, and the B left hanging until the C resolves it.
		Say{Phrase: "l1.2.letters"},
		&Hang{
			Phrase: "l1.2.play.letters", Want: scale[:7],
			Humpf: "l1.2.humpf", Resolve: Do, Thanks: "l1.2.thanks",
		},

		// C and F left of the groups, E and B right of them: said once,
		// the keys lit, no drill.
		Show{Phrase: "l1.2.trick", PCs: []int{int(Do), int(Fa), int(Mi), int(Si)}},

		// The motifs, from the C above, the left hand going down.
		Say{Phrase: "l1.2.motifs"},
		&Repeat{Phrase: "l1.2.motif.si", Keys: []int{c5, b4, c5}, Want: []Target{Do, Si, Do}},
		&Repeat{Phrase: "l1.2.motif.la", Keys: []int{c5, a4, c5}, Want: []Target{Do, La, Do}},
		&Repeat{Phrase: "l1.2.motif.sol", Keys: []int{c5, g4, c5}, Want: []Target{Do, Sol, Do}},

		// The seven notes to find, shuffled, each asked in either
		// notation.
		Say{Phrase: "l1.2.find"},
	}
	for _, i := range rand.Perm(len(sevenNotes)) {
		n := sevenNotes[i]
		phrase := findLetter[n]
		if rand.IntN(2) == 0 {
			phrase = findSolfege[n]
		}
		steps = append(steps, &Ask{Phrase: phrase, Want: n})
	}
	return append(steps, Say{Phrase: "l1.2.bravo"})
}

// The seven white notes, and the phrases that ask for each, in letters
// and in solfège.
var (
	sevenNotes = []Note{Do, Re, Mi, Fa, Sol, La, Si}
	fiveSharps = []Note{DoSharp, ReSharp, FaSharp, SolSharp, LaSharp}
	findLetter = map[Note]string{
		Do: "find.C", Re: "find.D", Mi: "find.E", Fa: "find.F", Sol: "find.G", La: "find.A", Si: "find.B",
		DoSharp: "find.Cs", ReSharp: "find.Ds", FaSharp: "find.Fs", SolSharp: "find.Gs", LaSharp: "find.As",
	}
	findSolfege = map[Note]string{
		Do: "find.do", Re: "find.re", Mi: "find.mi", Fa: "find.fa", Sol: "find.sol", La: "find.la", Si: "find.si",
		DoSharp: "find.dos", ReSharp: "find.res", FaSharp: "find.fas", SolSharp: "find.sols", LaSharp: "find.las",
	}
)

// sharps is lesson 1.3, "Les dièses" (see docs/debutants/chapitre-1.md):
// a note that goes up a notch, the five black keys by their sharps, E♯
// and B♯ to work out, then the twelve keys from C to C.
func sharps() []Step {
	steps := []Step{
		Say{Phrase: "l1.3.hello"},
		Say{Phrase: "l1.3.sign"},

		// A note, then its sharp: up a notch.
		&Repeat{Phrase: "l1.3.c.up", Keys: []int{c4, cs4}, Want: []Target{Do, DoSharp}},
		&Repeat{Phrase: "l1.3.f.up", Keys: []int{f4, fs4}, Want: []Target{Fa, FaSharp}},

		Say{Phrase: "l1.3.find"},
	}
	// The five black keys by their sharps, in solfège, then in letters,
	// each time in a new order.
	for _, find := range []map[Note]string{findSolfege, findLetter} {
		for _, i := range rand.Perm(len(fiveSharps)) {
			n := fiveSharps[i]
			steps = append(steps, &Ask{Phrase: find[n], Want: n})
		}
	}
	return append(steps,
		// E♯ and B♯, to work out, three chances each.
		&Guess{Phrase: "l1.3.guess.es", Want: MiSharp, Nope: "l1.3.nope", Answer: "l1.3.answer.es", Chances: 3},
		&Guess{Phrase: "l1.3.guess.bs", Want: SiSharp, Nope: "l1.3.nope", Answer: "l1.3.answer.bs", Chances: 3},
		Say{Phrase: "l1.3.enharmonic"},

		// The twelve keys, up from C, key by key.
		&Repeat{Phrase: "l1.3.chromatic", Want: []Target{
			Do, DoSharp, Re, ReSharp, Mi, Fa, FaSharp, Sol, SolSharp, La, LaSharp, Si, Do,
		}},
		Say{Phrase: "l1.3.bravo"},
	)
}

// reversed returns a reversed copy of `s`.
func reversed[T any](s []T) []T {
	out := make([]T, len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}
