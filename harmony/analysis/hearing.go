package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Hearing is a sequence heard as an analyst hears it, twice: a first
// time to install its tonic, a second time with the blocks read again
// on that tonic (see [Reread]).
type Hearing struct {
	Kinds   []harmony.ApproachKind // how each change prepares the next
	Blocks  []Block
	Phrases []Phrase

	// Heard is the tonality the first hearing concludes (see [Tune]),
	// Tune the one the degrees are counted in, the same unless the
	// reader chooses another.
	Heard, Tune []harmony.Tonality

	Sensed []Sensed
}

// Hear hears changes twice. choose picks the tonality to count the
// degrees in, given the one the first hearing concludes: the one the
// app declares, or one the reader imposes. A nil choose keeps the one
// heard.
func Hear(c Changes, choose func(heard []harmony.Tonality) []harmony.Tonality) Hearing {
	h := Hearing{Kinds: Approaches(c)}
	h.Blocks = Blocks(c, h.Kinds)
	h.Phrases = Phrases(c, h.Blocks)
	h.Heard = Tune(c, h.Phrases)
	h.Tune = h.Heard
	if choose != nil {
		h.Tune = choose(h.Heard)
	}
	h.Kinds, h.Blocks = Reread(c, h.Kinds, h.Blocks, Sense(c, h.Blocks, h.Phrases, h.Tune), h.Tune)
	h.Phrases = Phrases(c, h.Blocks)
	h.Sensed = Sense(c, h.Blocks, h.Phrases, h.Tune)
	return h
}
