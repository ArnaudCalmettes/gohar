package synth

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

// A SoundFont is a parsed SF2 file: recorded instruments, sampled note
// by note, that a Sampler plays.
//
// Parsing reads every sample into memory and allocates as much: load
// once, before the game starts, and share the result between samplers.
type SoundFont struct {
	sf *meltysynth.SoundFont
}

// A Preset names one instrument of a soundfont the way General MIDI
// does: a bank and a patch number, 0 to 127. Bank 128 holds the drum
// kits, where each key is a different percussion.
type Preset struct {
	Name  string
	Bank  int
	Patch int
}

// LoadSoundFont parses an SF2 file.
func LoadSoundFont(r io.Reader) (*SoundFont, error) {
	sf, err := meltysynth.NewSoundFont(r)
	if err != nil {
		return nil, fmt.Errorf("synth: reading the soundfont: %w", err)
	}
	return &SoundFont{sf: sf}, nil
}

// Presets lists the instruments of the soundfont, in file order.
func (s *SoundFont) Presets() []Preset {
	out := make([]Preset, len(s.sf.Presets))
	for i, p := range s.sf.Presets {
		out[i] = Preset{Name: p.Name, Bank: int(p.BankNumber), Patch: int(p.PatchNumber)}
	}
	return out
}

// samplerBlock is how many frames meltysynth computes at a time. A
// command applied between two calls only takes effect at the next
// block, so this is the granularity of a dated note: 16 frames, a third
// of a millisecond, well under the driver's own jitter. The Engine
// lands on the sample; a Sampler lands within a block of it.
const samplerBlock = 16

// A Sampler plays one preset of a soundfont: a double bass, a piano, a
// drum kit.
//
// It is an Instrument like the Engine, with the same queue, so it keeps
// dates and measures delays the same way. Each Sampler has its own
// meltysynth synthesizer, which allocates its voices once, when built:
// Read, NoteOn and NoteOff allocate nothing afterwards.
//
// Reverb and chorus are off: they cost more than every voice together,
// and a bass in a mix wants neither. To revisit with the piano.
type Sampler struct {
	queue

	synth *meltysynth.Synthesizer

	left, right [mixChunk]float32
}

var _ Instrument = (*Sampler)(nil)

// NewSampler builds an instrument that plays `p` from `sf`. A preset
// the soundfont does not have is an error rather than meltysynth's
// silent fallback to its first preset: a game that asked for a double
// bass and got a harpsichord would sound wrong without saying why.
func NewSampler(sf *SoundFont, p Preset) (*Sampler, error) {
	found := false
	for _, q := range sf.sf.Presets {
		if int(q.BankNumber) == p.Bank && int(q.PatchNumber) == p.Patch {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("synth: no preset at bank %d, patch %d", p.Bank, p.Patch)
	}

	settings := meltysynth.NewSynthesizerSettings(SampleRate)
	settings.BlockSize = samplerBlock
	settings.EnableReverbAndChorus = false
	ms, err := meltysynth.NewSynthesizer(sf.sf, settings)
	if err != nil {
		return nil, fmt.Errorf("synth: starting the sampler: %w", err)
	}

	// Everything plays on channel 0, set to the preset. Bank 128 needs
	// no channel 9: meltysynth looks presets up by bank and patch alone.
	const controller, programChange, bankSelect = 0xB0, 0xC0, 0x00
	ms.ProcessMidiMessage(0, controller, bankSelect, int32(p.Bank))
	ms.ProcessMidiMessage(0, programChange, int32(p.Patch), 0)

	return &Sampler{synth: ms}, nil
}

func (s *Sampler) Read(buf []byte) (int, error) {
	return s.render(s, buf, time.Now()), nil
}

func (s *Sampler) apply(c command) {
	if c.on {
		s.synth.NoteOn(0, int32(c.key), midiVelocity(c.velocity))
	} else {
		s.synth.NoteOff(0, int32(c.key))
	}
}

// fill renders the frames of `buf`, mixChunk at a time.
func (s *Sampler) fill(buf []byte) {
	frames := len(buf) / BytesPerFrame
	for done := 0; done < frames; {
		n := min(frames-done, mixChunk)
		s.synth.Render(s.left[:n], s.right[:n])
		out := buf[done*BytesPerFrame:]
		for i := range n {
			writeStereo(out[i*BytesPerFrame:], s.left[i], s.right[i])
		}
		done += n
	}
}

// midiVelocity turns a velocity between 0 and 1 into MIDI's 1 to 127.
// Never 0, which MIDI reads as a release.
const maxVelocity = 127

func midiVelocity(v float64) int32 {
	return int32(max(1, min(maxVelocity, math.Round(v*maxVelocity))))
}

// writeStereo writes a frame, clamped: a soundfont recorded hot, or a
// chord struck hard, may exceed full scale.
func writeStereo(buf []byte, l, r float32) {
	writeSample(buf[0:], float32(clamp(float64(l))))
	writeSample(buf[bytesPerSample:], float32(clamp(float64(r))))
}
