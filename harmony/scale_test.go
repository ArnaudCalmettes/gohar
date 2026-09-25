package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every pattern below is described against the major pattern, and only
// against it. Chaining the descriptions instead, minor from major then
// harmonic from minor, moves the reference from line to line and makes
// two patterns look further apart than they are.
//
// Held to one reference, the melodic minor comes out as the major with
// a single lowered third, which is both the shortest true description
// and the one a player can hear.
//
// The declared scale patterns are hand written bit masks. Their
// specification is the list of offsets below, not the literal: if the
// two disagree, the literal is what is wrong.
func TestScalePatternConstants(t *testing.T) {
	tests := []struct {
		name    string
		pattern harmony.ScalePattern
		offsets []harmony.Semitones
	}{
		{
			name:    "the major pattern is the reference the others are described against",
			pattern: harmony.ScaleMajor,
			offsets: []harmony.Semitones{0, 2, 4, 5, 7, 9, 11},
		},
		{
			name:    "the natural minor pattern is the major with a lowered third, sixth and seventh",
			pattern: harmony.ScaleNaturalMinor,
			offsets: []harmony.Semitones{0, 2, 3, 5, 7, 8, 10},
		},
		{
			name:    "the harmonic minor pattern is the major with a lowered third and sixth",
			pattern: harmony.ScaleHarmonicMinor,
			offsets: []harmony.Semitones{0, 2, 3, 5, 7, 8, 11},
		},
		{
			name:    "the melodic minor pattern is the major with a lowered third, and nothing else",
			pattern: harmony.ScaleMelodicMinor,
			offsets: []harmony.Semitones{0, 2, 3, 5, 7, 9, 11},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var offsets []harmony.Semitones
			var degrees []harmony.Degree
			for d, n := range tc.pattern.Offsets() {
				degrees = append(degrees, d)
				offsets = append(offsets, n)
			}

			assert.Equal(t, tc.offsets, offsets)
			assert.Equal(t, len(tc.offsets), tc.pattern.Len())
			assert.Equal(t,
				[]harmony.Degree{1, 2, 3, 4, 5, 6, 7}, degrees,
				"degrees are consecutive and start at one",
			)
		})
	}
}

func TestScalePatternInvariants(t *testing.T) {
	patterns := map[string]harmony.ScalePattern{
		"major":          harmony.ScaleMajor,
		"natural minor":  harmony.ScaleNaturalMinor,
		"harmonic minor": harmony.ScaleHarmonicMinor,
		"melodic minor":  harmony.ScaleMelodicMinor,
	}

	for name, p := range patterns {
		t.Run(name+" contains its own tonic", func(t *testing.T) {
			assert.True(t, p.Contains(0))
		})
		t.Run(name+" is heptatonic", func(t *testing.T) {
			assert.True(t, p.IsHeptatonic())
		})
		t.Run(name+" sets no bit above eleven", func(t *testing.T) {
			assert.Zero(t, uint16(p)&0xf000)
		})
	}
}

// A mode is a rotation, so it preserves the note count and re-anchors
// the pattern on one of its own members. The natural minor is the sixth
// mode of the major scale, which is the cheapest way to catch a
// rotation that went the wrong way.
func TestScalePatternMode(t *testing.T) {
	t.Run("the first mode returns the pattern unchanged", func(t *testing.T) {
		got, ok := harmony.ScaleMajor.Mode(1)
		require.True(t, ok)
		assert.Equal(t, harmony.ScaleMajor, got)
	})

	t.Run("the sixth mode of the major pattern is the natural minor", func(t *testing.T) {
		got, ok := harmony.ScaleMajor.Mode(6)
		require.True(t, ok)
		assert.Equal(t, harmony.ScaleNaturalMinor, got)
	})

	t.Run("every mode keeps the note count", func(t *testing.T) {
		for d := harmony.Degree(1); d <= 7; d++ {
			got, ok := harmony.ScaleMajor.Mode(d)
			require.True(t, ok)
			assert.Equal(t, 7, got.Len())
			assert.True(t, got.Contains(0), "a mode is anchored on its own tonic")
		}
	})

	t.Run("a degree beyond the pattern has no mode", func(t *testing.T) {
		_, ok := harmony.ScaleMajor.Mode(8)
		assert.False(t, ok)
	})
}

// ScalePattern.At leaves pattern space, NewScalePattern comes back into
// it. The pair has to round trip or the two representations drift.
func TestScalePatternAnchoring(t *testing.T) {
	t.Run("anchoring on a tonic and reading back returns the pattern", func(t *testing.T) {
		for tonic := harmony.PitchClass(0); tonic < harmony.PitchClassCount; tonic++ {
			set := harmony.ScaleMajor.At(tonic)
			back, err := harmony.NewScalePattern(set, tonic)
			require.NoError(t, err)
			assert.Equal(t, harmony.ScaleMajor, back, "tonic %d", tonic)
		}
	})

	t.Run("anchoring preserves the note count", func(t *testing.T) {
		for tonic := harmony.PitchClass(0); tonic < harmony.PitchClassCount; tonic++ {
			assert.Equal(t, harmony.ScaleMajor.Len(), harmony.ScaleMajor.At(tonic).Len())
		}
	})

	t.Run("the anchored set contains its tonic", func(t *testing.T) {
		for tonic := harmony.PitchClass(0); tonic < harmony.PitchClassCount; tonic++ {
			assert.True(t, harmony.ScaleMajor.At(tonic).Contains(tonic))
		}
	})

	t.Run("a tonic outside the set is refused", func(t *testing.T) {
		set, err := harmony.NewPitchSet(0, 2, 4, 5, 7, 9, 11)
		require.NoError(t, err)

		_, err = harmony.NewScalePattern(set, 1)
		assert.Error(t, err,
			"bit zero would come out clear, and a candidate tonic the set "+
				"does not contain is one to reject, not one to repair",
		)
	})
}

