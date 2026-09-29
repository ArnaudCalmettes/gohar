package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Link joins two II-Vs that follow each other, the first leading to
// the second rather than to its tonic: the consecutive II-Vs of En
// Harmonie (tome 1, chapter 8 §3.3, "Les enchaînements harmoniques
// fréquents"). Chaining the subdominant and the dominant of a key
// without resolving them breaks off, and the conclusion is awaited the
// more: a rebound from one cadence to the next, which be-bop is fond
// of.
//
// The book names the link by a step found from one bar to the next
// between one chord of each: a half tone, a tone, or the cycle of
// fifths, where the V of the first falls a fifth to the II of the
// second.
//
// Dominants chain the same way without their twos, each the V of the
// next: D7 G7 C7 F7, the bridge of the rhythm changes, by fifths; E7
// E♭7 D7 D♭7, by half tones, each the chromatic dominant of the next.
type Link struct {
	From, To int // the first chord of each: the two of a II-V, or a V
	Step     Step
}

// A Step is how a II-V follows the one before it.
type Step int

const (
	// HalfStep: a half tone between the two, chromatisms. Dm7 G7 |
	// E♭m7 A♭7, or Dm7 G7 | D♭m7 G♭7; C♯m7 F♯7 Cm7 F7 Bm7 E7 in Butch
	// And Butch.
	HalfStep Step = iota + 1

	// WholeStep: a tone, between the two IIs, Dm7 G7 | Em7 A7, or from
	// the V to the next II, Dm7 G7 | Am7 D7, Dm7 G7 | Fm7 B♭7.
	WholeStep

	// Fifths: the V falls a fifth to the next II, the cycle of fourths
	// or of fifths. In Dm7 G7 | Cm7 F7, Cm7 is the I the first cadence
	// could have concluded on, and the II of the next: Em7♭5 A7 Dm7 G7
	// Cm7 F7 in Confirmation.
	Fifths
)

var stepNames = [...]string{"", "½", "step", "5th"}

func (s Step) String() string {
	return stepNames[s]
}

// Links finds the consecutive II-Vs of a sequence: a two and its V,
// then right after them another two and its V, on another root. The
// II-Vs are read from the approaches, not the blocks: in Satin Doll,
// Dm7 G7 | Dm7 G7 | Em7 A7 | Em7 A7 goes up a tone whatever G7 before
// Dm7 is heard as. A II-V played again is no link but does not break
// the chain, and a step the book does not name is no link.
//
// The step between the twos comes first, the step from the V only when
// the twos give none: a chromatic dominant is always a half tone from
// the next two, and Gm7 G♭7 | Fm7 E7 in Autumn Leaves is a tone down.
func Links(c Changes, kinds []harmony.ApproachKind) []Link {
	var twos []int
	for i := range c.Chords {
		if n := c.Next(i); kinds[i].Has(harmony.TwoApproach) && n > i && !c.Chords[i].Silent && !c.Chords[n].Silent {
			twos = append(twos, i)
		}
	}
	root := func(i int) int { return int(c.Chords[i].Chord.Root) }
	var out []Link
	for k := 0; k+1 < len(twos); k++ {
		first, next := twos[k], twos[k+1]
		five := first + 1
		if next != c.Next(five) {
			continue
		}
		two := (root(next) - root(first) + 12) % 12
		fall := (root(next) - root(five) + 12) % 12
		var s Step
		switch {
		case two == 0:
		case fall == 5:
			s = Fifths
		case two == 1, two == 11:
			s = HalfStep
		case two == 2, two == 10:
			s = WholeStep
		case fall == 1, fall == 11:
			s = HalfStep
		case fall == 2, fall == 10:
			s = WholeStep
		}
		if s != 0 {
			out = append(out, Link{From: first, To: next, Step: s})
		}
	}
	return append(out, dominants(c, kinds)...)
}

// dominants finds the chains of dominants: a V followed by the dominant
// it prepares, a fifth below, or a half tone below as its chromatic
// dominant. G7 C7 in D7 G7 C7 F7; not G7 Cm7, a V going to a two.
func dominants(c Changes, kinds []harmony.ApproachKind) []Link {
	var out []Link
	for i, ch := range c.Chords {
		next := c.Next(i)
		if next <= i || ch.Silent || c.Chords[next].Silent || !isDominant(c.Chords[next].Chord.Pattern) {
			continue
		}
		switch {
		case kinds[i].Has(harmony.DominantApproach):
			out = append(out, Link{From: i, To: next, Step: Fifths})
		case kinds[i].Has(harmony.ChromaticApproach):
			out = append(out, Link{From: i, To: next, Step: HalfStep})
		}
	}
	return out
}
