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

func TestScaleAsNotes(t *testing.T) {
	testCases := []struct {
		Root    Note
		Pattern ScalePattern
		Want    []Note
	}{
		{
			NoteC, ScalePatternMajor,
			[]Note{NoteC, NoteD, NoteE, NoteF, NoteG, NoteA, NoteB},
		},
		{
			NoteD, ScalePatternMajor,
			[]Note{NoteD, NoteE, NoteF.Sharp(), NoteG, NoteA, NoteB, NoteC.Sharp().Octave(1)},
		},
		{
			NoteE, ScalePatternMelodicMinor,
			[]Note{NoteE, NoteF.Sharp(), NoteG, NoteA, NoteB, NoteC.Sharp().Octave(1), NoteD.Sharp().Octave(1)},
		},
		{
			NoteF, ScalePatternHarmonicMinor,
			[]Note{NoteF, NoteG, NoteA.Flat(), NoteB.Flat(), NoteC.Octave(1), NoteD.Flat().Octave(1), NoteE.Octave(1)},
		},
		{
			NoteC.Sharp(), ScalePatternHarmonicMajor,
			[]Note{NoteC.Sharp(), NoteD.Sharp(), NoteE.Sharp(), NoteF.Sharp(), NoteG.Sharp(), NoteA, NoteB.Sharp()},
		},
		{
			NoteA.Octave(-1), ScalePatternDoubleHarmonicMajor,
			[]Note{NoteA.Octave(-1), NoteB.Flat().Octave(-1), NoteC.Sharp(), NoteD, NoteE, NoteF, NoteG.Sharp()},
		},
	}

	for _, tc := range testCases {
		scale := Scale{tc.Root, tc.Pattern}
		have, _ := scale.AsNotes(nil)
		Expect(t, Equalf(
			fmt.Sprint(tc.Want),
			fmt.Sprint(have),
			"%s", scale,
		))
	}
}

func BenchmarkScaleAsNoteSlice(b *testing.B) {
	scale := Scale{NoteE.Flat(), ScalePatternHarmonicMinor}
	for i := 0; i < b.N; i++ {
		if _, err := scale.AsNotes(nil); err != nil {
			b.Fatal(err)
		}
	}
}

func TestScaleAsNoteSliceInto(t *testing.T) {
	cMajor := Scale{NoteC, ScalePatternMajor}
	_, err := cMajor.IntoNotes(nil, nil)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)
	_, err = cMajor.IntoNotes(make([]Note, 0), nil)
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkScaleAsNoteSliceInto(b *testing.B) {
	notes := make([]Note, 0, 12)
	scale := Scale{NoteE.Flat(), ScalePatternHarmonicMinor}
	for i := 0; i < b.N; i++ {
		_, err := scale.IntoNotes(notes, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestScaleAsPitchSlice(t *testing.T) {
	dMajor := Scale{NoteD, ScalePatternMajor}
	have := dMajor.AsPitches()
	Expect(t,
		Equal(
			[]Pitch{PitchD, PitchE, PitchFSharp, PitchG, PitchA, PitchB, PitchCSharp.Add(12)},
			have,
		),
	)
}

func BenchmarkScaleAsPitchSlice(b *testing.B) {
	scale := Scale{NoteD, ScalePatternMajor}
	for i := 0; i < b.N; i++ {
		if scale.AsPitches() == nil {
			b.Fatal()
		}
	}
}

func TestScaleAsPitchSliceInto(t *testing.T) {
	buffer := make([]Pitch, 0, 7)
	dMajor := Scale{NoteD, ScalePatternMajor}
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
	scale := Scale{NoteD, ScalePatternMajor}
	for i := 0; i < b.N; i++ {
		if _, err := scale.IntoPitches(buffer); err != nil {
			b.Fatal(err)
		}
	}
}
