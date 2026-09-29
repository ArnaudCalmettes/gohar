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
	// Start is the tonic the tune set out from, kept apart from the
	// ground to hear the return home after a bridge that modulated.
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
	// one there: see [tonicAt].
	Tonic []harmony.Tonality

	// Resolves is, on the V of a cadence that resolves on a tonic, the
	// tonalities of that tonic; Across tells that it lands across the
	// loop, on the first chord of the chart.
	Resolves []harmony.Tonality
	Across   bool
}

// Sense reads the sensed tonic at each change, given the blocks and the
// phrases of the sequence. Four mechanisms, each a method of [hearing]:
//
//   - home: until the first phrase concludes, the ground is a guess, the
//     first chord when it can be a tonic, else the first cadence that
//     resolves; the first phrase then gives the tune its home (see
//     [Home]). A blues has its tonic for ground from the start.
//   - cadence: a cadence that resolves elsewhere than on the ground makes
//     its target a local tonic, which lasts while the chords after it
//     hold in it or prepare a chord that does.
//   - modulation: the local tonic becomes the ground.
//   - return: coming home is easier than leaving.
//
// A key signature plays no part in what the ear hears: it belongs to
// the written score, the weakest clue to a tonality. The fiches of En
// Harmonie read the same without it.
//
// Each change sees the state left by those before it and the blocks,
// whose tonalities a target decides between: one change ahead, as the
// analysis of a performance allows.
//
// # Two hearings
//
// A chart loops: the chorus is heard again after itself, and that is
// the hearing to render. The second hearing goes on from the end of the
// first (or from the bar before the coda), so that a turnaround
// prepares the first chord, and knows the tune: its start tonic is the
// tonality of the tune. On the first hearing, a cadence across the
// loop has not sounded yet when its target, the first chord, does.
//
// # The tonality of the tune
//
// Which tonality the tune is in is the caller's to say, as tune: the
// one the analysis concludes ([Tune]), or one it is told, the key the
// chart declares or one the reader forces. The chords alone cannot
// always tell (see docs/grilles.md, "Le plafond des grilles seules"):
// In a Sentimental Mood and Lullaby Of Birdland give the same evidence,
// and the melody puts one in D minor, the other in A flat. Nil leaves
// the second hearing to go on as the first ended.
func Sense(c Changes, blocks []Block, phrases []Phrase, tune []harmony.Tonality) []Sensed {
	h := newHearing(c, blocks, phrases)
	seed := Sensed{}
	if t, ok := Blues(c); ok {
		seed.Start, seed.Ground = t, t
	}
	out := h.pass(seed, true)
	if !c.Loops {
		return out
	}
	end := len(out) - 1
	if c.Coda > 0 {
		end = c.Coda - 1
	}
	again := out[end]
	if tune != nil {
		again.Start, again.Ground, again.Local = tune, tune, nil
	}
	return h.pass(again, false)
}

// A hearing is an ear going through the changes, with what it has
// gathered so far.
type hearing struct {
	c      Changes
	blocks []Block
	r      roles
	bar    Ticks
	first  bool // the first hearing
	s      Sensed

	// The first phrase, which may open at rest, the first one that
	// concludes before the tune stops, nil when none does, and the home
	// they give (see [Home]).
	opens     Phrase
	concludes *Phrase
	home      []harmony.Tonality

	// What the local tonic has gathered toward becoming the ground.
	local struct {
		tonic    []harmony.Tonality // the tonalities its I is the tonic of
		since    int                // the first chord of its first cadence
		cadences int                // cadences resolved on it
		stable   Ticks              // time of its stable chords
	}
}

func newHearing(c Changes, blocks []Block, phrases []Phrase) *hearing {
	h := &hearing{c: c, blocks: blocks, r: rolesOf(c, blocks), bar: barOf(c), opens: Phrase{From: -1}}
	h.home = Home(c, phrases)
	if len(phrases) > 0 {
		h.opens = phrases[0]
	}
	if first, _ := concluding(phrases); first != nil && !first.Stops {
		h.concludes = first
	}
	return h
}

