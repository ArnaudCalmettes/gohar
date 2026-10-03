package synth

import (
	"encoding/binary"
	"io"
	"math"
	"testing"
)

// level is an input that writes the same sample everywhere: a
// stand-in for an instrument, easy to add up by hand.
type level float32

func (l level) Read(buf []byte) (int, error) {
	for i := 0; i+4 <= len(buf); i += 4 {
		binary.LittleEndian.PutUint32(buf[i:], math.Float32bits(float32(l)))
	}
	return len(buf), nil
}

// dry runs out after one frame.
type dry struct{ done bool }

func (d *dry) Read(buf []byte) (int, error) {
	if d.done {
		return 0, io.EOF
	}
	d.done = true
	return level(0.5).Read(buf[:BytesPerFrame])
}

func samples(buf []byte) []float32 {
	out := make([]float32, len(buf)/4)
	for i := range out {
		out[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf[i*4:]))
	}
	return out
}

func TestMixerSums(t *testing.T) {
	for _, tc := range []struct {
		name string
		add  func(m *Mixer)
		want float32
	}{
		{"two inputs at their level", func(m *Mixer) { m.Add(level(0.25), 1); m.Add(level(0.5), 1) }, 0.75},
		{"a gain", func(m *Mixer) { m.Add(level(0.5), 0.5) }, 0.25},
		{"too loud, clamped", func(m *Mixer) { m.Add(level(0.8), 1); m.Add(level(0.8), 1) }, 1},
		{"no input, silence", func(m *Mixer) {}, 0},
	} {
		var m Mixer
		tc.add(&m)
		buf := make([]byte, 64*BytesPerFrame)
		_, _ = m.Read(buf)
		for i, v := range samples(buf) {
			if v != tc.want {
				t.Errorf("%s: sample %d is %v, want %v", tc.name, i, v, tc.want)
				break
			}
		}
	}
}

// A buffer longer than the mixer's chunk is filled to the end.
func TestMixerLongBuffer(t *testing.T) {
	var m Mixer
	m.Add(level(0.5), 1)
	buf := make([]byte, (mixChunk+100)*BytesPerFrame)
	if n, _ := m.Read(buf); n != len(buf) {
		t.Fatalf("read %d bytes, want %d", n, len(buf))
	}
	s := samples(buf)
	if last := s[len(s)-1]; last != 0.5 {
		t.Errorf("last sample %v, want 0.5", last)
	}
}

// An input that runs dry is silence from there on, and the others play.
func TestMixerDryInput(t *testing.T) {
	var m Mixer
	m.Add(&dry{}, 1)
	m.Add(level(0.25), 1)
	buf := make([]byte, 4*BytesPerFrame)
	_, _ = m.Read(buf)
	s := samples(buf)
	if s[0] != 0.75 || s[ChannelCount] != 0.25 {
		t.Errorf("first frame %v, second %v, want 0.75 then 0.25", s[0], s[ChannelCount])
	}
}

func TestMixerDoesNotAllocate(t *testing.T) {
	var m Mixer
	m.Add(&Engine{Timbre: Square}, 0.5)
	m.Add(&Engine{Timbre: Noise}, 0.5)
	buf := make([]byte, 256*BytesPerFrame)
	allocs := testing.AllocsPerRun(50, func() { _, _ = m.Read(buf) })
	if allocs != 0 {
		t.Errorf("%v allocations per Read", allocs)
	}
}
