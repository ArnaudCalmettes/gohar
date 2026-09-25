package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Same contract as the scale constants: the offsets are the
// specification, the hexadecimal literal is only an encoding of them.
// This table is the reason the literals can be trusted.
func TestChordPatternConstants(t *testing.T) {
	tests := []struct {
		name    string
		pattern harmony.ChordPattern
		offsets []harmony.Semitones
	}{
		{
			name:    "the major triad stacks a major then a minor third",
			pattern: harmony.ChordMajorTriad,
			offsets: []harmony.Semitones{0, 4, 7},
		},
		{
			name:    "the minor triad stacks a minor then a major third",
			pattern: harmony.ChordMinorTriad,
			offsets: []harmony.Semitones{0, 3, 7},
		},
		{
			name:    "the diminished triad stacks two minor thirds",
			pattern: harmony.ChordDiminishedTriad,
			offsets: []harmony.Semitones{0, 3, 6},
		},
		{
			name:    "the augmented triad stacks two major thirds",
			pattern: harmony.ChordAugmentedTriad,
			offsets: []harmony.Semitones{0, 4, 8},
		},
		{
			name:    "the dominant seventh is a major triad with a minor seventh",
			pattern: harmony.ChordDominantSeventh,
			offsets: []harmony.Semitones{0, 4, 7, 10},
		},
		{
			name:    "the major seventh is a major triad with a major seventh",
			pattern: harmony.ChordMajorSeventh,
			offsets: []harmony.Semitones{0, 4, 7, 11},
		},
		{
			name:    "the minor seventh is a minor triad with a minor seventh",
			pattern: harmony.ChordMinorSeventh,
			offsets: []harmony.Semitones{0, 3, 7, 10},
		},
		{
			name:    "the half diminished is a diminished triad with a minor seventh",
			pattern: harmony.ChordHalfDiminished,
			offsets: []harmony.Semitones{0, 3, 6, 10},
		},
		{
			name:    "the diminished seventh stacks three minor thirds",
			pattern: harmony.ChordDiminishedSeventh,
			offsets: []harmony.Semitones{0, 3, 6, 9},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var offsets []harmony.Semitones
			for n := range tc.pattern.Offsets() {
				offsets = append(offsets, n)
			}
			assert.Equal(t, tc.offsets, offsets)
			assert.Equal(t, len(tc.offsets), tc.pattern.Len())

			built, err := harmony.NewChordPattern(tc.offsets...)
			require.NoError(t, err)
			assert.Equal(t, tc.pattern, built,
				"the literal encodes the offsets it claims to",
			)
		})
	}
}

func TestChordPatternInvariants(t *testing.T) {
	patterns := map[string]harmony.ChordPattern{
		"major triad":        harmony.ChordMajorTriad,
		"minor triad":        harmony.ChordMinorTriad,
		"diminished triad":   harmony.ChordDiminishedTriad,
		"augmented triad":    harmony.ChordAugmentedTriad,
		"dominant seventh":   harmony.ChordDominantSeventh,
		"major seventh":      harmony.ChordMajorSeventh,
		"minor seventh":      harmony.ChordMinorSeventh,
		"half diminished":    harmony.ChordHalfDiminished,
		"diminished seventh": harmony.ChordDiminishedSeventh,
	}

	for name, p := range patterns {
		t.Run(name+" contains its own root", func(t *testing.T) {
			assert.True(t, p.HasOffset(0),
				"Contains takes a pattern now, so p.Contains(0) would be "+
					"vacuously true and assert nothing")
		})
		t.Run(name+" sets no bit above twenty three", func(t *testing.T) {
			assert.Zero(t, uint32(p)&0xff000000)
		})
		t.Run(name+" folds to as many classes as it has notes", func(t *testing.T) {
			assert.Equal(t, p.Len(), p.Fold().Len(),
				"none of the declared patterns reaches past the octave, "+
					"so folding loses nothing",
			)
		})
	}
}

