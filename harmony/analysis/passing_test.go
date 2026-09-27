package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// passingOf writes the chords with an arrow after each passing chord,
// the way its bass walks, and spells a chord walking up with a sharp:
// "Fmaj7 F♯dim7↗ Gm7".
func passingOf(c analysis.Changes) string {
	var out []string
	for i, walk := range analysis.PassingChords(c) {
		name := changeName(c.Chords[i])
		switch {
		case walk > 0:
			name = sharpened(name) + "↗"
		case walk < 0:
			name += "↘"
		}
		out = append(out, name)
	}
	return strings.Join(out, " ")
}

func TestPassingChords(t *testing.T) {
	const (
		c, db, d, eb, e, f, gb, g, ab, a, bb harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
	)
	const (
		maj7 = harmony.ChordMajorSeventh
		min7 = harmony.ChordMinorSeventh
		dim7 = harmony.ChordDiminishedSeventh
		maj  = harmony.ChordMajorTriad
	)
	// inverted puts a chord over another bass.
	inverted := func(ch analysis.Changes, i int, bass harmony.PitchClass) analysis.Changes {
		ch.Chords[i].Bass = bass
		return ch
	}
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		// Mean to Me: I ♯Idim7 II ♯IIdim7 III.
		"going up": {
			changesOf(false, f, maj7, gb, dim7, g, min7, ab, dim7, a, min7),
			"Fmaj7 F♯dim7↗ Gm7 G♯dim7↗ Am7",
		},
		// Someday My Prince Will Come: B♭/D D♭dim7 Cm7.
		"going down, from an inverted chord": {
			inverted(changesOf(false, bb, maj, db, dim7, c, min7), 0, d),
			"B♭/D D♭dim7↘ Cm7",
		},
		// The same without the inversion: no walk, no passing chord.
		"no walk without the inversion": {
			changesOf(false, bb, maj, db, dim7, c, min7),
			"B♭ D♭dim7 Cm7",
		},
		// Across the loop: E♭dim7 walks from the Em7 before it down to
		// the Dm7 the chorus starts again on.
		"across the loop": {
			changesOf(true, d, min7, f, maj7, e, min7, eb, dim7),
			"Dm7 Fmaj7 Em7 E♭dim7↘",
		},
		"a leap is not a passage": {
			changesOf(false, c, maj7, gb, dim7, g, min7),
			"Cmaj7 F♯dim7 Gm7",
		},
	} {
		if got := passingOf(tc.changes); got != tc.want {
			t.Errorf("%s:\n got  %s\n want %s", name, got, tc.want)
		}
	}
}
