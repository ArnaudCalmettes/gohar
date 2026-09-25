package naming

import (
	"strconv"
	"strings"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Locale holds the words one language uses for letters, accidentals
// and modes.
//
// A Locale is a value, passed to the constructor of a [Namer] and never
// reachable through a package variable. The previous iteration of this
// library kept a mutable CurrentLocale pointer, which made the choice
// of language a hidden global: two goroutines could not disagree about
// it, and no test could set it without disturbing another. Rendering
// and game logic will run on separate goroutines here, so that would
// have been a data race waiting for a busy frame.
type Locale struct {
	// Letters names the seven letters, indexed by [Letter].
	Letters [LetterCount]string

	// Accidentals names the five signs as they appear in a note name,
	// from double flat to double sharp. Whether they are words, ASCII
	// or Unicode symbols is the locale's business.
	//
	// The natural entry is empty here. A natural note is named by its
	// letter alone: one says F, not F natural.
	Accidentals [5]string

	// DegreeSigns names the same five signs as they appear in a mode
	// name, indexed the same way. Used by [SignDegrees].
	//
	// A separate table from Accidentals, because the natural entry is
	// not empty here. In a note name a natural is silent; in a mode
	// name it is the whole point, phrygian natural 6 being named for
	// the very degree that is not altered relative to the major scale.
	// Reusing Accidentals would render it as bare phrygian 6, which
	// names a different thing, or as plain phrygian, which names
	// another mode outright.
	DegreeSigns [5]string

	// Intervals names the seven interval sizes, indexed by degree
	// minus one. Used by [IntervalDegrees].
	Intervals [7]string

	// PerfectQualities names the qualities a perfect interval takes,
	// indexed like DegreeSigns. Degrees 1, 4 and 5.
	//
	// The outer two are one word rather than a doubling: French says
	// sous-diminuée and suraugmentée, not doublement diminuée.
	PerfectQualities [5]string

	// ImperfectQualities names the qualities the other intervals take.
	// Degrees 2, 3, 6 and 7.
	//
	// The two families differ only at the lower end. An imperfect
	// interval lowered once is minor and lowered twice is diminished,
	// so its diminished sits where the perfect family holds its
	// sous-diminuée, and the imperfect family reaches no further down.
	ImperfectQualities [5]string

	// Degree renders one altered degree of a mode name.
	//
	// A function rather than a table, because the two languages differ
	// in strategy and not merely in vocabulary. English jazz usage
	// writes a sign and a number, lydian sharp 2. French names the
	// interval outright, phrygien sixte majeure, which is a
	// computation over the degree and its quality rather than a
	// lookup.
	//
	// Build one with [SignDegrees] or [IntervalDegrees].
	Degree DegreeRenderer

	// NaturalModes names the seven modes of the natural system,
	// indexed by [NaturalMode]. Altered mode names are built from
	// these plus their alterations.
	NaturalModes [7]string

	// Functions names the three harmonic functions, indexed by
	// [harmony.Function]. Index zero is the absent one and should render as
	// the empty string.
	Functions [4]string

	// ModeQualities names what a single altered third or fifth makes of
	// a mode, in this order: major, minor, augmented. The phrygian
	// natural 3 is said phrygien majeur, the lydian sharp 5 lydien
	// augmenté. Empty in a locale that writes signs instead.
	ModeQualities [3]string

	// Aliases holds the other names a mode goes by, keyed by its mother
	// scale and degree.
	//
	// Written out one by one rather than derived. A rule that produced
	// them would be one somebody breaks by wanting a name it does not
	// allow, and they should be free to: a locale is theirs to extend.
	// Copying a Locale shares this map, so build a new one rather than
	// editing French or English in place.
	Aliases map[ModeKey][]string

	// Tetrachords names the tetrachords practice gives a name to, in
	// the order of [namedTetrachords]. Only those five: any other
	// tetrachord is designated by its steps in every language.
	Tetrachords [5]string

	// Spell renders a note from its parts. A locale that writes the
	// sign after the letter and one that writes it before need
	// different code, not a different table.
	Spell func(letter, accidental string) string
}

// A DegreeRenderer turns one altered degree into the words a mode name
// uses for it.
type DegreeRenderer func(harmony.Degree, Quality) string

// SignDegrees renders a degree as a sign followed by its number, as in
// sharp 2. This is the English jazz register.
func SignDegrees(l Locale) DegreeRenderer {
	return func(d harmony.Degree, q Quality) string {
		return l.DegreeSigns[q+2] + strconv.Itoa(int(d))
	}
}

// IntervalDegrees renders a degree as the interval it names, as in
// sixte majeure. This is the French register.
//
// The quality depends on whether the interval is perfect. Degrees 1, 4
// and 5 are perfect and run from doubly diminished through perfect to
// augmented; the others run from diminished through minor and major to
// augmented. Naming a natural fourth major, or a natural sixth perfect,
// is the mistake this split exists to prevent.
//
// This is the interval quality naming that the core deliberately does
// not carry: an augmented fourth and a diminished fifth are the same
// distance and two names, and choosing between them needs the degree,
// which only a tonal context supplies.
func IntervalDegrees(l Locale) DegreeRenderer {
	return func(d harmony.Degree, q Quality) string {
		step := (int(d) - 1) % 7
		if step < 0 {
			step += 7
		}
		qualities := l.ImperfectQualities
		if step == 0 || step == 3 || step == 4 {
			qualities = l.PerfectQualities
		}
		return l.Intervals[step] + " " + qualities[q+2]
	}
}

// French names notes with the solfège syllables, and modes the way
// they are said aloud: see [Locale.ModeName].
var French = Locale{
	Letters:      [LetterCount]string{"do", "ré", "mi", "fa", "sol", "la", "si"},
	Accidentals:  [5]string{"double bémol", "bémol", "", "dièse", "double dièse"},
	NaturalModes: [7]string{"ionien", "dorien", "phrygien", "lydien", "mixolydien", "éolien", "locrien"},
	Spell:        spellApart,
	Functions:    [4]string{"", "tonique", "sous-dominante", "dominante"},
	Tetrachords:  [5]string{"majeur", "mineur", "phrygien", "lydien", "harmonique"},

	DegreeSigns:   [5]string{"♭♭", "♭", "♮", "♯", "♯♯"},
	ModeQualities: [3]string{"majeur", "mineur", "augmenté"},
	Aliases: map[ModeKey][]string{
		{harmony.MelodicMinor, 7}:  {"altéré"},
		{harmony.MelodicMinor, 4}:  {"lydien dominante"},
		{harmony.HarmonicMinor, 5}: {"phrygien dominante"},
	},

	Intervals: [7]string{
		"unisson", "seconde", "tierce", "quarte",
		"quinte", "sixte", "septième",
	},
	PerfectQualities: [5]string{
		"sous-diminuée", "diminuée", "juste", "augmentée", "suraugmentée",
	},
	ImperfectQualities: [5]string{
		"diminuée", "mineure", "majeure", "augmentée", "suraugmentée",
	},
}

// English names notes with letters and Unicode signs.
var English = Locale{
	Letters:      [LetterCount]string{"C", "D", "E", "F", "G", "A", "B"},
	Accidentals:  [5]string{"\u266d\u266d", "\u266d", "", "\u266f", "\u266f\u266f"},
	DegreeSigns:  [5]string{"\u266d\u266d", "\u266d", "\u266e", "\u266f", "\u266f\u266f"},
	NaturalModes: [7]string{"ionian", "dorian", "phrygian", "lydian", "mixolydian", "aeolian", "locrian"},
	Spell:        spellTogether,

	Aliases: map[ModeKey][]string{
		{harmony.MelodicMinor, 7}:  {"altered"},
		{harmony.MelodicMinor, 4}:  {"lydian dominant"},
		{harmony.HarmonicMinor, 5}: {"phrygian dominant"},
	},
	Functions:    [4]string{"", "tonic", "subdominant", "dominant"},
	Tetrachords:  [5]string{"major", "minor", "phrygian", "lydian", "harmonic"},
}

// Name returns the words for a spelled note in this locale.
func (l Locale) Name(n SpelledNote) string {
	letter := l.Letters[int(n.Letter)%LetterCount]
	accidental := l.Accidentals[n.Accidental+2]
	if l.Spell != nil {
		return l.Spell(letter, accidental)
	}
	return letter + accidental
}

// ModeName returns the words for a mode in this locale: the base
// natural mode followed by its altered degrees, each written with its
// sign from [Locale.DegreeSigns] and its degree number.
//
// The seven natural modes render as their bare name, having no altered
// degree to append.
//
// The altered degrees come out in the order the catalogue holds them,
// which is the order they are named in. That order is data, not a
// sort: lydian sharp 2 sharp 5 is named that way and would read wrong
// the other way round.
func (l Locale) ModeName(m Mode) string {
	name := l.NaturalModes[int(m.Base)%7]
	if len(m.Altered) == 0 {
		return name
	}
	if l.Degree == nil && l.Intervals[0] != "" {
		return name + " " + l.spokenAlterations(m.Altered)
	}
	render := l.degreeRenderer()
	for _, a := range m.Altered {
		name += " " + render(a.Degree, a.Quality)
	}
	return name
}

// ModeAliases returns the other names this locale gives a mode, or nil
// when it has none.
func (l Locale) ModeAliases(m Mode) []string {
	return l.Aliases[m.Key()]
}

// spokenAlterations renders the alterations of a mode the way they are
// said aloud.
//
// The procedure follows how a musician names a mode, and it looks at
// the whole mode rather than one degree at a time, because the first
// rule depends on how many degrees are altered:
//
//  1. Two alterations or more: each is spelled with its sign, ionien
//     ♯2 ♯5.
//  2. A single altered third: the quality it gives the mode, majeur or
//     mineur. A single raised fifth: augmenté.
//  3. A single alteration whose interval would be diminished: spelled
//     with its sign, ♭5, ♭4, ♭♭7.
//  4. Any other single alteration: its interval name, lydien seconde
//     augmentée, mixolydien sixte mineure.
//
// Rules 1 and 3 exist for the same reason. The word diminuée stands for
// a different sign depending on the degree, a double flat on a third or
// a seventh and a single flat on a fourth or a fifth, so saying the
// sign removes an ambiguity the word creates.
func (l Locale) spokenAlterations(altered []Alteration) string {
	if len(altered) >= 2 {
		parts := make([]string, len(altered))
		for i, a := range altered {
			parts[i] = l.spelledDegree(a)
		}
		return strings.Join(parts, " ")
	}

	a := altered[0]
	switch {
	case a.Degree == 3 && a.Quality == Natural && l.ModeQualities[0] != "":
		return l.ModeQualities[0]
	case a.Degree == 3 && a.Quality == Flat && l.ModeQualities[1] != "":
		return l.ModeQualities[1]
	case a.Degree == 5 && a.Quality == Sharp && l.ModeQualities[2] != "":
		return l.ModeQualities[2]
	case isDiminished(a):
		return l.spelledDegree(a)
	}
	return IntervalDegrees(l)(a.Degree, a.Quality)
}

// spelledDegree writes an alteration as its sign followed by its degree,
// as in ♯2.
func (l Locale) spelledDegree(a Alteration) string {
	return l.DegreeSigns[a.Quality+2] + strconv.Itoa(int(a.Degree))
}

// isDiminished reports whether an alteration makes its interval
// diminished: one flat on a perfect degree, two on the others.
func isDiminished(a Alteration) bool {
	step := (int(a.Degree) - 1) % 7
	if step == 0 || step == 3 || step == 4 {
		return a.Quality <= Flat
	}
	return a.Quality <= DoubleFlat
}

// FunctionName returns the word for a harmonic function in this
// locale, and the empty string when there is none.
func (l Locale) FunctionName(f harmony.Function) string {
	var parts []string
	for i, role := range []harmony.Function{harmony.Tonic, harmony.Subdominant, harmony.Dominant} {
		if f.Has(role) {
			parts = append(parts, l.Functions[i+1])
		}
	}
	if len(parts) == 0 {
		return l.Functions[0]
	}
	return strings.Join(parts, ", ")
}

// degreeRenderer resolves which register this locale writes mode names
// in.
//
// Resolved on demand rather than assigned in an init function, so that
// the declared locales stay immutable values. A locale that fills its
// interval tables gets the interval register; one that does not falls
// back to signs.
func (l Locale) degreeRenderer() DegreeRenderer {
	switch {
	case l.Degree != nil:
		return l.Degree
	case l.Intervals[0] != "":
		return IntervalDegrees(l)
	default:
		return SignDegrees(l)
	}
}

// spellApart writes the accidental as a separate word, which is what
// French does: fa dièse, not fa#.
func spellApart(letter, accidental string) string {
	if accidental == "" {
		return letter
	}
	return letter + " " + accidental
}

// spellTogether writes the accidental against the letter, which is what
// a symbol wants.
func spellTogether(letter, accidental string) string {
	return letter + accidental
}

// namedTetrachords lists the tetrachords that practice gives a name to,
// in the order a Locale's Tetrachords table follows.
//
// Which tetrachords are named is a fact about the practice, not about a
// language, so it is decided here once. A locale supplies words for
// these five and cannot name a sixth: a shape that has no name in one
// language has none in any.
var namedTetrachords = [5]harmony.Tetrachord{
	harmony.TetrachordMajor,
	harmony.TetrachordMinor,
	harmony.TetrachordPhrygian,
	harmony.TetrachordLydian,
	harmony.TetrachordHarmonic,
}

// TetrachordName returns the word this locale uses for a tetrachord, or
// its steps when it has none.
//
// The result is a qualifier, meant to follow the word for tetrachord in
// the sentence: « le tétracorde harmonique », « le tétracorde 1 1 3 ».
// Both read naturally, which is the point of designating the unnamed
// ones by their steps rather than leaving a gap.
func (l Locale) TetrachordName(t harmony.Tetrachord) string {
	for i, named := range namedTetrachords {
		if named == t && l.Tetrachords[i] != "" {
			return l.Tetrachords[i]
		}
	}
	return t.String()
}
