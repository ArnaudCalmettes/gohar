package analysis_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The plagal blocks read again on the ground: a V going to its VI, or
// back to its two, told from the ♭VII7 Im and the dorian IV7.
func TestReread(t *testing.T) {
	const (
		c, d, e, f, g, a, bb harmony.PitchClass = 0, 2, 4, 5, 7, 9, 10
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		// En Harmonie's C7 Dm7, a fifth lower, with a II: Dm7 G7 Am7
		// in C is II V VI.
		"in C, a V going to its VI": {
			barsOf(false, bar{c, maj7}, bar{d, min7, g, dom7}, bar{a, min7}, bar{d, min7, g, dom7}, bar{c, maj7}),
			"Dm7 G7 → ? : C M/m mel   Dm7 G7 → Cmaj7 : C",
		},
		// That Old Feeling, bars 7 to 9: the V alone is a block still.
		"without its two": {
			barsOf(false, bar{c, maj7}, bar{f, min7}, bar{g, dom7}, bar{a, min7}, bar{d, min7, g, dom7}, bar{c, maj7}),
			"G7 → ? : C M/m harm/mel   Dm7 G7 → Cmaj7 : C",
		},
		// In A minor, G7 Am7 is the aeolian ♭VII7 Im: it stands.
		"in A minor, the plagal reading stands": {
			barsOf(false, bar{a, min7}, bar{d, min7, g, dom7}, bar{a, min7}, bar{e, dom7}, bar{a, min7}),
			"Dm7 G7 → Am7 : Am   E7 → Am7 : Am harm/mel",
		},
		// Satin Doll: a two five played again, in C and then in D.
		"a V going back to its two": {
			barsOf(false, bar{d, min7, g, dom7}, bar{d, min7, g, dom7}, bar{e, min7, a, dom7}, bar{e, min7, a, dom7},
				bar{a, min7, d, dom7}, bar{d, min7, g, dom7}, bar{c, maj7}),
			"Dm7 G7 → ? : C M/m mel   Dm7 G7 → ? : C M/m mel   Em7 A7 → ? : D M/m mel   Em7 A7 → ? : D M/m mel   Am7 D7 → ? : G M/m mel   Dm7 G7 → Cmaj7 : C",
		},
		// Longer than four bars, a vamp: Am7 D7 in C is A dorian, not
		// a two five of G played again.
		"a vamp": {
			barsOf(false, bar{c, maj7}, bar{a, min7, d, dom7}, bar{a, min7, d, dom7}, bar{a, min7, d, dom7}, bar{a, min7, d, dom7},
				bar{a, min7, d, dom7}, bar{a, min7}, bar{d, min7, g, dom7}, bar{c, maj7}),
			"D7 → Am7 : Am mel   D7 → Am7 : Am mel   D7 → Am7 : Am mel   D7 → Am7 : Am mel   D7 → Am7 : Am mel   Dm7 G7 → Cmaj7 : C",
		},
		// Last Train Home: Gm7 Dm7 in B flat is a plagal IVm7 Im7, and
		// Gm7 no V going back to its two.
		"a IVm7 is no V": {
			barsOf(false, bar{bb, maj7}, bar{d, min7}, bar{g, min7}, bar{d, min7}, bar{g, min7}, bar{f, dom7}, bar{bb, maj7}),
			"Gm7 → Dm7 : Dm nat/harm   F7 → B♭maj7 : B♭",
		},
		// Mas Que Nada: in F minor, B♭7 is the dorian IV7 of Fm7.
		"the dorian IV7 stands": {
			barsOf(false, bar{f, min7, bb, dom7}, bar{f, min7, bb, dom7}, bar{f, min7}, bar{g, halfDim, c, dom7}, bar{f, min7}),
			"B♭7 → Fm7 : Fm mel   B♭7 → Fm7 : Fm mel   Gm7♭5 C7 → Fm7 : Fm harm",
		},
	} {
		kinds := analysis.Approaches(tc.changes)
		blocks := analysis.Blocks(tc.changes, kinds)
		phrases := analysis.Phrases(tc.changes, blocks)
		tune := analysis.Tune(tc.changes, phrases)
		sensed := analysis.Sense(tc.changes, blocks, phrases, tune)
		_, blocks = analysis.Reread(tc.changes, kinds, blocks, sensed, tune)
		if got := blocksName(tc.changes, blocks); got != tc.want {
			t.Errorf("%s:\n got %s\nwant %s", name, got, tc.want)
		}
	}
}
