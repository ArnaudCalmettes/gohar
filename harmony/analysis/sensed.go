package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Sensed is the tonic an ear expects once a chord has sounded, from
// what came before it and nothing after: the tonique pressentie of
// docs/grilles.md. At the third chord of Tenderly, two bars of E♭maj7
// A♭7 have installed E flat, and E♭m7 is heard as the tonic changing
// colour, not as the two of D flat.
//
// [Sense] reads it from left to right, the other way from the rest of
// the analysis: the region is what one concludes afterwards (see
// [Grounds]), the sensed tonic what one hears now. The gap between
// what it awaited and what comes is the surprise.
//
// Each set of tonalities is on one tonic, as a block announces them,
// and nil when there is none.
type Sensed struct {
	// Start is the first tonic installed, kept apart from the ground to
	// hear the return home after a bridge that modulated.
	Start []harmony.Tonality

	// Ground is the installed tonic, the one degrees count from, even
	// when a cadence tonicises another degree.
	Ground []harmony.Tonality

	// Since is the change the ground holds from, heard afterwards: the
	// first chord of the cadence that led to it, before the ear knew.
	Since int

	// Local is the tonic a cadence has just tonicised, when it is not
	// the ground; it lasts while the chords after it hold in it, or
	// prepare a chord that does.
	Local []harmony.Tonality

	// Awaited is what a cadence being played announces: the tonalities
	// of the block the chord belongs to.
	Awaited []harmony.Tonality

	// Tonic is what the chord is heard as the tonic of, when it can be
	// one there: see [Sense].
	Tonic []harmony.Tonality

	// Resolves is, on the V of a cadence that resolves on a tonic, the
	// tonalities of that tonic; Across tells that it lands across the
	// loop, on the first chord of the chart.
	Resolves []harmony.Tonality
	Across   bool
}

