package chart

import (
	"strings"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// A Symbol is a chord as a chart draws it, in three runs: the root and
// the kind of tetrad on the line, the altered fifth and the extensions
// raised, the bass back on the line. D7(♭13) is D7 and (♭13), Am7♭5 is
// Am7 and ♭5, D7/F♯ is D7 and /F♯.
type Symbol struct {
	Line, Raised, Bass string
}

// String writes the symbol on one line, for the recording: B♭7, Am7♭5,
// D7(♭13), D7/F♯, N.C.
func (s Symbol) String() string { return s.Line + s.Raised + s.Bass }

// SymbolOf cuts the chord `c` as a chart draws it, from the name
// naming gives: the raised run starts at the altered fifth or at the
// parentheses, whichever comes first.
//
// What stays on the line for now: a minor-major seventh, whose maj7 is
// the kind of tetrad though naming writes it in parentheses, Cm(maj7,9);
// the alt of a 7alt; the extensions that climb on the seventh, C9, C13,
// which naming writes in place of the 7.
func SymbolOf(c chordpro.Chord) Symbol {
	if c.NoChord {
		return Symbol{Line: "N.C."}
	}
	q, ok := naming.ChordStyle{}.Symbol(c.Pattern)
	if !ok {
		q = "?"
	}
	s := Symbol{Line: naming.English.Name(c.Root, naming.Signs)}
	if at := raisedAt(q); at >= 0 {
		s.Line, s.Raised = s.Line+q[:at], q[at:]
	} else {
		s.Line += q
	}
	if c.HasBass {
		s.Bass = "/" + naming.English.Name(c.Bass, naming.Signs)
	}
	return s
}

// raisedAt returns where the raised run of a quality `q` starts, or -1.
func raisedAt(q string) int {
	if strings.HasPrefix(q, "m(maj7") {
		return -1
	}
	at := -1
	for _, mark := range []string{"♭5", "♯5", "("} {
		if i := strings.Index(q, mark); i >= 0 && (at < 0 || i < at) {
			at = i
		}
	}
	return at
}
