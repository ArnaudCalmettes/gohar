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
			fmt.Sprint(scale.AsNotes()),
			"%s", scale,
		))
	}
}

func BenchmarkScaleAsNoteSlice(b *testing.B) {
	scale, _ := GetScale(NoteE.Flat(), "harmonic minor")
	for i := 0; i < b.N; i++ {
		notes := scale.AsNotes()
		if len(notes) == 0 {
			b.Fatal(notes)
		}
	}
}

func TestScaleAsNoteSliceInto(t *testing.T) {
	cMajor, _ := GetScale(NoteC, "major")
	_, err := cMajor.IntoNotes(nil)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)
	_, err = cMajor.IntoNotes(make([]Note, 0))
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkScaleAsNoteSliceInto(b *testing.B) {
	notes := make([]Note, 0, 12)
	scale, _ := GetScale(NoteE.Flat(), "harmonic minor")
	for i := 0; i < b.N; i++ {
		_, err := scale.IntoNotes(notes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestScaleAsPitchSlice(t *testing.T) {
	dMajor, _ := GetScale(NoteD, "major")
	have := dMajor.AsPitches()
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
		if scale.AsPitches() == nil {
			b.Fatal()
		}
	}
}

func TestScaleAsPitchSliceInto(t *testing.T) {
	buffer := make([]Pitch, 0, 7)
	dMajor, _ := GetScale(NoteD, "major")
	have, err := dMajor.IntoPitches(buffer)
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
		_, err = scale.IntoPitches(buffer)
		if err != nil {
			b.Fatal(err)
		}
	}
}
