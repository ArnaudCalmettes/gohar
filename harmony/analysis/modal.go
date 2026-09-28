package analysis

// A Plage is a stretch of changes where one chord is held with no
// cadence leading to it or away from it: a modal plage, which installs
// a colour rather than a tonic. Dm7 for sixteen bars in So What, the
// X7sus4 of kind of Maiden Voyage.
//
// A chart only says that the plage is modal, not which mode it is: a
// tetrad does not tell D dorian from D aeolian. The colour is in the
// melody, or in the analyst's knowledge of the tune, and a chart that
// writes only the tetrad is missing it. The analysis names the plage
// and leaves the mode to what the chord's provenance allows.
type Plage struct {
	From, To int // its first and last changes
}

// ModalBars is how long a chord must hold, with no cadence, to make a
// plage: a phrase's worth in the modal tunes (Maiden Voyage, Cantaloupe
// Island, Milestones), beyond the tonic held two bars in a tonal one.
const ModalBars = 4

// Modal finds the modal plages of a sequence: a chord, or the same
// chord written again, held for [ModalBars] bars or more, that no
// block leads to and that leads nowhere. A blues has none: its held
// I7 is its tonic.
func Modal(c Changes, blocks []Block) []Plage {
	if _, ok := Blues(c); ok {
		return nil
	}
	r := rolesOf(c, blocks)
	bar := barOf(c)
	var out []Plage
	for i := 0; i < len(c.Chords); {
		ch := c.Chords[i]
		to := i
		for to+1 < len(c.Chords) && !c.Chords[to+1].Silent && c.Chords[to+1].Chord == ch.Chord {
			to++
		}
		length := c.Chords[to].Start + c.Chords[to].Length - ch.Start
		// A turnaround back to the first chord, across the loop, is no
		// cadence to it.
		cadenced := r.target[i] >= 0 && blocks[r.target[i]].Five < i || r.member[to] >= 0
		if !ch.Silent && length >= Ticks(ModalBars)*bar && !cadenced {
			out = append(out, Plage{From: i, To: to})
		}
		i = to + 1
	}
	return out
}

// IsModal reports whether a sequence is a modal tune: its plages take
// up half of it or more. So What and Maiden Voyage are all plages, One
// Finger Snap twelve bars of twenty; a tonic held four bars in a tonal
// tune is not a modal tune.
func IsModal(c Changes, plages []Plage) bool {
	var in, all Ticks
	for _, ch := range c.Chords {
		all += ch.Length
	}
	for _, p := range plages {
		for _, ch := range c.Chords[p.From : p.To+1] {
			in += ch.Length
		}
	}
	return all > 0 && 2*in >= all
}
