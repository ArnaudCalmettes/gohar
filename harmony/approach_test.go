package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// Each chord of a progression against the next, as En Harmonie reads
// them in its chapter on preparing a chord.
func TestApproachOf(t *testing.T) {
	const (
		c, db, d, eb, e, f, gb, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
		fs                                      harmony.PitchClass = gb
	)
	var (
		major7  = []harmony.Semitones{0, 4, 7, 11}
		minor7  = []harmony.Semitones{0, 3, 7, 10}
		halfDim = []harmony.Semitones{0, 3, 6, 10}
		seventh = []harmony.Semitones{0, 4, 7, 10}
		sus     = []harmony.Semitones{0, 5, 7, 10}
		alt     = []harmony.Semitones{0, 4, 6, 10, 13, 15, 20}
		dim7    = []harmony.Semitones{0, 3, 6, 9}
		minor9  = []harmony.Semitones{0, 3, 7, 10, 14}
	)
	const (
		none  = harmony.NoApproach
		dom   = harmony.DominantApproach
		chrom = harmony.ChromaticApproach
		dimin = harmony.DiminishedApproach
		susp  = harmony.SuspensionApproach
		two   = harmony.TwoApproach
		plag  = harmony.PlagalApproach
	)
	type step struct {
		chord harmony.Chord
		kind  harmony.ApproachKind // what it does for the next
	}
	for name, steps := range map[string][]step{
		// Every chord preceded by its secondary dominant, back to the
		// start.
		"a dominant before each chord": {
			{chordOf(t, c, major7...), none},
			{chordOf(t, a, seventh...), dom},
			{chordOf(t, d, minor7...), none},
			{chordOf(t, b, seventh...), dom},
			{chordOf(t, e, minor7...), none},
			{chordOf(t, c, seventh...), dom},
			{chordOf(t, f, major7...), none},
			{chordOf(t, d, seventh...), dom},
			{chordOf(t, g, seventh...), none},
			{chordOf(t, e, seventh...), dom},
			{chordOf(t, a, minor7...), none},
			{chordOf(t, fs, seventh...), dom},
			{chordOf(t, b, halfDim...), none},
			{chordOf(t, g, seventh...), dom},
			{chordOf(t, c, major7...), none},
		},
		// The Night We Called It a Day: three chromatic dominants for a
		// bass that walks down. Each chord before them is a two, the
		// two of the dominant each one stands in for: En Harmonie's
		// original grid is C♯m7♭5 F♯7 | Bm7 E7 | Am7 D7 | Gmaj7.
		"chromatic dominants": {
			{chordOf(t, db, halfDim...), two},
			{chordOf(t, c, seventh...), chrom},
			{chordOf(t, b, minor7...), two},
			{chordOf(t, bb, seventh...), chrom},
			{chordOf(t, a, minor7...), two},
			{chordOf(t, ab, seventh...), chrom},
			{chordOf(t, g, major7...), none},
		},
		// Sophisticated Lady: the surface, a dominant a semitone above
		// the next each time, then a dominant a fifth above.
		"a chain of dominants": {
			{chordOf(t, gb, seventh...), chrom},
			{chordOf(t, f, seventh...), chrom},
			{chordOf(t, e, seventh...), chrom},
			{chordOf(t, eb, seventh...), dom},
			{chordOf(t, ab, major7...), none},
		},
		// No tritone, or the tritone going elsewhere.
		"not an approach": {
			{chordOf(t, g, seventh...), none}, // to the dominant a tone below
			{chordOf(t, f, seventh...), plag}, // ♭VII7-I: a plagal, see below
			{chordOf(t, g, major7...), none}, // ♭VIImaj7: a modal cadence, not read
			{chordOf(t, a, minor7...), none}, // a two needs a dominant after
			{chordOf(t, d, major7...), none},
		},
		// Cmaj7 C♯dim7 Dm7 D♯dim7 Em7: each diminished chord a semitone
		// under the next is its dominant seventh flat nine without root.
		"diminished chords": {
			{chordOf(t, c, major7...), none},
			{chordOf(t, db, dim7...), dimin},
			{chordOf(t, d, minor7...), none},
			{chordOf(t, eb, dim7...), dimin},
			{chordOf(t, e, minor7...), none},
		},
		// The same four notes, whatever root the chart writes: Ddim7 is
		// Bdim7, and prepares C as well.
		"a diminished chord by another name": {
			{chordOf(t, d, dim7...), dimin},
			{chordOf(t, c, major7...), none},
		},
		// Going down, B♭/D D♭dim7 Cm7: a passing chord, which two chords
		// cannot see, and no dominant of Cm7.
		"a diminished chord going down": {
			{chordOf(t, db, dim7...), none},
			{chordOf(t, c, minor7...), none},
		},
		// Em7 A7sus4 A7 Dm9 G7sus4 G7 Cmaj7: twos, suspensions, dominants.
		"twos and suspensions": {
			{chordOf(t, e, minor7...), none}, // the sus4 has no tritone
			{chordOf(t, a, sus...), susp},
			{chordOf(t, a, seventh...), dom},
			{chordOf(t, d, minor9...), none},
			{chordOf(t, g, sus...), susp},
			{chordOf(t, g, seventh...), dom},
			{chordOf(t, c, major7...), none},
		},
		// Dm7 | A♭m7 D♭7 | Cmaj7, then A♭m7 straight onto G7: the
		// chromatic subdominant, the two of the chromatic dominant.
		"chromatic subdominants": {
			{chordOf(t, d, minor7...), none},
			{chordOf(t, ab, minor7...), two},
			{chordOf(t, db, seventh...), chrom},
			{chordOf(t, c, major7...), none},
			{chordOf(t, ab, halfDim...), two},
			{chordOf(t, g, seventh...), dom},
			{chordOf(t, c, major7...), none},
		},
		// Plagal cadences: the IV of any quality before the tonic, and
		// the ♭VII7, a minor plagal in disguise, prepared by its IVm7.
		// Not before a chord that cannot be a tonic.
		"plagal cadences": {
			{chordOf(t, f, major7...), plag},
			{chordOf(t, c, major7...), none},
			{chordOf(t, f, minor7...), plag},
			{chordOf(t, c, minor7...), none},
			{chordOf(t, f, minor7...), two},
			{chordOf(t, bb, seventh...), plag},
			{chordOf(t, c, major7...), none},
			{chordOf(t, f, seventh...), plag},
			{chordOf(t, c, major7...), none},
			{chordOf(t, f, major7...), none}, // not before a dominant
			{chordOf(t, c, seventh...), none},
		},
		// The modal ♭VII-I, B♭maj7 from C mixolydian, B♭m7 from C
		// phrygian, are not read from two chords: they are the stepwise
		// motion of any tonal chart.
		"modal cadences, not read": {
			{chordOf(t, bb, major7...), none},
			{chordOf(t, c, major7...), none},
			{chordOf(t, bb, minor7...), none},
			{chordOf(t, c, minor7...), none},
		},
		"a two five one": {
			{chordOf(t, d, minor7...), two},
			{chordOf(t, g, seventh...), dom},
			{chordOf(t, c, major7...), none},
		},
		"the altered dominant is one": {
			{chordOf(t, g, alt...), dom},
			{chordOf(t, c, major7...), none},
		},
	} {
		for i := 0; i+1 < len(steps); i++ {
			if got := harmony.ApproachOf(steps[i].chord, steps[i+1].chord); got != steps[i].kind {
				t.Errorf("%s, chord %d: %v, want %v", name, i+1, got, steps[i].kind)
			}
		}
	}
}
