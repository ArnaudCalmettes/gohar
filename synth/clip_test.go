package synth

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"
)

// wav builds a mono 16 bit WAV file of `samples`, at `rate`.
func wav(rate uint32, samples ...int16) []byte {
	var b bytes.Buffer
	data := len(samples) * 2
	b.WriteString("RIFF")
	binary.Write(&b, binary.LittleEndian, uint32(36+data))
	b.WriteString("WAVEfmt ")
	for _, v := range []any{uint32(16), uint16(1), uint16(1), rate, rate * 2, uint16(2), uint16(16)} {
		binary.Write(&b, binary.LittleEndian, v)
	}
	b.WriteString("data")
	binary.Write(&b, binary.LittleEndian, uint32(data))
	binary.Write(&b, binary.LittleEndian, samples)
	return b.Bytes()
}

func TestLoadClip(t *testing.T) {
	c, err := LoadClip(bytes.NewReader(wav(SampleRate, 16384, -16384)))
	if err != nil {
		t.Fatal(err)
	}
	if len(c.frames) != 2 || c.frames[0] != [2]float32{0.5, 0.5} || c.frames[1] != [2]float32{-0.5, -0.5} {
		t.Errorf("frames %v, want half scale up then down, on both channels", c.frames)
	}
	if _, err := LoadClip(bytes.NewReader(wav(44100, 0))); err == nil {
		t.Error("a clip at 44.1 kHz accepted")
	}
	if _, err := LoadClip(bytes.NewReader([]byte("not a wav"))); err == nil {
		t.Error("garbage accepted")
	}
}

// A snap dated 10 ms ahead starts on frame 480, at its velocity, and
// rings to its end with no release.
func TestClipScheduled(t *testing.T) {
	c, err := LoadClip(bytes.NewReader(wav(SampleRate, 16384, 16384, 16384)))
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024*BytesPerFrame)
	c.ScheduleOn(60, 0.5, t0.Add(10*time.Millisecond))
	c.ScheduleOff(60, t0.Add(10*time.Millisecond))
	c.render(c, buf, t0)
	s := samples(buf)
	at := func(f int) float32 { return s[f*ChannelCount] }
	if at(479) != 0 || at(480) != 0.25 || at(482) != 0.25 || at(483) != 0 {
		t.Errorf("frames 479 to 483: %v %v %v %v, want silence, three at 0.25, silence",
			at(479), at(480), at(482), at(483))
	}
}

// Two snaps close together overlap rather than cut each other.
func TestClipOverlaps(t *testing.T) {
	c, err := LoadClip(bytes.NewReader(wav(SampleRate, 8192, 8192, 8192, 8192)))
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 8*BytesPerFrame)
	c.ScheduleOn(60, 1, t0)
	c.ScheduleOn(60, 1, t0.Add(2*time.Second/SampleRate))
	c.render(c, buf, t0)
	if got := samples(buf)[2*ChannelCount]; got != 0.5 {
		t.Errorf("frame 2, where both ring: %v, want 0.5", got)
	}
}

func TestClipDoesNotAllocate(t *testing.T) {
	c, err := LoadClip(bytes.NewReader(wav(SampleRate, make([]int16, 4800)...)))
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 256*BytesPerFrame)
	allocs := testing.AllocsPerRun(50, func() {
		c.ScheduleOn(60, 0.8, time.Now().Add(time.Millisecond))
		_, _ = c.Read(buf)
	})
	if allocs != 0 {
		t.Errorf("%v allocations per run", allocs)
	}
}
