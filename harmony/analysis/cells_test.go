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
// the I it promises, a deceptive cadence (the book's "V – …"), and
// "♭II" one whose X7 stand for their tritone twins.
func TestCells(t *testing.T) {
	const (
		c, db, d, eb, e, f, gb, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
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
		// The first V goes on to the III, the mediant: it stands for the
		// I in the cell that follows, but the V does not resolve there.
		"Have You Met Miss Jones": {
			changesOf(true, f, maj7, d, min7, g, min7, c, dom7, a, min7, d, min7, g, min7, c, dom7),
			"anatole… III-VI-II-V",
		},
		// Anthropology, bars 5-6: a III VI II V of D flat, heard as
		// such, its resolution avoided.
		"a III VI II V avoiding its I": {
			changesOf(true, f, min7, bb, dom7, eb, dom7, ab, dom7, d, min7), "III-VI-II-V…",
		},
		// A Beautiful Friendship: a II V I of D, then its IV7.
		"a major chord is no II": {
			changesOf(false, e, min7, a, dom7, d, maj7, g, dom7), "",
		},
		// Tritone substitutions, read on the chords they replace.
		"the V substituted": {
			changesOf(true, c, maj7, a, dom7, d, dom7, db, dom7), "anatole ♭II",
		},
		// Blue In Green, Too Young: III ♭III7 II V.
		"the VI substituted": {
			changesOf(true, e, min7, eb, dom7, d, min7, g, dom7, c, maj7), "III-VI-II-V ♭II",
		},
		// Take The A Train: I II7 IIm7 V, whose D7 would be the twin of
		// A♭7, a ♭VI7 that is not the V of the II.
		"a II7 is no VI": {
			changesOf(true, c, maj7, d, dom7, d, min7, g, dom7), "",
		},
		// Body And Soul, bars 23-24: the fifths become half tones.
		"a chromatic descent": {
			changesOf(false, d, min7, g, dom7, c, dom7, b, dom7, bb, dom7, eb, min7), "III-VI-II-V ♭II",
		},
		"all three substituted": {
			changesOf(true, c, maj7, eb, dom7, ab, dom7, db, dom7), "anatole ♭II",
		},
		// Sophisticated Lady restores C7 F7 B♭7 E♭7 under G♭7 F7 E7 E♭7,
		// all dominants: without a tonic chord first, no anatole.
		"a chromatic descent of dominants": {
			changesOf(true, gb, dom7, f, dom7, e, dom7, eb, dom7, ab, maj7), "",
		},
		// Not on a tonic chord: the cycle alone is no anatole.
		"a cycle of fifths": {
			changesOf(true, e, dom7, a, dom7, d, dom7, g, dom7), "",
		},
	} {
		var got []string
		for _, cl := range analysis.Cells(tc.changes) {
			name := cl.Kind.String()
			if cl.Substituted {
				name += " ♭II"
			}
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
