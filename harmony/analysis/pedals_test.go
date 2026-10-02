package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// pedalsOf names the pedals of changes as a musician would: "B♭ ped.,
// dominant, Fm7/B♭ to E♭maj7/B♭", one a line.
func pedalsOf(c analysis.Changes, sections []analysis.Section) string {
	kinds := map[analysis.PedalKind]string{
		analysis.OtherPedal:    "",
		analysis.TonicPedal:    ", tonic",
		analysis.DominantPedal: ", dominant",
	}
	var out []string
	for _, p := range analysis.Pedals(c, analysis.Grounds(c, sense(c)), sections) {
		line := noteNames[p.Bass] + " ped." + kinds[p.Kind]
		if p.General {
			line += ", general"
		}
		out = append(out, line+", "+changeName(c.Chords[p.From])+" to "+changeName(c.Chords[p.To]))
	}
	return strings.Join(out, "\n")
}

// over sets the bass of changes `at` to `bass`.
func over(c analysis.Changes, bass harmony.PitchClass, at ...int) analysis.Changes {
	for _, i := range at {
		c.Chords[i].Bass = bass
	}
	return c
}

func TestPedals(t *testing.T) {
	const c, d, eb, f, g, bb harmony.PitchClass = 0, 2, 3, 5, 7, 10
	const (
		maj7 = harmony.ChordMajorSeventh
		min7 = harmony.ChordMinorSeventh
		dom7 = harmony.ChordDominantSeventh
		half = harmony.ChordHalfDiminished
		min6 = harmony.ChordMinorSixth
		minT = harmony.ChordMinorTriad
	)
	for _, tc := range []struct {
		name string
		c    analysis.Changes
		want string
	}{{
		// The dominant pedal held onto the tonic, « entendu renversé
		// sur sa 5te » (En Harmonie, tome 2, chapter 5 §1.5).
		"dominant", over(changesOf(false, f, min7, bb, dom7, eb, maj7), bb, 0, 2),
		"B♭ ped., dominant, Fm7/B♭ to E♭maj7/B♭",
	}, {
		// The ostinato of Stolen Moments, before the cadence that
		// installs C minor.
		"tonic", over(changesOf(false, c, min7, d, min7, eb, maj7, d, min7, d, half, g, dom7, c, min6), c, 1, 2, 3),
		"C ped., tonic, Cm7 to Dm7/C",
	}, {
		// One chord written again over its third is no pedal.
		"one chord", over(changesOf(false, c, minT, c, minT, f, min7, bb, dom7, eb, maj7), eb, 0, 1),
		"",
	}} {
		if got := pedalsOf(tc.c, nil); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

// A pedal over a whole section is general, one over part of it passing.
func TestGeneralPedal(t *testing.T) {
	const c, d, eb, g harmony.PitchClass = 0, 2, 3, 7
	const (
		maj7 = harmony.ChordMajorSeventh
		min7 = harmony.ChordMinorSeventh
		dom7 = harmony.ChordDominantSeventh
		min6 = harmony.ChordMinorSixth
	)
	type bar = []any
	ostinato := over(barsOf(false, bar{c, min7}, bar{d, min7}, bar{eb, maj7}, bar{d, min7},
		bar{d, min7}, bar{g, dom7}, bar{c, min6}, bar{c, min6}), c, 1, 2, 3)
	sections := []analysis.Section{{Label: "A", From: 0, Bars: 4}, {Label: "B", From: 4, Bars: 4}}
	if got, want := pedalsOf(ostinato, sections), "C ped., tonic, general, Cm7 to Dm7/C"; got != want {
		t.Errorf("over a section: got %q, want %q", got, want)
	}
	sections = []analysis.Section{{Label: "A", From: 0, Bars: 8}}
	if got, want := pedalsOf(ostinato, sections), "C ped., tonic, Cm7 to Dm7/C"; got != want {
		t.Errorf("over half a section: got %q, want %q", got, want)
	}
}
