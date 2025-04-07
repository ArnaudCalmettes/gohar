package gohar

import (
	"errors"
	"fmt"
)

// A Scale is the association of a root Note and a ScalePattern.
// This pair can be extrapolated to an ordered set of notes or pitches within an octave.
//
// This representation allows for very efficient operation on scales, such as
// transposing, computing modes, or testing for inclusion.
//
// The whole structure fits within a 64bit word.
type Scale struct {
	// Root is the root note (or the tonic) of the scale.
	Root Note
	// Pattern is the scale's pattern.
	Pattern ScalePattern
}

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

// Deprecated: this should be handled by the Locale.
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

// Deprecated.
func GetScale(root Note, label string) (Scale, error) {
	if pattern, ok := ScalePatternMap[label]; !ok {
		return Scale{}, fmt.Errorf("%w %q", ErrUnknownScalePattern, label)
	} else {
		return Scale{root, pattern}, nil
	}
}

// String returns a string representation of the scale.
func (s Scale) String() string {
	return fmt.Sprintf("Scale(%s:%012b)", s.Root, s.Pattern)
}

func (s Scale) AsNotes() []Note {
	notes, _ := s.IntoNotes(
		make([]Note, 0, s.Pattern.CountNotes()),
	)
	return notes
}

func (s Scale) IntoNotes(target []Note) ([]Note, error) {
	intervals, err := s.Pattern.IntoIntervals(make([]Interval, 0, 12), nil)
	if err := CheckOutputBuffer(target, len(intervals)); err != nil {
		return nil, err
	}
	target = target[:0]
	for _, interval := range intervals {
		target = append(target, s.Root.Transpose(interval))
	}
	return target, err
}

func (s Scale) AsPitches() []Pitch {
	return s.Pattern.AsPitches(s.Root.Pitch())
}

func (s Scale) IntoPitches(target []Pitch) ([]Pitch, error) {
	return s.Pattern.IntoPitches(target, s.Root.Pitch())
}
