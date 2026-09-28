package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// phrasesText writes each phrase as its chords, then where it concludes:
// "→ Gm" when it comes to rest, "stops on A♭" where the tune stops, and
// nothing for the tail.
func phrasesText(c analysis.Changes) string {
	var lines []string
	for _, p := range phrasesOf(c) {
		var chords []string
		for i := p.From; i <= p.To; i++ {
			chords = append(chords, changeName(c.Chords[i]))
		}
		line := strings.Join(chords, " ")
		switch {
		case p.Tonic == nil:
			line += " (tail)"
		case p.Stops:
			line += ", stops on " + tonalityName(p.Tonic)
		default:
			line += " → " + tonalityName(p.Tonic)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func TestPhrases(t *testing.T) {
	const (
		c, d, eb, e, f, g, ab, a, bb harmony.PitchClass = 0, 2, 3, 4, 5, 7, 8, 9, 10
	)
	const (
		min     = harmony.ChordMinorTriad
		min6    = harmony.ChordMinorSixth
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    []string
	}{
		// B♭maj7 is passed through; Gm6, held two bars, is where the
		// phrase rests.
		"Autumn Leaves": {
			barsOf(true, bar{c, min7}, bar{f, dom7}, bar{bb, maj7}, bar{eb, maj7},
				bar{a, halfDim}, bar{d, dom7}, bar{g, min6}, bar{g, min6}),
			[]string{"Cm7 F7 B♭maj7 E♭maj7 Am7♭5 D7 Gm6 Gm6 → Gm nat/harm/mel"},
		},
		// Back to its opening Fm at bar 3; it stops on A♭maj7, and Gm7♭5
		// C7 after it is the tail, back to the first chord.
		"Lullaby Of Birdland": {
			barsOf(true, bar{f, min, d, halfDim}, bar{g, dom7, c, dom7}, bar{f, min}, bar{bb, min7, eb, dom7},
				bar{c, min7, f, min7}, bar{bb, min7, eb, dom7}, bar{ab, maj7}, bar{g, halfDim, c, dom7}),
			[]string{
				"Fm Dm7♭5 G7 C7 Fm → Fm nat/harm/mel",
				"B♭m7 E♭7 Cm7 Fm7 B♭m7 E♭7 A♭maj7, stops on A♭",
				"Gm7♭5 C7 (tail)",
			},
		},
	} {
		if got := phrasesText(tc.changes); got != strings.Join(tc.want, "\n") {
			t.Errorf("%s:\ngot\n%s\nwant\n%s", name, got, strings.Join(tc.want, "\n"))
		}
	}
}
