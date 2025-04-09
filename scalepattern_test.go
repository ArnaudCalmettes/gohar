package gohar

import (
	"fmt"
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
	_, err = ScalePatternMajor.IntoIntervals(make([]Interval, 0, 12), []int8{1, 2, 3, 5, 6})
	Expect(t,
		IsError(ErrInvalidDegree, err),
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

func TestScalePatternAsNotes(t *testing.T) {
	isNotes := func(want ...Note) CheckFunc[[]Note] {
		return func(have []Note, err error) error {
			if err := NoError(err); err != nil {
				return err
			}
			return Equal(want, have)
		}
	}
	testCases := []struct {
		ScalePattern
		Root    Note
		Degrees []int8
		Check   CheckFunc[[]Note]
	}{
		{
			ScalePatternMajor, NoteC, nil,
			isNotes(NoteC, NoteD, NoteE, NoteF, NoteG, NoteA, NoteB),
		},
		{
			ScalePatternMajor, NoteE.Flat(), nil,
			isNotes(
				NoteE.Flat(),
				NoteF,
				NoteG,
				NoteA.Flat(),
				NoteB.Flat(),
				NoteC.Octave(1),
				NoteD.Octave(1),
			),
		},
		{
			ScalePatternMajor, NoteC.Octave(-1), nil,
			isNotes(
				NoteC.Octave(-1),
				NoteD.Octave(-1),
				NoteE.Octave(-1),
				NoteF.Octave(-1),
				NoteG.Octave(-1),
				NoteA.Octave(-1),
				NoteB.Octave(-1),
			),
		},
		{
			ScalePatternMajor, NoteC.Octave(-2), nil,
			isNotes(
				NoteC.Octave(-2),
				NoteD.Octave(-2),
				NoteE.Octave(-2),
				NoteF.Octave(-2),
				NoteG.Octave(-2),
				NoteA.Octave(-2),
				NoteB.Octave(-2),
			),
		},
	}

	for _, tc := range testCases {
		have, err := tc.ScalePattern.AsNotes(tc.Root, tc.Degrees)
		Expect(t, tc.Check(have, err))
	}
}

func TestScalePatternIntoNotes(t *testing.T) {
	result, err := ScalePatternMajor.IntoNotes(nil, NoteC, nil)
	Expect(t,
		IsEmptySlice(result),
		IsError(ErrNilBuffer, err),
	)
	result, err = ScalePatternMajor.IntoNotes(make([]Note, 0), NoteC, nil)
	Expect(t,
		IsEmptySlice(result),
		IsError(ErrBufferOverflow, err),
	)
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
