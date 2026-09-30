package analysis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// conclusionName writes how a section ends as a musician says it: the
// tonic, the bar, strong or weak, and the chord the turnaround starts
// on. "C at bar 7, strong, then A7".
func conclusionName(c analysis.Changes, k analysis.Conclusion) string {
	if k.Arrives < 0 {
		return "none"
	}
	out := fmt.Sprintf("%s at bar %d", tonalityName(k.Tonic), c.Bar(c.Chords[k.Arrives].Start)+1)
	if k.Strong {
		out += ", strong"
	} else {
		out += ", weak"
	}
	if k.Loop >= 0 {
		out += ", then " + changeName(c.Chords[k.Loop])
	}
	return out
}

// Three sections of 8 bars: the first lands on its I at bar 7 and turns
// around from there, as Let's Cool One; the second at bar 8, as Lover
// Man; the third stops on its V, a half cadence (Siron, La partition
// intérieure, p. 390).
func TestConclusions(t *testing.T) {
	const c, d, e, f, g, a harmony.PitchClass = 0, 2, 4, 5, 7, 9
	const (
		maj7 = harmony.ChordMajorSeventh
		six  = harmony.ChordMajorSixth
		min7 = harmony.ChordMinorSeventh
		min6 = harmony.ChordMinorSixth
		dom7 = harmony.ChordDominantSeventh
	)
	type bar = []any
	ch := barsOf(true,
		bar{c, maj7}, bar{a, min7}, bar{d, min7}, bar{g, dom7},
		bar{e, min7, a, dom7}, bar{d, min7, g, dom7}, bar{c, maj7, a, dom7}, bar{d, min7, g, dom7},

		bar{c, maj7}, bar{a, min7}, bar{d, min7}, bar{g, dom7},
		bar{e, min7}, bar{a, dom7}, bar{d, min7, g, dom7}, bar{c, six},

		bar{f, maj7}, bar{f, min6}, bar{e, min7}, bar{a, dom7},
		bar{d, min7}, bar{d, min7}, bar{g, dom7}, bar{g, dom7})
	sections := []analysis.Section{{Label: "A", From: 0, Bars: 8}, {Label: "A", From: 8, Bars: 8}, {Label: "B", From: 16, Bars: 8}}
	var got []string
	for _, k := range analysis.Conclusions(ch, sections, analysis.Blocks(ch, analysis.Approaches(ch))) {
		got = append(got, conclusionName(ch, k))
	}
	want := []string{"C at bar 7, strong, then A7", "C at bar 16, weak", "none"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// A VI leant on is no conclusion. Yesterdays walks down the cycle onto
// B♭maj7 at bar 14, by F13 alone, and Em7♭5 A7 goes back to Dm: the
// section ends open on A7. Rosetta comes down the same kind of cycle
// onto F6, but on the strong bar 7 of its section, and concludes there
// before Bm7♭5 E7 goes to the Am of its bridge.
func TestLeaningOnTheSixth(t *testing.T) {
	const c, d, eb, e, f, g, a, bb, b harmony.PitchClass = 0, 2, 3, 4, 5, 7, 9, 10, 11
	const (
		min     = harmony.ChordMinorTriad
		minMaj7 = harmony.ChordMinorMajorSeventh
		maj7    = harmony.ChordMajorSeventh
		six     = harmony.ChordMajorSixth
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	eight := []analysis.Section{{Label: "A", From: 0, Bars: 8}, {Label: "B", From: 8, Bars: 8}}
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    []string
	}{
		"Yesterdays": {
			barsOf(true, bar{d, min}, bar{e, halfDim, a, dom7}, bar{d, min}, bar{e, halfDim, a, dom7},
				bar{d, min, d, minMaj7}, bar{d, min7}, bar{b, halfDim}, bar{e, dom7},
				bar{a, dom7}, bar{d, dom7}, bar{g, dom7}, bar{c, dom7},
				bar{f, dom7}, bar{bb, maj7}, bar{e, halfDim}, bar{a, dom7}),
			[]string{"none", "none"},
		},
		"Rosetta": {
			barsOf(true, bar{f, six}, bar{e, dom7}, bar{eb, dom7}, bar{d, dom7},
				bar{g, dom7}, bar{c, dom7}, bar{f, six}, bar{b, halfDim, e, dom7},
				bar{a, min}, bar{b, halfDim, e, dom7}, bar{a, min}, bar{d, min7, g, dom7},
				bar{c, maj7, a, min7}, bar{d, min7, g, dom7}, bar{g, min7}, bar{c, dom7}),
			[]string{"F at bar 7, strong, then Bm7♭5", "none"},
		},
	} {
		ch := tc.changes
		var got []string
		for _, k := range analysis.Conclusions(ch, eight, analysis.Blocks(ch, analysis.Approaches(ch))) {
			got = append(got, conclusionName(ch, k))
		}
		if strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
			t.Errorf("%s: got\n%s\nwant\n%s", name, strings.Join(got, "\n"), strings.Join(tc.want, "\n"))
		}
	}
}
