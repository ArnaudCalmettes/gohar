// Package naming attaches names to the numbers the core computes.
//
// The mode nomenclature follows the system taught by Bernard Maury.
// Five systems, seven modes each, thirty five modes with distinct
// patterns. Two of them cannot be confused by ear, and none of them
// share a pattern, which is what lets a played set resolve to a single
// name.
package naming

import (
	"slices"
	"strconv"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A NaturalMode is one of the seven modes of the natural system. Every
// other mode in the catalogue is named after one of these, plus a list
// of altered degrees.
type NaturalMode uint8

const (
	Ionian NaturalMode = iota
	Dorian
	Phrygian
	Lydian
	Mixolydian
	Aeolian
	Locrian
)

// A Quality is the state of a degree relative to the major scale.
//
// This is the crux of the nomenclature and the place a naive reading
// goes wrong. A quality is absolute, not a movement: it says what the
// degree is, not how far it moved from the base mode.
//
// Phrygian natural 6 carries Natural on its sixth degree, even though
// the phrygian sixth is flat and the alteration raises it. Locrian
// double flat 7 carries DoubleFlat, giving nine semitones, and not
// two semitones below the locrian flat seventh, which would give
// eight and is a different scale. Four modes in the catalogue
// distinguish the two readings, and all four carry a double flat.
type Quality int8

const (
	DoubleFlat Quality = -2
	Flat       Quality = -1
	Natural    Quality = 0
	Sharp      Quality = 1
)

// An Alteration is a degree together with the quality it takes.
type Alteration struct {
	Degree  harmony.Degree
	Quality Quality
}

// A Mode is one entry of the catalogue.
//
// Pattern, DCN and DCA are all derived rather than stored. The pattern
// follows from the base mode and the alterations, the DCN is inherited
// from the base mode, and the DCA is the alteration list itself. What
// stays irreducibly hand written is the name, the chord symbol and the
// function, none of which follow from anything.
type Mode struct {
	System  harmony.System
	Degree  harmony.Degree
	Base    NaturalMode
	Altered []Alteration

	// Chord is the tetrad the mode voices, with its extensions.
	// Empty for the two modes that voice no tetrad.
	Chord string

	// Extensions are recorded separately when a mode has them without
	// a tetrad to hang them on.
	Extensions string

	Function harmony.Function
}

// A ModeKey designates a mode by its mother scale and degree, which is
// how a musician points at one without naming it.
type ModeKey struct {
	System harmony.System
	Degree harmony.Degree
}

// Key returns the designation of m.
func (m Mode) Key() ModeKey {
	return ModeKey{System: m.System, Degree: m.Degree}
}

// Pattern derives the scale pattern of m.
//
// Start from the offsets of the base natural mode, then for each
// altered degree replace the offset with the major scale offset for
// that degree adjusted by the quality. Degrees not named keep the base
// mode's offset.
func (m Mode) Pattern() harmony.ScalePattern {
	offsets := naturalModeOffsets(m.Base)
	for _, a := range m.Altered {
		if a.Degree < 1 || a.Degree > 7 {
			continue
		}
		offsets[a.Degree-1] = majorOffsets[a.Degree-1] + harmony.Semitones(a.Quality)
	}

	var p harmony.ScalePattern
	for _, n := range offsets {
		p |= 1 << (((n % 12) + 12) % 12)
	}
	return p
}

// NaturalDegrees returns the DCN of m, inherited from its base mode.
//
// Every mode carries the DCN of the natural mode it is named after,
// including the twenty eight altered ones. The rule holds even where
// an alteration cancels a neighbouring degree.
func (m Mode) NaturalDegrees() []Alteration {
	return slices.Clone(naturalModeDCN[m.Base])
}

// AlteredDegrees returns the DCA of m, which is its alteration list.
//
// The seven natural modes have none.
func (m Mode) AlteredDegrees() []Alteration {
	return slices.Clone(m.Altered)
}

// naturalModeDCN holds the characteristic natural degrees of the seven
// natural modes. Every altered mode inherits from this table through
// its base mode.
var naturalModeDCN = map[NaturalMode][]Alteration{
	Ionian:     {{4, Natural}, {7, Natural}},
	Dorian:     {{6, Natural}},
	Phrygian:   {{2, Flat}},
	Lydian:     {{4, Sharp}},
	Mixolydian: {{7, Flat}},
	Aeolian:    {{6, Flat}},
	Locrian:    {{5, Flat}},
}

// catalogue holds the thirty five modes.
//
// Unexported and never mutated after declaration. Build a lookup from
// it with NewModeCatalog rather than reaching in here.
var catalogue = []Mode{
	// Natural system. These carry no alteration, so their DCA is empty
	// and their pattern is the base mode itself.
	{
		System: harmony.NaturalMajor, Degree: 1, Base: Ionian,
		Chord: "XMaj7 (9, 11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.NaturalMajor, Degree: 2, Base: Dorian,
		Chord: "Xm7 (9, 11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.NaturalMajor, Degree: 3, Base: Phrygian,
		Chord: "Xm7 (b9, 11, b13)", Function: harmony.Tonic,
	},
	{
		System: harmony.NaturalMajor, Degree: 4, Base: Lydian,
		Chord: "XMaj7 (9, #11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.NaturalMajor, Degree: 5, Base: Mixolydian,
		Chord: "X7 (9, 11, 13)", Function: harmony.Dominant,
	},
	{
		System: harmony.NaturalMajor, Degree: 6, Base: Aeolian,
		Chord: "Xm7 (9, 11, b13)", Function: harmony.Tonic,
	},
	{
		System: harmony.NaturalMajor, Degree: 7, Base: Locrian,
		Chord: "Xm7b5 (b9, 11, b13)", Function: harmony.Subdominant,
	},

	// Melodic minor system.
	{
		System: harmony.MelodicMinor, Degree: 1, Base: Ionian,
		Altered: []Alteration{{3, Flat}},
		Chord:   "XmMaj7 (9, 11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.MelodicMinor, Degree: 2, Base: Phrygian,
		Altered: []Alteration{{6, Natural}},
		Chord:   "X7sus4 (b9, b10, 13)", Function: harmony.Subdominant,
	},
	{
		System: harmony.MelodicMinor, Degree: 3, Base: Lydian,
		Altered: []Alteration{{5, Sharp}},
		Chord:   "XMaj7 (#5, 9, #11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.MelodicMinor, Degree: 4, Base: Lydian,
		Altered: []Alteration{{7, Flat}},
		Chord:   "X7 (9, #11, 13)", Function: harmony.Dominant,
	},
	{
		System: harmony.MelodicMinor, Degree: 5, Base: Mixolydian,
		Altered: []Alteration{{6, Flat}},
		Chord:   "X7 (9, 11, b13)", Function: harmony.Tonic,
	},
	{
		System: harmony.MelodicMinor, Degree: 6, Base: Locrian,
		Altered: []Alteration{{2, Natural}},
		Chord:   "Xm7b5 (9, 11, b13)", Function: harmony.Subdominant,
	},
	{
		System: harmony.MelodicMinor, Degree: 7, Base: Locrian,
		Altered: []Alteration{{4, Flat}},
		Chord:   "X7Alt (b5, b9, b10, b13)", Function: harmony.Dominant,
	},

	// Harmonic minor system.
	{
		System: harmony.HarmonicMinor, Degree: 1, Base: Aeolian,
		Altered: []Alteration{{7, Natural}},
		Chord:   "XmMaj7 (9, 11, b13)", Function: harmony.Tonic | harmony.Dominant,
	},
	{
		System: harmony.HarmonicMinor, Degree: 2, Base: Locrian,
		Altered: []Alteration{{6, Natural}},
		Chord:   "Xm7b5 (b9, 11, 13)", Function: harmony.Subdominant,
	},
	{
		System: harmony.HarmonicMinor, Degree: 3, Base: Ionian,
		Altered: []Alteration{{5, Sharp}},
		Chord:   "XMaj7 (#5, 9, 11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.HarmonicMinor, Degree: 4, Base: Dorian,
		Altered: []Alteration{{4, Sharp}},
		Chord:   "Xm7 (9, #11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.HarmonicMinor, Degree: 5, Base: Phrygian,
		Altered: []Alteration{{3, Natural}},
		Chord:   "X7 (b9, 11, b13)", Function: harmony.Dominant,
	},
	{
		System: harmony.HarmonicMinor, Degree: 6, Base: Lydian,
		Altered: []Alteration{{2, Sharp}},
		Chord:   "XMaj7 (#9, #11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.HarmonicMinor, Degree: 7, Base: Locrian,
		Altered: []Alteration{{4, Flat}, {7, DoubleFlat}},
		Chord:   "Xdim7 (b9, b11, b13)", Function: harmony.Dominant,
	},

	// Harmonic major system.
	{
		System: harmony.HarmonicMajor, Degree: 1, Base: Ionian,
		Altered: []Alteration{{6, Flat}},
		Chord:   "XMaj7 (9, 11, b13)", Function: harmony.Tonic | harmony.Dominant,
	},
	{
		System: harmony.HarmonicMajor, Degree: 2, Base: Dorian,
		Altered: []Alteration{{5, Flat}},
		Chord:   "Xm7b5 (9, 11, 13)", Function: harmony.Subdominant,
	},
	{
		System: harmony.HarmonicMajor, Degree: 3, Base: Phrygian,
		Altered: []Alteration{{4, Flat}},
		Chord:   "X7 (b9, b10, b13)", Function: harmony.Dominant,
	},
	{
		System: harmony.HarmonicMajor, Degree: 4, Base: Lydian,
		Altered: []Alteration{{3, Flat}},
		Chord:   "XmMaj7 (9, #11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.HarmonicMajor, Degree: 5, Base: Mixolydian,
		Altered: []Alteration{{2, Flat}},
		Chord:   "X7 (b9, 11, 13)", Function: harmony.Dominant,
	},
	{
		System: harmony.HarmonicMajor, Degree: 6, Base: Lydian,
		Altered: []Alteration{{2, Sharp}, {5, Sharp}},
		Chord:   "XMaj7 (#5, #9, #11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.HarmonicMajor, Degree: 7, Base: Locrian,
		Altered: []Alteration{{7, DoubleFlat}},
		Chord:   "Xdim7 (b9, 11, b13)", Function: harmony.Dominant,
	},

	// Double harmonic major system. Every mode here carries two
	// alterations, the system sitting two alterations from the natural.
	{
		System: harmony.DoubleHarmonicMajor, Degree: 1, Base: Ionian,
		Altered: []Alteration{{2, Flat}, {6, Flat}},
		Chord:   "XMaj7 (b9, 11, b13)", Function: harmony.Tonic,
	},
	{
		System: harmony.DoubleHarmonicMajor, Degree: 2, Base: Lydian,
		Altered: []Alteration{{2, Sharp}, {6, Sharp}},
		Chord:   "XMaj7 (#9, #11, #13)", Function: harmony.Tonic,
	},
	{
		System: harmony.DoubleHarmonicMajor, Degree: 3, Base: Phrygian,
		Altered: []Alteration{{4, Flat}, {7, DoubleFlat}},
		Function: harmony.NoFunction,
	},
	{
		System: harmony.DoubleHarmonicMajor, Degree: 4, Base: Aeolian,
		Altered: []Alteration{{4, Sharp}, {7, Natural}},
		Chord:   "XmMaj7 (9, #11, b13)", Function: harmony.Tonic,
	},
	{
		System: harmony.DoubleHarmonicMajor, Degree: 5, Base: Mixolydian,
		Altered: []Alteration{{2, Flat}, {5, Flat}},
		Chord:   "X7 (b5, b9, 11, 13)", Function: harmony.Dominant,
	},
	{
		System: harmony.DoubleHarmonicMajor, Degree: 6, Base: Ionian,
		Altered: []Alteration{{2, Sharp}, {5, Sharp}},
		Chord:   "XMaj7 (#5, #9, 11, 13)", Function: harmony.Tonic,
	},
	{
		System: harmony.DoubleHarmonicMajor, Degree: 7, Base: Locrian,
		Altered:    []Alteration{{3, DoubleFlat}, {7, DoubleFlat}},
		Extensions: "b9, 11, b13",
		Function:   harmony.NoFunction,
	},
}

// Modes returns the catalogue.
//
// The slice is derived at call time rather than kept in a package
// variable: the catalogue is data, and a shared mutable slice would be
// the kind of global the architecture rules out.
func Modes() []Mode {
	return slices.Clone(catalogue)
}

// Lookup returns the mode at the given degree of the given system.
//
// This is the designation that needs no name at all, the one a musician
// uses when saying that F lydian is the fourth degree of C natural
// major.
func Lookup(s harmony.System, d harmony.Degree) (Mode, bool) {
	for _, m := range catalogue {
		if m.System == s && m.Degree == d {
			return m, true
		}
	}
	return Mode{}, false
}

// String returns the name of m: the base natural mode followed by its
// altered degrees, as in lydian #2 or locrian b4 bb7.
//
// The seven natural modes render as their bare name.
func (m Mode) String() string {
	name := plainModeNames[m.Base]
	for _, a := range m.Altered {
		name += " " + plainSigns[a.Quality+2] + strconv.Itoa(int(a.Degree))
	}
	return name
}

// majorOffsets is the reference every degree quality is measured
// against.
var majorOffsets = [7]harmony.Semitones{0, 2, 4, 5, 7, 9, 11}

// naturalModeOffsets returns the offsets of one of the seven natural
// modes, read from its own tonic.
//
// Computed rather than tabulated: a table would be a second place for
// the major scale to live, and the two would drift.
func naturalModeOffsets(m NaturalMode) [7]harmony.Semitones {
	var out [7]harmony.Semitones
	root := majorOffsets[int(m)%7]
	for i := range out {
		out[i] = ((majorOffsets[(int(m)+i)%7] - root) + 12) % 12
	}
	return out
}

// plainModeNames and plainSigns render a mode without a locale, for
// debugging and for test names. Anything shown to a player goes
// through [Locale.ModeName].
var plainModeNames = [7]string{
	"ionian", "dorian", "phrygian", "lydian", "mixolydian", "aeolian", "locrian",
}

var plainSigns = [5]string{"bb", "b", "nat", "#", "##"}
