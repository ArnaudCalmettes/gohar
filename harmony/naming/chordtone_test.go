package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
)

// A note heard over a chord takes the letter of its interval from the
// root, as on a lead sheet.
func TestSpellOver(t *testing.T) {
	F := naming.SpelledNote{Letter: naming.LetterF, Accidental: naming.NaturalSign}
	Bb := naming.SpelledNote{Letter: naming.LetterB, Accidental: naming.FlatSign}
	C := naming.SpelledNote{Letter: naming.LetterC, Accidental: naming.NaturalSign}
	dim7 := harmony.ChordPattern(0x000249) // 0, 3, 6, 9

	for _, tc := range []struct {
		what  string
		root  naming.SpelledNote
		p     harmony.ChordPattern
		class harmony.PitchClass
		want  string
	}{
		{"the minor seventh of F7", F, harmony.ChordDominantSeventh, 3, "Eb"},
		{"the third of F7", F, harmony.ChordDominantSeventh, 9, "A"},
		{"the ninth of B♭7", Bb, harmony.ChordDominantSeventh, 0, "C"},
		{"the third of B♭7", Bb, harmony.ChordDominantSeventh, 2, "D"},
		{"the minor seventh of B♭7", Bb, harmony.ChordDominantSeventh, 8, "Ab"},
		{"a ♭10 over F7, never a ♯9", F, harmony.ChordDominantSeventh, 8, "Ab"},
		{"a ♭13 over C7, never a ♯5", C, harmony.ChordDominantSeventh, 8, "Ab"},
		{"the ♯11 over C7", C, harmony.ChordDominantSeventh, 6, "F#"},
		{"the ♭5 of Cm7♭5", C, harmony.ChordHalfDiminished, 6, "Gb"},
		{"the diminished seventh of Cdim7", C, dim7, 9, "Bbb"},
		{"the 13th over C7", C, harmony.ChordDominantSeventh, 9, "A"},
	} {
		got := naming.English.Name(naming.SpellOver(tc.root, tc.p, tc.class), naming.ASCII)
		assert.Equal(t, tc.want, got, tc.what)
	}
}
