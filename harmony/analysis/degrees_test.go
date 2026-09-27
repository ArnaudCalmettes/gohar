package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

func degreesOf(c analysis.Changes, home []harmony.Tonality) string {
	kinds := analysis.Approaches(c)
	var out []string
	for _, d := range analysis.Degrees(c, analysis.Blocks(c, kinds), analysis.PassingChords(c), home) {
		out = append(out, d.String())
	}
	return strings.Join(out, " ")
}

func TestDegrees(t *testing.T) {
	const (
		c, db, d, eb, f, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 5, 7, 8, 9, 10, 11
	)
	const (
		maj     = harmony.ChordMajorTriad
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
		dim7    = harmony.ChordDiminishedSeventh
	)
	flatNine, _ := harmony.NewChordPattern(0, 4, 7, 10, 13)
	for name, tc := range map[string]struct {
		changes analysis.Changes
		home    []harmony.Tonality
		want    string
	}{
		// Tenderly, bars 1 to 16. Against En Harmonie: bars 3 and 4 read
		// here as a two five of D flat that does not resolve, where
		// En Harmonie reads a borrowed I and a IV; bars 13 and 14 as a
		// two five of B flat, where En Harmonie reads VI and II7. Both
		// readings are true; En Harmonie's come with the cadences.
		"Tenderly": {
			changesOf(false, eb, maj7, ab, dom7, eb, min7, ab, dom7, f, min7, db, dom7, eb, maj7,
				g, halfDim, c, flatNine, f, halfDim, bb, dom7, f, halfDim, bb, dom7, b, dim7,
				c, min7, f, dom7, f, min7, bb, dom7),
			analysis.MajorTonalities(eb),
			"I IV7 II V II ♭VII7 I II V II V II V ♯Vdim7 II V II V",
		},
		// A chromatic dominant and its two, the chromatic subdominant.
		"♭VIm7 ♭II7 I": {
			changesOf(false, d, min7, ab, min7, db, dom7, c, maj7),
			analysis.MajorTonalities(c),
			"II ♭VIm7 ♭II7 I",
		},
		// In F minor, D♭maj7 is VI, not ♭VI.
		"a degree of the minor": {
			changesOf(false, f, harmony.ChordMinorTriad, db, maj7, g, halfDim, c, flatNine),
			analysis.MinorTonalities(f),
			"I VI II V",
		},
		// Someday My Prince Will Come: B♭/D D♭dim7 Cm7, a passing chord
		// written flat, going down; and the I in first inversion.
		"going down": {
			func() analysis.Changes {
				ch := changesOf(false, bb, maj, db, dim7, c, min7)
				ch.Chords[0].Bass = d
				return ch
			}(),
			analysis.MajorTonalities(bb),
			"I/3 ♭IIIdim7 II",
		},
		// A secondary dominant alone is written on its degree, with its
		// quality: VI7, the V7/II.
		"a secondary dominant": {
			changesOf(false, c, maj7, a, dom7, d, min7, g, dom7, c, maj7),
			analysis.MajorTonalities(c),
			"I VI7 II V I",
		},
	} {
		if got := degreesOf(tc.changes, tc.home); got != tc.want {
			t.Errorf("%s:\n got %s\nwant %s", name, got, tc.want)
		}
	}
}
