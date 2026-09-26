package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The twenty eight altered modes as they are said aloud in French, the
// refined register offered as an alternative to the systematic name.
//
// This table was read and approved line by line before being written
// down, so it is the reference and the procedure is what has to match
// it, not the other way round.
func TestFrenchSpokenModeNames(t *testing.T) {
	tests := []struct {
		system harmony.System
		degree harmony.Degree
		want   string
	}{
		{harmony.MelodicMinor, 1, "ionien mineur"},
		{harmony.MelodicMinor, 2, "phrygien sixte majeure"},
		{harmony.MelodicMinor, 3, "lydien augmenté"},
		{harmony.MelodicMinor, 4, "lydien septième mineure"},
		{harmony.MelodicMinor, 5, "mixolydien sixte mineure"},
		{harmony.MelodicMinor, 6, "locrien seconde majeure"},
		{harmony.MelodicMinor, 7, "locrien ♭4"},
		{harmony.HarmonicMinor, 1, "éolien septième majeure"},
		{harmony.HarmonicMinor, 2, "locrien sixte majeure"},
		{harmony.HarmonicMinor, 3, "ionien augmenté"},
		{harmony.HarmonicMinor, 4, "dorien quarte augmentée"},
		{harmony.HarmonicMinor, 5, "phrygien majeur"},
		{harmony.HarmonicMinor, 6, "lydien seconde augmentée"},
		{harmony.HarmonicMinor, 7, "locrien ♭4 𝄫7"},
		{harmony.HarmonicMajor, 1, "ionien sixte mineure"},
		{harmony.HarmonicMajor, 2, "dorien ♭5"},
		{harmony.HarmonicMajor, 3, "phrygien ♭4"},
		{harmony.HarmonicMajor, 4, "lydien mineur"},
		{harmony.HarmonicMajor, 5, "mixolydien seconde mineure"},
		{harmony.HarmonicMajor, 6, "lydien ♯2 ♯5"},
		{harmony.HarmonicMajor, 7, "locrien 𝄫7"},
		{harmony.DoubleHarmonicMajor, 1, "ionien ♭2 ♭6"},
		{harmony.DoubleHarmonicMajor, 2, "lydien ♯2 ♯6"},
		{harmony.DoubleHarmonicMajor, 3, "phrygien ♭4 𝄫7"},
		{harmony.DoubleHarmonicMajor, 4, "éolien ♯4 ♮7"},
		{harmony.DoubleHarmonicMajor, 5, "mixolydien ♭2 ♭5"},
		{harmony.DoubleHarmonicMajor, 6, "ionien ♯2 ♯5"},
		{harmony.DoubleHarmonicMajor, 7, "locrien 𝄫3 𝄫7"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			m, ok := naming.Lookup(tc.system, tc.degree)
			require.True(t, ok)
			got, ok := naming.French.SpokenModeName(m, naming.Signs)
			require.True(t, ok)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("its signs follow the notation", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.HarmonicMinor, 7)
		require.True(t, ok)
		got, _ := naming.French.SpokenModeName(m, naming.Words)
		assert.Equal(t, "locrien bémol 4 double bémol 7", got)
	})

	t.Run("English has no spoken register", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.MelodicMinor, 2)
		require.True(t, ok)
		_, ok = naming.English.SpokenModeName(m, naming.Signs)
		assert.False(t, ok)
	})
}

// The alternatives are the spoken name, when it differs, then the
// aliases: the phrygian dominant is found under all three.
func TestModeAlternatives(t *testing.T) {
	m, ok := naming.Lookup(harmony.HarmonicMinor, 5)
	require.True(t, ok)
	assert.Equal(t, "phrygien \u266e3", naming.French.ModeName(m, naming.Signs))
	assert.Equal(t,
		[]string{"phrygien majeur", "phrygien dominante"},
		naming.French.ModeAlternatives(m, naming.Signs))

	t.Run("a spoken name equal to the systematic one is not repeated", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.HarmonicMajor, 6)
		require.True(t, ok)
		assert.Empty(t, naming.French.ModeAlternatives(m, naming.Signs),
			"lydien ♯2 ♯5 is said the way it is written")
	})

	t.Run("a natural mode has none", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.NaturalMajor, 2)
		require.True(t, ok)
		assert.Empty(t, naming.French.ModeAlternatives(m, naming.Signs))
	})
}

// Aliases are written out, not derived, and they belong to a language.
func TestModeAliases(t *testing.T) {
	tests := []struct {
		system  harmony.System
		degree  harmony.Degree
		french  string
		english string
	}{
		{harmony.MelodicMinor, 7, "altéré", "altered"},
		{harmony.MelodicMinor, 4, "lydien dominante", "lydian dominant"},
		{harmony.HarmonicMinor, 5, "phrygien dominante", "phrygian dominant"},
	}

	for _, tc := range tests {
		t.Run(tc.french, func(t *testing.T) {
			m, ok := naming.Lookup(tc.system, tc.degree)
			require.True(t, ok)
			assert.Contains(t, naming.French.ModeAliases(m), tc.french)
			assert.Contains(t, naming.English.ModeAliases(m), tc.english)
		})
	}

	t.Run("a mode with no alias has none", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.NaturalMajor, 1)
		require.True(t, ok)
		assert.Empty(t, naming.French.ModeAliases(m))
	})

	t.Run("the phrygian flat 4 is not named after the function it has", func(t *testing.T) {
		m, ok := naming.Lookup(harmony.HarmonicMajor, 3)
		require.True(t, ok)
		assert.Empty(t, naming.French.ModeAliases(m),
			"its flat 4 is heard as a major third, but the name stays rigorous")
	})
}
