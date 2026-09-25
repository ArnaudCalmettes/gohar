package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func modeAt(t *testing.T, s harmony.System, d harmony.Degree) naming.Mode {
	t.Helper()
	m, ok := naming.Lookup(s, d)
	require.True(t, ok)
	return m
}

// The systematic name is the default in both languages: the base mode
// and each altered degree with its sign, nothing special cased. The
// refined names are alternatives, tested with the spoken register.
func TestModeNameIsSystematic(t *testing.T) {
	tests := []struct {
		what    string
		system  harmony.System
		degree  harmony.Degree
		english string
		french  string
		words   string // French, in the Words notation
	}{
		{
			what:   "a natural mode renders as its bare name",
			system: harmony.NaturalMajor, degree: 4,
			english: "lydian", french: "lydien", words: "lydien",
		},
		{
			what:   "a raised second",
			system: harmony.HarmonicMinor, degree: 6,
			english: "lydian \u266f2", french: "lydien \u266f2", words: "lydien dièse 2",
		},
		{
			what:   "a lowered third, not what it makes of the mode",
			system: harmony.MelodicMinor, degree: 1,
			english: "ionian \u266d3", french: "ionien \u266d3", words: "ionien bémol 3",
		},
		{
			what:   "a natural sixth keeps its sign",
			system: harmony.MelodicMinor, degree: 2,
			english: "phrygian \u266e6", french: "phrygien \u266e6", words: "phrygien bécarre 6",
		},
		{
			what:   "a raised fifth, not augmented",
			system: harmony.MelodicMinor, degree: 3,
			english: "lydian \u266f5", french: "lydien \u266f5", words: "lydien dièse 5",
		},
		{
			what:   "a double flat",
			system: harmony.HarmonicMajor, degree: 7,
			english: "locrian \u266d\u266d7", french: "locrien \u266d\u266d7", words: "locrien double bémol 7",
		},
		{
			what:   "two alterations, in the order the catalogue holds",
			system: harmony.HarmonicMajor, degree: 6,
			english: "lydian \u266f2 \u266f5", french: "lydien \u266f2 \u266f5", words: "lydien dièse 2 dièse 5",
		},
	}

	for _, tc := range tests {
		m := modeAt(t, tc.system, tc.degree)

		t.Run(tc.what+", in English", func(t *testing.T) {
			assert.Equal(t, tc.english, naming.English.ModeName(m, naming.Signs))
		})

		t.Run(tc.what+", in French", func(t *testing.T) {
			assert.Equal(t, tc.french, naming.French.ModeName(m, naming.Signs))
		})

		t.Run(tc.what+", in French words", func(t *testing.T) {
			assert.Equal(t, tc.words, naming.French.ModeName(m, naming.Words))
		})
	}

	t.Run("English words", func(t *testing.T) {
		m := modeAt(t, harmony.HarmonicMajor, 6)
		assert.Equal(t, "lydian sharp 2 sharp 5", naming.English.ModeName(m, naming.Words))
	})
}

// The perfect and imperfect split is the whole difficulty of the French
// register. A natural fourth is juste and never majeure; a natural
// sixth is majeure and never juste.
func TestPerfectAndImperfectQualitiesDoNotMix(t *testing.T) {
	render := naming.IntervalDegrees(naming.French)

	t.Run("perfect degrees are juste when natural", func(t *testing.T) {
		for _, d := range []harmony.Degree{1, 4, 5} {
			assert.Contains(t, render(d, naming.Natural), "juste",
				"degree %d is perfect", d)
		}
	})

	t.Run("imperfect degrees are majeure when natural", func(t *testing.T) {
		for _, d := range []harmony.Degree{2, 3, 6, 7} {
			assert.Contains(t, render(d, naming.Natural), "majeure",
				"degree %d is imperfect", d)
		}
	})

	t.Run("a flat is diminished on a perfect degree and minor elsewhere", func(t *testing.T) {
		assert.Contains(t, render(5, naming.Flat), "diminuée")
		assert.Contains(t, render(3, naming.Flat), "mineure")
	})

	t.Run("a sharp is augmented everywhere", func(t *testing.T) {
		for d := harmony.Degree(1); d <= 7; d++ {
			assert.Contains(t, render(d, naming.Sharp), "augmentée")
		}
	})

	t.Run("the outer qualities are one word, not a doubling", func(t *testing.T) {
		assert.Equal(t, "quinte sous-diminuée", render(5, naming.DoubleFlat),
			"a perfect degree lowered twice")
		assert.Equal(t, "tierce diminuée", render(3, naming.DoubleFlat),
			"an imperfect degree lowered twice reaches diminished, not lower")

		for _, q := range []string{
			naming.French.PerfectQualities[0],
			naming.French.PerfectQualities[4],
			naming.French.ImperfectQualities[4],
		} {
			assert.NotContains(t, q, "doublement")
		}
	})
}

