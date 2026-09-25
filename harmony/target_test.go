package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func chord(t *testing.T, root harmony.PitchClass, p harmony.ChordPattern) harmony.Chord {
	t.Helper()
	c, err := harmony.NewChord(root, p)
	require.NoError(t, err)
	return c
}

func approach(t *testing.T, chords ...harmony.Chord) harmony.Approach {
	t.Helper()
	a, err := harmony.NewApproach(chords...)
	require.NoError(t, err)
	return a
}

func toC(t *testing.T) harmony.Target {
	t.Helper()
	key, err := harmony.NewTonality(0, harmony.ScaleMajor)
	require.NoError(t, err)
	return harmony.Target{Root: 0, Key: key}
}

// The substitution is not a special case. What makes a dominant pull is
// a tritone, and two chords a tritone apart carry the same one, so both
// are found by the same test.
func TestDominantIsFoundByItsTritone(t *testing.T) {
	target := toC(t)

	t.Run("the dominant a fifth above the target", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 7, harmony.ChordDominantSeventh),
			chord(t, 0, harmony.ChordMajorSeventh),
		))
		assert.True(t, got.Dominant)
		assert.Equal(t, harmony.PitchClass(7), got.DominantRoot)
		assert.False(t, got.Substitute)
		assert.True(t, got.Arrived)
	})

	t.Run("the tritone substitution, without a table listing it", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 1, harmony.ChordDominantSeventh),
			chord(t, 0, harmony.ChordMajorSeventh),
		))
		assert.True(t, got.Dominant)
		assert.Equal(t, harmony.PitchClass(1), got.DominantRoot)
		assert.True(t, got.Substitute)
		assert.True(t, got.Arrived)
	})

	t.Run("exactly two roots are dominants of any target", func(t *testing.T) {
		for root := harmony.PitchClass(0); root < harmony.PitchClassCount; root++ {
			got := target.Resolve(approach(t,
				chord(t, root, harmony.ChordDominantSeventh),
			))
			expected := root == 7 || root == 1
			assert.Equal(t, expected, got.Dominant,
				"a dominant seventh on %d, resolving to C", root)
		}
	})

	t.Run("a chord with no tritone pulls nowhere", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 7, harmony.ChordMajorTriad),
			chord(t, 0, harmony.ChordMajorTriad),
		))
		assert.False(t, got.Dominant,
			"a plain triad on the fifth has no seventh and so no tritone")
		assert.True(t, got.Arrived)
	})
}

// Two chords, one degree. This is why an approach is a sequence: an
// evaluation counting chords would see a spurious extra one here.
func TestSuspensionIsOneDegreeNotTwoChords(t *testing.T) {
	target := toC(t)

	t.Run("a suspension resolving onto its own third is recognised", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 7, harmony.ChordDominantSeventhSus4),
			chord(t, 7, harmony.ChordDominantSeventh),
			chord(t, 0, harmony.ChordMajorSeventh),
		))
		assert.Equal(t, 0, got.SuspensionAt)
		assert.Equal(t, 3, got.Steps)
		assert.True(t, got.Dominant)
		assert.True(t, got.Arrived)
	})

	t.Run("a suspension that never resolves is not one", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 7, harmony.ChordDominantSeventhSus4),
			chord(t, 0, harmony.ChordMajorSeventh),
		))
		assert.Equal(t, -1, got.SuspensionAt,
			"no later chord on the same root supplied the third it stood for")
	})

	t.Run("a chord with a third is never a suspension", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 7, harmony.ChordDominantSeventh),
			chord(t, 7, harmony.ChordDominantSeventh),
			chord(t, 0, harmony.ChordMajorTriad),
		))
		assert.Equal(t, -1, got.SuspensionAt)
	})
}