func TestScalePatternOffset(t *testing.T) {
	t.Run("each degree of the major pattern sits where it should", func(t *testing.T) {
		want := []harmony.Semitones{0, 2, 4, 5, 7, 9, 11}
		for i, n := range want {
			got, ok := harmony.ScaleMajor.Offset(harmony.Degree(i + 1))
			require.True(t, ok, "degree %d", i+1)
			assert.Equal(t, n, got)
		}
	})

	t.Run("a degree beyond the pattern has no offset", func(t *testing.T) {
		_, ok := harmony.ScaleMajor.Offset(8)
		assert.False(t, ok)
	})

	t.Run("degree zero does not exist, since degrees start at one", func(t *testing.T) {
		_, ok := harmony.ScaleMajor.Offset(0)
		assert.False(t, ok)
	})
}

// From is the generic method that replaced six near identical ones. The
// two instantiations behave differently on purpose: on a class it folds
// into one octave, on a pitch it climbs.
func TestScalePatternFrom(t *testing.T) {
	t.Run("from a class it yields folded classes", func(t *testing.T) {
		var got []harmony.PitchClass
		for c := range harmony.ScaleMajor.From(harmony.PitchClass(0)) {
			got = append(got, c)
		}
		assert.Equal(t,
			[]harmony.PitchClass{0, 2, 4, 5, 7, 9, 11},
			got,
		)
	})

	t.Run("from a class it wraps around the octave", func(t *testing.T) {
		var got []harmony.PitchClass
		for c := range harmony.ScaleMajor.From(harmony.PitchClass(7)) {
			got = append(got, c)
		}
		assert.Equal(t,
			[]harmony.PitchClass{7, 9, 11, 0, 2, 4, 6},
			got,
			"the G major scale folded into one octave",
		)
	})

	t.Run("from a pitch it climbs instead of wrapping", func(t *testing.T) {
		var got []harmony.Pitch
		for p := range harmony.ScaleMajor.From(harmony.MiddleC) {
			got = append(got, p)
		}
		assert.Equal(t,
			[]harmony.Pitch{60, 62, 64, 65, 67, 69, 71},
			got,
		)
	})

	t.Run("both instantiations yield the same number of elements", func(t *testing.T) {
		classes := 0
		for range harmony.ScaleMajor.From(harmony.PitchClass(0)) {
			classes++
		}
		pitches := 0
		for range harmony.ScaleMajor.From(harmony.MiddleC) {
			pitches++
		}
		assert.Equal(t, harmony.ScaleMajor.Len(), classes)
		assert.Equal(t, harmony.ScaleMajor.Len(), pitches)
	})
}

func TestScale(t *testing.T) {
	dMajor, err := harmony.NewScale(2, harmony.ScaleMajor)
	require.NoError(t, err)

	t.Run("the classes run from the tonic upward and wrap", func(t *testing.T) {
		var got []harmony.PitchClass
		for c := range dMajor.Classes() {
			got = append(got, c)
		}
		assert.Equal(t,
			[]harmony.PitchClass{2, 4, 6, 7, 9, 11, 1},
			got,
		)
	})

	t.Run("the set holds the same classes regardless of order", func(t *testing.T) {
		want, err := harmony.NewPitchSet(1, 2, 4, 6, 7, 9, 11)
		require.NoError(t, err)
		assert.Equal(t, want, dMajor.Set())
	})

	t.Run("degrees resolve to classes", func(t *testing.T) {
		fifth, ok := dMajor.Degree(5)
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(9), fifth)
	})

	t.Run("a degree beyond the pattern resolves to nothing", func(t *testing.T) {
		_, ok := dMajor.Degree(8)
		assert.False(t, ok)
	})

	t.Run("an empty pattern makes no scale", func(t *testing.T) {
		_, err := harmony.NewScale(0, harmony.ScalePattern(0))
		assert.Error(t, err, "the bit zero invariant forbids it")
	})
}

// Pitches takes bounds as pitches rather than octave numbers, so a
// keyboard with an odd range needs no arithmetic at the call site.
func TestScalePitches(t *testing.T) {
	cMajor, err := harmony.NewScale(0, harmony.ScaleMajor)
	require.NoError(t, err)

	t.Run("both bounds are inclusive", func(t *testing.T) {
		var got []harmony.Pitch
		for p := range cMajor.Pitches(harmony.MiddleC, harmony.MiddleC.Transpose(12)) {
			got = append(got, p)
		}
		assert.Equal(t,
			[]harmony.Pitch{60, 62, 64, 65, 67, 69, 71, 72},
			got,
		)
	})

	t.Run("a bound off the scale is skipped rather than emitted", func(t *testing.T) {
		var got []harmony.Pitch
		for p := range cMajor.Pitches(61, 66) {
			got = append(got, p)
		}
		assert.Equal(t, []harmony.Pitch{62, 64, 65}, got)
	})

	t.Run("an empty range yields nothing", func(t *testing.T) {
		for range cMajor.Pitches(72, 60) {
			t.Fatal("a descending range yielded a pitch")
		}
	})
}
