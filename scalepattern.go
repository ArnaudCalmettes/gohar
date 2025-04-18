package gohar

import (
	"iter"
	"math/bits"
)

// A Scale pattern is represented bitwise as a 12-bit number that maps an octave.
// Each bit corresponds to a pitch relative to the root of the scale.
//
// For example the major scale (e.g. C D E F G A B) has the following pattern:
//
//	pattern:  101010110101
//	notes:    B A G FE D C
//
// This representation allows for very efficient set operations such as testing
// if a note or pitch belongs to a scale, because most of these operations can be
// carried out as a single native bitwise instruction on most target architectures.
type ScalePattern uint16

const (
	ScalePatternMajor               ScalePattern = 0b101010110101 // C D E F G A B
	ScalePatternMelodicMinor        ScalePattern = 0b101010101101 // C D Eb F G A B
	ScalePatternHarmonicMinor       ScalePattern = 0b100110101101 // C D Eb F G Ab B
	ScalePatternHarmonicMajor       ScalePattern = 0b100110110101 // C D E F G Ab B
	ScalePatternDoubleHarmonicMajor ScalePattern = 0b100110110011 // C Db E F G Ab B
)

// CountNotes returns the number of notes within the ScalePattern.
func (s ScalePattern) CountNotes() int {
	return bits.OnesCount16(uint16(s))
}

// Pitches iterates over the pitches of the scale pattern relative to given root.
func (s ScalePattern) Pitches(root Pitch) iter.Seq[Pitch] {
	return func(yield func(Pitch) bool) {
		for i := 0; i < 12; i++ {
			if (s & (1 << i)) != 0 {
				if !yield(Pitch(i) + root) {
					return
				}
			}
		}
	}
}

// Intervals converts the scale pattern into intervals relative to the tonic.
// This method assumes that the scale is stepwise.
func (s ScalePattern) Intervals() iter.Seq[Interval] {
	return func(yield func(Interval) bool) {
		d := 0
		for p := range s.Pitches(0) {
			if !yield(Interval{int8(d), p}) {
				return
			}
			d++
		}
	}
}

// IntervalsWithDegrees converts the scale pattern into intervals relative to the tonic.
// If degrees == nil, the scale is assumed to have stepwise motion, which is suitable
// for most common scales in western music (heptatonic scales).
// Otherwise, degrees describe the absolute pitch class intervals to use.
//
// Eg. for a major pentatonic scale (degrees are the same for minor):
//
//	// notes:          C        D    E    G    A
//	// intervals:      unisson, 2nd, 3rd, 5th, 6th
//	// degrees: []int8{1,       2,   3,   5,   6}
//
//	majorPentatonic := ScalePattern(0b1010010101)
//	majorPentatonic.IntervalsWithDegrees([]int8{1,2,3,5,6})
func (s ScalePattern) IntervalsWithDegrees(degrees []int8) iter.Seq[Interval] {
	return func(yield func(Interval) bool) {
		if len(degrees) != s.CountNotes() {
			return
		}
		i := 0
		for p := range s.Pitches(0) {
			if !yield(Interval{degrees[i] - 1, p}) {
				return
			}
			i++
		}
	}
}

// PitchClasses converts the scale pattern into PitchClasses starting with given root.
// This method assumes that the scale is stepwise.
func (s ScalePattern) PitchClasses(root PitchClass) iter.Seq[PitchClass] {
	return func(yield func(PitchClass) bool) {
		for i := range s.Intervals() {
			if !yield(root.Transpose(i)) {
				return
			}
		}
	}
}

// PitchClassesWithDegrees converts the scale pattern into PitchClasses starting with given root.
// degrees have the same meaning as in [ScalePattern.IntervalsWithDegrees].
func (s ScalePattern) PitchClassesWithDegrees(root PitchClass, degrees []int8) iter.Seq[PitchClass] {
	return func(yield func(PitchClass) bool) {
		for i := range s.IntervalsWithDegrees(degrees) {
			if !yield(root.Transpose(i)) {
				return
			}
		}
	}
}

// Notes converts the scale pattern into notes starting on given root.
func (s ScalePattern) Notes(root Note) iter.Seq[Note] {
	return func(yield func(Note) bool) {
		for i := range s.Intervals() {
			if !yield(root.Transpose(i)) {
				return
			}
		}
	}
}

// NotesWithDegrees converts the scale pattern into notes starting on given root.
// degrees have the same meaning as in [ScalePattern.IntervalsWithDegrees].
func (s ScalePattern) NotesWithDegrees(root Note, degrees []int8) iter.Seq[Note] {
	return func(yield func(Note) bool) {
		for i := range s.IntervalsWithDegrees(degrees) {
			if !yield(root.Transpose(i)) {
				return
			}
		}
	}
}

// Mode computes the n-th mode of the ScalePattern.
//
// ErrInvalidDegree is returned if degree is not in the range [1;s.CountNotes()].
func (s ScalePattern) Mode(degree int) (ScalePattern, error) {
	const mask = 0b0000111111111111 // 12 lowest bits
	var offset int
	if degree < 1 || degree > s.CountNotes() {
		return 0, wrapErrorf(ErrInvalidDegree, "%d", degree)
	}
	for d := degree; d > 1; d-- {
		offset = (offset + 1) % 12
		for s&(1<<offset) == 0 {
			offset = (offset + 1) % 12
		}
	}
	return s>>offset | (s<<(12-offset))&mask, nil
}