// Sense reads the sensed tonic at each change.
//
// # What installs a tonic
//
// The first chord gives the ground to start from when it can be a
// tonic chord (a major or minor triad, maj7, 6, m6, m(maj7), but not
// m7, almost always a subdominant), and else the first cadence that
// resolves. A cadence that resolves elsewhere than on the ground makes
// its target a local tonic.
//
// A blues is recognised by its form (see [Blues]), and its tonic is
// the ground from the start: its I7 is a seventh of kind, not a
// dominant waiting for a cadence.
//
// A key signature plays no part: it belongs to the written score, the
// weakest clue to a tonality, and an analyst who reads the chords
// finds the tonality for himself. The cadences are enough: the fiches
// of En Harmonie read the same with and without it.
//
// # Modulation
//
// A local tonic becomes the ground when it holds beyond one bar of
// stable chords, those that hold in it without preparing anything
// (B♭maj7 Gm7 in Tune Up), or when a second cadence resolves on it
// while it lasts (Black Orpheus, C major). Its I must be able to be a
// tonic: one does not modulate toward a subdominant, nor for one
// chord. A tonic held a single bar only tonicises: Along Came Betty
// dances from A flat to A and back, and never modulates.
//
// A plagal cadence (the IV, or the ♭VII7, before the tonic) concludes
// as a V-I does, and draws less: it confirms a tonic already there,
// the ground, home or the local tonic, and brings it back, but opens
// none.
//
// A cadence toward a degree of the local tonic that cannot be a tonic
// (A7 Dm7 while C major is local) does not interrupt it. Coming home
// is easier than leaving: one cadence that resolves on the start tonic
// brings the ground back to it, and so does its tonic chord alone,
// however it comes.
//
// # A two on the ground is its I
//
// A two five that does not resolve and whose two sits on the ground's
// tonic awaits nothing: its two is the tonic, borrowed from another
// scale. E♭m7 A♭7 in E flat is I then IV7, not a two five of D flat.
//
// # Local
//
// Each change sees the state left by those before it and the blocks,
// whose tonalities a target decides between: one change ahead, as the
// analysis of a performance allows.
func Sense(c Changes, blocks []Block) []Sensed {
	target := make([]int, len(c.Chords))
	member := make([]int, len(c.Chords))
	for i := range target {
		target[i], member[i] = -1, -1
	}
	for n, b := range blocks {
		if b.Target >= 0 {
			target[b.Target] = n
		}
		for _, i := range []int{b.Two, b.Sus, b.Five} {
			if i >= 0 {
				member[i] = n
			}
		}
	}
	bar := 4 * TicksPerBeat
	if len(c.Bars) > 1 {
		bar = c.Bars[1] - c.Bars[0]
	}

	// What a local tonic has gathered so far toward becoming the ground.
	var (
		tonic     []harmony.Tonality // the tonalities its I is the tonic of
		since     int                // the first chord of its first cadence
		cadences  int                // cadences resolved on it
		stable    Ticks              // time of its stable chords
		installed = func(s *Sensed) {
			s.Ground, s.Since, s.Local = tonic, since, nil
		}
	)

	// tonicAt tells what change i can be the tonic of: a tonic chord
	// (see [tonicOf]), or a minor seventh that a minor cadence resolves
	// on and that is not itself the two of a block. Charts write the
	// minor tonic m7 far more often than it is played (m6, m(maj7)):
	// Cm7 | Dm7♭5 G7♭9 | Cm7 in Softly, As In A Morning Sunrise.
	//
	// A cadence across the loop, a turnaround toward a first chord in m7,
	// makes it a tonic only when asked (across): as the chord heard, the
	// first chord of Fly Me To The Moon or All The Things You Are is a VI
	// or a II, not a tonic; as where the chart ends, the Cm7 that the last
	// G7 of Sugar goes back to is its tonic.
	tonicAt := func(i int, across bool) []harmony.Tonality {
		ch := c.Chords[i]
		if t := tonicOf(ch); t != nil {
			return t
		}
		n := target[i]
		if n < 0 || !across && blocks[n].Five > i || member[i] >= 0 ||
			ch.Chord.Pattern.Tetrad() != harmony.ChordMinorSeventh || ch.Inverted() {
			return nil
		}
		if a := blocks[n].Announced; len(a) > 0 && a[0].Tonic() == ch.Chord.Root && ModesOf(a)&Minor != 0 {
			return MinorTonalities(ch.Chord.Root)
		}
		return nil
	}

	// pass hears the chart once. On the first hearing, a cadence across
	// the loop has not sounded yet when its target, the first chord, does.
	pass := func(s Sensed, first bool) []Sensed {
		out := make([]Sensed, len(c.Chords))
		for i, ch := range c.Chords {
			s.Awaited, s.Tonic, s.Resolves, s.Across = nil, nil, nil, false
			if ch.Silent {
				out[i] = s
				continue
			}
			s.Tonic = tonicAt(i, false)
			if s.Ground == nil && i == 0 {
				s.Ground = s.Tonic
			}
			if t := s.Tonic; t != nil && s.Start != nil && s.Ground != nil &&
				sameTonic(t, s.Start) && !sameTonic(s.Ground, s.Start) {
				// Home is recognised on sight: its tonic chord brings
				// the ground back, with or without a cadence (F/C C at
				// the end of My Way). Heard afterwards, home starts
				// with the cadence when there is one.
				s.Ground, s.Since, s.Local = s.Start, i, nil
				if n := target[i]; n >= 0 {
					s.Since = firstOf(blocks[n])
				}
			}
			n := target[i]
			heard := n >= 0 && len(blocks[n].Announced) > 0 && !(first && blocks[n].Five > i)
			if heard && blocks[n].Kind == harmony.PlagalApproach {
				// A plagal cadence concludes on a tonic already there, the
				// ground, home or the local tonic, and opens none.
				t := blocks[n].Announced
				heard = sameTonic(t, s.Ground) || sameTonic(t, s.Start) || sameTonic(t, s.Local)
			}
			if heard {
				t := blocks[n].Announced
				from := firstOf(blocks[n])
				switch {
				case s.Ground == nil:
					s.Ground, s.Since, s.Local = t, from, nil
				case t[0].Tonic() == s.Ground[0].Tonic():
					s.Local = nil
				case s.Start != nil && t[0].Tonic() == s.Start[0].Tonic():
					s.Ground, s.Since, s.Local = s.Start, from, nil
				case s.Local != nil && t[0].Tonic() == s.Local[0].Tonic():
					cadences++
					if ti := s.Tonic; ti != nil {
						tonic = ti
					}
				case s.Local != nil && s.Tonic == nil && holds(s.Local, ch.Chord.Set()):
					// A degree of the local tonic tonicised in passing, as
					// Dm7 by A7 while C major is local: the local goes on.
				default:
					s.Local, since, cadences, stable = t, from, 1, 0
					tonic = s.Tonic
				}
			} else if s.Local != nil && !holds(s.Local, ch.Chord.Set()) && !prepares(c, blocks, member[i], s.Local) {
				s.Local = nil
			}
			if s.Local != nil && member[i] < 0 && holds(s.Local, ch.Chord.Set()) {
				stable += ch.Length
			}
			if s.Local != nil && tonic != nil && (stable > bar || cadences >= 2) {
				installed(&s)
			}
			if s.Start == nil {
				s.Start = s.Ground
			}
			if n := member[i]; n >= 0 && !borrowsTonic(c, blocks[n], s.Ground) {
				s.Awaited = blocks[n].Announced
			}
			if n := member[i]; n >= 0 && blocks[n].Five == i && blocks[n].Target >= 0 {
				s.Resolves, s.Across = tonicAt(blocks[n].Target, true), blocks[n].Target < i
			}
			out[i] = s
		}
		return out
	}

	// A chart loops: the chorus is heard again after itself, and that is
	// the hearing to render. The ground and the local tonic go on from
	// the end of the first chorus (or from the bar before the coda), so
	// that a turnaround prepares the first chord, and home is where the
	// tune ends: its last tonic, the coda's when it has one, installed
	// by being the last.
	seed := Sensed{}
	if t, ok := Blues(c); ok {
		seed.Start, seed.Ground = t, t
	}
	out := pass(seed, true)
	if !c.Loops {
		return out
	}
	end := len(out) - 1
	if c.Coda > 0 {
		end = c.Coda - 1
	}
	again := out[end]
	if home := Tune(c, out); home != nil {
		again.Start, again.Ground, again.Local = home, home, nil
	}
	return pass(again, false)
}

