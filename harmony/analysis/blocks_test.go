package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// blocksOf renders the blocks of changes as a musician would write
// them: the chords of the block, an arrow to the chord it lands on ("?"
// when it lands nowhere), then the tonality it announces, named as the
// analyse command does: see tonalityName.
func blocksOf(c analysis.Changes) string {
	return blocksName(c, analysis.Blocks(c, analysis.Approaches(c)))
}

// blocksName renders blocks as blocksOf does.
func blocksName(c analysis.Changes, blocks []analysis.Block) string {
	var out []string
	for _, b := range blocks {
		var chords []string
		for _, i := range []int{b.Two, b.Sus, b.Five} {
			if i >= 0 {
				chords = append(chords, changeName(c.Chords[i]))
			}
		}
		target := "?"
		if b.Target >= 0 {
			target = changeName(c.Chords[b.Target])
		}
		out = append(out, strings.Join(chords, " ")+" → "+target+" : "+tonalityName(b.Announced))
	}
	return strings.Join(out, "   ")
}

func TestBlocks(t *testing.T) {
	const (
		c, db, d, eb, e, f, g, ab, a harmony.PitchClass = 0, 1, 2, 3, 4, 5, 7, 8, 9
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
		dim7    = harmony.ChordDiminishedSeventh
	)
	var (
		sus, _      = harmony.NewChordPattern(0, 5, 7, 10)
		flatNine, _ = harmony.NewChordPattern(0, 4, 7, 10, 13)
		nine, _     = harmony.NewChordPattern(0, 4, 7, 10, 14)
		alt, _      = harmony.NewChordPattern(0, 4, 6, 10, 13, 15, 20)
	)
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		"a two five one": {
			changesOf(false, d, min7, g, dom7, c, maj7),
			"Dm7 G7 → Cmaj7 : C",
		},
		// A V7/II, then the two five one.
		"nested": {
			changesOf(false, a, dom7, d, min7, g, dom7, c, maj7),
			"A7 → Dm7 : Dm harm/mel   Dm7 G7 → Cmaj7 : C",
		},
		// Dm7 is the target of one block and the two of the next. A
		// minor seventh two onto a minor target: the melodic minor.
		"the double role": {
			changesOf(false, e, min7, a, dom7, d, min7, g, dom7, c, maj7),
			"Em7 A7 → Dm7 : Dm mel   Dm7 G7 → Cmaj7 : C",
		},
		// What Is This Thing Called Love: the two five announces F
		// harmonic minor, and lands on F major.
		"a minor two five onto a major chord": {
			changesOf(false, g, halfDim, c, dom7, f, maj7),
			"Gm7♭5 C7 → Fmaj7 : Fm harm",
		},
		"the suspension belongs to the block": {
			changesOf(false, d, min7, g, sus, g, dom7, c, maj7),
			"Dm7 G7sus4 G7 → Cmaj7 : C",
		},
		// Tenderly: a two five of D flat that does not resolve, and
		// announces its tonality all the same.
		"without resolution": {
			changesOf(false, eb, min7, ab, dom7, f, min7),
			"E♭m7 A♭7 → ? : D♭ M/m mel",
		},
		"a dominant alone going nowhere": {
			changesOf(false, c, maj7, g, dom7, f, maj7),
			"",
		},
		// The deceptive cadence of C, Dm7 G7 Am7, read alone is the
		// aeolian cadence of A minor, IVm7 ♭VII7 Im7: a plagal block.
		// Only the tonic of the passage tells them apart.
		"deceptive, or aeolian": {
			changesOf(false, d, min7, g, dom7, a, min7),
			"Dm7 G7 → Am7 : Am",
		},
		// One Finger Snap: C7alt takes its colours from outside F, and
		// its tritone still announces it, the two keeping the harmonic
		// minor.
		"an altered dominant": {
			changesOf(false, g, halfDim, c, alt, f, halfDim),
			"Gm7♭5 C7alt → Fm7♭5 : Fm harm",
		},
		// A7♭9 holds only D harmonic minor, and keeps it before D major;
		// E9 holds A major and A melodic minor, and the target decides.
		"the colours of the V decide without a two": {
			changesOf(false, c, maj7, a, flatNine, d, maj7, e, nine, a, maj7),
			"A7♭9 → Dmaj7 : Dm harm   E9 → Amaj7 : A",
		},
		"the target decides without colours": {
			changesOf(false, a, dom7, d, min7),
			"A7 → Dm7 : Dm harm/mel",
		},
		// The chromatic dominant and its two are not degrees of C: the
		// target alone decides.
		"a chromatic dominant with its two": {
			changesOf(false, ab, min7, db, dom7, c, maj7),
			"A♭m7 D♭7 → Cmaj7 : C",
		},
		"a diminished chord has no two": {
			changesOf(false, e, min7, db, dim7, d, min7),
			"D♭dim7 → Dm7 : Dm harm",
		},
		// A bare G7 onto a suspended chord: neither colour nor third.
		"nothing decides": {
			changesOf(false, g, dom7, c, sus),
			"G7 → C7sus4 : C M/m harm/mel",
		},
	} {
		if got := blocksOf(tc.changes); got != tc.want {
			t.Errorf("%s:\n got  %s\n want %s", name, got, tc.want)
		}
	}
}
