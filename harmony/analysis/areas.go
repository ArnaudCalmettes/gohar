package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A TonalArea is a stretch of a tune heard around a tonic other than the
// one it sets out from: a modulation as En Harmonie means it, a region
// or a ground a modulation has installed (see [Sensed]). It carries
// what Siron weighs to tell a true modulation from a transitory one (La
// partition intérieure, p. 380), measured and not yet weighed:
//
//   - the duration: « la durée de la modulation est importante pour
//     distinguer la sensation de modulation vraie d'une modulation
//     transitoire »;
//   - the memory: « la première tonalité entendue a toujours un énorme
//     poids »;
//   - the place in the form: « une modulation vraie est plus volontiers
//     située au début ou à la fin d'un cycle de mesures ou d'une phrase
//     mélodique. Au milieu d'une phrase harmonique, l'oreille entend
//     plutôt une modulation transitoire »;
//   - the distance: « la sensation de modulation est plus forte si les
//     tonalités sont éloignées », counted in key signatures on the cycle
//     of fifths, a tonality and its relative sharing one (p. 383).
type TonalArea struct {
	From, To int // its first and last changes

	// Tonic is the tonic heard, Leaves the one heard before it.
	Tonic, Leaves []harmony.Tonality

	// Bars is how long it lasts, in bars.
	Bars float64

	// First tells an area that leaves the first tonality heard.
	First bool

	// Opens tells an area that starts on the first bar of a section,
	// Closes one that ends with a section.
	Opens, Closes bool

	// Distance is the number of steps on the cycle of fifths between
	// the key signatures of Leaves and Tonic: 0 between relatives, 1
	// between neighbours, 6 at the tritone.
	Distance int

	// True tells a true modulation from a transitory one: the area holds
	// the first bar of a section and half the section at least. Siron
	// places a true modulation « plus volontiers […] au début ou à la fin
	// d'un cycle de mesures ou d'une phrase mélodique », and hears a
	// transitory one « au milieu d'une phrase harmonique » (p. 380); the
	// length weighs, « la durée de la modulation est importante ». His
	// true modulation, the bridge of Body and Soul, from D flat to D
	// (p. 382), fills the bridge; his transitory one, the C major of All
	// The Things You Are, a « respiration secondaire » (p. 379), ends the
	// first section without opening one. Where a section begins and
	// half its length are ours to set. The memory, his third criterion,
	// weighs nothing yet; the distance, his fourth, only keeps an area
	// from forming (see [TonalAreas]).
	True bool
}

