package naming

import "github.com/ArnaudCalmettes/gohar/harmony"

// namedScales are the scales a name designates on their own, apart from
// the modes they share their notes with: the major scale is not the
// ionian, and the natural minor, the reference of the minor keys, is
// not the aeolian, though it is also the mode of the sixth degree of
// the major.
//
// Which scales are named is a fact about the practice, decided here
// once, as for the tetrachords.
var namedScales = [5]harmony.ScalePattern{
	harmony.ScaleMajor,
	harmony.ScaleNaturalMinor,
	harmony.ScaleHarmonicMinor,
	harmony.ScaleMelodicMinor,
	harmony.ScaleHarmonicMajor,
}

// ScaleName returns the name of a scale on its tonic: ré majeur,
// fa♯ mineur harmonique; D major, F♯ harmonic minor. The tonic is
// spelled in the scale's own key.
//
// False for a pattern no name designates, and for a locale that has no
// word for it.
func (n *Namer) ScaleName(tonic harmony.PitchClass, p harmony.ScalePattern) (string, bool) {
	for i, named := range namedScales {
		if named != p || n.locale.Scales[i] == "" {
			continue
		}
		t, err := harmony.NewTonality(tonic, p)
		if err != nil {
			return "", false
		}
		return n.WithTonality(t).Name(tonic) + " " + n.locale.Scales[i], true
	}
	return "", false
}
