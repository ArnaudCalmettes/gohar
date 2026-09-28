package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// A Phrase runs from one rest to the next: from its first change to
// the one it concludes on, the tonic where the music comes to rest.
//
// The phrases tell where a tune sets out from (see [Home]) and where it
// stops (see [Tune]): the first of them, and the last.
type Phrase struct {
	From, To int // its first and last changes

	// Tonic is the tonic it concludes on, nil for the tail after the
	// last conclusion: the turnaround back to the first chord.
	Tonic []harmony.Tonality

	// Rest is, for the first phrase, the tonic it opens at rest on: a
	// tonic chord held longer than a bar, that its first cadence comes
	// back to. In a Sentimental Mood holds Dm for two bars and comes
	// back to it by A7; Blue Skies opens on the same line from Am, but
	// its first cadence goes to C6, and Just Friends opens on Cmaj7,
	// its IV, before Cm7 F7 Gmaj7.
	Rest []harmony.Tonality

	// Stops tells the last phrase, the one that concludes where the tune
	// stops rather than by coming to rest.
	Stops bool
}

// Phrases cuts a sequence into phrases. A phrase has no set length: it
// ends when it comes to rest, or when the tune stops.
//
// # Coming to rest
//
// A phrase comes to rest on a tonic chord that a cadence leads to,
// perfect or plagal, and that either
//
//   - brings back the tonic the tune opened on: How Insensitive is a
//     long sigh from Dm down to Dm again, fourteen bars later, and Fly
//     Me To The Moon comes back to Am7 at bar 8;
//   - or holds longer than a bar and longer than the chords that lead
//     to it: Gm6 in Autumn Leaves, two bars after Am7♭5 D7.
//
// A tonic passed through does not end the phrase (B♭maj7 in Autumn
// Leaves), nor the IV of the tonic before, however long: E♭maj7 two
// bars after B♭6 Fm7 B♭7 in Cherokee. Past the opening chord, a phrase
// rests on a tonic chord, not on a m7, almost always a subdominant:
// Cm7, held two bars in There Will Never Be Another You, is its VI. The
// rest lasts as long as the tonic holds.
//
// # Where the tune stops
//
// The last phrase concludes on the last tonic heard. A standard is
// played to it and no further: the turnaround after it goes back to
// the first chord, which it always points to, and says nothing of the
// tonality. Lullaby Of Birdland stops on A♭maj7 before Gm7♭5 C7 goes
// back to Fm. A tonic already installed, the one the tune opened on,
// the one the phrase before came to rest on or, with none, the one its
// first cadence went to, needs no cadence to come back: 'Round
// Midnight stops on E♭6, Chega De Saudade on D6. Another one comes after
// a two five one or a plagal cadence, not a lone V: Yesterdays goes
// through B♭maj7 in its cycle of dominants, and stops on Dm. And no
// tune stops on the IV of the tonic before: Virgo goes through B♭maj7,
// its IV, before Gm7 C7 goes back to F.
func Phrases(c Changes, blocks []Block) []Phrase {
	if opening(c) < 0 {
		return nil
	}
	ph := &phrasing{c: c, blocks: blocks, r: rolesOf(c, blocks), bar: barOf(c)}
	ph.opens = opensOn(c, blocks, ph.r)
	ph.home = ph.opens
	ph.rests()
	ph.stops()
	return ph.out
}

// A phrasing cuts changes into phrases.
type phrasing struct {
	c      Changes
	blocks []Block
	r      roles
	bar    Ticks
	opens  []harmony.Tonality // the tonic the tune opens on, if any

	out    []Phrase
	p      Phrase             // the phrase being cut
	before []harmony.Tonality // the tonic the phrase before came to rest on
	first  []harmony.Tonality // the tonic the first cadence resolves on
	home   []harmony.Tonality // the tonic opened on, else the first rested on
}

