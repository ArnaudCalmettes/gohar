package gohar

import (
	"fmt"
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestScaleStringer(t *testing.T) {
	testCases := []struct {
		Scale
		WantString string
	}{
		{Scale{NoteC, ScalePatternMajor}, "Scale(C0:101010110101)"},
		{Scale{NoteD, ScalePatternMelodicMinor}, "Scale(D0:101010101101)"},
	}

	for _, tc := range testCases {
		Expect(t,
			Equalf(tc.WantString, tc.Scale.String(), "string"),
		)
	}
}

func TestGetScale(t *testing.T) {
	_, err := GetScale(NoteC, "does not exist")
	Expect(t, IsError(ErrUnknownScalePattern, err))
}

func TestScaleAsNoteSlice(t *testing.T) {
	testCases := []struct {
		Root  Note
		Label string
		Want  []Note
	}{
		{
			NoteC, "major",
			[]Note{NoteC, NoteD, NoteE, NoteF, NoteG, NoteA, NoteB},
		},
		{
			NoteD, "major",
			[]Note{NoteD, NoteE, NoteF.Sharp(), NoteG, NoteA, NoteB, NoteC.Sharp().Octave(1)},
		},
		{
			NoteE, "melodic minor",
			[]Note{NoteE, NoteF.Sharp(), NoteG, NoteA, NoteB, NoteC.Sharp().Octave(1), NoteD.Sharp().Octave(1)},
		},
		{
			NoteF, "harmonic minor",
			[]Note{NoteF, NoteG, NoteA.Flat(), NoteB.Flat(), NoteC.Octave(1), NoteD.Flat().Octave(1), NoteE.Octave(1)},
		},
		{
			NoteC.Sharp(), "harmonic major",
			[]Note{NoteC.Sharp(), NoteD.Sharp(), NoteE.Sharp(), NoteF.Sharp(), NoteG.Sharp(), NoteA, NoteB.Sharp()},
		},
		{
			NoteA.Octave(-1), "double harmonic major",
			[]Note{NoteA.Octave(-1), NoteB.Flat().Octave(-1), NoteC.Sharp(), NoteD, NoteE, NoteF, NoteG.Sharp()},
		},
	}

	for _, tc := range testCases {
		scale, err := GetScale(tc.Root, tc.Label)
		Require(t, NoError(err))
		Expect(t, Equalf(
			fmt.Sprint(tc.Want),
			fmt.Sprint(scale.AsNoteSlice()),
			"%s", scale,
		))
	}
}

func BenchmarkScaleAsNoteSlice(b *testing.B) {
	scale, _ := GetScale(NoteE.Flat(), "harmonic minor")
	for i := 0; i < b.N; i++ {
		notes := scale.AsNoteSlice()
		if len(notes) == 0 {
			b.Fatal(notes)
		}
	}
}

func TestScaleAsNoteSliceInto(t *testing.T) {
	cMajor, _ := GetScale(NoteC, "major")
	_, err := cMajor.AsNoteSliceInto(nil)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)
	_, err = cMajor.AsNoteSliceInto(make([]Note, 0))
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkScaleAsNoteSliceInto(b *testing.B) {
	notes := make([]Note, 0, 12)
	scale, _ := GetScale(NoteE.Flat(), "harmonic minor")
	for i := 0; i < b.N; i++ {
		_, err := scale.AsNoteSliceInto(notes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestScaleAsPitchSlice(t *testing.T) {
	dMajor, _ := GetScale(NoteD, "major")
	have := dMajor.AsPitchSlice()
	Expect(t,
		Equal(
			[]Pitch{PitchD, PitchE, PitchFSharp, PitchG, PitchA, PitchB, PitchCSharp.Add(12)},
			have,
		),
	)
}

func BenchmarkScaleAsPitchSlice(b *testing.B) {
	scale, err := GetScale(NoteD, "major")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		if scale.AsPitchSlice() == nil {
			b.Fatal()
		}
	}
}

func TestScaleAsPitchSliceInto(t *testing.T) {
	buffer := make([]Pitch, 0, 7)
	dMajor, _ := GetScale(NoteD, "major")
	have, err := dMajor.AsPitchSliceInto(buffer)
	Expect(t,
		NoError(err),
		Equal(
			[]Pitch{PitchD, PitchE, PitchFSharp, PitchG, PitchA, PitchB, PitchCSharp.Add(12)},
			have,
		),
	)
}

func BenchmarkScaleAsPitchSliceInto(b *testing.B) {
	buffer := make([]Pitch, 0, 7)
	scale, err := GetScale(NoteD, "major")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = scale.AsPitchSliceInto(buffer)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestScalePatternAsPitchSlice(t *testing.T) {
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
		Expect(t, Equalf(tc.Want, tc.ScalePattern.AsPitchSlice(), "%s", tc.Name))
	}
}

func BenchmarkScalePatternAsPitchSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out := ScalePatternMajor.AsPitchSlice()
		if out == nil {
			b.Fatal()
		}
	}
}

func TestScalePatternAsPitchSliceInto(t *testing.T) {
	// Functionality is tested through "AsPitchSlice"
	// that always allocates the right amount of memory.
	// Here we make sure that error conditions are handled properly.
	scalePattern := ScalePatternMajor
	_, err := scalePattern.AsPitchSliceInto(nil)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)

	buffer := make([]Pitch, 0)
	_, err = scalePattern.AsPitchSliceInto(buffer)
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkScalePatternAsPitchSliceInto(b *testing.B) {
	buffer := make([]Pitch, 0, 7)
	for i := 0; i < b.N; i++ {
		_, err := ScalePatternMajor.AsPitchSliceInto(buffer)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestScalePatternAsIntervalSlice(t *testing.T) {
	Expect(t,
		Equal(
			[]Interval{
				IntUnisson, IntMajorSecond, IntMajorThird, IntPerfectFourth,
				IntPerfectFifth, IntMajorSixth, IntMajorSeventh,
			},
			ScalePatternMajor.AsIntervalSlice(),
		),
	)
}

func BenchmarkScalePatternAsIntervalSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if ScalePatternMajor.AsIntervalSlice() == nil {
			b.Fatal()
		}
	}
}

func TestScalePatternAsIntervalSliceInto(t *testing.T) {
	_, err := ScalePatternMajor.AsIntervalSliceInto(nil)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)
	_, err = ScalePatternMajor.AsIntervalSliceInto([]Interval{})
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkScalePatternAsIntervalSliceInto(b *testing.B) {
	buffer := make([]Interval, 0, 12)
	for i := 0; i < b.N; i++ {
		_, err := ScalePatternMajor.AsIntervalSliceInto(buffer)
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
			Want:  0b010101101011, // locrian
		},
	}

	for _, tc := range testCases {
		Expect(t,
			Equalf(tc.Want, tc.Scale.Mode(tc.Mode), "mode %d of %012b (expected %012b)", tc.Mode, tc.Scale, tc.Want),
		)
	}
}
