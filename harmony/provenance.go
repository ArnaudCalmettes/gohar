package harmony

// A NamedScale is a scale practice designates by a name on its tonic:
// D major, F sharp harmonic minor.
//
// # A scale is not a mode
//
// The major scale is not the ionian, and the natural minor, reference
// of the minor keys, is not the aeolian, though it is also the mode of
// the sixth degree of the major. A mode is designated by its mother
// scale and its degree, see [System]; a named scale by this type. The
// two share their notes and not their role, and the dex tracks them as
// distinct entries.
//
// Which scales are named is a fact about the practice, decided here
// once. The double harmonic is a mother scale but has no name of its
// own in tonal practice, and is left out.
type NamedScale uint8

const (
	NamedMajor NamedScale = iota
	NamedNaturalMinor
	NamedHarmonicMinor
	NamedMelodicMinor
	NamedHarmonicMajor
)

// NamedScaleCount is the number of named scales.
const NamedScaleCount = 5

var namedScalePatterns = [NamedScaleCount]ScalePattern{
	NamedMajor:         ScaleMajor,
	NamedNaturalMinor:  ScaleNaturalMinor,
	NamedHarmonicMinor: ScaleHarmonicMinor,
	NamedMelodicMinor:  ScaleMelodicMinor,
	NamedHarmonicMajor: ScaleHarmonicMajor,
}

// Pattern returns the pattern of `s` from its tonic, or the empty
// pattern for a value outside the five.
func (s NamedScale) Pattern() ScalePattern {
	if s >= NamedScaleCount {
		return 0
	}
	return namedScalePatterns[s]
}

// NamedScaleOf finds the named scale a pattern is, if any.
func NamedScaleOf(p ScalePattern) (NamedScale, bool) {
	for s, named := range namedScalePatterns {
		if named == p {
			return NamedScale(s), true
		}
	}
	return 0, false
}

// A Provenance is a named scale a chord can come from, and the degree
// its root holds there.
//
// It answers the player's first question on a chord of a chart: what
// to play over it. The answer is [Provenance.Mode], the scale read from
// the chord's root, and the provenance says where that mode comes from:
// D flat seven taken from A flat melodic minor is played in its fourth
// mode.
type Provenance struct {
	Scale  NamedScale
	Tonic  PitchClass
	Degree Degree // of the chord's root in the scale
}

// Tonality returns the scale the chord comes from, on its tonic.
func (p Provenance) Tonality() Tonality {
	t, _ := NewTonality(p.Tonic, p.Scale.Pattern())
	return t
}

// Mode returns the scale read from the chord's root.
func (p Provenance) Mode() ScalePattern {
	m, _ := p.Scale.Pattern().Mode(p.Degree)
	return m
}

// Provenances returns every named scale, on every tonic, that holds all
// the notes of `c`.
//
// # Every tonic, not the key's
//
// A borrowed chord comes from a scale on another tonic as often as from
// the key's own: in a tune in E flat, the D flat seven of Tenderly is
// borrowed from A flat melodic minor. So nothing here depends on a key.
// Choosing among the answers, the key of the tune first, then the one a
// cadence announces, is the analysis's business, which knows the
// context; this function only knows the chord.
//
// # The notes decide
//
// Every note of the pattern counts, extensions included: C seven sounds
// in seven scales, C seven flat nine in three. A caller who wants the
// widest answer passes the tetrad alone. The bass is not part of
// [Chord] and plays no part.
//
// A major scale and the natural minor on its sixth degree hold the same
// notes, and both are returned: they are two scales, not two names for
// one.
//
// The order is fixed: by named scale in the order of the constants,
// then by tonic counted up from the chord's root. It carries no
// preference, only determinism, and transposing the chord transposes
// the answer and nothing else.
func Provenances(c Chord) []Provenance {
	notes := c.Set()
	var out []Provenance
	for s := NamedScale(0); s < NamedScaleCount; s++ {
		for k := Semitones(0); k < 12; k++ {
			tonic := c.Root.Transpose(k)
			if !notes.IsSubsetOf(s.Pattern().At(tonic)) {
				continue
			}
			p := Provenance{Scale: s, Tonic: tonic}
			p.Degree, _ = p.Tonality().DegreeOf(c.Root)
			out = append(out, p)
		}
	}
	return out
}
