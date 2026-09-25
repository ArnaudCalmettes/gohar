package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func offsetsOf(p harmony.ScalePattern) []harmony.Semitones {
	var out []harmony.Semitones
	for _, n := range p.Offsets() {
		out = append(out, n)
	}
	return out
}

func TestMotherScales(t *testing.T) {
	tests := []struct {
		name   string
		system harmony.System
		want   []harmony.Semitones
	}{
		{"natural major", harmony.NaturalMajor, []harmony.Semitones{0, 2, 4, 5, 7, 9, 11}},
		{"melodic minor", harmony.MelodicMinor, []harmony.Semitones{0, 2, 3, 5, 7, 9, 11}},
		{"harmonic minor", harmony.HarmonicMinor, []harmony.Semitones{0, 2, 3, 5, 7, 8, 11}},
		{"harmonic major", harmony.HarmonicMajor, []harmony.Semitones{0, 2, 4, 5, 7, 8, 11}},
		{"double harmonic major", harmony.DoubleHarmonicMajor, []harmony.Semitones{0, 1, 4, 5, 7, 8, 11}},
	}

	for _, tc := range tests {
		t.Run("the "+tc.name+" system has the expected offsets", func(t *testing.T) {
			assert.Equal(t, tc.want, offsetsOf(tc.system.Pattern()))
		})
	}
}

// The operation that matters most in daily use: from a pattern, find
// the mother scale and the degree. It has to work without any naming.
func TestModeOf(t *testing.T) {
	t.Run("every mode of every system is found at its own degree", func(t *testing.T) {
		for s := harmony.System(0); s < harmony.SystemCount; s++ {
			for d := harmony.Degree(1); d <= 7; d++ {
				p, ok := s.Mode(d)
				require.True(t, ok)

				gotS, gotD, ok := harmony.ModeOf(p)
				require.True(t, ok, "system %d degree %d", s, d)
				assert.Equal(t, s, gotS)
				assert.Equal(t, d, gotD)
			}
		}
	})

	t.Run("F lydian is the fourth degree of the natural system", func(t *testing.T) {
		lydian, ok := harmony.ScaleMajor.Mode(4)
		require.True(t, ok)

		s, d, ok := harmony.ModeOf(lydian)
		require.True(t, ok)
		assert.Equal(t, harmony.NaturalMajor, s)
		assert.Equal(t, harmony.Degree(4), d)
	})

	t.Run("the natural minor is found as a mode, not as a system", func(t *testing.T) {
		s, d, ok := harmony.ModeOf(harmony.ScaleNaturalMinor)
		require.True(t, ok)
		assert.Equal(t, harmony.NaturalMajor, s)
		assert.Equal(t, harmony.Degree(6), d, "it is the aeolian")
	})

	t.Run("a pattern outside the five systems is not found", func(t *testing.T) {
		wholeTone, err := harmony.NewScalePattern(
			harmony.PitchSet(0b010101_010101), 0)
		require.NoError(t, err)

		_, _, ok := harmony.ModeOf(wholeTone)
		assert.False(t, ok, "a whole tone scale has no mother scale")
	})

	t.Run("a value outside the five systems has no pattern", func(t *testing.T) {
		assert.Zero(t, harmony.System(harmony.SystemCount).Pattern())
	})
}
