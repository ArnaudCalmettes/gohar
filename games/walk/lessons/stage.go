package lessons

// A Stage is what the game lends a lesson: the walker's bubble, the
// keyboard on the screen, the walker's own playing, and the casserole.
// The lesson never draws nor sounds anything itself, so that it runs in
// any game, and in a test without a screen (see "Le code" in
// docs/debutants.md).
type Stage interface {
	// Say shows the bubble of the game's phrase `phrase`, in the
	// player's language: the lesson knows the phrase's ID only.
	Say(phrase string)

	// Light lights every key of the pitch classes `pcs` (0 for C, up to
	// 11 for B), in every octave. Nil puts them out.
	Light(pcs []int)

	// Play has the walker play `keys`, MIDI numbers, one after the
	// other. The game calls Runner.PhraseEnded once the last one is over.
	Play(keys []int)

	// Miss answers a wrong note: the casserole sounds, and the right
	// keys, every key of the pitch classes `pcs`, are shown.
	Miss(pcs []int)

	// Hit answers what was asked, played right to its end: a sign that
	// it is right, before the next step.
	Hit()

	// Tease answers misses in a row where nobody misses in good faith:
	// the walker's annoyance, as far as `level`. At FedUp he walks out,
	// and the game resumes the step with Runner.Again once he is back.
	Tease(level Annoyance)
}

// An Annoyance is how far the walker's patience has gone.
type Annoyance int

const (
	Hey       Annoyance = iota + 1 // "Hé…"
	OnPurpose                      // "Tu le fais exprès ?"
	FedUp                          // he walks out: GAME OVER
)