// firstOf returns the first chord of a block: its two, its suspension
// or its V.
func firstOf(b Block) int {
	switch {
	case b.Two >= 0:
		return b.Two
	case b.Sus >= 0:
		return b.Sus
	}
	return b.Five
}

// prepares reports whether block n, the one a change belongs to,
// resolves on a chord that holds in the tonalities: A7 before Dm7 while
// C major is local keeps it.
func prepares(c Changes, blocks []Block, n int, ts []harmony.Tonality) bool {
	if n < 0 || blocks[n].Target < 0 {
		return false
	}
	return holds(ts, c.Chords[blocks[n].Target].Chord.Set())
}

// Grounds reads the ground of each change afterwards: a ground the ear
// installs at some chord holds, heard back, from the first chord of the
// cadence that led to it. In Tune Up, C major is installed on the
// second bar of Cmaj7, and holds from Dm7 on. The chords heard before
// any ground, a tune starting on its two, belong to the first one.
func Grounds(c Changes, sensed []Sensed) [][]harmony.Tonality {
	out := make([][]harmony.Tonality, len(sensed))
	for i, s := range sensed {
		if s.Ground != nil && (i == 0 || sensed[i-1].Ground == nil) {
			for k := range i {
				out[k] = s.Ground
			}
		}
		out[i] = s.Ground
		if i == 0 || sameTonic(s.Ground, sensed[i-1].Ground) {
			continue
		}
		// A cadence across the loop starts at the end of the chorus.
		for k := s.Since; k != i; k = (k + 1) % len(out) {
			out[k] = s.Ground
		}
	}
	return out
}

func sameTonic(a, b []harmony.Tonality) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a[0].Tonic() == b[0].Tonic()
}

