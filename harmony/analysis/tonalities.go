package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// The analysis reads chords in tonalities, and often in several at
// once on one tonic: F harmonic minor; D flat major or D flat melodic
// minor, when nothing decides between them. A slice of tonalities
// sharing a tonic says it, each one a named scale: the scale is what
// the player improvises in, which is why En Harmonie writes "borrowed
// from F harmonic minor", not "from F minor".

// A Mode is major, minor, or both when tonalities of each remain.
type Mode uint8

const (
	Major Mode = 1 << iota
	Minor
)

// ModesOf tells whether tonalities have a major or a minor third.
func ModesOf(ts []harmony.Tonality) Mode {
	var m Mode
	for _, t := range ts {
		if t.Pattern().Contains(4) {
			m |= Major
		} else {
			m |= Minor
		}
	}
	return m
}

// holds reports whether one of the tonalities holds every class of
// `set`.
func holds(ts []harmony.Tonality, set harmony.PitchSet) bool {
	for _, t := range ts {
		if set.IsSubsetOf(t.Set()) {
			return true
		}
	}
	return false
}

// tonalitiesOf puts named scales on a tonic. They are heptatonic, so
// none fails.
func tonalitiesOf(tonic harmony.PitchClass, scales ...harmony.NamedScale) []harmony.Tonality {
	out := make([]harmony.Tonality, len(scales))
	for i, s := range scales {
		out[i], _ = harmony.NewTonality(tonic, s.Pattern())
	}
	return out
}

// MajorTonalities and MinorTonalities are how a tune in major or in
// minor is read: a minor tune in all three minors, the natural one for
// its key signature, the harmonic and the melodic for its cadences.
func MajorTonalities(tonic harmony.PitchClass) []harmony.Tonality {
	return tonalitiesOf(tonic, harmony.NamedMajor)
}

func MinorTonalities(tonic harmony.PitchClass) []harmony.Tonality {
	return tonalitiesOf(tonic, harmony.NamedNaturalMinor, harmony.NamedHarmonicMinor, harmony.NamedMelodicMinor)
}

// announcing are the scales a cadence announces a tonality in: the
// four of tonal practice. The harmonic major holds the chords of some
// cadences too (a half diminished two before a major tonic) but it is
// a scale one borrows from, not a tonality one announces; it stays in
// the provenance of chords.
var announcing = []harmony.NamedScale{
	harmony.NamedMajor,
	harmony.NamedNaturalMinor,
	harmony.NamedHarmonicMinor,
	harmony.NamedMelodicMinor,
}

// announce decides the tonalities a block announces.
//
// The tonic is the target, or a fifth under the V when the block does
// not resolve. The scales are those that hold the whole preparation:
// the two, the suspension, and the V when it is a dominant or a
// diminished chord. A chromatic dominant is not a degree of its
// tonality, nor is its two, the chromatic subdominant: they hold
// nothing, and the target alone decides. So a minor seventh two keeps
// the major and the melodic minor, a half diminished one the harmonic
// minor, a V seven flat thirteen the minors.
//
// Then the target decides between what is left, by its third, and
// only then: a preparation that holds no scale of the target's mode
// keeps its own. Gm7 C7 onto Fm7 is F melodic minor; Gm7♭5 C7 onto F
// major stays F harmonic minor, En Harmonie's mixed cadence.
//
// When the extensions leave no scale at all, the tetrads alone are
// tried: a colour the tonal scales do not hold, a sharp eleven, should
// not hide the tonality.
func announce(c Changes, b Block) []harmony.Tonality {
	tonic := c.Chords[b.Five].Chord.Root.Transpose(5)
	if b.Target >= 0 {
		tonic = c.Chords[b.Target].Chord.Root
	}
	var members []harmony.Chord
	five := c.Chords[b.Five].Chord
	if !b.Kind.Has(harmony.ChromaticApproach) {
		members = append(members, five)
		if b.Sus >= 0 {
			members = append(members, c.Chords[b.Sus].Chord)
		}
		if two := c.Chords[max(b.Two, 0)].Chord; b.Two >= 0 && two.Root == five.Root.Transpose(7) {
			members = append(members, two)
		}
	}

	var out []harmony.Tonality
	for _, tetrads := range []bool{false, true} {
		for _, t := range tonalitiesOf(tonic, announcing...) {
			fits := true
			for _, ch := range members {
				if tetrads {
					ch.Pattern = ch.Pattern.Tetrad()
				}
				fits = fits && ch.Set().IsSubsetOf(t.Set())
			}
			if fits {
				out = append(out, t)
			}
		}
		if len(out) > 0 {
			break
		}
	}

	if b.Target >= 0 && len(out) > 1 {
		third := c.Chords[b.Target].Chord.Pattern
		var want Mode
		switch {
		case third.HasOffset(4):
			want = Major
		case third.HasOffset(3):
			want = Minor
		}
		var kept []harmony.Tonality
		for _, t := range out {
			if ModesOf([]harmony.Tonality{t}) == want {
				kept = append(kept, t)
			}
		}
		if len(kept) > 0 {
			out = kept
		}
	}
	return out
}
