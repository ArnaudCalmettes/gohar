package gohar

import (
	"fmt"
	"math/bits"
)

// A chord pattern is represented bitwise as a 24-bit number.
// Each bit corresponds to a discrete pitch, within two 12-pitch octaves.
// When unpacked, a chord is preferably laid out using the first octave for the
// base chord (triad or tetrad) and the second octave for extensions (9th and above).
//
// This representation keeps the chord in a format that is efficiently computable,
// especially for set operations such as testing whether a chord contains certain
// degrees or extensions, whilst remaining close to most consistent naming conventions.
type ChordPattern uint32

const (
	// Major triad (C: C E G)
	ChordPatternMajor = ChordPattern(0b000010010001)
	// Minor triad (Cm: C Eb G)
	ChordPatternMinor = ChordPattern(0b000010001001)
	// Diminished triad (C dim: C Eb Gb)
	ChordPatternDiminished = ChordPattern(0b000001001001)
	// Augmented triad (C aug: C E G#)
	ChordPatternAugmented = ChordPattern(0b000100010001)
	// sus4 chord (Csus4: C F G)
	ChordPatternSus4 = ChordPattern(0b000010100001)
	// Major 7 chord (CMaj7: C E G B)
	ChordPatternMajor7 = ChordPattern(0b100010010001)
	// Major 7 (omit 5) chord (CMaj7 No5: C E B)
	ChordPatternMajor7No5 = ChordPattern(0b100000010001)
	// Dominant 7 chord (C7: C E G Bb)
	ChordPattern7 = ChordPattern(0b010010010001)
	// Dominant 7 (omit5) chord (C7 No5: C E Bb)
	ChordPattern7No5 = ChordPattern(0b010000010001)
	// Minor 7 chord (Cm7: C Eb G Bb)
	ChordPatternMinor7 = ChordPattern(0b010010001001)
	// Minor 7 No5 chord (Cm7 No5: C Eb Bb)
	ChordPatternMinor7No5 = ChordPattern(0b010000001001)
	// Minor 7 b5 chord (Cm7 b5: C Eb Gb Bb)
	ChordPatternMinor7Flat5 = ChordPattern(0b010001001001)
	// Diminished 7 chord (Cdim7: C Eb Gb Bbb)
	ChordPatternDiminished7 = ChordPattern(0b001001001001)
)

// String returns a string representation of the chord.
func (c ChordPattern) String() string {
	return fmt.Sprintf("ChordPrint(%b)", c)
}

// CountNotes counts the notes in the chord.
func (c ChordPattern) CountNotes() int {
	return bits.OnesCount32(uint32(c))
}

// Add adds a pitch to the chord. This has no effect if
// the pitch is already present.
func (c ChordPattern) Add(p Pitch) ChordPattern {
	return c | 1<<int(p)
}

// Omit removes a pitch from the chord. This has no effect
// if the pitch is already absent.
func (c ChordPattern) Omit(p Pitch) ChordPattern {
	if !c.HasDegree(p) {
		return c
	}
	return c ^ 1<<int(p)
}

// Contains return if other is a subset of the current chord pattern.
func (c ChordPattern) Contains(o ChordPattern) bool {
	return c&o == o
}

// HasDegree returns true if the chord contains given degree as a pitch.
func (c ChordPattern) HasDegree(p Pitch) bool {
	return c&(1<<int(p)) != 0
}

// HasAnyDegree returns true if the chord contains any of the pitches.
func (c ChordPattern) HasAnyDegree(pitches ...Pitch) bool {
	for _, p := range pitches {
		if c.HasDegree(p) {
			return true
		}
	}
	return false
}

// HasAllDegree returns true if the chord contains all of the pitches.
func (c ChordPattern) HasAllDegrees(pitches ...Pitch) bool {
	if len(pitches) == 0 {
		return false
	}
	for _, p := range pitches {
		if !c.HasDegree(p) {
			return false
		}
	}
	return true
}

func (c ChordPattern) AsIntervals() []Interval {
	chord, _ := c.IntoIntervals(make([]Interval, 0, c.CountNotes()))
	return chord
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

func (c ChordPattern) IntoIntervals(out []Interval) ([]Interval, error) {
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
	if c.Contains(ChordPatternDiminished7) {
		// Not a major 6th
		chord[3] = IntDiminishedSeventh
	}

	// Altered chord: #9 -> b10
	if c.Contains(ChordPattern7No5) && c.HasAllDegrees(PitchDiffAugmentedNinth, PitchDiffMinorThirteenth) {
		chord[3] = IntMinorTenth
	}
	return chord, nil
}

func (c ChordPattern) Unpack() ChordPattern {
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
	if c.Contains(ChordPatternDiminished7) {
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
	if c.Contains(ChordPattern7No5) && c.HasDegree(PitchDiffAugmentedNinth) {
		c = c.moveUp(PitchDiffAugmentedFifth)
	}
	return c
}

// move a degree up an octave (make it an extension)
func (c ChordPattern) moveUp(degree Pitch) ChordPattern {
	if c.HasDegree(degree) {
		return c.swap(degree, degree+PitchDiffOctave)
	}
	return c
}

// swap two degrees in the ChordPrint
func (c ChordPattern) swap(a, b Pitch) ChordPattern {
	mask := ChordPattern(1<<int(a) | 1<<int(b))
	if bits.OnesCount32(uint32(c&mask)) != 1 {
		return c
	}
	return c ^ ChordPattern(mask)
}
