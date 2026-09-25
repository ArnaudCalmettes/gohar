package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The literal is the designation: 0x131 is the tetrachord 1 3 1.
func TestTetrachordReadsAsItsSteps(t *testing.T) {
	tests := []struct {
		name  string
		t     harmony.Tetrachord
		steps string
		span  harmony.Semitones
	}{
		{"major", harmony.TetrachordMajor, "2 2 1", 5},
		{"minor", harmony.TetrachordMinor, "2 1 2", 5},
		{"phrygian", harmony.TetrachordPhrygian, "1 2 2", 5},
		{"harmonic", harmony.TetrachordHarmonic, "1 3 1", 5},
		{"lydian", harmony.TetrachordLydian, "2 2 2", 6},
	}

	for _, tc := range tests {
		t.Run("the "+tc.name+" tetrachord is "+tc.steps, func(t *testing.T) {
			assert.Equal(t, tc.steps, tc.t.String())
			assert.Equal(t, tc.span, tc.t.Span())
			assert.True(t, tc.t.IsValid())
		})
	}

	t.Run("the lydian is the only named one wider than a fourth", func(t *testing.T) {
		assert.False(t, harmony.TetrachordLydian.IsFourth())
		for _, tc := range tests[:4] {
			assert.True(t, tc.t.IsFourth(), tc.name)
		}
	})

	t.Run("building from steps gives the same value as the literal", func(t *testing.T) {
		got, err := harmony.NewTetrachord(1, 3, 1)
		require.NoError(t, err)
		assert.Equal(t, harmony.TetrachordHarmonic, got)
	})
}

func TestNewTetrachordRefusesWhatIsNotOne(t *testing.T) {
	t.Run("a step of zero would double a note", func(t *testing.T) {
		_, err := harmony.NewTetrachord(2, 0, 3)
		assert.Error(t, err)
	})

	t.Run("a step the encoding cannot hold", func(t *testing.T) {
		_, err := harmony.NewTetrachord(1, 16, 1)
		assert.Error(t, err)
	})

	t.Run("the zero value is not a tetrachord", func(t *testing.T) {
		assert.False(t, harmony.Tetrachord(0).IsValid())
	})
}

// The five mother scales read as pairs of tetrachords, and three of them
// are named after theirs.
func TestMotherScalesReadAsTetrachords(t *testing.T) {
	tests := []struct {
		name         string
		system       harmony.System
		lower, upper harmony.Tetrachord
	}{
		{"natural, major over major", harmony.NaturalMajor,
			harmony.TetrachordMajor, harmony.TetrachordMajor},
		{"melodic minor, minor over major", harmony.MelodicMinor,
			harmony.TetrachordMinor, harmony.TetrachordMajor},
		{"harmonic minor, minor over harmonic", harmony.HarmonicMinor,
			harmony.TetrachordMinor, harmony.TetrachordHarmonic},
		{"harmonic major, major over harmonic", harmony.HarmonicMajor,
			harmony.TetrachordMajor, harmony.TetrachordHarmonic},
		{"double harmonic, harmonic over harmonic", harmony.DoubleHarmonicMajor,
			harmony.TetrachordHarmonic, harmony.TetrachordHarmonic},
	}

	for _, tc := range tests {
		t.Run("the "+tc.name, func(t *testing.T) {
			got, ok := tc.system.Pattern().Tetrachords()
			require.True(t, ok)
			assert.Equal(t, tc.lower, got.Lower)
			assert.Equal(t, tc.upper, got.Upper)
			assert.Equal(t, harmony.Semitones(2), got.Gap,
				"a whole tone between two tetrachords that each span a fourth")
		})
	}
}

// The inter-tetrachordal space shrinks when a tetrachord is wider than a
// fourth. The lydian tetrachord reaches up to the tritone and leaves a
// semitone.
func TestGapShrinksAroundTheLydianTetrachord(t *testing.T) {
	t.Run("the lydian mode is lydian, a semitone, then major", func(t *testing.T) {
		p, ok := harmony.NaturalMajor.Mode(4)
		require.True(t, ok)

		got, ok := p.Tetrachords()
		require.True(t, ok)
		assert.Equal(t, harmony.TetrachordLydian, got.Lower)
		assert.Equal(t, harmony.Semitones(1), got.Gap)
		assert.Equal(t, harmony.TetrachordMajor, got.Upper)
	})

	t.Run("the locrian is phrygian, a semitone, then lydian", func(t *testing.T) {
		p, ok := harmony.NaturalMajor.Mode(7)
		require.True(t, ok)

		got, ok := p.Tetrachords()
		require.True(t, ok)
		assert.Equal(t, harmony.TetrachordPhrygian, got.Lower)
		assert.Equal(t, harmony.Semitones(1), got.Gap)
		assert.Equal(t, harmony.TetrachordLydian, got.Upper)
	})
}

// Seven notes force the split, so the reading always succeeds on a
// heptatonic pattern. In the altered systems it often produces a shape
// with no name, and it must say so by designating it by its steps.
func TestTetrachordsAreMechanicalOutsideTheNaturalSystem(t *testing.T) {
	t.Run("the seventh mode of the harmonic minor starts on 1 2 1", func(t *testing.T) {
		p, ok := harmony.HarmonicMinor.Mode(7)
		require.True(t, ok)

		got, ok := p.Tetrachords()
		require.True(t, ok)
		assert.Equal(t, "1 2 1", got.Lower.String())
		assert.False(t, got.Lower.IsFourth(), "it spans only a major third")
	})

	t.Run("every mode of every system splits into an octave", func(t *testing.T) {
		for s := harmony.System(0); s < harmony.SystemCount; s++ {
			for d := harmony.Degree(1); d <= 7; d++ {
				p, ok := s.Mode(d)
				require.True(t, ok)

				got, ok := p.Tetrachords()
				require.True(t, ok, "system %d degree %d", s, d)
				assert.Equal(t, harmony.SemitonesPerOctave,
					got.Lower.Span()+got.Gap+got.Upper.Span(),
					"system %d degree %d", s, d)
			}
		}
	})
}

func TestTetrachordsNeedsSevenNotes(t *testing.T) {
	pentatonic, err := harmony.NewScalePattern(harmony.PitchSet(0b001010_010101), 0)
	require.NoError(t, err)

	_, ok := pentatonic.Tetrachords()
	assert.False(t, ok, "five notes cannot be read as two tetrachords")
}