// Position matters below the octave as well as above it. These two
// cases are the ones the twenty four bit encoding exists for, and the
// ones a fold to twelve bits destroys.
func TestChordPatternPositionIsNotClass(t *testing.T) {
	ninth, err := harmony.NewChordPattern(0, 4, 7, 10, 14)
	require.NoError(t, err)

	t.Run("a ninth is not a second", func(t *testing.T) {
		assert.True(t, ninth.HasOffset(14))
		assert.False(t, ninth.HasOffset(2))
		assert.True(t, ninth.Has(harmony.IntMajorNinth))
		assert.False(t, ninth.Has(harmony.IntMajorSecond))
	})

	t.Run("folding a ninth yields the class of a second", func(t *testing.T) {
		assert.True(t, ninth.Fold().Contains(2))
	})

	t.Run("a ninth and an added second fold to the same classes", func(t *testing.T) {
		added, err := harmony.NewChordPattern(0, 2, 4, 7, 10)
		require.NoError(t, err)
		assert.Equal(t, ninth.Fold(), added.Fold(),
			"the core stops here: telling them apart is the scoring layer's job",
		)
		assert.NotEqual(t, ninth, added,
			"but the patterns themselves stay distinct",
		)
	})
}

// Has answers on the semitone count, since a bit mask holds no degree.
// That makes it enharmonic by necessity, and the necessity is real: one
// hears a raised ninth in an altered chord and calls it a minor tenth
// only once the chord has been identified, so both have to match the
// same bit until then.
func TestChordPatternHas(t *testing.T) {
	t.Run("a plain seventh chord holds no ninth of any kind", func(t *testing.T) {
		for _, i := range []harmony.Interval{
			harmony.IntMinorNinth, harmony.IntMajorNinth, harmony.IntAugmentedNinth,
		} {
			assert.False(t, harmony.ChordDominantSeventh.Has(i), "%v", i)
		}
	})

	t.Run("an unaltered ninth sits fourteen semitones above the root", func(t *testing.T) {
		p, err := harmony.NewChordPattern(0, 4, 7, 10, 14)
		require.NoError(t, err)
		assert.True(t, p.Has(harmony.IntMajorNinth))
		assert.False(t, p.Has(harmony.IntMinorNinth))
	})

	t.Run("a minor ninth sits one semitone below it", func(t *testing.T) {
		p, err := harmony.NewChordPattern(0, 4, 7, 10, 13)
		require.NoError(t, err)
		assert.True(t, p.Has(harmony.IntMinorNinth))
		assert.False(t, p.Has(harmony.IntMajorNinth))
	})

	t.Run("a raised ninth and a minor tenth match the same bit", func(t *testing.T) {
		p, err := harmony.NewChordPattern(0, 4, 10, 15)
		require.NoError(t, err)
		assert.True(t, p.Has(harmony.IntAugmentedNinth))
		assert.True(t, p.Has(harmony.IntMinorTenth))
		assert.True(t, harmony.IntAugmentedNinth.IsEnharmonic(harmony.IntMinorTenth))
	})

	t.Run("a pattern can hold both altered ninths at once", func(t *testing.T) {
		p, err := harmony.NewChordPattern(0, 4, 10, 13, 15)
		require.NoError(t, err)
		assert.True(t, p.HasAll(harmony.IntMinorNinth, harmony.IntAugmentedNinth))
	})

	t.Run("HasAny is true on the first match, HasAll on none missing", func(t *testing.T) {
		assert.True(t, harmony.ChordDominantSeventh.HasAny(
			harmony.IntMajorThird, harmony.IntMinorThird))
		assert.False(t, harmony.ChordDominantSeventh.HasAll(
			harmony.IntMajorThird, harmony.IntMinorThird))
		assert.False(t, harmony.ChordDominantSeventh.HasAll(),
			"asking for nothing is not a match")
	})
}

func TestNewChordPatternRejectsOutOfRangeOffsets(t *testing.T) {
	_, err := harmony.NewChordPattern(0, 4, 24)
	assert.Error(t, err, "twenty four semitones is past the two octave span")
}

