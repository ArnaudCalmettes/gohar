package analysis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The evidence a tune gives for each tonic, at which chords: Jordu
// comes back to C minor more often than it goes to E flat.
func TestCandidates(t *testing.T) {
	const c, d, eb, f, g, ab, bb harmony.PitchClass = 0, 2, 3, 5, 7, 8, 10
	const (
		min6 = harmony.ChordMinorSixth
		maj7 = harmony.ChordMajorSeventh
		dom7 = harmony.ChordDominantSeventh
	)
	type bar = []any
	jordu := barsOf(false, bar{d, dom7, g, dom7}, bar{c, min6}, bar{f, dom7, bb, dom7}, bar{eb, maj7},
		bar{d, dom7, g, dom7}, bar{c, min6}, bar{ab, dom7}, bar{ab, dom7},
		bar{d, dom7, g, dom7}, bar{c, min6}, bar{f, dom7, bb, dom7}, bar{eb, maj7},
		bar{d, dom7, g, dom7}, bar{c, min6}, bar{ab, dom7}, bar{ab, dom7})
	blocks := analysis.Blocks(jordu, analysis.Approaches(jordu))
	var got []string
	for _, cd := range analysis.Candidates(jordu, blocks, analysis.Phrases(jordu, blocks)) {
		var proofs []string
		for _, e := range cd.Evidence {
			proofs = append(proofs, fmt.Sprintf("%s %s", e.Proof, changeName(jordu.Chords[e.At])))
		}
		got = append(got, tonalityName(cd.Tonic)+": "+strings.Join(proofs, ", "))
	}
	want := []string{
		"Cm nat/harm/mel: returns Cm6, returns Cm6, returns Cm6, stops Cm6",
		"E♭: returns E♭maj7, returns E♭maj7",
	}
	if g, w := strings.Join(got, "\n"), strings.Join(want, "\n"); g != w {
		t.Errorf("got\n%s\nwant\n%s", g, w)
	}
}
