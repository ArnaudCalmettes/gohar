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
func TonalAreas(c Changes, blocks []Block, sensed []Sensed, sections []Section) []TonalArea {
	heard := func(s Sensed) []harmony.Tonality {
		if s.Region != nil {
			return s.Region
		}
		return s.Ground
	}
	type run struct {
		from, to int
		tonic    []harmony.Tonality
	}
	var runs []run
	for i := 0; i < len(sensed); {
		t := heard(sensed[i])
		j := i
		for j+1 < len(sensed) && sameTonic(heard(sensed[j+1]), t) {
			j++
		}
		if t != nil {
			runs = append(runs, run{i, j, t})
		}
		i = j + 1
	}
	r := rolesOf(c, blocks)
	for k := 1; k < len(runs); k++ {
		n := r.target[runs[k].from]
		if n < 0 || !sameTonic(blocks[n].Announced, runs[k].tonic) {
			continue
		}
		if from := firstOf(blocks[n]); from > runs[k-1].from && from < runs[k].from {
			runs[k].from, runs[k-1].to = from, from-1
		}
	}
	var out []TonalArea
	if len(runs) == 0 {
		return nil
	}
	first := runs[0].tonic
	before := first // the tonic heard before, a passing one left out
	for _, rn := range runs[1:] {
		var held Ticks
		for i := rn.from; i <= rn.to; i++ {
			held += c.Chords[i].Length
		}
		switch {
		case sameTonic(rn.tonic, first):
			before = first
		case held >= 2*barOf(c) && len(out) > 0 && sameTonic(rn.tonic, before):
			// The same area, after a passing tonic left out.
			last := &out[len(out)-1]
			*last = area(c, sections, last.From, rn.to, rn.tonic, last.Leaves, first)
		case held >= 2*barOf(c):
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
	}
	if before != nil {
		d := signature(t) - signature(before)
		d = (d%12 + 12) % 12
		a.Distance = min(d, 12-d)
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