// pass hears the changes once, from the state seed.
func (h *hearing) pass(seed Sensed, first bool) []Sensed {
	h.s, h.first = seed, first
	out := make([]Sensed, len(h.c.Chords))
	for i, ch := range h.c.Chords {
		h.s.Awaited, h.s.Tonic, h.s.Resolves, h.s.Across = nil, nil, nil, false
		if ch.Silent {
			out[i] = h.s
			continue
		}
		h.s.Tonic = tonicAt(h.c, h.blocks, h.r, i)
		if h.s.Ground == nil && i == 0 {
			h.s.Ground = h.s.Tonic
		}
		h.returns(i)
		h.cadence(i)
		h.modulation(i)
		if first {
			h.homeAt(i)
		}
		if h.s.Start == nil {
			h.s.Start = h.s.Ground
		}
		h.expect(i)
		out[i] = h.s
	}
	return out
}

// homeAt installs home where the first phrases give it: at the opening
// when the tune opens at rest, In a Sentimental Mood on Dm; else where
// the first phrase concludes, when it concludes at home: Just Friends,
// which opens on Cmaj7, its IV, on Gmaj7.
func (h *hearing) homeAt(i int) {
	switch {
	case i == h.opens.From && h.opens.Rest != nil:
		h.s.Start, h.s.Ground = h.opens.Rest, h.opens.Rest
	case h.concludes != nil && i == h.concludes.To && sameTonic(h.concludes.Tonic, h.home):
		h.s.Start, h.s.Local = h.home, nil
		if !sameTonic(h.s.Ground, h.home) {
			h.s.Ground, h.s.Since = h.home, h.since(i)
		}
	}
}

// returns brings the ground back home on sight: the tonic chord of the
// start tonic, with or without a cadence (F/C C at the end of My Way).
// Heard afterwards, home starts with the cadence when there is one.
func (h *hearing) returns(i int) {
	s := &h.s
	if t := s.Tonic; t == nil || s.Start == nil || s.Ground == nil ||
		!sameTonic(t, s.Start) || sameTonic(s.Ground, s.Start) {
		return
	}
	s.Ground, s.Since, s.Local = s.Start, h.since(i), nil
}

// since returns where a ground reached at change i holds from: the
// first chord of the cadence that led to it, else the change itself.
func (h *hearing) since(i int) int {
	if n := h.r.target[i]; n >= 0 {
		return firstOf(h.blocks[n])
	}
	return i
}

// cadence hears the cadence that resolves on change i, if any: it
// installs the first ground, brings home back, confirms the local
// tonic, or opens a new one. Without one, the local tonic ends when the
// chord neither holds in it nor prepares a chord that does.
//
// A plagal cadence (the IV, or the ♭VII7, before the tonic) concludes
// as a V-I does, and draws less: it confirms a tonic already there,
// the ground, home or the local tonic, and opens none. A cadence toward
// a degree of the local tonic that cannot be a tonic (A7 Dm7 while C
// major is local) does not interrupt it. Nor does a chord on the local
// tonic itself, whatever its colour: Gm7 Gm(maj7) Gm7 Gm(maj7) in My
// Lucky Star is G minor with a seventh that wanders, not a chord out of
// G melodic minor.
func (h *hearing) cadence(i int) {
	s, ch := &h.s, h.c.Chords[i]
	n := h.r.target[i]
	heard := n >= 0 && len(h.blocks[n].Announced) > 0 && !(h.first && h.blocks[n].Five > i)
	if heard && h.blocks[n].Kind == harmony.PlagalApproach {
		t := h.blocks[n].Announced
		heard = sameTonic(t, s.Ground) || sameTonic(t, s.Start) || sameTonic(t, s.Local)
	}
	if !heard {
		if s.Local != nil && !holds(s.Local, ch.Chord.Set()) && !onTonic(s.Local, ch) &&
			!prepares(h.c, h.blocks, h.r.member[i], s.Local) {
			s.Local = nil
		}
		return
	}
	t, from := h.blocks[n].Announced, firstOf(h.blocks[n])
	switch {
	case s.Ground == nil:
		s.Ground, s.Since, s.Local = t, from, nil
	case sameTonic(t, s.Ground):
		s.Local = nil
	case sameTonic(t, s.Start):
		s.Ground, s.Since, s.Local = s.Start, from, nil
	case sameTonic(t, s.Local):
		h.local.cadences++
		if s.Tonic != nil {
			h.local.tonic = s.Tonic
		}
	case s.Local != nil && s.Tonic == nil && holds(s.Local, ch.Chord.Set()):
		// A degree of the local tonic tonicised in passing.
	default:
		s.Local = t
		h.local.tonic, h.local.since, h.local.cadences, h.local.stable = s.Tonic, from, 1, 0
	}
}

