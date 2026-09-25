package analysis_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func held(pitches ...harmony.Pitch) analysis.Snapshot {
	return analysis.Snapshot{Pitches: pitches}
}

func newRecognizer(t *testing.T) *analysis.Recognizer {
	t.Helper()
	r, err := analysis.NewRecognizer(analysis.CommonTetrads...)
	require.NoError(t, err)
	return r
}

func noContext() analysis.Context { return analysis.Context{} }

func TestIdentifiesTheObviousChords(t *testing.T) {
	r := newRecognizer(t)

	tests := []struct {
		name   string
		held   analysis.Snapshot
		root   harmony.PitchClass
		tetrad harmony.ChordPattern
	}{
		{
			name: "a major triad in root position",
			held: held(60, 64, 67),
			root: 0, tetrad: harmony.ChordMajorTriad,
		},
		{
			name: "a dominant seventh in root position",
			held: held(55, 59, 62, 65),
			root: 7, tetrad: harmony.ChordDominantSeventh,
		},
		{
			name: "a major seventh is not taken for a dominant",
			held: held(60, 64, 67, 71),
			root: 0, tetrad: harmony.ChordMajorSeventh,
		},
		{
			name: "a fifthless voicing is a chord in its own right",
			held: held(60, 64, 70),
			root: 0, tetrad: harmony.ChordDominantSeventhNo5,
		},
		{
			name: "a sus4 is not a triad with a displaced third",
			held: held(60, 65, 67),
			root: 0, tetrad: harmony.ChordSus4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := r.Best(tc.held, noContext())
			require.True(t, ok)
			assert.Equal(t, tc.root, got.Root)
			assert.Equal(t, tc.tetrad, got.Tetrad)
			assert.True(t, got.BassIsRoot)
		})
	}
}

// The heart of the redesign. Whether a note is a second or a ninth is
// determined by the construction rules, before anything is compared.
// There is nothing here for a weighting to arbitrate.
func TestConstructionRulesDetermineTheReading(t *testing.T) {
	r := newRecognizer(t)

	t.Run("a D above a C major triad is a ninth, not a second", func(t *testing.T) {
		got, ok := r.Best(held(60, 62, 64, 67), noContext())
		require.True(t, ok)

		assert.Equal(t, harmony.ChordMajorTriad, got.Tetrad,
			"the tetrad is still the triad; the D colours it")
		assert.True(t, got.Extensions.Has(harmony.IntMajorNinth))
		assert.False(t, got.Pattern.Has(harmony.IntMajorSecond),
			"the second was pushed up, so it no longer sits at offset two")
	})

	t.Run("a D with no third stays a second, making a sus2", func(t *testing.T) {
		got, ok := r.Best(held(60, 62, 67), noContext())
		require.True(t, ok)
		assert.Equal(t, harmony.ChordSus2, got.Tetrad)
	})

	t.Run("an F above a C major triad is an eleventh, not a fourth", func(t *testing.T) {
		got, ok := r.Best(held(60, 64, 65, 67), noContext())
		require.True(t, ok)
		assert.Equal(t, harmony.ChordMajorTriad, got.Tetrad)
		assert.True(t, got.Extensions.Has(harmony.IntPerfectEleventh))
	})

	t.Run("no chord holds a minor second", func(t *testing.T) {
		got, ok := r.Best(held(60, 61, 64, 67), noContext())
		require.True(t, ok)
		assert.True(t, got.Extensions.Has(harmony.IntMinorNinth))
		assert.False(t, got.Pattern.Has(harmony.IntMinorSecond))
	})

	t.Run("a chord holding both thirds has the minor one as a raised ninth", func(t *testing.T) {
		got, ok := r.Best(held(60, 63, 64, 70), noContext())
		require.True(t, ok)
		assert.True(t, got.Pattern.Has(harmony.IntMajorThird),
			"the diminished fourth is heard as the major third of the tetrad")
		assert.True(t, got.Extensions.HasAny(harmony.IntAugmentedNinth, harmony.IntMinorTenth))
	})

	t.Run("a sixth alongside a seventh becomes a thirteenth", func(t *testing.T) {
		got, ok := r.Best(held(60, 64, 67, 69, 70), noContext())
		require.True(t, ok)
		assert.Equal(t, harmony.ChordDominantSeventh, got.Tetrad)
		assert.True(t, got.Extensions.Has(harmony.IntMajorThirteenth))
	})
}

