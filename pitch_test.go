package gohar

import (
	"fmt"
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func pitchEqual(want, got Pitch) error {
	if want != got {
		return fmt.Errorf("expected Pitch %d, got %d", want, got)
	}
	return nil
}

func pitchDiffEqual(want, got Pitch) error {
	if want != got {
		return fmt.Errorf("expected Pitch %d, got %d", want, got)
	}
	return nil
}

func TestPitchNormalize(t *testing.T) {
	testCases := []struct {
		Name  string
		Input Pitch
		Want  Pitch
	}{
		{"zero", 0, 0},
		{"1oct below", -7, 5},
		{"2oct below", -19, 5},
		{"1oct above", 17, 5},
		{"2oct above", 29, 5},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			got := tc.Input.Normalize()
			Expect(t, Equalf(tc.Want, got, "%s", tc.Name))
		})
	}
}

func TestPitchAdd(t *testing.T) {
	testCases := []struct {
		Input    Pitch
		Interval Pitch
		Want     Pitch
	}{
		{PitchC, PitchDiffUnisson, PitchC},
		{PitchDFlat, PitchDiffHalfStep, PitchD},
		{PitchAFlat, PitchDiffMajorThird, PitchC + PitchDiffOctave},
	}

	for _, tc := range testCases {
		label := fmt.Sprintf("%d(%+d)", tc.Input, tc.Interval)
		t.Run(label, func(t *testing.T) {
			Expect(t, Equalf(tc.Want, tc.Input.Add(tc.Interval), "%s", label))
		})
	}
}
