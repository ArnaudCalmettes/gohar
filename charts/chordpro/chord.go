package chordpro

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// A Chord is a chord symbol as musicians write it: "Bbmaj7", "Am7b5",
// "D7(b9)", "C7alt", "F/A", "N.C.".
type Chord struct {
	Root    naming.SpelledNote
	Pattern harmony.ChordPattern
	Bass    naming.SpelledNote
	HasBass bool

	NoChord  bool // N.C.
	Optional bool // between parentheses: the player's to add
}

// Harmony anchors the pattern on its root.
func (c Chord) Harmony() harmony.Chord {
	return harmony.Chord{Root: c.Root.Class(), Pattern: c.Pattern}
}

// same tells two chords that sound the same: a chord written again
// goes on, as in charts/ireal.
func (c Chord) same(d Chord) bool {
	return c.NoChord == d.NoChord && c.Root.Class() == d.Root.Class() &&
		c.Pattern == d.Pattern && c.HasBass == d.HasBass && c.Bass.Class() == d.Bass.Class()
}

// noChord is how a silence is written.
const noChord = "N.C."

// handSpellings turn what musicians write and charts/ireal does not
// read into what it does: the signs of Unicode, the 4 of a sus4, the
// sus of a sus2, the slash of a 6/9.
var handSpellings = strings.NewReplacer(
	"♭", "b", "♯", "#", "♮", "nat",
	"sus4", "sus", "sus2", "2", "6/9", "69",
)

// ReadChord reads a chord symbol, optional between parentheses.
func ReadChord(s string) (Chord, error) {
	var c Chord
	if inner, ok := strings.CutPrefix(s, "("); ok {
		if inner, ok = strings.CutSuffix(inner, ")"); !ok {
			return c, fmt.Errorf("chordpro: chord %q: unclosed parenthesis", s)
		}
		s, c.Optional = inner, true
	}
	if s == noChord {
		c.NoChord = true
		return c, nil
	}
	s = handSpellings.Replace(s)
	symbol, bass, slash := strings.Cut(s, "/")
	root, quality, ok := cutRoot(symbol)
	if !ok {
		return c, fmt.Errorf("chordpro: chord %q: root", s)
	}
	c.Root = root
	if c.Pattern, ok = readQuality(quality); !ok {
		return c, fmt.Errorf("chordpro: chord %q: quality %q", s, quality)
	}
	if slash {
		b, rest, ok := cutRoot(bass)
		if !ok || rest != "" {
			return c, fmt.Errorf("chordpro: chord %q: bass %q", s, bass)
		}
		c.Bass, c.HasBass = b, true
	}
	return c, nil
}

// extensions are the notes written in parentheses after a quality, as
// naming writes them: C7♭5(♭9), C7sus4(♭9), and the major seventh of a
// minor chord, which opens its parentheses: Cm(maj7,9).
var extensions = map[string]harmony.Semitones{
	"maj7": 11,
	"b9":   13, "9": 14, "#9": 15, "11": 17, "#11": 18, "b13": 20, "13": 21, "nat14": 23,
}

// readQuality reads a quality. The table of charts/ireal reads it whole
// when it can; otherwise the head before the parentheses is read alone,
// and the extensions between them added to it, in any order: the table
// knows 7b9sus, naming writes 7sus4(♭9), and both are the same chord.
func readQuality(q string) (harmony.ChordPattern, bool) {
	if p, ok := ireal.HandQuality(q); ok {
		return p, true
	}
	head, rest, ok := strings.Cut(q, "(")
	if !ok {
		return 0, false
	}
	list, ok := strings.CutSuffix(rest, ")")
	if !ok {
		return 0, false
	}
	p, ok := ireal.HandQuality(head)
	if !ok {
		return 0, false
	}
	offsets := slices.Collect(p.Offsets())
	for _, e := range strings.Split(list, ",") {
		n, ok := extensions[strings.TrimSpace(e)]
		if !ok {
			return 0, false
		}
		offsets = append(offsets, n)
	}
	p, err := harmony.NewChordPattern(offsets...)
	return p, err == nil
}

var letters = map[byte]naming.Letter{
	'C': naming.LetterC, 'D': naming.LetterD, 'E': naming.LetterE,
	'F': naming.LetterF, 'G': naming.LetterG, 'A': naming.LetterA,
	'B': naming.LetterB,
}

// cutRoot reads the note at the start of `s`, a letter and an optional
// b or #, and returns what follows it.
func cutRoot(s string) (naming.SpelledNote, string, bool) {
	if s == "" {
		return naming.SpelledNote{}, "", false
	}
	l, ok := letters[s[0]]
	if !ok {
		return naming.SpelledNote{}, "", false
	}
	n, rest := naming.SpelledNote{Letter: l}, s[1:]
	switch {
	case strings.HasPrefix(rest, "b"):
		n.Accidental, rest = naming.FlatSign, rest[1:]
	case strings.HasPrefix(rest, "#"):
		n.Accidental, rest = naming.SharpSign, rest[1:]
	}
	return n, rest, true
}

// style writes the qualities as En Harmonie does, the 6/9 with the slash
// of the keyboard.
var style = naming.ChordStyle{Slash: true}

// asciiSigns write the signs of a symbol with the characters of any
// keyboard, as the profile writes them.
var asciiSigns = strings.NewReplacer("♭", "b", "♯", "#", "♮", "nat")

// String writes the chord in ASCII, as the profile writes it: "Bbmaj7",
// "D7(b9,#11)", "F/A", "(Db7)". A pattern no symbol names is written
// with its offsets, which no reader takes back: better an error to see
// than a chord guessed.
func (c Chord) String() string {
	var s string
	switch q, ok := style.Symbol(c.Pattern); {
	case c.NoChord:
		s = noChord
	case !ok:
		s = naming.English.Name(c.Root, naming.ASCII) + "?" + c.Pattern.String()
	default:
		s = naming.English.Name(c.Root, naming.ASCII) + asciiSigns.Replace(q)
		if c.HasBass {
			s += "/" + naming.English.Name(c.Bass, naming.ASCII)
		}
	}
	if c.Optional {
		s = "(" + s + ")"
	}
	return s
}
