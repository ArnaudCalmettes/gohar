package gohar

import (
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestParseNote(t *testing.T) {
	isError := HasError[Note]
	isNote := func(pc PitchClass, oct int8) CheckFunc[Note] {
		return AsCheckFunc(noteEqual)(Note{pc, oct})
	}

	testCases := []struct {
		Input string
		Check CheckFunc[Note]
	}{
		{"", isError(ErrCannotParseNote)},
		{"A", isNote(PitchClassA, 0)},
		{"B3", isNote(PitchClassB, 3)},
		{"G-2", isNote(PitchClassG, -2)},
		{"F#", isNote(PitchClassF.Sharp(), 0)},
		{"F#4", isNote(PitchClassF.Sharp(), 4)},
		{"Fb-1", isNote(PitchClassF.Flat(), -1)},
		{"G𝄫-2", isNote(PitchClassG.DoubleFlat(), -2)},
		{"A##+3", isNote(PitchClassA.DoubleSharp(), 3)},
		{"ebb", isNote(PitchClassE.DoubleFlat(), 0)},
	}

	for _, tc := range testCases {
		got, err := ParseNote(tc.Input)
		Expect(t, tc.Check(got, err))
	}
}

func TestParsePitch(t *testing.T) {
	isError := HasError[Pitch]
	isPitch := AsCheckFunc(pitchEqual)

	testCases := []struct {
		Input string
		Check CheckFunc[Pitch]
	}{
		{"C", isPitch(PitchC)},
		{"D#", isPitch(PitchDSharp)},
		{"Eb", isPitch(PitchEFlat)},
		{"F##", isPitch(PitchFDoubleSharp)},
		{"Gbb", isPitch(PitchGDoubleFlat)},
		{"A♯0", isPitch(PitchASharp)},
		{"B♭+0", isPitch(PitchBFlat)},
		{"C♮+1", isPitch(PitchC + 12)},
		{"D𝄫1", isPitch(PitchDDoubleFlat + 12)},
		{"E𝄪-2", isPitch(PitchEDoubleSharp - 24)},
		{"", isError(ErrCannotParseNote)},
		{"X#", isError(ErrCannotParseNote)},
	}

	for _, tc := range testCases {
		got, err := ParsePitch(tc.Input)
		Expect(t, tc.Check(got, err))
	}
}

func TestParseAlteration(t *testing.T) {
	isError := HasError[Pitch]
	isInterval := AsCheckFunc(pitchDiffEqual)

	testCases := []struct {
		Input string
		Check CheckFunc[Pitch]
	}{
		{"rubbish", isError(ErrUnknownAlteration)},
		{"", isInterval(0)},
		{"n", isInterval(0)},
		{"#", isInterval(1)},
		{"b", isInterval(-1)},
		{"##", isInterval(2)},
		{"bb", isInterval(-2)},
		{AltSharp, isInterval(1)},
		{AltFlat, isInterval(-1)},
		{AltNatural, isInterval(0)},
		{AltDoubleSharp, isInterval(+2)},
		{AltDoubleFlat, isInterval(-2)},
	}

	for _, tc := range testCases {
		got, err := ParseAlteration(tc.Input)
		Expect(t, tc.Check(got, err))
	}
}
