package analysis

import (
	"slices"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Phrase runs from one rest to the next: from its first change to
// the one it concludes on, the tonic where the music comes to rest.
//
// The phrases tell where a tune sets out from (see [Home]) and where it
// stops (see [Tune]): the first of them, and the last.
type Phrase struct {
	From, To int // its first and last changes

	// Tonic is the tonic it concludes on, nil for the tail after the
	// last conclusion: the turnaround back to the first chord. Arrives
	// is the change where the tonic is reached, To being where it stops
	// holding; -1 for the tail. A tune that ends across the loop
	// arrives on its first change (see "Where the tune stops").
	Tonic   []harmony.Tonality
	Arrives int

	// Stops tells the phrase the tune ends on, on its last chorus.
	Stops bool
}

// Phrases cuts a sequence into phrases. A phrase ends when it comes to
// rest, or when the tune stops.
//
// # Coming to rest
//
// « La cadence conclusive est le point d'arrivée d'une phrase
// harmonique », and in a tune that keeps to the carrure, the conclusive
// cadences mark the groups of bars (Siron, La partition intérieure,
// p. 390). A phrase of a sequence with bars comes to rest where a
// section of the form concludes (see [Sections], [Conclusions]), and
// nowhere else, unless on the IV of the tonic it opens on (see
// onTheFourth):
//
//   - a tonic reached in the middle of a section confirms the tonic
//     under way, and ends no phrase: the E♭maj7 at bar 3 of Let's Cool
//     One;
//   - a section that ends on its two five, its I on the first bar of
//     the next one, ends open on its V, the half cadence of En Harmonie
//     (tome 1, chapter 8): the I opens the next phrase. So the words
//     say in A Fine Romance.
//
// A sequence without bars, a performance heard live, has no form to
// read yet: there, a phrase comes to rest on a tonic chord that a
// cadence leads to, perfect or plagal, and that either
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
// On its last chorus, a tune ends where the chart says, when it says:
// after its coda, at its "Fine", or on the chord it holds under a
// fermata (see [Changes]). A sequence that does not loop ends on its
// last chord. The tune stops on the tonic heard there: that chord's, in
// root position or on its fifth, or the one a cadence installs on it.
//
// A chart that says nothing ends as its last section does:
//
//   - on its conclusion, when it has one: the turnaround after it goes
//     back to the first chord and is not played at the end. Lullaby Of
//     Birdland stops on A♭maj7 before Gm7♭5 C7;
//   - across the loop, when it ends open on its V: the last time, the V
//     resolves on the first chord, and the tune stops there. Yesterdays
//     ends on A7, and stops on the Dm it goes back to; Sugar on G7, and
//     stops on Cm7.
//
// When its last section does neither, and in a sequence without bars,
// which has no form, the last phrase stops on the last tonic heard. A tonic already installed, the one the tune
// opened on, the one the phrase before came to rest on or, with none,
// the one its first cadence went to, needs no cadence to come back:
// 'Round Midnight stops on E♭6, Chega De Saudade on D6. Another one
// comes after a two five one or a plagal cadence, not a lone V, and not
// on the IV of the tonic before.
func Phrases(c Changes, blocks []Block) []Phrase {
	if opening(c) < 0 {
		return nil
	}
	ph := &phrasing{c: c, blocks: blocks, r: rolesOf(c, blocks), bar: barOf(c)}
	ph.opens = opensOn(c, blocks, ph.r)
	ph.home = ph.opens
	sections := Sections(c)
	if len(sections) > 0 {
		ph.concludes = map[int]bool{}
		for _, k := range Conclusions(c, sections, blocks) {
			if k.Arrives >= 0 && !ph.onTheFourth(k.Tonic) {
				ph.concludes[k.Arrives] = true
			}
		}
	}
	ph.rests()
	if ph.concludes != nil {
		ph.ends(sections[len(sections)-1])
	} else {
		ph.stops()
	}
	return ph.out
}

// A phrasing cuts changes into phrases.
type phrasing struct {
	c      Changes
	blocks []Block
	r      roles
	bar    Ticks
	opens  []harmony.Tonality // the tonic the tune opens on, if any

	// concludes holds the changes where a section concludes, nil for a
	// sequence without bars (see "Coming to rest").
	concludes map[int]bool

	out    []Phrase
	p      Phrase             // the phrase being cut
	before []harmony.Tonality // the tonic the phrase before came to rest on
	first  []harmony.Tonality // the tonic the first cadence resolves on
	home   []harmony.Tonality // the tonic opened on, else the first rested on
}

// cadenced returns the tonic a cadence within the chorus leads change i
// to, and that cadence: see [cadencedAt].
func (ph *phrasing) cadenced(i int) ([]harmony.Tonality, *Block) {
	return cadencedAt(ph.c, ph.blocks, ph.r, i)
}

// restsAt tells the tonic a phrase comes to rest on at change i, nil
// when it does not: see "Coming to rest".
func (ph *phrasing) restsAt(i int) []harmony.Tonality {
	t, b := ph.cadenced(i)
	if ph.concludes != nil {
		if !ph.concludes[i] {
			return nil
		}
		return t
	}
	switch {
	case t == nil:
		return nil
	case sameTonality(t, ph.opens):
		return t
	case tonicOf(ph.c.Chords[i]) == nil, ph.subdominant(t):
		return nil
	}
	if !ph.holds(i, b) {
		return nil
	}
	return t
}

// holds reports whether the tonic at change i, that block b leads to,
// holds longer than a bar and longer than each chord of b.
func (ph *phrasing) holds(i int, b *Block) bool {
	h := held(ph.c, i)
	if h <= ph.bar {
		return false
	}
	for _, k := range []int{b.Two, b.Sus, b.Five} {
		if k >= 0 && ph.c.Chords[k].Length >= h {
			return false
		}
	}
	return true
}

// rests cuts the phrases that come to rest, from the opening on.
func (ph *phrasing) rests() {
	o := opening(ph.c)
	ph.p = Phrase{From: o}
	ph.before = tonicOf(ph.c.Chords[o])
	for i := o + 1; i < len(ph.c.Chords); i++ {
		if t, _ := ph.cadenced(i); ph.first == nil && t != nil {
			ph.first = t
		}
		if t := ph.restsAt(i); t != nil {
			to := heldTo(ph.c, i)
			ph.conclude(i, to, t, false)
			i = to
		}
	}
}

// stopsAt tells the tonic the tune stops on at change i, nil when it
// does not: see "Where the tune stops". An installed tonic comes back
// in any position, its second inversion included: B♭maj7/F before G7
// Cm7 F7 at the end of Someday My Prince Will Come. Any other bass makes
// another chord: Fmaj7/G at the end of Only Trust Your Heart is the V
// of C with its fourth, not F.
func (ph *phrasing) stopsAt(i int, installed []harmony.Tonality) []harmony.Tonality {
	root := ph.c.Chords[i]
	if above := (int(root.Bass) - int(root.Chord.Root) + 12) % 12; above == 7 {
		root.Bass = root.Chord.Root
	}
	comesBack := func(t []harmony.Tonality) bool {
		return t != nil && (sameTonic(t, ph.opens) || sameTonic(t, installed))
	}
	if t := tonicOf(root); comesBack(t) {
		return t
	}
	if t := minorSeventhTonic(ph.c, ph.blocks, ph.r, i); comesBack(t) {
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

// stops cuts the last phrase, from the last rest to the last tonic
// heard, and the tail after it. It reports whether it found one.
func (ph *phrasing) stops() bool {
	if ph.p.From >= len(ph.c.Chords) {
		return false
	}
	installed := ph.before
	if installed == nil {
		installed = ph.first
	}
	found := false
	for i := len(ph.c.Chords) - 1; i >= ph.p.From; i-- {
		if t := ph.stopsAt(i, installed); t != nil {
			ph.conclude(i, i, t, true)
			found = true
			break
		}
	}
	ph.tail()
	return found
}

// ends marks where a tune with bars ends on its last chorus, whose
// last section is `last`: see "Where the tune stops".
func (ph *phrasing) ends(last Section) {
	c := ph.c
	e := -1
	switch {
	case c.End > 0:
		e = c.End
	case c.Coda > 0, !c.Loops:
		e = len(c.Chords) - 1
	}
	if e >= 0 {
		if t := ph.endsOn(e); t != nil && !ph.onTheFourth(t) {
			ph.stopAt(e, t)
			return
		}
	}
	n := len(ph.out)
	concluded := n > 0 && c.Bar(c.Chords[ph.out[n-1].Arrives].Start) >= last.From
	if !concluded && c.Loops && ph.p.From < len(c.Chords) {
		if o, t := ph.acrossTheLoop(); t != nil {
			ph.p.Arrives, ph.p.To, ph.p.Tonic, ph.p.Stops = o, len(c.Chords)-1, t, true
			ph.out = append(ph.out, ph.p)
			return
		}
	}
	if concluded {
		ph.out[n-1].Stops = true
		ph.tail()
		return
	}
	if !ph.stops() && n > 0 {
		ph.out[n-1].Stops = true
	}
}

// onTheFourth reports whether tonic `t` is the IV of the tonic the tune
// opens on, when a cadence comes back to that tonic within the chorus:
// no phrase rests there, and no tune stops there, even where the chart
// marks its end. Unforgettable goes from G to Cmaj7 as a blues goes
// from I to IV, its phrase said again a fourth higher, and its chart
// ends on that Cmaj7 at bar 31: Am7 D7 after it goes back to G, where
// the tune stops. Somewhere opens on B♭, but B♭ turns into B♭7 and no
// cadence comes back to it: its E♭ is no IV.
func (ph *phrasing) onTheFourth(t []harmony.Tonality) bool {
	if ph.opens == nil || t[0].Tonic() != ph.opens[0].Tonic().Transpose(5) {
		return false
	}
	for i, ch := range ph.c.Chords {
		if back, _ := ph.cadenced(i); back != nil && sameTonic(back, ph.opens) && tonicOf(ch) != nil {
			return true
		}
	}
	return false
}

// acrossTheLoop returns the first change, and the tonic the cadence at
// the end of the chart resolves on there, going round: A7 Dm in
// Yesterdays, G7 Cm7 in Sugar. The first chord must be a tonic the tune
// opens on (see opensOn): a m7 that no cadence resolves on inside the
// chorus is a VI or a II, the Fm7 of All The Things You Are. nil when
// the last chords lead nowhere, or elsewhere.
func (ph *phrasing) acrossTheLoop() (int, []harmony.Tonality) {
	o := opening(ph.c)
	n := ph.r.target[o]
	if n < 0 || ph.blocks[n].Five <= o || !sameTonic(ph.blocks[n].Announced, ph.opens) {
		return o, nil
	}
	return o, ph.opens
}

// endsOn returns the tonic the tune ends on at change i, the end the
// chart marks: the tonic chord there, in root position or on its fifth,
// or the tonic a cadence installs; nil when there is none. A m7 alone
// is not enough: Wave and Triste end their charts on a m7 that is a II.
func (ph *phrasing) endsOn(i int) []harmony.Tonality {
	root := ph.c.Chords[i]
	if above := (int(root.Bass) - int(root.Chord.Root) + 12) % 12; above == 7 {
		root.Bass = root.Chord.Root
	}
	if t := tonicOf(root); t != nil {
		return t
	}
	if t := minorSeventhTonic(ph.c, ph.blocks, ph.r, i); t != nil {
		return t
	}
	t, _ := ph.cadenced(i)
	return t
}

// stopAt makes the tune stop on tonic `t` at change `e`: the phrase
// already cut there when it rests on that tonic, else a phrase that
// ends there, cut from the one that holds `e`. The phrases after it are
// played on the choruses before the last. Somewhere rests on A♭ at the
// end of its last section, and the chart ends it two bars before, on
// E♭, which a phrase of its own concludes.
func (ph *phrasing) stopAt(e int, t []harmony.Tonality) {
	to := heldTo(ph.c, e)
	for k := range ph.out {
		p := ph.out[k]
		switch {
		case e < p.From || e > p.To:
			continue
		case sameTonic(p.Tonic, t):
			ph.out[k].Stops = true
		default:
			stop := Phrase{From: p.From, To: min(to, p.To), Tonic: t, Arrives: e, Stops: true}
			p.From = stop.To + 1
			ph.out = slices.Insert(ph.out, k, stop)
			if p.From <= p.To {
				ph.out[k+1] = p
			} else {
				ph.out = slices.Delete(ph.out, k+1, k+2)
			}
		}
		ph.tail()
		return
	}
	ph.conclude(e, to, t, true)
	ph.tail()
}

// tail adds the phrase after the last conclusion, if any is left.
func (ph *phrasing) tail() {
	if ph.p.From < len(ph.c.Chords) {
		ph.p.To, ph.p.Arrives = len(ph.c.Chords)-1, -1
		ph.out = append(ph.out, ph.p)
	}
}

// conclude ends the phrase being cut on tonic t, reached at change
// arrives and held to change to.
func (ph *phrasing) conclude(arrives, to int, t []harmony.Tonality, stops bool) {
	ph.p.Arrives, ph.p.To, ph.p.Tonic, ph.p.Stops = arrives, to, t, stops
	ph.out = append(ph.out, ph.p)
	ph.p, ph.before = Phrase{From: to + 1}, t
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

// Home reads where the tune sets out from: the tonic its first phrase
// concludes on, nil when none does. Autumn Leaves sets out from G
// minor, Fly Me To The Moon from A minor, and Just Friends from G,
// which it confirms twice before its first conclusion, at bar 31.
func Home(c Changes, phrases []Phrase) []harmony.Tonality {
	if t, ok := Blues(c); ok {
		return t
	}
	if first, _ := concluding(phrases); first != nil {
		return first.Tonic
	}
	return nil
}

// Tune reads the tonality of the tune as En Harmonie does (tome 1,
// chapter 8 §1.2, p. 99 and 100): its first and its last chord,
// turnaround left out, confirm it when they give the same tonality,
// Blame It On My Youth and Angel Eyes. When they part, « la
// prédominance de l'une ou l'autre des deux tonalités durant le
// morceau » decides: My Funny Valentine starts on Cm, ends on E♭6, and
// is in C minor. Neither the book nor Siron says what predominance is:
// here, the one heard the longer (see [heardFor]).
//
// The first chord counts when it is the tonic where it stands (see
// [firstTonic]); the last one is where the tune stops (see
// [Phrases]), in the major or the three minors as its tonic's third
// says. A picardy third does not make a minor tune major (see
// [Picardy]). A blues is in its own tonic, found by its form. A tune
// that stops nowhere is in its home.
func Tune(c Changes, phrases []Phrase) []harmony.Tonality {
	if t, ok := Blues(c); ok {
		return t
	}
	end := stopping(phrases)
	if end == nil {
		return Home(c, phrases)
	}
	last := end.Tonic
	if Picardy(c, phrases) {
		last = MinorTonalities(last[0].Tonic())
	}
	blocks := Blocks(c, Approaches(c))
	first := firstTonic(c, blocks)
	if first == nil || sameTonic(first, last) {
		return last
	}
	heard := heardFor(c, blocks)
	if heard[first[0].Tonic()] > heard[last[0].Tonic()] {
		return first
	}
	return last
}

// firstTonic returns the tonic of the first chord, when it is the tonic
// in its context, nil otherwise. The book asks for the first chord and
// says no more; what makes it the tonic where it stands is ours to
// say: the tune's first cadence resolves on its tonic. My Funny
// Valentine opens on Cm6, and Dm7♭5 G7♭9 goes back to Cm7; In a
// Sentimental Mood on Dm, and A7 goes back to it. Just Friends opens on
// Cmaj7, and its first cadence goes to G: C is its IV. Blue Skies opens
// on Am, and its first cadence goes to C. The Fm7 that opens All The
// Things You Are is a VI, the Em7 of Tune Up a II.
func firstTonic(c Changes, blocks []Block) []harmony.Tonality {
	r := rolesOf(c, blocks)
	first := opensOn(c, blocks, r)
	if first == nil {
		return nil
	}
	for i := range c.Chords {
		if t, _ := cadencedAt(c, blocks, r, i); t != nil {
			if sameTonic(t, first) {
				return first
			}
			return nil
		}
	}
	return nil
}

// heardFor returns how long each tonic is heard: the tonal areas for
// theirs (see [TonalAreas]), and the rest of the tune for the tonic it
// is first heard in. A tonicisation, shorter than an area, counts for
// the tonic around it. The tune is heard without its phrases, which
// would install the home where the first one concludes: In a
// Sentimental Mood would be heard in F from its first bar.
func heardFor(c Changes, blocks []Block) map[harmony.PitchClass]Ticks {
	sensed := Sense(c, blocks, nil, nil)
	out := map[harmony.PitchClass]Ticks{}
	var first []harmony.Tonality
	for _, s := range sensed {
		if first = s.Region; first == nil {
			first = s.Ground
		}
		if first != nil {
			break
		}
	}
	if first == nil {
		return out
	}
	for _, ch := range c.Chords {
		out[first[0].Tonic()] += ch.Length
	}
	for _, a := range TonalAreas(c, blocks, sensed, Sections(c)) {
		for i := a.From; i <= a.To; i++ {
			out[first[0].Tonic()] -= c.Chords[i].Length
			out[a.Tonic[0].Tonic()] += c.Chords[i].Length
		}
	}
	return out
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
//
// Nor is it one when a cadence of the major has installed the major
// tonic before the end: a two five whose two is a m7, heard as the II
// of the major where a m7♭5 is the II of the minor. Somewhere passes
// through E♭m by a IV-I in the middle of a phrase, and goes to E♭ by
// Fm7 B♭7: it is not a minor tune.
func Picardy(c Changes, phrases []Phrase) bool {
	if _, ok := Blues(c); ok {
		return false
	}
	end := stopping(phrases)
	if end == nil || ModesOf(end.Tonic) != Major {
		return false
	}
	if majorInstalled(c, end) {
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

// majorInstalled reports whether a cadence announcing the major leads
// to the major tonic chord of `end` before the tune stops there.
func majorInstalled(c Changes, end *Phrase) bool {
	blocks := Blocks(c, Approaches(c))
	r := rolesOf(c, blocks)
	for i := range c.Chords[:end.Arrives] {
		t, b := cadencedAt(c, blocks, r, i)
		if t != nil && sameTonic(t, end.Tonic) && ModesOf(t) == Major && ModesOf(b.Announced)&Major != 0 {
			return true
		}
	}
	return false
}

// stopping returns the phrase the tune stops on, else the last that
// concludes, nil when none does.
func stopping(phrases []Phrase) *Phrase {
	for i := range phrases {
		if phrases[i].Stops {
			return &phrases[i]
		}
	}
	_, last := concluding(phrases)
	return last
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
