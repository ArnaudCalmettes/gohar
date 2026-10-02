package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Pedal is a bass held under chords on different roots: the C under
// Cm7 Dm7/C E♭maj7/C in the introduction of Stolen Moments, the B♭
// under Fm9/B♭ B♭7 E♭maj9/B♭.
//
// # What the source says
//
// En Harmonie analyses a passage over a pedal « toujours en deux
// temps : on considère d'abord le type de pédale, puis les accords
// présents », the bass being then « une note étrangère » (tome 2,
// chapter 5 §1.6). It sorts pedals by the degree held, tonic or
// dominant, and by their reach, general or passing.
//
// # How gohar reads it
//
// The chords keep their own reading: each is figured on its root, and
// its function holds over the pedal (the B13/B♭ of I Fall In Love Too
// Easily is a chromatic dominant). The pedal is named on top of them.
//
// Its degree is read in the ground heard at its first chord. The book
// gives no measure of reach, so ours is set here: a pedal is general
// when it holds over a whole section of the form found by [Sections],
// passing otherwise.
//
// What the book says of the double pedal, a fifth held to frame a
// melody, is advice for an arrangement: a chart of chords alone has
// nothing to show of it.
type Pedal struct {
	From, To int // its first and last changes
	Bass     harmony.PitchClass
	Kind     PedalKind
	General  bool
}

// PedalKind tells which degree of the ground a pedal holds.
type PedalKind int

const (
	OtherPedal    PedalKind = iota // any degree but the tonic and the dominant
	TonicPedal                     // the 1st degree
	DominantPedal                  // the 5th degree
)

// pedalFifth is the dominant, a fifth above the tonic.
const pedalFifth = 7

// Pedals finds the pedals of a sequence: a bass held under two changes
// or more in a row, on two roots at least, so that one of them at least
// is inverted. Cm/E♭ Cm/E♭ is one chord written again; C Am/C is the
// shortest pedal. A silence, or the end of the chorus, stops one.
// `grounds` are those of [Grounds], `sections` those of [Sections]: nil
// sections leave every pedal passing.
func Pedals(c Changes, grounds [][]harmony.Tonality, sections []Section) []Pedal {
	var out []Pedal
	for i := 0; i < len(c.Chords); {
		j, roots := i+1, map[harmony.PitchClass]bool{c.Chords[i].Chord.Root: true}
		for j < len(c.Chords) && c.Next(j-1) == j && !c.Chords[j].Silent && c.Chords[j].Bass == c.Chords[i].Bass {
			roots[c.Chords[j].Chord.Root] = true
			j++
		}
		if !c.Chords[i].Silent && len(roots) >= 2 {
			p := Pedal{From: i, To: j - 1, Bass: c.Chords[i].Bass}
			if g := grounds[i]; g != nil {
				switch (int(p.Bass) - int(g[0].Tonic()) + 12) % 12 {
				case 0:
					p.Kind = TonicPedal
				case pedalFifth:
					p.Kind = DominantPedal
				}
			}
			p.General = holdsSection(c, p, sections)
			out = append(out, p)
		}
		i = j
	}
	return out
}

// holdsSection reports whether a pedal holds over a whole section, a
// pickup aside.
func holdsSection(c Changes, p Pedal, sections []Section) bool {
	from := c.Chords[p.From].Start
	to := c.Chords[p.To].Start + c.Chords[p.To].Length
	last := c.Chords[len(c.Chords)-1]
	for _, s := range sections {
		if s.Label == "-" {
			continue
		}
		end := last.Start + last.Length
		if next := s.From + s.Bars; next < len(c.Bars) {
			end = c.Bars[next]
		}
		if from <= c.Bars[s.From] && end <= to {
			return true
		}
	}
	return false
}
