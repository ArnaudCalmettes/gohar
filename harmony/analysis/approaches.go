package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// Approaches reads, for each change, how it prepares the change that
// follows it: [harmony.ApproachOf] across the whole sequence, and across
// the bar line where it loops.
//
// A silent change prepares nothing and is prepared by nothing. The last
// change of a sequence that does not loop prepares nothing yet: in a
// performance, what it prepares is not known until the next chord
// sounds.
func Approaches(c Changes) []harmony.ApproachKind {
	out := make([]harmony.ApproachKind, len(c.Chords))
	for i, ch := range c.Chords {
		next := c.Next(i)
		if next < 0 || ch.Silent || c.Chords[next].Silent {
			continue
		}
		out[i] = harmony.ApproachOf(ch.Chord, c.Chords[next].Chord)
	}
	return out
}

// A Chain is a run of changes each preparing the next, and the change
// they lead to: A7 Dm7 G7 before Cmaj7, a secondary dominant, a two and
// a dominant.
//
// It is read from right to left, as En Harmonie teaches: from the arrival,
// back through its preparations for as long as each one prepares the
// next. The arrival is a change that prepares nothing itself, so a
// chain holds everything that leads to it, En Harmonie's nested readings
// included: A7 is the dominant of Dm7, and Dm7 G7 the two five of C.
// Cutting the chain into those blocks is the next layer's business.
type Chain struct {
	Steps   []int // the preparing changes, in the order they sound
	Arrival int
}

// Chains finds every chain of a sequence, given its approaches, in the
// order of their arrivals.
//
// A looping sequence where every change prepares the next, a cycle of
// dominants that never lands, has no arrival and no chain.
func Chains(c Changes, kinds []harmony.ApproachKind) []Chain {
	var out []Chain
	for j := range c.Chords {
		if kinds[j] != harmony.NoApproach || c.Chords[j].Silent {
			continue
		}
		var steps []int
		for p := c.Prev(j); p >= 0 && p != j && kinds[p] != harmony.NoApproach; p = c.Prev(p) {
			if len(steps) == len(c.Chords) {
				break
			}
			steps = append(steps, p)
		}
		if len(steps) == 0 {
			continue
		}
		for l, r := 0, len(steps)-1; l < r; l, r = l+1, r-1 {
			steps[l], steps[r] = steps[r], steps[l]
		}
		out = append(out, Chain{Steps: steps, Arrival: j})
	}
	return out
}
