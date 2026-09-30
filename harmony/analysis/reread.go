package analysis

import (
	"slices"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// Reread reads the plagal blocks again on the ground heard at their V,
// once a first hearing has installed it, and returns the approaches and
// blocks re-read, leaving those given as they are.
//
// Read alone, an X7 before a minor chord a whole tone above it, or a
// fourth below, is plagal: G7 Am7 the ♭VII7 Im of A minor, G7 Dm7 the
// IV7 Im of D melodic minor. Nothing in the two chords tells otherwise.
// The ground does: where C is heard, G7 is its V, and it is not
// preparing Am7 or Dm7 as their subdominant. So does the chord before,
// when the V goes back to the two it came from: Em7 A7 Em7 is a two
// five of D played again, whatever the ground.
//
//   - G7 Am7 is a deceptive cadence, the V going to its VI ("cadence
//     rompue", V – … in En Harmonie, tome 1, chapter 8: C7 Dm7).
//   - G7 Dm7 is a V going back to its two, the two five played again of
//     Satin Doll, Dm7 G7 | Dm7 G7.
//
// Either way the block becomes a V that does not resolve, announcing
// its tonic, with its two when it has one. It stays a block even
// without its two, as C7 Dm7 in the book: a V whose tonic is heard
// does not go nowhere, it promises that tonic. Where the ground is the
// minor chord itself, the plagal reading stands: Fm7 B♭7 Fm7 in Mas Que
// Nada is the dorian I IV7 of F minor. So it does in a vamp, the two
// chords alternating over more than four bars, and in an aeolian
// cadence, its ♭VImaj7 before it: Amaj7 B7 C♯m7 in Guile's Theme is no
// V of E going to its VI.
//
// The ground is the one of the first hearing, which reads each change
// in turn, or before it installs one, the tonality the tune is analysed
// in, as Sense has it: a second hearing on the blocks re-read is as
// live as the first. [Hear] runs both.
func Reread(c Changes, kinds []harmony.ApproachKind, blocks []Block, sensed []Sensed, tune []harmony.Tonality) ([]harmony.ApproachKind, []Block) {
	kinds, blocks = slices.Clone(kinds), slices.Clone(blocks)
	for n, b := range blocks {
		if b.Kind != harmony.PlagalApproach || b.Aeolian || b.Target < 0 || !isDominant(c.Chords[b.Five].Chord.Pattern) {
			continue // a IVm7 is no V
		}
		ground := sensed[b.Five].Ground
		if ground == nil {
			ground = tune
		}
		five, next := c.Chords[b.Five].Chord, c.Chords[b.Target].Chord
		// The VI and the II are minor with their fifth: a m7♭5 or a
		// dim7 there is something else. And a minor chord that is the
		// ground itself is its tonic.
		if !isMinor(next.Pattern) || !next.Pattern.HasOffset(7) || ground != nil && next.Root == ground[0].Tonic() {
			continue
		}
		if !ofTheGround(five, next, ground) && !playedAgain(c, b) {
			continue
		}
		kinds[b.Five] &^= harmony.PlagalApproach
		b.Kind, b.Target = 0, -1 // as a two five going elsewhere
		if b.Two < 0 {
			p := c.Prev(b.Five)
			if b.Sus >= 0 {
				p = c.Prev(b.Sus)
			}
			if p >= 0 && p != b.Five && isTwoOf(c.Chords[p], c.Chords[b.Five]) {
				b.Two = p
			}
		}
		b.Announced = announce(c, b)
		blocks[n] = b
	}
	return kinds, blocks
}

// ofTheGround reports whether five is the V of a major ground, and next
// its VI or its II.
func ofTheGround(five, next harmony.Chord, ground []harmony.Tonality) bool {
	if ground == nil || ModesOf(ground)&Major == 0 {
		return false
	}
	tonic := ground[0].Tonic()
	return five.Root == tonic.Transpose(7) && (next.Root == tonic.Transpose(9) || next.Root == tonic.Transpose(2))
}

// playedAgain reports whether the block's V goes back to the two it
// came from, whatever the ground: Em7 A7 | Em7 A7 in Satin Doll, a two
// five of D played again while C is heard.
func playedAgain(c Changes, b Block) bool {
	before := c.Prev(b.Five)
	if b.Sus >= 0 {
		before = c.Prev(b.Sus)
	}
	if before < 0 || c.Chords[before].Silent {
		return false
	}
	two, next, five := c.Chords[before].Chord, c.Chords[b.Target].Chord, c.Chords[b.Five].Chord
	return two.Root == next.Root && isMinor(two.Pattern) && next.Root == five.Root.Transpose(7) &&
		!vamp(c, two.Root, five.Root, b.Five)
}

// vamp reports whether a minor chord on two and a dominant on five
// alternate around change at over more than four bars: a vamp, modal,
// rather than a two five played again.
func vamp(c Changes, two, five harmony.PitchClass, at int) bool {
	alternates := func(i int) bool {
		ch := c.Chords[i]
		return !ch.Silent && (ch.Chord.Root == two && isMinor(ch.Chord.Pattern) || ch.Chord.Root == five && isDominant(ch.Chord.Pattern))
	}
	from, to := at, at
	for from > 0 && alternates(from-1) {
		from--
	}
	for to+1 < len(c.Chords) && alternates(to+1) {
		to++
	}
	start, end := c.Chords[from].Start, c.Chords[to].Start+c.Chords[to].Length
	if len(c.Bars) == 0 {
		return end-start > 4*4*TicksPerBeat
	}
	return c.Bar(end-1)-c.Bar(start)+1 > 4
}
