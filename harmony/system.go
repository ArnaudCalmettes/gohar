package harmony

// A System is one of the five mother scales the modes are drawn from.
//
// A mode is designated without any name by its mother scale and its
// degree: F lydian is the fourth degree of C natural major. That
// designation is structural, which is why it lives here. Attaching a
// name to it is the naming package's business; knowing which mother
// scale and which degree a pattern comes from is an ordinary operation
// on scales and has to be available without one.
//
// The five mother scales read as pairs of tetrachords, and three of them
// are named after theirs: the harmonic minor is a minor tetrachord under
// a harmonic one, the harmonic major a major under a harmonic, and the
// double harmonic two harmonics. The natural is major over major, the
// melodic minor over major.
type System uint8

const (
	NaturalMajor System = iota
	MelodicMinor
	HarmonicMinor
	HarmonicMajor
	DoubleHarmonicMajor
)

// SystemCount is the number of mother scales.
const SystemCount = 5

var systemPatterns = [SystemCount]ScalePattern{
	NaturalMajor:        ScaleMajor,
	MelodicMinor:        ScaleMelodicMinor,
	HarmonicMinor:       ScaleHarmonicMinor,
	HarmonicMajor:       ScaleHarmonicMajor,
	DoubleHarmonicMajor: ScaleDoubleHarmonicMajor,
}

// Pattern returns the mother scale of s, read from its own tonic.
//
// Returns the empty pattern for a value outside the five, rather than
// panicking: a System that came from outside the package deserves an
// answer that callers can test.
func (s System) Pattern() ScalePattern {
	if s >= SystemCount {
		return 0
	}
	return systemPatterns[s]
}

// Mode returns the pattern of degree d of s, and whether d exists.
func (s System) Mode(d Degree) (ScalePattern, bool) {
	return s.Pattern().Mode(d)
}

// ModeOf finds the mother scale and degree a pattern comes from.
//
// Thirty five comparisons of sixteen bit integers, at most: two scales
// compare in a single instruction, so there is nothing to index and no
// table to build. The thirty five modes of the five systems are all
// distinct, so a match is unique.
//
// Reports false for a pattern that is not one of them. Most seven note
// patterns are not, and none with another number of notes is.
func ModeOf(p ScalePattern) (System, Degree, bool) {
	for s := System(0); s < SystemCount; s++ {
		for d := Degree(1); d <= 7; d++ {
			if m, ok := s.Mode(d); ok && m == p {
				return s, d, true
			}
		}
	}
	return 0, 0, false
}
