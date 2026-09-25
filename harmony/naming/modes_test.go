package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests exist to catch transcription errors, not to discover
// theory. Every property below was given as true by someone who studied
// the system, so a failure means the table was typed wrong.
//
// That is the point. A wrong offset produces a pattern that stays
// unique, so uniqueness alone would not catch it; the properties that
// tie a mode to its name and to its mother scale would.

func TestCatalogueHasThirtyFiveModes(t *testing.T) {
	assert.Len(t, naming.Modes(), 35)
}

// Thirty five distinct patterns. Two modes sharing one would be
// indistinguishable by ear, which the nomenclature does not allow.
func TestModePatternsAreUnique(t *testing.T) {
	seen := make(map[harmony.ScalePattern]naming.Mode, 35)

	for _, m := range naming.Modes() {
		p := m.Pattern()
		if previous, clash := seen[p]; clash {
			t.Errorf(
				"pattern %s claimed by system %d degree %d and by system %d degree %d",
				p, previous.System, previous.Degree, m.System, m.Degree,
			)
			continue
		}
		seen[p] = m
	}

	assert.Len(t, seen, 35)
}

// Every mode is a rotation of its own mother scale. This is the
// strongest automatic check available: it fails on any offset error
// that is not simultaneously a coherent error on the mother scale.
func TestEveryModeRotatesItsOwnSystem(t *testing.T) {
	for _, m := range naming.Modes() {
		t.Run(m.String()+" is a rotation of its mother scale", func(t *testing.T) {
			want, ok := m.System.Mode(m.Degree)
			require.True(t, ok, "degree %d of the mother scale", m.Degree)
			assert.Equal(t, want, m.Pattern())
		})
	}
}

// The name reproduces the pattern: start from the base natural mode,
// then set each named degree to the major scale offset for that degree
// adjusted by its quality.
//
// The quality is absolute rather than a movement from the base mode.
// Four modes distinguish the two readings, all carrying a double flat,
// and they are the reason this test exists in this form.
func TestNameReproducesPattern(t *testing.T) {
	major := collectOffsets(t, harmony.NaturalMajor.Pattern())

	for _, m := range naming.Modes() {
		t.Run(m.String()+" is exactly what its name describes", func(t *testing.T) {
			base := naming.Mode{System: harmony.NaturalMajor, Base: m.Base}
			offsets := collectOffsets(t, base.Pattern())

			for _, alt := range m.Altered {
				offsets[alt.Degree-1] = major[alt.Degree-1] + harmony.Semitones(alt.Quality)
			}

			assert.Equal(t, collectOffsets(t, m.Pattern()), offsets)
		})
	}
}

// The four modes where an absolute reading and a relative one diverge.
// Pinned explicitly so that a future refactor of Pattern cannot quietly
// switch back to subtracting from the base mode.
func TestDoubleFlatIsAbsoluteNotRelative(t *testing.T) {
	tests := []struct {
		name   string
		system harmony.System
		degree harmony.Degree
		want   []harmony.Semitones
	}{
		{
			name:   "locrian double flat 7 of the harmonic major puts its seventh at nine",
			system: harmony.HarmonicMajor, degree: 7,
			want: []harmony.Semitones{0, 1, 3, 5, 6, 8, 9},
		},
		{
			name:   "locrian flat 4 double flat 7 of the harmonic minor puts its seventh at nine",
			system: harmony.HarmonicMinor, degree: 7,
			want: []harmony.Semitones{0, 1, 3, 4, 6, 8, 9},
		},
		{
			name:   "phrygian flat 4 double flat 7 of the double harmonic puts its seventh at nine",
			system: harmony.DoubleHarmonicMajor, degree: 3,
			want: []harmony.Semitones{0, 1, 3, 4, 7, 8, 9},
		},
		{
			name:   "locrian double flat 3 double flat 7 puts its third at two and its seventh at nine",
			system: harmony.DoubleHarmonicMajor, degree: 7,
			want: []harmony.Semitones{0, 1, 2, 5, 6, 8, 9},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, ok := naming.Lookup(tc.system, tc.degree)
			require.True(t, ok)
			assert.Equal(t, tc.want, collectOffsets(t, m.Pattern()))
		})
	}
}

