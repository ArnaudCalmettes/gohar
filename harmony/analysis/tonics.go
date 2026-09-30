package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// What a chord can be the tonic of, and how long it holds: the words
// the phrases and the sensed tonic share.

// roles tells, for each change, the block it is the target of and the
// block it belongs to, -1 for none.
type roles struct {
	target, member []int
}

func rolesOf(c Changes, blocks []Block) roles {
	r := roles{make([]int, len(c.Chords)), make([]int, len(c.Chords))}
	for i := range c.Chords {
		r.target[i], r.member[i] = -1, -1
	}
	for n, b := range blocks {
		if b.Target >= 0 {
			r.target[b.Target] = n
		}
		for _, i := range []int{b.Two, b.Sus, b.Five} {
			if i >= 0 {
				r.member[i] = n
			}
		}
	}
	return r
}

// tonicOf returns the tonalities a change is the tonic of, nil when it
// cannot be one: a chord that can be a tonic (a triad, maj7, 6, m6,
// m(maj7), not m7 nor a seventh), in root position or with its third in
// the bass. With its fifth in the bass, it is read for now as a chord
// over a pedal: F/C C in My Way is the IV of C over a tonic pedal. With
// any other bass, E♭maj7/F over the pedal of The Look Of Love, the bass
// is what holds.
func tonicOf(change Change) []harmony.Tonality {
	ch := change.Chord
	if change.Silent {
		return nil
	}
	if above := (int(change.Bass) - int(ch.Root) + 12) % 12; change.Inverted() && above != 3 && above != 4 {
		return nil
	}
	switch ch.Pattern.Tetrad() {
	case harmony.ChordMajorTriad, harmony.ChordMajorSeventh, harmony.ChordMajorSixth:
		return MajorTonalities(ch.Root)
	case harmony.ChordMinorTriad, harmony.ChordMinorSixth, harmony.ChordMinorMajorSeventh:
		return MinorTonalities(ch.Root)
	}
	return nil
}

// minorSeventhTonic returns the minor tonalities of a m7 that a minor
// cadence resolves on by its V, when it is not itself the two of a
// block: nil otherwise. Charts write the minor tonic m7 far more often
// than it is played (m6, m(maj7)): Cm7 | Dm7♭5 G7♭9 | Cm7 in Softly, As
// In A Morning Sunrise. A plagal IV7 does not make a m7 a tonic: C7
// Gm7 in Honeysuckle Rose is a V going back to its two. An aeolian
// cadence does: Amaj7 B7 C♯m7 in Guile's Theme installs C♯ minor.
func minorSeventhTonic(c Changes, blocks []Block, r roles, i int) []harmony.Tonality {
	ch := c.Chords[i]
	n := r.target[i]
	if n < 0 || r.member[i] >= 0 || ch.Silent || ch.Inverted() ||
		blocks[n].Kind == harmony.PlagalApproach && !blocks[n].Aeolian ||
		ch.Chord.Pattern.Tetrad() != harmony.ChordMinorSeventh {
		return nil
	}
	if a := blocks[n].Announced; len(a) > 0 && a[0].Tonic() == ch.Chord.Root && ModesOf(a)&Minor != 0 {
		return MinorTonalities(ch.Chord.Root)
	}
	return nil
}

// tonicAt tells what change i is heard as the tonic of: a tonic chord,
// or a m7 that a minor cadence within the chorus resolves on. A cadence
// across the loop, a turnaround toward a first chord in m7, does not
// count: as the chord heard, the first chord of Fly Me To The Moon or
// All The Things You Are is a VI or a II, not a tonic.
func tonicAt(c Changes, blocks []Block, r roles, i int) []harmony.Tonality {
	if t := tonicOf(c.Chords[i]); t != nil {
		return t
	}
	if n := r.target[i]; n >= 0 && blocks[n].Five > i {
		return nil
	}
	return minorSeventhTonic(c, blocks, r, i)
}

