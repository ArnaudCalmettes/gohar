package gohar

import (
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
	Root PitchClass
	// Pattern is the scale's pattern.
	Pattern ScalePattern
}

// String returns a string representation of the scale.
func (s Scale) String() string {
	return fmt.Sprintf("Scale(%s:%012b)", s.Root, s.Pattern)
}
