package gohar

import (
	"errors"
	"fmt"
	"math/bits"
)

type Scale struct {
	Root    Note
	Pattern ScalePattern
}

type ScalePattern uint16

const (
	// C D E F G A B
	ScalePatternMajor ScalePattern = 0b101010110101

	// C D Eb F G A B
	ScalePatternMelodicMinor ScalePattern = 0b101010101101

	// C D Eb F G Ab B
	ScalePatternHarmonicMinor ScalePattern = 0b100110101101

	// C D E F G Ab B
	ScalePatternHarmonicMajor ScalePattern = 0b100110110101

	// C Db E F G Ab B
	ScalePatternDoubleHarmonicMajor ScalePattern = 0b100110110011
)

var (
	ScalePatternMap = map[string]ScalePattern{
		"major":                 ScalePatternMajor,
		"natural major":         ScalePatternMajor,
		"melodic minor":         ScalePatternMelodicMinor,
		"harmonic minor":        ScalePatternHarmonicMinor,
		"harmonic major":        ScalePatternHarmonicMajor,
		"double harmonic major": ScalePatternDoubleHarmonicMajor,
	}

	KnownScalePatterns = map[ScalePattern]string{
		ScalePatternMajor:               "major",
		ScalePatternMelodicMinor:        "melodic minor",
		ScalePatternHarmonicMinor:       "harmonic minor",
		ScalePatternHarmonicMajor:       "harmonic major",
		ScalePatternDoubleHarmonicMajor: "double harmonic major",
	}

	ErrUnknownScalePattern = errors.New("unknown scale pattern")
)

func GetScale(root Note, label string) (Scale, error) {
	if pattern, ok := ScalePatternMap[label]; !ok {
		return Scale{}, fmt.Errorf("%w %q", ErrUnknownScalePattern, label)
	} else {
		return Scale{root, pattern}, nil
	}
}

func (s Scale) String() string {
	return fmt.Sprintf("Scale(%s:%012b)", s.Root, s.Pattern)
}

func (s Scale) AsNoteSlice() []Note {
	notes, _ := s.AsNoteSliceInto(
		make([]Note, 0, s.Pattern.CountNotes()),
	)
	return notes
}

func (s Scale) AsNoteSliceInto(notes []Note) ([]Note, error) {
	buffer := make([]Interval, 0, 12)
	intervals, err := s.Pattern.AsIntervalSliceInto(buffer)
	if err := CheckOutputBuffer(notes, len(intervals)); err != nil {
		return nil, err
	}
	notes = notes[:0]
	for _, interval := range intervals {
		notes = append(notes, s.Root.Transpose(interval))
	}
	return notes, err
}

func (s Scale) AsPitchSlice() []Pitch {
	pitches := s.Pattern.AsPitchSlice()
	root := s.Root.Pitch()
	for i := range pitches {
		pitches[i] += root
	}
	return pitches
}

func (s Scale) AsPitchSliceInto(pitches []Pitch) ([]Pitch, error) {
	pitches, err := s.Pattern.AsPitchSliceInto(pitches)
	root := s.Root.Pitch()
	for i := range pitches {
		pitches[i] += root
	}
	return pitches, err
}

func (s ScalePattern) AsPitchSlice() []Pitch {
	ps, _ := s.AsPitchSliceInto(
		make([]Pitch, 0, s.CountNotes()),
	)
	return ps
}

func (s ScalePattern) AsPitchSliceInto(pitches []Pitch) ([]Pitch, error) {
	if err := CheckOutputBuffer(pitches, s.CountNotes()); err != nil {
		return nil, err
	}
	pitches = pitches[:0]
	for i := range 12 {
		if s&(1<<i) != 0 {
			pitches = append(pitches, Pitch(i))
		}
	}
	return pitches, nil
}

func (s ScalePattern) AsIntervalSlice() []Interval {
	is, _ := s.AsIntervalSliceInto(
		make([]Interval, 0, s.CountNotes()),
	)
	return is
}

func (s ScalePattern) AsIntervalSliceInto(intervals []Interval) ([]Interval, error) {
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
	degree = wrap(degree-1, s.CountNotes()) + 1
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