// Tune reads the tonality of the tune from its end, in the major or the
// three minors as its tonic's third says. Walking back from the last
// chord, the first of these decides:
//
//   - a chord heard as the tonic of the ground: F/C C at the end of My
//     Way, a plagal amen on C; Cmaj7 before the turnaround of Fly Me To
//     The Moon;
//   - the V of a cadence that resolves on a tonic: E♭7 A♭maj7 at the
//     end of Along Came Betty, before the turnaround Bm7 E7;
//   - across the loop, a cadence that goes back to the ground: G7 at
//     the end of Sugar goes back to Cm7 and makes it a tune in C minor,
//     not in the F minor that G♭7 Fm7 tonicises in bar 11. A turnaround
//     toward a first chord that is not the ground (E7 Am7 in Fly Me To
//     The Moon) does not count.
//
// A tierce picarde does not make a minor tune major: when the tonic
// chord that decides is major but the tune spends longer on its minor
// tonic than on its major one, the tune is in the minor (see
// [Picardy]).
//
// A blues is in its own tonic, and with none of these, the tune is in
// the ground at the end.
func Tune(c Changes, sensed []Sensed) []harmony.Tonality {
	t, _ := tune(c, sensed)
	return t
}

// Picardy reports whether a minor tune ends on its tonic made major,
// the tierce picarde of church music: a major chord rings with fewer
// clashing partials under the resonance of a great organ.
func Picardy(c Changes, sensed []Sensed) bool {
	_, p := tune(c, sensed)
	return p
}

func tune(c Changes, sensed []Sensed) ([]harmony.Tonality, bool) {
	if t, ok := Blues(c); ok {
		return t, false
	}
	// picardy keeps the minor under a major tonic chord, when the tune
	// heard its minor tonic longer than its major one.
	picardy := func(t []harmony.Tonality) ([]harmony.Tonality, bool) {
		if ModesOf(t) != Major {
			return t, false
		}
		var major, minor Ticks
		for i, s := range sensed {
			if s.Tonic == nil || !sameTonic(s.Tonic, t) {
				continue
			}
			if ModesOf(s.Tonic) == Major {
				major += c.Chords[i].Length
			} else {
				minor += c.Chords[i].Length
			}
		}
		if minor > major {
			return MinorTonalities(t[0].Tonic()), true
		}
		return t, false
	}
	for i := len(sensed) - 1; i >= 0; i-- {
		s := sensed[i]
		switch {
		case c.Chords[i].Silent:
		case s.Tonic != nil && sameTonic(s.Tonic, s.Ground):
			return picardy(s.Tonic)
		case s.Resolves != nil && (!s.Across || sameTonic(s.Resolves, s.Ground)):
			return picardy(s.Resolves)
		}
	}
	for i := len(sensed) - 1; i >= 0; i-- {
		if sensed[i].Ground != nil {
			return sensed[i].Ground, false
		}
	}
	return nil, false
}

// tonicOf returns the tonalities a change is the tonic of, nil when it
// cannot be one: a chord that can be a tonic (a triad, maj7, 6, m6,
// m(maj7), not m7 nor a seventh), in root position or with its third in
// the bass. With its fifth in the bass, it is a six-four over a pedal:
// F/C C in My Way is C, the F a neighbour over the bass.
func tonicOf(change Change) []harmony.Tonality {
	ch := change.Chord
	if change.Inverted() && (int(change.Bass)-int(ch.Root)+12)%12 == 7 {
		return nil
	}
	switch ch.Pattern.Tetrad() {
	case harmony.ChordMajorTriad, harmony.ChordMajorSeventh, harmony.ChordMajorSixth:
		return MajorTonalities(ch.Root)
	case harmony.ChordMinorTriad, harmony.ChordMinorSixth, harmony.ChordMinorMajorSeventh:
		return MinorTonalities(ch.Root)
	}
	return nil
}

// borrowsTonic reports whether a block is a two five that does not
// resolve and whose two sits on the ground's tonic: its two is then
// the tonic borrowed from another scale, not a two.
func borrowsTonic(c Changes, b Block, ground []harmony.Tonality) bool {
	return b.Two >= 0 && b.Target < 0 && ground != nil &&
		c.Chords[b.Two].Chord.Root == ground[0].Tonic()
}
