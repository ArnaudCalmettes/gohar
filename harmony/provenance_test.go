package harmony_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

func chordOf(t *testing.T, root harmony.PitchClass, offsets ...harmony.Semitones) harmony.Chord {
	t.Helper()
	p, err := harmony.NewChordPattern(offsets...)
	if err != nil {
		t.Fatal(err)
	}
	return harmony.Chord{Root: root, Pattern: p}
}

func TestNamedScales(t *testing.T) {
	for s := harmony.NamedScale(0); s < harmony.NamedScaleCount; s++ {
		got, ok := harmony.NamedScaleOf(s.Pattern())
		if !ok || got != s {
			t.Errorf("scale %d: its pattern names %d, %v", s, got, ok)
		}
	}
	dorian, _ := harmony.ScaleMajor.Mode(2)
	if _, ok := harmony.NamedScaleOf(dorian); ok {
		t.Error("the dorian is taken for a named scale")
	}
}

// The borrowings of Tenderly, bars 1 to 16, as En Harmonie analyses
// them: each chord comes, among others, from the scale the book names.
func TestProvenancesOfTenderly(t *testing.T) {
	const (
		c, db, d, eb, e, f, gb, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	var (
		seventh    = []harmony.Semitones{0, 4, 7, 10}
		minor7     = []harmony.Semitones{0, 3, 7, 10}
		halfDim    = []harmony.Semitones{0, 3, 6, 10}
		diminished = []harmony.Semitones{0, 3, 6, 9}
		flatNine   = []harmony.Semitones{0, 4, 7, 10, 13}
	)
	for _, tc := range []struct {
		bar   string
		chord harmony.Chord
		want  harmony.Provenance
	}{
		{"2, Ab7", chordOf(t, ab, seventh...), harmony.Provenance{Scale: harmony.NamedMelodicMinor, Tonic: eb, Degree: 4}},
		{"3, Ebm7", chordOf(t, eb, minor7...), harmony.Provenance{Scale: harmony.NamedNaturalMinor, Tonic: eb, Degree: 1}},
		{"6, Db7", chordOf(t, db, seventh...), harmony.Provenance{Scale: harmony.NamedMelodicMinor, Tonic: ab, Degree: 4}},
		{"8, Gm7b5", chordOf(t, g, halfDim...), harmony.Provenance{Scale: harmony.NamedHarmonicMinor, Tonic: f, Degree: 2}},
		{"8, C7b9", chordOf(t, c, flatNine...), harmony.Provenance{Scale: harmony.NamedHarmonicMinor, Tonic: f, Degree: 5}},
		{"9, Fm7b5", chordOf(t, f, halfDim...), harmony.Provenance{Scale: harmony.NamedHarmonicMinor, Tonic: eb, Degree: 2}},
		{"10, Bb7", chordOf(t, bb, seventh...), harmony.Provenance{Scale: harmony.NamedHarmonicMinor, Tonic: eb, Degree: 5}},
		{"12, Bdim7", chordOf(t, b, diminished...), harmony.Provenance{Scale: harmony.NamedHarmonicMinor, Tonic: c, Degree: 7}},
	} {
		if got := harmony.Provenances(tc.chord); !slices.Contains(got, tc.want) {
			t.Errorf("bar %s: %v, without %v", tc.bar, got, tc.want)
		}
	}
}

func TestProvenances(t *testing.T) {
	const (
		c, d, e, f, g, ab, a harmony.PitchClass = 0, 2, 4, 5, 7, 8, 9
	)
	var (
		major7   = []harmony.Semitones{0, 4, 7, 11}
		seventh  = []harmony.Semitones{0, 4, 7, 10}
		flatNine = []harmony.Semitones{0, 4, 7, 10, 13}
		alt      = []harmony.Semitones{0, 4, 6, 10, 13, 15, 20}
		cluster  = []harmony.Semitones{0, 1, 2, 3}
	)
	const (
		major  = harmony.NamedMajor
		minor  = harmony.NamedNaturalMinor
		harmMi = harmony.NamedHarmonicMinor
		melMi  = harmony.NamedMelodicMinor
		harmMa = harmony.NamedHarmonicMajor
	)
	from := func(s harmony.NamedScale, tonic harmony.PitchClass, deg harmony.Degree) harmony.Provenance {
		return harmony.Provenance{Scale: s, Tonic: tonic, Degree: deg}
	}
	for name, tc := range map[string]struct {
		chord harmony.Chord
		want  []harmony.Provenance
	}{
		"Cmaj7": {chordOf(t, c, major7...), []harmony.Provenance{
			from(major, c, 1), from(major, g, 4),
			from(minor, e, 6), from(minor, a, 3),
			from(harmMi, e, 6),
			from(harmMa, c, 1),
		}},
		"C7": {chordOf(t, c, seventh...), []harmony.Provenance{
			from(major, f, 5),
			from(minor, d, 7),
			from(harmMi, f, 5),
			from(melMi, f, 5), from(melMi, g, 4),
			from(harmMa, f, 5), from(harmMa, ab, 3),
		}},
		// The extensions narrow the answer.
		"C7b9": {chordOf(t, c, flatNine...), []harmony.Provenance{
			from(harmMi, f, 5),
			from(harmMa, f, 5), from(harmMa, ab, 3),
		}},
		// 7alt, locrien ♭4 on G: the whole of A flat melodic minor.
		"G7alt": {chordOf(t, g, alt...), []harmony.Provenance{
			from(melMi, ab, 7),
		}},
		"a cluster no scale holds": {chordOf(t, c, cluster...), nil},
	} {
		if got := harmony.Provenances(tc.chord); !slices.Equal(got, tc.want) {
			t.Errorf("%s: %v, want %v", name, got, tc.want)
		}
	}
}

// Transposing the chord transposes the tonics, and nothing else.
func TestProvenancesTranspose(t *testing.T) {
	base := harmony.Provenances(chordOf(t, 0, 0, 3, 7, 10))
	for n := harmony.Semitones(1); n < 12; n++ {
		got := harmony.Provenances(chordOf(t, harmony.PitchClass(n), 0, 3, 7, 10))
		want := make([]harmony.Provenance, len(base))
		for i, p := range base {
			p.Tonic = p.Tonic.Transpose(n)
			want[i] = p
		}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("+%d: %v, want %v", n, got, want)
		}
	}
}
