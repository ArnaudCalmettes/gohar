package band

import (
	_ "embed"
	"sync/atomic"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
	"github.com/ArnaudCalmettes/gohar/synth"
)

// snapWAV is the finger snap, from newagesoup on Freesound, CC0 (see
// sounds/CREDITS.md). Its attack falls 3 ms into the file, well under
// the spread of the humanizing to come.
//
//go:embed sounds/snap.wav
var snapWAV []byte

// bowlWAV is the singing bowl of the walker's meditation, synthesized
// by sounds/bowl.go (see sounds/CREDITS.md).
//
//go:embed sounds/bowl.wav
var bowlWAV []byte

// walkSF2 is GeneralUser GS slimmed down to what the band plays: the
// double bass, the piano, and the keys of the Jazz kit below (see
// sounds/CREDITS.md and `make slim`).
//
//go:embed sounds/walk.sf2
var walkSF2 []byte

// Lookahead is how far ahead the band is scheduled. Ebiten's loop
// ticks every 16 ms, too coarse to play a beat on time: notes are
// queued a window ahead and the synth applies them to the sample (see
// "Le son" in docs/walk.md).
const Lookahead = 100 * time.Millisecond

// CountIn is the count-in, in beats: two bars, the time to move the
// hands from the space bar to the keyboard.
const CountIn = 2 * mark.BeatsPerBar

// The drums: the Jazz kit of the soundfont (see "Le son" in
// docs/walk.md), the hi-hat closed by the foot, and the ride.
var jazzKit = synth.Preset{Bank: 128, Patch: 32}

const (
	pedalHat = 44
	ride1    = 51 // or ride 2, 59: they differ only in pitch

	// The calibration's "ta, ta, ta, TI": the low wood block, then the
	// high one, of the General MIDI map every kit follows.
	hiWood  = 76
	lowWood = 77

	// The casserole of the lessons, a wrong note's "clong": the cowbell.
	// walk.sf2 also keeps the bass drum, 36, the snare, 38, and the hi-hat
	// closed by the stick, 42, for the plain kit of chapter 2 to come
	// (docs/debutants/chapitre-2.md).
	Casserole    = 56
	CasseroleVel = 0.9
)

var acousticGrand = synth.Preset{Bank: 0, Patch: 0}

// The velocities of the calibration's bar, 0 to 1, set by ear; the
// band's parts have theirs in parts.go.
const (
	taVel = 0.8 // its beats, well in front
	tiVel = 1   // its fourth beat
)

// The gains of the mix. With a soundfont, the balance of the renders
// validated by ear, then the snap lowered: alone in the high range, it
// was all one heard at a low volume;
// meltysynth plays at half scale, hence the large gains of the samplers.
// The kit's is a first guess, to set by ear.
const (
	bassGain  = 3
	pianoGain = 2
	kitGain   = 2
	snapGain  = 0.25
	bowlGain  = 0.3 // a first guess, to set by ear

	chipHandsGain = 0.6
	chipNoiseGain = 0.8
)

// The keys of the chip sounds, the fallback without a soundfont: the
// noise is pitched by the key, high for a short hiss.
const (
	noiseKey = 84
	chipRide = 96
	hold     = 30 * time.Millisecond

	// The bowl of the chip: a long triangle, high.
	chipBowl     = 86
	chipBowlHold = 1500 * time.Millisecond
)

// BowlVel is the singing bowl's velocity.
const BowlVel = 0.8

// A Band is everything the game sounds: the drums, the snaps when the
// walker snaps, the reference bass when it plays, and the player's two
// hands.
//
// Two flavours. With a soundfont, a double bass, a piano and a kit from
// it, and the recorded snap. Without, the chip voices of ear: a
// triangle for the hands, a noise for the drums and the snap. Either
// way the hands sound at once, NoteOn, and the band is scheduled ahead,
// ScheduleOn.
type Band struct {
	kit, snap   synth.Instrument
	bowl        synth.Instrument // the singing bowl, a clip; a triangle on the chip
	bass, Piano synth.Instrument // the player's hands, split at `Split`
	Split       int
	BassHand    atomic.Bool  // the split holds: only while a grid is played
	Silent      atomic.Bool  // the keys do not sound: they turn a lesson's page
	hanging     atomic.Int32 // a key kept sounding once up (see Hold), 0 for none
	hat, ride   int          // the keys of the kit
	ta, ti      int          // the calibration's
	chip        bool         // the kit is a noise: each stroke needs its release

	Demo   bool      // the band plays the bass itself, the reference line
	held   int       // the key the reference bass holds, 0 for none
	heldAt time.Time // when it starts: a lookahead from now, at most
}

