package gohar

import (
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestLocaleSprintNote(t *testing.T) {
	isError := HasError[string]
	isString := AsCheckFunc(func(a, b string) error {
		return Equal(a, b)
	})
	testCases := []struct {
		Loc   *Locale
		Note  Note
		Check CheckFunc[string]
	}{
		{&LocaleFrench, Note{}, isError(ErrInvalidPitchClass)},
		{&LocaleFrench, NoteC.DoubleSharp().Sharp(), isError(ErrInvalidAlteration)},
		{&LocaleFrench, NoteC.Sharp(), isString("do" + AltSharp)},
		{&LocaleFrench, NoteA.Flat(), isString("la" + AltFlat)},
		{&LocaleFrench, NoteB.DoubleFlat().Octave(-2), isString("si" + AltDoubleFlat)},
	}

	for _, tc := range testCases {
		have, err := tc.Loc.NoteName(tc.Note)
		Expect(t, tc.Check(have, err))
	}
}
