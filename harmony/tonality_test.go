package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustScalePattern(t *testing.T, offsets ...harmony.PitchClass) harmony.ScalePattern {
	t.Helper()
	s, err := harmony.NewPitchSet(offsets...)
	require.NoError(t, err)
	p, err := harmony.NewScalePattern(s, 0)
	require.NoError(t, err)
	return p
}

// The heptatonic constraint is the reason the type exists. A pattern
// that is not heptatonic has no tonality, and the refusal is an answer
// rather than an obstacle.
func TestNewTonalityRequiresSevenNotes(t *testing.T) {
	accepted := map[string]harmony.ScalePattern{
		"major":          harmony.ScaleMajor,
		"natural minor":  harmony.ScaleNaturalMinor,
		"harmonic minor": harmony.ScaleHarmonicMinor,
		"melodic minor":  harmony.ScaleMelodicMinor,
	}

	for name, p := range accepted {
		t.Run("the "+name+" pattern makes a tonality", func(t *testing.T) {
			got, err := harmony.NewTonality(0, p)
			require.NoError(t, err)
			assert.False(t, got.IsZero())
		})
	}

	t.Run("a pentatonic pattern has no tonality", func(t *testing.T) {
		_, err := harmony.NewTonality(0, mustScalePattern(t, 0, 2, 4, 7, 9))
		assert.Error(t, err)
	})

	t.Run("a whole tone pattern has no tonality", func(t *testing.T) {
		_, err := harmony.NewTonality(0, mustScalePattern(t, 0, 2, 4, 6, 8, 10))
		assert.Error(t, err)
	})

	t.Run("an octatonic pattern has no tonality", func(t *testing.T) {
		_, err := harmony.NewTonality(0, mustScalePattern(t, 0, 2, 3, 5, 6, 8, 9, 11))
		assert.Error(t, err)
	})
}

// The zero value means no context inferred. Treating it as C major
// would make the engine claim a key it never found, and an interface
// built on that lies to the player.
func TestZeroTonalityIsNotCMajor(t *testing.T) {
	var none harmony.Tonality

	cMajor, err := harmony.NewTonality(0, harmony.ScaleMajor)
	require.NoError(t, err)

	t.Run("the zero value reports itself as absent", func(t *testing.T) {
		assert.True(t, none.IsZero())
	})

	t.Run("C major is present", func(t *testing.T) {
		assert.False(t, cMajor.IsZero())
	})

	t.Run("the two are not equal", func(t *testing.T) {
		assert.NotEqual(t, cMajor, none)
	})

	t.Run("the zero value holds no class", func(t *testing.T) {
		assert.True(t, none.Set().IsEmpty())
		for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
			assert.False(t, none.Contains(c), "class %d", c)
		}
	})
}

func TestTonalityDegreeOf(t *testing.T) {
	cMajor, err := harmony.NewTonality(0, harmony.ScaleMajor)
	require.NoError(t, err)

	t.Run("the tonic sits at the first degree", func(t *testing.T) {
		got, ok := cMajor.DegreeOf(0)
		require.True(t, ok)
		assert.Equal(t, harmony.Degree(1), got)
	})

	t.Run("the fifth degree is seven semitones above the tonic", func(t *testing.T) {
		got, ok := cMajor.DegreeOf(7)
		require.True(t, ok)
		assert.Equal(t, harmony.Degree(5), got)
	})

	t.Run("a foreign class has no degree", func(t *testing.T) {
		_, ok := cMajor.DegreeOf(1)
		assert.False(t, ok,
			"an accidental is a note foreign to the key, not a degree of it",
		)
	})

	t.Run("degree and class round trip in both directions", func(t *testing.T) {
		for d := harmony.Degree(1); d <= 7; d++ {
			c, ok := cMajor.Degree(d)
			require.True(t, ok, "degree %d", d)

			back, ok := cMajor.DegreeOf(c)
			require.True(t, ok, "class %d", c)
			assert.Equal(t, d, back)
		}
	})

	t.Run("there is no eighth degree", func(t *testing.T) {
		_, ok := cMajor.Degree(8)
		assert.False(t, ok)
	})
}

func TestTonalityTranspose(t *testing.T) {
	cMajor, err := harmony.NewTonality(0, harmony.ScaleMajor)
	require.NoError(t, err)

	t.Run("the pattern survives and the tonic moves", func(t *testing.T) {
		got := cMajor.Transpose(7)
		assert.Equal(t, harmony.PitchClass(7), got.Tonic())
		assert.Equal(t, cMajor.Pattern(), got.Pattern())
	})

	t.Run("the result stays a usable tonality at every distance", func(t *testing.T) {
		for n := harmony.Semitones(-12); n <= 12; n++ {
			got := cMajor.Transpose(n)
			assert.False(t, got.IsZero(), "distance %d", n)
			assert.True(t, got.Pattern().IsHeptatonic(), "distance %d", n)
		}
	})

	t.Run("degrees move with the tonic", func(t *testing.T) {
		got := cMajor.Transpose(7)
		fifth, ok := got.Degree(5)
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(2), fifth)
	})
}

// Widening to a Scale is one way: a Scale accepts any pattern, so
// coming back has to go through the check again.
func TestTonalityWidensToScale(t *testing.T) {
	cMajor, err := harmony.NewTonality(0, harmony.ScaleMajor)
	require.NoError(t, err)

	got := cMajor.Scale()
	assert.Equal(t, harmony.PitchClass(0), got.Tonic)
	assert.Equal(t, harmony.ScaleMajor, got.Pattern)
	assert.Equal(t, cMajor.Set(), got.Set())
}
