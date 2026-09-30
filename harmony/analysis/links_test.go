package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The consecutive II-Vs of En Harmonie, tome 1, chapter 8 §3.3, each
// link written as the step from one II-V to the next.
func TestLinks(t *testing.T) {
	const (
		c, db, d, eb, e, f, gb, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
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
		"a half tone up": {
			barsOf(false, bar{d, min7, g, dom7}, bar{eb, min7, ab, dom7}), "½",
		},
		"a half tone down": {
			barsOf(false, bar{d, min7, g, dom7}, bar{db, min7, gb, dom7}), "½",
		},
		"a tone between the twos": {
			barsOf(false, bar{d, min7, g, dom7}, bar{e, min7, a, dom7}), "step",
		},
		"a tone from the V": {
			barsOf(false, bar{d, min7, g, dom7}, bar{a, min7, d, dom7}), "step",
		},
		// Autumn Leaves: the twos a tone apart, whatever the chromatic
		// dominant between them.
		"a tone, by chromatic dominants": {
			barsOf(false, bar{g, min7, gb, dom7}, bar{f, min7, e, dom7}), "step",
		},
		// Autumn Leaves, bars 25 to 28: D7 falls a fifth to Gm7, but G♭7
		// breaks the cycle, and the twos are a tone apart.
		"a fifth, the cycle broken": {
			barsOf(false, bar{a, halfDim}, bar{d, dom7}, bar{g, min7, gb, dom7}, bar{f, min7, e, dom7}),
			"step step",
		},
		"the cycle of fifths": {
			barsOf(false, bar{d, min7, g, dom7}, bar{c, min7, f, dom7}), "5th",
		},
		// Fifths, then a half tone, toward Cmaj7.
		"combined": {
			barsOf(false, bar{gb, min7, b, dom7}, bar{e, min7, a, dom7}, bar{ab, min7, db, dom7}, bar{c, maj7}),
			"5th ½",
		},
		"Confirmation": {
			barsOf(false, bar{f, maj7}, bar{e, halfDim, a, dom7}, bar{d, min7, g, dom7}, bar{c, min7, f, dom7}),
			"5th 5th",
		},
		// The bridge of the rhythm changes, Anthropology: dominants
		// down the cycle of fifths, two bars each.
		"a chain of dominants": {
			barsOf(false, bar{d, dom7}, bar{d, dom7}, bar{g, dom7}, bar{g, dom7},
				bar{c, dom7}, bar{c, dom7}, bar{f, dom7}, bar{f, dom7}, bar{bb, maj7}),
			"5th 5th 5th",
		},
		"a chain of chromatic dominants": {
			barsOf(false, bar{e, dom7}, bar{eb, dom7}, bar{d, dom7}, bar{db, dom7}, bar{c, maj7}),
			"½ ½ ½",
		},
		// A V going to a two is no chain of dominants.
		"a V to a two": {
			barsOf(false, bar{g, dom7}, bar{c, min7, f, dom7}, bar{bb, maj7}), "",
		},
		// Satin Doll: a tone, then a half tone, a II-V played again
		// between them.
		"Satin Doll": {
			barsOf(false, bar{d, min7, g, dom7}, bar{d, min7, g, dom7}, bar{e, min7, a, dom7}, bar{e, min7, a, dom7},
				bar{a, min7, d, dom7}, bar{ab, min7, db, dom7}, bar{c, maj7}),
			"step ½",
		},
	} {
		ch := tc.changes
		var got []string
		for _, l := range analysis.Links(ch, analysis.Approaches(ch)) {
			got = append(got, l.Step.String())
		}
		if g := strings.Join(got, " "); g != tc.want {
			t.Errorf("%s: %q, want %q", name, g, tc.want)
		}
	}
}
