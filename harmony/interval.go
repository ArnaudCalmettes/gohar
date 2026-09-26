package harmony

// Degrees is a distance measured in scale steps, as [Semitones] is a
// distance measured in semitones.
//
// It is to [Degree] what Semitones is to [Pitch]: a span, not a
// position. Degree 3 is the third note of a scale; Degrees(2) is the
// gap that separates it from the first. Mixing them silently produces
// results that are wrong by one, which is the hardest kind of error to
// see in a table of intervals.
type Degrees int8

// An Interval is the gap between two notes, held as a degree span and a
// semitone count.
//
// # Why both numbers
//
// Neither alone identifies an interval. Six semitones is an augmented
// fourth or a diminished fifth, and the semitone count cannot tell
// them apart. Three scale steps is a fourth of some quality, and the
// step count cannot say which. Together they are exact: the augmented
// fourth is three steps and six semitones, the diminished fifth is
// four steps and six semitones.
//
// # This is not spelling
//
// An earlier draft of this package left Interval out, on the grounds
// that choosing between augmented fourth and diminished fifth was a
// naming decision and belonged in the naming package. That was wrong.
// The choice is carried by the degree span, which is structural: it
// says which rung of the stack the note occupies, and a chord where
// the fourth is raised is a different chord from one where the fifth
// is lowered, whatever either is called.
//
// What does belong to naming is the word. The core computes
// [Interval.Alteration] and hands over a number; naming turns it into
// augmented, diminished, minor or major, with the perfect and
// imperfect split its own business.
type Interval struct {
	Degrees   Degrees
	Semitones Semitones
}

// The intervals in common use, from the unison to the major
// fourteenth.
//
// There are gaps past the octave. A minor twelfth or a minor
// fourteenth is never useful in practice, the version below the octave
// being what one reaches for, so they are absent rather than declared
// and unused.
var (
	IntUnison            = Interval{0, 0}
	IntMinorSecond       = Interval{1, 1}
	IntMajorSecond       = Interval{1, 2}
	IntAugmentedSecond   = Interval{1, 3}
	IntMinorThird        = Interval{2, 3}
	IntMajorThird        = Interval{2, 4}
	IntDiminishedFourth  = Interval{3, 4}
	IntPerfectFourth     = Interval{3, 5}
	IntAugmentedFourth   = Interval{3, 6}
	IntDiminishedFifth   = Interval{4, 6}
	IntPerfectFifth      = Interval{4, 7}
	IntAugmentedFifth    = Interval{4, 8}
	IntMinorSixth        = Interval{5, 8}
	IntMajorSixth        = Interval{5, 9}
	IntAugmentedSixth    = Interval{5, 10}
	IntDiminishedSeventh = Interval{6, 9}
	IntMinorSeventh      = Interval{6, 10}
	IntMajorSeventh      = Interval{6, 11}
	IntOctave            = Interval{7, 12}
	IntMinorNinth        = Interval{8, 13}
	IntMajorNinth        = Interval{8, 14}
	IntAugmentedNinth    = Interval{8, 15}
	IntMinorTenth        = Interval{9, 15}
	IntMajorTenth        = Interval{9, 16}
	IntPerfectEleventh   = Interval{10, 17}
	IntAugmentedEleventh = Interval{10, 18}
	IntMinorThirteenth   = Interval{12, 20}
	IntMajorThirteenth   = Interval{12, 21}
	IntMajorFourteenth   = Interval{13, 23}
)

// IsPerfect reports whether the degree span of `i` takes perfect
// qualities rather than major and minor ones.
//
// The unison, the fourth and the fifth, and their compounds. The split
// is not a convention: those intervals emerge first from the harmonic
// series, which is why they behave differently, and why a natural
// fourth is perfect and never major.
func (i Interval) IsPerfect() bool {
	d := i.step()
	return d == 0 || d == 3 || d == 4
}

// Natural returns the semitone count that `i` would have at its natural
// quality: perfect for a perfect degree span, major for the others.
//
// This is the reference every quality is measured against. It is the
// major scale, extended past the octave.
func (i Interval) Natural() Semitones {
	d := i.step()
	octaves := (int(i.Degrees) - d) / 7
	return majorOffsets[d] + Semitones(12*octaves)
}

// Alteration returns how far `i` sits from its natural quality, in
// semitones. Zero for a major or perfect interval, minus one for a
// minor or diminished one, and so on.
//
// A number, not a word. Turning minus one into minor on an imperfect
// degree and into diminished on a perfect one is naming's job, and it
// needs the tables the core does not have.
func (i Interval) Alteration() Semitones {
	return i.Semitones - i.Natural()
}

// Fold reduces `i` below the octave, keeping its quality.
//
// A major ninth folds to a major second, a perfect eleventh to a
// perfect fourth. Both numbers move together, which is what keeps the
// quality intact: folding the semitones alone would turn a major ninth
// into something two steps off.
func (i Interval) Fold() Interval {
	d := i.step()
	return Interval{Degrees(d), majorOffsets[d] + i.Alteration()}
}

// Up returns `i`, and Down returns the same interval taken downward.
func (i Interval) Down() Interval {
	return Interval{-i.Degrees, -i.Semitones}
}

// AddOctave returns `i` raised by one octave, a second becoming a ninth.
func (i Interval) AddOctave() Interval {
	return Interval{i.Degrees + 7, i.Semitones + 12}
}

// IsEnharmonic reports whether `i` and `other` span the same distance in
// semitones while differing in degree.
//
// An augmented fourth and a diminished fifth are enharmonic. So are a
// raised ninth and a minor tenth, which is the pair that makes an
// altered chord awkward to describe: one hears the ninth, and once the
// chord is identified one calls it a tenth.
func (i Interval) IsEnharmonic(other Interval) bool {
	return i.Semitones == other.Semitones && i.Degrees != other.Degrees
}

// majorOffsets holds the semitone count of each degree of the major
// scale, which is the reference every quality is measured against.
var majorOffsets = [7]Semitones{0, 2, 4, 5, 7, 9, 11}

// step reduces the degree span of `i` to a rung between 0 and 6, folding
// downward spans the same way as upward ones so that the quality of a
// descending interval reads like that of its ascending twin.
func (i Interval) step() int {
	d := int(i.Degrees) % 7
	if d < 0 {
		d += 7
	}
	return d
}
