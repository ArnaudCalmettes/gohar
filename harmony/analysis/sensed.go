package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Sensed is the tonic an ear expects once a chord has sounded, from
// what came before it and nothing after: the tonique pressentie of
// docs/grilles.md. At the third chord of Tenderly, two bars of E♭maj7
// A♭7 have installed E flat, and E♭m7 is heard as the tonic changing
// colour, not as the two of D flat.
//
// [Sense] reads it from left to right, the other way from the rest of
// the analysis: the region is what one concludes afterwards, the
// sensed tonic what one hears now. The gap between what it awaited and
// what comes is the surprise.
//
// Each field is a set of tonalities on one tonic, as a block announces
// them, and nil when there is none.
type Sensed struct {
	// Start is the first tonic installed, kept apart from the ground to
	// hear the return home after a bridge that modulated.
	Start []harmony.Tonality

	// Ground is the installed tonic, the one degrees count from, even
	// when a cadence tonicises another degree.
	Ground []harmony.Tonality

	// Local is the tonic a cadence has just tonicised, when it is not
	// the ground; it lasts while the chords after it hold in it.
	Local []harmony.Tonality

	// Awaited is what a cadence being played announces: the tonalities
	// of the block the chord belongs to.
	Awaited []harmony.Tonality
}

// Sense reads the sensed tonic at each change, starting from the
// tonalities of a key signature, nil when there is none.
//
// # What installs a tonic
//
// The signature gives the ground to start from. Without one, the first
// chord gives it when it can be a tonic chord (a major or minor triad,
// maj7, 6, m6, m(maj7), but not m7, too often a two), and else the
// first cadence that resolves. A cadence that resolves elsewhere than
// on the ground makes its target a local tonic.
//
// For now a local tonic never becomes the ground: when it does, after
// how long and how many cadences, is the threshold of a modulation,
// still to set (docs/grilles.md). [Tune] reads the tonality of the tune
// from the last change, which covers a tune ending away from where it
// started.
//
// # A two on the ground is its I
//
// A two five that does not resolve and whose two sits on the ground's
// tonic awaits nothing: its two is the tonic, borrowed from another
// scale. E♭m7 A♭7 in E flat is I then IV7, not a two five of D flat.
//
// # Local
//
// Each change sees the state left by those before it and the blocks,
// whose tonalities a target decides between: one change ahead, as the
// analysis of a performance allows.
func Sense(c Changes, blocks []Block, signature []harmony.Tonality) []Sensed {
	target := make([]int, len(c.Chords))
	member := make([]int, len(c.Chords))
	for i := range target {
		target[i], member[i] = -1, -1
	}
	for n, b := range blocks {
		if b.Target >= 0 {
			target[b.Target] = n
		}
		for _, i := range []int{b.Two, b.Sus, b.Five} {
			if i >= 0 {
				member[i] = n
			}
		}
	}

	out := make([]Sensed, len(c.Chords))
	s := Sensed{Start: signature, Ground: signature}
	for i, ch := range c.Chords {
		s.Awaited = nil
		if ch.Silent {
			out[i] = s
			continue
		}
		if s.Ground == nil && i == 0 {
			s.Ground = tonicOf(ch.Chord)
		}
		if n := target[i]; n >= 0 && len(blocks[n].Announced) > 0 {
			t := blocks[n].Announced
			switch {
			case s.Ground == nil:
				s.Ground, s.Local = t, nil
			case t[0].Tonic() == s.Ground[0].Tonic():
				s.Local = nil
			default:
				s.Local = t
			}
		} else if s.Local != nil && !holds(s.Local, ch.Chord.Set()) {
			s.Local = nil
		}
		if s.Start == nil {
			s.Start = s.Ground
		}
		if n := member[i]; n >= 0 && !borrowsTonic(c, blocks[n], s.Ground) {
			s.Awaited = blocks[n].Announced
		}
		out[i] = s
	}
	return out
}

// Tune reads the tonality of the tune at its last change: the local
// tonic when the tune ends on one, else the ground, in the major or in
// the three minors as the last chord's third says.
func Tune(c Changes, sensed []Sensed) []harmony.Tonality {
	for i := len(c.Chords) - 1; i >= 0; i-- {
		if c.Chords[i].Silent {
			continue
		}
		ts := sensed[i].Local
		if ts == nil {
			ts = sensed[i].Ground
		}
		if ts == nil {
			return nil
		}
		if t := tonicOf(c.Chords[i].Chord); t != nil && t[0].Tonic() == ts[0].Tonic() {
			return t
		}
		return ts
	}
	return nil
}

// tonicOf returns the tonalities a chord is the tonic of, nil when it
// cannot be one.
func tonicOf(ch harmony.Chord) []harmony.Tonality {
	switch ch.Pattern.Tetrad() {
	case harmony.ChordMajorTriad, harmony.ChordMajorSeventh, harmony.ChordMajorSixth:
		return MajorTonalities(ch.Root)
	case harmony.ChordMinorTriad, harmony.ChordMinorSixth, harmony.ChordMinorMajorSeventh:
		return MinorTonalities(ch.Root)
	}
	return nil
}

// borrowsTonic reports whether a block is a two five that does not
// resolve and whose two sits on the ground's tonic: its two is then
// the tonic borrowed from another scale, not a two.
func borrowsTonic(c Changes, b Block, ground []harmony.Tonality) bool {
	return b.Two >= 0 && b.Target < 0 && ground != nil &&
		c.Chords[b.Two].Chord.Root == ground[0].Tonic()
}
