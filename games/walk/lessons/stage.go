package lessons

// RepeatBar is the bar of a grid that repeats the one before, as a chart
// writes it: "%".
const RepeatBar = "%"

// A Stage is what the game lends a lesson: the walker's bubble, the
// keyboard on the screen, the walker's own playing, and the casserole.
// The lesson never draws nor sounds anything itself, so that it runs in
// any game, and in a test without a screen (see "Le code" in
// docs/debutants.md).
type Stage interface {
	// Say shows the bubble of the game's phrase `phrase`, in the
	// player's language: the lesson knows the phrase's ID only.
	Say(phrase string)

	// Add adds the phrase `phrase` under the bubble's: a reminder, the
	// step's own phrase still there.
	Add(phrase string)

	// Light lights every key of the pitch classes `pcs` (0 for C, up to
	// 11 for B), in every octave. Nil puts them out.
	Light(pcs []int)

	// Play has the walker play `keys`, MIDI numbers, one after the
	// other. The game calls Runner.PhraseEnded once the last one is over.
	Play(keys []int)

	// PlayBars is Play, a note a bar, held through it: the roots of a
	// line of a grid.
	PlayBars(keys []int)

	// Miss answers a wrong note: the casserole sounds, and the right
	// keys, every key of the pitch classes `pcs`, are shown. `beat` is
	// where the wrong note fell in the phrase asked, from 0: the game
	// takes the player as playing in rhythm, a beat a note, four to the
	// bar, and shows the way again on the next bar, as Play does after
	// a Miss.
	Miss(pcs []int, beat int)

	// Hit answers what was asked, played right to its end: a sign that
	// it is right, before the next step.
	Hit()

	// Hold keeps MIDI key `key` sounding after the player lets it go,
	// until Release: a note left hanging. The keyboard on the screen
	// does not show it held.
	Hold(key int)

	// Release lets go of the key Hold kept sounding.
	Release()

	// Blink has MIDI key `key` blink on the keyboard on the screen, a
	// key asked for without a word; 0 stops it.
	Blink(key int)

	// Chords writes chord symbols, in ChordPro ("Eb7"), side by side,
	// in the hand of the charts. Nil wipes them.
	Chords(chords []string)

	// Line shows a line of a grid, a chord a bar, the bar `bar` shaded
	// as the bar played, -1 for none; RepeatBar for a bar repeats the one
	// before. Nil wipes it.
	Line(chords []string, bar int)

	// Bass has the player play the bass: under the split, his keys and
	// the walker's phrases sound a double bass, and the keyboard on the
	// screen shows the low octaves.
	Bass(on bool)

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
