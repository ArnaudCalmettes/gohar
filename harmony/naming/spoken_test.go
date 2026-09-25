package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The twenty eight altered modes as they are said aloud in French.
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
		{harmony.HarmonicMinor, 7, "locrien ♭4 ♭♭7"},
		{harmony.HarmonicMajor, 1, "ionien sixte mineure"},
		{harmony.HarmonicMajor, 2, "dorien ♭5"},
		{harmony.HarmonicMajor, 3, "phrygien ♭4"},
		{harmony.HarmonicMajor, 4, "lydien mineur"},
		{harmony.HarmonicMajor, 5, "mixolydien seconde mineure"},
		{harmony.HarmonicMajor, 6, "lydien ♯2 ♯5"},
		{harmony.HarmonicMajor, 7, "locrien ♭♭7"},
		{harmony.DoubleHarmonicMajor, 1, "ionien ♭2 ♭6"},
		{harmony.DoubleHarmonicMajor, 2, "lydien ♯2 ♯6"},
		{harmony.DoubleHarmonicMajor, 3, "phrygien ♭4 ♭♭7"},
		{harmony.DoubleHarmonicMajor, 4, "éolien ♯4 ♮7"},
		{harmony.DoubleHarmonicMajor, 5, "mixolydien ♭2 ♭5"},
		{harmony.DoubleHarmonicMajor, 6, "ionien ♯2 ♯5"},
		{harmony.DoubleHarmonicMajor, 7, "locrien ♭♭3 ♭♭7"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			m, ok := naming.Lookup(tc.system, tc.degree)
			require.True(t, ok)
			assert.Equal(t, tc.want, naming.French.ModeName(m))
		})
	}
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
