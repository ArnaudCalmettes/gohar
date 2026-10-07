//go:build ignore

// bowl writes bowl.wav, the singing bowl of the walker's meditation
// (the easter egg of the lessons): `go run bowl.go` in this directory.
//
// A struck bowl rings in a few inharmonic modes, each split in two
// close frequencies by the bowl's imperfect roundness, which beat. The
// ratios are those measured on Tibetan bowls, around 1, 2.8, 5.4 and
// 8.9; the higher the mode, the shorter it rings.
package main

import (
	"encoding/binary"
	"log"
	"math"
	"os"
)

const (
	rate   = 48000
	length = 5.0   // seconds
	fade   = 0.4   // the last seconds, faded out
	attack = 0.004 // the strike, ramped in
	root   = 293.0 // Hz, about a D4
	peak   = 0.8
)

// A mode is one ring of the bowl: its ratio to the root, its loudness,
// how fast it dies (seconds to fall by e), and the beat of its pair.
type mode struct {
	ratio, gain, tau, beat float64
}

var modes = []mode{
	{1, 1, 1.8, 1.3},
	{2.8, 0.55, 1.1, 2.1},
	{5.4, 0.25, 0.6, 3.2},
	{8.9, 0.12, 0.35, 4.5},
}

func main() {
	n := int(length * rate)
	s := make([]float64, n)
	top := 0.0
	for i := range s {
		t := float64(i) / rate
		v := 0.0
		for _, m := range modes {
			f := root * m.ratio
			a := m.gain * math.Exp(-t/m.tau)
			v += a * (math.Sin(2*math.Pi*(f-m.beat/2)*t) + math.Sin(2*math.Pi*(f+m.beat/2)*t)) / 2
		}
		v *= min(t/attack, 1)
		if r := length - t; r < fade {
			v *= r / fade
		}
		s[i] = v
		top = max(top, math.Abs(v))
	}

	f, err := os.Create("bowl.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	data := uint32(2 * n)
	w := func(v any) {
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			log.Fatal(err)
		}
	}
	f.WriteString("RIFF")
	w(36 + data)
	f.WriteString("WAVEfmt ")
	w(uint32(16))
	w(uint16(1)) // PCM
	w(uint16(1)) // mono
	w(uint32(rate))
	w(uint32(2 * rate))
	w(uint16(2))
	w(uint16(16))
	f.WriteString("data")
	w(data)
	for _, v := range s {
		w(int16(math.Round(v / top * peak * math.MaxInt16)))
	}
}
