package analysis

import (
	"fmt"
	"slices"
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// An Engine turns a stream of key events into a stable reading.
//
// # Why it takes the time instead of reading a clock
//
// The caller drives. A game calls Advance once per frame with the frame
// time, and the engine never sleeps, never ticks and never starts a
// goroutine. There is no clock port here and no fake clock in the
// tests: a test that wants to skip four seconds passes a time four
// seconds later, which is both simpler and more honest than a clock
// interface that exists only to be faked.
//
// The same holds for input. The game reads its MIDI device, or replays
// a fixed sequence in a training mode, and calls NoteOn and NoteOff.
// Which of those it is, the engine has no reason to know.
//
// # Why there are two time constants
//
// A chord changes every second or two; a key changes every thirty. An
// engine with one inertia either lags behind the chords or lets a
// single passing note rewrite the key. So the chord reading and the
// tonality accumulate separately, and the slow one only moves on
// weight of evidence.
//
// The Engine is not safe for concurrent use. A game feeding it from an
// input goroutine while reading it from the render loop must guard it,
// or better, feed and read it from the same loop.
type Engine struct {
	recognizer *Recognizer
	config     Config

	// sounding holds the keys currently down; released holds the ones
	// let go, with the moment they were, so that Gather can expire
	// them.
	sounding map[harmony.Pitch]bool
	released map[harmony.Pitch]time.Time

	reading Reading
	hasRead bool

	// pending is the reading waiting out Settle before it replaces the
	// current one.
	pending      Reading
	hasPending   bool
	pendingSince time.Time

	tonality harmony.Tonality

	// pinned holds a key the caller fixed. While it is set the
	// accumulator keeps running but does not change what is reported.
	pinned   harmony.Tonality
	isPinned bool

	// history is the sequence of readings adopted, in order, with no
	// immediate repeat. It is what [Engine.Progression] turns into a
	// shape.
	history []harmony.PitchClass
	tetrads []harmony.ChordPattern

	// events is what the tonality is inferred from: every class struck,
	// with when. Pruned to TonalityWindow on each Advance.
	events []event

	// lastTime clamps out of order timestamps forward rather than
	// dropping the note they carry.
	lastTime time.Time

	scratch []harmony.Pitch
}

type event struct {
	class harmony.PitchClass
	at    time.Time
}

// Config tunes the temporal behaviour.
//
// These are data because there is no correct value, only values that
// feel right at the keyboard, and no test can validate them: an
// assertion can prove the engine waits two hundred milliseconds, it
// cannot prove two hundred is the right number. Expect to change them
// by playing, not by reasoning.
type Config struct {
	// Gather is how long a note keeps counting toward the current
	// reading after it stops sounding.
	//
	// This is what lets an arpeggio read as a chord. Too short and a
	// rolled voicing never forms one; too long and the reading trails
	// behind the player. Somewhere near half a second is a starting
	// point.
	Gather time.Duration

	// Settle is how long a new reading must stay on top before it
	// replaces the current one.
	//
	// Half the hysteresis, and the half that guards against a passing
	// note. The other half needs no number: a reading that is still
	// valid is simply kept, so the four roots of a diminished seventh chord
	// cannot trade places while the chord is held. See [Context.Current].
	Settle time.Duration

	// TonalityWindow is how long the tonality accumulator remembers.
	//
	// An order of magnitude above Gather, and that ratio matters more
	// than either value. A key that moves as fast as a chord is not a
	// key.
	TonalityWindow time.Duration

	// TonalityMargin is how far ahead a candidate key must be before
	// it replaces the current one.
	//
	// High on purpose. An accidental note is not a modulation, and the
	// cost of switching late is a few seconds of stale spelling, while
	// the cost of switching early is a display that lies.
	TonalityMargin int
}

// DefaultConfig is a starting point, not an answer. The ratios between
// the values are the part worth keeping.
var DefaultConfig = Config{
	Gather:         400 * time.Millisecond,
	Settle:         150 * time.Millisecond,
	TonalityWindow: 12 * time.Second,
	TonalityMargin: 8,
}

// NewEngine builds an engine over a recognizer.
func NewEngine(r *Recognizer, c Config) (*Engine, error) {
	if r == nil {
		return nil, fmt.Errorf("analysis: an engine needs a recognizer")
	}
	if c.Gather <= 0 || c.Settle < 0 || c.TonalityWindow <= 0 {
		return nil, fmt.Errorf("analysis: config has a non-positive duration")
	}
	return &Engine{
		recognizer: r,
		config:     c,
		sounding:   make(map[harmony.Pitch]bool),
		released:   make(map[harmony.Pitch]time.Time),
	}, nil
}

// NoteOn records that a pitch started sounding.
//
// Out of order timestamps are tolerated and clamped rather than
// rejected: a MIDI device that stamps events on its own schedule will
// occasionally hand over a note that predates the last one, and
// dropping it would silently lose a note the player heard.
func (e *Engine) NoteOn(p harmony.Pitch, at time.Time) {
	at = e.clamp(at)
	e.sounding[p] = true
	delete(e.released, p)
	e.events = append(e.events, event{class: p.Class(), at: at})
}

// NoteOff records that a pitch stopped sounding.
//
// The pitch keeps counting toward the reading for [Config.Gather]
// after this, which is what lets an arpeggio form a chord.
//
// A note off for a pitch that was never on is ignored. Devices send
// them after a reconnect, and treating it as an error would take down
// a game over something inaudible.
func (e *Engine) NoteOff(p harmony.Pitch, at time.Time) {
	at = e.clamp(at)
	if !e.sounding[p] {
		return
	}
	delete(e.sounding, p)
	e.released[p] = at
}

// Advance moves the engine to `now`: expires what has fallen out of the
// window, reads the chord and the key again, and applies the
// hysteresis.
//
// Idempotent for a given time, so a caller may advance twice on one
// frame without changing anything.
func (e *Engine) Advance(now time.Time) {
	now = e.clamp(now)

	for p, off := range e.released {
		if now.Sub(off) > e.config.Gather {
			delete(e.released, p)
		}
	}

	cutoff := now.Add(-e.config.TonalityWindow)
	kept := e.events[:0]
	for _, ev := range e.events {
		if !ev.at.Before(cutoff) {
			kept = append(kept, ev)
		}
	}
	e.events = kept

	e.inferTonality()

	best, ok := e.recognizer.Best(e.Snapshot(), Context{
		Tonality: e.tonality,
		Current:  e.reading,
	})
	if !ok {
		e.reading, e.hasRead = Reading{}, false
		e.pending, e.hasPending = Reading{}, false
		return
	}

	if !e.hasRead {
		e.adopt(best)
		e.pending, e.hasPending = Reading{}, false
		return
	}
	if sameReading(best, e.reading) {
		e.pending, e.hasPending = Reading{}, false
		return
	}

	if !e.hasPending || !sameReading(best, e.pending) {
		e.pending, e.hasPending, e.pendingSince = best, true, now
		return
	}
	if now.Sub(e.pendingSince) >= e.config.Settle {
		e.adopt(best)
		e.pending, e.hasPending = Reading{}, false
	}
}

// Reading returns the current stable reading, and whether there is one.
//
// Stable is the point: this is what survived [Config.Settle] and the
// stickiness rule, not the first entry of the last identification. A
// caller that wants every reading calls the recognizer directly on
// [Engine.Snapshot], passing the current one as [Context.Current].
func (e *Engine) Reading() (Reading, bool) {
	return e.reading, e.hasRead
}

// Snapshot returns what the engine currently considers to be sounding,
// gathered notes included.
//
// The returned slice is reused between calls. Copy it to keep it.
func (e *Engine) Snapshot() Snapshot {
	e.scratch = e.scratch[:0]
	for p := range e.sounding {
		e.scratch = append(e.scratch, p)
	}
	for p := range e.released {
		e.scratch = append(e.scratch, p)
	}
	slices.Sort(e.scratch)
	return Snapshot{Pitches: e.scratch}
}

// Tonality returns the key the engine believes it is hearing, or the
// zero value when it has inferred none.
//
// The zero value is not C major and a caller must not render it as a
// key. A player working outside any key, and an engine that has not
// heard enough yet, both land here, and telling them apart is not the
// engine's job: neither is in a key, and that is the honest answer to
// display.
func (e *Engine) Tonality() harmony.Tonality {
	if e.isPinned {
		return e.pinned
	}
	return e.tonality
}

// PinTonality fixes the key and stops inferring it.
//
// # Why both regimes exist
//
// An exercise usually sets the key itself. Asked to play a dorian
// starting where they like, the player fixes the tonic with their
// first note, and there is nothing for the engine to guess: inferring
// it would only add a delay and a chance of getting it wrong.
//
// Free play is the other regime, where nobody has said what the key is
// and the engine has to accumulate evidence. That is the one the
// inference exists for, and probably not the more common one in a
// teaching game.
//
// While pinned, the tonality accumulator still runs but does not
// change the reported key. So releasing gives an answer immediately
// rather than starting from nothing.
func (e *Engine) PinTonality(t harmony.Tonality) {
	e.pinned, e.isPinned = t, true
}

// ReleaseTonality resumes inference from whatever evidence has
// accumulated while pinned.
func (e *Engine) ReleaseTonality() {
	e.pinned, e.isPinned = harmony.Tonality{}, false
}

// TonalityIsPinned reports whether the key was set by the caller
// rather than inferred.
//
// An interface showing a key should be able to say which: told and
// inferred are different claims, and a player deserves to know when
// the program is guessing.
func (e *Engine) TonalityIsPinned() bool {
	return e.isPinned
}

// Progression returns the chords read so far as a relative
// [harmony.Progression], ready to compare against an expected shape.
//
// Relative, so a player who transposed the exercise still matches it.
// The key they chose is not lost: it is the root of the first reading,
// and [harmony.Progression.Compare] returns the distance from whatever
// key was expected.
func (e *Engine) Progression() (harmony.Progression, bool) {
	if len(e.history) == 0 {
		return harmony.Progression{}, false
	}
	p, err := harmony.NewProgression(e.tetrads, e.history)
	if err != nil {
		return harmony.Progression{}, false
	}
	return p, true
}

// Reset clears everything: held notes, reading, accumulated evidence,
// the chord history and any pinned key.
//
// For a game moving between exercises, where carrying the key inferred
// from the last one into the next would be worse than starting blind.
func (e *Engine) Reset() {
	clear(e.sounding)
	clear(e.released)
	e.reading, e.hasRead = Reading{}, false
	e.pending, e.hasPending = Reading{}, false
	e.tonality = harmony.Tonality{}
	e.pinned, e.isPinned = harmony.Tonality{}, false
	e.history = e.history[:0]
	e.tetrads = e.tetrads[:0]
	e.events = e.events[:0]
	e.scratch = e.scratch[:0]
}

// clamp moves a timestamp forward to the last one seen.
//
// Devices stamp events on their own schedule and occasionally hand over
// one that predates the last. Dropping it would lose a note the player
// heard, which is worse than placing it a few milliseconds late.
func (e *Engine) clamp(at time.Time) time.Time {
	if at.Before(e.lastTime) {
		return e.lastTime
	}
	e.lastTime = at
	return at
}

// inferTonality picks the key that best explains what has been struck
// inside the window.
//
// A count, not a classification: each event adds one for a key that
// contains its class and takes one away for a key that does not. This
// is the key, which really is a matter of evidence; chords are
// recognised by equality, see [Recognizer]. The best key has to clear
// TonalityMargin outright before anything is claimed, which is what
// keeps a single chord from being called a key, and it has to beat the
// current one by that margin again before it replaces it.
//
// # A deliberate limitation
//
// Only major patterns are considered. A natural minor scale holds the
// same classes as its relative major, so admitting both would make
// every inference a tie between two answers that the classes alone
// cannot separate. Telling them apart needs what the player does with
// the notes, cadences above all, and that is not in this counting.
//
// The consequence is that a piece in A minor is reported as C major.
// The spelling that follows is right either way, since the two share
// their notes; only the name is wrong, so an interface that shows the
// key should say so carefully until this is improved.
func (e *Engine) inferTonality() {
	if len(e.events) == 0 {
		return
	}

	best, bestScore := harmony.Tonality{}, 0
	for tonic := harmony.PitchClass(0); tonic < harmony.PitchClassCount; tonic++ {
		candidate, err := harmony.NewTonality(tonic, harmony.ScaleMajor)
		if err != nil {
			continue
		}
		score := 0
		for _, ev := range e.events {
			if candidate.Contains(ev.class) {
				score++
			} else {
				score--
			}
		}
		if score > bestScore {
			best, bestScore = candidate, score
		}
	}

	if bestScore <= e.config.TonalityMargin {
		return
	}
	if e.tonality.IsZero() {
		e.tonality = best
		return
	}

	current := 0
	for _, ev := range e.events {
		if e.tonality.Contains(ev.class) {
			current++
		} else {
			current--
		}
	}
	if bestScore > current+e.config.TonalityMargin {
		e.tonality = best
	}
}

// adopt makes a reading current and records it.
//
// A chord held across several frames is one chord, so an adoption that
// repeats the last entry does not extend the history. Without that, a
// player holding a tonic for four bars would have played a four chord
// progression.
func (e *Engine) adopt(r Reading) {
	e.reading, e.hasRead = r, true

	n := len(e.history)
	if n > 0 && e.history[n-1] == r.Root && e.tetrads[n-1] == r.Tetrad {
		return
	}
	e.history = append(e.history, r.Root)
	e.tetrads = append(e.tetrads, r.Tetrad)
}
