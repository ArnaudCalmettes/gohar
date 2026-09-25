package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
)

func TestFunctionOf(t *testing.T) {
	t.Run("a dominant seventh is a dominant", func(t *testing.T) {
		assert.Equal(t, harmony.Dominant, harmony.FunctionOf(harmony.ChordDominantSeventh))
	})

	t.Run("a plain major triad can play all three roles", func(t *testing.T) {
		f := harmony.FunctionOf(harmony.ChordMajorTriad)
		assert.True(t, f.Has(harmony.Tonic))
		assert.True(t, f.Has(harmony.Subdominant))
		assert.True(t, f.Has(harmony.Dominant))
	})

	t.Run("a half diminished is a subdominant, as the locrian is", func(t *testing.T) {
		assert.Equal(t, harmony.Subdominant, harmony.FunctionOf(harmony.ChordHalfDiminished))
	})

	// The one that matters most, because it is the one a game would get
	// wrong by reflex. Removing the tritone is what made it plagal.
	t.Run("a sus chord is never a dominant", func(t *testing.T) {
		for _, p := range []harmony.ChordPattern{
			harmony.ChordSus2,
			harmony.ChordSus4,
			harmony.ChordDominantSeventhSus2,
			harmony.ChordDominantSeventhSus4,
		} {
			f := harmony.FunctionOf(p)
			assert.False(t, f.Has(harmony.Dominant), "%v", p)
			assert.True(t, f.Has(harmony.Subdominant), "%v", p)
		}
	})

	t.Run("every recognised tetrad has at least one role", func(t *testing.T) {
		for _, p := range []harmony.ChordPattern{
			harmony.ChordMajorTriad, harmony.ChordMinorTriad,
			harmony.ChordDiminishedTriad, harmony.ChordAugmentedTriad,
			harmony.ChordSus2, harmony.ChordSus4,
			harmony.ChordMajorSixth, harmony.ChordMinorSixth,
			harmony.ChordMajorSeventh, harmony.ChordMajorSeventhNo5,
			harmony.ChordMajorSeventhSharp5,
			harmony.ChordDominantSeventh, harmony.ChordDominantSeventhNo5,
			harmony.ChordDominantSeventhFlat5, harmony.ChordDominantSeventhSharp5,
			harmony.ChordMinorSeventh, harmony.ChordMinorSeventhNo5,
			harmony.ChordHalfDiminished,
			harmony.ChordMinorMajorSeventh, harmony.ChordMinorMajorSeventhNo5,
			harmony.ChordDiminishedSeventh,
			harmony.ChordDominantSeventhSus2, harmony.ChordDominantSeventhSus4,
		} {
			assert.NotEqual(t, harmony.NoFunction, harmony.FunctionOf(p), "%v", p)
		}
	})

	t.Run("an unknown shape has none", func(t *testing.T) {
		assert.Equal(t, harmony.NoFunction, harmony.FunctionOf(harmony.ChordPattern(1)))
	})
}
