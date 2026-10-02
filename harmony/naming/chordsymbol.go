package naming

import (
	"strings"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A ChordStyle chooses among the usual ways of writing the quality of a
// chord symbol. Its zero value is the default, as En Harmonie writes
// them: Cmaj7, Cm7, Cm7♭5, Cdim7.
//
// Lead sheets disagree on these signs, and a reader is used to one set
// or another: a choice of presentation, like [Notation], and not a
// matter of language, since English and French charts write the same
// symbols.
type ChordStyle struct {
	// MajorSeventh writes the major seventh: maj7 by default, ♮7 or Δ7.
	MajorSeventh MajorSeventhSign

	// Minus writes a minor chord C-7 rather than Cm7.
	Minus bool

	// HalfDiminishedSign writes the half-diminished seventh Cø rather
	// than Cm7♭5.
	HalfDiminishedSign bool

	// DiminishedSign writes the diminished chords C° and C°7 rather than
	// Cdim and Cdim7.
	DiminishedSign bool

	// Slash writes the sixth chord with its ninth C6/9, with the slash
	// of the keyboard, rather than C6⁄9, with the fraction slash.
	Slash bool
}

// fractionSlash separates the 6 and the 9 of a sixth chord with its
// ninth, as the fraction ⁶⁄₉ of an engraved chart does, on one line.
// The slash of the keyboard reads as well, a 9 being no bass note:
// the option [ChordStyle.Slash].
const fractionSlash = "\u2044"

// A MajorSeventhSign is a way of writing the major seventh in a chord
// symbol.
type MajorSeventhSign uint8

// The natural sign reads well only raised, as on an engraved chart: on
// one line of text, C♮9 could be read as a C with a natural ninth. The
// terminal keeps maj by default.
const (
	MajSeventh     MajorSeventhSign = iota // Cmaj7, Cmaj9, Cm(maj7): the default
	NaturalSeventh                         // C♮7, C♮9, Cm(♮7), for a raised sign
	DeltaSeventh                           // CΔ7, CΔ9, C-(Δ7)
)

// The extensions, as offsets in the second octave of a normalised
// pattern (see [harmony.ChordPattern.Normalize]).
const (
	flatNinth        = 13
	ninth            = 14
	sharpNinth       = 15
	eleventh         = 17
	sharpEleventh    = 18
	flatThirteenth   = 20
	thirteenth       = 21
	majorFourteenth  = 23
	extensionsFrom   = 12
	extensionsBefore = 24
)

// extensionNames names each extension as it is written in parentheses.
// The major fourteenth is the major seventh over a diminished chord,
// raised an octave by Normalize: Cdim7(♮14), from the scale 2-1.
var extensionNames = map[harmony.Semitones]string{
	flatNinth:       "♭9",
	ninth:           "9",
	sharpNinth:      "♯9",
	eleventh:        "11",
	sharpEleventh:   "♯11",
	flatThirteenth:  "♭13",
	thirteenth:      "13",
	majorFourteenth: "♮14",
}

// altered is the 7alt: the chord of the altered mode, locrian ♭4,
// C7(♭5, ♭9, ♭10, ♭13). Its ♭10 sounds where a ♯9 would.
var altered = mustPattern(0, 4, 6, 10, 13, 15, 20)

// Symbol writes the quality of a chord symbol for a pattern, the root
// and the bass left to the caller: "m7", "maj7(♯11)", "13", "7alt". It
// reports false for a pattern whose first octave is no chord it names.
//
// The pattern is read as harmony normalises it (see
// [harmony.ChordPattern.Normalize]): a ♭5 is a fifth and a
// ♯11 an extension, a ♯5 a fifth and a ♭13 an extension, so the two are
// never confused.
//
// The fifth, when altered, stays against the seventh, as charts write
// it: C7♭5, Cm7♭5. The natural extensions climb on the seventh while
// they are stacked, as charts write them too: C9, C13 (the 9 and the 13,
// the 11 left out, as on a major third it is avoided), Cm11 (the 9 and
// the 11), Cm13 (the 9, the 11 and the 13). The 11 only climbs on a
// minor third: on a major one, C9(11) or C13(11), it is the note a
// chart writes apart (see [climbing]). What is left goes in parentheses,
// from the lowest: C7(♭9,♯11). The minor-major seventh opens its own
// parentheses, which its extensions join: Cm(maj7,9).
//
// A chord without a seventh adds its ninth (Cadd9, Cmadd9), a sixth
// chord stacks it (C6/9). A note that a voicing leaves out is not shown:
// a chord without its fifth is written as the chord.
func (s ChordStyle) Symbol(p harmony.ChordPattern) (string, bool) {
	p = p.Normalize()
	if p == altered {
		return "7alt", true
	}
	m := "m"
	if s.Minus {
		m = "-"
	}
	ext := extensionsOf(p)
	climb := func(minor bool) string { return climbing(ext, minor) }

	var head string
	switch p.Tetrad() {
	case harmony.ChordMajorTriad:
		head = added(ext, "")
	case harmony.ChordMinorTriad:
		head = added(ext, m)
	case harmony.ChordDiminishedTriad:
		head = s.diminished("")
	case harmony.ChordAugmentedTriad:
		head = "+"
	case harmony.ChordSus2:
		head = "sus2"
	case harmony.ChordSus4:
		head = "sus4"
	case fifthOnly:
		head = "5"
	case harmony.ChordMajorSixth:
		head = s.sixNine(ext, "6")
	case harmony.ChordMinorSixth:
		head = s.sixNine(ext, m+"6")
	case minorFlatSixth:
		head = m + "♭6"
	case minorSharpFifth:
		head = m + "♯5"
	case harmony.ChordMajorSeventh, harmony.ChordMajorSeventhNo5:
		head = s.major(climb(false))
	case harmony.ChordMajorSeventhSharp5:
		head = s.major(climb(false)) + "♯5"
	case harmony.ChordDominantSeventh, harmony.ChordDominantSeventhNo5:
		head = climb(false)
	case harmony.ChordDominantSeventhFlat5:
		head = climb(false) + "♭5"
	case harmony.ChordDominantSeventhSharp5:
		head = climb(false) + "♯5"
	case harmony.ChordMinorSeventh, harmony.ChordMinorSeventhNo5:
		head = m + climb(true)
	case harmony.ChordHalfDiminished:
		if n := climb(true); s.HalfDiminishedSign {
			head = "ø" + strings.TrimPrefix(n, "7")
		} else {
			head = m + n + "♭5"
		}
	case harmony.ChordMinorMajorSeventh, harmony.ChordMinorMajorSeventhNo5:
		// Its seventh opens the parentheses, and every extension goes in
		// them after it, none climbing: Cm(maj7,9,♯11).
		return m + "(" + strings.Join(append([]string{s.major("7")}, ext.names()...), ",") + ")", true
	case harmony.ChordDiminishedSeventh:
		head = s.diminished("7")
	case harmony.ChordDominantSeventhSus4:
		head = climb(false) + "sus4"
	case harmony.ChordDominantSeventhSus2:
		head = climb(false) + "sus2"
	default:
		return "", false
	}
	return head + ext.rest(), true
}

// The first octaves harmony has no constant for, that the app writes:
// C5, Cm♭6, Cm♯5.
var (
	fifthOnly       = mustPattern(0, 7)
	minorFlatSixth  = mustPattern(0, 3, 7, 8)
	minorSharpFifth = mustPattern(0, 3, 8)
)

// major writes a major seventh, or the extension that climbed on it, in
// the style's sign: maj7, ♮9, Δ13.
func (s ChordStyle) major(n string) string {
	switch s.MajorSeventh {
	case DeltaSeventh:
		return "Δ" + n
	case NaturalSeventh:
		return "♮" + n
	}
	return "maj" + n
}

// diminished writes a diminished chord in the style's sign: dim, °7.
func (s ChordStyle) diminished(seventh string) string {
	if s.DiminishedSign {
		return "°" + seventh
	}
	return "dim" + seventh
}

// extensions are the extensions of a pattern not yet written.
type extensions map[harmony.Semitones]bool

func extensionsOf(p harmony.ChordPattern) extensions {
	e := extensions{}
	for n := range p.Offsets() {
		if n >= extensionsFrom && n < extensionsBefore {
			e[n] = true
		}
	}
	return e
}

// take reports whether the extension `n` is there, and marks it written.
func (e extensions) take(n harmony.Semitones) bool {
	had := e[n]
	delete(e, n)
	return had
}

// names names the extensions left, from the lowest.
func (e extensions) names() []string {
	var out []string
	for n := harmony.Semitones(extensionsFrom); n < extensionsBefore; n++ {
		if e[n] {
			out = append(out, extensionNames[n])
		}
	}
	return out
}

// rest writes the extensions left in parentheses, from the lowest:
// "(♭9,♯11)", or nothing. No space after the comma: on a chart, a
// space separates two chords.
func (e extensions) rest() string {
	names := e.names()
	if len(names) == 0 {
		return ""
	}
	return "(" + strings.Join(names, ",") + ")"
}

// climbing returns the number a seventh is written with, the natural
// extensions stacked on it taken: 7, 9, 11, 13. On a minor third, the
// stack has no gap: 11 holds the 9, 13 the 9 and the 11, and a 13
// without its 11 stays apart, Cm9(13). On a major third, the 11 never
// climbs: 13 holds the 9 alone, and the 11 stays apart, C13(11).
func climbing(e extensions, minor bool) string {
	if !e.take(ninth) {
		return "7"
	}
	if minor {
		if !e.take(eleventh) {
			return "9"
		}
		if e.take(thirteenth) {
			return "13"
		}
		return "11"
	}
	if e.take(thirteenth) {
		return "13"
	}
	return "9"
}

// added writes a triad with its added ninth: add9, madd9.
func added(e extensions, triad string) string {
	if e.take(ninth) {
		return triad + "add9"
	}
	return triad
}

// sixNine writes a sixth chord with its ninth stacked: 6⁄9, m6⁄9.
func (s ChordStyle) sixNine(e extensions, six string) string {
	if !e.take(ninth) {
		return six
	}
	if s.Slash {
		return six + "/9"
	}
	return six + fractionSlash + "9"
}

func mustPattern(offsets ...harmony.Semitones) harmony.ChordPattern {
	p, err := harmony.NewChordPattern(offsets...)
	if err != nil {
		panic(err)
	}
	return p
}
