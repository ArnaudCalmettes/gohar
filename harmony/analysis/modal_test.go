package analysis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// plagesText writes each modal plage as its chord and its bars: "Dm7,
// bars 1 to 16".
func plagesText(c analysis.Changes, plages []analysis.Plage) string {
	var out []string
	for _, p := range plages {
		from, to := c.Bar(c.Chords[p.From].Start)+1, c.Bar(c.Chords[p.To].Start+c.Chords[p.To].Length-1)+1
		out = append(out, fmt.Sprintf("%s, bars %d to %d", changeName(c.Chords[p.From]), from, to))
	}
	return strings.Join(out, "; ")
}

// A modal plage is a chord held four bars or more with no cadence to it
// or from it. A tune is modal when its plages take half of it.
func TestModal(t *testing.T) {
	const c, d, eb, f, g, bb harmony.PitchClass = 0, 2, 3, 5, 7, 10
	const (
		maj7 = harmony.ChordMajorSeventh
		min7 = harmony.ChordMinorSeventh
		dom7 = harmony.ChordDominantSeventh
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		plages  string
		modal   bool
	}{
		// So What: Dm7 sixteen bars, E♭m7 eight, Dm7 eight.
		"So What": {
			barsOf(true, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7},
				bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7},
				bar{eb, min7}, bar{eb, min7}, bar{eb, min7}, bar{eb, min7}, bar{eb, min7}, bar{eb, min7}, bar{eb, min7}, bar{eb, min7},
				bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}, bar{d, min7}),
			"Dm7, bars 1 to 16; E♭m7, bars 17 to 24; Dm7, bars 25 to 32",
			true,
		},
		// A tonal tune whose tonic holds four bars after its cadence is
		// not modal: the tonic is what the cadence led to.
		"a tonic held four bars": {
			barsOf(true, bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{c, maj7}, bar{c, maj7}, bar{c, maj7},
				bar{f, maj7}, bar{d, min7, g, dom7}),
			"",
			false,
		},
		// Held chords with no cadence: a plage each, but only one of
		// them takes half the tune.
		"a vamp before a tune": {
			barsOf(true, bar{bb, min7}, bar{bb, min7}, bar{bb, min7}, bar{bb, min7},
				bar{c, maj7}, bar{d, min7}, bar{g, dom7}, bar{c, maj7}),
			"B♭m7, bars 1 to 4",
			true,
		},
	} {
		c := tc.changes
		plages := analysis.Modal(c, analysis.Blocks(c, analysis.Approaches(c)))
		if got := plagesText(c, plages); got != tc.plages {
			t.Errorf("%s: %q, want %q", name, got, tc.plages)
		}
		if got := analysis.IsModal(c, plages); got != tc.modal {
			t.Errorf("%s: modal %v, want %v", name, got, tc.modal)
		}
	}
}
