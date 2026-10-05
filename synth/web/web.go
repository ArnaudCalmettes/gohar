//go:build js && wasm

// Package web plays the synth's instruments in the browser, on the audio
// thread, rather than on the main thread where Go runs.
//
// On the desktop, synth computes the samples in Go and hands them to
// oto. In the browser, oto takes them from the main thread, which also
// draws the game and runs Go: the reserves it needs to ride over a frame
// add up to a tenth of a second, too late to play. Here FluidSynth,
// compiled to WebAssembly by js-synthesizer, runs in an AudioWorklet and
// renders the soundfont on the audio thread; Go only sends it the notes,
// at once for the player's keys, dated through FluidSynth's sequencer
// for the band's. A recorded sound is a Web Audio buffer, started at its
// date by the browser itself.
//
// The page loads js-synthesizer.min.js, and serves libfluidsynth and
// js-synthesizer's worklet next to it (see make wasm). Everything here
// talks to them through syscall/js, and nothing else in the repository
// does.
package web

import (
	"errors"
	"fmt"
	"syscall/js"
)

// seventhOrder is FluidSynth's finest interpolation of the samples, its
// fourth order being the default; cheap on the audio thread.
const seventhOrder = 7

// The scripts the AudioWorklet loads, served next to the page.
const (
	fluidSynthScript = "libfluidsynth-2.4.6.js"
	workletScript    = "js-synthesizer.worklet.min.js"
)

// Options tune the synthesiser.
type Options struct {
	// Gain is FluidSynth's master gain; zero means DefaultGain. The
	// instruments' volumes are relative to it.
	Gain float64
}

// DefaultGain is a starting point, to set by ear against the desktop's
// mix.
const DefaultGain = 0.5

// A Synth is FluidSynth on the audio thread, one soundfont loaded, and
// the audio context the recorded sounds share with it.
type Synth struct {
	ctx   js.Value // the AudioContext
	synth js.Value // js-synthesizer's AudioWorkletNodeSynthesizer
	seq   js.Value // its sequencer, for the dated notes
	sfont int      // the soundfont's id in FluidSynth
}

// Open starts the audio context, loads FluidSynth into its worklet and
// the soundfont `sf2` into FluidSynth. It blocks until all of it is
// ready: call it from a goroutine that may wait, main, never from a
// JavaScript callback.
//
// The browser starts no sound before a gesture of the player: the page
// starts the program on a click, which is that gesture.
func Open(sf2 []byte, opts Options) (*Synth, error) {
	jssynth := js.Global().Get("JSSynth")
	if !jssynth.Truthy() {
		return nil, errors.New("web: JSSynth not found, the page must load js-synthesizer.min.js")
	}
	if opts.Gain == 0 {
		opts.Gain = DefaultGain
	}
	ctx := js.Global().Get("AudioContext").New(map[string]any{"latencyHint": "interactive"})
	worklet := ctx.Get("audioWorklet")
	for _, script := range []string{fluidSynthScript, workletScript} {
		if _, err := await(worklet.Call("addModule", script)); err != nil {
			return nil, fmt.Errorf("web: loading %s: %w", script, err)
		}
	}
	s := &Synth{ctx: ctx, synth: jssynth.Get("AudioWorkletNodeSynthesizer").New()}
	s.synth.Call("init", ctx.Get("sampleRate"))
	// No reverb, no chorus, as on the desktop (see "Ni réverbération ni
	// chorus" in docs/architecture.md): FluidSynth turns both on by
	// default, and their halo blurs the bass and the kit. The settings go
	// to the node: the worklet's synthesiser is created with it.
	settings := map[string]any{"reverbActive": false, "chorusActive": false}
	s.synth.Call("createAudioNode", ctx, settings).Call("connect", ctx.Get("destination"))
	s.synth.Call("setGain", opts.Gain)
	s.synth.Call("setInterpolation", seventhOrder)

	buf := js.Global().Get("Uint8Array").New(len(sf2))
	js.CopyBytesToJS(buf, sf2)
	id, err := await(s.synth.Call("loadSFont", buf.Get("buffer")))
	if err != nil {
		return nil, fmt.Errorf("web: loading the soundfont: %w", err)
	}
	s.sfont = id.Int()

	if s.seq, err = await(s.synth.Call("createSequencer")); err != nil {
		return nil, fmt.Errorf("web: creating the sequencer: %w", err)
	}
	if _, err := await(s.seq.Call("registerSynthesizer", s.synth)); err != nil {
		return nil, fmt.Errorf("web: registering the synthesiser: %w", err)
	}
	return s, nil
}

// Latency returns what the browser says of its own delays, in seconds:
// the audio context's, and the output's when it knows it.
func (s *Synth) Latency() (base, output float64) {
	base = s.ctx.Get("baseLatency").Float()
	if o := s.ctx.Get("outputLatency"); o.Type() == js.TypeNumber {
		output = o.Float()
	}
	return base, output
}

// Close stops the audio context.
func (s *Synth) Close() error {
	s.ctx.Call("close")
	return nil
}

// await waits for the promise `p` and returns what it resolves with, or
// what it rejects with as an error.
func await(p js.Value) (js.Value, error) {
	done := make(chan js.Value, 1)
	failed := make(chan error, 1)
	then := js.FuncOf(func(_ js.Value, args []js.Value) any {
		v := js.Undefined()
		if len(args) > 0 {
			v = args[0]
		}
		done <- v
		return nil
	})
	defer then.Release()
	catch := js.FuncOf(func(_ js.Value, args []js.Value) any {
		msg := "rejected"
		if len(args) > 0 {
			msg = args[0].Call("toString").String()
		}
		failed <- errors.New(msg)
		return nil
	})
	defer catch.Release()
	p.Call("then", then, catch)
	select {
	case v := <-done:
		return v, nil
	case err := <-failed:
		return js.Undefined(), err
	}
}
