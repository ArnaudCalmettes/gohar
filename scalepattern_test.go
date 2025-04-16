package gohar

import (
	"fmt"
	"slices"
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestScalePatternPitches(t *testing.T) {
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
		Expect(t,
			Equalf(
				tc.Want,
				slices.Collect(tc.ScalePattern.Pitches(PitchC)),
				"%s", tc.Name,
			),
		)
	}
}

func BenchmarkScalePatternPitches(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for p := range ScalePatternMajor.Pitches(0) {
			if p < 0 {
				b.Fail()
			}
		}
	}
}

func TestScalePatternIntervals(t *testing.T) {
	t.Run("major scale", func(t *testing.T) {
		Expect(t,
			Equal(
				[]Interval{
					IntUnisson, IntMajorSecond, IntMajorThird, IntPerfectFourth,
					IntPerfectFifth, IntMajorSixth, IntMajorSeventh,
				},
				slices.Collect(ScalePatternMajor.Intervals(nil)),
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
				slices.Collect(ScalePattern(0b1010010101).Intervals([]int8{1, 2, 3, 5, 6})),
			),
		)
	})
}

func TestScalePatternPitchClasses(t *testing.T) {
	testCases := []struct {
		ScalePattern
		Root    PitchClass
		Degrees []int8
		Want    []PitchClass
	}{
		{
			ScalePatternMajor, PitchClassC, nil,
			[]PitchClass{
				PitchClassC,
				PitchClassD,
				PitchClassE,
				PitchClassF,
				PitchClassG,
				PitchClassA,
				PitchClassB,
			},
		},
		{
			ScalePatternMelodicMinor, PitchClassG, nil,
			[]PitchClass{
				PitchClassG,
				PitchClassA,
				PitchClassB.Flat(),
				PitchClassC,
				PitchClassD,
				PitchClassE,
				PitchClassF.Sharp(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			Expect(t,
				Equal(
					tc.Want,
					slices.Collect(
						tc.ScalePattern.PitchClasses(tc.Root, tc.Degrees),
					),
				),
			)
		})
	}

}

func TestScalePatternMode(t *testing.T) {
	isScalePattern := AsCheckFunc(func(want, have ScalePattern) error {
		return Equal(want, have)
	})
	isError := HasError[ScalePattern]
	testCases := []struct {
		Scale ScalePattern
		Mode  int
		Check CheckFunc[ScalePattern]
	}{
		{
			Scale: ScalePatternMajor,
			Mode:  1,
			Check: isScalePattern(ScalePatternMajor),
		},
		{
			Scale: ScalePatternMajor,
			Mode:  4,
			Check: isScalePattern(0b101011010101), // lydian
		},
		{
			Scale: ScalePatternMajor,
			Mode:  7,
			Check: isScalePattern(0b010101101011), // locrian
		},
		{
			Scale: ScalePatternMajor,
			Mode:  0,
			Check: isError(ErrInvalidDegree),
		},
		{
			Scale: ScalePatternMajor,
			Mode:  8,
			Check: isError(ErrInvalidDegree),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v %d", tc.Scale, tc.Mode), func(t *testing.T) {
			got, err := tc.Scale.Mode(tc.Mode)
			Expect(t, tc.Check(got, err))
		})
	}
}
