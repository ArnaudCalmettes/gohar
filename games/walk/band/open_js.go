//go:build js && wasm

package band

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

// Open builds the band in the browser: FluidSynth on the audio
// thread plays the embedded soundfont, Web Audio the snap and the bowl
// (see synth/web, and "Le navigateur" in docs/architecture.md). The volumes
// keep the desktop mix's proportions, the loudest at full volume. The
// chip sounds, another soundfont and the presets' listing stay on the
// desktop; so does the buffer, which the browser sizes itself.
func Open(f Flags) (*Band, func(), error) {
	s, err := web.Open(walkSF2, web.Options{Gain: f.Gain})
	if err != nil {
		return nil, nil, err
	}
	base, out := s.Latency()
	fmt.Fprintf(os.Stderr, "audio: base %.1f ms, output %.1f ms\n", base*1000, out*1000)
	loudest := float64(bassGain)
	b := &Band{Split: f.Split, Demo: f.Demo, hat: pedalHat, ride: ride1, ta: lowWood, ti: hiWood}
	b.bass = s.Instrument(f.Bass.Bank, f.Bass.Patch, bassChannel, bassGain/loudest)
	b.Piano = s.Instrument(acousticGrand.Bank, acousticGrand.Patch, pianoChannel, pianoGain/loudest)
	b.kit = s.Instrument(jazzKit.Bank, jazzKit.Patch, 0, kitGain/loudest)
	if b.snap, err = s.Clip(snapWAV, snapGain); err != nil {
		return nil, nil, err
	}
	if b.bowl, err = s.Clip(bowlWAV, bowlGain); err != nil {
		return nil, nil, err
	}
	return b, func() { s.Close() }, nil
}
