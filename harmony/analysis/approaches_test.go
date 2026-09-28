package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// approachesOf writes each chord with what it does for the next in
// brackets, as the analyse command labels it: "Dm7 (II) G7 (V) Cmaj7".
func approachesOf(c analysis.Changes, kinds []harmony.ApproachKind) string {
	var out []string
	for i, ch := range c.Chords {
		name := changeName(ch)
		if kinds[i] != harmony.NoApproach {
			name += " (" + kinds[i].String() + ")"
		}
		out = append(out, name)
	}
	return strings.Join(out, " ")
}

func TestApproaches(t *testing.T) {
	const (
		c, db, d, eb, e, f, gb, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	var silence any
	for name, tc := range map[string]struct {
		changes    analysis.Changes
		approaches string
	}{
		// The turnaround prepares the first chord, across the loop.
		"I VI II V, looping": {
			changesOf(true, c, maj7, a, min7, d, min7, g, dom7),
			"Cmaj7 Am7 Dm7 (II) G7 (V)",
		},
		// A performance: what G7 prepares is not known yet.
		"a two five being played": {
			changesOf(false, d, min7, g, dom7),
			"Dm7 (II) G7",
		},
		"a silence prepares nothing": {
			changesOf(false, c, dom7, f, silence, f, maj7),
			"C7 N.C. Fmaj7",
		},
		// En Harmonie, a secondary dominant before each chord.
		"a dominant before each chord": {
			changesOf(true, c, maj7, a, dom7, d, min7, b, dom7, e, min7, c, dom7, f, maj7,
				d, dom7, g, dom7, e, dom7, a, min7, gb, dom7, b, halfDim, g, dom7),
			"Cmaj7 A7 (V) Dm7 B7 (V) Em7 C7 (V) Fmaj7 D7 (V) G7 E7 (V) Am7 F♯7 (V) Bm7♭5 G7 (V)",
		},
		// Dominants round the cycle of fifths forever: nothing lands.
		"a cycle that never lands": {
			changesOf(true, a, dom7, d, dom7, g, dom7, c, dom7, f, dom7, bb, dom7,
				eb, dom7, ab, dom7, db, dom7, gb, dom7, b, dom7, e, dom7),
			"A7 (V) D7 (V) G7 (V) C7 (V) F7 (V) B♭7 (V) E♭7 (V) A♭7 (V) D♭7 (V) F♯7 (V) B7 (V) E7 (V)",
		},
	} {
		kinds := analysis.Approaches(tc.changes)
		if got := approachesOf(tc.changes, kinds); got != tc.approaches {
			t.Errorf("%s, approaches:\n got  %s\n want %s", name, got, tc.approaches)
		}
	}
}
