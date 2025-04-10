package gohar

import (
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

// AsPitches returns the scale pattern as a new slice of pitches relative
// to given root.
//
// This method dynamically allocates a new slice of pitches.
// See [ScalePattern.IntoPitches] for one that doesn't.
func (s ScalePattern) AsPitches(root Pitch) []Pitch {
	ps, _ := s.IntoPitches(
		make([]Pitch, 0, s.CountNotes()),
		root,
	)
	return ps
}

// IntoPitches converts the scale pattern into pitches relative to given root
// and writes them into the target slice.
//
// ErrBufferOverflow is returned if the target slice doesn't have enough capacity.
func (s ScalePattern) IntoPitches(target []Pitch, root Pitch) ([]Pitch, error) {
	if err := checkOutputBuffer(target, s.CountNotes()); err != nil {
		return nil, err
	}
	target = target[:0]
	for i := range 12 {
		if s&(1<<i) != 0 {
			target = append(target, Pitch(i)+root)
		}
	}
	return target, nil
}

// AsIntervals converts the scale pattern into intervals relative to the tonic.
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
//	majorPentatonic.AsIntervals([]int8{1,2,3,5,6})
//
// This method dynamically allocates a new slice of intervals.
// See [ScalePattern.IntoIntervals] if you need control over memory allocation.
func (s ScalePattern) AsIntervals(degrees []int8) []Interval {
	is, _ := s.IntoIntervals(
		make([]Interval, 0, s.CountNotes()),
		degrees,
	)
	return is
}

var range12 = []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// IntoIntervals converts the scale pattern into intervals relative to the tonic
// and writes them into the target slice.
//
// ErrBufferOverflow is returned if the target slice doesn't have enough capacity.
func (s ScalePattern) IntoIntervals(target []Interval, degrees []int8) ([]Interval, error) {
	if err := checkOutputBuffer(target, s.CountNotes()); err != nil {
		return nil, err
	}
	if degrees == nil {
		degrees = range12[:s.CountNotes()]
	}
	if len(degrees) != s.CountNotes() {
		return nil, wrapErrorf(ErrInvalidDegree,
			"expected []int8 with length %d: got %d", s.CountNotes(), len(degrees),
		)
	}
	target = target[:0]
	pitches, err := s.IntoPitches(make([]Pitch, 0, 12), 0)
	for i, pitch := range pitches {
		target = append(target, Interval{degrees[i] - 1, pitch})
	}
	return target, err
}

// AsNotes applies the scale pattern to given root note.
// degrees can be nil and has the same meaning as in [ScalePattern.AsIntervals].
//
// This method dynamically allocates a slice of notes.
// See [ScalePattern.IntoNotes] if you need control over memory allocation.
func (s ScalePattern) AsNotes(root Note, degrees []int8) ([]Note, error) {
	return s.IntoNotes(make([]Note, 0, s.CountNotes()), root, degrees)
}

// IntoNotes applies the scale pattern to given root note and writes the results in the
// target Note slice.
//
// ErrBufferOverflow is returned if the target slice doesn't have enough capacity.
func (s ScalePattern) IntoNotes(target []Note, root Note, degrees []int8) ([]Note, error) {
	if err := checkOutputBuffer(target, s.CountNotes()); err != nil {
		return nil, err
	}
	target = target[:0]
	intervals, err := s.IntoIntervals(make([]Interval, 0, 12), degrees)
	for _, inter := range intervals {
		target = append(target, root.Transpose(inter))
	}
	return target, err
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
