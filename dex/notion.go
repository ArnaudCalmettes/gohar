// Package dex holds the collection of musical notions a player builds
// up, shared by every game.
//
// # What it is for
//
// It celebrates what the player knows and never shows what is left.
// Learning music is a lifetime's work and the total means nothing, so
// nothing here computes a global denominator: no total, no percentage,
// no progress bar. Counting inside an ensemble that music itself bounds
// is the one exception, seven degrees in a system or twelve tonics for
// a notion, because those are theoretical facts rather than
// administrative targets.
//
// # What it does not know
//
// Any game. Games send facts, the dex decides what they mark, and it
// answers questions in return. Nothing in here mentions a score, a
// combo or a difficulty: the day a fact carries one, the dex stops
// being common ground.
package dex

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Kind says which family a notion belongs to.
type Kind uint8

const (
	// KindMode is a mode, designated by its mother scale and degree.
	KindMode Kind = iota + 1

	// KindTetrad is a chord shape, collapsed into one octave.
	KindTetrad

	// KindTetrachord is a four note cell.
	KindTetrachord

	// KindProgression is an entry of the progression catalogue.
	KindProgression

	// KindInterval is the gap between two notes, from the unison to
	// the octave and beyond. An elementary notion: see
	// [Notion.Elementary].
	KindInterval
)

// A ProgressionID designates one entry of the progression catalogue.
//
// The catalogue is not written yet. It will hold the two five one in
// major, the same in minor, the suspended two five, and whatever else
// the games need, each form counting as an entry of its own rather than
// as a variant of another.
//
// Unlike the modes, this catalogue has no natural bound: it grows with
// the repertoire. That is why nothing ever displays a count over it.
type ProgressionID uint16

// A Notion is what a dex entry is about.
//
// # Why a struct and not a name
//
// A notion is designated the way the library designates it, never by a
// free string. Two games writing "ii-V-I" and "2-5-1" would fill the
// collection with entries that ignore each other, and the whole point
// of the theory layer is that the same shape is the same value.
//
// Only the fields its kind uses carry meaning. The type stays
// comparable so that it can key a map.
//
// # Persistence
//
// [Notion.String] is the stable form to write down, not the numeric
// values of the fields. Those are Go constants and may be reordered;
// the text may not.
type Notion struct {
	Kind Kind

	// System and Degree designate a mode.
	System harmony.System
	Degree harmony.Degree

	// Tetrad designates a chord shape.
	Tetrad harmony.ChordPattern

	// Tetrachord designates a four note cell.
	Tetrachord harmony.Tetrachord

	// Progression designates a catalogue entry.
	Progression ProgressionID

	// Interval designates a gap, degree span and semitones both, so
	// that an augmented fourth and a diminished fifth stay two notions.
	Interval harmony.Interval
}

// ModeOf designates a mode by its mother scale and degree.
func ModeOf(system harmony.System, degree harmony.Degree) Notion {
	return Notion{Kind: KindMode, System: system, Degree: degree}
}

// TetradOf designates a chord shape.
func TetradOf(p harmony.ChordPattern) Notion {
	return Notion{Kind: KindTetrad, Tetrad: p}
}

// TetrachordOf designates a four note cell.
func TetrachordOf(t harmony.Tetrachord) Notion {
	return Notion{Kind: KindTetrachord, Tetrachord: t}
}

// ProgressionOf designates a catalogue entry.
func ProgressionOf(id ProgressionID) Notion {
	return Notion{Kind: KindProgression, Progression: id}
}

// IntervalOf designates an interval.
func IntervalOf(i harmony.Interval) Notion {
	return Notion{Kind: KindInterval, Interval: i}
}

// Elementary reports whether `n` is one of the small steps a beginner
// is rewarded for and a musician stops counting: an interval.
//
// An elementary notion never receives FactSounded. A turnaround of
// altered dominants sounds dozens of intervals, and saying the player
// heard a minor third there would mean nothing. It is marked only by
// an activity that asks for it, so its marks stop growing, quietly,
// once the player has moved on to what it builds. Nor does it show as
// cooling: a third mastered needs no review, the scales it makes up do.
func (n Notion) Elementary() bool {
	return n.Kind == KindInterval
}

// IsZero reports whether `n` designates nothing.
func (n Notion) IsZero() bool {
	return n.Kind == 0
}

// String returns the stable designation of `n`, as in "mode:2/4".
//
// This is what persistence writes and reads back, which is why it must
// not change once a dex exists in the wild. It is not a name: turning a
// notion into words a player reads is the naming package's business,
// and it needs a locale this package does not carry.
func (n Notion) String() string {
	switch n.Kind {
	case KindMode:
		return fmt.Sprintf("mode:%d/%d", n.System, n.Degree)
	case KindTetrad:
		return fmt.Sprintf("tetrad:%06x", uint32(n.Tetrad))
	case KindTetrachord:
		return fmt.Sprintf("tetrachord:%03x", uint16(n.Tetrachord))
	case KindProgression:
		return fmt.Sprintf("progression:%d", n.Progression)
	case KindInterval:
		return fmt.Sprintf("interval:%d/%d", n.Interval.Degrees, n.Interval.Semitones)
	}
	return "none"
}

