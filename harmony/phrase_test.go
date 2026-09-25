package harmony_test

import (
	"slices"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func phrase(t *testing.T, pitches ...harmony.Pitch) harmony.Phrase {
	t.Helper()
	p, err := harmony.NewPhrase(pitches...)
	require.NoError(t, err)
	return p
}

// The premise again: a player without perfect pitch starts where they
// like, so the shape is what is being asked for.
func TestPhraseIgnoresTranspositionWithoutLosingIt(t *testing.T) {
	expected := phrase(t, 60, 62, 64, 60)

	t.Run("the same shape from another note is the same shape", func(t *testing.T) {
		got := expected.Compare(phrase(t, 67, 69, 71, 67))
		assert.True(t, got.Same)
		assert.Equal(t, -1, got.FirstDivergence)
		assert.Equal(t, harmony.Semitones(7), got.Shift)
		assert.True(t, got.Transposed)
	})

	t.Run("the same start reports no shift", func(t *testing.T) {
		got := expected.Compare(expected)
		assert.True(t, got.Same)
		assert.False(t, got.Transposed)
		assert.Equal(t, harmony.Unison, got.Shift)
	})

	t.Run("an octave up is a transposition of twelve, not of nothing", func(t *testing.T) {
		got := expected.Compare(phrase(t, 72, 74, 76, 72))
		assert.True(t, got.Same)
		assert.Equal(t, harmony.Semitones(12), got.Shift,
			"a melodic register is audible, unlike a chord root's octave")
	})

	t.Run("a shape declared with no performance reports no shift", func(t *testing.T) {
		declared, err := harmony.NewPhraseOffsets(0, 2, 4, 0)
		require.NoError(t, err)

		got := declared.Compare(phrase(t, 67, 69, 71, 67))
		assert.True(t, got.Same)
		assert.False(t, got.Transposed)
	})
}

// The one rule that separates a phrase from a progression. Folding the
// offsets, which is right for chord roots, would make these two the
// same shape.
func TestPhraseOffsetsAreSignedAndUnfolded(t *testing.T) {
	up := phrase(t, 60, 71)
	down := phrase(t, 60, 59)

	t.Run("a major seventh up is not a minor second down", func(t *testing.T) {
		assert.False(t, up.Compare(down).Same,
			"one leaps and the other steps, and they land on the same class")
		assert.Equal(t, 1, up.Compare(down).FirstDivergence)
	})

	t.Run("an offset may exceed an octave", func(t *testing.T) {
		wide := phrase(t, 60, 84)
		var got []harmony.Semitones
		for _, n := range wide.Offsets() {
			got = append(got, n)
		}
		assert.Equal(t, []harmony.Semitones{0, 24}, got)
	})

	t.Run("a descending phrase keeps its sign", func(t *testing.T) {
		var got []harmony.Semitones
		for _, n := range phrase(t, 60, 55, 48).Offsets() {
			got = append(got, n)
		}
		assert.Equal(t, []harmony.Semitones{0, -5, -12}, got)
	})
}

func TestPhraseLocatesTheDivergence(t *testing.T) {
	expected := phrase(t, 60, 62, 64, 65)

	t.Run("the index is where it parted, not how much", func(t *testing.T) {
		got := expected.Compare(phrase(t, 67, 69, 72, 72))
		assert.False(t, got.Same)
		assert.Equal(t, 2, got.FirstDivergence,
			"the first two notes were right, and that is what to tell the player")
	})

	t.Run("a wrong first move parts at index one", func(t *testing.T) {
		got := expected.Compare(phrase(t, 67, 68, 69, 70))
		assert.Equal(t, 1, got.FirstDivergence,
			"index zero is the origin and always matches")
	})

	t.Run("a shorter performance parts past its end", func(t *testing.T) {
		got := expected.Compare(phrase(t, 67, 69))
		assert.False(t, got.Same)
		assert.Equal(t, 2, got.FirstDivergence)
	})
}

func TestPhraseAt(t *testing.T) {
	shape, err := harmony.NewPhraseOffsets(0, 2, 4, 0)
	require.NoError(t, err)

	t.Run("voicing from a pitch gives that performance", func(t *testing.T) {
		assert.Equal(t,
			[]harmony.Pitch{67, 69, 71, 67},
			slices.Collect(shape.At(67)),
		)
	})

	t.Run("voicing then rebuilding gives the shape back", func(t *testing.T) {
		for root := harmony.Pitch(48); root <= 72; root++ {
			back, err := harmony.NewPhrase(slices.Collect(shape.At(root))...)
			require.NoError(t, err)
			assert.True(t, shape.Compare(back).Same, "root %d", root)
		}
	})
}

// The contour is what a device with no pitch can express, and the least
// a listener can be asked to hear.
func TestPhraseContour(t *testing.T) {
	t.Run("a contour has one fewer element than the phrase has notes", func(t *testing.T) {
		p := phrase(t, 60, 62, 64, 60)
		assert.Len(t, slices.Collect(p.Contour()), p.Len()-1)
	})

	t.Run("directions read up, down and neither", func(t *testing.T) {
		assert.Equal(t,
			[]harmony.Direction{harmony.Up, harmony.Level, harmony.Down},
			slices.Collect(phrase(t, 60, 64, 64, 60).Contour()),
		)
	})

	t.Run("many phrases share a contour, and that is the point", func(t *testing.T) {
		gentle := phrase(t, 60, 62, 64)
		steep := phrase(t, 60, 67, 79)

		assert.False(t, gentle.Compare(steep).Same,
			"the distances differ")
		assert.True(t, gentle.CompareContour(steep).Same,
			"the directions do not, which is what a beginner is asked for first")
	})

	t.Run("a wrong direction is located at its step", func(t *testing.T) {
		got := phrase(t, 60, 62, 64).CompareContour(phrase(t, 60, 62, 59))
		assert.False(t, got.Same)
		assert.Equal(t, 1, got.FirstDivergence)
	})
}

func TestNewPhraseRejectsWhatCannotBeAShape(t *testing.T) {
	t.Run("a phrase with no note", func(t *testing.T) {
		_, err := harmony.NewPhrase()
		assert.Error(t, err)
	})

	t.Run("offsets that do not start at the origin", func(t *testing.T) {
		_, err := harmony.NewPhraseOffsets(3, 5)
		assert.Error(t, err,
			"the first note is the origin, so it sits at zero by construction")
	})

	t.Run("a single note is a shape, if a dull one", func(t *testing.T) {
		p, err := harmony.NewPhrase(60)
		require.NoError(t, err)
		assert.Equal(t, 1, p.Len())
		assert.Empty(t, slices.Collect(p.Contour()))
	})
}
