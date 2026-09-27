package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// degreesOf writes both readings of a sequence: on the installed
// tonic, and bracketed as En Harmonie prints it.
func degreesOf(c analysis.Changes, home []harmony.Tonality) (string, string) {
	blocks := analysis.Blocks(c, analysis.Approaches(c))
	sensed := analysis.Sense(c, blocks, home)
	passing := analysis.PassingChords(c)
	write := func(ds []analysis.Degree) string {
		var out []string
		for _, d := range ds {
			out = append(out, d.String())
		}
		return strings.Join(out, " ")
	}
	return write(analysis.Degrees(c, passing, sensed)), write(analysis.Bracketed(c, blocks, passing, sensed))
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
		ground  string // on the installed tonic
		bracket string // as En Harmonie prints it
	}{
		// Tenderly, bars 1 to 16. On the tonic, bar 8 is the three six
		// of a three six two five one; bracketed, as En Harmonie prints
		// it, a two five of F minor. Bars 3 and 4, E♭m7 A♭7, a two five
		// that does not resolve with its two on the tonic, are a
		// borrowed I and a IV7 either way. Bars 13 and 14 are VI II7 on
		// the tonic, as En Harmonie reads them, but a bracketed two five
		// of B flat: its reading will come with the cadences.
		"Tenderly": {
			changesOf(false, eb, maj7, ab, dom7, eb, min7, ab, dom7, f, min7, db, dom7, eb, maj7,
				g, halfDim, c, flatNine, f, halfDim, bb, dom7, f, halfDim, bb, dom7, b, dim7,
				c, min7, f, dom7, f, min7, bb, dom7),
			analysis.MajorTonalities(eb),
			"I IV7 Im7 IV7 II ♭VII7 I IIIm7♭5 VI7 IIm7♭5 V IIm7♭5 V ♯Vdim7 VI II7 II V",
			"I IV7 Im7 IV7 II ♭VII7 I II V II V II V ♯Vdim7 II V II V",
		},
		// There Will Never Be Another You: a two five tonicises the VI.
		"a two five toward the VI": {
			changesOf(false, eb, maj7, d, halfDim, g, dom7, c, min7),
			analysis.MajorTonalities(eb),
			"I VII III7 VI",
			"I II V VI",
		},
		// A chromatic dominant and its two, the chromatic subdominant.
		"♭VIm7 ♭II7 I": {
			changesOf(false, d, min7, ab, min7, db, dom7, c, maj7),
			analysis.MajorTonalities(c),
			"II ♭VIm7 ♭II7 I",
			"II ♭VIm7 ♭II7 I",
		},
		// In F minor, D♭maj7 is VI, not ♭VI.
		"a degree of the minor": {
			changesOf(false, f, harmony.ChordMinorTriad, db, maj7, g, halfDim, c, flatNine),
			analysis.MinorTonalities(f),
			"I VI II V",
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
			"I/3 ♭IIIdim7 II",
		},
		// A secondary dominant alone is written on its degree, with its
		// quality: VI7, the V7/II.
		"a secondary dominant": {
			changesOf(false, c, maj7, a, dom7, d, min7, g, dom7, c, maj7),
			analysis.MajorTonalities(c),
			"I VI7 II V I",
			"I VI7 II V I",
		},
	} {
		ground, bracket := degreesOf(tc.changes, tc.home)
		if ground != tc.ground {
			t.Errorf("%s, on the tonic:\n got  %s\n want %s", name, ground, tc.ground)
		}
		if bracket != tc.bracket {
			t.Errorf("%s, bracketed:\n got  %s\n want %s", name, bracket, tc.bracket)
		}
	}
}