// The natural sign is silent in a note name and spoken in a mode name.
// Sharing one table between the two would render phrygian natural 6 as
// bare phrygian, which is a different mode entirely.
func TestNaturalSignIsSilentInNotesAndSpokenInModes(t *testing.T) {
	t.Run("a natural note is named by its letter alone", func(t *testing.T) {
		for _, n := range []naming.Notation{naming.Signs, naming.Words} {
			assert.Equal(t, "F", naming.English.Name(natural(naming.LetterF), n))
			assert.Equal(t, "fa", naming.French.Name(natural(naming.LetterF), n))
		}
	})

	t.Run("a natural degree is spoken in a mode name", func(t *testing.T) {
		m := modeAt(t, harmony.MelodicMinor, 2)
		assert.Equal(t, "phrygian \u266e6", naming.English.ModeName(m, naming.Signs))
		assert.Equal(t, "phrygien bécarre 6", naming.French.ModeName(m, naming.Words))
		assert.Contains(t, naming.French.ModeAlternatives(m, naming.Signs), "phrygien sixte majeure",
			"the spoken register says the same thing in another form")
	})

	t.Run("dropping the sign would name a different mode", func(t *testing.T) {
		altered := modeAt(t, harmony.MelodicMinor, 2)
		bare := modeAt(t, harmony.NaturalMajor, 3)

		require.NotEqual(t, altered.Pattern(), bare.Pattern())
		assert.NotEqual(t,
			naming.English.ModeName(bare, naming.Signs),
			naming.English.ModeName(altered, naming.Signs),
			"phrygian and phrygian natural 6 are two modes and must read as two names",
		)
	})
}

// The thirty five names have to be as distinct as the thirty five
// patterns, or two modes become indistinguishable on screen even though
// they are not by ear.
func TestModeNamesAreUnique(t *testing.T) {
	for _, locale := range map[string]naming.Locale{
		"English": naming.English,
		"French":  naming.French,
	} {
		for _, notation := range []naming.Notation{naming.Signs, naming.Words} {
			seen := make(map[string]naming.Mode, 35)

			for _, m := range naming.Modes() {
				name := locale.ModeName(m, notation)
				previous, clash := seen[name]
				assert.False(t, clash,
					"name %q claimed by system %d degree %d and by system %d degree %d",
					name, previous.System, previous.Degree, m.System, m.Degree,
				)
				seen[name] = m
			}

			assert.Len(t, seen, 35)
		}
	}
}

func TestFunctionName(t *testing.T) {
	t.Run("the three functions have words in both locales", func(t *testing.T) {
		for _, f := range []harmony.Function{harmony.Tonic, harmony.Subdominant, harmony.Dominant} {
			assert.NotEmpty(t, naming.English.FunctionName(f))
			assert.NotEmpty(t, naming.French.FunctionName(f))
		}
	})

	t.Run("no function renders as nothing at all", func(t *testing.T) {
		assert.Empty(t, naming.English.FunctionName(harmony.NoFunction))
		assert.Empty(t, naming.French.FunctionName(harmony.NoFunction))
	})

	t.Run("the two modes without a tetrad render no function", func(t *testing.T) {
		for _, m := range naming.Modes() {
			if m.Chord != "" {
				continue
			}
			assert.Empty(t, naming.French.FunctionName(m.Function), m.String())
		}
	})
}
