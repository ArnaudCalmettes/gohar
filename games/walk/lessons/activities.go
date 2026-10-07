package lessons

import (
	"math/rand/v2"
	"slices"
)

// The activities of the lessons written so far, by the lesson's ID (see
// "L'activité de la leçon" in docs/debutants.md). A lesson is given
// once; its activity is played as often as the player likes, and draws
// its order anew each time from `rng`.
var activities = map[string]func(rng *rand.Rand) []Step{
	"1.1": func(rng *rand.Rand) []Step { return motifActivity(rng, motifsUp, runUp, runDown) },
	"1.2": func(rng *rand.Rand) []Step { return motifActivity(rng, slices.Concat(motifsUp, motifsDown), scaleUp) },
	"1.3": func(rng *rand.Rand) []Step { return accidentalActivity(rng, "a1.3.hello", sharpsPlain, sharpsTraps) },
	"1.4": func(rng *rand.Rand) []Step { return accidentalActivity(rng, "a1.4.hello", flatsPlain, flatsTraps) },
}

// Activity returns the steps of the activity of lesson `id`, drawn
// from `rng`, or nil while it is not written.
func Activity(id string, rng *rand.Rand) []Step {
	if a, ok := activities[id]; ok {
		return a(rng)
	}
	return nil
}

// A motif is played from C and back: played, the walker plays `keys`;
// spoken, the bubble gives `letters`, the phrase of its letter names.
type motif struct {
	keys    []int
	want    []Target
	letters string
}

// played and spoken are the two ways of asking for a motif.
func (m motif) played() Step {
	return &Repeat{Phrase: "activity.after", Keys: m.keys, Want: m.want}
}

func (m motif) spoken() Step {
	return &Repeat{Phrase: m.letters, Want: m.want}
}

// The motifs of the lessons and their runs: lesson 1.1's up from middle
// C, lesson 1.2's down from the C above. Said, C G C is the same both
// ways.
var (
	motifsUp = []motif{
		{[]int{c4, d4, c4}, []Target{Do, Re, Do}, "letters.CDC"},
		{[]int{c4, e4, c4}, []Target{Do, Mi, Do}, "letters.CEC"},
		{[]int{c4, f4, c4}, []Target{Do, Fa, Do}, "letters.CFC"},
		{[]int{c4, g4, c4}, []Target{Do, Sol, Do}, "letters.CGC"},
	}
	motifsDown = []motif{
		{[]int{c5, b4, c5}, []Target{Do, Si, Do}, "letters.CBC"},
		{[]int{c5, a4, c5}, []Target{Do, La, Do}, "letters.CAC"},
		{[]int{c5, g4, c5}, []Target{Do, Sol, Do}, "letters.CGC"},
	}
	runUp   = motif{[]int{c4, d4, e4, f4, g4}, []Target{Do, Re, Mi, Fa, Sol}, "letters.CDEFG"}
	runDown = motif{[]int{g4, f4, e4, d4, c4}, []Target{Sol, Fa, Mi, Re, Do}, "letters.GFEDC"}
	scaleUp = motif{
		[]int{c4, d4, e4, f4, g4, a4, b4, c5},
		[]Target{Do, Re, Mi, Fa, Sol, La, Si, Do},
		"letters.CDEFGABC",
	}
)

// motifActivity is the model of the activities of chapter 1 (see
// docs/debutants/chapitre-1.md): `motifs`, played then said in
// letters, in order then shuffled, then both ways mixed, with `runs`,
// longer sequences, among them. The lessons named the notes in both
// notations; from here on, the game asks for letters only, and the
// keyboard keeps both written on its keys.
func motifActivity(rng *rand.Rand, motifs []motif, runs ...motif) []Step {
	// way plays the motifs one way, in order then shuffled.
	way := func(ask func(motif) Step) []Step {
		var steps []Step
		for _, m := range motifs {
			steps = append(steps, ask(m))
		}
		for _, i := range rng.Perm(len(motifs)) {
			steps = append(steps, ask(motifs[i]))
		}
		return steps
	}

	steps := []Step{Say{Phrase: "activity.hello"}}
	steps = append(steps, way(motif.played)...)
	steps = append(steps, Say{Phrase: "activity.letters"})
	steps = append(steps, way(motif.spoken)...)
	steps = append(steps, Say{Phrase: "activity.sing"}, Say{Phrase: "activity.mixed"})

	// Both ways mixed: each motif once each way, and each run once each
	// way, all shuffled. The advice to sing comes back halfway.
	var mixed []Step
	for _, m := range slices.Concat(motifs, runs) {
		mixed = append(mixed, m.played(), m.spoken())
	}
	rng.Shuffle(len(mixed), func(i, j int) { mixed[i], mixed[j] = mixed[j], mixed[i] })
	half := len(mixed) / 2
	steps = append(steps, mixed[:half]...)
	steps = append(steps, Say{Phrase: "activity.sing.again"})
	steps = append(steps, mixed[half:]...)
	return append(steps, Say{Phrase: "activity.bravo"})
}

