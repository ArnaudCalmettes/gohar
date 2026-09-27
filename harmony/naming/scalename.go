package naming

import "github.com/ArnaudCalmettes/gohar/harmony"

// ScaleName returns the name of a scale on its tonic: ré majeur,
// fa♯ mineur harmonique; D major, F♯ harmonic minor. The tonic is
// spelled in the scale's own key.
//
// False for a pattern no name designates, and for a locale that has no
// word for it.
func (n *Namer) ScaleName(tonic harmony.PitchClass, p harmony.ScalePattern) (string, bool) {
	s, ok := harmony.NamedScaleOf(p)
	if !ok || n.locale.Scales[s] == "" {
		return "", false
	}
	t, err := harmony.NewTonality(tonic, p)
	if err != nil {
		return "", false
	}
	return n.WithTonality(t).Name(tonic) + " " + n.locale.Scales[s], true
}
