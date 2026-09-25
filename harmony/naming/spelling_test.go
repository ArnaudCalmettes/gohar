package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func namerIn(t *testing.T, tonic harmony.PitchClass, p harmony.ScalePattern, spelled naming.SpelledNote) *naming.Namer {
	t.Helper()

	tonality, err := harmony.NewTonality(tonic, p)
	require.NoError(t, err)

	n, err := naming.NewNamer(naming.English)
	require.NoError(t, err)
	return n.WithSpelledTonality(tonality, spelled)
}

func sharp(l naming.Letter) naming.SpelledNote {
	return naming.SpelledNote{Letter: l, Accidental: naming.SharpSign}
}

func natural(l naming.Letter) naming.SpelledNote {
	return naming.SpelledNote{Letter: l, Accidental: naming.NaturalSign}
}

func flat(l naming.Letter) naming.SpelledNote {
	return naming.SpelledNote{Letter: l, Accidental: naming.FlatSign}
}

// One letter per degree, in order. Everything else in this file is a
// consequence of that single rule.
func TestEachLetterAppearsOnce(t *testing.T) {
	patterns := map[string]harmony.ScalePattern{
		"major":          harmony.ScaleMajor,
		"natural minor":  harmony.ScaleNaturalMinor,
		"harmonic minor": harmony.ScaleHarmonicMinor,
		"melodic minor":  harmony.ScaleMelodicMinor,
	}

	for name, p := range patterns {
		for tonic := harmony.PitchClass(0); tonic < harmony.PitchClassCount; tonic++ {
			t.Run(name+" spells seven distinct letters from every tonic", func(t *testing.T) {
				tonality, err := harmony.NewTonality(tonic, p)
				require.NoError(t, err)

				n, err := naming.NewNamer(naming.English)
				require.NoError(t, err)

				seen := map[naming.Letter]bool{}
				for _, note := range n.WithTonality(tonality).Scale() {
					assert.False(t, seen[note.Letter],
						"letter reused in %s from tonic %d", name, tonic)
					seen[note.Letter] = true
				}
				assert.Len(t, seen, 7)
			})
		}
	}
}

// The awkward keys, where the rule produces spellings a naive
// implementation gets wrong.
func TestSpellingIsForcedNotChosen(t *testing.T) {
	t.Run("F sharp major reaches E sharp rather than F", func(t *testing.T) {
		n := namerIn(t, 6, harmony.ScaleMajor, sharp(naming.LetterF))
		assert.Equal(t,
			[]naming.SpelledNote{
				sharp(naming.LetterF), sharp(naming.LetterG), sharp(naming.LetterA),
				natural(naming.LetterB), sharp(naming.LetterC), sharp(naming.LetterD),
				sharp(naming.LetterE),
			},
			n.Scale(),
			"the F letter is taken by the tonic, so the seventh degree must be E sharp",
		)
	})

	t.Run("G sharp major reaches F double sharp", func(t *testing.T) {
		n := namerIn(t, 8, harmony.ScaleMajor, sharp(naming.LetterG))

		got := n.Scale()
		require.Len(t, got, 7)
		assert.Equal(t,
			naming.SpelledNote{Letter: naming.LetterF, Accidental: naming.DoubleSharpSign},
			got[6],
			"a double sharp is legitimate, not a sign the algorithm broke",
		)
	})

	t.Run("B flat major spells flats, never sharps", func(t *testing.T) {
		n := namerIn(t, 10, harmony.ScaleMajor, flat(naming.LetterB))
		for _, note := range n.Scale() {
			assert.NotEqual(t, naming.SharpSign, note.Accidental)
		}
	})

	t.Run("the same class spells differently in two keys", func(t *testing.T) {
		inG := namerIn(t, 7, harmony.ScaleMajor, natural(naming.LetterG))
		inDFlat := namerIn(t, 1, harmony.ScaleMajor, flat(naming.LetterD))

		assert.Equal(t, sharp(naming.LetterF), inG.Note(6),
			"the raised fourth of G major")
		assert.Equal(t, flat(naming.LetterG), inDFlat.Note(6),
			"the fourth of D flat major: one pitch class, two spellings, "+
				"and only the context decides",
		)
	})
}

