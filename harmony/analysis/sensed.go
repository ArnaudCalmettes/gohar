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

	// Home is where the tune sets out from, once its first phrase has
	// landed: see [Home].
	Home []harmony.Tonality

	// Lands tells that the music comes to rest here, on Tonic: a phrase
	// ends. The first phrase to land gives the tune its home (see
	// [Tune]).
	Lands bool
}

// Sense reads the sensed tonic at each change.
//
// # What installs a tonic
//
// Home is where the first phrase comes to rest. A phrase has no set
// length: it goes on until it lands on a tonic chord that a cadence
// leads to, and that either
//
//   - brings back the chord the tune opened on, when that one could be
//     a tonic: How Insensitive is a long sigh from Dm down to Dm again,
//     fourteen bars later, and Fly Me To The Moon comes back to Am7 at
//     bar 8;
//   - or holds longer than a bar and longer than the chords that lead
//     to it: Gm6 in Autumn Leaves, two bars after Am7♭5 D7.
//
// A tonic passed through does not end the phrase (B♭maj7 in Autumn
// Leaves), nor does a subdominant, however long: E♭maj7 two bars after
// B♭6 Fm7 B♭7 in Cherokee is its IV. Past the opening chord, a phrase
// rests on a tonic chord, not on a m7, almost always a subdominant:
// Cm7, held two bars in There Will Never Be Another You, is its VI. And
// a m7 is not led to its rest by a plagal cadence: C7 Gm7 in
// Honeysuckle Rose goes back and forth between a two and its five.
//
// A tune can also open at rest, on a tonic chord held longer than a bar
// (E♭maj7 in There Will Never Be Another You), unless the first phrase
// then lands a fifth above it: Just Friends opens on Cmaj7, its IV, and
// lands on Gmaj7 after Cm7 F7, a plagal cadence stretched out.
//
// Until the first phrase lands, the ground is a guess: the first chord
// when it can be a tonic chord (a major or minor triad, maj7, 6, m6,
// m(maj7), but not m7), and else the first cadence that resolves.
// Landing makes it home. A cadence that resolves elsewhere than on the
// ground makes its target a local tonic.
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

	// opening is the first chord that sounds, the one a phrase can come
	// back to when it could be a tonic.
	opening := opening(c)
	openingTonic := false
	if opening >= 0 {
		o := c.Chords[opening]
		openingTonic = tonicOf(o) != nil
		if o.Chord.Pattern.Tetrad() == harmony.ChordMinorSeventh {
			// A m7 opens on a tonic unless it is a two, whether or not a
			// block claims it: Gm7 C7 in Honeysuckle Rose.
			n := c.Next(opening)
			openingTonic = member[opening] < 0 && !(n >= 0 && c.Chords[n].Chord.Root == o.Chord.Root.Transpose(5) &&
				c.Chords[n].Chord.Pattern.Tetrad() == harmony.ChordDominantSeventh)
		}
	}

	held := func(i int) Ticks { return held(c, i) }

	// landsAt tells what tonic the music comes to rest on at change i,
	// nil when a phrase does not end there: see "What installs a tonic".
	landsAt := func(i int, ground []harmony.Tonality) []harmony.Tonality {
		n := target[i]
		if n < 0 || blocks[n].Five > i || len(blocks[n].Announced) == 0 {
			return nil
		}
		t := tonicAt(i, false)
		if t == nil || blocks[n].Announced[0].Tonic() != t[0].Tonic() ||
			tonicOf(c.Chords[i]) == nil && blocks[n].Kind == harmony.PlagalApproach {
			return nil
		}
		if o := c.Chords[opening].Chord; openingTonic && o.Root == t[0].Tonic() &&
			isMinor(o.Pattern) == (ModesOf(t) == Minor) {
			return t
		}
		if tonicOf(c.Chords[i]) == nil ||
			ground != nil && t[0].Tonic() == ground[0].Tonic().Transpose(5) {
			return nil
		}
		h := held(i)
		if h <= bar {
			return nil
		}
		for _, k := range []int{blocks[n].Two, blocks[n].Sus, blocks[n].Five} {
			if k >= 0 && c.Chords[k].Length >= h {
				return nil
			}
		}
		return t
	}

	// pass hears the chart once. On the first hearing, a cadence across
	// the loop has not sounded yet when its target, the first chord, does.
	pass := func(s Sensed, first bool) []Sensed {
		out := make([]Sensed, len(c.Chords))
		landed := !first
		atRest := false // the tune opened at rest, and has not landed since
		for i, ch := range c.Chords {
			s.Awaited, s.Tonic, s.Resolves, s.Across, s.Lands = nil, nil, nil, false, false
			if ch.Silent {
				out[i] = s
				continue
			}
			before := s.Ground
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
			if t := landsAt(i, before); t != nil {
				s.Lands = true
				wasAtRest := atRest
				atRest = false
				if !landed || wasAtRest && t[0].Tonic() == s.Start[0].Tonic().Transpose(7) {
					// The first phrase lands: the guess gives way to home.
					landed = true
					s.Start, s.Home, s.Local = t, t, nil
					if !sameTonic(t, s.Ground) {
						s.Ground, s.Since = t, firstOf(blocks[n])
					}
				}
			} else if i == opening && !landed && s.Tonic != nil && tonicOf(ch) != nil && held(i) > bar {
				// The tune opens at rest.
				s.Lands, landed, atRest = true, true, true
				s.Start, s.Ground, s.Home = s.Tonic, s.Tonic, s.Tonic
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
	// that a turnaround prepares the first chord. Heard a second time,
	// the tune is known: its start tonic is where it stops (see [Tune]).
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

// Tune reads the tonality of the tune, in the major or the three
// minors as its tonic's third says: the last tonic heard, where the
// last phrase stops. A standard is played to that tonic and no
// further: the turnaround after it goes back to the first chord, which
// it always points to, and says nothing of the tonality. Lullaby Of
// Birdland stops on A♭maj7 before Gm7♭5 C7 goes back to Fm; All The
// Things You Are is only settled by its last A♭maj7. Walking back from
// the last chord, the first of these decides:
//
//   - a chord heard as the tonic of the ground: F/C C at the end of My
//     Way, a plagal amen on C; Cmaj7 before the turnaround of Fly Me To
//     The Moon;
//   - a two five one within the chorus: B♭m7 E♭7 A♭maj7 at the end of
//     Lullaby Of Birdland or Along Came Betty, before their
//     turnarounds. A lone V does not stop a tune: Yesterdays goes
//     through B♭maj7 in its cycle of dominants, and stops in D minor.
//     Nor does a m7, almost always a subdominant, unless it is the
//     ground's tonic: F♯7 Fm7 near the end of Sugar tonicises its IV,
//     and Sugar stops in C minor. A V does, on the last chord of the
//     chart: B♭7 E♭7 A♭6 at the end of Sweet Georgia Brown.
//
// Some charts do not write the last tonic, and stop on the IV of home
// before a turnaround back to it that is the last cadence:
// Unforgettable, from G, stops on Cmaj7 before Am7 D7 Gmaj7. The IV is
// then skipped.
//
// A tune that opens at rest, on its tonic held longer than a bar, and
// whose first cadence comes back to it, has installed its home before
// leaving it, and is in it wherever it stops: In a Sentimental Mood
// holds Dm for two bars, returns to it by A7, and is in D minor though
// it concludes on Gm7 C7♭9 Fmaj7, in its relative major. Blue Skies
// opens on the same line from Am, but its first cadence goes to C6:
// it is in C. A tune that ends on the same tonic keeps the mode it
// ends in: Chega De Saudade holds D minor for its first half, D major
// for its second, stops on D6, and is in D major.
// Elsewhere, where the tune sets out from (see [Home]) is not always
// where it stops: Lullaby Of Birdland sets out from F minor, on a
// chord held half a bar. A picardy third does not make a minor tune
// major: the tonic that decides may turn major on the final chord, and
// the tune stays minor (see [Picardy]).
//
// A blues is in its own tonic, found by its form. A tune with none of
// these is in its home, and else in the ground at the end.
func Tune(c Changes, sensed []Sensed) []harmony.Tonality {
	if t, ok := Blues(c); ok {
		return t
	}
	t, _ := ending(c, sensed)
	if h := Home(c, sensed); opensAtRest(c, h) && sameTonic(firstCadence(c, sensed), h) && !sameTonic(t, h) {
		return h
	}
	return t
}

// Home reads where the tune sets out from: the tonic its first phrase
// comes to rest on (see [Sense]), nil when no phrase lands. Autumn
// Leaves sets out from G minor, How Insensitive from D minor.
func Home(c Changes, sensed []Sensed) []harmony.Tonality {
	if t, ok := Blues(c); ok {
		return t
	}
	if len(sensed) == 0 {
		return nil
	}
	return sensed[len(sensed)-1].Home
}

// hasTwo reports whether the V at change i comes after its two, a m7
// or m7♭5 a fifth above it.
func hasTwo(c Changes, i int) bool {
	if i == 0 {
		return false
	}
	p := c.Chords[i-1]
	t := p.Chord.Pattern.Tetrad()
	return !p.Silent && p.Chord.Root == c.Chords[i].Chord.Root.Transpose(7) &&
		(t == harmony.ChordMinorSeventh || t == harmony.ChordHalfDiminished)
}

// held returns how long a tonic holds from change i: the chord and the
// chords on the same root and third after it, tonic chords or, for a
// minor tonic, its m7 (Gm6 | Gm6, the line of Dm Dm(maj7) Dm7 Dm6 in
// In a Sentimental Mood; not Gm6 | G7, nor E♭maj7 | E♭m6 in Candy, the
// IV turning minor).
func held(c Changes, i int) Ticks {
	var d Ticks
	first := c.Chords[i].Chord
	for k := i; k < len(c.Chords); k++ {
		ch := c.Chords[k]
		keeps := ch.Chord.Pattern == first.Pattern ||
			isMinor(ch.Chord.Pattern) == isMinor(first.Pattern) && (tonicOf(ch) != nil ||
				isMinor(first.Pattern) && ch.Chord.Pattern.Tetrad() == harmony.ChordMinorSeventh)
		if ch.Silent || ch.Chord.Root != first.Root || k > i && !keeps {
			break
		}
		d += ch.Length
	}
	return d
}

// firstCadence returns the tonic the first cadence within the chorus
// resolves on, nil when none does.
func firstCadence(c Changes, sensed []Sensed) []harmony.Tonality {
	for i, s := range sensed {
		if s.Resolves != nil && !s.Across && tonicOf(c.Chords[c.Next(i)]) != nil {
			return s.Resolves
		}
	}
	return nil
}

// opensAtRest reports whether the tune opens on a tonic chord held
// longer than a bar, the home it sets out from.
func opensAtRest(c Changes, home []harmony.Tonality) bool {
	o := opening(c)
	if o < 0 || home == nil {
		return false
	}
	bar := 4 * TicksPerBeat
	if len(c.Bars) > 1 {
		bar = c.Bars[1] - c.Bars[0]
	}
	t := tonicOf(c.Chords[o])
	return t != nil && sameTonic(t, home) && held(c, o) > bar
}

// opening returns the first change that sounds, -1 when none does.
func opening(c Changes) int {
	for i, ch := range c.Chords {
		if !ch.Silent {
			return i
		}
	}
	return -1
}

// Picardy reports whether a minor tune ends on its tonic made major,
// the picardy third of church music: a major chord rings with fewer
// clashing partials under the resonance of a great organ. It is on the
// last chord: a major tonic heard before is a modulation to it, unless
// a section closed there as the tune does and the minor came back
// ('Round Midnight).
func Picardy(c Changes, sensed []Sensed) bool {
	if _, ok := Blues(c); ok {
		return false
	}
	_, p := ending(c, sensed)
	return p
}

// ending returns the tonality the tune ends in (see [Tune]) and whether
// a major tonic there is a picardy third.
func ending(c Changes, sensed []Sensed) ([]harmony.Tonality, bool) {
	// picardy keeps the minor under a major tonic at change at, the last
	// tonic heard, when it is a picardy third: the tonic, minor so
	// far and only minor, turns major on the final chord, and the tune
	// was in the minor when it got there. A section may
	// close with the same ending before: 'Round Midnight ends its second
	// A as it ends, on E♭6, and its last A starts again on E♭m. Any other
	// major tonic, or the minor not coming back, is a modulation: the
	// second half of Chega De Saudade, in D major.
	picardy := func(t []harmony.Tonality, at int) ([]harmony.Tonality, bool) {
		if ModesOf(t) != Major || at == 0 || ModesOf(sensed[at-1].Ground) != Minor ||
			!sameTonic(sensed[at-1].Ground, t) {
			return t, false
		}
		last := at
		for last > 0 && c.Chords[last-1].Chord.Root == t[0].Tonic() && !c.Chords[last-1].Silent {
			last--
		}
		// same tells whether change i closes as the tune does: the same
		// chord before the same tonic.
		same := func(i int) bool {
			return i > 0 && last > 0 && c.Chords[i].Chord == c.Chords[last].Chord &&
				c.Chords[i-1].Chord == c.Chords[last-1].Chord
		}
		minor, closed := false, false // closed: a section just closed in the major
		for i, s := range sensed[:last] {
			if s.Tonic == nil || !sameTonic(s.Tonic, t) || c.Chords[i].Silent {
				continue
			}
			if ModesOf(s.Tonic) != Major {
				minor, closed = true, false
			} else if closed || !same(i) {
				return t, false
			} else if c.Chords[i-1].Chord.Root != t[0].Tonic() {
				closed = true
			}
		}
		if minor && !closed {
			return MinorTonalities(t[0].Tonic()), true
		}
		return t, false
	}
	home := Home(c, sensed)
	// back tells whether the last cadence goes back to home across the
	// loop, and subdominant a tonic that is then the IV of home: the
	// chart stops there, and leaves the last cadence to the first chord.
	back := false
	for i := len(sensed) - 1; i >= 0; i-- {
		if s := sensed[i]; s.Resolves != nil {
			back = s.Across && home != nil && sameTonic(s.Resolves, home)
			break
		}
	}
	subdominant := func(t []harmony.Tonality) bool {
		return back && t[0].Tonic() == home[0].Tonic().Transpose(5)
	}
	for i := len(sensed) - 1; i >= 0; i-- {
		s := sensed[i]
		switch {
		case c.Chords[i].Silent:
		case s.Tonic != nil && sameTonic(s.Tonic, s.Ground) && !subdominant(s.Tonic):
			return picardy(s.Tonic, i)
		case s.Resolves != nil && !s.Across && !subdominant(s.Resolves) &&
			tonicOf(c.Chords[c.Next(i)]) != nil && (hasTwo(c, i) || c.Next(i) == len(c.Chords)-1):
			return picardy(s.Resolves, c.Next(i))
		}
	}
	if h := home; h != nil {
		return h, false
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

// isMinor reports whether a chord has a minor third and no major one.
func isMinor(p harmony.ChordPattern) bool {
	return p.HasOffset(3) && !p.HasOffset(4)
}
