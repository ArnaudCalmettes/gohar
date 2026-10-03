package main

import (
	"bytes"
	_ "embed"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/synth"
)

// snapWAV is the finger snap, from newagesoup on Freesound, CC0 (see
// sounds/CREDITS.md). Its attack falls 3 ms into the file, well under
// the spread of the humanizing to come.
//
//go:embed sounds/snap.wav
var snapWAV []byte

// The drums: the Jazz kit of the soundfont (see "Le son" in
// docs/walk.md), the hi-hat closed by the foot, and the ride.
var jazzKit = synth.Preset{Bank: 128, Patch: 32}

const (
	pedalHat = 44
	ride1    = 51 // or ride 2, 59: they differ only in pitch
)

var acousticGrand = synth.Preset{Bank: 0, Patch: 0}

// The velocities of the band, 0 to 1, set by ear.
const (
	countVel    = 0.7 // the hi-hat of the count-in
	backbeatVel = 0.7 // the ride on 2 and 4
	beatVel     = 0.4 // the ride on 1 and 3
	andVel      = 0.5 // the ride on the "and" of 2 and 4
	hatVel      = 0.6 // the hi-hat on 2 and 4
	endVel      = 0.7 // the ride on the demo's last note
	snapVel     = 1
	bassVel     = 0.8 // the reference bass
)

// The gains of the mix. With a soundfont, the balance of the renders
// validated by ear, bass to snap as 4 to 0.6, a little lower overall;
// meltysynth plays at half scale, hence the large gains of the samplers.
// The kit's is a first guess, to set by ear.
const (
	bassGain  = 3
	pianoGain = 2
	kitGain   = 2
	snapGain  = 0.45

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
	chip        bool // the kit is a noise: each stroke needs its release

	demo bool // the band plays the bass itself, the reference line
	held int  // the key the reference bass holds, 0 for none
}

// newBand builds the band on `mix`. A nil soundfont gives the chip
// sounds.
func newBand(mix *synth.Mixer, sf *synth.SoundFont, bassPreset synth.Preset, split int, demo bool) (*band, error) {
	b := &band{split: split, demo: demo, hat: pedalHat, ride: ride1}
	if sf == nil {
		hands, err := synth.NewEngine("triangle", true)
		if err != nil {
			return nil, err
		}
		noise, err := synth.NewEngine("noise", true)
		if err != nil {
			return nil, err
		}
		b.bass, b.piano, b.snap, b.kit = hands, hands, noise, noise
		b.hat, b.ride, b.chip = noiseKey, chipRide, true
		mix.Add(hands, chipHandsGain)
		mix.Add(noise, chipNoiseGain)
		return b, nil
	}

	bass, err := synth.NewSampler(sf, bassPreset)
	if err != nil {
		return nil, err
	}
	piano, err := synth.NewSampler(sf, acousticGrand)
	if err != nil {
		return nil, err
	}
	kit, err := synth.NewSampler(sf, jazzKit)
	if err != nil {
		return nil, err
	}
	snap, err := synth.LoadClip(bytes.NewReader(snapWAV))
	if err != nil {
		return nil, err
	}
	b.bass, b.piano, b.kit, b.snap = bass, piano, kit, snap
	mix.Add(bass, bassGain)
	mix.Add(piano, pianoGain)
	mix.Add(kit, kitGain)
	mix.Add(snap, snapGain)
	return b, nil
}

// beat schedules what the band plays on the beat at `p`, at `at`, `and`
// being the swung "and" after it:
//   - in the count-in, the hi-hat on 1 and 3 of its first bar, then on
//     every beat of the second: "1, 3, 1, 2, 3, 4";
//   - from bar 1, the ride on every beat and on the "and" of 2 and 4,
//     the accents on 2 and 4, with the hi-hat on 2 and 4: the
//     metronome;
//   - a finger snap on 2 and 4 when `snap`, the walker snapping his
//     fingers: the sound of the juice;
//   - in demo, `key`, the reference line, legato. No key, 0, in the
//     count-in or out of demo.
func (bd *band) beat(p Position, key int, at, and time.Time, snap bool) {
	backbeat := p.Beat%2 == 0
	switch {
	case p.Bar < 0 && backbeat:
		// the first bar of the count-in: 1 and 3 only
	case p.Bar <= 0:
		bd.strike(bd.hat, countVel, at)
	default:
		vel := beatVel
		if backbeat {
			vel = backbeatVel
			bd.strike(bd.hat, hatVel, at)
			bd.strike(bd.ride, andVel, and)
		}
		bd.strike(bd.ride, vel, at)
	}
	if snap && backbeat {
		bd.snap.ScheduleOn(noiseKey, snapVel, at)
		bd.snap.ScheduleOff(noiseKey, at.Add(hold)) // a clip ignores it
	}
	if bd.demo && key != 0 {
		bd.walk(key, at)
	}
}

// end plays the demo's last note, `key` at `at`, with a stroke of the
// ride: the chord the turnaround leads to (see Ending).
func (bd *band) end(key int, at time.Time) {
	bd.strike(bd.ride, endVel, at)
	bd.walk(key, at)
}

// walk sounds `key` on the reference bass, legato, as Siskind asks: each
// note holds until the next one, released on the same date it is
// replaced.
func (bd *band) walk(key int, at time.Time) {
	if bd.held != 0 {
		bd.bass.ScheduleOff(bd.held, at)
	}
	bd.bass.ScheduleOn(key, bassVel, at)
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
