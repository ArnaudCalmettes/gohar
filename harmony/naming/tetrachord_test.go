package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/stretchr/testify/assert"
)

func TestTetrachordName(t *testing.T) {
	// The table of En Harmonie, tome 2, in its order.
	named := []struct {
		t       harmony.Tetrachord
		french  string
		english string
	}{
		{harmony.TetrachordMajor, "majeur", "major"},
		{harmony.TetrachordMinor, "mineur", "minor"},
		{harmony.TetrachordPhrygian, "phrygien", "phrygian"},
		{harmony.TetrachordLydian, "lydien", "lydian"},
		{harmony.TetrachordDiminished, "diminué", "diminished"},
		{harmony.TetrachordHarmonic, "harmonique", "harmonic"},
		{harmony.TetrachordLydianSharp2, "lydien ♯2", "lydian ♯2"},
		{harmony.TetrachordMinorSharp4, "mineur ♯4", "minor ♯4"},
		{harmony.TetrachordPhrygianDoubleFlat3, "phrygien 𝄫3", "phrygian 𝄫3"},
		{harmony.TetrachordMajorSharp2, "majeur ♯2", "major ♯2"},
	}

	for _, tc := range named {
		t.Run("the tetrachord "+tc.t.String()+" is "+tc.french, func(t *testing.T) {
			assert.Equal(t, tc.french, naming.French.TetrachordName(tc.t, naming.Signs))
			assert.Equal(t, tc.english, naming.English.TetrachordName(tc.t, naming.Signs))
		})
	}

	t.Run("the altered degree is said in words when asked", func(t *testing.T) {
		assert.Equal(t, "lydien dièse 2",
			naming.French.TetrachordName(harmony.TetrachordLydianSharp2, naming.Words))
		assert.Equal(t, "phrygien double bémol 3",
			naming.French.TetrachordName(harmony.TetrachordPhrygianDoubleFlat3, naming.Words))
		assert.Equal(t, "harmonique",
			naming.French.TetrachordName(harmony.TetrachordHarmonic, naming.Words),
			"a name without a degree has nothing to say in words")
	})

	// The rule the whole feature exists for: no name is invented.
	t.Run("a tetrachord the book does not name is designated by its steps", func(t *testing.T) {
		for _, steps := range []harmony.Tetrachord{0x132, 0x231, 0x123, 0x411, 0x141} {
			assert.Equal(t, steps.String(), naming.French.TetrachordName(steps, naming.Signs))
			assert.Equal(t, steps.String(), naming.English.TetrachordName(steps, naming.Signs),
				"a shape with no name in one language has none in any")
		}
	})

	t.Run("both kinds read naturally after the word for tetrachord", func(t *testing.T) {
		assert.Equal(t, "le tétracorde mineur ♯4",
			"le tétracorde "+naming.French.TetrachordName(harmony.TetrachordMinorSharp4, naming.Signs))
		assert.Equal(t, "le tétracorde 1 3 2",
			"le tétracorde "+naming.French.TetrachordName(0x132, naming.Signs))
	})
}
