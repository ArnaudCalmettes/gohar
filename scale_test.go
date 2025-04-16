package gohar

import (
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestScaleStringer(t *testing.T) {
	testCases := []struct {
		Scale
		WantString string
	}{
		{Scale{PitchClassC, ScalePatternMajor}, "Scale(C:101010110101)"},
		{Scale{PitchClassD.Sharp(), ScalePatternMelodicMinor}, "Scale(D♯:101010101101)"},
	}

	for _, tc := range testCases {
		Expect(t,
			Equalf(tc.WantString, tc.Scale.String(), "string"),
		)
	}
}
