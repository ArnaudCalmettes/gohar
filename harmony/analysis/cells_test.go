package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The anatole and the III VI II V, in the forms En Harmonie gives
// (tome 1, chapter 8 §3.3, and chapter 9 for the secondary dominants),
// looping back to their first chord; "…" marks a cell whose V avoids
// the I it promises.
func TestCells(t *testing.T) {
	const (
		c, db, d, eb, e, f, g, ab, a, bb harmony.PitchClass = 0, 1, 2, 3, 4, 5, 7, 8, 9, 10
	)
	const (
		maj     = harmony.ChordMajorTriad
		minor   = harmony.ChordMinorTriad
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		"I Got Rhythm, in B flat": {
			changesOf(true, bb, maj, g, min7, c, min7, f, dom7), "anatole",
		},
		"in F harmonic minor": {
			changesOf(true, f, minor, db, maj7, g, halfDim, c, dom7), "anatole",
		},
		"Softly, in C minor": {
			changesOf(true, c, minor, a, halfDim, d, halfDim, g, dom7), "anatole",
		},
		"a minor tonic written m7": {
			changesOf(true, c, min7, a, halfDim, d, halfDim, g, dom7), "anatole",
		},
		"a secondary dominant on the VI": {
			changesOf(true, c, maj7, a, dom7, d, min7, g, dom7), "anatole",
		},
		"and on the II": {
			changesOf(true, c, maj7, a, dom7, d, dom7, g, dom7), "anatole",
		},
		"Have You Met Miss Jones": {
			changesOf(true, f, maj7, d, min7, g, min7, c, dom7, a, min7, d, min7, g, min7, c, dom7),
			"anatole III-VI-II-V",
		},
		// Anthropology, bars 5-6: a III VI II V of D flat, heard as
		// such, its resolution avoided.
		"a III VI II V avoiding its I": {
			changesOf(true, f, min7, bb, dom7, eb, dom7, ab, dom7, d, min7), "III-VI-II-V…",
		},
		// Not on a tonic chord: the cycle alone is no anatole.
		"a cycle of fifths": {
			changesOf(true, e, dom7, a, dom7, d, dom7, g, dom7), "",
		},
	} {
		var got []string
		for _, cl := range analysis.Cells(tc.changes) {
			name := cl.Kind.String()
			if !cl.Resolves {
				name += "…"
			}
			got = append(got, name)
		}
		if g := strings.Join(got, " "); g != tc.want {
			t.Errorf("%s: %q, want %q", name, g, tc.want)
		}
	}
}
