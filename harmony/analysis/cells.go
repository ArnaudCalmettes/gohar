package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Cell is a formula of a few chords that standards are built from, and
// that a musician hears as one thing: the anatole, the III VI II V. En
// Harmonie presents them with the frequent chord sequences (tome 1,
// chapter 8 §3.3).
type Cell struct {
	From, To int // its first and last changes
	Kind     CellKind

	// Resolves tells whether the V goes on to its I, or to the III
	// standing for it: Dm7 after B♭ G7 Cm7 F7 in Anthropology. False
	// when the I promised is avoided.
	Resolves bool
}

// A CellKind is which formula a cell is.
type CellKind int

const (
	// Anatole: I VI II V, "known in France as the anatole" (the cell, not
	// the form) "and as rhythm changes in English-speaking countries":
	// B♭ Gm7 Cm7 F7 in I Got Rhythm. It cycles along the fifths, and
	// "comes in many forms, in major and in minor, through borrowings and
	// substitutions": Fm D♭maj7 Gm7♭5 C7 in F harmonic minor, Cm Am7♭5
	// Dm7♭5 G7 in Softly, As In A Morning Sunrise, and with secondary
	// dominants, Cmaj7 A7 Dm7 G7, Cmaj7 A7 D7 G7 (chapter 9, "Modifier un
	// enchaînement harmonique").
	Anatole CellKind = iota + 1

	// ThreeSixTwoFive: III VI II V, the book's III VI II V I, "a simple
	// variant of the anatole", the III standing for the I, to play the
	// anatole twice without playing the I again: Fmaj7 Dm7 Gm7 C7 Am7 Dm7
	// Gm7 C7 in Have You Met Miss Jones.
	ThreeSixTwoFive
)

var cellNames = [...]string{"", "anatole", "III-VI-II-V"}

func (k CellKind) String() string {
	return cellNames[k]
}

// Cells finds the cells of a sequence, from its roots and its chords'
// thirds, not from the tonality: four changes whose roots are, from the
// tonic the V points to, I (or III), VI, II and V; the VI on the major
// sixth or, in minor, the minor sixth (Fm D♭maj7), the III on either
// third. The first is a tonic chord for an anatole, a minor tonic
// written m7 included (Cm7 Am7♭5 Dm7♭5 G7), and a minor chord for a III
// VI II V; the VI and the II may be of any quality with a third,
// secondary dominants included; the last is a dominant. A change held
// longer counts once: B♭ | G7 | Cm7 | F7 and B♭ G7 | Cm7 F7 alike.
//
// A cell is heard by its shape, whether or not its V resolves: Fm7 B♭7
// E♭7 A♭7 in Anthropology is a III VI II V of D flat, and the D flat it
// promises is avoided, A♭7 going to Dm7. The cell says so (Resolves);
// the surprise itself belongs to the expectation.
func Cells(c Changes) []Cell {
	var out []Cell
	for i := 0; i+3 < len(c.Chords); i++ {
		ch := c.Chords[i : i+4]
		if ch[0].Silent || ch[1].Silent || ch[2].Silent || ch[3].Silent || !isDominant(ch[3].Chord.Pattern) {
			continue
		}
		tonic := ch[3].Chord.Root.Transpose(5)
		degree := func(k int) int { return (int(ch[k].Chord.Root) - int(tonic) + 12) % 12 }
		third := func(k int) bool { return ch[k].Chord.Pattern.HasOffset(3) || ch[k].Chord.Pattern.HasOffset(4) }
		if six := degree(1); six != 8 && six != 9 || degree(2) != 2 || !third(1) || !third(2) {
			continue
		}
		var kind CellKind
		switch first := degree(0); {
		case first == 0 && (tonicOf(ch[0]) != nil || ch[0].Chord.Pattern.Tetrad() == harmony.ChordMinorSeventh):
			kind = Anatole
		case (first == 3 || first == 4) && isMinor(ch[0].Chord.Pattern):
			kind = ThreeSixTwoFive
		default:
			continue
		}
		out = append(out, Cell{From: i, To: i + 3, Kind: kind, Resolves: resolves(c, i+3, tonic)})
		i += 3
	}
	return out
}

// resolves reports whether the V at change five goes on to the tonic,
// or to its III, the minor chord that stands for it.
func resolves(c Changes, five int, tonic harmony.PitchClass) bool {
	next := c.Next(five)
	if next < 0 || c.Chords[next].Silent {
		return false
	}
	ch := c.Chords[next].Chord
	third := (int(ch.Root) - int(tonic) + 12) % 12
	return ch.Root == tonic || (third == 3 || third == 4) && isMinor(ch.Pattern)
}
