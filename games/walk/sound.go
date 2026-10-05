package main

import (
	_ "embed"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/synth"
)

// snapWAV is the finger snap, from newagesoup on Freesound, CC0 (see
// sounds/CREDITS.md). Its attack falls 3 ms into the file, well under
// the spread of the humanizing to come.
//
//go:embed sounds/snap.wav
var snapWAV []byte

// walkSF2 is GeneralUser GS slimmed down to what the band plays: the
// double bass, the piano, and the four keys of the Jazz kit below (see
// sounds/CREDITS.md and `make slim`).
//
//go:embed sounds/walk.sf2
var walkSF2 []byte

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

	chipHandsGain = 0.6
	chipNoiseGain = 0.8
)

// The keys of the chip sounds, the fallback without a soundfont: the
// noise is pitched by the key, high for a short hiss.
const (
	noiseKey = 84
	chipRide = 96
	hold     = 30 * time.Millisecond
)

// A band is everything the game sounds: the drums, the snaps when the
// walker snaps, the reference bass when it plays, and the player's two
// hands.
//
// Two flavours. With a soundfont, a double bass, a piano and a kit from
// it, and the recorded snap. Without, the chip voices of ear: a
// triangle for the hands, a noise for the drums and the snap. Either
// way the hands sound at once, NoteOn, and the band is scheduled ahead,
// ScheduleOn.
type band struct {
	kit, snap   synth.Instrument
	bass, piano synth.Instrument // the player's hands, split at `split`
	split       int
	hat, ride   int  // the keys of the kit
	ta, ti      int  // the calibration's
	chip        bool // the kit is a noise: each stroke needs its release

	demo bool // the band plays the bass itself, the reference line
	held int  // the key the reference bass holds, 0 for none
}

// play schedules `strokes`, the parts of beat `n` of `m`, each at its
// place in the beat: the band only plays them, the parts decide (see
// parts.go).
func (bd *band) play(strokes []stroke, n int, m tempo.Metronome) {
	for _, s := range strokes {
		at := m.AtBeats(float64(n) + s.at)
		switch s.sound {
		case hatSound:
			bd.strike(bd.hat, s.vel, at)
		case rideSound:
			bd.strike(bd.ride, s.vel, at)
		case snapSound:
			bd.snap.ScheduleOn(noiseKey, s.vel, at)
			bd.snap.ScheduleOff(noiseKey, at.Add(hold)) // a clip ignores it
		case bassSound:
			bd.walk(s.key, s.vel, at)
		}
	}
}

// tick schedules a beat of the calibration at `at`: "ta", or "TI" when
// `accent`. The calibration schedules it itself, on its own beats.
func (bd *band) tick(at time.Time, accent bool) {
	if accent {
		bd.strike(bd.ti, tiVel, at)
		return
	}
	bd.strike(bd.ta, taVel, at)
}

// walk sounds `key` on the reference bass at `vel`, legato, as Siskind
// asks: each note holds until the next one, released on the same date
// it is replaced.
func (bd *band) walk(key int, vel float64, at time.Time) {
	if bd.held != 0 {
		bd.bass.ScheduleOff(bd.held, at)
	}
	bd.bass.ScheduleOn(key, vel, at)
	bd.held = key
}

// strike schedules a stroke of the kit on `key`. A cymbal rings until
// it is struck again: the stroke before is released on the same date,
// as the bass is. The noise of the chip has no ring of its own, and is
// released at once.
func (bd *band) strike(key int, vel float64, at time.Time) {
	if bd.chip {
		bd.kit.ScheduleOn(key, vel, at)
		bd.kit.ScheduleOff(key, at.Add(hold))
		return
	}
	bd.kit.ScheduleOff(key, at)
	bd.kit.ScheduleOn(key, vel, at)
}

// stop releases what the reference bass holds, at `at`.
func (bd *band) stop(at time.Time) {
	if bd.held != 0 {
		bd.bass.ScheduleOff(bd.held, at)
		bd.held = 0
	}
}

// key sounds the player's key at once: the double bass below the
// split, the piano above. Called from the MIDI goroutine: the
// instruments' queues are safe for it.
func (bd *band) key(e keyboard.Event) {
	inst := bd.piano
	if e.Key < bd.split {
		inst = bd.bass
	}
	if e.Down {
		inst.NoteOn(e.Key, e.Velocity, e.At)
	} else {
		inst.NoteOff(e.Key, e.At)
	}
}
