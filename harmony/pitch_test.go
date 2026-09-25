package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
)

func TestPitchMidiAnchor(t *testing.T) {
	t.Run("middle C is MIDI sixty", func(t *testing.T) {
		assert.Equal(t, harmony.Pitch(60), harmony.MiddleC)
	})

	t.Run("middle C is class zero in octave four", func(t *testing.T) {
		assert.Equal(t, harmony.PitchClass(0), harmony.MiddleC.Class())
		assert.Equal(t, int8(4), harmony.MiddleC.Octave())
	})

	t.Run("the octave below middle C is octave three", func(t *testing.T) {
		assert.Equal(t, int8(3), harmony.MiddleC.Transpose(-12).Octave())
	})
}

func TestPitchIsValid(t *testing.T) {
	tests := []struct {
		name  string
		pitch harmony.Pitch
		valid bool
	}{
		{"the bottom of the MIDI range is valid", harmony.PitchMin, true},
		{"the top of the MIDI range is valid", harmony.PitchMax, true},
		{"middle C is valid", harmony.MiddleC, true},
		{"a negative pitch is out of range", -1, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.valid, tc.pitch.IsValid())
		})
	}
}

// Transpose is documented as not clamping, because an out of range
// intermediate value is legitimate during a computation. This pins that
// down so nobody adds a clamp later thinking it an improvement.
func TestPitchTransposeDoesNotClamp(t *testing.T) {
	t.Run("moving below the range yields an invalid pitch, not zero", func(t *testing.T) {
		got := harmony.PitchMin.Transpose(-1)
		assert.False(t, got.IsValid())
		assert.NotEqual(t, harmony.PitchMin, got)
	})

	t.Run("an excursion out of range and back is lossless", func(t *testing.T) {
		assert.Equal(t, harmony.MiddleC, harmony.MiddleC.Transpose(-70).Transpose(70))
	})
}

func TestPitchSub(t *testing.T) {
	t.Run("the distance is signed and oriented from the argument", func(t *testing.T) {
		above := harmony.MiddleC.Transpose(7)
		assert.Equal(t, harmony.Semitones(7), above.Sub(harmony.MiddleC))
		assert.Equal(t, harmony.Semitones(-7), harmony.MiddleC.Sub(above))
	})

	t.Run("a pitch is at no distance from itself", func(t *testing.T) {
		assert.Equal(t, harmony.Unison, harmony.MiddleC.Sub(harmony.MiddleC))
	})

	t.Run("the distance can exceed one octave", func(t *testing.T) {
		assert.Equal(t,
			harmony.Semitones(19),
			harmony.MiddleC.Transpose(19).Sub(harmony.MiddleC),
		)
	})
}

// Class is defined for every Pitch, including the out of range ones
// that Transpose is allowed to produce.
func TestPitchClassIsTotal(t *testing.T) {
	t.Run("a negative pitch still has a class", func(t *testing.T) {
		assert.Equal(t, harmony.PitchClass(11), harmony.Pitch(-1).Class())
	})

	t.Run("pitches an octave apart share a class", func(t *testing.T) {
		for n := harmony.Semitones(-24); n <= 24; n += 12 {
			assert.Equal(t,
				harmony.MiddleC.Class(),
				harmony.MiddleC.Transpose(n).Class(),
				"distance %d", n,
			)
		}
	})
}

func TestSemitonesFold(t *testing.T) {
	tests := []struct {
		name     string
		distance harmony.Semitones
		want     harmony.PitchClass
	}{
		{"the unison folds to class zero", harmony.Unison, 0},
		{"an octave folds to class zero", harmony.SemitonesPerOctave, 0},
		{"a fifth folds to class seven", 7, 7},
		{"a descending semitone folds to eleven, not to an error", -1, 11},
		{"a descending octave folds to zero", -12, 0},
		{"a distance past the octave folds into it", 19, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.distance.Fold())
		})
	}
}

func TestSemitonesAbs(t *testing.T) {
	t.Run("a descending distance becomes ascending", func(t *testing.T) {
		assert.Equal(t, harmony.Semitones(7), harmony.Semitones(-7).Abs())
	})

	t.Run("an ascending distance is unchanged", func(t *testing.T) {
		assert.Equal(t, harmony.Semitones(7), harmony.Semitones(7).Abs())
	})
}
