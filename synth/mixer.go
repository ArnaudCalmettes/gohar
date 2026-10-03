package synth

import (
	"encoding/binary"
	"io"
	"math"
)

// mixChunk is how many frames the mixer handles at a time: its buffers
// are allocated once, at this size, and a longer Read is done in
// several passes. 85 ms, far above any buffer the device asks for.
const mixChunk = 4096

// A Mixer sums several instruments into the one stream a Device plays:
// a Device takes a single source, and an Engine plays a single timbre.
//
// Each input has its gain. The sum is clamped like a dense chord in an
// Engine, so that several instruments at full level saturate instead
// of wrapping around.
//
// Inputs are added before the mixer is handed to Open; adding one while
// it plays is a data race. Read follows the rules of the audio
// goroutine: no allocation, no lock, no wait. An input that fails or
// runs dry is heard as silence.
type Mixer struct {
	inputs  []input
	scratch [mixChunk * BytesPerFrame]byte
	sum     [mixChunk * ChannelCount]float32
}

type input struct {
	src  io.Reader
	gain float32
}

// Add mixes `src` in at `gain`, 1 for its own level.
func (m *Mixer) Add(src io.Reader, gain float64) {
	m.inputs = append(m.inputs, input{src: src, gain: float32(gain)})
}

func (m *Mixer) Read(buf []byte) (int, error) {
	frames := len(buf) / BytesPerFrame
	for done := 0; done < frames; {
		n := min(frames-done, mixChunk)
		m.mix(buf[done*BytesPerFrame : (done+n)*BytesPerFrame])
		done += n
	}
	return frames * BytesPerFrame, nil
}

// mix fills `out`, at most mixChunk frames, with the sum of the inputs.
func (m *Mixer) mix(out []byte) {
	samples := len(out) / 4
	sum := m.sum[:samples]
	clear(sum)

	for _, in := range m.inputs {
		s := m.scratch[:len(out)]
		got, _ := io.ReadFull(in.src, s)
		for i := range got / 4 {
			sum[i] += in.gain * math.Float32frombits(binary.LittleEndian.Uint32(s[i*4:]))
		}
	}

	for i, v := range sum {
		binary.LittleEndian.PutUint32(out[i*4:], math.Float32bits(float32(clamp(float64(v)))))
	}
}
