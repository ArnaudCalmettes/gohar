package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func steps(t *testing.T, s ...harmony.Step) harmony.Progression {
	t.Helper()
	p, err := harmony.NewProgressionSteps(s...)
	require.NoError(t, err)
	return p
}

// The premise: a player without perfect pitch transposes. A progression
// is therefore a shape, and the key is a fact about one performance.
func TestCompareIgnoresTranspositionWithoutLosingIt(t *testing.T) {
	expected := harmony.TwoFiveOneMajor

	t.Run("the same shape a tone higher is the same shape", func(t *testing.T) {
		played := steps(t,
			harmony.Step{Offset: 0, Tetrad: harmony.ChordMinorSeventh},
			harmony.Step{Offset: 5, Tetrad: harmony.ChordDominantSeventh},
			harmony.Step{Offset: 10, Tetrad: harmony.ChordMajorSeventh},
		)

		got := expected.Compare(played)
		assert.True(t, got.Same)
		assert.Equal(t, -1, got.FirstDivergence)
	})

	t.Run("the shift is reported even when the shape matches", func(t *testing.T) {
		inD, err := harmony.NewProgression(
			[]harmony.ChordPattern{
				harmony.ChordMinorSeventh,
				harmony.ChordDominantSeventh,
				harmony.ChordMajorSeventh,
			},
			[]harmony.PitchClass{4, 9, 2},
		)
		require.NoError(t, err)

		inC, err := harmony.NewProgression(
			[]harmony.ChordPattern{
				harmony.ChordMinorSeventh,
				harmony.ChordDominantSeventh,
				harmony.ChordMajorSeventh,
			},
			[]harmony.PitchClass{2, 7, 0},
		)
		require.NoError(t, err)

		got := inC.Compare(inD)
		assert.True(t, got.Same, "a two five one is a two five one")
		assert.True(t, got.Transposed)
		assert.Equal(t, harmony.Semitones(2), got.Shift,
			"ignoring the key is the caller's decision, so the number has to be there",
		)
	})

	t.Run("the same key reports no shift", func(t *testing.T) {
		got := expected.Compare(expected)
		assert.True(t, got.Same)
		assert.False(t, got.Transposed)
		assert.Equal(t, harmony.Unison, got.Shift)
	})
}

// One shift for the whole sequence. Resolving it per chord would accept
// a performance whose chords were each some transposition of the right
// one while the relations between them were wrong, which is exactly the
// mistake worth catching.
func TestOneShiftForTheWholeSequence(t *testing.T) {
	expected := harmony.TwoFiveOneMajor

	t.Run("a chord transposed on its own breaks the shape", func(t *testing.T) {
		played := steps(t,
			harmony.Step{Offset: 0, Tetrad: harmony.ChordMinorSeventh},
			harmony.Step{Offset: 7, Tetrad: harmony.ChordDominantSeventh},
			harmony.Step{Offset: 10, Tetrad: harmony.ChordMajorSeventh},
		)

		got := expected.Compare(played)
		assert.False(t, got.Same,
			"each chord is right, the motion between them is not")
		assert.Equal(t, 1, got.FirstDivergence)
	})

	t.Run("the divergence is located, not counted", func(t *testing.T) {
		played := steps(t,
			harmony.Step{Offset: 0, Tetrad: harmony.ChordMinorSeventh},
			harmony.Step{Offset: 5, Tetrad: harmony.ChordDominantSeventh},
			harmony.Step{Offset: 10, Tetrad: harmony.ChordMinorSeventh},
		)

		got := expected.Compare(played)
		assert.False(t, got.Same)
		assert.Equal(t, 2, got.FirstDivergence,
			"the first two chords were right, and that is what to tell the player")
	})
}

// Extensions colour a chord without changing which chord it is.
func TestExtensionsDoNotChangeTheShape(t *testing.T) {
	expected := harmony.TwoFiveOneMajor

	ninth, err := harmony.NewChordPattern(0, 4, 7, 10, 14)
	require.NoError(t, err)

	played := steps(t,
		harmony.Step{Offset: 0, Tetrad: harmony.ChordMinorSeventh},
		harmony.Step{Offset: 5, Tetrad: ninth.Tetrad()},
		harmony.Step{Offset: 10, Tetrad: harmony.ChordMajorSeventh},
	)

	got := expected.Compare(played)
	assert.True(t, got.Same,
		"voicing a ninth on the dominant is still playing the two five one")
}

