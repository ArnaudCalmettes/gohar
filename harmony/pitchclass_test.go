package harmony_test

import (
	"strconv"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPitchClassIsValid(t *testing.T) {
	t.Run("every class below twelve is valid", func(t *testing.T) {
		for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
			assert.True(t, c.IsValid(), "class %d", c)
		}
	})

	t.Run("twelve is out of range", func(t *testing.T) {
		assert.False(t, harmony.PitchClass(harmony.PitchClassCount).IsValid())
	})
}

// Transpose on a class is total and closed: it has no failure mode and
// never leaves the domain. That is what lets the rest of the package
// call it without checking.
func TestPitchClassTransposeIsTotalAndClosed(t *testing.T) {
	t.Run("every class and every distance yields a valid class", func(t *testing.T) {
		for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
			for n := harmony.Semitones(-24); n <= 24; n++ {
				assert.True(t, c.Transpose(n).IsValid(),
					"class %d moved by %d", c, n)
			}
		}
	})

	t.Run("a descending semitone from class zero is eleven", func(t *testing.T) {
		assert.Equal(t, harmony.PitchClass(11), harmony.PitchClass(0).Transpose(-1))
	})

	t.Run("an ascending semitone from class eleven is zero", func(t *testing.T) {
		assert.Equal(t, harmony.PitchClass(0), harmony.PitchClass(11).Transpose(1))
	})

	t.Run("an octave in either direction is the identity", func(t *testing.T) {
		for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
			assert.Equal(t, c, c.Transpose(harmony.SemitonesPerOctave))
			assert.Equal(t, c, c.Transpose(-harmony.SemitonesPerOctave))
		}
	})
}

// Up is deliberately not symmetric, and the asymmetry is the whole
// reason it is not called Distance.
func TestPitchClassUp(t *testing.T) {
	t.Run("the ascending distance from zero to eleven is eleven", func(t *testing.T) {
		assert.Equal(t, harmony.Semitones(11), harmony.PitchClass(0).Up(11))
	})

	t.Run("the ascending distance from eleven to zero is one", func(t *testing.T) {
		assert.Equal(t, harmony.Semitones(1), harmony.PitchClass(11).Up(0))
	})

	t.Run("a class is at no distance from itself", func(t *testing.T) {
		for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
			assert.Equal(t, harmony.Unison, c.Up(c))
		}
	})

	t.Run("the result always falls below one octave", func(t *testing.T) {
		for a := harmony.PitchClass(0); a < harmony.PitchClassCount; a++ {
			for b := harmony.PitchClass(0); b < harmony.PitchClassCount; b++ {
				got := a.Up(b)
				assert.GreaterOrEqual(t, got, harmony.Unison)
				assert.Less(t, got, harmony.SemitonesPerOctave)
			}
		}
	})

	t.Run("the two directions add up to an octave", func(t *testing.T) {
		for a := harmony.PitchClass(0); a < harmony.PitchClassCount; a++ {
			for b := harmony.PitchClass(0); b < harmony.PitchClassCount; b++ {
				if a == b {
					continue
				}
				assert.Equal(t, harmony.SemitonesPerOctave, a.Up(b)+b.Up(a))
			}
		}
	})
}

func TestPitchClassToPitch(t *testing.T) {
	t.Run("class zero in octave four is middle C", func(t *testing.T) {
		assert.Equal(t, harmony.MiddleC, harmony.PitchClass(0).Pitch(4))
	})

	t.Run("building a pitch and reading its class round trips", func(t *testing.T) {
		for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
			assert.Equal(t, c, c.Pitch(4).Class())
		}
	})
}

func TestPitchClassSet(t *testing.T) {
	t.Run("a class yields the set holding only itself", func(t *testing.T) {
		for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
			got := c.Set()
			assert.Equal(t, 1, got.Len())
			assert.True(t, got.Contains(c))
		}
	})
}

// The core carries no note names, and String is the place where one
// would sneak back in. A String that guessed a spelling would be a
// naming decision taken without the context needed to justify it.
func TestPitchClassStringCarriesNoName(t *testing.T) {
	for c := harmony.PitchClass(0); c < harmony.PitchClassCount; c++ {
		t.Run("class "+strconv.Itoa(int(c))+" renders as its number", func(t *testing.T) {
			got := c.String()
			n, err := strconv.Atoi(got)
			require.NoError(t, err, "String must stay numeric, got %q", got)
			assert.Equal(t, int(c), n)
		})
	}
}
