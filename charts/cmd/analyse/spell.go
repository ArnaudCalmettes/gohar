package main

import (
	"fmt"
	"strings"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// A namer spells the tonics the report names as the chords under them
// are spelled, so that the two agree (see [respell]). The zero namer
// spells them with the fewest accidentals.
type namer struct {
	roots  []naming.SpelledNote                      // the root of each change, unset for a silence
	tonics map[harmony.PitchClass]naming.SpelledNote // the tonics of the zones, then the roots
}

// unset marks a change left as the chart spells it, a silence or a chord
// without a degree. It cannot be the zero value, which is C natural: no
// note is written with three sharps.
var unset = naming.SpelledNote{Accidental: 3}

// filled returns `n` notes, all unset.
func filled(n int) []naming.SpelledNote {
	out := make([]naming.SpelledNote, n)
	for i := range out {
		out[i] = unset
	}
	return out
}

// A respelling is a chart's chords as [respell] writes them.
type respelling struct {
	roots, basses []naming.SpelledNote // unset where the chart's own spelling stays
	names         namer                // the tonics, spelled as the chords are
	cost          cost
	smells        []int // the changes whose root or bass smells, once each
	stained       int   // the changes that smell or tolerate an enharmony

	// tolerated holds the changes whose root or bass needs a double
	// accidental by degree, shown instead on the next letter, and strict
	// their roots and basses by degree: the B𝄫 of Cdim7/B𝄫, shown A.
	tolerated          []int
	strict, strictBass []naming.SpelledNote
}

// A cost ranks the spellings of the chords of a zone, its tonic given:
// the movements left unreadable (P0), then the names that leave the
// figuring (P1), then the accidentals written, and last, for a chord
// that prepares a resolution, a name off its degree in the zone, which
// only breaks the ties its target leaves. It is compared component by
// component, the first that differs deciding. The smells do not enter
// it: they weigh the tonics of the zones (P2, see [respell]).
type cost [costs]int

// The components of a cost, in their order of priority.
const (
	motionCost     = iota // movements left unreadable (P0)
	figuringCost          // names that leave the figuring (P1)
	accidentalCost        // accidentals written
	tieCost               // chords preparing a resolution, named off their degree
	costs
)

func (a cost) plus(b cost) cost {
	for i := range a {
		a[i] += b[i]
	}
	return a
}

func (a cost) less(b cost) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// respell writes the chords of the timeline as the analysis reads them,
// in place. The rules are ours, by order of priority.
//
//	P0. The movements read as plainly as they can (see [motion]): a
//	    bass moving by a fifth moves by a perfect fifth, by a third a
//	    third, by a step a second. A semitone may also be chromatic, on
//	    the same letter, and a chromatic note follows the direction of
//	    the line, sharp going up, flat going down: E7 E♭7 D7 D♭7 C,
//	    C C♯dim7 Dm7.
//	P1. A root takes the letter of its degree in the zone, the ground
//	    the degrees are counted on (see analysis.Grounds): a ♭II is
//	    named as a ♭II. A root on a degree of the zone is the frame the
//	    movements are read in, and keeps its letter whatever they say;
//	    P0 rules the others, the chromatic roots and the chords that
//	    prepare a resolution, which take their name from their target
//	    (see [approaches]). P1 then only breaks their ties.
//	P2. The tonic of each zone is spelled with the fewest accidentals in
//	    its key signature (naming.TonicSpelling) or as its enharmonic, if
//	    a signature can write it, whichever gives the grid the fewest
//	    smells, then the fewest accidentals written by degree: C♯ major
//	    shows its E♯ as F, but still writes seven sharps.
//	P3. A tonicisation is named as the chord it tonicises (see [namer.tonic]).
//
// They apply in two stages. For each choice of the tonics of the zones,
// the chords are spelled by P0, then P1, the best path through the
// changes (see [cost] and [spellWith]); P2 then picks, among those
// grids, the one that smells least.
//
// A bass takes the letter of its interval with the root, except a held
// bass, a pedal, spelled once by its degree in the zone and kept while
// it is held (see [pedals]).
//
// Usage then simplifies the names shown, never the figuring nor the
// movements (see [usage]): an E♯, F♭, B♯ or C♭ is shown on the next
// letter, B for the C♭ of a ♭II7 in B♭ minor.
//
// A smell is a root or a bass shown with a double accidental, or as an
// E♯, F♭, B♯ or C♭ that usage keeps. It counts once for each change,
// and the report gives the count (see [smells]).
func respell(tl ireal.Timeline, c analysis.Changes, sensed []analysis.Sensed, blocks []analysis.Block, ds []analysis.Degree, ps []analysis.Pedal) respelling {
	grounds := analysis.Grounds(c, sensed)
	zones, firsts := zonesOf(grounds)
	approach := approaches(c, blocks)
	// Every choice of tonics, one per zone, while they stay few enough
	// to weigh together; the plain tonics otherwise.
	combos := [][]naming.SpelledNote{nil}
	for _, z := range firsts {
		var next [][]naming.SpelledNote
		for _, combo := range combos {
			for _, t := range tonicCandidates(grounds[z][0]) {
				next = append(next, append(append([]naming.SpelledNote(nil), combo...), t))
			}
		}
		if len(next) > 64 {
			next = [][]naming.SpelledNote{plainTonics(grounds, firsts)}
			combos = next
			break
		}
		combos = next
	}
	// The tonics are compared by the grid they give, the simplest
	// winning (P2): the fewest changes that smell or tolerate an
	// enharmony, then the fewest accidentals; the movements and the
	// figuring only break the ties, each choice of tonics having a
	// figuring of its own.
	simpler := func(a, b respelling) bool {
		if a.stained != b.stained {
			return a.stained < b.stained
		}
		for _, k := range []int{accidentalCost, motionCost, figuringCost, tieCost} {
			if a.cost[k] != b.cost[k] {
				return a.cost[k] < b.cost[k]
			}
		}
		return false
	}
	var best respelling
	for k, combo := range combos {
		if s := spellWith(c, ds, ps, grounds, zones, approach, combo); k == 0 || simpler(s, best) {
			best = s
		}
	}

	best.names = namer{roots: best.roots, tonics: map[harmony.PitchClass]naming.SpelledNote{}}
	for _, i := range firsts {
		if t, ok := zoneTonic(best, c, grounds, zones, i); ok {
			best.names.tonics[t.Class()] = t
		}
	}
	for _, r := range best.roots {
		if _, ok := best.names.tonics[r.Class()]; !ok && r != unset {
			best.names.tonics[r.Class()] = r
		}
	}
	for i := range tl.Spans {
		if i >= len(best.roots) || tl.Spans[i].NoChord || best.roots[i] == unset {
			continue
		}
		tl.Spans[i].Chord.Root = name(best.roots[i])
		if tl.Spans[i].Chord.Bass != "" && best.basses[i] != unset {
			tl.Spans[i].Chord.Bass = name(best.basses[i])
		}
		if symmetric(c.Chords[i]) && tl.Spans[i].Chord.Bass != "" {
			tl.Spans[i].Chord.Root, tl.Spans[i].Chord.Bass = tl.Spans[i].Chord.Bass, ""
		}
	}
	return best
}

// zonesOf numbers the zones, the runs of changes on the same ground,
// and returns the zone of each change, -1 without a ground, and the
// first change of each zone.
func zonesOf(grounds [][]harmony.Tonality) ([]int, []int) {
	zones := make([]int, len(grounds))
	var firsts []int
	for i, g := range grounds {
		switch {
		case len(g) == 0:
			zones[i] = -1
		case i > 0 && len(grounds[i-1]) > 0 && grounds[i-1][0] == g[0]:
			zones[i] = zones[i-1]
		default:
			zones[i] = len(firsts)
			firsts = append(firsts, i)
		}
	}
	return zones, firsts
}

// tonicCandidates spells the tonic of a zone with the fewest
// accidentals, then as its enharmonic when a signature can write it, seven
// accidentals at most: G♯m and A♭m, F♯ and G♭, but F minor alone, E♯
// minor needing eight sharps.
func tonicCandidates(t harmony.Tonality) []naming.SpelledNote {
	plain := naming.TonicSpelling(t)
	out := []naming.SpelledNote{plain}
	for _, step := range []int{-1, 1} {
		if s, ok := naming.SpellAbove(plain, step, plain.Class()); ok && count(s) <= 1 && naming.KeyAccidentals(t, s) <= 7 {
			out = append(out, s)
		}
	}
	return out
}

func plainTonics(grounds [][]harmony.Tonality, firsts []int) []naming.SpelledNote {
	out := make([]naming.SpelledNote, len(firsts))
	for z, i := range firsts {
		out[z] = naming.TonicSpelling(grounds[i][0])
	}
	return out
}

// zoneTonic returns the tonic of the zone starting at change `i` as the
// spelling chose it, read back from a root on the tonic, or from the
// reference spelling of a root by its degree.
func zoneTonic(s respelling, c analysis.Changes, grounds [][]harmony.Tonality, zones []int, i int) (naming.SpelledNote, bool) {
	tonic := grounds[i][0].Tonic()
	for j := i; j < len(c.Chords) && zones[j] == zones[i]; j++ {
		if r := s.roots[j]; r != unset && r.Class() == tonic {
			return r, true
		}
	}
	return unset, false
}

// approaches marks the chords that prepare a resolution, the two, the
// suspension and the five of a block that resolves: their name comes
// from their target, through the movements (P0 of [respell]), never
// from their degree in the zone. In E♭ minor, the Bm7 E7 that resolves
// on Am7 is spelled from A, not as the C♭m7 of the sixth degree.
func approaches(c analysis.Changes, blocks []analysis.Block) []bool {
	out := make([]bool, len(c.Chords))
	for _, b := range blocks {
		if b.Target < 0 {
			continue
		}
		for _, i := range []int{b.Two, b.Sus, b.Five} {
			if i >= 0 {
				out[i] = true
			}
		}
	}
	return out
}

// spellWith spells every chord with the tonics of the zones given, by
// dynamic programming over the changes: each root among its possible
// spellings, each movement between two roots weighed by [motion], the
// best path kept.
func spellWith(c analysis.Changes, ds []analysis.Degree, ps []analysis.Pedal, grounds [][]harmony.Tonality, zones []int, approach []bool, tonics []naming.SpelledNote) respelling {
	n := len(c.Chords)
	held := pedals(c, ps, grounds, zones, tonics)
	type state struct {
		root, bass   naming.SpelledNote // as the movements read them
		shown, under naming.SpelledNote // as usage writes them
		cost         cost
		from         int
	}
	states := make([][]state, n)
	for i, ch := range c.Chords {
		if ch.Silent {
			continue
		}
		var ref naming.SpelledNote
		hasRef := false
		if z := zones[i]; z >= 0 && ds[i].Number > 0 {
			ref, hasRef = naming.SpellAbove(tonics[z], int(ds[i].Number)-1, ch.Chord.Root)
		}
		// A chord that prepares a resolution takes its name from its
		// target: its degree in the zone only breaks the ties.
		tie := hasRef && approach[i]
		hasRef = hasRef && !approach[i]
		var prev []state
		if i > 0 {
			prev = states[i-1]
		}
		// A root on a degree of the zone is the frame the others move in:
		// it takes the letter of its degree, whatever the movements. A
		// chromatic root is free, and so is a chord that prepares a
		// resolution, named from its target.
		candidates := spellings(ch.Chord.Root)
		if _, diatonic := degreeIn(grounds[i], ch.Chord.Root); hasRef && diatonic {
			candidates = []naming.SpelledNote{ref}
		} else {
			candidates = usual(candidates)
		}
		for _, root := range candidates {
			bass := root
			if ch.Inverted() {
				if held[i] != unset {
					bass = held[i]
				} else if b, ok := naming.SpellAbove(root, bassSteps(ch), ch.Bass); ok {
					bass = b
				} else {
					bass = unset
				}
			}
			// The costs weigh the spelling by degree, which the figuring
			// and the movements rest on; usage only changes the name shown.
			shown, under := usage(root, grounds[i]), usage(bass, grounds[i])
			own := cost{accidentalCost: count(root)}
			if ch.Inverted() && bass != unset {
				own[accidentalCost] += count(bass)
			}
			if hasRef && root != ref {
				own[figuringCost] = 1
			}
			if tie && root != ref {
				own[tieCost] = 1
			}
			st := state{root: root, bass: bass, shown: shown, under: under, cost: own, from: -1}
			for k, p := range prev {
				step := cost{motionCost: motion(p.bass, bass, i-1, i, grounds)}
				if total := p.cost.plus(own).plus(step); st.from < 0 || total.less(st.cost) {
					st.cost, st.from = total, k
				}
			}
			states[i] = append(states[i], st)
		}
	}

	s := respelling{roots: filled(n), basses: filled(n), strict: filled(n), strictBass: filled(n)}
	// Walk back from the best end of each run of chords between silences.
	for end := n - 1; end >= 0; end-- {
		if len(states[end]) == 0 || end+1 < n && len(states[end+1]) > 0 {
			continue
		}
		k := 0
		for j, st := range states[end] {
			if st.cost.less(states[end][k].cost) {
				k = j
			}
		}
		s.cost = s.cost.plus(states[end][k].cost)
		for i := end; i >= 0 && len(states[i]) > 0; i-- {
			st := states[i][k]
			s.roots[i], s.basses[i] = st.shown, st.under
			s.strict[i], s.strictBass[i] = st.root, st.bass
			inverted := c.Chords[i].Inverted()
			if smelly(st.shown) || inverted && smelly(st.under) {
				s.smells = append([]int{i}, s.smells...)
			}
			if double(st.root) || inverted && double(st.bass) {
				s.tolerated = append([]int{i}, s.tolerated...)
			}
			// A tolerated enharmony weighs as a smell when the zones are
			// chosen, though the report lists it apart.
			if smelly(st.shown) || inverted && smelly(st.under) || double(st.root) || inverted && double(st.bass) {
				s.stained++
			}
			k = st.from
			if k < 0 {
				break
			}
		}
	}
	return s
}

// usual keeps the spellings usage writes, without a double accidental
// and never as E♯, F♭, B♯ or C♭: B for C♭, E for F♭. A free root, one
// that P1 does not hold, is only ever spelled so; a root on a degree of
// the zone keeps its letter, the C♭ of G♭ major as much as any other.
// When no usual spelling is left, as for no class there is, the list
// comes back whole.
func usual(all []naming.SpelledNote) []naming.SpelledNote {
	var out []naming.SpelledNote
	for _, s := range all {
		if !smelly(s) {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return all
	}
	return out
}

// spellings lists the ways to write a class with one accidental at
// most, then with a double: B and C♭, C♯ and D♭, then A𝄪.
func spellings(c harmony.PitchClass) []naming.SpelledNote {
	var out []naming.SpelledNote
	for l := naming.Letter(0); l < naming.LetterCount; l++ {
		if s, ok := naming.SpellAbove(naming.SpelledNote{Letter: l}, 0, c); ok {
			out = append(out, s)
		}
	}
	return out
}

// motion is 1 when the movement from bass `a`, change `i`, to bass `b`,
// change `j`, reads poorly, 0 otherwise: the bass, the root of a chord
// in root position, since it is the line the ear follows. G♯dim7
// F/A is G♯ up to A, not G♯ down to F. A movement reads plainly when
// its letters span the interval it sounds: a fifth or a fourth is
// perfect, a third spans three letters, a step two. A semitone may also
// stay on its letter, a chromatic step; and a note outside the scale of
// its ground goes the way of the line, never a flat going up nor a
// sharp going down. The same root twice keeps its name.
func motion(a, b naming.SpelledNote, i, j int, grounds [][]harmony.Tonality) int {
	up := (int(b.Class()) - int(a.Class()) + 12) % 12
	d := (int(b.Letter) - int(a.Letter) + 7) % 7
	ok := false
	switch up {
	case 0:
		ok = a == b
	case 1:
		ok = (d == 1 || d == 0) && !against(a, 1, grounds[i]) && !against(b, 1, grounds[j])
	case 11:
		ok = (d == 6 || d == 0) && !against(a, -1, grounds[i]) && !against(b, -1, grounds[j])
	case 2:
		ok = d == 1
	case 10:
		ok = d == 6
	case 3, 4:
		ok = d == 2
	case 8, 9:
		ok = d == 5
	case 5:
		ok = d == 3
	case 7:
		ok = d == 4
	case 6:
		ok = d == 3 || d == 4
	}
	if ok {
		return 0
	}
	return 1
}

// against reports whether `n`, outside the scales of its ground, is
// written against a line going `dir`: a flat going up, a sharp going
// down.
func against(n naming.SpelledNote, dir int, ground []harmony.Tonality) bool {
	if _, in := degreeIn(ground, n.Class()); in {
		return false
	}
	return dir > 0 && n.Accidental < 0 || dir < 0 && n.Accidental > 0
}

// pedals spells the held basses (see [analysis.Pedals]): each takes the
// letter of its degree in the zone, or failing that its plainest
// spelling, and keeps it while it is held. The E♭ under A♭/E♭ B♭m/E♭
// Cm/E♭ D♭/E♭ E♭7sus in A♭.
func pedals(c analysis.Changes, ps []analysis.Pedal, grounds [][]harmony.Tonality, zones []int, tonics []naming.SpelledNote) []naming.SpelledNote {
	out := filled(len(c.Chords))
	for _, p := range ps {
		spelled := spellings(p.Bass)[0]
		for _, s := range spellings(p.Bass) {
			if count(s) < count(spelled) || count(s) == count(spelled) && s.Accidental < spelled.Accidental {
				spelled = s
			}
		}
		if z := zones[p.From]; z >= 0 {
			if n, ok := degreeIn(grounds[p.From], p.Bass); ok {
				if s, ok := naming.SpellAbove(tonics[z], int(n)-1, p.Bass); ok {
					spelled = s
				}
			}
		}
		for k := p.From; k <= p.To; k++ {
			out[k] = spelled
		}
	}
	return out
}

// degreeIn returns the degree class `c` is in the first scale of `ts`
// that holds it.
func degreeIn(ts []harmony.Tonality, c harmony.PitchClass) (harmony.Degree, bool) {
	for _, t := range ts {
		s := harmony.Semitones((int(c) - int(t.Tonic()) + 12) % 12)
		for n, off := range t.Pattern().Offsets() {
			if off == s {
				return n, true
			}
		}
	}
	return 0, false
}

// usage writes a note as usage does once the spelling is settled, on the
// next letter when it carries a double accidental, a tolerated enharmony
// the report lists (see [tolerated]), or when it is an E♯, F♭, B♯ or
// C♭: A for B𝄫, B for C♭, E for F♭. The figuring and the movements
// keep the spelling by degree; only the name shown changes, and with it
// what smells. The raised leading tone of a minor, the seventh degree
// of its harmonic or melodic scale, keeps its letter when it takes a
// single sharp: in F♯ minor, E♯ goes up to F♯.
func usage(n naming.SpelledNote, ground []harmony.Tonality) naming.SpelledNote {
	if double(n) {
		return nextLetter(n)
	}
	if count(n) != 1 || n.Class() != 0 && n.Class() != 4 && n.Class() != 5 && n.Class() != 11 {
		return n
	}
	if n.Accidental > 0 {
		for _, t := range ground {
			p := t.Pattern()
			if (p == harmony.ScaleHarmonicMinor || p == harmony.ScaleMelodicMinor) && n.Class() == t.Tonic().Transpose(11) {
				return n
			}
		}
	}
	return nextLetter(n)
}

// nextLetter writes `n` on the next letter toward its accidental, with
// one sign fewer: A for B𝄫, G for F𝄪, B for C♭.
func nextLetter(n naming.SpelledNote) naming.SpelledNote {
	step := 1
	if n.Accidental < 0 {
		step = -1
	}
	if m, ok := naming.SpellAbove(n, step, n.Class()); ok {
		return m
	}
	return n
}

// double reports whether `n` is written with a double accidental.
func double(n naming.SpelledNote) bool {
	return count(n) == 2
}

// smelly reports whether `n` is written with a double accidental, or as
// E♯, F♭, B♯ or C♭.
func smelly(n naming.SpelledNote) bool {
	switch {
	case count(n) == 2:
		return true
	case count(n) == 1:
		switch n.Class() {
		case 0, 4, 5, 11:
			return true
		}
	}
	return false
}

func name(n naming.SpelledNote) string {
	return naming.English.Name(n, naming.Signs)
}

// count is the number of signs a note writes, two for a double.
func count(n naming.SpelledNote) int {
	return max(int(n.Accidental), -int(n.Accidental))
}

// bassSteps counts the letters from the root of a chord to its bass, by
// the interval the bass makes with the root: the third, a letter of two
// steps, the diminished fifth of a m7♭5 on four, the diminished seventh
// of a dim7 on six where the sixth of a C6/A takes five. A bass outside
// the chord is a ninth or an eleventh: the G of F/G is the ninth of F,
// one letter up.
func bassSteps(ch analysis.Change) int {
	p := ch.Chord.Pattern
	switch (int(ch.Bass) - int(ch.Chord.Root) + 12) % 12 {
	case 1, 2:
		return 1
	case 3, 4:
		return 2
	case 5:
		return 3
	case 6, 7:
		return 4
	case 8:
		// The ♭6 of a minor chord with its fifth, the ♯5 of an augmented
		// one without.
		if p.HasOffset(7) {
			return 5
		}
		return 4
	case 9:
		if p == harmony.ChordDiminishedSeventh {
			return 6
		}
		return 5
	case 10, 11:
		return 6
	}
	return 0
}

// symmetric reports whether a change is a chord that divides the octave
// equally, a dim7 or an augmented triad, over one of its own notes:
// such a chord is the same chord on its bass, and usage writes it so,
// Adim7 for Cdim7/A, E+ for C+/E. Only the name shown changes.
func symmetric(ch analysis.Change) bool {
	p := ch.Chord.Pattern
	if !ch.Inverted() || p != harmony.ChordDiminishedSeventh && p != harmony.ChordAugmentedTriad {
		return false
	}
	return p.HasOffset(harmony.Semitones((int(ch.Bass) - int(ch.Chord.Root) + 12) % 12))
}

// tonic spells the tonic of tonalities. Heard around change `at`, it is
// spelled as the chords there are, the root of the nearest change on the
// tonic, so that a tonicisation and the chord it tonicises bear the same
// name, [F♯m] under F♯m7 (P3 of [respell]). Heard over the whole tune,
// `at` negative, it is the tonic of a zone as the spelling chose it, or
// a root on that class as it is written first, or the fewest
// accidentals.
func (n namer) tonic(ts []harmony.Tonality, at int) string {
	c := ts[0].Tonic()
	if at >= 0 {
		for d := 0; d < len(n.roots); d++ {
			for _, j := range []int{at + d, at - d} {
				if j >= 0 && j < len(n.roots) && n.roots[j] != unset && n.roots[j].Class() == c {
					return name(n.roots[j])
				}
			}
		}
	}
	if t, ok := n.tonics[c]; ok {
		return name(t)
	}
	return name(naming.TonicSpelling(ts[0]))
}

// of names tonalities as [short] does, their tonic spelled by [tonic]:
// "F♯m", "D♭(m)".
func (n namer) of(ts []harmony.Tonality, at int) string {
	return short(n.tonic(ts, at), ts)
}

// smells counts the changes whose root or bass smells, and says where:
// "smells: 2, C♭7 bar 4, C♭7♯11 bar 6", then the tolerated enharmonies:
// "tolerated enharmony: Adim7 for Cdim7/B𝄫 bar 12". An empty string when
// there is neither.
func smells(s respelling, tl ireal.Timeline) string {
	var b strings.Builder
	if len(s.smells) > 0 {
		var where []string
		for _, i := range s.smells {
			where = append(where, fmt.Sprintf("%s bar %d", symbol(tl.Spans[i].Chord), tl.Spans[i].Bar+1))
		}
		fmt.Fprintf(&b, "smells: %d, %s\n", len(s.smells), strings.Join(where, ", "))
	}
	if len(s.tolerated) > 0 {
		var where []string
		for _, i := range s.tolerated {
			// The chord by degree, inverted again when usage showed a
			// symmetric chord on its bass.
			strict := tl.Spans[i].Chord
			strict.Root = name(s.strict[i])
			if c := s.strictBass[i]; c != unset && c != s.strict[i] {
				strict.Bass = name(c)
			}
			where = append(where, fmt.Sprintf("%s for %s bar %d", symbol(tl.Spans[i].Chord), symbol(strict), tl.Spans[i].Bar+1))
		}
		fmt.Fprintf(&b, "tolerated enharmony: %s\n", strings.Join(where, ", "))
	}
	return b.String()
}
