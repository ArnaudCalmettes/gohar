package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
)

func TestTetrachordName(t *testing.T) {
	named := []struct {
		t       harmony.Tetrachord
		french  string
		english string
	}{
		{harmony.TetrachordMajor, "majeur", "major"},
		{harmony.TetrachordMinor, "mineur", "minor"},
		{harmony.TetrachordPhrygian, "phrygien", "phrygian"},
		{harmony.TetrachordLydian, "lydien", "lydian"},
		{harmony.TetrachordHarmonic, "harmonique", "harmonic"},
	}

	for _, tc := range named {
		t.Run("the tetrachord "+tc.t.String()+" has a name in both locales", func(t *testing.T) {
			assert.Equal(t, tc.french, naming.French.TetrachordName(tc.t))
			assert.Equal(t, tc.english, naming.English.TetrachordName(tc.t))
		})
	}

	// The rule the whole feature exists for: no name is invented.
	t.Run("an unnamed tetrachord is designated by its steps", func(t *testing.T) {
		for _, steps := range []harmony.Tetrachord{0x113, 0x311, 0x121, 0x312, 0x213} {
			assert.Equal(t, steps.String(), naming.French.TetrachordName(steps))
			assert.Equal(t, steps.String(), naming.English.TetrachordName(steps),
				"a shape with no name in one language has none in any")
		}
	})

	t.Run("both kinds read naturally after the word for tetrachord", func(t *testing.T) {
		assert.Equal(t, "le tétracorde harmonique",
			"le tétracorde "+naming.French.TetrachordName(harmony.TetrachordHarmonic))
		assert.Equal(t, "le tétracorde 1 1 3",
			"le tétracorde "+naming.French.TetrachordName(0x113))
	})
}
