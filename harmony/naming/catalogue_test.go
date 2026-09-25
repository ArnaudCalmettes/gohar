package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The two modes that carry a V7 flat 9 on I alongside their tonic role.
// The double role used to live in a prose note that was removed with
// the Usage field; the bit field now holds it as data.
func TestDoubleFunctionModes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		system harmony.System
	}{
		{"the ionian flat 6 of the harmonic major", harmony.HarmonicMajor},
		{"the aeolian natural 7 of the harmonic minor", harmony.HarmonicMinor},
	} {
		t.Run(tc.name+" is both tonic and dominant", func(t *testing.T) {
			m, ok := naming.Lookup(tc.system, 1)
			require.True(t, ok)
			assert.True(t, m.Function.Has(harmony.Tonic))
			assert.True(t, m.Function.Has(harmony.Dominant))
		})
	}
}

// The mirror is a structural check on the transcription, not a naming
// rule: a mistyped offset would very likely send a mode outside the
// catalogue.
func TestCatalogueIsClosedUnderMirror(t *testing.T) {
	patterns := make(map[harmony.ScalePattern]naming.Mode, 35)
	for _, m := range naming.Modes() {
		patterns[m.Pattern()] = m
	}

	t.Run("every mirror lands on a mode of the catalogue", func(t *testing.T) {
		for _, m := range naming.Modes() {
			_, ok := patterns[m.Pattern().Mirror()]
			assert.True(t, ok, "the mirror of %s is not a known mode", m)
		}
	})

	t.Run("the mirror is its own inverse", func(t *testing.T) {
		for _, m := range naming.Modes() {
			assert.Equal(t, m.Pattern(), m.Pattern().Mirror().Mirror(), m.String())
		}
	})

	t.Run("exactly three modes are their own mirror", func(t *testing.T) {
		var palindromes []string
		for _, m := range naming.Modes() {
			if m.Pattern().Mirror() == m.Pattern() {
				palindromes = append(palindromes, m.String())
			}
		}
		assert.ElementsMatch(t,
			[]string{"dorian", "mixolydian b6", "ionian b2 b6"},
			palindromes,
			"a mode that mirrors onto itself is one whose inverse sounds like it",
		)
	})

	t.Run("the harmonic minor and harmonic major mirror each other", func(t *testing.T) {
		for _, m := range naming.Modes() {
			mirrored := patterns[m.Pattern().Mirror()]
			switch m.System {
			case harmony.HarmonicMinor:
				assert.Equal(t, harmony.HarmonicMajor, mirrored.System, m.String())
			case harmony.HarmonicMajor:
				assert.Equal(t, harmony.HarmonicMinor, mirrored.System, m.String())
			default:
				assert.Equal(t, m.System, mirrored.System,
					"%s should mirror within its own system", m)
			}
		}
	})

	t.Run("the ionian mirrors onto the phrygian", func(t *testing.T) {
		ionian, ok := naming.Lookup(harmony.NaturalMajor, 1)
		require.True(t, ok)
		phrygian, ok := naming.Lookup(harmony.NaturalMajor, 3)
		require.True(t, ok)
		assert.Equal(t, phrygian.Pattern(), ionian.Pattern().Mirror())
	})
}
