package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Degree is a chord read in a tonality, as En Harmonie writes it under a
// chart: II, V7 of a borrowed quality, ♭VII7, I/3.
//
// # Diatonic or borrowed
//
// A chord whose notes all belong to the tonality is written by its
// degree alone: II, not IIm7, the quality being the tonality's. A chord that
// borrows is written with its quality, which is what tells the reader
// it came from elsewhere: IVm7, ♭II7, ♯IVdim7. Quality is then the
// chord's tetrad, and zero for a diatonic chord.
//
// The degree counts in the tonality's own scale: in F minor, D♭maj7 is VI,
// not ♭VI. A root outside every scale of the tonality is a degree raised or
// lowered, lowered by preference (♭II, ♭III, ♭VI, ♭VII) but raised
// below the fourth and the fifth (♯IV, ♯III in minor), as En Harmonie
// writes them. A passing chord follows its bass instead: raised going
// up (♯Idim7, ♯Vdim7), lowered going down (♭IIIdim7).
type Degree struct {
	Number     harmony.Degree // 1 to 7, 0 for a silence or no tonality
	Accidental int            // -1 lowered, +1 raised
	Quality    harmony.ChordPattern

	// Inversion is the chord degree in the bass, 3, 5 or 7, and 0 in
	// root position: I/3.
	Inversion int
}

// DegreeOf reads a change in tonalities that share a tonic.
func DegreeOf(ch Change, ts []harmony.Tonality) Degree {
	return degreeOf(ch, ts, 0)
}

// degreeOf reads a change in tonalities, naming an altered root by the way
// a passing bass walks, when it does.
func degreeOf(ch Change, ts []harmony.Tonality, walk int) Degree {
	if ch.Silent || len(ts) == 0 {
		return Degree{}
	}
	var d Degree
	s := harmony.Semitones((int(ch.Chord.Root) - int(ts[0].Tonic()) + 12) % 12)
	found := false
	for _, t := range ts {
		for n, off := range t.Pattern().Offsets() {
			if off == s {
				d.Number, found = n, true
			}
		}
		if found {
			break
		}
	}
	// A passing chord is named by its walking bass, even when its root
	// is a degree under another name: D♯dim7 between Dm7 and C7/E in C
	// minor is ♯IIdim7, not IIIdim7.
	if !found || walk != 0 {
		d.Number, d.Accidental = altered(ts[0].Pattern(), s, walk)
	}
	if !holds(ts, ch.Chord.Set()) {
		d.Quality = ch.Chord.Pattern.Tetrad()
	}
	switch (int(ch.Bass) - int(ch.Chord.Root) + 12) % 12 {
	case 3, 4:
		d.Inversion = 3
	case 6, 7, 8:
		d.Inversion = 5
	case 9, 10, 11:
		d.Inversion = 7
	}
	if !ch.Inverted() {
		d.Inversion = 0
	}
	return d
}

// altered names a root outside a scale by the degree it alters.
func altered(p harmony.ScalePattern, s harmony.Semitones, walk int) (harmony.Degree, int) {
	above, below := harmony.Degree(0), harmony.Degree(0)
	for n, off := range p.Offsets() {
		switch off {
		case (s + 1) % 12:
			above = n
		case (s + 11) % 12:
			below = n
		}
	}
	switch {
	case walk > 0 && below > 0:
		return below, 1
	case walk < 0 && above > 0:
		return above, -1
	case below > 0 && (above == 4 || above == 5):
		return below, 1
	case above > 0:
		return above, -1
	}
	return below, 1
}

// String writes a degree as En Harmonie does: the numeral, then the
// quality of a borrowed chord, then the inversion. For tests and tools,
// until naming renders degrees in every language.
func (d Degree) String() string {
	if d.Number == 0 {
		return ""
	}
	s := [...]string{"♭", "", "♯"}[d.Accidental+1] +
		[...]string{"", "I", "II", "III", "IV", "V", "VI", "VII"}[d.Number]
	if d.Quality != 0 {
		q, ok := qualities[d.Quality]
		if !ok {
			q = "?"
		}
		s += q
	}
	if d.Inversion > 0 {
		s += "/" + string(rune('0'+d.Inversion))
	}
	return s
}