// modulation makes the local tonic the ground when it holds beyond one
// bar of stable chords, or when a second cadence resolves on it while
// it lasts (Black Orpheus, C major). A stable chord sits on the local
// tonic, or holds in it and not in the ground, without preparing
// anything: B♭maj7 Gm7 in Tune Up, coming from D. A chord the ground
// owns says nothing for the local tonic: in Only Trust Your Heart, B7
// Em7 tonicises the III of C, and the Am7 after it is the VI of C, not
// the IV of E minor, before Dm7 G7 Cmaj7. Its I must be able to be a
// tonic: one does not modulate toward a subdominant, nor for one
// chord. A tonic held a single bar only tonicises: Along Came Betty
// dances from A flat to A and back, and never modulates.
//
// Nor does one modulate to the II of the ground, the subdominant
// degree: the II tonicised is the three six two five of the tonality,
// however long it holds. In My Lucky Star, Am7 D7 Gm7 then four bars
// of Gm7 Gm(maj7) are in F, before Dm7 G7 C7sus C7 goes back to F6.
func (h *hearing) modulation(i int) {
	s, ch := &h.s, h.c.Chords[i]
	if s.Local == nil {
		return
	}
	own := holds(s.Local, ch.Chord.Set()) && (s.Ground == nil || !holds(s.Ground, ch.Chord.Set()))
	if h.r.member[i] < 0 && (own || onTonic(s.Local, ch)) {
		h.local.stable += ch.Length
	}
	two := s.Ground != nil && h.local.tonic != nil && h.local.tonic[0].Tonic() == s.Ground[0].Tonic().Transpose(2)
	if h.local.tonic != nil && !two && (h.local.stable > h.bar || h.local.cadences >= 2) {
		s.Ground, s.Since, s.Local = h.local.tonic, h.local.since, nil
	}
}

// expect records what the chord makes the ear wait for: the tonality
// its block announces, and, on a V, the tonic it resolves on.
//
// A two five that does not resolve and whose two sits on the ground's
// tonic awaits nothing: its two is the tonic, borrowed from another
// scale. E♭m7 A♭7 in E flat is I then IV7, not a two five of D flat.
func (h *hearing) expect(i int) {
	n := h.r.member[i]
	if n < 0 {
		return
	}
	b := h.blocks[n]
	if !borrowsTonic(h.c, b, h.s.Ground) {
		h.s.Awaited = b.Announced
	}
	if b.Five == i && b.Target >= 0 {
		// Across the loop, the m7 the turnaround goes back to is where
		// the chart ends: the Cm7 after the last G7 of Sugar.
		t := tonicOf(h.c.Chords[b.Target])
		if t == nil {
			t = minorSeventhTonic(h.c, h.blocks, h.r, b.Target)
		}
		h.s.Resolves, h.s.Across = t, b.Target < i
	}
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

// borrowsTonic reports whether a block is a two five that does not
// resolve and whose two sits on the ground's tonic: its two is then
// the tonic borrowed from another scale, not a two.
func borrowsTonic(c Changes, b Block, ground []harmony.Tonality) bool {
	return b.Two >= 0 && b.Target < 0 && ground != nil &&
		c.Chords[b.Two].Chord.Root == ground[0].Tonic()
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