// ParseNotion reads back what [Notion.String] wrote.
//
// Strict on purpose. A dex read from disk is the player's collection,
// and a designation that half parses would silently move a mark onto
// another notion. Anything that is not exactly a written form is an
// error rather than a guess.
//
// A mode is checked against the catalogue, since its two numbers
// designate something that either exists or does not. A tetrad and a
// tetrachord are only checked for width: an unusual shape is still a
// notion worth writing down, and it is the recognizer's table, not
// persistence, that decides which ones a game knows.
func ParseNotion(s string) (Notion, error) {
	if s == "none" {
		return Notion{}, nil
	}

	kind, rest, ok := strings.Cut(s, ":")
	if !ok {
		return Notion{}, fmt.Errorf("dex: %q is not a notion", s)
	}

	switch kind {
	case "mode":
		left, right, ok := strings.Cut(rest, "/")
		if !ok {
			return Notion{}, fmt.Errorf("dex: %q has no degree", s)
		}
		system, err := strconv.ParseUint(left, 10, 8)
		if err != nil {
			return Notion{}, fmt.Errorf("dex: %q has no system: %w", s, err)
		}
		degree, err := strconv.ParseUint(right, 10, 8)
		if err != nil {
			return Notion{}, fmt.Errorf("dex: %q has no degree: %w", s, err)
		}
		n := ModeOf(harmony.System(system), harmony.Degree(degree))
		// Mode reports false on an unknown system as well as on a degree
		// outside the seven, which covers both fields in one call.
		if _, ok := n.System.Mode(n.Degree); !ok {
			return Notion{}, fmt.Errorf("dex: %q designates no mode", s)
		}
		return n, nil

	case "tetrad":
		v, err := strconv.ParseUint(rest, 16, 32)
		if err != nil {
			return Notion{}, fmt.Errorf("dex: %q has no tetrad: %w", s, err)
		}
		return TetradOf(harmony.ChordPattern(v)), nil

	case "tetrachord":
		v, err := strconv.ParseUint(rest, 16, 16)
		if err != nil {
			return Notion{}, fmt.Errorf("dex: %q has no tetrachord: %w", s, err)
		}
		return TetrachordOf(harmony.Tetrachord(v)), nil

	case "progression":
		v, err := strconv.ParseUint(rest, 10, 16)
		if err != nil {
			return Notion{}, fmt.Errorf("dex: %q has no entry: %w", s, err)
		}
		return ProgressionOf(ProgressionID(v)), nil

	case "interval":
		left, right, ok := strings.Cut(rest, "/")
		if !ok {
			return Notion{}, fmt.Errorf("dex: %q has no semitones", s)
		}
		degrees, err := strconv.ParseUint(left, 10, 8)
		if err != nil || degrees > 14 {
			return Notion{}, fmt.Errorf("dex: %q has no degree span", s)
		}
		semitones, err := strconv.ParseUint(right, 10, 8)
		if err != nil || semitones > 24 {
			return Notion{}, fmt.Errorf("dex: %q has no semitones", s)
		}
		return IntervalOf(harmony.Interval{
			Degrees: harmony.Degrees(degrees), Semitones: harmony.Semitones(semitones),
		}), nil
	}

	return Notion{}, fmt.Errorf("dex: %q is of no known kind", s)
}

// Components returns the notions a player needs to know before this one
// can be shown to them, even in silhouette.
//
// This is what drives the reveal: the minor two five one appears once
// the player owns the minor seventh flat five, the dominant and the
// minor major seventh. Before that it does not exist for them.
//
// # What a mode is made of
//
// Its two tetrachords, the lower from the tonic to the fourth note and
// the upper from the fifth to the octave: that is the order they are
// learnt in, the tetrachords of a system before its modes. One
// component when both are the same shape, the dorian being two minor
// tetrachords.
//
// Every other kind has none yet: nil, which reveals nothing.
func (n Notion) Components() []Notion {
	if n.Kind != KindMode {
		return nil
	}
	p, ok := n.System.Mode(n.Degree)
	if !ok {
		return nil
	}
	split, ok := p.Tetrachords()
	if !ok {
		return nil
	}
	if split.Lower == split.Upper {
		return []Notion{TetrachordOf(split.Lower)}
	}
	return []Notion{TetrachordOf(split.Lower), TetrachordOf(split.Upper)}
}
