package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

func TestScaleName(t *testing.T) {
	for _, c := range []struct {
		locale naming.Locale
		tonic  harmony.PitchClass
		p      harmony.ScalePattern
		want   string
	}{
		{naming.French, 2, harmony.ScaleMajor, "ré majeur"},
		{naming.French, 9, harmony.ScaleNaturalMinor, "la mineur"},
		{naming.French, 6, harmony.ScaleHarmonicMinor, "fa♯ mineur harmonique"},
		{naming.English, 2, harmony.ScaleMajor, "D major"},
		{naming.English, 6, harmony.ScaleHarmonicMinor, "F♯ harmonic minor"},
	} {
		n, err := naming.NewNamer(c.locale)
		if err != nil {
			t.Fatal(err)
		}
		if got, ok := n.ScaleName(c.tonic, c.p); !ok || got != c.want {
			t.Errorf("%q, %v; want %q", got, ok, c.want)
		}
	}
}

// A mode is not a scale: the dorian has no scale name.
func TestModesAreNotScales(t *testing.T) {
	n, _ := naming.NewNamer(naming.French)
	dorian, _ := harmony.ScaleMajor.Mode(2)
	if got, ok := n.ScaleName(2, dorian); ok {
		t.Errorf("dorian named %q as a scale", got)
	}
}

func TestIntervalName(t *testing.T) {
	for _, c := range []struct {
		locale naming.Locale
		i      harmony.Interval
		want   string
	}{
		{naming.French, harmony.IntUnison, "unisson"},
		{naming.French, harmony.IntMajorThird, "tierce majeure"},
		{naming.French, harmony.IntPerfectFourth, "quarte juste"},
		{naming.French, harmony.IntAugmentedFourth, "quarte augmentée"},
		{naming.French, harmony.IntDiminishedFifth, "quinte diminuée"},
		{naming.French, harmony.IntMinorSeventh, "septième mineure"},
		{naming.English, harmony.IntMajorThird, "major third"},
		{naming.English, harmony.IntPerfectFifth, "perfect fifth"},
	} {
		if got, ok := c.locale.IntervalName(c.i); !ok || got != c.want {
			t.Errorf("%v: %q, %v; want %q", c.i, got, ok, c.want)
		}
	}
}

// English names intervals and still has no spoken register.
func TestEnglishHasNoSpokenRegister(t *testing.T) {
	for _, m := range naming.Modes() {
		if s, ok := naming.English.SpokenModeName(m, naming.Signs); ok {
			t.Fatalf("%v spoken in English as %q", m, s)
		}
	}
}
