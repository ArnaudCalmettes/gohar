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

func TestPitchesWithClasses(t *testing.T) {
	type pair struct {
		Pitch
		PitchClass
	}
	testCases := []struct {
		From Pitch
		To   Pitch
		PCS  []PitchClass
		Want []pair
	}{
		{
			0, 12, nil,
			[]pair{},
		},
		{
			0, 12, []PitchClass{PitchClassC},
			[]pair{
				{0, PitchClassC},
				{12, PitchClassC},
			},
		},
		{
			0, 24, []PitchClass{PitchClassE.Flat(), PitchClassA},
			[]pair{
				{3, PitchClassE.Flat()},
				{9, PitchClassA},
				{15, PitchClassE.Flat()},
				{21, PitchClassA},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result := make([]pair, 0, 64)
			for p, pc := range PitchesWithClasses(tc.From, tc.To, tc.PCS) {
				result = append(result, pair{p, pc})
			}
			Expect(t,
				Equal(tc.Want, result),
			)
		})
	}
}
func BenchmarkPitchClassPitches(b *testing.B) {
	for i := 0; i < b.N; i++ {
		count := 0
		for range PitchClassC.Pitches(-100, 100) {
			count++
		}
		if count == 0 {
			b.Fatal()
		}
	}
}

func pitchClassEqual(want, have PitchClass) error {
	if want != have {
		return fmt.Errorf("want %0x, have %0x", uint8(want), uint8(have))
	}
	return nil
}

func BenchmarkPitchesWithClasses(b *testing.B) {
	benchCases := []*struct {
		Name         string
		PitchClasses []PitchClass
	}{
		{"default", nil},
		{"chromatic", defaultPitchClasses[:12]},
		{
			"heptatonic",
			[]PitchClass{
				PitchClassC,
				PitchClassD,
				PitchClassE,
				PitchClassF,
				PitchClassG,
				PitchClassA,
				PitchClassB,
			},
		},
		{
			"pentatonic",
			[]PitchClass{
				PitchClassC,
				PitchClassD,
				PitchClassE,
				PitchClassG,
				PitchClassA,
			},
		},
		{"tetrad", []PitchClass{PitchClassC, PitchClassE, PitchClassG, PitchClassB}},
		{"triad", []PitchClass{PitchClassC, PitchClassE, PitchClassG}},
		{"interval", []PitchClass{PitchClassC, PitchClassF.Sharp()}},
		{"note", []PitchClass{PitchClassC}},
	}
	for _, bc := range benchCases {
		b.Run(bc.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, pc := range PitchesWithClasses(-100, 100, bc.PitchClasses) {
					if !pc.IsValid() {
						b.Fatal(pc)
					}
				}
			}
		})
	}
}
