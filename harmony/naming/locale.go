package naming

import (
	"strconv"
	"strings"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Notation says how accidentals are written: as signs or as words.
//
// A separate choice from the language. Both languages can write si♭ or
// si bémol, B♭ or B flat, and the choice is one of presentation:
// signs are what a musician reads, words are what a speech synthesiser
// can say. A game may well hold both, one for the screen and one for a
// screen reader, which is why it lives on the [Namer] and not in the
// [Locale].
type Notation uint8

const (
	// Signs writes ♭ ♮ ♯ against the letter or the degree. The default.
	Signs Notation = iota

	// Words writes them out in the locale's language, apart from the
	// letter: si bémol, phrygien bécarre 6.
	Words
)

// accidentalSigns and degreeSigns are the same five signs in every
// language, from double flat to double sharp.
//
// Two tables because the natural is silent in a note name and spoken in
// a mode name: phrygian natural 6 is named for the very degree that is
// not altered relative to the major scale. Dropping it would render
// bare phrygian 6, which names a different thing, or plain phrygian,
// which names another mode outright.
//
// The doubles are signs of their own, 𝄫 and 𝄪, and not two simple signs
// side by side: that is how they are engraved and read.
var (
	accidentalSigns = [5]string{"\U0001D12B", "\u266d", "", "\u266f", "\U0001D12A"}
	degreeSigns     = [5]string{"\U0001D12B", "\u266d", "\u266e", "\u266f", "\U0001D12A"}
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

	// AccidentalWords names the five accidentals as words, from double
	// flat to double sharp, for the [Words] notation. The signs are the
	// same in every language and are not the locale's business.
	//
	// The natural entry is empty: a natural note is named by its letter
	// alone, one says F, not F natural.
	AccidentalWords [5]string

	// DegreeWords names the same five in a mode name, where the natural
	// is not empty: phrygien bécarre 6.
	DegreeWords [5]string

	// Intervals names the seven interval sizes, indexed by degree
	// minus one. Used by [IntervalDegrees] and the spoken register of
	// [Locale.SpokenModeName]; empty in a locale that has none.
	Intervals [7]string

	// PerfectQualities names the qualities a perfect interval takes,
	// indexed like DegreeWords. Degrees 1, 4 and 5.
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

	// NaturalModes names the seven modes of the natural system,
	// indexed by [NaturalMode]. Altered mode names are built from
	// these plus their alterations.
	NaturalModes [7]string

	// Functions names the three harmonic functions, indexed by
	// [harmony.Function]. Index zero is the absent one and should render as
	// the empty string.
	Functions [4]string

	// ModeQualities names what a single altered third or fifth makes of
	// a mode in the spoken register, in this order: major, minor,
	// augmented. The phrygian natural 3 is said phrygien majeur, the
	// lydian sharp 5 lydien augmenté. Empty in a locale that has no
	// spoken register.
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
}

// A DegreeRenderer turns one altered degree into the words a mode name
// uses for it.
type DegreeRenderer func(harmony.Degree, Quality) string

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

// French names notes with the solfège syllables. It also has a spoken
// register for modes, the way they are said aloud: see
// [Locale.SpokenModeName].
var French = Locale{
	Letters:         [LetterCount]string{"do", "ré", "mi", "fa", "sol", "la", "si"},
	AccidentalWords: [5]string{"double bémol", "bémol", "", "dièse", "double dièse"},
	DegreeWords:     [5]string{"double bémol", "bémol", "bécarre", "dièse", "double dièse"},
	NaturalModes:    [7]string{"ionien", "dorien", "phrygien", "lydien", "mixolydien", "éolien", "locrien"},
	Functions:       [4]string{"", "tonique", "sous-dominante", "dominante"},
	Tetrachords:     [5]string{"majeur", "mineur", "phrygien", "lydien", "harmonique"},

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

// English names notes with letters. It has no spoken register: its
// usage says the sign and the number, which is the systematic name.
var English = Locale{
	Letters:         [LetterCount]string{"C", "D", "E", "F", "G", "A", "B"},
	AccidentalWords: [5]string{"double flat", "flat", "", "sharp", "double sharp"},
	DegreeWords:     [5]string{"double flat", "flat", "natural", "sharp", "double sharp"},
	NaturalModes:    [7]string{"ionian", "dorian", "phrygian", "lydian", "mixolydian", "aeolian", "locrian"},

	Aliases: map[ModeKey][]string{
		{harmony.MelodicMinor, 7}:  {"altered"},
		{harmony.MelodicMinor, 4}:  {"lydian dominant"},
		{harmony.HarmonicMinor, 5}: {"phrygian dominant"},
	},
	Functions:   [4]string{"", "tonic", "subdominant", "dominant"},
	Tetrachords: [5]string{"major", "minor", "phrygian", "lydian", "harmonic"},
}

// Name returns a spelled note in this locale and notation: fa♯ or fa
// dièse, F♯ or F sharp.
//
// Signs go against the letter and words apart from it, whatever the
// language: that follows from the notation, not from the locale.
func (l Locale) Name(n SpelledNote, notation Notation) string {
	letter := l.Letters[int(n.Letter)%LetterCount]
	i := int(n.Accidental) + 2
	if n.Accidental == NaturalSign {
		return letter
	}
	if notation == Words {
		return letter + " " + l.AccidentalWords[i]
	}
	return letter + accidentalSigns[i]
}

// ModeName returns the systematic name of a mode: the base natural
// mode followed by each altered degree, its sign or word then its
// number. Phrygien ♮6, lydian ♯2 ♯5; in words, phrygien bécarre 6.
//
// The name every mode has in every locale, and the default. Nothing is
// special cased, which is what makes it predictable to read and to
// say. The refined names, spoken register and aliases alike, are
// alternatives: see [Locale.ModeAlternatives].
//
// The altered degrees come out in the order the catalogue holds them,
// which is the order they are named in. That order is data, not a
// sort: lydian sharp 2 sharp 5 is named that way and would read wrong
// the other way round.
func (l Locale) ModeName(m Mode, notation Notation) string {
	name := l.NaturalModes[int(m.Base)%7]
	for _, a := range m.Altered {
		name += " " + l.degree(a, notation)
	}
	return name
}

// degree writes one alteration with its sign or word, then its number.
func (l Locale) degree(a Alteration, notation Notation) string {
	i := int(a.Quality) + 2
	if notation == Words {
		return l.DegreeWords[i] + " " + strconv.Itoa(int(a.Degree))
	}
	return degreeSigns[i] + strconv.Itoa(int(a.Degree))
}

// ModeAliases returns the other names this locale gives a mode, or nil
// when it has none.
func (l Locale) ModeAliases(m Mode) []string {
	return l.Aliases[m.Key()]
}

// ModeAlternatives returns every other name of a mode: its spoken name
// when it differs from the systematic one, then its aliases. Nil when
// the systematic name is the only one.
//
// What they are for is rigour and refinement, and recognition: a player
// who knows the phrygian dominant under that name should find it.
func (l Locale) ModeAlternatives(m Mode, notation Notation) []string {
	var out []string
	if spoken, ok := l.SpokenModeName(m, notation); ok && spoken != l.ModeName(m, notation) {
		out = append(out, spoken)
	}
	return append(out, l.ModeAliases(m)...)
}

// SpokenModeName returns the name of a mode as it is said aloud in
// this locale's refined register, and false when the locale has none
// or the mode has no alteration to say.
//
// French has one: phrygien sixte majeure, lydien augmenté. It is an
// alternative, never the default name.
func (l Locale) SpokenModeName(m Mode, notation Notation) (string, bool) {
	if len(m.Altered) == 0 || l.Intervals[0] == "" {
		return "", false
	}
	return l.NaturalModes[int(m.Base)%7] + " " + l.spokenAlterations(m.Altered, notation), true
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
//     with its sign, ♭5, ♭4, 𝄫7.
//  4. Any other single alteration: its interval name, lydien seconde
//     augmentée, mixolydien sixte mineure.
//
// Signs or words in rules 1 and 3 follow the notation.
//
// Rules 1 and 3 exist for the same reason. The word diminuée stands for
// a different sign depending on the degree, a double flat on a third or
// a seventh and a single flat on a fourth or a fifth, so saying the
// sign removes an ambiguity the word creates.
func (l Locale) spokenAlterations(altered []Alteration, notation Notation) string {
	if len(altered) >= 2 {
		parts := make([]string, len(altered))
		for i, a := range altered {
			parts[i] = l.degree(a, notation)
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
		return l.degree(a, notation)
	}
	return IntervalDegrees(l)(a.Degree, a.Quality)
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
