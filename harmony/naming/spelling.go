package naming

import (
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Letter is one of the seven letters a note can take.
//
// The order matters: spelling by degree walks the letters in sequence,
// wrapping after the seventh, so that a heptatonic scale gets each
// letter exactly once.
type Letter uint8

const (
	LetterC Letter = iota
	LetterD
	LetterE
	LetterF
	LetterG
	LetterA
	LetterB
	LetterCount = 7
)

// Natural returns the pitch class the letter denotes with no
// accidental: C is 0, D is 2, E is 4, F is 5, G is 7, A is 9, B is 11.
func (l Letter) Natural() harmony.PitchClass {
	return letterNaturals[int(l)%LetterCount]
}

// Next returns the letter after `l`, wrapping from B back to C.
func (l Letter) Next() Letter {
	return Letter((int(l) + 1) % LetterCount)
}

// An Accidental is the sign written on a letter, counted in semitones.
//
// # Not the same thing as a Quality
//
// The two types hold similar numbers and mean different things. A
// Quality measures a degree against the major scale: the third degree
// of any major scale is a natural third, whatever letter it lands on.
// An Accidental measures a note against its own letter: the third
// degree of D major is F sharp, a natural third carrying a sharp.
//
// Keeping them apart is what stops the spelling algorithm from
// applying a degree quality as if it were a written sign.
type Accidental int8

const (
	DoubleFlatSign  Accidental = -2
	FlatSign        Accidental = -1
	NaturalSign     Accidental = 0
	SharpSign       Accidental = 1
	DoubleSharpSign Accidental = 2
)

// A SpelledNote is a letter with an accidental: the written, sayable
// form of a pitch class.
//
// This is what the core deliberately does not carry. C sharp and D flat
// are one PitchClass and two SpelledNotes, and choosing between them
// needs the tonal context the core has no business guessing.
//
// The name follows the English sense of spelling, which is first an
// oral act: one spells a note aloud without any staff being involved.
// That is the use this library is built for.
type SpelledNote struct {
	Letter     Letter
	Accidental Accidental
}

// Class returns the pitch class `n` sounds.
//
// Many to one: E sharp and F return the same class, which is the whole
// reason the core works in classes and this package in spellings.
func (n SpelledNote) Class() harmony.PitchClass {
	return n.Letter.Natural().Transpose(harmony.Semitones(n.Accidental))
}

// Enharmonic reports whether `n` and `other` sound the same class while
// being written differently.
func (n SpelledNote) Enharmonic(other SpelledNote) bool {
	return n != other && n.Class() == other.Class()
}

// spellDegrees assigns a letter and an accidental to each degree of a
// tonality, starting from the letter given for its tonic.
//
// # The algorithm
//
// A heptatonic scale takes each of the seven letters exactly once, in
// order. So the letter of degree d follows from the tonic's letter and
// nothing else: walk the sequence, wrapping after the seventh. The
// accidental is then forced, being the distance from the letter's
// natural class to the class the scale actually holds.
//
// Nothing is chosen. Given a tonic spelling and a heptatonic pattern,
// every other spelling in the scale is determined. That is why the
// nomenclature is rigorous rather than conventional, and why a lydian
// sharp 2 has a sharp second rather than a flat third: the third letter
// is already taken.
//
// # When the accidental runs out of range
//
// Some patterns push a degree more than two semitones from its letter,
// which no accidental can write. Rather than fail, the spelling falls
// back to the nearest letter that can carry the note. The result is no
// longer one letter per degree, and it is not diatonic in the strict
// sense, but it is readable and the caller gets a note rather than an
// error. Correctness in the theoretical sense is not worth a failed
// display in the middle of a game.
func spellDegrees(t harmony.Tonality, tonic SpelledNote) []SpelledNote {
	if t.IsZero() {
		return nil
	}

	out := make([]SpelledNote, 0, 7)
	letter := tonic.Letter
	for _, offset := range t.Pattern().Offsets() {
		class := t.Tonic().Transpose(offset)
		if note, ok := spellOn(letter, class); ok {
			out = append(out, note)
		} else {
			out = append(out, defaultSpelling(class))
		}
		letter = letter.Next()
	}
	return out
}

// defaultSpelling gives the spelling a class takes with no tonality to
// place it in.
//
// The seven classes of the C natural scale take their own letter with
// no accidental. The five others take the letter below with a sharp,
// following the convention every MIDI tool already shows the player.
//
// This is a spelling convention of last resort, not a tonality. It has
// no tonic, it names no key, and nothing built on it may claim that the
// player is in C major.
func defaultSpelling(c harmony.PitchClass) SpelledNote {
	for l := Letter(0); l < LetterCount; l++ {
		if l.Natural() == c {
			return SpelledNote{Letter: l, Accidental: NaturalSign}
		}
	}
	for l := Letter(0); l < LetterCount; l++ {
		if l.Natural().Transpose(1) == c {
			return SpelledNote{Letter: l, Accidental: SharpSign}
		}
	}
	return SpelledNote{}
}

// letterNaturals holds the class each letter denotes with no
// accidental, in the order the letters are declared.
var letterNaturals = [LetterCount]harmony.PitchClass{0, 2, 4, 5, 7, 9, 11}

// spellOn writes a class on a given letter, and reports whether an
// accidental can reach it at all.
//
// The distance is folded into the range that an accidental can express.
// Beyond a double sharp or a double flat there is nothing to write, and
// the caller falls back rather than inventing a triple sign.
func spellOn(l Letter, c harmony.PitchClass) (SpelledNote, bool) {
	alt := int(l.Natural().Up(c))
	if alt > 6 {
		alt -= 12
	}
	if alt < int(DoubleFlatSign) || alt > int(DoubleSharpSign) {
		return SpelledNote{}, false
	}
	return SpelledNote{Letter: l, Accidental: Accidental(alt)}, true
}

// nearestSpelling writes a class as an alteration of the nearest note
// already spelled in the scale.
//
// Used for a class the tonality does not contain, where no degree owns
// the letter. Ties go to sharpening the note below, matching the
// convention [defaultSpelling] uses and the direction a chromatic
// approach usually takes.
func nearestSpelling(scale []SpelledNote, c harmony.PitchClass) (SpelledNote, bool) {
	best, found := SpelledNote{}, false
	bestDistance := 99

	for _, note := range scale {
		for _, delta := range []int{1, -1, 2, -2} {
			if note.Class().Transpose(harmony.Semitones(delta)) != c {
				continue
			}
			alt := int(note.Accidental) + delta
			if alt < int(DoubleFlatSign) || alt > int(DoubleSharpSign) {
				continue
			}
			distance := delta
			if distance < 0 {
				distance = -distance
			}
			if !found || distance < bestDistance {
				best = SpelledNote{Letter: note.Letter, Accidental: Accidental(alt)}
				bestDistance, found = distance, true
			}
		}
	}
	return best, found
}

// TonicSpelling returns the spelling of the tonic of `t` that writes
// its scale with the fewest accidentals: the rule of least alteration.
// F sharp minor, not G flat minor, which would need a double flat; C
// sharp minor, not D flat minor; D flat major, not C sharp major.
//
// The count is taken on the scale of the key signature, every degree
// spelled from the tonic (see spellDegrees), a double accidental
// counting twice. The harmonic and melodic minors take the signature of
// the natural minor, and the harmonic major that of the major, as they
// are written: G sharp harmonic minor and A flat harmonic minor both
// write six accidentals, the F double sharp of the first counting
// twice, but the signature of G sharp minor has five sharps and that of
// A flat minor seven flats. Any other scale of the catalogue, a mode of
// the major or of the melodic minor, is counted on itself.
//
// Two tonics tie, at six accidentals each way: F sharp or G flat major,
// D sharp or E flat minor. The tie goes to the flats, as jazz writes
// them. The rule and its tie are ours, from common sense and common use.
func TonicSpelling(t harmony.Tonality) SpelledNote {
	if t.IsZero() {
		return SpelledNote{}
	}
	best, bestCount := defaultSpelling(t.Tonic()), -1
	for l := Letter(0); l < LetterCount; l++ {
		tonic, ok := spellOn(l, t.Tonic())
		if !ok {
			continue
		}
		count := KeyAccidentals(t, tonic)
		switch {
		case bestCount < 0, count < bestCount,
			count == bestCount && tonic.Accidental < best.Accidental:
			best, bestCount = tonic, count
		}
	}
	return best
}

// Accidentals counts the accidentals written in `notes`, a double sharp
// or a double flat counting twice.
func Accidentals(notes []SpelledNote) int {
	n := 0
	for _, note := range notes {
		a := int(note.Accidental)
		if a < 0 {
			a = -a
		}
		n += a
	}
	return n
}

// KeyAccidentals counts the accidentals of the key signature of `t`
// with its tonic written `tonic`: 5 for G sharp minor, 7 for A flat
// minor, 8 for G sharp major, a key no signature writes. A double
// counts twice.
func KeyAccidentals(t harmony.Tonality, tonic SpelledNote) int {
	return Accidentals(spellDegrees(signature(t), tonic))
}

// signature returns the tonality whose scale gives the key signature of
// `t`: the natural minor for the harmonic and melodic minors, the major
// for the harmonic major, `t` itself otherwise.
func signature(t harmony.Tonality) harmony.Tonality {
	p := t.Pattern()
	switch p {
	case harmony.ScaleHarmonicMinor, harmony.ScaleMelodicMinor:
		p = harmony.ScaleNaturalMinor
	case harmony.ScaleHarmonicMajor:
		p = harmony.ScaleMajor
	default:
		return t
	}
	s, err := harmony.NewTonality(t.Tonic(), p)
	if err != nil {
		return t
	}
	return s
}

// SpellAbove writes class `c` on the letter `steps` letters above the
// letter of `from`, below when `steps` is negative, and reports whether
// an accidental can reach it.
//
// This is how a chord is spelled in a tonality: its root takes the
// letter of its degree, counted from the tonic, and the notes of the
// chord the letters of their intervals, counted from the root. In D
// flat, the ♭VI is a B double flat, the sixth letter from D (steps 5),
// and the E7 that is the VI7 of G is spelled E G♯ B D, never E A♭ B D.
// A degree lowered or raised keeps the letter of the degree it alters:
// ♯IVdim7 in C is F♯dim7, ♭VI7 is A♭7.
//
// With `c` the class of `from` itself, it gives the enharmonic of a
// note one letter away: SpellAbove(D♭, -1, C♯) is C♯.
func SpellAbove(from SpelledNote, steps int, c harmony.PitchClass) (SpelledNote, bool) {
	l := Letter(((int(from.Letter)+steps)%LetterCount + LetterCount) % LetterCount)
	return spellOn(l, c)
}
