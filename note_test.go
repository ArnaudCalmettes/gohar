package gohar

import (
	"fmt"
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func noteEqual(want, got Note) error {
	if want != got {
		return fmt.Errorf("expected Note %q, got %q", want, got)
	}
	return nil
}

func TestNoteName(t *testing.T) {
	testCases := []struct {
		Note
		Want string
	}{
		{Note{}, ""},
		{NoteC, "C"},
		{NoteC.Flat().Octave(-2), "C" + AltFlat},
	}

	for _, tc := range testCases {
		Expect(t, Equal(tc.Want, tc.Name()))
	}
}

func TestNoteStringer(t *testing.T) {
	testCases := []struct {
		Note
		Want string
	}{
		{Note{}, "Note{}"},
		{NoteC.Natural(), "C0"},
		{NoteF.Sharp(), "F" + AltSharp + "0"},
		{NoteB.Flat(), "B" + AltFlat + "0"},
		{NoteC.Octave(3), "C3"},
		{NoteD.DoubleFlat().Octave(-1), "D" + AltDoubleFlat + "-1"},
		{NoteC.DoubleSharp().Octave(3), "C" + AltDoubleSharp + "3"},
		{NoteG.Sharp().Sharp().Sharp().Octave(-5), "G(+3)-5"},
	}

	for _, tc := range testCases {
		Expect(t, Equal(tc.Want, tc.String()))
	}
}

func TestNotePitch(t *testing.T) {
	Expect(t, ShouldPanic(func() {
		Note{Base: 'H'}.Pitch()
	}))

	isPitch := AsCheckFunc(pitchEqual)
	testCases := []struct {
		Note
		Check CheckFunc[Pitch]
	}{
		{NoteC, isPitch(0)},
		{NoteD.Sharp(), isPitch(3)},
		{NoteE.Flat(), isPitch(3)},
		{NoteF.Octave(1), isPitch(17)},
		{NoteG.Octave(-1), isPitch(-5)},
		{NoteA.Sharp().Octave(1), isPitch(22)},
		{NoteB.Flat().Octave(-1), isPitch(-2)},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.Note), func(t *testing.T) {
			got := tc.Note.Pitch()
			Expect(t, tc.Check(got, nil))
		})
	}
}

func TestNoteAlterations(t *testing.T) {
	Expect(t,
		Equal(0, NoteC.Alt),
		Equal(-1, NoteD.Flat().Alt),
		Equal(-2, NoteB.DoubleFlat().Alt),
		Equal(1, NoteF.Sharp().Alt),
		Equal(2, NoteG.DoubleSharp().Alt),
	)

	Expect(t,
		noteEqual(NoteC.Flat(), Note{'C', -1, -1}),
		noteEqual(NoteD.Flat(), Note{'D', -1, 0}),
		noteEqual(NoteC.DoubleFlat(), Note{'C', -2, -1}),
		noteEqual(NoteD.DoubleFlat(), Note{'D', -2, 0}),
		noteEqual(NoteB.Sharp(), Note{'B', 1, 1}),
		noteEqual(NoteF.Sharp(), Note{'F', 1, 0}),
		noteEqual(NoteB.DoubleSharp(), Note{'B', 2, 1}),
		noteEqual(NoteF.DoubleSharp(), Note{'F', 2, 0}),
	)
}

func TestNoteTranspose(t *testing.T) {
	Expect(t,
		noteEqual(NoteC, NoteC.Transpose(IntUnisson)),
		noteEqual(NoteD.Sharp(), NoteC.Transpose(IntAugmentedSecond)),
		noteEqual(NoteE.Flat(), NoteC.Transpose(IntMinorThird)),
		noteEqual(NoteB.Flat().Octave(1), NoteD.Transpose(IntMinorThirteenth)),
		noteEqual(NoteF.Octave(-1), NoteC.Transpose(IntPerfectFifth.Down())),
		noteEqual(NoteC.Octave(1), NoteE.Flat().Transpose(IntMajorSixth)),
		noteEqual(NoteC.Octave(-2), NoteC.Octave(-2).Transpose(IntUnisson)),
	)
}

func TestNoteIsEnharmonic(t *testing.T) {
	Expect(t,
		Equal(false, NoteC.IsEnharmonic(NoteD)),
		Equal(true, NoteC.IsEnharmonic(NoteC)),
		Equal(true, NoteC.IsEnharmonic(NoteD.DoubleFlat())),
	)
}

func TestNoteWithPitch(t *testing.T) {
	isNote := AsCheckFunc(noteEqual)
	isError := HasError[Note]
	testCases := []struct {
		Base byte
		Pitch
		Check CheckFunc[Note]
	}{
		{'C', 0, isNote(NoteC)},
		{'B', 0, isNote(NoteB.Sharp().Octave(0))},
		{'C', -1, isNote(NoteC.Flat().Octave(-1))},
		{'E', 3, isNote(NoteE.Flat().Octave(0))},
		{'C', -12, isNote(NoteC.Octave(-1))},
		{'D', -10, isNote(NoteD.Octave(-1))},
		{'C', -24, isNote(NoteC.Octave(-2))},
		{'D', -22, isNote(NoteD.Octave(-2))},
		{'H', 13, isError(ErrInvalidBaseNote)},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%c @ %d", tc.Base, tc.Pitch), func(t *testing.T) {
			Expect(t, tc.Check(NoteWithPitch(tc.Base, tc.Pitch)))
		})
	}
}

func TestFindClosestNote(t *testing.T) {
	Expect(t,
		noteEqual(NoteC, FindClosestNote(0)),
		noteEqual(NoteD.Flat(), FindClosestNote(1)),
		noteEqual(NoteD, FindClosestNote(2)),
		noteEqual(NoteE.Flat(), FindClosestNote(3)),
		noteEqual(NoteE, FindClosestNote(4)),
		noteEqual(NoteF, FindClosestNote(5)),
		noteEqual(NoteG.Flat(), FindClosestNote(6)),
		noteEqual(NoteG, FindClosestNote(7)),
		noteEqual(NoteA.Flat(), FindClosestNote(8)),
		noteEqual(NoteA, FindClosestNote(9)),
		noteEqual(NoteB.Flat(), FindClosestNote(10)),
		noteEqual(NoteB, FindClosestNote(11)),
		noteEqual(NoteG.Flat().Octave(1), FindClosestNote(18)),
		noteEqual(NoteD.Flat().Octave(-1), FindClosestNote(-11)),
		noteEqual(NoteC.Octave(-1), FindClosestNote(-12)),
	)
}