var qualities = map[harmony.ChordPattern]string{
	harmony.ChordMajorTriad:            "",
	harmony.ChordMinorTriad:            "m",
	harmony.ChordDiminishedTriad:       "dim",
	harmony.ChordAugmentedTriad:        "+",
	harmony.ChordSus2:                  "sus2",
	harmony.ChordSus4:                  "sus4",
	harmony.ChordMajorSixth:            "6",
	harmony.ChordMinorSixth:            "m6",
	harmony.ChordMajorSeventh:          "maj7",
	harmony.ChordMajorSeventhNo5:       "maj7",
	harmony.ChordMajorSeventhSharp5:    "maj7♯5",
	harmony.ChordDominantSeventh:       "7",
	harmony.ChordDominantSeventhNo5:    "7",
	harmony.ChordDominantSeventhFlat5:  "7♭5",
	harmony.ChordDominantSeventhSharp5: "7♯5",
	harmony.ChordMinorSeventh:          "m7",
	harmony.ChordMinorSeventhNo5:       "m7",
	harmony.ChordHalfDiminished:        "m7♭5",
	harmony.ChordMinorMajorSeventh:     "m(maj7)",
	harmony.ChordMinorMajorSeventhNo5:  "m(maj7)",
	harmony.ChordDiminishedSeventh:     "dim7",
	harmony.ChordDominantSeventhSus2:   "7sus2",
	harmony.ChordDominantSeventhSus4:   "7sus4",
}

// Degrees reads every change of a sequence in its ground, heard
// afterwards (see [Grounds]), with its quality when it borrows: what
// each chord is in the tonality installed there. In E flat, Gm7♭5
// C7♭9 before Fm7♭5 is IIIm7♭5 VI7, the three six of a three six two
// five one; A7 before Dm7 in C is VI7, a secondary dominant; a passing
// chord is ♯Idim7 or ♭IIIdim7, even when it is also the dominant
// without root of the next chord. Where the tune modulates, the
// degrees count from the new tonic from the cadence that led there on:
// in Tune Up, Dm7 G7 Cmaj7 is II V I.
//
// [Bracketed] gives the other reading, the one En Harmonie prints.
func Degrees(c Changes, passing []int, sensed []Sensed) []Degree {
	grounds := Grounds(c, sensed)
	out := make([]Degree, len(c.Chords))
	for i, ch := range c.Chords {
		out[i] = degreeOf(ch, grounds[i], passing[i])
	}
	return out
}

// Bracketed reads a sequence as En Harmonie prints it, a bracket under
// each two five toward its target: the chords of a two five in the
// tonalities it announces, the others as [Degrees] does. Dm7♭5 G7
// before Cm7 in E flat is II V, then VI.
//
// A plagal cadence is not bracketed: its IVm7 ♭VII7 reads on the tonic
// it concludes on. A two five that does not resolve and whose two sits
// on the ground's tonic is not bracketed either: E♭m7 A♭7 in E flat is
// Im7 IV7, a borrowed tonic.
func Bracketed(c Changes, blocks []Block, passing []int, sensed []Sensed) []Degree {
	out := Degrees(c, passing, sensed)
	for _, b := range blocks {
		if len(b.Announced) == 0 || b.Two < 0 || b.Kind == harmony.PlagalApproach ||
			borrowsTonic(c, b, sensed[b.Two].Ground) {
			continue
		}
		for _, i := range []int{b.Two, b.Sus, b.Five} {
			if i >= 0 && passing[i] == 0 {
				out[i] = degreeOf(c.Chords[i], b.Announced, 0)
			}
		}
	}
	return out
}
