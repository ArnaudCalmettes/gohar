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

// The keys of the chip sounds, the fallback without a soundfont.
const (
	noiseKey = 84 // the noise is pitched by the key: high, a short hiss
	hold     = 30 * time.Millisecond
)

// A band is everything the game sounds: the snaps on 2 and 4, the
// reference bass when it plays, and the player's two hands.
//
// Two flavours. With a soundfont, a double bass and a piano from it,
// and the recorded snap. Without, the chip voices of ear: a triangle
// for the hands, a noise for the snap. Either way the hands sound at
// once, NoteOn, and the band is scheduled ahead, ScheduleOn.
type band struct {
	snap        synth.Instrument
	bass, piano synth.Instrument // the player's hands, split at `split`
	split       int

	demo bool // the band plays the bass itself, the reference line
	held int  // the key the reference bass holds, 0 for none
}

// newBand builds the band on `mix`. A nil soundfont gives the chip
// sounds.
func newBand(mix *synth.Mixer, sf *synth.SoundFont, bassPreset synth.Preset, split int, demo bool) (*band, error) {
	b := &band{split: split, demo: demo}
	if sf == nil {
		hands, err := synth.NewEngine("triangle", true)
		if err != nil {
			return nil, err
		}
		noise, err := synth.NewEngine("noise", true)
		if err != nil {
			return nil, err
		}
		b.bass, b.piano, b.snap = hands, hands, noise
		mix.Add(hands, 0.6)
		mix.Add(noise, 0.8)
		return b, nil
	}

	bass, err := synth.NewSampler(sf, bassPreset)
	if err != nil {
		return nil, err
	}
	piano, err := synth.NewSampler(sf, synth.Preset{Bank: 0, Patch: 0})
	if err != nil {
		return nil, err
	}
	snap, err := synth.LoadClip(bytes.NewReader(snapWAV))
	if err != nil {
		return nil, err
	}
	b.bass, b.piano, b.snap = bass, piano, snap
	// The balance of the renders validated by ear, bass to snap as 4 to
	// 0.6, a little lower overall. meltysynth plays at half scale, hence
	// the large gains of the soundfont.
	mix.Add(bass, 3)
	mix.Add(piano, 2)
	mix.Add(snap, 0.45)
	return b, nil
}

// beat schedules what the band plays on beat `n`, at `at`: snaps on 2
// and 4, from the count-in on, and in demo the root of `b` on every
// beat, legato. A beat of the count-in has no chord: `b` is nil.
func (bd *band) beat(p Position, b *Beat, at time.Time) {
	if p.Beat%2 == 0 {
		bd.snap.ScheduleOn(noiseKey, 1, at)
		bd.snap.ScheduleOff(noiseKey, at.Add(hold)) // a clip ignores it
	}
	if !bd.demo || b == nil {
		return
	}
	root := bassKey(b.Chord.Bass)
	// Legato, as Siskind asks: each note holds until the next one,
	// released on the same date it is replaced.
	if bd.held != 0 {
		bd.bass.ScheduleOff(bd.held, at)
	}
	bd.bass.ScheduleOn(root, 0.8, at)
	bd.held = root
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
