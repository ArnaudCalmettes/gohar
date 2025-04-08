package gohar

import (
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestChordPrintStringer(t *testing.T) {
	Expect(t,
		Equal("ChordPrint(0)", ChordPattern(0).String()),
		Equal("ChordPrint(10010001)", ChordPatternMajor.String()),
	)
}

func TestChordPrintHasDegree(t *testing.T) {
	testCases := []struct {
		Chord  ChordPattern
		Degree Pitch
		Want   bool
	}{
		{Chord: 0b0, Degree: 0, Want: false},
		{Chord: 0b1, Degree: 0, Want: true},
		{Chord: 0b10, Degree: 1, Want: true},
		{Chord: 0b10, Degree: 0, Want: false},
		{
			Chord:  ChordPatternMajor,
			Degree: PitchDiffPerfectFifth,
			Want:   true,
		},
		{
			Chord:  ChordPatternMinor,
			Degree: PitchDiffMajorThird,
			Want:   false,
		},
	}

	for _, tc := range testCases {
		have := tc.Chord.HasDegree(tc.Degree)
		Expect(t,
			Equalf(tc.Want, have, "%s has %d", tc.Chord, tc.Degree),
		)
	}
}

func TestChordPrintHasAnyDegree(t *testing.T) {
	testCases := []struct {
		Chord   ChordPattern
		Degrees []Pitch
		Want    bool
	}{
		{Chord: 0b0, Degrees: []Pitch{0}, Want: false},
		{Chord: 0b1, Degrees: nil, Want: false},
		{
			Chord: ChordPatternMajor,
			Degrees: []Pitch{
				PitchDiffMinorThird,
				PitchDiffMajorThird,
			},
			Want: true,
		},
		{
			Chord: ChordPatternMajor,
			Degrees: []Pitch{
				PitchDiffMajorSecond,
				PitchDiffMinorThird,
			},
			Want: false,
		},
	}

	for _, tc := range testCases {
		have := tc.Chord.HasAnyDegree(tc.Degrees...)
		Expect(t,
			Equalf(tc.Want, have, "%s has any %v", tc.Chord, tc.Degrees),
		)
	}
}

func TestChordPrintHasAllDegrees(t *testing.T) {
	testCases := []struct {
		ChordPattern
		Degrees []Pitch
		Want    bool
	}{
		{0b0, []Pitch{0}, false},
		{0b1, nil, false},
		{0b1, []Pitch{0}, true},
		{0b10010001, []Pitch{0, 4, 7}, true},
		{0b10010001, []Pitch{0, 3, 7}, false},
	}

	for _, tc := range testCases {
		have := tc.HasAllDegrees(tc.Degrees...)
		Expect(t,
			Equalf(tc.Want, have, "%s has all %v", tc.ChordPattern, tc.Degrees),
		)
	}
}

func TestChordPrintContains(t *testing.T) {
	Expect(t,
		Equal(true, ChordPatternMajor7.Contains(ChordPatternMajor)),
		Equal(false, ChordPatternMinor.Contains(ChordPatternMajor7)),
	)
}

func TestChordPrintUnpackAsIntervalSlice(t *testing.T) {
	testCases := []struct {
		Name  string
		Chord ChordPattern
		Want  []Interval
	}{
		{
			Name:  "major",
			Chord: ChordPatternMajor,
			Want:  []Interval{IntUnisson, IntMajorThird, IntPerfectFifth},
		},
		{
			Name:  "minor",
			Chord: ChordPatternMinor,
			Want:  []Interval{IntUnisson, IntMinorThird, IntPerfectFifth},
		},
		{
			Name:  "dim",
			Chord: ChordPatternDiminished,
			Want:  []Interval{IntUnisson, IntMinorThird, IntDiminishedFifth},
		},
		{
			Name:  "aug",
			Chord: ChordPatternAugmented,
			Want:  []Interval{IntUnisson, IntMajorThird, IntAugmentedFifth},
		},
		{
			Name:  "maj7 (#11)",
			Chord: ChordPatternMajor7.Add(PitchDiffAugmentedEleventh),
			Want: []Interval{
				IntUnisson,
				IntMajorThird,
				IntPerfectFifth,
				IntMajorSeventh,
				IntAugmentedEleventh,
			},
		},
		{
			Name:  "maj7 (omit5)",
			Chord: ChordPatternMajor7.Omit(PitchDiffPerfectFifth),
			Want: []Interval{
				IntUnisson,
				IntMajorThird,
				IntMajorSeventh,
			},
		},
		{
			Name:  "7 (b10,b13)",
			Chord: ChordPattern7No5.Add(PitchDiffMinorThird).Add(PitchDiffAugmentedFifth),
			Want: []Interval{
				IntUnisson,
				IntMajorThird,
				IntMinorSeventh,
				IntMinorTenth,
				IntMinorThirteenth,
			},
		},
		{
			Name: "6/9 #11",
			Chord: ChordPatternMajor.
				Add(PitchDiffMajorSixth).
				Add(PitchDiffMajorSecond).
				Add(PitchDiffAugmentedFourth),
			Want: []Interval{
				IntUnisson,
				IntMajorThird,
				IntPerfectFifth,
				IntMajorSixth,
				IntMajorNinth,
				IntAugmentedEleventh,
			},
		},
		{
			Name:  "dim7 add14",
			Chord: ChordPatternDiminished7.Add(PitchDiffMajorSeventh),
			Want: []Interval{
				IntUnisson,
				IntMinorThird,
				IntDiminishedFifth,
				IntDiminishedSeventh,
				IntMajorFourteenth,
			},
		},
		{
			Name:  "9 sus4",
			Chord: ChordPatternSus4.Add(PitchDiffMinorSeventh).Add(PitchDiffMajorSecond),
			Want: []Interval{
				IntUnisson,
				IntPerfectFourth,
				IntPerfectFifth,
				IntMinorSeventh,
				IntMajorNinth,
			},
		},
		{
			Name: "7 omit 5",
			// omit the fifth twice and see what happens
			Chord: ChordPattern7No5.Omit(PitchDiffPerfectFifth),
			Want: []Interval{
				IntUnisson,
				IntMajorThird,
				IntMinorSeventh,
			},
		},
	}

	for _, tc := range testCases {
		Expect(t, Equalf(tc.Want, tc.Chord.Unpack().AsIntervals(), "case %q", tc.Name))
	}
}

func BenchmarkChordPrintAsIntervalSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out := ChordPatternMajor.
			Add(PitchDiffMajorSixth).
			Add(PitchDiffMajorNinth).
			Add(PitchDiffAugmentedFourth).
			AsIntervals()
		if out == nil {
			b.Fatal(out)
		}
	}
}

func TestChordPrintAsIntervalSlice(t *testing.T) {
	_, err := ChordPatternMajor.IntoIntervals(nil)
	Expect(t,
		IsError(ErrNilBuffer, err),
	)

	_, err = ChordPatternMajor.IntoIntervals(make([]Interval, 0))
	Expect(t,
		IsError(ErrBufferOverflow, err),
	)
}

func BenchmarkChordPrintAsIntervalSliceInto(b *testing.B) {
	out := make([]Interval, 0, 10)
	var err error
	for i := 0; i < b.N; i++ {
		out, err = ChordPatternMajor.
			Add(PitchDiffMajorSixth).
			Add(PitchDiffMajorNinth).
			Add(PitchDiffAugmentedFourth).
			IntoIntervals(out)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestSwap(t *testing.T) {
	c := ChordPattern(0b101)
	Expect(t,
		Equal(c, c.swap(0, 2)), // Both are set
		Equal(ChordPattern(0b110), c.swap(0, 1)),
	)
}
