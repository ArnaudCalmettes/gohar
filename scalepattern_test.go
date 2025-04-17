package gohar

import (
	"fmt"
	"slices"
	"testing"

	. "github.com/ArnaudCalmettes/gohar/test/helpers"
)

func TestScalePatternPitches(t *testing.T) {
	t.Run("break", func(t *testing.T) {
		for pitch := range ScalePatternMajor.Pitches(0) {
			if pitch != 0 {
				break
			}
		}
		// nothing blew up: success! \o/
	})

	testCases := []struct {
		Name string
		ScalePattern
		Want []Pitch
	}{
		{
			"major", ScalePatternMajor,
			[]Pitch{PitchC, PitchD, PitchE, PitchF, PitchG, PitchA, PitchB},
		}, {
			"melodic minor", ScalePatternMelodicMinor,
			[]Pitch{PitchC, PitchD, PitchEFlat, PitchF, PitchG, PitchA, PitchB},
		}, {
			"harmonic minor", ScalePatternHarmonicMinor,
			[]Pitch{PitchC, PitchD, PitchEFlat, PitchF, PitchG, PitchAFlat, PitchB},
		}, {
			"harmonic major", ScalePatternHarmonicMajor,
			[]Pitch{PitchC, PitchD, PitchE, PitchF, PitchG, PitchAFlat, PitchB},
		}, {
			"double harmonic major", ScalePatternDoubleHarmonicMajor,
			[]Pitch{PitchC, PitchDFlat, PitchE, PitchF, PitchG, PitchAFlat, PitchB},
		},
	}
	for _, tc := range testCases {
		Expect(t,
			Equalf(
				tc.Want,
				slices.Collect(tc.ScalePattern.Pitches(PitchC)),
				"%s", tc.Name,
			),
		)
	}
}

func BenchmarkScalePatternPitches(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for p := range ScalePatternMajor.Pitches(0) {
			if p < 0 {
				b.Fail()
			}
		}
	}
}

func TestScalePatternIntervals(t *testing.T) {
	t.Run("break", func(t *testing.T) {
		for i := range ScalePatternDoubleHarmonicMajor.Intervals() {
			if i != IntUnisson {
				break
			}
		}
		// nothing blew up...
	})

	t.Run("major scale", func(t *testing.T) {
		Expect(t,
			Equal(
				[]Interval{
					IntUnisson, IntMajorSecond, IntMajorThird, IntPerfectFourth,
					IntPerfectFifth, IntMajorSixth, IntMajorSeventh,
				},
				slices.Collect(ScalePatternMajor.Intervals()),
			),
		)
	})
	t.Run("major pentatonic", func(t *testing.T) {
		Expect(t,
			Equal(
				[]Interval{
					IntUnisson, IntMajorSecond, IntMajorThird,
					IntPerfectFifth, IntMajorSixth,
				},
				slices.Collect(ScalePattern(0b1010010101).IntervalsWithDegrees([]int8{1, 2, 3, 5, 6})),
			),
		)
	})
}

func BenchmarkScalePatternIntervals(b *testing.B) {
	for b.Loop() {
		for interval := range ScalePatternMajor.Intervals() {
			if interval.ScaleDiff < 0 {
				b.Fatal()
			}
		}
	}
}

func TestScalePatternIntervalsWithDegrees(t *testing.T) {
	pentatonic := ScalePattern(0b001010010101)
	t.Run("nil", func(t *testing.T) {
		for range pentatonic.IntervalsWithDegrees(nil) {
			t.Fatal("this line shouldn't be reached")
		}
	})

	t.Run("invalid degrees", func(t *testing.T) {
		for range pentatonic.IntervalsWithDegrees([]int8{1, 2, 3, 4, 5, 6, 7}) {
			t.Fatal("this line shouldn't be reached")
		}
	})

	t.Run("nominal", func(t *testing.T) {
		have := slices.Collect(pentatonic.IntervalsWithDegrees([]int8{1, 2, 3, 5, 6}))
		Expect(t,
			Equal(
				[]Interval{
					IntUnisson, IntMajorSecond, IntMajorThird,
					IntPerfectFifth, IntMajorSixth,
				},
				have,
			),
		)
	})

	t.Run("break", func(t *testing.T) {
		for i := range pentatonic.IntervalsWithDegrees([]int8{1, 2, 3, 5, 6}) {
			if i != IntUnisson {
				break
			}
		}
		// ... and nothings's broken. That's good news!
	})
}

