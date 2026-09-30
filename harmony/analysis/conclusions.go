package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Conclusion is how a section of a tune ends: the cadence that
// concludes it, and the turnaround after it.
//
// # What the sources say
//
// « La cadence conclusive est le point d'arrivée d'une phrase
// harmonique », and in a tune that keeps to the carrure, the conclusive
// cadences mark the groups of bars (Siron, La partition intérieure,
// p. 390). In a phrase of 8 bars, it lands on a tonic chord, or one
// that stands for it. Two places:
//
//   - strong (« masculine »), on the last bar but one, the most frequent:
//     Let's Cool One lands on its I at bar 7;
//   - weak (« féminine »), on an even bar, most often the last: Lover
//     Man lands on its I at bar 8.
//
// The conclusive chord starts the cadence-boucle, the turnaround, that
// resolves on the first chord of the next section: « une anacrouse de
// la phrase harmonique suivante » (p. 353). En Harmonie says the same of
// the Ier degré that ends a tune (tome 1, chapter 8, p. 125).
//
// # In gohar
//
// The conclusive cadence of a section is the last cadence that resolves
// on a tonic sounding in its last three bars, cadencedAt's: a V-I, a
// plagal cadence, on a tonic chord or on a m7 it installs. The tonic may
// be reached before and held into them as a tonic chord: Sweet Sue lands
// on G6 at bar 5 and holds it to the end of its section. The Dm Dm(maj7)
// of Yesterdays at bar 5 goes on as Dm7 at bar 6, a line down to Bm7♭5
// E7: the tonic chord is over before the last bars. A section that ends
// on its V, a half cadence, or that goes through its last bars without
// resolving, concludes on none.
//
// One landing concludes nothing, though a V leads to it (see
// [leansOnSixth]): the VI of the minor tonic the turnaround goes back
// to, on a weak bar. It is heard as leaning on that chord for a bar,
// before the section ends open on its V.
type Conclusion struct {
	// Arrives is the change the conclusive cadence lands on, -1 when
	// the section concludes on none.
	Arrives int

	// Tonic is the tonic it lands on.
	Tonic []harmony.Tonality

	// Strong tells a cadence landing on the last bar but one from one
	// landing later or earlier.
	Strong bool

	// Loop is the first change of the turnaround, the chords after the
	// conclusive one and before the next section: -1 when the tonic
	// holds to the end of the section.
	Loop int
}

// Conclusions returns how each section ends, given the blocks of the
// changes (see [Conclusion]). A pickup, or a section shorter than 4
// bars, concludes on none.
func Conclusions(c Changes, sections []Section, blocks []Block) []Conclusion {
	r := rolesOf(c, blocks)
	out := make([]Conclusion, len(sections))
	for n, s := range sections {
		out[n] = Conclusion{Arrives: -1, Loop: -1}
		if s.Label == "-" || s.Bars < 4 {
			continue
		}
		end := s.From + s.Bars
		for i, ch := range c.Chords {
			bar := c.Bar(ch.Start)
			if bar < s.From || bar >= end || !soundsFrom(c, i, end-3) {
				continue
			}
			if t, b := cadencedAt(c, blocks, r, i); t != nil && !leansOnSixth(c, b, t, bar, end) {
				out[n] = Conclusion{Arrives: i, Tonic: t, Strong: bar == end-2, Loop: -1}
			}
		}
		if i := out[n].Arrives; i >= 0 {
			if k := heldTo(c, i) + 1; k < len(c.Chords) && c.Bar(c.Chords[k].Start) < end {
				out[n].Loop = k
			}
		}
	}
	return out
}

// leansOnSixth reports whether tonic `t`, that block `b` lands on at
// `bar` of a section ending before `end`, is only leant on: reached by
// a lone V, with no two, no sus and no plagal cadence, on a weak bar,
// and the VI of the minor tonic the next section opens on. Yesterdays
// walks down the cycle, D9 G13 C9 F13, onto B♭maj7 at bar 14 of 16,
// then Em7♭5 A7 goes back to Dm: the ear hears no end of phrase on
// B♭, the section ends open on A7. Our reading, from the ear: no source
// says it.
//
// Each condition keeps a real conclusion. Rosetta comes down the same
// kind of cycle onto F6, its home, but on the strong bar 15, before
// Bm7♭5 E7 goes to Am. Lover Man lands on Fmaj7 at bar 16, the VI of
// the Am of its bridge, by Gm7 C7, a two five.
func leansOnSixth(c Changes, b *Block, t []harmony.Tonality, bar, end int) bool {
	if b.Two >= 0 || b.Sus >= 0 || b.Kind == harmony.PlagalApproach || bar == end-2 {
		return false
	}
	next := nextOpening(c, end)
	return next >= 0 && tonicOf(c.Chords[next]) != nil && isMinor(c.Chords[next].Chord.Pattern) &&
		c.Chords[next].Chord.Root == t[0].Tonic().Transpose(4)
}

// nextOpening returns the first change of the bar `end`, the one the
// next section opens on, or across the loop the first chord of the
// chart: -1 when there is none.
func nextOpening(c Changes, end int) int {
	for i, ch := range c.Chords {
		if c.Bar(ch.Start) == end {
			return i
		}
	}
	if c.Loops {
		return opening(c)
	}
	return -1
}

// soundsFrom reports whether the tonic reached at change `i` still
// sounds as a tonic chord at bar `from` or after: the change itself, or
// one of the chords that hold it (see heldTo).
func soundsFrom(c Changes, i, from int) bool {
	for k := i; k <= heldTo(c, i); k++ {
		ch := c.Chords[k]
		if (k == i || tonicOf(ch) != nil) && c.Bar(ch.Start+ch.Length-1) >= from {
			return true
		}
	}
	return false
}
