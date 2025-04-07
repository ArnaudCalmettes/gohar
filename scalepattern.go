package gohar

import "math/bits"

// A Scale pattern is a bitwise representation of an octave.
// Each bit maps to a pitch relative to the root of the scale.
//
// For example the major scale (e.g. C D E F G A B) corresponds to the
// corresponding pattern:
//
//	pattern:  101010110101
//	notes:    B A G FE D C
type ScalePattern uint16

// AsPitches returns the scale pattern as a new slice of pitches relative
// to given root.
func (s ScalePattern) AsPitches(root Pitch) []Pitch {
	ps, _ := s.IntoPitches(
		make([]Pitch, 0, s.CountNotes()),
		root,
	)
	return ps
}

// IntoPitches converts the scale pattern into pitches relative to given root
// and writes them into the target slice.
// ErrBufferOverflow is returned if the target slice doesn't have enough capacity.
func (s ScalePattern) IntoPitches(target []Pitch, root Pitch) ([]Pitch, error) {
	if err := CheckOutputBuffer(target, s.CountNotes()); err != nil {
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
//	notes:          C        D          E          G            A
//	intervals:      unisson, major 2nd, major 3rd, perfect 5th, major 6th
//	degrees: []int8{0,       1,         2,         4,           5}
func (s ScalePattern) AsIntervals(degrees []int8) []Interval {
	is, _ := s.IntoIntervals(
		make([]Interval, 0, s.CountNotes()),
		nil,
	)
	return is
}

var range12 = []int8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// IntoIntervals converts the scale pattern into intervals relative to the tonic
// and writes them into the target slice.
// ErrBufferOverflow is returned if the target slice doesn't have enough capacity.
func (s ScalePattern) IntoIntervals(target []Interval, degrees []int8) ([]Interval, error) {
	if err := CheckOutputBuffer(target, s.CountNotes()); err != nil {
		return nil, err
	}
	if degrees == nil {
		degrees = range12[:s.CountNotes()]
	}
	target = target[:0]
	var d int8
	for i := 0; i < 12; i++ {
		if s&(1<<i) != 0 {
			target = append(target, Interval{degrees[d], Pitch(i)})
			d++
		}
	}
	return target, nil
}

// CountNotes returns the number of notes within the ScalePattern.
func (s ScalePattern) CountNotes() int {
	return bits.OnesCount16(uint16(s))
}

// Mode computes the n-th mode of the ScalePattern.
// n is expected to be between 1 and s.CountNotes() included, otherwise ScalePattern(0) is returned.
func (s ScalePattern) Mode(n int) ScalePattern {
	const mask = 0b0000111111111111 // 12 lowest bits
	var offset int
	if n < 1 || n > s.CountNotes() {
		return 0
	}
	for d := n; d > 1; d-- {
		offset = wrap(offset+1, 12)
		for s&(1<<offset) == 0 {
			offset = wrap(offset+1, 12)
		}
	}
	return s>>offset | (s<<(12-offset))&mask
}

func wrap(n, mod int) int {
	n = n % mod
	if n < 0 {
		n += mod
	}
	return n
}
