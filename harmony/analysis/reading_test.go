package analysis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// readingName writes what the tonality of a tune rests on as a musician
// reads it: "first D7, not the tonic, the first cadence goes to Cm;
// last Cm6 at bar 14, Cm: Cm".
func readingName(c analysis.Changes, r analysis.TuneReading) string {
	var parts []string
	if r.Opens >= 0 {
		first := "not the tonic, no cadence resolves"
		switch {
		case r.First != nil:
			first = tonalityName(r.First)
		case r.Cadence != nil:
			first = "not the tonic, the first cadence goes to " + tonalityName(r.Cadence)
		}
		parts = append(parts, fmt.Sprintf("first %s, %s", changeName(c.Chords[r.Opens]), first))
	}
	if r.Stops >= 0 {
		parts = append(parts, fmt.Sprintf("last %s at bar %d, %s", changeName(c.Chords[r.Stops]),
			c.Bar(c.Chords[r.Stops].Start)+1, tonalityName(r.Last)))
	}
	if r.Heard != nil {
		f, l := r.First[0].Tonic(), r.Last[0].Tonic()
		parts = append(parts, fmt.Sprintf("heard %g bars against %g", r.Heard[f], r.Heard[l]))
	}
	return strings.Join(parts, "; ") + ": " + tonalityName(r.Tonality)
}

// What the tonality of a tune rests on (En Harmonie, tome 1, chapter 8
// §1.2, p. 99 and 100): its first chord and its last one, and when
// they part, how long each of their tonics is heard.
func TestReadTune(t *testing.T) {
	const c, d, eb, f, g, ab, bb harmony.PitchClass = 0, 2, 3, 5, 7, 8, 10
	const (
		min6    = harmony.ChordMinorSixth
		min7    = harmony.ChordMinorSeventh
		maj7    = harmony.ChordMajorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		// Jordu opens on D7 G7, a V of V, and comes back to Cm6 at bar 6
		// of each section, the weak conclusive cadence Siron finds in it
		// (La partition intérieure, p. 390). Played once, without going
		// round, it stops on the second.
		"Jordu": {
			barsOf(false, bar{d, dom7, g, dom7}, bar{c, min6}, bar{f, dom7, bb, dom7}, bar{eb, maj7},
				bar{d, dom7, g, dom7}, bar{c, min6}, bar{ab, dom7}, bar{ab, dom7},
				bar{d, dom7, g, dom7}, bar{c, min6}, bar{f, dom7, bb, dom7}, bar{eb, maj7},
				bar{d, dom7, g, dom7}, bar{c, min6}, bar{ab, dom7}, bar{ab, dom7}),
			"first D7, not the tonic, the first cadence goes to Cm nat/harm/mel; last Cm6 at bar 14, Cm nat/harm/mel: Cm nat/harm/mel",
		},
		// My Funny Valentine starts on Cm and ends on E♭: the two part,
		// and C minor, heard the longer, is the tonality of the tune. The
		// E♭ of the last bar forms no tonal area: it counts for the tonic
		// around it.
		"first and last part": {
			barsOf(false, bar{c, min6}, bar{d, halfDim, g, dom7}, bar{c, min7}, bar{c, min6},
				bar{f, min7}, bar{d, halfDim, g, dom7}, bar{c, min7}, bar{bb, dom7}, bar{eb, maj7}),
			"first Cm6, Cm nat/harm/mel; last E♭maj7 at bar 9, E♭; heard 9 bars against 0: Cm nat/harm/mel",
		},
	} {
		ch := tc.changes
		blocks := analysis.Blocks(ch, analysis.Approaches(ch))
		if got := readingName(ch, analysis.ReadTune(ch, analysis.Phrases(ch, blocks))); got != tc.want {
			t.Errorf("%s:\n got %s\nwant %s", name, got, tc.want)
		}
	}
}
