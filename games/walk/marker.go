package main

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// The marker is the second brick of the game (see "Les briques" in
// docs/walk.md): notes in, marks out, with no clock of its own. It
// takes note of what was played and never grades the player.
//
// Each note gets two marks, one for its timing and one for its pitch,
// against the beat it is nearest to. Each beat where a chord arrives
// gets a mark of its own once its window has closed: the root landed,
// or it did not.

// A Note is a key pressed in the bass zone, at the instant it sounded:
// the MIDI timestamp, already corrected for the calibrated latency.
type Note struct {
	Key int
	At  time.Time
}

// Timing says how close to its beat a note fell.
type Timing int

const (
	OnTime  Timing = iota
	Early          // within the loose window, before the beat
	Late           // within the loose window, after it
	Between        // outside both windows: no beat claims the note
)

// Pitch says what a note is to the chord sounding on its beat.
type Pitch int

const (
	Root      Pitch = iota // the bass of the change: its root, or the written bass of a slash chord
	ChordTone              // another note of the chord
	Outside                // not in the chord
)

// A NoteMark is what a note was. A note Between beats has no pitch
// mark: there is no chord to measure it against.
type NoteMark struct {
	Note
	Beat   int           // the nearest beat
	Off    time.Duration // negative when early
	Timing Timing
	Pitch  Pitch
}

// BeatKind says what became of a beat.
type BeatKind int

const (
	// Landed: the chord arrived and one note claimed the beat, its root,
	// on time or within the loose window.
	Landed BeatKind = iota

	// Missed: the chord arrived and no root claimed the beat.
	Missed

	// Doubled: two notes or more claimed the beat. On an arrival, it
	// counts as missed: the player is to sound the one note meant, and
	// faced with two, the marker refuses to guess which.
	Doubled
)

// A BeatMark is what became of a beat, once its window closed.
type BeatMark struct {
	Beat int
	Kind BeatKind
}

// Rules are what the marker expects. Data, not code: a level or a
// constraint of a run is a different value (see "Ce qu'on attend, temps
// par temps" in docs/walk.md).
type Rules struct {
	// OnTime and Loose are the half widths of the windows around a beat.
	// The defaults are starting points, to be set by ear.
	OnTime, Loose time.Duration

	// Split is the lowest key of the right hand: below it, the bass.
	Split int

	// Inversions lets another chord tone land an arrival, for a player
	// who knows what to do with it. Off for a beginner, who plays roots.
	Inversions bool
}

// FirstPalier is the first palier: roots, on the changes, nothing else.
var FirstPalier = Rules{
	OnTime: 30 * time.Millisecond,
	Loose:  100 * time.Millisecond,
	Split:  55, // G3, sol2 in French
}

// A Marker marks the notes of one run.
type Marker struct {
	rules Rules
	m     Metronome
	beats []Beat

	claims []claim // per beat: the notes that claimed it
	closed int     // beats before this one have their mark
}

type claim struct {
	notes  int
	landed bool // a note that may land the beat was among them
}

// NewMarker marks a run over `beats`, laid out for `m`.
func NewMarker(r Rules, m Metronome, beats []Beat) *Marker {
	return &Marker{rules: r, m: m, beats: beats, claims: make([]claim, len(beats))}
}

// Play marks a note. False for a note the marker ignores: in the right
// hand's zone, or outside the run, in the count-in or after the end.
func (k *Marker) Play(n Note) (NoteMark, bool) {
	if n.Key >= k.rules.Split {
		return NoteMark{}, false
	}
	beat, off := k.m.Nearest(n.At)
	if beat < 0 || beat >= len(k.beats) {
		return NoteMark{}, false
	}
	mark := NoteMark{Note: n, Beat: beat, Off: off, Timing: k.timing(off)}
	if mark.Timing == Between {
		return mark, true
	}

	b := k.beats[beat]
	pc := harmony.PitchClass(n.Key % 12)
	switch {
	case !b.Chord.Silent && pc == b.Chord.Bass:
		mark.Pitch = Root
	case !b.Chord.Silent && b.Chord.Chord.Set().Contains(pc):
		mark.Pitch = ChordTone
	default:
		mark.Pitch = Outside
	}

	c := &k.claims[beat]
	c.notes++
	if mark.Pitch == Root || mark.Pitch == ChordTone && k.rules.Inversions {
		c.landed = true
	}
	return mark, true
}

// Close marks the beats whose window has closed by `now`, in order. A
// beat is only marked when something happened to it: an arrival, or
// two notes where one was meant.
func (k *Marker) Close(now time.Time) []BeatMark {
	var marks []BeatMark
	for ; k.closed < len(k.beats); k.closed++ {
		n := k.closed
		if !now.After(k.m.At(n).Add(k.rules.Loose)) {
			break
		}
		c, b := k.claims[n], k.beats[n]
		switch {
		case c.notes > 1:
			marks = append(marks, BeatMark{n, Doubled})
		case !b.Arrives || b.Chord.Silent:
		case c.landed:
			marks = append(marks, BeatMark{n, Landed})
		default:
			marks = append(marks, BeatMark{n, Missed})
		}
	}
	return marks
}

func (k *Marker) timing(off time.Duration) Timing {
	a := off
	if a < 0 {
		a = -a
	}
	switch {
	case a <= k.rules.OnTime:
		return OnTime
	case a > k.rules.Loose:
		return Between
	case off < 0:
		return Early
	}
	return Late
}
