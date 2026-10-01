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
// fifths. The step is measured between the twos, or between the fives,
// as it is heard and as players describe it, never from a V to the next
// two: Dm7 G7 | Cm7 F7, another II-V a tone below, is a step, whatever
// G7 falling a fifth to Cm7. Following the cycle of fifths is the
// bridge of the rhythm changes, the twos a fifth apart: Am7 D7 | Dm7 G7
// | Gm7 C7 | Cm7 F7.
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

	// WholeStep: a tone between the two IIs, Dm7 G7 | Em7 A7 up, or
	// Dm7 G7 | Cm7 F7 down, where Cm7 is the I the first cadence could
	// have concluded on and the II of the next: Em7♭5 A7 Dm7 G7 Cm7 F7
	// in Confirmation.
	WholeStep

	// Fifths: the cycle of fifths, the twos a fifth apart, and the fives
	// too: Am7 D7 | Dm7 G7, the bridge of the rhythm changes with its
	// twos, or Em7 A7 | Am7 D7 in Satin Doll; or a chain of dominants,
	// each the V of the next, D7 G7 C7 F7.
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
// The step between the twos comes first, and the step between the fives
// only when the twos give none: Bm7 E7 | Am7 D7 in Straight Street is a
// step down, whatever E7 falling a fifth to Am7, and Gm7 G♭7 | Fm7 E7 in
// Autumn Leaves is a tone down, its fives chromatic dominants.
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
		if root(next) == root(first) {
			continue // a II-V played again
		}
		s := step(root(next) - root(first))
		if after := c.Next(next); s == 0 && after > next && !c.Chords[after].Silent {
			s = step(root(after) - root(five))
		}
		if s != 0 {
			out = append(out, Link{From: first, To: next, Step: s})
		}
	}
	return append(out, dominants(c, kinds)...)
}

// step names the gap between two roots, `d` semitones up: a half tone,
// a tone, or a fifth down, the cycle of fifths. Zero for a gap the book
// does not name.
func step(d int) Step {
	switch (d%12 + 12) % 12 {
	case 1, 11:
		return HalfStep
	case 2, 10:
		return WholeStep
	case 5:
		return Fifths
	}
	return 0
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