// Every mode is heptatonic, which is what makes a Tonality out of it
// and what lets it be spelled one letter per degree.
func TestEveryModeIsHeptatonic(t *testing.T) {
	for _, m := range naming.Modes() {
		t.Run(m.String()+" has seven notes", func(t *testing.T) {
			assert.True(t, m.Pattern().IsHeptatonic())
		})
	}
}

func TestCharacteristicDegrees(t *testing.T) {
	t.Run("the seven natural modes carry no altered degree", func(t *testing.T) {
		for _, m := range naming.Modes() {
			if m.System != harmony.NaturalMajor {
				continue
			}
			assert.Empty(t, m.AlteredDegrees(), m.String())
		}
	})

	t.Run("the twenty eight altered modes each carry at least one", func(t *testing.T) {
		count := 0
		for _, m := range naming.Modes() {
			if m.System == harmony.NaturalMajor {
				continue
			}
			count++
			assert.NotEmpty(t, m.AlteredDegrees(), m.String())
		}
		assert.Equal(t, 28, count)
	})

	t.Run("every mode inherits the natural degrees of its base mode", func(t *testing.T) {
		for _, m := range naming.Modes() {
			base, ok := naming.Lookup(harmony.NaturalMajor, harmony.Degree(m.Base)+1)
			require.True(t, ok, m.String())
			assert.Equal(t, base.NaturalDegrees(), m.NaturalDegrees(), m.String())
		}
	})

	t.Run("the ionian carries a natural fourth and a natural seventh", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.NaturalMajor, 1)
		require.True(t, ok)
		assert.Equal(t,
			[]naming.Alteration{{Degree: 4, Quality: naming.Natural},
				{Degree: 7, Quality: naming.Natural}},
			m.NaturalDegrees(),
			"the tritone between them is what pulls toward the tonic",
		)
	})

	t.Run("the lydian sharp 2 keeps the sharp fourth and adds the sharp second", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.HarmonicMinor, 6)
		require.True(t, ok)
		assert.Equal(t,
			[]naming.Alteration{{Degree: 4, Quality: naming.Sharp}},
			m.NaturalDegrees(),
		)
		assert.Equal(t,
			[]naming.Alteration{{Degree: 2, Quality: naming.Sharp}},
			m.AlteredDegrees(),
		)
	})
}

// Two modes voice no tetrad, and one of them still carries extensions
// with no chord to hang them on.
func TestModesWithoutTetrad(t *testing.T) {
	t.Run("exactly two modes voice no tetrad", func(t *testing.T) {
		var without []naming.Mode
		for _, m := range naming.Modes() {
			if m.Chord == "" {
				without = append(without, m)
			}
		}
		require.Len(t, without, 2)
		for _, m := range without {
			assert.Equal(t, harmony.DoubleHarmonicMajor, m.System)
			assert.Equal(t, harmony.NoFunction, m.Function)
		}
	})
}

// There are three functions and no more. A mode voicing a minor tetrad
// in a tonic role carries Tonic: the minor lives in the chord symbol,
// and a separate value for it would split one function into two that
// behave alike.
//
// Closed is not the same as exclusive. The first degrees of harmonic
// minor and of harmonic major carry Tonic and Dominant at once, because
// they are also read as a dominant over a tonic pedal. What the
// catalogue may not do is invent a fourth role.
func TestFunctionsAreClosed(t *testing.T) {
	t.Run("no mode carries a role outside the three", func(t *testing.T) {
		known := harmony.Tonic | harmony.Subdominant | harmony.Dominant
		for _, m := range naming.Modes() {
			assert.Zero(t, m.Function&^known, m.String())
		}
	})

	t.Run("only the two modes without a tetrad carry no function", func(t *testing.T) {
		for _, m := range naming.Modes() {
			if m.Function == harmony.NoFunction {
				assert.Empty(t, m.Chord, m.String())
			}
		}
	})

	t.Run("the dorian and the aeolian are tonic, minor chord notwithstanding", func(t *testing.T) {
		for _, degree := range []harmony.Degree{2, 6} {
			m, ok := naming.Lookup(harmony.NaturalMajor, degree)
			require.True(t, ok)
			assert.Equal(t, harmony.Tonic, m.Function)
			assert.Contains(t, m.Chord, "Xm7")
		}
	})
}
