//go:build js && wasm

package main

import (
	"fmt"
	"os"

	"github.com/ArnaudCalmettes/gohar/synth/web"
)

// The channels of FluidSynth the band plays on; the kit takes the drum
// channel of its own (see synth/web).
const (
	bassChannel  = 0
	pianoChannel = 1
)

// openAudio builds the band in the browser: FluidSynth on the audio
// thread plays the embedded soundfont, Web Audio the recorded snap (see
// synth/web, and "Le navigateur" in docs/architecture.md). The volumes
// keep the desktop mix's proportions, the loudest at full volume. The
// chip sounds, another soundfont and the presets' listing stay on the
// desktop; so does the buffer, which the browser sizes itself.
func openAudio(f audioFlags) (*band, func(), error) {
	s, err := web.Open(walkSF2, web.Options{Gain: f.gain})
	if err != nil {
		return nil, nil, err
	}
	base, out := s.Latency()
	fmt.Fprintf(os.Stderr, "audio: base %.1f ms, output %.1f ms\n", base*1000, out*1000)
	loudest := float64(bassGain)
	b := &band{split: f.split, demo: f.demo, hat: pedalHat, ride: ride1, ta: lowWood, ti: hiWood}
	b.bass = s.Instrument(f.bass.Bank, f.bass.Patch, bassChannel, bassGain/loudest)
	b.piano = s.Instrument(acousticGrand.Bank, acousticGrand.Patch, pianoChannel, pianoGain/loudest)
	b.kit = s.Instrument(jazzKit.Bank, jazzKit.Patch, 0, kitGain/loudest)
	if b.snap, err = s.Clip(snapWAV, snapGain); err != nil {
		return nil, nil, err
	}
	return b, func() { s.Close() }, nil
}
