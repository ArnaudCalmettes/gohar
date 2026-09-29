package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Cell is a formula of a few chords that standards are built from, and
// that a musician hears as one thing: the anatole, the III VI II V. En
// Harmonie presents them with the frequent chord sequences (tome 1,
// chapter 8 §3.3).
type Cell struct {
	From, To int // its first and last changes
	Kind     CellKind

	// Resolves tells whether the V goes on to its I. False when the I
	// promised is avoided, a deceptive cadence ("cadence rompue", V – …
	// in En Harmonie), the III included: in Anthropology, F7 goes to
	// Dm7, the mediant, which stands for the I in the III VI II V that
	// follows but is no resolution of the V.
	Resolves bool

	// Substituted tells whether some X7 of the cell stands for its
	// tritone twin, the cell being read on the chords it replaces:
	// C E♭7 A♭7 D♭7 is the anatole C A7 D7 G7. "Any X7 can be substituted
	// by another X7, whatever its function" (chapter 9, p. 141).
	Substituted bool
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
// VI II V; the VI may be of any quality with a third, a secondary
// dominant included; the II is minor or a dominant, a major chord there
// being the arrival of a cadence (Em7 A7 Dmaj7 G7 is a II V I of D,
// not a III VI II V of C); the last is a dominant. A change held
// longer counts once: B♭ | G7 | Cm7 | F7 and B♭ G7 | Cm7 F7 alike.
//
// A cell is heard by its shape, whether or not its V resolves: Fm7 B♭7
// E♭7 A♭7 in Anthropology is a III VI II V of D flat, and the D flat it
// promises is avoided, A♭7 going to Dm7: a deceptive cadence, which the
// cell notes (Resolves).
//
// Its X7 may have been substituted by their tritone twins. En Harmonie
// analyses such a passage "from right to left", to "find the initial
// cadences" again (chapter 9, pp. 142-143): Cells reads the four
// changes as written first, then with one X7 restored, then two, then
// three, and keeps the first reading that makes a cell. The fifths
// become half tones: in Body And Soul, Dm7 G7 C7 B7 B♭7 is a III VI II
// V of B flat walking down chromatically, B7 for F7.
func Cells(c Changes) []Cell {
	var out []Cell
	for i := 0; i+3 < len(c.Chords); i++ {
		ch := c.Chords[i : i+4]
		if ch[0].Silent || ch[1].Silent || ch[2].Silent || ch[3].Silent {
			continue
		}
		for _, twins := range restorations {
			if cl, ok := cellOf(ch, twins); ok {
				cl.From, cl.To = i, i+3
				cl.Resolves = resolves(c, i+3, cl.tonic)
				out = append(out, cl.Cell)
				i += 3
				break
			}
		}
	}
	return out
}

// restorations are the X7 among the last three changes of a cell that
// may be read as their tritone twin, one bit each, the fewest first:
// the chords as written win over any restoration.
var restorations = [...]uint8{0, 0b001, 0b010, 0b100, 0b011, 0b101, 0b110, 0b111}

// A reading is a cell found, with the tonic its V points to.
type reading struct {
	Cell
	tonic harmony.PitchClass
}

// cellOf reads four changes as a cell, the X7 marked in twins read as
// their tritone twin. See [Cells] for the shape.
func cellOf(ch []Change, twins uint8) (cl reading, ok bool) {
	var roots [4]harmony.PitchClass
	for k, change := range ch {
		roots[k] = change.Chord.Root
		if k > 0 && twins&(1<<(k-1)) != 0 {
			if !isDominant(change.Chord.Pattern) {
				return cl, false // only an X7 has a twin
			}
			roots[k] = roots[k].Transpose(6)
		}
	}
	if !isDominant(ch[3].Chord.Pattern) {
		return cl, false
	}
	tonic := roots[3].Transpose(5)
	degree := func(k int) int { return (int(roots[k]) - int(tonic) + 12) % 12 }
	third := func(k int) bool { return ch[k].Chord.Pattern.HasOffset(3) || ch[k].Chord.Pattern.HasOffset(4) }
	if six := degree(1); six != 8 && six != 9 || degree(2) != 2 || !third(1) || !third(2) {
		return cl, false
	}
	if two := ch[2].Chord.Pattern; !isMinor(two) && !isDominant(two) {
		return cl, false
	}
	if twins&0b001 != 0 && degree(1) != 9 {
		return cl, false // restored, the VI is the V of the II
	}
	switch first := degree(0); {
	case first == 0 && (tonicOf(ch[0]) != nil || ch[0].Chord.Pattern.Tetrad() == harmony.ChordMinorSeventh):
		cl.Kind = Anatole
	case (first == 3 || first == 4) && isMinor(ch[0].Chord.Pattern):
		cl.Kind = ThreeSixTwoFive
	default:
		return cl, false
	}
	cl.tonic, cl.Substituted = tonic, twins != 0
	return cl, true
}

// resolves reports whether the V at change five goes on to the tonic.
func resolves(c Changes, five int, tonic harmony.PitchClass) bool {
	next := c.Next(five)
	return next >= 0 && !c.Chords[next].Silent && c.Chords[next].Chord.Root == tonic
}