// Play schedules `strokes`, the parts of beat `n` of `m`, each at its
// place in the beat: the band only plays them, the parts decide (see
// parts.go).
func (bd *Band) Play(strokes []Stroke, n int, m tempo.Metronome) {
	for _, s := range strokes {
		at := m.AtBeats(float64(n) + s.at)
		switch s.sound {
		case hatSound:
			bd.Strike(bd.hat, s.vel, at)
		case rideSound:
			bd.Strike(bd.ride, s.vel, at)
		case snapSound:
			bd.snap.ScheduleOn(noiseKey, s.vel, at)
			bd.snap.ScheduleOff(noiseKey, at.Add(hold)) // a clip ignores it
		case bassSound:
			bd.walk(s.key, s.vel, at)
		}
	}
}

// Tick schedules a beat of the calibration at `at`: "ta", or "TI" when
// `accent`. The calibration schedules it itself, on its own beats.
func (bd *Band) Tick(at time.Time, accent bool) {
	if accent {
		bd.Strike(bd.ti, tiVel, at)
		return
	}
	bd.Strike(bd.ta, taVel, at)
}

// walk sounds `key` on the reference bass at `vel`, legato, as Siskind
// asks: each note holds until the next one, released on the same date
// it is replaced.
func (bd *Band) walk(key int, vel float64, at time.Time) {
	if bd.held != 0 {
		bd.bass.ScheduleOff(bd.held, at)
	}
	bd.bass.ScheduleOn(key, vel, at)
	bd.held, bd.heldAt = key, at
}

// Strike schedules a stroke of the kit on `key`. A cymbal rings until
// it is struck again: the stroke before is released on the same date,
// as the bass is. The noise of the chip has no ring of its own, and is
// released at once.
func (bd *Band) Strike(key int, vel float64, at time.Time) {
	if bd.chip {
		bd.kit.ScheduleOn(key, vel, at)
		bd.kit.ScheduleOff(key, at.Add(hold))
		return
	}
	bd.kit.ScheduleOff(key, at)
	bd.kit.ScheduleOn(key, vel, at)
}

// SnapAt sounds the walker's snap at `at`, at `vel`, out of any part:
// his sign of a right answer in a lesson.
func (bd *Band) SnapAt(vel float64, at time.Time) {
	bd.snap.ScheduleOn(noiseKey, vel, at)
	bd.snap.ScheduleOff(noiseKey, at.Add(hold)) // a clip ignores it
}

// Bowl strikes the singing bowl at `at`: the walker meditates.
func (bd *Band) Bowl(at time.Time) {
	bd.bowl.ScheduleOn(chipBowl, BowlVel, at)
	bd.bowl.ScheduleOff(chipBowl, at.Add(chipBowlHold)) // a clip ignores it
}

// Stop releases what the reference bass holds, at `at`. The last note
// may be scheduled to start later, within the lookahead, and nothing
// takes back a date given: released at `at`, before it starts, it would
// then ring on alone. It is released once started, `hold` later, a
// pluck barely heard.
func (bd *Band) Stop(at time.Time) {
	if bd.held != 0 {
		bd.bass.ScheduleOff(bd.held, maxTime(at, bd.heldAt.Add(hold)))
		bd.held = 0
	}
}

// maxTime is the later of `a` and `b`.
func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// Key sounds the player's key at once: the piano, and while a grid is
// played, `BassHand`, the double bass below the split; nothing while
// `Silent`, when a lesson waits for a bubble to be read and the key
// turns the page. Called from the MIDI goroutine: the instruments'
// queues are safe for it, and the flags are atomic.
func (bd *Band) Key(e keyboard.Event) {
	if e.Down && bd.Silent.Load() {
		return
	}
	if !e.Down && int32(e.Key) == bd.hanging.Load() {
		return // left hanging, until Unhold
	}
	inst := bd.Piano
	if e.Key < bd.Split && bd.BassHand.Load() {
		inst = bd.bass
	}
	if e.Down {
		inst.NoteOn(e.Key, e.Velocity, e.At)
	} else {
		inst.NoteOff(e.Key, e.At)
	}
}

// Pluck schedules a note of the double bass, `key` at `vel` from `on`
// to `off`, out of any part: the walker's phrases in a lesson on the
// bass.
func (bd *Band) Pluck(key int, vel float64, on, off time.Time) {
	bd.bass.ScheduleOn(key, vel, on)
	bd.bass.ScheduleOff(key, off)
}

// Hold keeps key `k` sounding once the player lets it go, until
// Unhold: a note left hanging in a lesson. The scene asks for it as it
// hears the key go down, within a frame: well before a finger lifts.
func (bd *Band) Hold(k int) { bd.hanging.Store(int32(k)) }

// Unhold lets go, at `at`, of the key Hold kept sounding, if any.
func (bd *Band) Unhold(at time.Time) {
	if k := bd.hanging.Swap(0); k != 0 {
		bd.Piano.ScheduleOff(int(k), at)
	}
}
