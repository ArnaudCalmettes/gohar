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
	Root Note
	// Pattern is the scale's pattern.
	Pattern ScalePattern
}

// String returns a string representation of the scale.
func (s Scale) String() string {
	return fmt.Sprintf("Scale(%s:%012b)", s.Root, s.Pattern)
}

// AsNotes returns the scale as a slice of notes.
//
// degrees have the same meaning as in [ScalePattern.AsIntervals].
//
// This method dynamically allocates a new slice of notes.
// Check [Scale.IntoNotes] if you need to keep control over memory allocation.
func (s Scale) AsNotes(degrees []int8) ([]Note, error) {
	return s.Pattern.AsNotes(s.Root, degrees)
}

// IntoNotes converts the scale into Notes and writes them into the target slice.
//
// degrees have the same meaning as in [ScalePattern.AsIntervals].
//
// ErrBufferOverflow is returned if the target slice doesn't have enough capacity.
func (s Scale) IntoNotes(target []Note, degrees []int8) ([]Note, error) {
	return s.Pattern.IntoNotes(target, s.Root, degrees)
}

// AsPitches returns the scale as a slice of pitches.
//
// This method dynamically allocates a new slice of notes.
// Check [Scale.IntoPitches] if you need to keep control over memory allocation.
func (s Scale) AsPitches() []Pitch {
	return s.Pattern.AsPitches(s.Root.Pitch())
}

// IntoPitches converts the scale into pitches and writes them into the target slice.
//
// ErrBufferOverflow is returned if the target slice doesn't have enough capacity.
func (s Scale) IntoPitches(target []Pitch) ([]Pitch, error) {
	return s.Pattern.IntoPitches(target, s.Root.Pitch())
}
