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
func (s ScalePattern) AsIntervals() []Interval {
	is, _ := s.IntoIntervals(
		make([]Interval, 0, s.CountNotes()),
	)
	return is
}

func (s ScalePattern) IntoIntervals(intervals []Interval) ([]Interval, error) {
	if err := CheckOutputBuffer(intervals, s.CountNotes()); err != nil {
		return nil, err
	}
	intervals = intervals[:0]
	var degree int8
	for i := 0; i < 12; i++ {
		if s&(1<<i) != 0 {
			intervals = append(intervals, Interval{degree, Pitch(i)})
			degree++
		}
	}
	return intervals, nil
}

func (s ScalePattern) CountNotes() int {
	return bits.OnesCount16(uint16(s))
}

const octaveMask = 0b111111111111

func (s ScalePattern) Mode(degree int) ScalePattern {
	var offset int
	if degree < 1 || degree > s.CountNotes() {
		return 0
	}
	for d := degree; d > 1; d-- {
		offset = wrap(offset+1, 12)
		for s&(1<<offset) == 0 {
			offset = wrap(offset+1, 12)
		}
	}
	return s>>offset | (s<<(12-offset))&octaveMask
}

func wrap(n, mod int) int {
	n = n % mod
	if n < 0 {
		n += mod
	}
	return n
}