// TonalAreas returns the stretches of the changes heard around another
// tonic than the first one, given the blocks and what was sensed at
// each change: its region when there is one, else its ground. Heard
// afterwards, as [Grounds] does, an area starts with the cadence that
// led to it: in Tune Up, C major from Dm7 on.
//
// An area is kept when it lasts two bars or more, the cadence that
// leads to it counted in. « Une cadence prépare généralement la
// modulation, celle-ci étant confirmée si la durée est assez longue
// pour que la nouvelle tonalité soit installée » (En Harmonie, tome 1,
// chapter 10 §1.8, p. 159): the book gives no length, and two bars are
// ours to set. Gm7 C7 | Fm7, at bars 30 and 31 of There Will Never Be
// Another You, gone by the next chord, is a tonicisation, as the book
// brackets it; the four bars of C major in Tune Up, and the seven of
// Black Orpheus, are modulations, and so is F♯m7 | B7 | Emaj7 in the
// bridge of All The Things You Are, one of Siron's (La partition
// intérieure, p. 387).
//
// Two rules keep distant centres from installing anything. The
// tonalities a dominant borrows are, « dans une harmonie peu
// chromatique », the neighbours of the tonality (Siron, p. 342); a two
// five that resolves on the tonic chord of a distant one, two steps
// away or more on the cycle of fifths, is heard as a centre of its
// own. It cuts the area it sounds in, and is no area itself, a
// tonicisation « ne porte que sur un accord » (p. 379). And
// « l'enchaînement rapproché de centres tonaux éloignés détruit la
// sensation d'une véritable modulation » (p. 380): a stretch of two
// bars at most, between two distant centres as short as it is, is
// left out. Giant Steps, three major tonalities a major third apart,
// « n'est composé que de cadences » (p. 533): none of them installs
// its tonality. Two steps, and the length of two bars, are ours to set.
func TonalAreas(c Changes, blocks []Block, sensed []Sensed, sections []Section) []TonalArea {
	heard := func(i int) ([]harmony.Tonality, bool) {
		s, base := sensed[i], sensed[i].Ground
		if s.Region != nil {
			base = s.Region
		}
		if t := s.Tonicised; t != nil && base != nil && !sameTonic(t, base) &&
			!isDominant(c.Chords[i].Chord.Pattern) && steps(t, base) > 1 {
			return t, true
		}
		return base, false
	}
	type run struct {
		from, to  int
		tonic     []harmony.Tonality
		tonicised bool // a distant tonicisation, heard as a centre
	}
	var runs []run
	for i := 0; i < len(sensed); {
		t, tonicised := heard(i)
		j := i
		for j+1 < len(sensed) {
			if next, _ := heard(j + 1); !sameTonic(next, t) {
				break
			}
			j++
		}
		if t != nil {
			runs = append(runs, run{i, j, t, tonicised})
		}
		i = j + 1
	}
	r := rolesOf(c, blocks)
	for k := 1; k < len(runs); k++ {
		n := r.target[runs[k].from]
		if n < 0 || !sameTonic(blocks[n].Announced, runs[k].tonic) {
			continue
		}
		from := firstOf(blocks[n])
		switch {
		case from > runs[k-1].from && from < runs[k].from:
			runs[k].from, runs[k-1].to = from, from-1
		case from == runs[k-1].from && k > 1:
			// The run before is only the cadence: Cm7 F7 heard in D, in
			// Tune Up, between a region of C and one of B flat.
			runs[k].from = from
			runs = append(runs[:k-1], runs[k:]...)
			k--
		}
	}
	var out []TonalArea
	if len(runs) == 0 {
		return nil
	}
	held := make([]Ticks, len(runs))
	for k, rn := range runs {
		for i := rn.from; i <= rn.to; i++ {
			held[k] += c.Chords[i].Length
		}
	}
	if last := len(runs) - 1; c.Loops && last > 0 {
		// Across the loop, the cadence that leads back to the first
		// chord ends the tune: the C♯m7 F♯7 of the last bar of Giant
		// Steps, toward its Bmaj7.
		if n := r.target[runs[0].from]; n >= 0 && sameTonic(blocks[n].Announced, runs[0].tonic) {
			if from := firstOf(blocks[n]); from > runs[last].from {
				for i := from; i <= runs[last].to; i++ {
					held[last] -= c.Chords[i].Length
					held[0] += c.Chords[i].Length
				}
				runs[last].to = from - 1
			}
		}
	}
	passing := func(k int) bool {
		short := func(k int) bool { return held[k] <= 2*barOf(c) }
		near := func(j int) bool {
			if c.Loops {
				j = (j + len(runs)) % len(runs)
			}
			return j >= 0 && j < len(runs) && j != k &&
				short(j) && steps(runs[j].tonic, runs[k].tonic) > 1
		}
		return short(k) && near(k-1) && near(k+1)
	}
	first := runs[0].tonic
	before := first // the tonic heard before, a passing one left out
	for k, rn := range runs {
		if k == 0 {
			continue
		}
		switch {
		case rn.tonicised:
			// A tonicisation « ne porte que sur un accord »: it cuts the
			// area it sounds in and is none.
		case passing(k):
			// A centre between two distant ones, as short as they are.
		case sameTonic(rn.tonic, first):
			before = first
		case held[k] >= 2*barOf(c) && len(out) > 0 && sameTonic(rn.tonic, before):
			// The same area, after a passing tonic left out.
			last := &out[len(out)-1]
			*last = area(c, sections, last.From, rn.to, rn.tonic, last.Leaves, first)
		case held[k] >= 2*barOf(c):
			out = append(out, area(c, sections, rn.from, rn.to, rn.tonic, before, first))
			before = rn.tonic
		}
	}
	return out
}

// area measures the tonal area of changes i to j around tonic t, after
// tonic `before`, in a tune that set out from `first`.
func area(c Changes, sections []Section, i, j int, t, before, first []harmony.Tonality) TonalArea {
	a := TonalArea{From: i, To: j, Tonic: t, Leaves: before, First: sameTonic(before, first)}
	var length Ticks
	for k := i; k <= j; k++ {
		length += c.Chords[k].Length
	}
	a.Bars = float64(length) / float64(barOf(c))
	start := c.Bar(c.Chords[i].Start)
	end := c.Bar(c.Chords[j].Start + c.Chords[j].Length - 1)
	for _, s := range sections {
		a.Opens = a.Opens || start == s.From
		a.Closes = a.Closes || end == s.From+s.Bars-1
		if start <= s.From && end >= s.From {
			a.True = a.True || 2*(min(end, s.From+s.Bars-1)-s.From+1) >= s.Bars
		}
	}
	if before != nil {
		a.Distance = steps(t, before)
	}
	return a
}

// signature places a tonality on the cycle of fifths, by its key
// signature: C major and A minor at 0, G major at 1, F major at 11. A
// minor tonic is placed by its relative major.
func signature(t []harmony.Tonality) int {
	tonic := int(t[0].Tonic())
	if ModesOf(t) == Minor {
		tonic += 3
	}
	return tonic * 7 % 12
}

// steps counts the steps between two tonalities on the cycle of fifths,
// by their key signatures: 0 between relatives, 6 at the tritone.
func steps(a, b []harmony.Tonality) int {
	d := signature(a) - signature(b)
	d = (d%12 + 12) % 12
	return min(d, 12-d)
}
