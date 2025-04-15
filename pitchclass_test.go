package gohar

import (
	"fmt"
	"slices"
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestNewPitchClass(t *testing.T) {
	isError := HasError[PitchClass]
	isPitchClass := AsCheckFunc(pitchClassEqual)
	testCases := []struct {
		Base  byte
		Alt   Pitch
		Check CheckFunc[PitchClass]
	}{
		{
			'#', 0,
			isError(ErrInvalidPitchClass),
		},
		{
			'A', 3,
			isError(ErrInvalidAlteration),
		},
		{
			'B', 0,
			isPitchClass(PitchClassB),
		},
		{
			'C', 1,
			isPitchClass(PitchClassC.Sharp()),
		},
		{
			'D', -1,
			isPitchClass(PitchClassD.Flat()),
		},
		{
			'E', -2,
			isPitchClass(PitchClassE.DoubleFlat()),
		},
		{
			'F', +2,
			isPitchClass(PitchClassF.DoubleSharp()),
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			have, err := NewPitchClassFromChar(tc.Base, tc.Alt)
			Expect(t, tc.Check(have, err))
		})
	}
}

func TestDefaultPitchClass(t *testing.T) {
	testCases := []struct {
		Pitch
		Want PitchClass
	}{
		{0, PitchClassC},
		{12, PitchClassC},
		{-12, PitchClassC},
		{24, PitchClassC},
		{-24, PitchClassC},
		{48, PitchClassC},
		{-48, PitchClassC},
		{6, PitchClassG.Flat()},
		{-6, PitchClassG.Flat()},
		{18, PitchClassG.Flat()},
		{-18, PitchClassG.Flat()},
		{30, PitchClassG.Flat()},
		{-30, PitchClassG.Flat()},
		{42, PitchClassG.Flat()},
		{-42, PitchClassG.Flat()},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			Expect(t,
				pitchClassEqual(
					tc.Want,
					DefaultPitchClass(tc.Pitch),
				),
			)
		})
	}
}

func TestPitchClassNaming(t *testing.T) {
	testCases := []struct {
		Expr PitchClass
		Want PitchClass
	}{
		{PitchClassA.Sharp(), 5 | PitchClassSharp},
		{PitchClassB.Flat(), 6 | PitchClassFlat},
		{PitchClassC.DoubleSharp(), 0 | PitchClassDoubleSharp},
		{PitchClassD.DoubleFlat(), 1 | PitchClassDoubleFlat},
	}
	for _, tc := range testCases {
		Expect(t,
			pitchClassEqual(tc.Want, tc.Expr),
		)
	}
}

func TestPitchClassProperties(t *testing.T) {
	Expect(t,
		Equal('D', PitchClassD.BaseName()),
		Equal(1, PitchClassD.Base()),
		Equal(-1, PitchClassE.Flat().Alt()),
		Equal(-2, PitchClassB.DoubleFlat().Alt()),
	)
}

func TestPitchClassIsEnharmonic(t *testing.T) {
	Expect(t,
		Equal(false, PitchClassC.IsEnharmonic(PitchClassC.Flat())),
		Equal(true, PitchClassC.Sharp().IsEnharmonic(PitchClassD.Flat())),
	)
}

func TestPitchClassPitch(t *testing.T) {
	testCases := []struct {
		PitchClass
		Octave int8
		Want   Pitch
	}{
		{PitchClassC, 0, 0},
		{PitchClassD, 1, 14},
		{PitchClassE.Flat(), 2, 27},
		{PitchClassF.Sharp(), 3, 42},
		{PitchClassG.DoubleFlat(), -1, -7},
		{PitchClassA.DoubleSharp(), -2, -13},
		{PitchClassB, -3, -25},
	}
	for _, tc := range testCases {
		Expect(t,
			Equal(tc.Want, tc.PitchClass.Pitch(tc.Octave)),
		)
	}
}

func TestPitchClassTranspose(t *testing.T) {
	testCases := []struct {
		PitchClass
		Interval
		Want PitchClass
	}{
		{
			PitchClassC, IntUnisson,
			PitchClassC,
		},
		{
			PitchClassC, IntAugmentedSecond,
			PitchClassD.Sharp(),
		},
		{
			PitchClassC, IntMinorThird,
			PitchClassE.Flat(),
		},
		{
			PitchClassD, IntMinorThirteenth,
			PitchClassB.Flat(),
		},
		{
			PitchClassC, IntPerfectFifth.Down(),
			PitchClassF,
		},
		{
			PitchClassE.Flat(), IntMajorSixth,
			PitchClassC,
		},
	}
	for _, tc := range testCases {
		label := fmt.Sprintf("%s + (%d,%d)", tc.PitchClass, tc.Interval.ScaleDiff, tc.Interval.PitchDiff)
		t.Run(label, func(t *testing.T) {
			Expect(t,
				pitchClassEqual(
					tc.Want,
					tc.PitchClass.Transpose(tc.Interval),
				),
			)
		})
	}
}

func TestPitchClassPitches(t *testing.T) {
	testCases := []struct {
		PitchClass
		From Pitch
		To   Pitch
		Want []Pitch
	}{
		{
			PitchClassC, -12, 12,
			[]Pitch{-12, 0, 12},
		},
		{
			PitchClassD.Sharp(), -24, 24,
			[]Pitch{-21, -9, 3, 15},
		},
		{
			PitchClassE.Flat(), 7, 36,
			[]Pitch{15, 27},
		},
		{
			PitchClassF, 12, -12,
			nil,
		},
	}
	for _, tc := range testCases {
		Expect(t,
			Equal(tc.Want, slices.Collect(
				tc.PitchClass.Pitches(tc.From, tc.To),
			)),
		)
	}
}

func TestPitchClassString(t *testing.T) {
	testCases := []struct {
		PitchClass
		Want string
	}{
		{0, "<invalid>"},
		{PitchClassC, "C"},
		{PitchClassD.Sharp(), "D" + AltSharp},
		{PitchClassE.Flat(), "E" + AltFlat},
		{PitchClassF.DoubleSharp(), "F" + AltDoubleSharp},
		{PitchClassB.DoubleFlat(), "B" + AltDoubleFlat},
	}

	for _, tc := range testCases {
		t.Run(tc.PitchClass.String(), func(t *testing.T) {
			Expect(t,
				Equal(tc.Want, tc.PitchClass.String()),
			)
		})
	}
}

func pitchClassEqual(want, have PitchClass) error {
	if want != have {
		return fmt.Errorf("want %0x, have %0x", uint8(want), uint8(have))
	}
	return nil
}
