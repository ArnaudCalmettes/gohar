package gohar

import (
	"fmt"
	"math/bits"
)

type ChordPrint uint32

const (
	ChordPrintMajor       = ChordPrint(0b000010010001)
	ChordPrintMinor       = ChordPrint(0b000010001001)
	ChordPrintDiminished  = ChordPrint(0b000001001001)
	ChordPrintAugmented   = ChordPrint(0b000100010001)
	ChordPrintSus4        = ChordPrint(0b000010100001)
	ChordPrintMajor7      = ChordPrint(0b100010010001)
	ChordPrintMajor7No5   = ChordPrint(0b100000010001)
	ChordPrint7           = ChordPrint(0b010010010001)
	ChordPrint7No5        = ChordPrint(0b010000010001)
	ChordPrintMinor7      = ChordPrint(0b010010001001)
	ChordPrintMinor7No5   = ChordPrint(0b010000001001)
	ChordPrintDiminished7 = ChordPrint(0b001001001001)
)

func (c ChordPrint) String() string {
	return fmt.Sprintf("ChordPrint(%b)", c)
}

func (c ChordPrint) CountNotes() int {
	return bits.OnesCount32(uint32(c))
}

func (c ChordPrint) Add(pitchDiff Pitch) ChordPrint {
	return c | 1<<int(pitchDiff)
}

func (c ChordPrint) Omit(pitchDiff Pitch) ChordPrint {
	if !c.HasDegree(pitchDiff) {
		return c
	}
	return c ^ 1<<int(pitchDiff)
}

func (c ChordPrint) Contains(other ChordPrint) bool {
	return c&other == other
}

func (c ChordPrint) HasDegree(pitchDiff Pitch) bool {
	return c&(1<<int(pitchDiff)) != 0
}

func (c ChordPrint) HasAnyDegree(pitchDiffs ...Pitch) bool {
	for _, d := range pitchDiffs {
		if c.HasDegree(d) {
			return true
		}
	}
	return false
}

func (c ChordPrint) HasAllDegrees(pitchDiffs ...Pitch) bool {
	if len(pitchDiffs) == 0 {
		return false
	}
	for _, d := range pitchDiffs {
		if !c.HasDegree(d) {
			return false
		}
	}
	return true
}

func (c ChordPrint) AsIntervalSlice() (chord []Interval) {
	chord, _ = c.AsIntervalSliceInto(
		make([]Interval, 0, c.CountNotes()),
	)
	return
}

var (
	sieve = []Interval{
		IntUnisson,
		IntMajorSecond,
		IntMinorThird,
		IntMajorThird,
		IntPerfectFourth,
		IntDiminishedFifth,
		IntPerfectFifth,
		IntAugmentedFifth,
		IntMajorSixth,
		IntMinorSeventh,
		IntMajorSeventh,
		IntMinorNinth,
		IntMajorNinth,
		IntAugmentedNinth,
		IntPerfectEleventh,
		IntAugmentedEleventh,
		IntMinorThirteenth,
		IntMajorThirteenth,
		IntMajorFourteenth,
	}
)

func (c ChordPrint) AsIntervalSliceInto(out []Interval) ([]Interval, error) {
	if err := CheckOutputBuffer(out, c.CountNotes()); err != nil {
		return nil, err
	}
	chord := out[:0]

	for _, interval := range sieve {
		if c.HasDegree(interval.PitchDiff) {
			chord = append(chord, interval)
		}
	}

	// Irregular cases
	if c.Contains(ChordPrintDiminished7) {
		// Not a major 6th
		chord[3] = IntDiminishedSeventh
	}

	// Altered chord: #9 -> b10
	if c.Contains(ChordPrint7No5) && c.HasAllDegrees(PitchDiffAugmentedNinth, PitchDiffMinorThirteenth) {
		chord[3] = IntMinorTenth
	}
	return chord, nil
}

func (c ChordPrint) Unpack() ChordPrint {
	// b2 always becomes b9
	c = c.moveUp(PitchDiffMinorSecond)
	// the third is present -> 2 and 4 become 9 and 11
	if c.HasAnyDegree(PitchDiffMajorThird, PitchDiffMinorThird) {
		c = c.moveUp(PitchDiffMajorSecond)
		c = c.moveUp(PitchDiffPerfectFourth)
		// #2 or b3 becomes #9 or b10
		if c.HasAllDegrees(PitchDiffMajorThird, PitchDiffMinorThird) {
			c = c.moveUp(PitchDiffMinorThird)
		}
	}
	// chord is sus4 -> 2 becomes 9
	if c.HasDegree(PitchDiffPerfectFourth) {
		c = c.moveUp(PitchDiffMajorSecond)
	}
	// major triad -> move #11 up (it's not a b5)
	if c.HasAllDegrees(PitchDiffMajorThird, PitchDiffPerfectFifth) {
		c = c.moveUp(PitchDiffAugmentedFourth)
	}
	// If we have a diminished tetrad, major 7th becomes a major 14th
	if c.Contains(ChordPrintDiminished7) {
		c = c.moveUp(PitchDiffMajorSeventh)
	}
	// Otherwise, when a 7th is present, 6 and b6 become 13 and b13
	if c.HasAnyDegree(PitchDiffMajorSeventh, PitchDiffMinorSeventh) {
		c = c.moveUp(PitchDiffMajorSixth)
		c = c.moveUp(PitchDiffMinorSixth)
		if c.HasDegree(PitchDiffMajorSeventh) {
			c = c.moveUp(PitchDiffDiminishedFifth)
		}
	}
	// Altered chord: #5 -> b13
	if c.Contains(ChordPrint7No5) && c.HasDegree(PitchDiffAugmentedNinth) {
		c = c.moveUp(PitchDiffAugmentedFifth)
	}
	return c
}

// move a degree up an octave (make it an extension)
func (c ChordPrint) moveUp(degree Pitch) ChordPrint {
	if c.HasDegree(degree) {
		return c.swap(degree, degree+PitchDiffOctave)
	}
	return c
}

// swap two degrees in the ChordPrint
func (c ChordPrint) swap(a, b Pitch) ChordPrint {
	mask := ChordPrint(1<<int(a) | 1<<int(b))
	if bits.OnesCount32(uint32(c&mask)) != 1 {
		return c
	}
	return c ^ ChordPrint(mask)
}
