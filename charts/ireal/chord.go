package ireal

import (
	"errors"
	"fmt"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// A ReadChord is a chord symbol read: its root and bass as written,
// sharp or flat, and the chord it names.
type ReadChord struct {
	Root    naming.SpelledNote
	Pattern harmony.ChordPattern
	Bass    naming.SpelledNote
	HasBass bool
}

// Chord anchors the pattern on its root, for whoever does not care how
// the root is spelled.
func (r ReadChord) Chord() harmony.Chord {
	return harmony.Chord{Root: r.Root.Class(), Pattern: r.Pattern}
}

// ErrInvisibleRoot is returned for a chord whose root is "W": it takes
// the root of the chord before it, which only the reader of the whole
// chart knows.
var ErrInvisibleRoot = errors.New("ireal: invisible root")

// Read reads a chord symbol into a chord.
//
// It fails on a quality neither in the app's list nor, when typed
// freely, recognisable once rewritten in the app's spelling: better an
// error to report than a chord guessed.
func (c ChordSymbol) Read() (ReadChord, error) {
	if c.Root == "W" {
		return ReadChord{}, ErrInvisibleRoot
	}
	root, ok := spell(c.Root)
	if !ok {
		return ReadChord{}, fmt.Errorf("ireal: root %q", c.Root)
	}
	o, ok := offsets(c.Quality, c.Custom)
	if !ok {
		return ReadChord{}, fmt.Errorf("ireal: quality %q", c.Quality)
	}
	st := make([]harmony.Semitones, len(o))
	for i, n := range o {
		st[i] = harmony.Semitones(n)
	}
	p, err := harmony.NewChordPattern(st...)
	if err != nil {
		return ReadChord{}, err
	}
	r := ReadChord{Root: root, Pattern: p}
	if c.Bass != "" {
		if r.Bass, ok = spell(c.Bass); !ok {
			return ReadChord{}, fmt.Errorf("ireal: bass %q", c.Bass)
		}
		r.HasBass = true
	}
	return r, nil
}

var letters = map[byte]naming.Letter{
	'C': naming.LetterC, 'D': naming.LetterD, 'E': naming.LetterE,
	'F': naming.LetterF, 'G': naming.LetterG, 'A': naming.LetterA,
	'B': naming.LetterB,
}

// spell reads a note as the app writes it: a letter and an optional #
// or b.
func spell(s string) (naming.SpelledNote, bool) {
	if s == "" {
		return naming.SpelledNote{}, false
	}
	l, ok := letters[s[0]]
	if !ok {
		return naming.SpelledNote{}, false
	}
	n := naming.SpelledNote{Letter: l}
	switch s[1:] {
	case "":
	case "#":
		n.Accidental = naming.SharpSign
	case "b":
		n.Accidental = naming.FlatSign
	default:
		return naming.SpelledNote{}, false
	}
	return n, true
}
