package gohar

import (
	"fmt"
)

// A Note is an abstract music note that can be manipulated before it is converted to pitch.
//
// Each note has:
// * a base name (A, B, C, D, E, F, G),
// * an optional alteration (number of sharps or flats) usually between -2 and +2.
// * an optional octave (for example C#1 is one octave above C#0). Default is 0.
//
// Multiple notes can correspond to the same pitch (e.g.: E# and F). Notes that have the same pitch but different
// names are called "enharmonic".
type Note struct {
	Base byte
	Alt  Pitch
	Oct  int8
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
	return NoteWithPitch(
		moveBaseNote(n.Base, int(i.ScaleDiff)),
		n.Pitch()+i.PitchDiff,
	)
}

// IsEnharmonic returns true if both notes have the same pitch.
func (n Note) IsEnharmonic(note Note) bool {
	return note.Pitch() == n.Pitch()
}

// IsHigherThan returns true if the current note is higher than the argument.
func (n Note) IsHigherThan(note Note) bool {
	return n.Pitch() > note.Pitch()
}

// Print the note's name regardless of the octave.
func (n Note) Name() string {
	var zero Note
	if n == zero {
		return ""
	}
	return fmt.Sprintf("%c%s", n.Base, altToString(n.Alt))
}

// Print the full note as a string.
func (n Note) String() string {
	var zero Note
	if n == zero {
		return "Note{}"
	}
	return fmt.Sprintf("%c%s%d", n.Base, altToString(n.Alt), n.Oct)
}

// Get the note's pitch.
//
// This method never fails: an invalid basename is treated as a C
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

func basePitch(b byte) Pitch {
	switch b {
	case 'A':
		return PitchA
	case 'B':
		return PitchB
	case 'C':
		return PitchC
	case 'D':
		return PitchD
	case 'E':
		return PitchE
	case 'F':
		return PitchF
	case 'G':
		return PitchG
	}
	return 0
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

func NoteWithPitch(name byte, pitch Pitch) Note {
	note := Note{
		Base: name,
		Oct:  int8(pitch) / 12,
	}
	if pitch < 0 {
		note.Oct--
	}
	diff := (pitch - note.Pitch()) % 12
	if diff < -6 {
		diff += 12
	}
	if diff > 6 {
		diff -= 12
	}
	note.Alt = diff
	return note
}

type FindOption uint

var (
	FindOptionPreferSharps FindOption = 0x01
)

func FindClosestNote(pitch Pitch, options ...FindOption) Note {
	var opt FindOption
	for _, o := range options {
		opt |= o
	}
	withOpt := func(o FindOption) bool {
		return opt&o == o
	}
	oct := int8(pitch) / 12
	if pitch < 0 {
		oct--
	}

	var note Note
	switch pitch.Normalize() {
	case 0:
		note = NoteC
	case 1:
		if withOpt(FindOptionPreferSharps) {
			note = NoteC.Sharp()
		} else {
			note = NoteD.Flat()
		}
	case 2:
		note = NoteD
	case 3:
		if withOpt(FindOptionPreferSharps) {
			note = NoteD.Sharp()
		} else {
			note = NoteE.Flat()
		}
	case 4:
		note = NoteE
	case 5:
		note = NoteF
	case 6:
		if withOpt(FindOptionPreferSharps) {
			note = NoteF.Sharp()
		} else {
			note = NoteG.Flat()
		}
	case 7:
		note = NoteG
	case 8:
		if withOpt(FindOptionPreferSharps) {
			note = NoteG.Sharp()
		} else {
			note = NoteA.Flat()
		}

	case 9:
		note = NoteA
	case 10:
		if withOpt(FindOptionPreferSharps) {
			note = NoteA.Sharp()
		} else {
			note = NoteB.Flat()
		}
	case 11:
		note = NoteB
	}
	return note.Octave(oct)
}

// Sorting order
type ByPitch []Note

func (b ByPitch) Len() int {
	return len(b)
}

func (b ByPitch) Less(i, j int) bool {
	return b[j].IsHigherThan(b[i])
}

func (b ByPitch) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
