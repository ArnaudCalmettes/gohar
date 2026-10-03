package main

import (
	"math/rand/v2"
	"slices"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// The reference bass (see "La basse de référence" in docs/walk.md): a
// walking line generated from the grid, bar by bar, with the rules of
// Jeremy Siskind, Jazz Piano Fundamentals, Book 2, Unit 10, p. 204 to
// 206. Which rule applies depends on the harmonic rhythm of the bar and
// on where the next chord lies; when a rule offers several paths, one
// is drawn at random, so that the line does not loop on one formula.
//
// The generator and the marker share their rules: a generated line
// must land every arrival, and the tests check it does.

// Walk returns a key per beat of `beats`, in the register of the double
// bass, E1 to G2, climbing to C3 at most. `rng` draws among the paths a rule allows: a fixed
// seed gives the same line every time.
//
// Each bar offers its paths, as pitch classes, each placed on keys after
// the line so far (see step), from every octave its first note may
// take; one of those that fit wins the draw (see choose).
func Walk(beats []Beat, rng *rand.Rand) []int {
	keys := make([]int, len(beats))
	prev, dir := 0, 0
	played := map[harmony.PitchClass][]int{} // the last bar over each root
	for bar := 0; bar < len(beats); bar += perBar {
		end := min(bar+perBar, len(beats))
		root := beats[bar].Chord.Bass
		p, placed := choose(paths(beats[bar:end], target(beats, end)), beats[bar:end], prev, dir, played[root], rng)
		copy(keys[bar:end], placed)
		prev, dir = placed[len(placed)-1], p.dir
		played[root] = placed
	}
	return keys
}

// A path is a bar of the line, as pitch classes. `dir` tells a scale
// run, +1 up and -1 down, which the next bar may carry on; 0 for the
// other paths.
type path struct {
	pcs []harmony.PitchClass
	dir int
}

// Ending returns the key the line ends on, after its last beat, or 0
// for none. A chart that loops and does not say where it ends (`end`,
// analysis.Changes.End, is 0) would stop on its turnaround, on a
// dominant: the line lands on the chord the turnaround leads to, where
// the chorus would loop back, unless the last chord is already it.
func Ending(beats []Beat, line []int, end int) int {
	last := len(beats) - 1
	if end != 0 || last < 0 {
		return 0
	}
	pc := beats[last].Next.Bass
	if pc == beats[last].Chord.Bass {
		return 0
	}
	k, _ := step(pc, line[last], 0, pc, anyLeap)
	return k
}

// target is the root the line heads for after the beats before `end`:
// the bass of the chord on the next beat, or the chord after the last
// one when the run ends there.
func target(beats []Beat, end int) harmony.PitchClass {
	if end < len(beats) {
		return beats[end].Chord.Bass
	}
	return beats[end-1].Next.Bass
}

// carryOn is how often a run carries on into the next bar when a path
// allows it: always would replay the same walk down bar after bar.
const carryOn = 2.0 / 3

// choose draws a bar among the paths `ps` placed after `prev`:
//   - never `last`, the bar played over the same root before, when
//     anything else fits: a demo is to show the variety of a line;
//   - two times in three, one carrying on the run `dir` when there is.
//
// When nothing fits, a path placed with the looser rules of step.
func choose(ps []path, bar []Beat, prev, dir int, last []int, rng *rand.Rand) (path, []int) {
	type fit struct {
		p    path
		keys []int
	}
	var fits, fresh, runs []fit
	for _, p := range ps {
		for _, keys := range place(p, bar, prev, strict) {
			fits = append(fits, fit{p, keys})
		}
	}
	for _, f := range fits {
		if !slices.Equal(f.keys, last) {
			fresh = append(fresh, f)
		}
	}
	if len(fresh) > 0 {
		fits = fresh
	}
	for _, f := range fits {
		if dir != 0 && f.p.dir == dir {
			runs = append(runs, f)
		}
	}
	if len(runs) > 0 && rng.Float64() < carryOn {
		fits = runs
	}
	if len(fits) > 0 {
		f := fits[rng.IntN(len(fits))]
		return f.p, f.keys
	}
	p := ps[rng.IntN(len(ps))]
	for _, l := range []looseness{wide, anyLeap} {
		if placed := place(p, bar, prev, l); len(placed) > 0 {
			return p, placed[0]
		}
	}
	panic("walk: no key for a path") // anyLeap always finds one
}

// paths lists what a bar may play, by the rule its harmonic rhythm
// calls for.
func paths(bar []Beat, next harmony.PitchClass) []path {
	if len(bar) < perBar {
		pcs := make([]harmony.PitchClass, len(bar))
		for i, b := range bar {
			pcs[i] = b.Chord.Bass
		}
		return []path{{pcs: pcs}}
	}

	switch changes := arrivals(bar); {
	case changes == 1:
		return whole(bar[0], next)

	case changes == 2 && bar[2].Arrives:
		// Rule 2, two chords a bar: each root, then the fifth, the third,
		// or a half step into the next root.
		var ps []path
		for _, a := range half(bar[0], bar[2].Chord.Bass) {
			for _, b := range half(bar[2], next) {
				ps = append(ps, path{pcs: []harmony.PitchClass{bar[0].Chord.Bass, a, bar[2].Chord.Bass, b}})
			}
		}
		return ps

	default:
		// Rule 1, four chords a bar: the roots. Any other rhythm, a chord
		// on beat 2 or 4, gets the same treatment: the root where a chord
		// arrives, its fifth elsewhere.
		pcs := make([]harmony.PitchClass, len(bar))
		for i, b := range bar {
			pcs[i] = b.Chord.Bass
			if i > 0 && !b.Arrives {
				pcs[i] = fifth(b)
			}
		}
		return []path{{pcs: pcs}}
	}
}

// arrivals counts the chords that start in the bar, beat 1 included
// even when the chord carries on from the bar before.
func arrivals(bar []Beat) int {
	n := 0
	for i, b := range bar {
		if b.Arrives || i == 0 {
			n++
		}
	}
	return n
}

// half lists what may follow the root of `b` on the second of two
// beats, heading for `next`.
func half(b Beat, next harmony.PitchClass) []harmony.PitchClass {
	return []harmony.PitchClass{fifth(b), third(b), next.Transpose(1), next.Transpose(-1)}
}

// whole lists the paths of four beats under one chord, heading for
// `next`.
func whole(b Beat, next harmony.PitchClass) []path {
	r := b.Chord.Bass
	s := scaleOf(b)
	up := path{[]harmony.PitchClass{r, above(s, r, 1), above(s, r, 2), above(s, r, 3)}, +1}
	down := path{[]harmony.PitchClass{r, below(s, r, 1), below(s, r, 2), below(s, r, 3)}, -1}
	switch {
	case r.Up(next) == 5:
		// Rule 3, the cycle of fifths: the next root a fifth below. The
		// walk up, a tone then three half steps (G A A♯ B, C); the walk
		// down, four degrees of the chord's scale (D C B A, G); the
		// triad, then a half step above the target (E♭ G B♭ A, A♭).
		return []path{
			{[]harmony.PitchClass{r, r.Transpose(2), r.Transpose(3), r.Transpose(4)}, +1},
			down,
			{pcs: []harmony.PitchClass{r, third(b), fifth(b), next.Transpose(1)}},
		}

	case next == r:
		// Rule 4, the same chord over several bars: the triad arpeggio,
		// or the scale climbing toward the fifth. The scale coming down
		// to the fifth below is ours, the mirror, so that a walk down can
		// carry on into a held chord.
		return []path{
			{pcs: []harmony.PitchClass{r, third(b), fifth(b), third(b)}},
			up,
			down,
		}

	default:
		// Rule 5, any other motion: chord tones and a half step into the
		// target. Siskind's bass in two doubled, R R 5 5, would leap an
		// octave on the fifth: only the root is repeated here.
		var ps []path
		for _, a := range []harmony.PitchClass{next.Transpose(1), next.Transpose(-1)} {
			ps = append(ps,
				path{pcs: []harmony.PitchClass{r, r, fifth(b), a}},
				path{pcs: []harmony.PitchClass{r, third(b), fifth(b), a}})
		}
		return ps
	}
}

// third and fifth are the chord's own, minor or major, diminished or
// perfect: what its pattern holds.
func third(b Beat) harmony.PitchClass {
	return chordTone(b, 3, 4)
}

func fifth(b Beat) harmony.PitchClass {
	return chordTone(b, 7, 6, 8)
}

// chordTone returns the first interval of `options` the chord holds,
// above its root; the first option when it holds none.
func chordTone(b Beat, options ...harmony.Semitones) harmony.PitchClass {
	r := b.Chord.Chord.Root
	for _, o := range options {
		if b.Chord.Chord.Pattern.HasOffset(o) {
			return r.Transpose(o)
		}
	}
	return r.Transpose(options[0])
}

// scaleOf is the scale the bassist walks over a chord: Siron's
// "conservative" choice (La partition intérieure, p. 693), dorian on a
// minor seventh, mixolydian on a dominant, ionian on a major chord,
// locrian on a half-diminished one. As offsets from the root.
func scaleOf(b Beat) []harmony.Semitones {
	p := b.Chord.Chord.Pattern
	switch {
	case p.HasOffset(3) && p.HasOffset(6):
		return []harmony.Semitones{0, 1, 3, 5, 6, 8, 10} // locrian
	case p.HasOffset(3):
		return []harmony.Semitones{0, 2, 3, 5, 7, 9, 10} // dorian
	case p.HasOffset(11) || !p.HasOffset(10):
		return []harmony.Semitones{0, 2, 4, 5, 7, 9, 11} // ionian
	}
	return []harmony.Semitones{0, 2, 4, 5, 7, 9, 10} // mixolydian
}

// above and below step `n` degrees of `scale` up or down from its root
// `r`.
func above(scale []harmony.Semitones, r harmony.PitchClass, n int) harmony.PitchClass {
	return r.Transpose(scale[n%len(scale)])
}

func below(scale []harmony.Semitones, r harmony.PitchClass, n int) harmony.PitchClass {
	return r.Transpose(scale[(len(scale)-n%len(scale))%len(scale)])
}

// Register of the line: the four strings of the double bass, E1 to G2
// (Siskind, p. 116 and 204), preferred. The line may climb past them up
// to C3, where the right hand starts (p. 118): otherwise an F2 or a G2,
// near the top, leaves a bar nothing but a walk down. An octave on the
// root may also go down to C1, the lowest Siskind allows (p. 117).
const (
	lowE1  = 28
	highG2 = 43

	lowC1  = 24
	highC3 = 48
)

// The leaps the line allows: a fifth at most, an octave on the root.
const (
	maxLeap = 7
	octave  = 12
)

// A looseness is how far step may bend the rules when no path fits.
type looseness int

const (
	strict  looseness = iota
	wide              // down to C1 for any note
	anyLeap           // any leap, any direction, C1 to C3
)

// place puts the pitch classes of `p` on keys after `prev`, once for
// each key its first note may take: F1 or F2 after a C2, say, both a
// fifth or less away. The other notes follow, each on the key step
// prefers. Nil when no placement fits.
func place(p path, bar []Beat, prev int, l looseness) [][]int {
	var out [][]int
next:
	for _, first := range candidates(p.pcs[0], prev, 0, bar[0].Chord.Bass, l) {
		keys := []int{first}
		for i := 1; i < len(p.pcs); i++ {
			k, ok := step(p.pcs[i], keys[i-1], p.dir, bar[i].Chord.Bass, l)
			if !ok {
				continue next
			}
			keys = append(keys, k)
		}
		out = append(out, keys)
	}
	return out
}

// step returns the key candidates prefers for `pc` after `prev`.
func step(pc harmony.PitchClass, prev, dir int, root harmony.PitchClass, l looseness) (int, bool) {
	ks := candidates(pc, prev, dir, root, l)
	if len(ks) == 0 {
		return 0, false
	}
	return ks[0], true
}

// candidates returns the keys of `pc` the rules allow after `prev`, the
// preferred first: within the four strings, then the nearest. A fifth
// away at most, in direction `dir` when it is not 0, from E1 to C3. The
// same key twice in a row is never allowed (Siskind, p. 116); the
// octave is, when `pc` is `root`, and it may go down to C1 for that.
// Without a previous note, any key of the four strings.
func candidates(pc harmony.PitchClass, prev, dir int, root harmony.PitchClass, l looseness) []int {
	var ks []int
	for k := lowC1; k <= highC3; k++ {
		if k%12 != int(pc) || k == prev {
			continue
		}
		in := k >= lowE1 && k <= highG2
		if prev == 0 {
			if in {
				ks = append(ks, k)
			}
			continue
		}
		leap := abs(k - prev)
		eighth := leap == octave && pc == root
		switch {
		case l < anyLeap && dir*(k-prev) < 0:
			continue
		case l < anyLeap && leap > maxLeap && !eighth:
			continue
		case l == strict && !eighth && (k < lowE1 || k > highC3):
			continue
		}
		ks = append(ks, k)
	}
	inside := func(k int) bool { return k >= lowE1 && k <= highG2 }
	slices.SortStableFunc(ks, func(a, b int) int {
		if inside(a) != inside(b) {
			if inside(a) {
				return -1
			}
			return 1
		}
		return abs(a-prev) - abs(b-prev)
	})
	return ks
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
