package lessons

import "math/rand/v2"

// The activities of the lessons written so far, by the lesson's ID (see
// "L'activité de la leçon" in docs/debutants.md). A lesson is given
// once; its activity is played as often as the player likes, and draws
// its order anew each time from `rng`.
var activities = map[string]func(rng *rand.Rand) []Step{
	"1.1": firstStepsActivity,
}

// Activity returns the steps of the activity of lesson `id`, drawn
// from `rng`, or nil while it is not written.
func Activity(id string, rng *rand.Rand) []Step {
	if a, ok := activities[id]; ok {
		return a(rng)
	}
	return nil
}

// A motif is played from C and back: shown, the walker plays `keys`;
// said, the bubble gives `letters`, the phrase of its letter names.
type motif struct {
	keys    []int
	want    []Target
	letters string
}

// played and spoken are the two ways of asking for a motif.
func (m motif) played() Step {
	return &Repeat{Phrase: "a1.1.after", Keys: m.keys, Want: m.want}
}

func (m motif) spoken() Step {
	return &Repeat{Phrase: m.letters, Want: m.want}
}

// firstStepsActivity is the activity of lesson 1.1 (see
// docs/debutants/chapitre-1.md): its motifs, shown then said in letters,
// in order then shuffled, then both ways mixed, with the whole run of
// the lesson, C D E F G, among them. The lesson named the notes in both
// notations; from here on, the game asks for letters only, and the
// keyboard keeps both written on its keys.
func firstStepsActivity(rng *rand.Rand) []Step {
	motifs := []motif{
		{[]int{c4, d4, c4}, []Target{Do, Re, Do}, "a1.1.cdc"},
		{[]int{c4, e4, c4}, []Target{Do, Mi, Do}, "a1.1.cec"},
		{[]int{c4, f4, c4}, []Target{Do, Fa, Do}, "a1.1.cfc"},
		{[]int{c4, g4, c4}, []Target{Do, Sol, Do}, "a1.1.cgc"},
	}
	run := motif{[]int{c4, d4, e4, f4, g4}, []Target{Do, Re, Mi, Fa, Sol}, "a1.1.cdefg"}

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

	steps := []Step{Say{Phrase: "a1.1.hello"}}
	steps = append(steps, way(motif.played)...)
	steps = append(steps, Say{Phrase: "a1.1.letters"})
	steps = append(steps, way(motif.spoken)...)
	steps = append(steps, Say{Phrase: "a1.1.sing"}, Say{Phrase: "a1.1.mixed"})

	// Both ways mixed: each motif once each way, and the run once each
	// way, all shuffled. The advice to sing comes back halfway.
	var mixed []Step
	for _, m := range append(motifs, run) {
		mixed = append(mixed, m.played(), m.spoken())
	}
	rng.Shuffle(len(mixed), func(i, j int) { mixed[i], mixed[j] = mixed[j], mixed[i] })
	half := len(mixed) / 2
	steps = append(steps, mixed[:half]...)
	steps = append(steps, Say{Phrase: "a1.1.sing.again"})
	steps = append(steps, mixed[half:]...)
	return append(steps, Say{Phrase: "a1.1.bravo"})
}
