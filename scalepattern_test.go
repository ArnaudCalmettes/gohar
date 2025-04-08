package gohar

import (
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestScalePatternAsPitches(t *testing.T) {
	testCases := []struct {
		Name string
		ScalePattern
		Want []Pitch
	}{
		{
			"major", ScalePatternMajor,
			[]Pitch{PitchC, PitchD, PitchE, PitchF, PitchG, PitchA, PitchB},
		}, {
			"melodic minor", ScalePatternMelodicMinor,
			[]Pitch{PitchC, PitchD, PitchEFlat, PitchF, PitchG, PitchA, PitchB},
		}, {
			"harmonic minor", ScalePatternHarmonicMinor,
			[]Pitch{PitchC, PitchD, PitchEFlat, PitchF, PitchG, PitchAFlat, PitchB},
		}, {
			"harmonic major", ScalePatternHarmonicMajor,
			[]Pitch{PitchC, PitchD, PitchE, PitchF, PitchG, PitchAFlat, PitchB},
		}, {
			"double harmonic major", ScalePatternDoubleHarmonicMajor,
			[]Pitch{PitchC, PitchDFlat, PitchE, PitchF, PitchG, PitchAFlat, PitchB},
		},
	}
	for _, tc := range testCases {
		Expect(t, Equalf(tc.Want, tc.ScalePattern.AsPitches(PitchC), "%s", tc.Name))
	}
}

func BenchmarkScalePatternAsPitches(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out := ScalePatternMajor.AsPitches(0)
		if out == nil {
			b.Fatal()
		}
	}
}

func TestScalePatternIntoPitches(t *testing.T) {
	// Functionality is tested through "AsPitchSlice"
	// that always allocates the right amount of memory.
	// Here we make sure that error conditions are handled properly.
	scalePattern := ScalePatternMajor
	_, err := scalePattern.IntoPitches(nil, 0)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)

	buffer := make([]Pitch, 0)
	_, err = scalePattern.IntoPitches(buffer, 0)
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkScalePatternIntoPitches(b *testing.B) {
	buffer := make([]Pitch, 0, 7)
	for i := 0; i < b.N; i++ {
		_, err := ScalePatternMajor.IntoPitches(buffer, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestScalePatternAsIntervals(t *testing.T) {
	t.Run("major scale", func(t *testing.T) {
		Expect(t,
			Equal(
				[]Interval{
					IntUnisson, IntMajorSecond, IntMajorThird, IntPerfectFourth,
					IntPerfectFifth, IntMajorSixth, IntMajorSeventh,
				},
				ScalePatternMajor.AsIntervals(nil),
			),
		)
	})
	t.Run("major pentatonic", func(t *testing.T) {
		Expect(t,
			Equal(
				[]Interval{
					IntUnisson, IntMajorSecond, IntMajorThird,
					IntPerfectFifth, IntMajorSixth,
				},
				ScalePattern(0b1010010101).AsIntervals([]int8{1, 2, 3, 5, 6}),
			),
		)
	})
}

func BenchmarkScalePatternAsIntervals(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if ScalePatternMajor.AsIntervals(nil) == nil {
			b.Fatal()
		}
	}
}

func TestScalePatternIntoIntervals(t *testing.T) {
	_, err := ScalePatternMajor.IntoIntervals(nil, nil)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)
	_, err = ScalePatternMajor.IntoIntervals([]Interval{}, nil)
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkScalePatternIntoIntervals(b *testing.B) {
	buffer := make([]Interval, 0, 12)
	for i := 0; i < b.N; i++ {
		_, err := ScalePatternMajor.IntoIntervals(buffer, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestScalePatternMode(t *testing.T) {
	testCases := []struct {
		Scale ScalePattern
		Mode  int
		Want  ScalePattern
	}{
		{
			Scale: ScalePatternMajor,
			Mode:  1,
			Want:  ScalePatternMajor,
		},
		{
			Scale: ScalePatternMajor,
			Mode:  4,
			Want:  0b101011010101, // lydian
		},
		{
			Scale: ScalePatternMajor,
			Mode:  7,
			Want:  0b010101101011, // locrian
		},
		{
			Scale: ScalePatternMajor,
			Mode:  0,
			Want:  0, // error
		},
	}

	for _, tc := range testCases {
		Expect(t,
			Equalf(tc.Want, tc.Scale.Mode(tc.Mode), "mode %d of %012b (expected %012b)", tc.Mode, tc.Scale, tc.Want),
		)
	}
}
