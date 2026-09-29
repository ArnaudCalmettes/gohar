package analysis

import (
	"slices"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// The tonality of a tune is not one fact but a verdict, and no single
// rule reaches it: the end decides All The Things You Are, the start In
// a Sentimental Mood, and Jordu comes back to C minor more often than
// it goes anywhere else. So the tonics a tune gives are gathered as
// candidates, each with its evidence, the events that make it heard as
// the reference. The verdict can then be both informed and explained.
//
// Evidence accumulates as the tune goes: the tonic with the most of it
// at a given bar is the one felt as the tune's at that bar, and the one
// with the most at the end is the tune's. The local tonic of [Sense] is
// another matter: a passage can modulate for a while, and tonicise a
// degree here and there, without the reference changing.

// A Proof is a kind of event that makes a tonic heard as the reference.
type Proof int

const (
	// Returns: a cadence within the chorus resolves on the tonic, a
	// tonic chord. The turnaround back to the first chord does not
	// count: going round again is no new arrival. Nor does a cadence on
	// a m7, most often a degree tonicised by its V, the II or the VI: A7
	// Dm7 in C, the V of the II. The m7 tonic of Softly, As In A Morning
	// Sunrise counts where the tune opens or stops on it.
	Returns Proof = iota

	// Opens: the tune opens on the tonic (see [Phrases]).
	Opens

	// Rests: a phrase comes to rest on the tonic (see [Phrases]).
	Rests

	// Stops: the tune stops on the tonic (see [Phrases]).
	Stops

	// Form: the tune is a blues on the tonic, found by its form (see
	// [Blues]). A blues gives no other evidence.
	Form
)

var proofNames = [...]string{"returns", "opens", "rests", "stops", "form"}

func (p Proof) String() string {
	return proofNames[p]
}

// Evidence is an event that makes a tonic heard as the reference.
type Evidence struct {
	Proof Proof
	At    int // the change where it is heard
}

// A Candidate is a tonic the tune could be in, with its evidence in the
// order it is heard. A tonic and its parallel are two candidates: D
// minor and D major in Chega De Saudade.
type Candidate struct {
	Tonic    []harmony.Tonality
	Evidence []Evidence
}

// Candidates gathers the tonics a tune gives and their evidence, one
// proof per event: a cadence that brings a phrase to rest is a rest,
// not a rest and a return; a tonic held eight bars is one event. The
// candidates come with the most evidence first, then in the order they
// are first heard.
//
// Candidates reads the evidence, it does not decide: [Tune] still does.
func Candidates(c Changes, blocks []Block, phrases []Phrase) []Candidate {
	if t, ok := Blues(c); ok {
		return []Candidate{{Tonic: t, Evidence: []Evidence{{Form, 0}}}}
	}
	r := rolesOf(c, blocks)
	type event struct {
		Evidence
		tonic []harmony.Tonality
	}
	events := map[int]event{}
	record := func(i int, p Proof, t []harmony.Tonality) {
		if e, ok := events[i]; !ok || p > e.Proof {
			events[i] = event{Evidence{p, i}, t}
		}
	}
	for i := range c.Chords {
		if t, _ := cadencedAt(c, blocks, r, i); t != nil && tonicOf(c.Chords[i]) != nil {
			record(i, Returns, t)
		}
	}
	if t := opensOn(c, blocks, r); t != nil {
		record(opening(c), Opens, t)
	}
	for _, p := range phrases {
		switch {
		case p.Tonic == nil:
		case p.Stops:
			record(p.Arrives, Stops, p.Tonic)
		default:
			record(p.Arrives, Rests, p.Tonic)
		}
	}

	var out []Candidate
	for i := range c.Chords {
		e, ok := events[i]
		if !ok {
			continue
		}
		k := slices.IndexFunc(out, func(cd Candidate) bool { return sameTonality(cd.Tonic, e.tonic) })
		if k < 0 {
			out, k = append(out, Candidate{Tonic: e.tonic}), len(out)
		}
		out[k].Evidence = append(out[k].Evidence, e.Evidence)
	}
	slices.SortStableFunc(out, func(a, b Candidate) int {
		return len(b.Evidence) - len(a.Evidence)
	})
	return out
}
