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
// chord's tetrad, the whole 7alt for an altered chord (VI7alt, not the
// VI7♭5 its tetrad would say), and zero for a diatonic chord. A
// dominant seventh on another degree than V keeps its quality, diatonic
// as it may be: in A minor, G7 is VII7, not the dominant of the key.
//
// The degree counts in the tonality's own scale: in F minor, D♭maj7 is
// VI, not ♭VI, as En Harmonie's tables of VI V I write it (tome 1,
// chapter 8, "Fonctions des accords": G∆9 is the VI of B harmonic
// minor). A root outside every scale of the tonality is a degree raised
// or lowered, lowered by preference (♭II, ♭III, ♭VI, ♭VII) but raised
// below the fourth and the fifth (♯IV, ♯III in minor). A passing chord
// follows its bass instead: raised going up (♯Idim7, ♯Vdim7), lowered
// going down (♭IIIdim7).
//
// Writing 7alt is a choice no source settles, open to correction by an
// expert.
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
	// A borrowed chord keeps its quality, and so does a suspension,
	// diatonic as it is: the V7sus4 is a subdominant chord, as the Bill
	// Evans Piano Academy teaches it, and writing V for it would make
	// it the V.
	p := ch.Chord.Pattern
	dominant := p.HasOffset(4) && p.HasOffset(10) && d.Number != 5
	switch {
	case p == sevenAlt:
		d.Quality = sevenAlt
	case !holds(ts, ch.Chord.Set()), isSus(p), dominant:
		d.Quality = p.Tetrad()
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

// sevenAlt is the chord of the altered mode, locrian ♭4: C7(♭5, ♭9,
// ♭10, ♭13), its ♭4 heard as the major third.
var sevenAlt, _ = harmony.NewChordPattern(0, 4, 6, 10, 13, 15, 20)

var qualities = map[harmony.ChordPattern]string{
	sevenAlt:                           "7alt",
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
//
// Nor is the first step of a march of two fives, whose V7 turns into
// the two of the next one on the same root: Cm7 F7 Fm7 B♭7 in E flat,
// bars 13 to 16 of Tenderly, is VI II7 II V. F7 is a two altered, as
// En Harmonie puts it, not the V of a B flat that never comes. Nor one
// whose V lands on the suspension of the next V: Dm7 G7 C7sus4 C7 in F
// is VI II7 V7sus4 V, the suspension a two hiding.
func Bracketed(c Changes, blocks []Block, passing []int, sensed []Sensed) []Degree {
	out := Degrees(c, passing, sensed)
	for _, b := range blocks {
		if len(b.Announced) == 0 || b.Two < 0 || b.Kind == harmony.PlagalApproach ||
			borrowsTonic(c, b, sensed[b.Two].Ground) || marches(c, blocks, b) {
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

// marches reports whether a two five is the first step of a march: its
// V7 turns into the minor seventh two of the next block, on the same
// root (F7 Fm7 B♭7), or lands on the suspension of the next V, a two
// hiding under the root of that V (Dm7 G7 C7sus4 C7 F in My Lucky
// Star, a three six two five of F).
func marches(c Changes, blocks []Block, b Block) bool {
	if b.Five < 0 {
		return false
	}
	n := c.Next(b.Five)
	if n < 0 {
		return false
	}
	next := c.Chords[n]
	for _, o := range blocks {
		switch {
		case b.Target < 0 && o.Two == n && next.Chord.Root == c.Chords[b.Five].Chord.Root &&
			next.Chord.Pattern.Tetrad() == harmony.ChordMinorSeventh,
			b.Target == n && o.Sus == n:
			return true
		}
	}
	return false
}

// isSus reports whether a chord has a fourth or a second in place of
// its third.
func isSus(p harmony.ChordPattern) bool {
	return !p.HasAny(harmony.IntMajorThird, harmony.IntMinorThird) &&
		p.HasAny(harmony.IntPerfectFourth, harmony.IntMajorSecond)
}