func TestDifferentLengthsNeverMatch(t *testing.T) {
	expected := harmony.TwoFiveOneMajor

	short := steps(t,
		harmony.Step{Offset: 0, Tetrad: harmony.ChordMinorSeventh},
		harmony.Step{Offset: 5, Tetrad: harmony.ChordDominantSeventh},
	)

	got := expected.Compare(short)
	assert.False(t, got.Same)
	assert.Equal(t, 2, got.FirstDivergence,
		"the divergence is the first index past the shorter one")
}

// The degree is what a player is told, and it needs a tonality.
func TestDegreesNameWhatTheAbsoluteRootCannot(t *testing.T) {
	cMajor, err := harmony.NewTonality(0, harmony.ScaleMajor)
	require.NoError(t, err)

	t.Run("a two five one in C names its degrees two, five and one", func(t *testing.T) {
		assert.Equal(t,
			[]harmony.Degree{2, 5, 1},
			harmony.TwoFiveOneMajor.Degrees(2, cMajor),
		)
	})

	t.Run("the same shape in another key names the same degrees", func(t *testing.T) {
		eMajor, err := harmony.NewTonality(4, harmony.ScaleMajor)
		require.NoError(t, err)

		assert.Equal(t,
			harmony.TwoFiveOneMajor.Degrees(2, cMajor),
			harmony.TwoFiveOneMajor.Degrees(6, eMajor),
			"the whole point: the feedback survives transposition",
		)
	})

	t.Run("a root outside the key has no degree", func(t *testing.T) {
		got := harmony.TwoFiveOneMajor.Degrees(1, cMajor)
		require.Len(t, got, 3)
		assert.Equal(t, harmony.Degree(0), got[0],
			"zero means no degree, and degrees count from one")
	})
}

func TestProgressionAnchoring(t *testing.T) {
	t.Run("anchoring yields the chords of one performance", func(t *testing.T) {
		var roots []harmony.PitchClass
		var tetrads []harmony.ChordPattern
		for chord := range harmony.TwoFiveOneMajor.At(2) {
			roots = append(roots, chord.Root)
			tetrads = append(tetrads, chord.Pattern)
		}

		assert.Equal(t, []harmony.PitchClass{2, 7, 0}, roots)
		assert.Equal(t,
			[]harmony.ChordPattern{
				harmony.ChordMinorSeventh,
				harmony.ChordDominantSeventh,
				harmony.ChordMajorSeventh,
			},
			tetrads,
		)
	})

	t.Run("anchoring on every root gives the same shape back", func(t *testing.T) {
		for root := harmony.PitchClass(0); root < harmony.PitchClassCount; root++ {
			var tetrads []harmony.ChordPattern
			var roots []harmony.PitchClass
			for chord := range harmony.TwoFiveOneMajor.At(root) {
				tetrads = append(tetrads, chord.Pattern)
				roots = append(roots, chord.Root)
			}

			back, err := harmony.NewProgression(tetrads, roots)
			require.NoError(t, err)
			assert.True(t, harmony.TwoFiveOneMajor.Compare(back).Same, "root %d", root)
		}
	})

	t.Run("a progression has no Transpose, since it is already relative", func(t *testing.T) {
		assert.Equal(t, harmony.Semitones(0), harmony.TwoFiveOneMajor.Compare(
			harmony.TwoFiveOneMajor).Shift)
	})
}

func TestNewProgressionStepsRejectsAnOriginItLacks(t *testing.T) {
	_, err := harmony.NewProgressionSteps(
		harmony.Step{Offset: 3, Tetrad: harmony.ChordMinorSeventh},
		harmony.Step{Offset: 8, Tetrad: harmony.ChordDominantSeventh},
	)
	assert.Error(t, err,
		"the first step is the origin, so it sits at zero by construction")
}