// said is a sequence only said, in letters: `phrase` gives `want`.
func said(phrase string, want ...Target) motif { return motif{want: want, letters: phrase} }

// The sequences of the activities of lessons 1.3 and 1.4: a white key,
// a black key by its sharp or its flat, then the next white key, up and
// down, every black key that way; and the traps, E♯ and B♯, or F♭ and
// C♭, among them.
var (
	sharpsPlain = slices.Concat(sharpsUp, sharpsDown)
	flatsPlain  = slices.Concat(flatsUp, flatsDown)

	sharpsUp = []motif{
		said("letters.C-Cs-D", Do, DoSharp, Re),
		said("letters.D-Ds-E", Re, ReSharp, Mi),
		said("letters.F-Fs-G", Fa, FaSharp, Sol),
		said("letters.G-Gs-A", Sol, SolSharp, La),
		said("letters.A-As-B", La, LaSharp, Si),
	}
	sharpsDown = []motif{
		said("letters.D-Cs-C", Re, DoSharp, Do),
		said("letters.E-Ds-D", Mi, ReSharp, Re),
		said("letters.G-Fs-F", Sol, FaSharp, Fa),
		said("letters.A-Gs-G", La, SolSharp, Sol),
		said("letters.B-As-A", Si, LaSharp, La),
	}
	sharpsTraps = []motif{
		said("letters.E-Es-Fs", Mi, MiSharp, FaSharp),
		said("letters.B-Bs-Cs", Si, SiSharp, DoSharp),
		said("letters.Cs-Ds-Es", DoSharp, ReSharp, MiSharp),
		said("letters.Gs-As-Bs", SolSharp, LaSharp, SiSharp),
	}
)

var (
	flatsUp = []motif{
		said("letters.C-Db-D", Do, ReFlat, Re),
		said("letters.D-Eb-E", Re, MiFlat, Mi),
		said("letters.F-Gb-G", Fa, SolFlat, Sol),
		said("letters.G-Ab-A", Sol, LaFlat, La),
		said("letters.A-Bb-B", La, SiFlat, Si),
	}
	flatsDown = []motif{
		said("letters.D-Db-C", Re, ReFlat, Do),
		said("letters.E-Eb-D", Mi, MiFlat, Re),
		said("letters.G-Gb-F", Sol, SolFlat, Fa),
		said("letters.A-Ab-G", La, LaFlat, Sol),
		said("letters.B-Bb-A", Si, SiFlat, La),
	}
	flatsTraps = []motif{
		said("letters.F-Fb-Eb", Fa, FaFlat, MiFlat),
		said("letters.C-Cb-Bb", Do, DoFlat, SiFlat),
		said("letters.Gb-Fb-Eb", SolFlat, FaFlat, MiFlat),
		said("letters.Db-Cb-Bb", ReFlat, DoFlat, SiFlat),
	}
)

// accidentalActivity is the activity of lessons 1.3 and 1.4 (see
// docs/debutants/chapitre-1.md): `plain` sequences, said in letters, in
// a new order each time; then again, with the `traps` slipped in.
func accidentalActivity(rng *rand.Rand, hello string, plain, traps []motif) []Step {
	shuffled := func(ms []motif) []Step {
		var steps []Step
		for _, i := range rng.Perm(len(ms)) {
			steps = append(steps, ms[i].spoken())
		}
		return steps
	}
	steps := []Step{Say{Phrase: hello}}
	steps = append(steps, shuffled(plain)...)
	steps = append(steps, Say{Phrase: "activity.traps"})
	steps = append(steps, shuffled(slices.Concat(plain, traps))...)
	return append(steps, Say{Phrase: "activity.bravo"})
}