func BenchmarkScalePatternIntervalsWithDegrees(b *testing.B) {
	for b.Loop() {
		for interval := range ScalePatternMajor.IntervalsWithDegrees([]int8{1, 2, 3, 4, 5, 6, 7}) {
			if interval.ScaleDiff < 0 {
				b.Fatal()
			}
		}
	}

}

func TestScalePatternPitchClasses(t *testing.T) {
	t.Run("break", func(t *testing.T) {
		count := 0
		for pc := range ScalePatternMajor.PitchClasses(PitchClassF) {
			if pc != PitchClassF {
				break
			}
			count++
		}
		Expect(t, Equal(1, count))
	})

	testCases := []struct {
		ScalePattern
		Root Note
		Want []Note
	}{
		{
			ScalePatternMajor, NoteC,
			[]Note{
				NoteC,
				NoteD,
				NoteE,
				NoteF,
				NoteG,
				NoteA,
				NoteB,
			},
		},
		{
			ScalePatternMelodicMinor, NoteG,
			[]Note{
				NoteG,
				NoteA,
				NoteB.Flat(),
				NoteC.Octave(1),
				NoteD.Octave(1),
				NoteE.Octave(1),
				NoteF.Sharp().Octave(1),
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			Expect(t,
				Equal(
					tc.Want,
					slices.Collect(
						tc.ScalePattern.Notes(tc.Root),
					),
				),
			)
		})
	}
}

func BenchmarkScalePatternPitchClasses(b *testing.B) {
	for b.Loop() {
		for pc := range ScalePatternMajor.PitchClasses(PitchClassC) {
			if !pc.IsValid() {
				b.Fatal()
			}
		}
	}
}

func TestScalePatternPitchClassesWithDegrees(t *testing.T) {
	pentatonic := ScalePattern(0b001010010101)
	t.Run("break", func(t *testing.T) {
		count := 0
		for pc := range pentatonic.PitchClassesWithDegrees(PitchClassC, []int8{1, 2, 3, 5, 6}) {
			if pc != PitchClassC {
				break
			}
			count++
		}
		Expect(t, Equal(1, count))
	})
}

func TestScalePatternNotes(t *testing.T) {
	t.Run("break", func(t *testing.T) {
		count := 0
		for pc := range ScalePatternMajor.Notes(NoteF) {
			if pc != NoteF {
				break
			}
			count++
		}
		Expect(t, Equal(1, count))
	})
}
func BenchmarkScalePatternNotes(b *testing.B) {
	for b.Loop() {
		for pc := range ScalePatternMajor.Notes(NoteC) {
			if !pc.IsValid() {
				b.Fatal()
			}
		}
	}
}

func TestScalePatternNotesWithDegrees(t *testing.T) {
	pentatonic := ScalePattern(0b001010010101)
	t.Run("break", func(t *testing.T) {
		count := 0
		for pc := range pentatonic.NotesWithDegrees(NoteC, []int8{1, 2, 3, 5, 6}) {
			if pc != NoteC {
				break
			}
			count++
		}
		Expect(t, Equal(1, count))
	})
}

func TestScalePatternMode(t *testing.T) {
	isScalePattern := AsCheckFunc(func(want, have ScalePattern) error {
		return Equal(want, have)
	})
	isError := HasError[ScalePattern]
	testCases := []struct {
		Scale ScalePattern
		Mode  int
		Check CheckFunc[ScalePattern]
	}{
		{
			Scale: ScalePatternMajor,
			Mode:  1,
			Check: isScalePattern(ScalePatternMajor),
		},
		{
			Scale: ScalePatternMajor,
			Mode:  4,
			Check: isScalePattern(0b101011010101), // lydian
		},
		{
			Scale: ScalePatternMajor,
			Mode:  7,
			Check: isScalePattern(0b010101101011), // locrian
		},
		{
			Scale: ScalePatternMajor,
			Mode:  0,
			Check: isError(ErrInvalidDegree),
		},
		{
			Scale: ScalePatternMajor,
			Mode:  8,
			Check: isError(ErrInvalidDegree),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v %d", tc.Scale, tc.Mode), func(t *testing.T) {
			got, err := tc.Scale.Mode(tc.Mode)
			Expect(t, tc.Check(got, err))
		})
	}
}
