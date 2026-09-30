package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The anatole and the III VI II V, in the forms En Harmonie gives
// (tome 1, chapter 8 §3.3, and chapter 9 for the secondary dominants),
// looping back to their first chord, and the modal cadence of the
// aeolian (tome 2, chapter 2 §5.2); "…" marks a cell whose V avoids
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
	type bar = []any
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
		// The fifths become half tones, a chord a bar: B♭7 comes on the
		// first beat, an arrival.
		"a chromatic descent": {
			changesOf(false, d, min7, g, dom7, c, dom7, b, dom7, bb, dom7, eb, min7), "III-VI-II-V ♭II",
		},
		// Autumn Leaves, bars 25 to 28: Am7♭5 | D7 | Gm7 G♭7 changes its
		// rhythm, and is no cell; Gm7 G♭7 | Fm7 E7 is one, G♭7 for C7
		// and E7 for B♭7.
		"a steady harmonic rhythm": {
			barsOf(false, bar{a, halfDim}, bar{d, dom7}, bar{g, min7, gb, dom7}, bar{f, min7, e, dom7}, bar{eb, maj7}),
			"III-VI-II-V ♭II",
		},
		// Body And Soul, bars 23 and 24, as played: B♭7 comes in the
		// middle of the bar, one more link in a chain of dominants
		// walking down to E♭m.
		"a chain of dominants": {
			barsOf(false, bar{d, min7, g, dom7}, bar{c, dom7, b, dom7, bb, dom7}, bar{eb, min7}), "",
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
		// The II V I of C, its approach replaced by the modal cadence of
		// C aeolian, as En Harmonie does it.
		"the modal cadence of C aeolian": {
			changesOf(false, ab, maj7, bb, dom7, c, maj7), "♭VI-♭VII-I",
		},
		// Guile's Theme (Street Fighter II), on its minor I: the
		// cadence in its own mode.
		"Guile's Theme": {
			changesOf(false, a, maj7, b, dom7, db, min7), "♭VI-♭VII-I",
		},
		// Route 209 (Pokémon Diamond and Pearl), on triads: the
		// borrowing that resolves in major, the Mario Cadence.
		"Route 209": {
			changesOf(false, f, maj, g, dom7, a, maj), "♭VI-♭VII-I",
		},
		// Fm7 B♭7 Cmaj7: a disguised minor plagal, which is a block.
		"the ♭VII7 I without its ♭VI": {
			changesOf(false, f, min7, bb, dom7, c, maj7), "",
		},
		// A♭maj7 B♭maj7 Cmaj7: parallel chords, the ♭VII no dominant.
		"a ♭VIImaj7 is no ♭VII7": {
			changesOf(false, ab, maj7, bb, maj7, c, maj7), "",
		},
		// The cadence ends on the I an anatole starts from.
		"the cadence, then an anatole": {
			changesOf(false, ab, maj7, bb, dom7, c, maj7, a, min7, d, min7, g, dom7, c, maj7),
			"♭VI-♭VII-I anatole",
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
