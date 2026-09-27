package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// PassingChords finds the diminished passing chords of a sequence: for
// each change, +1 when it is one on a bass walking up, -1 on a bass
// walking down, 0 otherwise.
//
// A passing chord links two chords by a chromatic bass: its bass is a
// semitone from the bass before it and a semitone from the bass after,
// in the same direction. Up, I ♯Idim7 II (Fmaj7 F♯dim7 Gm7, Mean to
// Me); down, I/3 ♭IIIdim7 II (B♭/D D♭dim7 Cm7, Someday My Prince Will
// Come), where the walk exists only because the first chord is
// inverted.
//
// Three changes and their basses: what [harmony.ApproachOf] cannot see
// from two chords. Going up, the diminished chord is also the dominant
// without its root of the chord it leads to, and the two readings stand
// together, each from its own calculation.
//
// Only diminished chords, as En Harmonie teaches it. A chromatic bass
// walking under other chords, chromatic dominants and inversions in It
// Never Entered My Mind, is a bass line, another reading to come.
func PassingChords(c Changes) []int {
	out := make([]int, len(c.Chords))
	for i, ch := range c.Chords {
		prev, next := c.Prev(i), c.Next(i)
		if ch.Silent || prev < 0 || next < 0 || ch.Chord.Pattern.Tetrad() != harmony.ChordDiminishedSeventh {
			continue
		}
		before, after := c.Chords[prev], c.Chords[next]
		if before.Silent || after.Silent {
			continue
		}
		switch {
		case before.Bass.Transpose(1) == ch.Bass && ch.Bass.Transpose(1) == after.Bass:
			out[i] = 1
		case before.Bass.Transpose(-1) == ch.Bass && ch.Bass.Transpose(-1) == after.Bass:
			out[i] = -1
		}
	}
	return out
}
