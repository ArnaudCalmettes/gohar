package mark

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The marker is the second brick of the game (see "Les briques" in
// docs/walk.md): notes in, marks out, with no clock of its own. It
// takes note of what was played and never grades the player.
//
// Each note gets two marks, one for its timing and one for its pitch,
// against the beat it is nearest to. Each beat where a chord arrives,
// and beat 1 of each bar where it carries on, gets a mark of its own
// once its window has closed: a note landed it, or none did.

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
	// Landed: one note claimed the beat, on time or within the loose
	// window, and the rules accept it there: the root of a chord that
	// arrives; any note of one that carries on.
	Landed BeatKind = iota

	// Missed: a note was expected and none landed the beat.
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
	// OnTime and Loose are the half widths of the windows around a beat,
	// in fractions of a beat, so that they widen as the tempo slows. The
	// defaults are generous starting points, to be set by ear. The notes
	// come corrected for the calibrated latency (see Note).
	OnTime, Loose float64

	// Split is the lowest key of the right hand: below it, the bass.
	Split int

	// Inversions lets another chord tone land an arrival, for a player
	// who knows what to do with it. Off for a beginner, who plays roots.
	Inversions bool

	// HeldChordTone lets any note of the chord land beat 1 of a bar
	// where the chord carries on, besides its root: the bass moves
	// within the chord it holds.
	HeldChordTone bool
}

// FirstPalier is the first palier: roots, on the changes, nothing else.
var FirstPalier = Rules{
	OnTime: 1.0 / 8, // 75 ms at 100
	Loose:  1.0 / 3, // 200 ms at 100
	Split:  55,      // G3, sol2 in French

	HeldChordTone: true,
}

// Palier0 is the palier of the beginners, without tempo (see
// docs/debutants/chapitre-1.md): the root of every bar, a bar repeated
// included.
var Palier0 = Rules{Split: FirstPalier.Split}

// A Marker marks the notes of one run.
type Marker struct {
	rules Rules
	m     tempo.Metronome
	beats []Beat

	claims []claim // per beat: the notes that claimed it
	closed int     // beats before this one have their mark
}

type claim struct {
	notes  int
	Landed bool // a note that may land the beat was among them
}

// NewMarker marks a run over `beats`, laid out for `m`.
func NewMarker(r Rules, m tempo.Metronome, beats []Beat) *Marker {
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
	mark.Pitch = pitchOf(n.Key, b.Chord)
	c := &k.claims[beat]
	c.notes++
	if k.rules.lands(b, mark.Pitch) {
		c.Landed = true
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
		if !now.After(k.m.At(n).Add(k.window(k.rules.Loose))) {
			break
		}
		c, b := k.claims[n], k.beats[n]
		switch {
		case c.notes > 1:
			marks = append(marks, BeatMark{n, Doubled})
		case !b.Arrives && !b.Holds || b.Chord.Silent:
		case c.Landed:
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
	case a <= k.window(k.rules.OnTime):
		return OnTime
	case a > k.window(k.rules.Loose):
		return Between
	case off < 0:
		return Early
	}
	return Late
}

// pitchOf says what `key` is to the chord `ch`: its bass, another of
// octave is the MIDI keys from a note to the same, an octave up.
const octave = 12

// its notes, or neither. Nothing is in a silence.
func pitchOf(key int, ch analysis.Change) Pitch {
	pc := harmony.PitchClass(key % octave)
	switch {
	case ch.Silent:
		return Outside
	case pc == ch.Bass:
		return Root
	case ch.Chord.Set().Contains(pc):
		return ChordTone
	}
	return Outside
}

// lands tells whether a note of pitch `p` lands beat `b`: its root
// always; another chord tone when inversions are allowed, or on beat 1
// of a held chord when the rules allow it.
func (r Rules) lands(b Beat, p Pitch) bool {
	switch {
	case p == Root:
		return true
	case p != ChordTone:
		return false
	}
	return r.Inversions || b.Holds && r.HeldChordTone
}

// window turns a fraction of a beat into a duration at the run's tempo.
func (k *Marker) window(f float64) time.Duration {
	return time.Duration(f * float64(k.m.Beat()))
}
