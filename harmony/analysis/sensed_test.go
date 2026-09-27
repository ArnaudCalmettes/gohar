package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// sensedOf writes what the ear expects after each chord, one chord to
// a line: the ground it is heard in, the local tonic a cadence has
// just tonicised, and what the cadence being played awaits.
//
//	Gm7♭5 in E♭, awaiting Fm harm
func sensedOf(c analysis.Changes, signature []harmony.Tonality) string {
	var out []string
	for i, s := range analysis.Sense(c, analysis.Blocks(c, analysis.Approaches(c)), signature) {
		line := changeName(c.Chords[i]) + " in "
		if s.Ground == nil {
			line += "?"
		} else {
			line += tonalityName(s.Ground)
		}
		if s.Local != nil {
			line += ", on " + tonalityName(s.Local)
		}
		if s.Awaited != nil {
			line += ", awaiting " + tonalityName(s.Awaited)
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func TestSense(t *testing.T) {
	const (
		c, db, d, eb, f, g, ab, a, bb harmony.PitchClass = 0, 1, 2, 3, 5, 7, 8, 9, 10
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		min6    = harmony.ChordMinorSixth
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	flatNine, _ := harmony.NewChordPattern(0, 4, 7, 10, 13)
	for name, tc := range map[string]struct {
		changes   analysis.Changes
		signature []harmony.Tonality
		want      []string
	}{
		// Tenderly, bars 1 to 9, as docs/grilles.md hears them. E♭m7
		// A♭7 awaits nothing: its two is the tonic, borrowed. The
		// cadence onto Fm7♭5 tonicises F minor for one chord, and
		// Fm7♭5 is already the two of E flat minor.
		"Tenderly": {
			changesOf(false, eb, maj7, ab, dom7, eb, min7, ab, dom7, f, min7, db, dom7, eb, maj7,
				g, halfDim, c, flatNine, f, halfDim, bb, flatNine),
			analysis.MajorTonalities(eb),
			[]string{
				"E♭maj7 in E♭",
				"A♭7 in E♭",
				"E♭m7 in E♭",
				"A♭7 in E♭",
				"Fm7 in E♭",
				"D♭7 in E♭",
				"E♭maj7 in E♭",
				"Gm7♭5 in E♭, awaiting Fm harm",
				"C7♭9 in E♭, awaiting Fm harm",
				"Fm7♭5 in E♭, on Fm harm, awaiting E♭m harm",
				"B♭7♭9 in E♭, awaiting E♭m harm",
			},
		},
		// The same two five on the tonic, resolving this time: it
		// tonicises D flat.
		"a two five on the tonic that resolves": {
			changesOf(false, eb, maj7, eb, min7, ab, dom7, db, maj7),
			analysis.MajorTonalities(eb),
			[]string{
				"E♭maj7 in E♭",
				"E♭m7 in E♭, awaiting D♭",
				"A♭7 in E♭, awaiting D♭",
				"D♭maj7 in E♭, on D♭",
			},
		},
		// Autumn Leaves without a signature: a two chord installs
		// nothing, the first cadence installs B flat, and G minor comes
		// as a local tonic. Setting when it becomes the ground is the
		// threshold of a modulation, still to come.
		"Autumn Leaves, the relative first": {
			changesOf(false, c, min7, f, dom7, bb, maj7, eb, maj7, a, halfDim, d, dom7, g, min6),
			nil,
			[]string{
				"Cm7 in ?, awaiting B♭",
				"F7 in ?, awaiting B♭",
				"B♭maj7 in B♭",
				"E♭maj7 in B♭",
				"Am7♭5 in B♭, awaiting Gm harm",
				"D7 in B♭, awaiting Gm harm",
				"Gm6 in B♭, on Gm harm",
			},
		},
	} {
		if got, want := sensedOf(tc.changes, tc.signature), strings.Join(tc.want, "\n"); got != want {
			t.Errorf("%s:\ngot\n%s\nwant\n%s", name, got, want)
		}
	}
}

// The tonality of the tune is read at its end: Autumn Leaves, begun in
// B flat, ends in G minor.
func TestTune(t *testing.T) {
	const (
		c, d, eb, f, g, a, bb harmony.PitchClass = 0, 2, 3, 5, 7, 9, 10
	)
	leaves := changesOf(false, c, harmony.ChordMinorSeventh, f, harmony.ChordDominantSeventh,
		bb, harmony.ChordMajorSeventh, eb, harmony.ChordMajorSeventh,
		a, harmony.ChordHalfDiminished, d, harmony.ChordDominantSeventh, g, harmony.ChordMinorSixth)
	sensed := analysis.Sense(leaves, analysis.Blocks(leaves, analysis.Approaches(leaves)), nil)
	if got := tonalityName(analysis.Tune(leaves, sensed)); got != "Gm nat/harm/mel" {
		t.Errorf("%s, want Gm nat/harm/mel", got)
	}
}