// A mode name says which degree is altered, and the spelling has to
// agree with it. Lydian sharp 2 must reach a raised second and not a
// lowered third, because the third letter is already spoken for.
func TestModeNameAgreesWithSpelling(t *testing.T) {
	tonality, err := harmony.NewTonality(0, mustLookupPattern(t, harmony.HarmonicMinor, 6))
	require.NoError(t, err)

	n, err := naming.NewNamer(naming.English)
	require.NoError(t, err)

	got := n.WithSpelledTonality(tonality, natural(naming.LetterC)).Scale()
	require.Len(t, got, 7)

	assert.Equal(t, sharp(naming.LetterD), got[1],
		"the second degree is a raised second, not a lowered third")
	assert.Equal(t, natural(naming.LetterE), got[2],
		"the third letter is taken by the third degree")
}

// With no tonality there is still a spelling, and it is not a key.
func TestDefaultSpellingIsNotAKey(t *testing.T) {
	n, err := naming.NewNamer(naming.English)
	require.NoError(t, err)

	t.Run("the seven natural classes take their own letter", func(t *testing.T) {
		for class, letter := range map[harmony.PitchClass]naming.Letter{
			0: naming.LetterC, 2: naming.LetterD, 4: naming.LetterE,
			5: naming.LetterF, 7: naming.LetterG, 9: naming.LetterA,
			11: naming.LetterB,
		} {
			assert.Equal(t, natural(letter), n.Note(class))
		}
	})

	t.Run("the five others take a sharp", func(t *testing.T) {
		for _, class := range []harmony.PitchClass{1, 3, 6, 8, 10} {
			assert.Equal(t, naming.SharpSign, n.Note(class).Accidental)
		}
	})

	t.Run("no tonality is claimed", func(t *testing.T) {
		assert.True(t, n.Tonality().IsZero(),
			"the fallback is a spelling convention, never a key the player is in")
		assert.Nil(t, n.Scale(),
			"there is no scale to spell, and nil says so better than seven notes of C")
	})
}

func TestEnharmonic(t *testing.T) {
	t.Run("E sharp and F sound alike and are written differently", func(t *testing.T) {
		assert.True(t, sharp(naming.LetterE).Enharmonic(natural(naming.LetterF)))
	})

	t.Run("a note is not enharmonic with itself", func(t *testing.T) {
		assert.False(t, natural(naming.LetterF).Enharmonic(natural(naming.LetterF)),
			"enharmonic means same sound and different writing")
	})

	t.Run("C sharp and D flat sound alike", func(t *testing.T) {
		assert.True(t, sharp(naming.LetterC).Enharmonic(flat(naming.LetterD)))
		assert.Equal(t, sharp(naming.LetterC).Class(), flat(naming.LetterD).Class())
	})
}

func TestLocaleRendersInItsOwnWords(t *testing.T) {
	t.Run("the English locale writes letters and Unicode signs", func(t *testing.T) {
		assert.Equal(t, "F\u266f", naming.English.Name(sharp(naming.LetterF)))
	})

	t.Run("the French locale writes solfège syllables", func(t *testing.T) {
		assert.Equal(t, "fa dièse", naming.French.Name(sharp(naming.LetterF)))
	})

	t.Run("a natural carries no sign in either locale", func(t *testing.T) {
		assert.Equal(t, "C", naming.English.Name(natural(naming.LetterC)))
		assert.Equal(t, "do", naming.French.Name(natural(naming.LetterC)))
	})
}

func mustLookupPattern(t *testing.T, s harmony.System, d harmony.Degree) harmony.ScalePattern {
	t.Helper()
	m, ok := naming.Lookup(s, d)
	require.True(t, ok)
	return m.Pattern()
}
