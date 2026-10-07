//go:build !js

package band

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ArnaudCalmettes/gohar/synth"
)

// Open builds the band and opens the audio output: on the desktop,
// the synth computes every sample in Go and oto plays them. It returns a
// nil band when -list has only listed the presets, and what closes the
// output once the game is over.
func Open(f Flags) (*Band, func(), error) {
	var sf *synth.SoundFont
	if !f.Chip {
		sf = load(f.SF2)
		if f.List {
			for _, p := range sf.Presets() {
				fmt.Printf("%3d:%-3d %s\n", p.Bank, p.Patch, p.Name)
			}
			return nil, nil, nil
		}
	}
	var mix synth.Mixer
	bd, err := newBand(&mix, sf, f.Bass, f.Split, f.Demo)
	if err != nil {
		return nil, nil, err
	}
	out, err := synth.Open(&mix, synth.Options{Device: f.Device})
	if err != nil {
		return nil, nil, err
	}
	return bd, func() {
		out.Close()
		if err := out.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "audio:", err)
		}
	}, nil
}

// load reads the soundfont at `path`, or the one embedded when `path`
// is empty. A path given that cannot be read stops the game: the player
// asked for those sounds.
func load(path string) *synth.SoundFont {
	var r io.Reader = bytes.NewReader(walkSF2)
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r = f
	}
	sf, err := synth.LoadSoundFont(r)
	if err != nil {
		log.Fatal(err)
	}
	return sf
}

// newBand builds the band on `mix`. A nil soundfont gives the chip
// sounds.
func newBand(mix *synth.Mixer, sf *synth.SoundFont, bassPreset synth.Preset, split int, demo bool) (*Band, error) {
	b := &Band{Split: split, Demo: demo, hat: pedalHat, ride: ride1, ta: lowWood, ti: hiWood}
	if sf == nil {
		hands, err := synth.NewEngine("triangle", true)
		if err != nil {
			return nil, err
		}
		noise, err := synth.NewEngine("noise", true)
		if err != nil {
			return nil, err
		}
		b.bass, b.Piano, b.snap, b.kit = hands, hands, noise, noise
		b.hat, b.ride, b.chip = noiseKey, chipRide, true
		b.ta, b.ti = noiseKey, chipRide
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
	b.bass, b.Piano, b.kit, b.snap = bass, piano, kit, snap
	mix.Add(bass, bassGain)
	mix.Add(piano, pianoGain)
	mix.Add(kit, kitGain)
	mix.Add(snap, snapGain)
	return b, nil
}
