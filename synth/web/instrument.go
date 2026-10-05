//go:build js && wasm

package web

import (
	"math"
	"syscall/js"
	"time"

	"github.com/ArnaudCalmettes/gohar/synth"
)

// drumChannel is General MIDI's drum channel, where FluidSynth plays a
// kit.
const drumChannel = 9

// An instrument is a preset of the soundfont on a channel of FluidSynth,
// a synth.Instrument as the desktop's Sampler is. Its notes now go
// straight to FluidSynth; its dated notes go through the sequencer,
// dated relative to now, to the millisecond.
type instrument struct {
	s       *Synth
	channel int
}

// Instrument plays the preset `bank`:`patch` on `channel`, at `volume`,
// from 0 to 1, the channel's share of the master gain. A kit, bank 128,
// takes the drum channel whatever `channel` says.
func (s *Synth) Instrument(bank, patch, channel int, volume float64) synth.Instrument {
	if bank == 128 {
		channel = drumChannel
		s.synth.Call("setChannelType", channel, true)
	}
	s.synth.Call("midiProgramSelect", channel, s.sfont, bank, patch)
	s.synth.Call("midiControl", channel, 7, midi7(volume)) // 7: the channel's volume
	return &instrument{s: s, channel: channel}
}

func (i *instrument) NoteOn(key int, velocity float64, _ time.Time) {
	i.s.synth.Call("midiNoteOn", i.channel, key, midi7(velocity))
}

func (i *instrument) NoteOff(key int, _ time.Time) {
	i.s.synth.Call("midiNoteOff", i.channel, key)
}

func (i *instrument) ScheduleOn(key int, velocity float64, at time.Time) {
	i.s.sendAt(at, map[string]any{"type": "noteon", "channel": i.channel, "key": key, "vel": midi7(velocity)})
}

func (i *instrument) ScheduleOff(key int, at time.Time) {
	i.s.sendAt(at, map[string]any{"type": "noteoff", "channel": i.channel, "key": key})
}

// Read is never called: FluidSynth renders on the audio thread, and
// nothing in Go mixes it.
func (i *instrument) Read([]byte) (int, error) { return 0, nil }

// Delays and Histogram measure nothing yet: the delay of a key is
// FluidSynth's, out of Go's sight.
func (i *instrument) Delays() (last, worst time.Duration) { return 0, 0 }
func (i *instrument) Histogram() synth.Histogram          { return synth.Histogram{} }

// sendAt hands `event` to the sequencer, to play at `at`. The date is
// turned into a delay from now on the main thread; the message reaches
// the audio thread a moment later, the same moment for every note, a
// constant the calibration measures with the rest. A date already past
// plays at once, as on the desktop.
func (s *Synth) sendAt(at time.Time, event map[string]any) {
	ms := max(time.Until(at).Milliseconds(), 0)
	s.seq.Call("sendEventAt", js.ValueOf(event), ms, false)
}

// midi7 turns a level from 0 to 1 into a MIDI value, 1 to 127, as the
// desktop's Sampler does for velocities.
func midi7(v float64) int {
	return int(max(1, min(127, math.Round(v*127))))
}
