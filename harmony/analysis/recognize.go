package analysis

import (
	"fmt"
	"sort"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Recognizer identifies chords in snapshots.
//
// Built once over a table of tetrads, then called on every frame.
// Immutable after construction and safe to share.
type Recognizer struct {
	tetrads map[harmony.ChordPattern]struct{}
	scratch []harmony.Pitch
}

// NewRecognizer builds a recognizer over the given tetrads.
//
// The caller supplies the table. A game teaching triads passes triads
// and gets fewer readings for it; passing everything means a player
// holding three notes sees them explained several ways, which is
// correct and rarely what a beginner needs.
//
// The patterns are taken as tetrads, so they must be given collapsed
// into one octave. [CommonTetrads] is the usual set.
func NewRecognizer(tetrads ...harmony.ChordPattern) (*Recognizer, error) {
	if len(tetrads) == 0 {
		return nil, fmt.Errorf("analysis: a recognizer needs at least one tetrad")
	}

	table := make(map[harmony.ChordPattern]struct{}, len(tetrads))
	for _, t := range tetrads {
		if !t.HasOffset(0) {
			return nil, fmt.Errorf("analysis: a tetrad must contain its own root")
		}
		if t.Tetrad() != t {
			return nil, fmt.Errorf(
				"analysis: tetrad %v reaches past the octave; give it collapsed", t)
		}
		table[t] = struct{}{}
	}
	return &Recognizer{tetrads: table}, nil
}

// CommonTetrads is the table a keyboard game will want by default.
//
// The omit-five variants are in it on purpose. A fifthless voicing is
// not a chord with a note missing, it is how the chord is played, and
// a table that lacks it would report no chord where a jazz pianist
// hears an obvious one.
var CommonTetrads = []harmony.ChordPattern{
	harmony.ChordMajorTriad,
	harmony.ChordMinorTriad,
	harmony.ChordDiminishedTriad,
	harmony.ChordAugmentedTriad,
	harmony.ChordSus2,
	harmony.ChordSus4,
	harmony.ChordMajorSixth,
	harmony.ChordMinorSixth,
	harmony.ChordMajorSeventh,
	harmony.ChordMajorSeventhNo5,
	harmony.ChordMajorSeventhSharp5,
	harmony.ChordDominantSeventh,
	harmony.ChordDominantSeventhNo5,
	harmony.ChordDominantSeventhFlat5,
	harmony.ChordDominantSeventhSharp5,
	harmony.ChordMinorSeventh,
	harmony.ChordMinorSeventhNo5,
	harmony.ChordHalfDiminished,
	harmony.ChordMinorMajorSeventh,
	harmony.ChordMinorMajorSeventhNo5,
	harmony.ChordDiminishedSeventh,
	harmony.ChordDominantSeventhSus2,
	harmony.ChordDominantSeventhSus4,
}

// Identify returns the readings of s, best first.
//
// # How one reading is built
//
// Take a sounding pitch as candidate root. Gather every sounding note
// into the octave above it, giving a raw pattern. Normalise it, which
// moves each member onto the rung harmony assigns it. Isolate the
// tetrad and look it up. A hit is a reading; a miss is nothing at all,
// and this root is simply not the root.
//
// # How readings are ordered
//
// Rules, in this order, each breaking the ties the previous left:
//
//  1. A reading equal to [Context.Current] comes first. A still valid
//     answer is not replaced by an equal one, which is what keeps a
//     display still under a diminished seventh chord.
//  2. Root position beats inversion. With no tonal context this is
//     about the best that can be said, and it is what a listener does.
//  3. A root belonging to [Context.Tonality] beats one outside it.
//     Ordering only: a chord foreign to the key is still identified,
//     just later in the list.
//  4. Fewer extensions beats more. A reading that explains the notes
//     as a plain seventh chord is preferred to one that needs an
//     altered thirteenth to say the same thing.
//  5. Table order, which is stable and arbitrary. A caller that cares
//     past this point breaks the tie itself.
//
// # What it does not do
//
// It does not decide. Four roots explain a diminished seventh chord and the
// notes do not say which the player meant. All four come back.
//
// It does not allocate on the hot path. The returned slice is reused
// between calls; copy a reading to keep it.
func (r *Recognizer) Identify(s Snapshot, ctx Context) []Reading {
	r.scratch = r.scratch[:0]
	if s.IsEmpty() {
		return nil
	}

	var readings []Reading
	seen := harmony.EmptyPitchSet
	for _, p := range s.Pitches {
		if seen.Contains(p.Class()) {
			continue
		}
		seen = seen.With(p.Class())
		if reading, ok := r.IdentifyRoot(p, s); ok {
			readings = append(readings, reading)
		}
	}

	sort.SliceStable(readings, func(i, j int) bool {
		return rank(readings[i], ctx) < rank(readings[j], ctx)
	})
	return readings
}

// Best returns the first reading, and whether there is one.
//
// It hides exactly the ambiguity the list exists to show, so prefer
// the list wherever the difference between one answer and four matters
// to the player.
func (r *Recognizer) Best(s Snapshot, ctx Context) (Reading, bool) {
	readings := r.Identify(s, ctx)
	if len(readings) == 0 {
		return Reading{}, false
	}
	return readings[0], true
}

// IdentifyRoot returns the reading of s under the given root, and
// whether the notes form a chord on it at all.
//
// The single step [Recognizer.Identify] runs once per sounding note.
// Exposed because the engine wants to re-test one root cheaply: when
// checking whether the current reading still holds, there is no reason
// to search the other eleven.
//
// The root has to sound. A chord pattern carries its root as an implied
// bit, so without that rule any class at all could be made the root of
// a chord it takes no part in: C E G read from a C sharp gives a
// diminished triad with the C as an upper structure, which is a reading
// nobody hears. Identify never meets the case, since it only ever tries
// sounding notes. This is also what the engine wants: a reading whose
// root has been released is over.
func (r *Recognizer) IdentifyRoot(root harmony.Pitch, s Snapshot) (Reading, bool) {
	if s.IsEmpty() {
		return Reading{}, false
	}

	rootClass := root.Class()
	if !s.Classes().Contains(rootClass) {
		return Reading{}, false
	}

	raw := harmony.ChordPattern(1)
	for _, p := range s.Pitches {
		raw = raw.With(rootClass.Up(p.Class()))
	}

	normalized := raw.Normalize()
	tetrad := normalized.Tetrad()
	if !r.Knows(tetrad) {
		return Reading{}, false
	}

	bass, _ := s.Bass()
	return Reading{
		Root:       rootClass,
		Pattern:    normalized,
		Tetrad:     tetrad,
		Extensions: normalized &^ tetrad,
		BassIsRoot: bass.Class() == rootClass,
	}, true
}

// Knows reports whether the recognizer's table holds a tetrad.
func (r *Recognizer) Knows(tetrad harmony.ChordPattern) bool {
	_, ok := r.tetrads[tetrad]
	return ok
}

// rank orders readings by the rules stated on [Recognizer.Identify].
//
// A single integer rather than a comparison chain, so that the order of
// the rules is written once and reads top to bottom. Lower comes first.
// The last rule, table order, is left to the stability of the sort
// rather than encoded here.
func rank(r Reading, ctx Context) int {
	score := 0
	if !sameReading(r, ctx.Current) {
		score += 8
	}
	if !r.BassIsRoot {
		score += 4
	}
	if !ctx.Tonality.IsZero() && !ctx.Tonality.Contains(r.Root) {
		score += 2
	}
	if !r.Extensions.IsSubsetOf(0) {
		score++
	}
	return score
}

func sameReading(a, b Reading) bool {
	return !b.IsZero() && a.Root == b.Root && a.Pattern == b.Pattern
}