// cadenced returns the tonic a cadence within the chorus leads change i
// to, and that cadence.
func (ph *phrasing) cadenced(i int) ([]harmony.Tonality, *Block) {
	n := ph.r.target[i]
	if n < 0 || ph.blocks[n].Five > i || len(ph.blocks[n].Announced) == 0 {
		return nil, nil
	}
	t := tonicAt(ph.c, ph.blocks, ph.r, i)
	if t == nil || !sameTonic(ph.blocks[n].Announced, t) {
		return nil, nil
	}
	return t, &ph.blocks[n]
}

// restsAt tells the tonic a phrase comes to rest on at change i, nil
// when it does not: see "Coming to rest".
func (ph *phrasing) restsAt(i int) []harmony.Tonality {
	t, b := ph.cadenced(i)
	switch {
	case t == nil:
		return nil
	case sameTonality(t, ph.opens):
		return t
	case tonicOf(ph.c.Chords[i]) == nil, ph.subdominant(t):
		return nil
	}
	h := held(ph.c, i)
	if h <= ph.bar {
		return nil
	}
	for _, k := range []int{b.Two, b.Sus, b.Five} {
		if k >= 0 && ph.c.Chords[k].Length >= h {
			return nil
		}
	}
	return t
}

// rests cuts the phrases that come to rest, from the opening on.
func (ph *phrasing) rests() {
	o := opening(ph.c)
	ph.p = Phrase{From: o}
	ph.before = tonicOf(ph.c.Chords[o])
	for i := o + 1; i < len(ph.c.Chords); i++ {
		if t, _ := ph.cadenced(i); ph.first == nil && t != nil {
			ph.first = t
			if rest := tonicOf(ph.c.Chords[o]); sameTonic(t, rest) && held(ph.c, o) > ph.bar {
				ph.p.Rest = rest
			}
		}
		if t := ph.restsAt(i); t != nil {
			i = heldTo(ph.c, i)
			ph.conclude(i, t, false)
		}
	}
}

// stopsAt tells the tonic the tune stops on at change i, nil when it
// does not: see "Where the tune stops". An installed tonic comes back
// in any position, the six-four included: B♭maj7/F before G7 Cm7 F7
// at the end of Someday My Prince Will Come. Any other bass makes
// another chord: Fmaj7/G at the end of Only Trust Your Heart is the V
// of C with its fourth, not F.
func (ph *phrasing) stopsAt(i int, installed []harmony.Tonality) []harmony.Tonality {
	root := ph.c.Chords[i]
	if above := (int(root.Bass) - int(root.Chord.Root) + 12) % 12; above == 7 {
		root.Bass = root.Chord.Root
	}
	if t := tonicOf(root); t != nil && (sameTonic(t, ph.opens) || sameTonic(t, installed)) {
		return t
	}
	if t := minorSeventhTonic(ph.c, ph.blocks, ph.r, i); t != nil && (sameTonic(t, ph.opens) || sameTonic(t, installed)) {
		return t
	}
	t, b := ph.cadenced(i)
	switch {
	case t == nil,
		tonicOf(ph.c.Chords[i]) == nil && !sameTonality(t, ph.opens),
		ph.subdominant(t),
		b.Two < 0 && b.Sus < 0 && b.Kind != harmony.PlagalApproach:
		return nil
	}
	return t
}

// stops cuts the last phrase, from the last rest to where the tune
// stops, and the tail after it.
func (ph *phrasing) stops() {
	if ph.p.From >= len(ph.c.Chords) {
		return
	}
	installed := ph.before
	if installed == nil {
		installed = ph.first
	}
	for i := len(ph.c.Chords) - 1; i >= ph.p.From; i-- {
		if t := ph.stopsAt(i, installed); t != nil {
			ph.conclude(i, t, true)
			break
		}
	}
	if ph.p.From < len(ph.c.Chords) {
		ph.p.To = len(ph.c.Chords) - 1
		ph.out = append(ph.out, ph.p)
	}
}

// conclude ends the phrase being cut at change i, on tonic t.
func (ph *phrasing) conclude(i int, t []harmony.Tonality, stops bool) {
	ph.p.To, ph.p.Tonic, ph.p.Stops = i, t, stops
	ph.out = append(ph.out, ph.p)
	ph.p, ph.before = Phrase{From: i + 1}, t
	if ph.home == nil {
		ph.home = t
	}
}

