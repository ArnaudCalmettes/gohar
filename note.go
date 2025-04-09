package gohar

import (
	"fmt"
)

// A Note is an abstract music note that can be converted to pitch.
//
// Multiple notes can correspond to the same pitch (e.g.: E# and F). Notes that have the same pitch but different
// names are called "enharmonic".
type Note struct {
	// Base is the base pitch class of the note, represented as a byte between 'A' and 'G'.
	Base byte
	// Alt is the note's alteration, represented as a pitch offset between -2 (double flat) and +2 (double sharp).
	Alt Pitch
	// Oct transposes the note up (when Oct > 0) or down (when Oct < 0) relative to the default octave (0).
	Oct int8
}

var (
	NoteA = Note{Base: 'A'}
	NoteB = Note{Base: 'B'}
	NoteC = Note{Base: 'C'}
	NoteD = Note{Base: 'D'}
	NoteE = Note{Base: 'E'}
	NoteF = Note{Base: 'F'}
	NoteG = Note{Base: 'G'}
)

// Sharp returns the note with an added sharp.
func (n Note) Sharp() Note {
	if n.Base == 'B' && n.Alt == 0 {
		n.Oct++
	}
	n.Alt++
	return n
}

// Flat returns the note with an added flat.
func (n Note) Flat() Note {
	if n.Base == 'C' && n.Alt == 0 {
		n.Oct--
	}
	n.Alt--
	return n
}

// DoubleSharp makes the note "double sharp".
func (n Note) DoubleSharp() Note {
	if n.Base == 'B' && n.Alt == 0 {
		n.Oct++
	}
	n.Alt = 2
	return n
}

// DoubleFlat makes the note "double flat".
func (n Note) DoubleFlat() Note {
	if n.Base == 'C' && n.Alt == 0 {
		n.Oct--
	}
	n.Alt = -2
	return n
}

// Natural resets any alterations on the note.
func (n Note) Natural() Note {
	n.Alt = 0
	return n
}

// Octave transposes the note to desired octave.
func (n Note) Octave(oct int8) Note {
	n.Oct = oct
	return n
}

// Transpose returns the note transposed by given interval
func (n Note) Transpose(i Interval) Note {
	n, _ = NoteWithPitch(
		moveBaseNote(n.Base, int(i.ScaleDiff)),
		n.Pitch()+i.PitchDiff,
	)
	return n
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
	var zero Note
	if n == zero {
		return ""
	}
	return fmt.Sprintf("%c%s", n.Base, altToString(n.Alt))
}

// String returns the full note as a string.
func (n Note) String() string {
	var zero Note
	if n == zero {
		return "Note{}"
	}
	return fmt.Sprintf("%c%s%d", n.Base, altToString(n.Alt), n.Oct)
}

// Get the note's pitch.
//
// This method panics if the note is malformed.
func (n Note) Pitch() Pitch {
	return basePitch(n.Base).
		Add(n.Alt).
		Normalize().
		Add(Pitch(n.Oct) * PitchDiffOctave)
}

func moveBaseNote(base byte, diff int) byte {
	idx := (int(base) - int('A') + diff) % 7
	if idx < 0 {
		idx += 7
	}
	return byte(idx) + 'A'
}

var basePitches = [7]Pitch{PitchA, PitchB, PitchC, PitchD, PitchE, PitchF, PitchG}

func basePitch(b byte) Pitch {
	if b < 'A' || 'G' < b {
		panic(wrapErrorf(ErrInvalidBaseNote, "%c", b))
	}
	return basePitches[int(b-'A')]
}

func altToString(alt Pitch) string {
	switch alt {
	case -2:
		return AltDoubleFlat
	case -1:
		return AltFlat
	case 0:
		return ""
	case 1:
		return AltSharp
	case 2:
		return AltDoubleSharp
	default:
		return fmt.Sprintf("(%+d)", alt)
	}
}

// NoteWithPitch builds a note with given base and any
// octaves and alterations needed so that the note has
// given pitch.
func NoteWithPitch(base byte, pitch Pitch) (Note, error) {
	if base < 'A' || 'G' < base {
		return Note{}, wrapErrorf(ErrInvalidBaseNote, "%c", base)
	}
	note := Note{
		Base: base,
	}
	if pitch >= 0 {
		note.Oct = int8(pitch) / 12
	} else {
		note.Oct = int8(pitch+1)/12 - 1
	}

	diff := (pitch - note.Pitch()) % 12
	if diff < -6 {
		diff += 12
	}
	if diff > 6 {
		diff -= 12
	}
	note.Alt = diff
	return note, nil
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