// cadencedAt returns the tonic a cadence within the chorus leads change
// i to, and that cadence; nil when none does. A turnaround across the
// loop leads nowhere: going round again is no new arrival.
func cadencedAt(c Changes, blocks []Block, r roles, i int) ([]harmony.Tonality, *Block) {
	n := r.target[i]
	if n < 0 || blocks[n].Five > i || len(blocks[n].Announced) == 0 {
		return nil, nil
	}
	t := tonicAt(c, blocks, r, i)
	if t == nil || !sameTonic(blocks[n].Announced, t) {
		return nil, nil
	}
	return t, &blocks[n]
}

// opensOn returns the tonic the sequence opens on, the one a phrase can
// come back to, nil when the first chord cannot be a tonic. A m7 opens
// on a tonic when a minor cadence resolves on that chord somewhere in
// the chorus: the Am7 of Fly Me To The Moon, which E7 comes back to at
// bar 8; not the Gm7 of Honeysuckle Rose or the Dm7 of Satin Doll, a
// two that nothing ever resolves on.
func opensOn(c Changes, blocks []Block, r roles) []harmony.Tonality {
	o := opening(c)
	if o < 0 {
		return nil
	}
	open := c.Chords[o]
	if t := tonicOf(open); t != nil {
		return t
	}
	if open.Chord.Pattern.Tetrad() != harmony.ChordMinorSeventh || r.member[o] >= 0 {
		return nil
	}
	for i, ch := range c.Chords {
		if ch.Chord == open.Chord && minorSeventhTonic(c, blocks, r, i) != nil {
			return MinorTonalities(open.Chord.Root)
		}
	}
	return nil
}

// heldTo returns the last change of the tonic held from change i: the
// chords on the same root and third after it, tonic chords or, for a
// minor tonic, its m7 (Gm6 | Gm6, the line of Dm Dm(maj7) Dm7 Dm6 in
// In a Sentimental Mood; not Gm6 | G7, nor E♭maj7 | E♭m6 in Candy, the
// IV turning minor).
func heldTo(c Changes, i int) int {
	first := c.Chords[i].Chord
	k := i + 1
	for ; k < len(c.Chords); k++ {
		ch := c.Chords[k]
		keeps := ch.Chord.Pattern == first.Pattern ||
			isMinor(ch.Chord.Pattern) == isMinor(first.Pattern) && (tonicOf(ch) != nil ||
				isMinor(first.Pattern) && ch.Chord.Pattern.Tetrad() == harmony.ChordMinorSeventh)
		if ch.Silent || ch.Chord.Root != first.Root || !keeps {
			break
		}
	}
	return k - 1
}

// held returns how long a tonic holds from change i (see [heldTo]).
func held(c Changes, i int) Ticks {
	var d Ticks
	for k := i; k <= heldTo(c, i); k++ {
		d += c.Chords[k].Length
	}
	return d
}

// barOf returns the length of a bar, four beats when the changes have
// no bars.
func barOf(c Changes) Ticks {
	if len(c.Bars) > 1 {
		return c.Bars[1] - c.Bars[0]
	}
	return 4 * TicksPerBeat
}

// opening returns the first change that sounds, -1 when none does.
func opening(c Changes) int {
	for i, ch := range c.Chords {
		if !ch.Silent {
			return i
		}
	}
	return -1
}

// isMinor reports whether a chord has a minor third and no major one.
func isMinor(p harmony.ChordPattern) bool {
	return p.HasOffset(3) && !p.HasOffset(4)
}

func sameTonic(a, b []harmony.Tonality) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a[0].Tonic() == b[0].Tonic()
}

// sameTonality reports whether two sets of tonalities have the same
// tonic and the same mode.
func sameTonality(a, b []harmony.Tonality) bool {
	return sameTonic(a, b) && (a == nil || ModesOf(a) == ModesOf(b))
}

// onTonic reports whether a change sits on the tonic of tonalities with
// their third: the tonic changing colour, Gm7 or Gm(maj7) on G minor.
func onTonic(ts []harmony.Tonality, ch Change) bool {
	return ts != nil && !ch.Silent && ch.Chord.Root == ts[0].Tonic() &&
		isMinor(ch.Chord.Pattern) == (ModesOf(ts) == Minor)
}