// subdominant reports whether a tonic is the IV of the tonic the phrase
// before came to rest on: a phrase neither rests nor stops there, unless
// it is home. E♭maj7 two bars after B♭6 Fm7 B♭7 in Cherokee; B♭maj7
// before Gm7 C7 goes back to F in Virgo. Dm after a rest on Am in Chega
// De Saudade is its IV, and the tonic it set out from.
func (ph *phrasing) subdominant(t []harmony.Tonality) bool {
	return ph.before != nil && t[0].Tonic() == ph.before[0].Tonic().Transpose(5) && !sameTonic(t, ph.home)
}

// Home reads where the tune sets out from: the tonic it opens at rest
// on, else the one its first phrase comes to rest on, nil when none
// does before the tune stops. Autumn Leaves sets out from G minor, How
// Insensitive from D minor, Just Friends from G, where its first phrase
// rests after opening on Cmaj7, its IV.
func Home(c Changes, phrases []Phrase) []harmony.Tonality {
	if t, ok := Blues(c); ok {
		return t
	}
	if len(phrases) == 0 {
		return nil
	}
	if rest := phrases[0].Rest; rest != nil {
		return rest
	}
	for _, p := range phrases {
		if p.Tonic != nil && !p.Stops {
			return p.Tonic
		}
	}
	return nil
}

// Tune reads the tonality of the tune: where its last phrase stops (see
// [Phrases]), in the major or the three minors as its tonic's third
// says. All The Things You Are is only settled by its last A♭maj7.
//
// A tune that opens at rest has installed its home before leaving it,
// and is in it wherever it stops, unless it stops on the same tonic: In
// a Sentimental Mood is in D minor though it concludes on Gm7 C7♭9
// Fmaj7. Chega De Saudade holds D minor for its first half and D major
// for its second, stops on D6, and reads as D major here: the tune is
// in D, as much minor as major, and the reading only says where it
// stops.
//
// A picardy third does not make a minor tune major (see [Picardy]). A
// blues is in its own tonic, found by its form. A tune where nothing
// concludes is in its home.
func Tune(c Changes, phrases []Phrase) []harmony.Tonality {
	if t, ok := Blues(c); ok {
		return t
	}
	home := Home(c, phrases)
	_, end := concluding(phrases)
	switch {
	case end == nil,
		phrases[0].Rest != nil && !sameTonic(end.Tonic, home):
		return home
	case Picardy(c, phrases):
		return MinorTonalities(end.Tonic[0].Tonic())
	}
	return end.Tonic
}

// Picardy reports whether a minor tune ends on its tonic made major,
// the picardy third of church music: a major chord rings with fewer
// clashing partials under the resonance of a great organ.
//
// The tonic is minor where it is first heard and where it is last heard
// before the end, and major on the last chord. 'Round Midnight closes
// its second A on E♭6, then starts its last A again on E♭m: a picardy
// third. Chega De Saudade holds D major for its whole second half: not
// a picardy third, a tune as much major as minor.
func Picardy(c Changes, phrases []Phrase) bool {
	if _, ok := Blues(c); ok {
		return false
	}
	_, end := concluding(phrases)
	if end == nil || ModesOf(end.Tonic) != Major {
		return false
	}
	var first, last []harmony.Tonality
	for i := range c.Chords[:end.To] {
		t := tonicOf(c.Chords[i])
		if !sameTonic(t, end.Tonic) || heldTo(c, i) >= end.To {
			continue // not the tonic, or the final one held
		}
		if first == nil {
			first = t
		}
		last = t
	}
	return first != nil && ModesOf(first) != Major && ModesOf(last) != Major
}

// concluding returns the first and the last phrases that conclude, nil
// when none does.
func concluding(phrases []Phrase) (first, last *Phrase) {
	for i := range phrases {
		if phrases[i].Tonic == nil {
			continue
		}
		if first == nil {
			first = &phrases[i]
		}
		last = &phrases[i]
	}
	return first, last
}
