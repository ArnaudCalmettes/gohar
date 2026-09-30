package harmony

import "fmt"

// A Tetrachord is a cell of four notes, designated by the three steps
// that separate them.
//
// # Named by the source, designated by its steps otherwise
//
// Ten tetrachords have a name, the ten that En Harmonie (tome 2,
// « Récapitulatif des modes ») finds in the five systems, declared
// below. Every other one is designated by its steps and nothing else:
// the tetrachord 1 3 2, the tetrachord 4 1 1. The name comes from the
// source, never from this library. Naming the rest would mean borrowing
// a vocabulary from another musical culture, or leaning on a lineage
// that does not hold up, and a shape that has no name in the practice
// this library serves is better left without one than given a false
// one.
//
// # Why the literal reads as the steps
//
// Each step occupies four bits, the first step highest, so a hexadecimal
// literal spells the tetrachord the way it is said aloud: 0x131 is the
// tetrachord 1 3 1. The type stays a sixteen bit integer, so two
// tetrachords compare in one instruction and can be declared as
// constants, which an array could not.
//
// A step of ten semitones or more would print as a letter. No
// tetrachord met in practice comes near it.
//
// # Invariant
//
// Every step is at least one semitone, and the four bits above the
// three steps are zero. A step of zero would be a doubled note, not a
// tetrachord.
type Tetrachord uint16

// The ten tetrachords En Harmonie names, in the order of its table.
// Their words, in any language, belong to the naming package; these
// identifiers only designate them in code, as the scale constants do.
//
// The last four are named after one of the first six with an altered
// degree, counted from the major scale as mode names are: the phrygian
// 𝄫3 lowers the E♭ of the phrygian tetrachord of C to E𝄫.
const (
	// The natural system, and the only ones the ear game asks for yet.
	TetrachordMajor    Tetrachord = 0x221
	TetrachordMinor    Tetrachord = 0x212
	TetrachordPhrygian Tetrachord = 0x122
	TetrachordLydian   Tetrachord = 0x222

	// Spans a diminished fourth: melodic minor.
	TetrachordDiminished Tetrachord = 0x121

	// Spans a perfect fourth: harmonic minor and harmonic major.
	TetrachordHarmonic Tetrachord = 0x131

	// Span an augmented fourth: harmonic minor.
	TetrachordLydianSharp2 Tetrachord = 0x312
	TetrachordMinorSharp4  Tetrachord = 0x213

	// Span a perfect fourth: double harmonic major.
	TetrachordPhrygianDoubleFlat3 Tetrachord = 0x113
	TetrachordMajorSharp2         Tetrachord = 0x311
)

// NewTetrachord builds a tetrachord from its three steps.
//
// Refuses a step below one semitone, which would double a note, and one
// above fifteen, which the encoding cannot hold.
func NewTetrachord(a, b, c Semitones) (Tetrachord, error) {
	for _, s := range [3]Semitones{a, b, c} {
		if s < 1 || s > 15 {
			return 0, fmt.Errorf("gohar: tetrachord step %d is out of range", s)
		}
	}
	return tetrachordOf(a, b, c), nil
}

// tetrachordOf packs three steps already known to be valid.
func tetrachordOf(a, b, c Semitones) Tetrachord {
	return Tetrachord(uint16(a)<<8 | uint16(b)<<4 | uint16(c))
}

// Steps returns the three steps of `t`, from its lowest note up.
func (t Tetrachord) Steps() (a, b, c Semitones) {
	return Semitones(t >> 8 & 0xf), Semitones(t >> 4 & 0xf), Semitones(t & 0xf)
}

// Span returns the distance from the lowest note of `t` to its highest.
//
// Five semitones for a tetrachord that spans a perfect fourth, which is
// the classical case. Not every tetrachord does: the lydian spans a
// tritone, and some shapes found in the altered systems span only a
// major third.
func (t Tetrachord) Span() Semitones {
	a, b, c := t.Steps()
	return a + b + c
}

// IsFourth reports whether `t` spans a perfect fourth.
func (t Tetrachord) IsFourth() bool {
	return t.Span() == 5
}

// IsValid reports whether `t` holds three steps of at least one semitone
// and nothing above them.
func (t Tetrachord) IsValid() bool {
	a, b, c := t.Steps()
	return t>>12 == 0 && a >= 1 && b >= 1 && c >= 1
}

// String returns the steps of `t`, as in 1 3 1. This is its designation,
// not a name, which is why it may live here.
func (t Tetrachord) String() string {
	a, b, c := t.Steps()
	return fmt.Sprintf("%d %d %d", a, b, c)
}

// A TetrachordSplit reads a heptatonic scale as two tetrachords and the
// space between them.
type TetrachordSplit struct {
	// Lower runs from the tonic up through the fourth note.
	Lower Tetrachord

	// Upper runs from the fifth note up to the octave.
	Upper Tetrachord

	// Gap is the space between the two, from the top of the lower
	// tetrachord to the bottom of the upper: the inter-tetrachordal
	// space.
	//
	// A whole tone when both tetrachords span a fourth. Less when one
	// of them is wider, as a lydian tetrachord reaching up to the
	// tritone leaves only a semitone. More when one is narrower. In
	// every case the two spans and the gap add up to an octave.
	Gap Semitones
}

// Tetrachords reads `p` as two tetrachords, and reports whether `p` is
// heptatonic, which is the only case where the reading applies.
//
// # A forced split, and what it is worth
//
// Seven notes leave no choice: four go to the lower tetrachord and
// three, with the octave, to the upper. So the reading always succeeds
// on a heptatonic pattern.
//
// On the 35 modes of the five systems, it always lands on one of the
// ten tetrachords En Harmonie names: that is where its table comes
// from. The harmonic minor is a minor tetrachord under a harmonic one,
// the locrian a phrygian under a lydian. Not every half spans a fourth,
// which is why the gap between them varies. On any other heptatonic
// pattern, the shapes it produces may have no name, and come back
// designated by their steps, which is all this reading can honestly say
// about them.
func (p ScalePattern) Tetrachords() (TetrachordSplit, bool) {
	if !p.IsHeptatonic() || !p.Contains(0) {
		return TetrachordSplit{}, false
	}

	var o [7]Semitones
	i := 0
	for _, n := range p.Offsets() {
		o[i] = n
		i++
	}

	return TetrachordSplit{
		Lower: tetrachordOf(o[1]-o[0], o[2]-o[1], o[3]-o[2]),
		Upper: tetrachordOf(o[5]-o[4], o[6]-o[5], SemitonesPerOctave-o[6]),
		Gap:   o[4] - o[3],
	}, true
}
