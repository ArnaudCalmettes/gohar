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
