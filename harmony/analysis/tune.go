package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// The tonality of the tune, read from its phrases as En Harmonie reads
// it: its first chord, its last one, and when they part, the
// predominance of one or the other (tome 1, chapter 8 §1.2). Where the
// tune stops comes from the phrases (see [Phrases]).

// Tune reads the tonality of the tune: see [ReadTune].
func Tune(c Changes, phrases []Phrase) []harmony.Tonality {
	return ReadTune(c, phrases).Tonality
}

// A TuneReading is what the tonality of a tune rests on, as [ReadTune]
// weighs it.
type TuneReading struct {
	// Tonality is the tonality of the tune, nil when none is heard.
	Tonality []harmony.Tonality

	// Blues tells a blues, in its own tonic, found by its form: nothing
	// else is weighed.
	Blues bool

	// Opens is the first chord that sounds, -1 when none does; First
	// the tonic it gives, nil when the first cadence of the tune
	// resolves elsewhere or on nothing (see [FirstTonic]); and Cadence
	// the tonic that cadence resolves on, nil when none resolves.
	Opens   int
	First   []harmony.Tonality
	Cadence []harmony.Tonality

	// Stops is the change the tune stops on, -1 when it stops nowhere,
	// and Last its tonic, in the minor when the last chord is a picardy
	// third (see [Picardy]).
	Stops   int
	Last    []harmony.Tonality
	Picardy bool

	// Heard is how long, in bars, the first and the last tonics are
	// heard, when they part (see heardFor); nil when they do not.
	Heard map[harmony.PitchClass]float64
}

// ReadTune reads the tonality of the tune as En Harmonie does (tome 1,
// chapter 8 §1.2, p. 99 and 100): its first and its last chord,
// turnaround left out, confirm it when they give the same tonality,
// Blame It On My Youth and Angel Eyes. When they part, « la
// prédominance de l'une ou l'autre des deux tonalités durant le
// morceau » decides: My Funny Valentine starts on Cm, ends on E♭6, and
// is in C minor. Neither the book nor Siron says what predominance is:
// here, the one heard the longer (see heardFor).
//
// The first chord counts when it is the tonic where it stands (see
// [FirstTonic]); the last one is where the tune stops (see
// [Phrases]), in the major or the three minors as its tonic's third
// says. A picardy third does not make a minor tune major (see
// [Picardy]). A blues is in its own tonic, found by its form. A tune
// that stops nowhere has no tonality heard.
func ReadTune(c Changes, phrases []Phrase) TuneReading {
	r := TuneReading{Opens: opening(c), Stops: -1}
	if t, ok := Blues(c); ok {
		r.Tonality, r.Blues = t, true
		return r
	}
	blocks := Blocks(c, Approaches(c))
	r.First, r.Cadence = FirstTonic(c, blocks), firstCadence(c, blocks)
	end := stopping(phrases)
	if end == nil {
		return r
	}
	r.Stops, r.Last = end.Arrives, end.Tonic
	if r.Picardy = Picardy(c, phrases); r.Picardy {
		r.Last = MinorTonalities(r.Last[0].Tonic())
	}
	r.Tonality = r.Last
	if r.First == nil || sameTonic(r.First, r.Last) {
		return r
	}
	heard := heardFor(c, blocks)
	bar := float64(barOf(c))
	f, l := r.First[0].Tonic(), r.Last[0].Tonic()
	r.Heard = map[harmony.PitchClass]float64{f: float64(heard[f]) / bar, l: float64(heard[l]) / bar}
	if heard[f] > heard[l] {
		r.Tonality = r.First
	}
	return r
}

// FirstTonic returns the tonic of the first chord, when it is the tonic
// in its context, nil otherwise: the tonality a tune sets out from, one
// of the two [Tune] weighs. The book asks for the first chord and
// says no more; what makes it the tonic where it stands is ours to
// say: the tune's first cadence resolves on its tonic. My Funny
// Valentine opens on Cm6, and Dm7♭5 G7♭9 goes back to Cm7; In a
// Sentimental Mood on Dm, and A7 goes back to it. Just Friends opens on
// Cmaj7, and its first cadence goes to G: C is its IV. Blue Skies opens
// on Am, and its first cadence goes to C. The Fm7 that opens All The
// Things You Are is a VI, the Em7 of Tune Up a II.
func FirstTonic(c Changes, blocks []Block) []harmony.Tonality {
	first := opensOn(c, blocks, rolesOf(c, blocks))
	if first == nil || !sameTonic(first, firstCadence(c, blocks)) {
		return nil
	}
	return first
}

// firstCadence returns the tonic the first cadence of the tune resolves
// on, nil when none does.
func firstCadence(c Changes, blocks []Block) []harmony.Tonality {
	r := rolesOf(c, blocks)
	for i := range c.Chords {
		if t, _ := cadencedAt(c, blocks, r, i); t != nil {
			return t
		}
	}
	return nil
}

// heardFor returns how long each tonic is heard: the tonal areas for
// theirs (see [TonalAreas]), and the rest of the tune for the tonic it
// is first heard in. A tonicisation, shorter than an area, counts for
// the tonic around it. The tune is heard without a tonality given.
func heardFor(c Changes, blocks []Block) map[harmony.PitchClass]Ticks {
	sensed := Sense(c, blocks, nil)
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
// And the tune's first cadence resolves on that tonic, minor or
// already major: a tune whose first resolution goes elsewhere does not
// set out from it. 'Round Midnight resolves on E♭m after its
// introduction. Somewhere first resolves on A♭, hints at E♭ by avoided
// cadences, touches E♭m in its bridge, and resolves on E♭ at the very
// end: it is in E♭, its melody leaving no doubt.
func Picardy(c Changes, phrases []Phrase) bool {
	if _, ok := Blues(c); ok {
		return false
	}
	end := stopping(phrases)
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
	return first != nil && ModesOf(first) != Major && ModesOf(last) != Major &&
		sameTonic(firstCadence(c, Blocks(c, Approaches(c))), first)
}
