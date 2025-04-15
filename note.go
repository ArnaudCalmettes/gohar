package gohar

import (
	"fmt"
)

// A Note is an abstract music note that can be converted to pitch.
//
// Multiple notes can correspond to the same pitch (e.g.: E# and F). Notes that have the same pitch but different
// names are called "enharmonic".
type Note struct {
	// PitchClass is the pitch class of the note
	PitchClass
	// Oct transposes the note up (when Oct > 0) or down (when Oct < 0) relative to the default octave (0).
	Oct int8
}

var (
	NoteA = Note{PitchClassA, 0}
	NoteB = Note{PitchClassB, 0}
	NoteC = Note{PitchClassC, 0}
	NoteD = Note{PitchClassD, 0}
	NoteE = Note{PitchClassE, 0}
	NoteF = Note{PitchClassF, 0}
	NoteG = Note{PitchClassG, 0}
)

// Sharp returns the note with an added sharp.
func (n Note) Sharp() Note {
	oct := n.Oct
	if n.PitchClass.Pitch(0) == PitchB {
		oct++
	}
	return Note{n.PitchClass.Sharp(), oct}
}

// Flat returns the note with an added flat.
func (n Note) Flat() Note {
	oct := n.Oct
	if n.PitchClass.Pitch(0) == PitchC {
		oct--
	}
	return Note{n.PitchClass.Flat(), oct}
}

// DoubleSharp makes the note "double sharp".
func (n Note) DoubleSharp() Note {
	oct := n.Oct
	if n.PitchClass.Pitch(0) >= PitchBFlat {
		oct++
	}
	return Note{n.PitchClass.DoubleSharp(), oct}
}

// DoubleFlat makes the note "double flat".
func (n Note) DoubleFlat() Note {
	oct := n.Oct
	if n.PitchClass.Pitch(0) <= PitchCSharp {
		oct--
	}
	return Note{n.PitchClass.DoubleFlat(), oct}
}

// Octave transposes the note to desired octave.
func (n Note) Octave(oct int8) Note {
	n.Oct = oct
	return n
}

// Transpose returns the note transposed by given interval
func (n Note) Transpose(i Interval) Note {
	pitch := n.Pitch() + i.PitchDiff
	return Note{
		n.PitchClass.Transpose(i),
		pitch.GetOctave(),
	}
}

// IsEnharmonic returns true if both notes have the same pitch.
func (n Note) IsEnharmonic(note Note) bool {
	return note.Pitch() == n.Pitch()
}

// IsHigherThan returns true if the current note is higher than the argument.
func (n Note) IsHigherThan(note Note) bool {
	return n.Pitch() > note.Pitch()
}

// Name returns the note's name regardless of the octave.
func (n Note) Name() string {
	return n.PitchClass.String()
}

// String returns the full note as a string.
func (n Note) String() string {
	return fmt.Sprintf("%s%d", n.PitchClass, n.Oct)
}

// Get the note's pitch.
//
// This method panics if the note is malformed.
func (n Note) Pitch() Pitch {
	return n.PitchClass.Pitch(n.Oct)
}

// NoteWithPitch builds a note with given base and any
// octaves and alterations needed so that the note has
// given pitch.
func NoteWithPitch(pc PitchClass, pitch Pitch) Note {
	return Note{pitchClassWithPitch(pc, pitch), pitch.GetOctave()}
}

var closestNote = [12]Note{
	NoteC,
	NoteD.Flat(),
	NoteD,
	NoteE.Flat(),
	NoteE,
	NoteF,
	NoteG.Flat(),
	NoteG,
	NoteA.Flat(),
	NoteA,
	NoteB.Flat(),
	NoteB,
}

// FindClosestNote returns the closest note with given pitch.
// When a pitch corresponds to an altered note ("black key"),
// it is always assumed to be the flattened note above it.
func FindClosestNote(pitch Pitch) Note {
	var oct int8
	if pitch >= 0 {
		oct = int8(pitch) / 12
	} else {
		oct = int8(pitch+1)/12 - 1
	}
	return closestNote[int(pitch.Normalize())].Octave((oct))
}