func TestArrival(t *testing.T) {
	target := toC(t)

	t.Run("landing on the target root arrives", func(t *testing.T) {
		got := target.Resolve(approach(t, chord(t, 0, harmony.ChordMajorTriad)))
		assert.True(t, got.Arrived)
	})

	t.Run("landing elsewhere does not", func(t *testing.T) {
		got := target.Resolve(approach(t, chord(t, 5, harmony.ChordMajorTriad)))
		assert.False(t, got.Arrived)
	})

	t.Run("a target naming a tetrad demands that tetrad", func(t *testing.T) {
		named := target
		named.Tetrad = harmony.ChordMajorSeventh

		assert.True(t, named.Resolve(approach(t,
			chord(t, 0, harmony.ChordMajorSeventh))).Arrived)
		assert.False(t, named.Resolve(approach(t,
			chord(t, 0, harmony.ChordMinorSeventh))).Arrived,
			"the root is right and the colour is not")
	})

	t.Run("a target naming no tetrad leaves the colour open", func(t *testing.T) {
		for _, p := range []harmony.ChordPattern{
			harmony.ChordMajorTriad,
			harmony.ChordMajorSeventh,
			harmony.ChordMajorSixth,
		} {
			assert.True(t, target.Resolve(approach(t, chord(t, 0, p))).Arrived)
		}
	})

	t.Run("only the last chord decides", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 0, harmony.ChordMajorTriad),
			chord(t, 5, harmony.ChordMajorTriad),
		))
		assert.False(t, got.Arrived,
			"passing through the target on the way is not arriving at it")
	})
}

// A colour that belongs and one that does not are different events, and
// the type keeps them apart so that a game cannot reward the pile-up.
func TestExtensionsAreRecordedThenJudged(t *testing.T) {
	target := toC(t)

	ninth, err := harmony.NewChordPattern(0, 4, 7, 10, 14)
	require.NoError(t, err)
	flatNinth, err := harmony.NewChordPattern(0, 4, 7, 10, 13)
	require.NoError(t, err)

	t.Run("an extension inside the key is recorded and kept", func(t *testing.T) {
		got := target.Resolve(approach(t, chord(t, 7, ninth)))
		assert.False(t, got.Extensions.IsEmpty())
		assert.Equal(t, got.Extensions, got.InKey,
			"the ninth of G is A, which is in C major")
	})

	t.Run("an extension outside the key is recorded and not kept", func(t *testing.T) {
		got := target.Resolve(approach(t, chord(t, 7, flatNinth)))
		assert.False(t, got.Extensions.IsEmpty())
		assert.True(t, got.InKey.IsEmpty(),
			"the minor ninth of G is A flat, foreign to C major")
	})

	t.Run("with no key nothing is kept, rather than everything", func(t *testing.T) {
		open := harmony.Target{Root: 0}
		got := open.Resolve(approach(t, chord(t, 7, ninth)))
		assert.False(t, got.Extensions.IsEmpty())
		assert.True(t, got.InKey.IsEmpty(),
			"an unjudged colour is not a validated one")
	})

	t.Run("a plain tetrad sounds no extension at all", func(t *testing.T) {
		got := target.Resolve(approach(t, chord(t, 7, harmony.ChordDominantSeventh)))
		assert.True(t, got.Extensions.IsEmpty())
	})
}

func TestMotionMeasuresConnection(t *testing.T) {
	target := toC(t)

	t.Run("a single chord moves nothing", func(t *testing.T) {
		got := target.Resolve(approach(t, chord(t, 0, harmony.ChordMajorTriad)))
		assert.Equal(t, harmony.Unison, got.Motion)
	})

	t.Run("repeating a chord moves nothing", func(t *testing.T) {
		got := target.Resolve(approach(t,
			chord(t, 0, harmony.ChordMajorTriad),
			chord(t, 0, harmony.ChordMajorTriad),
		))
		assert.Equal(t, harmony.Unison, got.Motion)
	})

	t.Run("chords sharing notes move less than chords sharing none", func(t *testing.T) {
		near := target.Resolve(approach(t,
			chord(t, 0, harmony.ChordMajorTriad),
			chord(t, 9, harmony.ChordMinorTriad),
		))
		far := target.Resolve(approach(t,
			chord(t, 0, harmony.ChordMajorTriad),
			chord(t, 6, harmony.ChordMajorTriad),
		))
		assert.Less(t, near.Motion, far.Motion,
			"C and A minor share two notes, C and F sharp share none")
	})
}

// Nothing here refuses. An approach that reaches nothing comes back
// with what was observed, because telling a player they never landed is
// worth more than telling them they scored badly.
func TestResolveNeverRefuses(t *testing.T) {
	target := toC(t)

	got := target.Resolve(approach(t,
		chord(t, 6, harmony.ChordDiminishedTriad),
		chord(t, 11, harmony.ChordAugmentedTriad),
	))

	assert.False(t, got.Arrived)
	assert.False(t, got.Dominant)
	assert.Equal(t, -1, got.SuspensionAt)
	assert.Equal(t, 2, got.Steps)
}

func TestNewApproachRejectsNothing(t *testing.T) {
	_, err := harmony.NewApproach()
	assert.Error(t, err, "an approach with no chord is not an approach")
}
