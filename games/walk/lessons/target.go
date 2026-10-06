package lessons

import "slices"

// A Verdict is how the keys held answer what a step asks.
type Verdict int

const (
	Miss    Verdict = iota // a wrong note: the casserole
	Partial                // on the way: a pair's first key, the second to come
	Hit                    // the answer
)

// A Target is what a step asks for: a note in any octave, the two black
// keys of a pair together. The player aims at it deliberately, so any
// key off it is a miss (see "Jouer délibérément" in docs/debutants.md).
type Target interface {
	// Judge tells how `held`, the keys down now in the order they went
	// down, answers the target.
	Judge(held []int) Verdict

	// Hint is the pitch classes to show after a miss.
	Hint() []int
}

// pc is the pitch class of MIDI key `k`: 0 for C, up to 11 for B.
func pc(k int) int { return (k%12 + 12) % 12 }

// Note asks for one note, in any octave: "trouve un do". Its pitch
// class is 0 for C, up to 11 for B.
type Note int

// The notes, by their names in both notations, for the lessons'
// scripts.
const (
	Do  Note = 0
	Re  Note = 2
	Mi  Note = 4
	Fa  Note = 5
	Sol Note = 7
	La  Note = 9
	Si  Note = 11
)

// Judge looks at the key that went down last only: a key still held,
// the C under the left hand while the right one plays, is no mistake.
func (n Note) Judge(held []int) Verdict {
	if len(held) > 0 && pc(held[len(held)-1]) == int(n) {
		return Hit
	}
	return Miss
}

func (n Note) Hint() []int { return []int{int(n)} }

// A Group asks for a group of black keys pressed together: the pair,
// C♯ and D♯, or the trio, F♯, G♯ and A♯, in any octave. The keys need
// not go down at the very same time: until the last one, the group is
// Partial.
type Group []int

// The groups of black keys, by their pitch classes from the lowest.
var (
	Pair = Group{1, 3}
	Trio = Group{6, 8, 10}
)

func (g Group) Judge(held []int) Verdict {
	if len(held) == 0 || len(held) > len(g) {
		return Miss
	}
	base := -1 // the key of the group's lowest note, the same for every key held
	seen := map[int]bool{}
	for _, k := range held {
		i := slices.Index(g, pc(k))
		if i < 0 || seen[k] {
			return Miss
		}
		seen[k] = true
		b := k - (g[i] - g[0])
		if base >= 0 && b != base {
			return Miss
		}
		base = b
	}
	if len(held) == len(g) {
		return Hit
	}
	return Partial
}

func (g Group) Hint() []int { return g }
