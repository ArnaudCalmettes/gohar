package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Block is a cadence: [II] [sus4] V → target.
//
// The V is the one chord it needs: a chord that prepares the next as
// its dominant, its chromatic dominant or its diminished chord. Before
// it may come its suspension, the X7sus4 on the same root, and before
// that its two. The target is the chord after the V, when the V
// prepares it; the target is not part of the block, only pointed at.
//
// # Nested blocks
//
// Each chord belongs to one block at most, and a block may target a
// chord of another. A7 Dm7 G7 before Cmaj7 is two blocks, [Dm7 G7] →
// Cmaj7 and [A7] → Dm7: En Harmonie's reading, a V7/II then a two five
// one. The bebop double role follows with nothing more: in Em7 A7 Dm7
// G7 Cmaj7, Dm7 is the target of [Em7 A7] and the two of [Dm7 G7].
//
// # Without a target
//
// A two and a dominant that goes elsewhere still make a block: the two
// five without resolution of En Harmonie, which announces its tonality
// all the same, a fifth under its V. A dominant alone that prepares
// nothing is not a block.
//
// # Plagal
//
// A plagal cadence is a block too, its subdominant in the place of the
// V (Kind is harmony.PlagalApproach): the IV before the tonic, or the
// ♭VII7 with its IVm7 as its two (Fm7 B♭7 Cmaj7). It concludes and
// defines the tonality as a V-I does (see [Sense]). With its ♭VImaj7
// before it, the ♭VII7 I is an aeolian cadence (Aeolian).
//
// # Deceptive
//
// The deceptive cadence (cadence rompue), the V going to the VI of its
// tonality, is not a block's to tell: Dm7 G7 Am7 read alone is the
// aeolian IVm7 ♭VII7 Im7 of A minor, a plagal block, and it is the
// tonic of the passage, C, that makes it a V going to its VI. See
// [Reread], which reads the blocks again once that tonic is heard.
type Block struct {
	Two, Sus int // indices of the changes, -1 when absent
	Five     int
	Target   int                  // -1 when the V does not resolve
	Kind     harmony.ApproachKind // how the V prepares the target

	// The tonalities the block announces, on one tonic, whether or not
	// it lands there: in What Is This Thing Called Love, a two five
	// announces F harmonic minor and resolves on F major. See
	// [announce].
	Announced []harmony.Tonality

	// Aeolian tells a plagal ♭VII7 I that ends an aeolian cadence, its
	// ♭VImaj7 before it (see [AeolianCadence]). A modal cadence is
	// never read again as a V going to its VI (see [Reread]): Amaj7 B7
	// C♯m7 in Guile's Theme is no V of E.
	Aeolian bool
}

// fives are the approaches that make a V.
const fives = harmony.DominantApproach | harmony.ChromaticApproach | harmony.DiminishedApproach

// Blocks finds the blocks of a sequence, given its approaches, in the
// order of their V.
//
// Each block is read from its V backwards, over two changes at most:
// local, as the analysis of a performance needs.
func Blocks(c Changes, kinds []harmony.ApproachKind) []Block {
	var out []Block
	for f, five := range c.Chords {
		b := Block{Two: -1, Sus: -1, Five: f, Target: -1}
		switch {
		case five.Silent:
			continue
		case kinds[f]&fives != 0:
			b.Target, b.Kind = c.Next(f), kinds[f]&fives
		case kinds[f].Has(harmony.PlagalApproach):
			b.Target, b.Kind = c.Next(f), harmony.PlagalApproach
		case !isDominant(five.Chord.Pattern):
			continue
		}

		p := c.Prev(f)
		b.Aeolian = b.Kind == harmony.PlagalApproach && b.Target >= 0 && p >= 0 && p != f &&
			aeolian(c.Chords[p], five, c.Chords[b.Target])
		if p >= 0 && kinds[p].Has(harmony.SuspensionApproach) {
			b.Sus, p = p, c.Prev(p)
		}
		// A plagal IV has no two; the ♭VII7 has its IVm7 (Fm7 B♭7 C).
		backdoor := b.Target >= 0 && c.Chords[b.Target].Chord.Root == five.Chord.Root.Transpose(2)
		if p >= 0 && p != f && !b.Kind.Has(harmony.DiminishedApproach) &&
			(b.Kind != harmony.PlagalApproach || backdoor) && isTwoOf(c.Chords[p], five) {
			b.Two = p
		}
		if b.Target < 0 && b.Two < 0 {
			continue // a dominant alone, going nowhere
		}

		b.Announced = announce(c, b)
		out = append(out, b)
	}
	return out
}

// isTwoOf reports whether a change is the two of a V: before it, or
// before its suspension, which the two approaches as it would the V.
func isTwoOf(two, five Change) bool {
	if two.Silent {
		return false
	}
	v := harmony.Chord{Root: five.Chord.Root, Pattern: harmony.ChordDominantSeventh}
	return harmony.ApproachOf(two.Chord, v).Has(harmony.TwoApproach)
}

// isDominant reports whether a pattern holds the tritone of a dominant.
func isDominant(p harmony.ChordPattern) bool {
	return p.HasOffset(4) && p.HasOffset(10)
}