// A shape that matches nothing is not a chord. This is what a scoring
// layer could never say, and it is the answer a teaching game needs.
func TestNonChordsAreRefused(t *testing.T) {
	r := newRecognizer(t)

	t.Run("a cluster is not a chord", func(t *testing.T) {
		_, ok := r.Best(held(60, 61, 62), noContext())
		assert.False(t, ok)
	})

	t.Run("nothing held identifies nothing", func(t *testing.T) {
		_, ok := r.Best(held(), noContext())
		assert.False(t, ok)
		assert.Empty(t, r.Identify(held(), noContext()))
	})

	t.Run("two notes are an interval, not a chord", func(t *testing.T) {
		_, ok := r.Best(held(60, 67), noContext())
		assert.False(t, ok,
			"a power chord is a doubled fifth, and no third means no quality")
	})
}

func TestBassOrdersTheReadings(t *testing.T) {
	r := newRecognizer(t)

	t.Run("root position comes before inversion", func(t *testing.T) {
		got, ok := r.Best(held(57, 60, 64, 67), noContext())
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(9), got.Root)
		assert.Equal(t, harmony.ChordMinorSeventh, got.Tetrad)
	})

	t.Run("the same classes over C read as a sixth chord", func(t *testing.T) {
		got, ok := r.Best(held(60, 64, 67, 69), noContext())
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(0), got.Root)
		assert.Equal(t, harmony.ChordMajorSixth, got.Tetrad)
	})

	t.Run("an inversion is identified and marked as one", func(t *testing.T) {
		got, ok := r.Best(held(64, 67, 72), noContext())
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(0), got.Root)
		assert.False(t, got.BassIsRoot)
	})
}

// Four roots explain a diminished seventh chord equally. All four come back,
// and the stickiness rule stops the answer from walking between them.
func TestSymmetricChordsStayAmbiguous(t *testing.T) {
	r := newRecognizer(t)

	all := r.Identify(held(60, 63, 66, 69), noContext())

	roots := map[harmony.PitchClass]bool{}
	for _, reading := range all {
		if reading.Tetrad == harmony.ChordDiminishedSeventh {
			roots[reading.Root] = true
		}
	}
	require.Equal(t,
		map[harmony.PitchClass]bool{0: true, 3: true, 6: true, 9: true},
		roots,
	)

	t.Run("the current reading is kept when it is still valid", func(t *testing.T) {
		first, ok := r.Best(held(60, 63, 66, 69), noContext())
		require.True(t, ok)

		other := all[len(all)-1]
		require.NotEqual(t, first.Root, other.Root)

		got, ok := r.Best(held(60, 63, 66, 69), analysis.Context{Current: other})
		require.True(t, ok)
		assert.Equal(t, other.Root, got.Root,
			"a still valid answer is not replaced by an equal one")
	})
}

func TestTonalityOrdersWithoutForcing(t *testing.T) {
	r := newRecognizer(t)

	cMajor, err := harmony.NewTonality(0, harmony.ScaleMajor)
	require.NoError(t, err)
	inC := analysis.Context{Tonality: cMajor}

	t.Run("a chord foreign to the key is still identified", func(t *testing.T) {
		got, ok := r.Best(held(61, 65, 68), inC)
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(1), got.Root,
			"unusual in C major is not inaudible")
	})

	t.Run("the key never overrides the bass rule", func(t *testing.T) {
		withKey, ok := r.Best(held(57, 60, 64, 67), inC)
		require.True(t, ok)
		without, ok := r.Best(held(57, 60, 64, 67), noContext())
		require.True(t, ok)
		assert.Equal(t, without.Root, withKey.Root)
	})
}

func TestIdentifyRootTestsOneRootOnly(t *testing.T) {
	r := newRecognizer(t)

	t.Run("the right root yields the reading", func(t *testing.T) {
		got, ok := r.IdentifyRoot(60, held(60, 64, 67))
		require.True(t, ok)
		assert.Equal(t, harmony.ChordMajorTriad, got.Tetrad)
	})

	t.Run("a wrong root yields nothing rather than a poor reading", func(t *testing.T) {
		_, ok := r.IdentifyRoot(61, held(60, 64, 67))
		assert.False(t, ok)
	})
}
