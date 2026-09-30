package analysis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// areaName writes a tonal area as a musician reads it: "C, bars 5 to
// 8, 4 bars, leaving D, the first tonality, 2 steps".
func areaName(c analysis.Changes, a analysis.TonalArea) string {
	last := c.Chords[a.To]
	out := fmt.Sprintf("%s, bars %d to %d, %g bars, leaving %s", tonalityName(a.Tonic),
		c.Bar(c.Chords[a.From].Start)+1, c.Bar(last.Start+last.Length-1)+1, a.Bars, tonalityName(a.Leaves))
	if a.First {
		out += ", the first tonality"
	}
	if a.Opens {
		out += ", opening a section"
	}
	if a.Closes {
		out += ", closing a section"
	}
	return out + fmt.Sprintf(", %d steps", a.Distance)
}

// Tune Up: D major, then C and B flat, each a tone lower, two steps
// each on the cycle of fifths, twice, at bars 5 to 12 as En Harmonie
// reads them (tome 1, p. 160).
func TestTonalAreas(t *testing.T) {
	const c, d, e, f, g, a, bb harmony.PitchClass = 0, 2, 4, 5, 7, 9, 10
	const (
		maj7 = harmony.ChordMajorSeventh
		min7 = harmony.ChordMinorSeventh
		dom7 = harmony.ChordDominantSeventh
	)
	type bar = []any
	tuneUp := barsOf(true, bar{e, min7}, bar{a, dom7}, bar{d, maj7}, bar{d, maj7},
		bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{c, maj7},
		bar{c, min7}, bar{f, dom7}, bar{bb, maj7}, bar{g, min7},
		bar{e, min7}, bar{f, dom7}, bar{bb, maj7}, bar{a, dom7},
		bar{e, min7}, bar{a, dom7}, bar{d, maj7}, bar{d, maj7},
		bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{c, maj7},
		bar{c, min7}, bar{f, dom7}, bar{bb, maj7}, bar{g, min7},
		bar{e, min7}, bar{a, dom7}, bar{d, maj7}, bar{d, maj7})
	blocks := analysis.Blocks(tuneUp, analysis.Approaches(tuneUp))
	phrases := analysis.Phrases(tuneUp, blocks)
	sensed := analysis.Sense(tuneUp, blocks, phrases, analysis.Tune(tuneUp, phrases))
	var got []string
	for _, ar := range analysis.TonalAreas(tuneUp, blocks, sensed, analysis.Sections(tuneUp)) {
		got = append(got, areaName(tuneUp, ar))
	}
	want := []string{
		"C, bars 5 to 8, 4 bars, leaving D, the first tonality, 2 steps",
		"B♭, bars 9 to 12, 4 bars, leaving C, 2 steps",
		"C, bars 21 to 24, 4 bars, leaving D, the first tonality, 2 steps",
		"B♭, bars 25 to 28, 4 bars, leaving C, 2 steps",
	}
	if g, w := strings.Join(got, "\n"), strings.Join(want, "\n"); g != w {
		t.Errorf("got\n%s\nwant\n%s", g, w)
	}
}
