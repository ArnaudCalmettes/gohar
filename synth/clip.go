package synth

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"time"
)

// maxClipVoices is how many takes of a clip may ring at once. A snap
// lasts a fraction of a beat: four overlapping already means a tempo no
// one plays at.
const maxClipVoices = 4

// A Clip plays one recorded sound, a percussion: a finger snap. Every
// press starts it again from the top, at the press's velocity, and it
// rings out on its own; the key does not matter, and a release does
// nothing.
//
// It exists for the sounds no soundfont has. Like the Sampler, it is an
// Instrument with the common queue, so it keeps dates the same way.
type Clip struct {
	queue

	frames [][2]float32
	voices [maxClipVoices]clipVoice
}

type clipVoice struct {
	pos    int // next frame to play; len(frames) when silent
	gain   float32
	active bool
}

var _ Instrument = (*Clip)(nil)

// LoadClip reads a WAV file: PCM in 16 or 24 bits, or 32 bit floats,
// mono or stereo, at SampleRate. Another rate is an error rather than
// a resampling, the package's rule.
func LoadClip(r io.Reader) (*Clip, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	frames, err := decodeWAV(data)
	if err != nil {
		return nil, fmt.Errorf("synth: reading the clip: %w", err)
	}
	return &Clip{frames: frames}, nil
}

func (c *Clip) Read(buf []byte) (int, error) {
	return c.render(c, buf, time.Now()), nil
}

func (c *Clip) apply(cmd command) {
	if !cmd.on {
		return
	}
	// A free voice, or else the one furthest into its sound.
	slot := 0
	for i := range c.voices {
		if !c.voices[i].active {
			slot = i
			break
		}
		if c.voices[i].pos > c.voices[slot].pos {
			slot = i
		}
	}
	c.voices[slot] = clipVoice{gain: float32(cmd.velocity), active: true}
}

func (c *Clip) fill(buf []byte) {
	for f := range len(buf) / BytesPerFrame {
		var l, r float32
		for i := range c.voices {
			v := &c.voices[i]
			if !v.active {
				continue
			}
			l += v.gain * c.frames[v.pos][0]
			r += v.gain * c.frames[v.pos][1]
			v.pos++
			if v.pos == len(c.frames) {
				*v = clipVoice{}
			}
		}
		writeStereo(buf[f*BytesPerFrame:], l, r)
	}
}

// The formats of a WAV file this decoder reads, from its "fmt " chunk.
const (
	wavPCM        = 1
	wavFloat      = 3
	wavExtensible = 0xFFFE // the real format follows, further in the chunk
)

// decodeWAV returns the frames of a WAV file, mono spread over both
// channels. Only the chunks it needs are read, "fmt " and "data"; the
// others, metadata, are skipped.
func decodeWAV(b []byte) ([][2]float32, error) {
	if len(b) < 12 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WAVE" {
		return nil, errors.New("not a WAV file")
	}
	var format, channels, bits int
	var rate uint32
	var data []byte
	for p := 12; p+8 <= len(b); {
		id, size := string(b[p:p+4]), int(binary.LittleEndian.Uint32(b[p+4:]))
		body := b[p+8 : min(len(b), p+8+size)]
		switch id {
		case "fmt ":
			if len(body) < 16 {
				return nil, errors.New("short format chunk")
			}
			format = int(binary.LittleEndian.Uint16(body[0:]))
			channels = int(binary.LittleEndian.Uint16(body[2:]))
			rate = binary.LittleEndian.Uint32(body[4:])
			bits = int(binary.LittleEndian.Uint16(body[14:]))
			if format == wavExtensible && len(body) >= 26 {
				format = int(binary.LittleEndian.Uint16(body[24:]))
			}
		case "data":
			data = body
		}
		p += 8 + size + size%2 // chunks are padded to an even size
	}

	switch {
	case data == nil || channels == 0:
		return nil, errors.New("no format or no data")
	case rate != SampleRate:
		return nil, fmt.Errorf("sampled at %d Hz, want %d", rate, SampleRate)
	case channels > 2:
		return nil, fmt.Errorf("%d channels, want 1 or 2", channels)
	}

	var sample func([]byte) float32
	switch {
	case format == wavPCM && bits == 16:
		sample = func(s []byte) float32 { return float32(int16(binary.LittleEndian.Uint16(s))) / (1 << 15) }
	case format == wavPCM && bits == 24:
		sample = func(s []byte) float32 {
			v := int32(uint32(s[0])<<8|uint32(s[1])<<16|uint32(s[2])<<24) >> 8
			return float32(v) / (1 << 23)
		}
	case format == wavFloat && bits == 32:
		sample = func(s []byte) float32 { return math.Float32frombits(binary.LittleEndian.Uint32(s)) }
	default:
		return nil, fmt.Errorf("format %d in %d bits, want PCM 16 or 24, or float 32", format, bits)
	}

	width := bits / 8
	frames := make([][2]float32, len(data)/(width*channels))
	for i := range frames {
		s := data[i*width*channels:]
		frames[i][0] = sample(s)
		frames[i][1] = frames[i][0]
		if channels == 2 {
			frames[i][1] = sample(s[width:])
		}
	}
	return frames, nil
}
