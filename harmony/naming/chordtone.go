package naming

import "github.com/ArnaudCalmettes/gohar/harmony"

// SpellOver spells a note heard over a chord, counted from its root:
// the letter of the interval the note makes with the root, as a lead
// sheet writes the degrees of a chord. Over F7, the minor seventh is an
// E♭ and never a D♯; over B♭7, the ninth is a C.
//
// The interval names follow En Harmonie: three semitones are a ♭10
// (or ♭3), never a ♯9; eight are a ♭13 (or ♭6), never a ♯5. Two depend
// on the chord: six semitones are a ♭5 when the chord has no fifth to
// raise (Cm7♭5, C7♭5) and a ♯11 otherwise; nine are the diminished
// seventh of a dim7, its ♭♭7, and a 13th otherwise.
//
// `root` is the root as the chord symbol spells it. A spelling out of
// reach of a double accidental falls back on the letter nearest by
// default; it does not happen from a root with one accidental at most.
func SpellOver(root SpelledNote, p harmony.ChordPattern, c harmony.PitchClass) SpelledNote {
	steps := stepsOver[root.Class().Up(c)]
	switch root.Class().Up(c) {
	case 6:
		if p.HasOffset(6) && !p.HasOffset(7) {
			steps = 4 // the ♭5
		}
	case 9:
		if p.HasOffset(3) && p.HasOffset(6) && p.HasOffset(9) && !p.HasOffset(10) {
			steps = 6 // the diminished seventh of a dim7
		}
	}
	if n, ok := SpellAbove(root, steps, c); ok {
		return n
	}
	return defaultSpelling(c)
}

// stepsOver is how many letters above the root each interval sits, in
// semitones: 0 the root, 1 the ♭9, 2 the 9, 3 the ♭10, 4 the third, 5
// the fourth, 6 the ♯11, 7 the fifth, 8 the ♭13, 9 the 13, 10 the minor
// seventh, 11 the major seventh.
var stepsOver = [12]int{0, 1, 1, 2, 2, 3, 3, 4, 5, 5, 6, 6}
