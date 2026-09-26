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