// On a pitch, From preserves the stack: a ninth lands an octave and a
// tone above the root, which is the whole reason the pattern spans
// twenty four bits. On a class it folds and the stack is lost.
func TestChordPatternFrom(t *testing.T) {
	ninth, err := harmony.NewChordPattern(0, 4, 7, 10, 14)
	require.NoError(t, err)

	t.Run("from a pitch the ninth stays above the octave", func(t *testing.T) {
		var got []harmony.Pitch
		for p := range ninth.From(harmony.MiddleC) {
			got = append(got, p)
		}
		assert.Equal(t,
			[]harmony.Pitch{60, 64, 67, 70, 74},
			got,
		)
	})

	t.Run("from a class the octave is lost but the order is not", func(t *testing.T) {
		var got []harmony.PitchClass
		for c := range ninth.From(harmony.PitchClass(0)) {
			got = append(got, c)
		}
		assert.Equal(t,
			[]harmony.PitchClass{0, 4, 7, 10, 2},
			got,
			"the ninth still comes last, folded onto the class of a second, "+
				"so the sequence wraps rather than climbs",
		)
	})

	t.Run("the folded iteration visits the same classes as Fold", func(t *testing.T) {
		var got harmony.PitchSet
		for c := range ninth.From(harmony.PitchClass(0)) {
			got = got.With(c)
		}
		assert.Equal(t, ninth.Fold(), got,
			"the same classes, compared as a set: From has an order and a "+
				"set does not, so comparing the sequences would assert "+
				"something neither one promises",
		)
	})

	t.Run("a triad on a pitch stays within one octave", func(t *testing.T) {
		var got []harmony.Pitch
		for p := range harmony.ChordMajorTriad.From(harmony.MiddleC) {
			got = append(got, p)
		}
		assert.Equal(t, []harmony.Pitch{60, 64, 67}, got)
	})
}

func TestChordPatternWith(t *testing.T) {
	t.Run("adding a seventh to a major triad gives the dominant seventh", func(t *testing.T) {
		assert.Equal(t,
			harmony.ChordDominantSeventh,
			harmony.ChordMajorTriad.With(10),
		)
	})

	t.Run("adding a member already present changes nothing", func(t *testing.T) {
		assert.Equal(t,
			harmony.ChordMajorTriad,
			harmony.ChordMajorTriad.With(4),
		)
	})
}

func TestChordPatternIsSubsetOf(t *testing.T) {
	t.Run("a triad is contained in the seventh chord built on it", func(t *testing.T) {
		assert.True(t, harmony.ChordMajorTriad.IsSubsetOf(harmony.ChordDominantSeventh))
	})

	t.Run("the seventh chord is not contained in its triad", func(t *testing.T) {
		assert.False(t, harmony.ChordDominantSeventh.IsSubsetOf(harmony.ChordMajorTriad))
	})

	t.Run("every pattern contains itself", func(t *testing.T) {
		assert.True(t, harmony.ChordMajorTriad.IsSubsetOf(harmony.ChordMajorTriad))
	})

	t.Run("comparison happens at full resolution, not on folded forms", func(t *testing.T) {
		ninth, err := harmony.NewChordPattern(0, 4, 7, 10, 14)
		require.NoError(t, err)
		added, err := harmony.NewChordPattern(0, 2, 4, 7, 10)
		require.NoError(t, err)

		assert.False(t, added.IsSubsetOf(ninth),
			"the second sits at offset two and the ninth at fourteen, "+
				"so neither contains the other despite folding alike",
		)
		assert.False(t, ninth.IsSubsetOf(added))
	})
}

func TestChord(t *testing.T) {
	g7, err := harmony.NewChord(7, harmony.ChordDominantSeventh)
	require.NoError(t, err)

	t.Run("the set holds the classes the chord sounds", func(t *testing.T) {
		want, err := harmony.NewPitchSet(7, 11, 2, 5)
		require.NoError(t, err)
		assert.Equal(t, want, g7.Set())
	})

	t.Run("an empty pattern makes no chord", func(t *testing.T) {
		_, err := harmony.NewChord(0, harmony.ChordPattern(0))
		assert.Error(t, err, "the bit zero invariant forbids it")
	})

	t.Run("the pitches climb and stay inside the bounds", func(t *testing.T) {
		var got []harmony.Pitch
		for p := range g7.Pitches(harmony.MiddleC, harmony.MiddleC.Transpose(24)) {
			got = append(got, p)
		}
		assert.Equal(t, []harmony.Pitch{67, 71, 74, 77}, got)
	})
}
