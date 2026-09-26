package synth

import (
	"encoding/binary"
	"math"
	"time"
)

const (
	// SampleRate is fixed for the whole package. One rate, one context,
	// no resampling anywhere: a game that needs another one is not a
	// case we have.
	SampleRate = 48000

	// ChannelCount is two because oto's driver works in stereo floats
	// whatever we ask, so mono would only mean writing the same sample
	// twice further down.
	ChannelCount = 2

	// BytesPerFrame is one float32 per channel.
	BytesPerFrame = ChannelCount * 4
)

const (
	// maxVoices is how many keys may sound at once. Ten fingers and a
	// sustain pedal do not reach it, and a fixed array keeps Read free
	// of allocation.
	maxVoices = 24

	// maxCommands is the depth of the handover queue. A dropped note is
	// better than a control goroutine blocking on the audio one.
	maxCommands = 64
)

const headroom = 1.0 / 6

// An Envelope shapes a note over time: it rises to full level in
// Attack, falls to Sustain in Decay, holds there while the key is down,
// and dies away in Release once it is up.
//
// The zero Envelope is the default: a rise short enough to be heard as
// immediate and long enough not to click, a full sustain, and a short
// release. It is what the sine has always had.
type Envelope struct {
	Attack  time.Duration
	Decay   time.Duration
	Sustain float64 // between 0 and 1
	Release time.Duration
}

// DefaultEnvelope is what the zero Envelope stands for.
var DefaultEnvelope = Envelope{
	Attack:  3 * time.Millisecond,
	Sustain: 1,
	Release: 40 * time.Millisecond,
}

// ChipEnvelope suits the pulse timbres: a quick attack, a fall to a
// held level, a short tail, the shape the old sound chips were driven
// with.
var ChipEnvelope = Envelope{
	Attack:  2 * time.Millisecond,
	Decay:   150 * time.Millisecond,
	Sustain: 0.6,
	Release: 80 * time.Millisecond,
}

type voiceState uint8

const (
	voiceOff voiceState = iota
	voiceAttack
	voiceDecay
	voiceSustain
	voiceRelease
)

type voice struct {
	key   int
	osc   oscillator
	gain  float64 // the velocity, times headroom
	env   float64 // the envelope's level, between 0 and 1
	fall  float64 // how much env loses per frame once released
	state voiceState
}

// An Engine turns key presses into samples, one voice per key, in one
// of a handful of timbres.
//
// # Two goroutines, one handover
//
// NoteOn and NoteOff are called from wherever the keys come from, a
// MIDI callback or a game loop. Read is called by the audio driver on
// its own goroutine and must never wait, never allocate and never
// block on anything.
//
// So the two sides share a queue of commands and nothing else, and the
// voices are touched by the audio goroutine alone.
//
// # Settings
//
// Tuning, Timbre, Smooth and Envelope are read by the audio goroutine:
// set them before handing the engine to Open. Changing them while it
// plays is a data race.
//
// The zero Engine is usable: a sine at 440 with the default envelope.
type Engine struct {
	queue

	Tuning Tuning
	Timbre Timbre

	// Smooth rounds the jumps of the pulse timbres and draws the
	// triangle as a line; left false, they keep the aliasing and the
	// steps of the consoles. See oscillator.next.
	Smooth bool

	Envelope Envelope

	voices  [maxVoices]voice
	pending [maxCommands]command
}

var _ Instrument = (*Engine)(nil)

// NewEngine builds an engine for a timbre given by name, with the
// envelope that suits it: the default one for the sine, the chip one
// for the others.
func NewEngine(timbre string, smooth bool) (*Engine, error) {
	t, err := ParseTimbre(timbre)
	if err != nil {
		return nil, err
	}
	e := &Engine{Timbre: t, Smooth: smooth}
	if t != Sine {
		e.Envelope = ChipEnvelope
	}
	return e, nil
}

func (e *Engine) envelope() Envelope {
	if e.Envelope == (Envelope{}) {
		return DefaultEnvelope
	}
	return e.Envelope
}

// frames converts a duration into a count of frames, never below one so
// that a zero duration is an immediate step rather than a division by
// zero.
func frames(d time.Duration) float64 {
	return max(1, float64(SampleRate)*d.Seconds())
}

func (e *Engine) Read(buf []byte) (int, error) {
	n := e.take(time.Now(), &e.pending)
	for i := range n {
		e.apply(e.pending[i])
	}

	env := e.envelope()
	rise := 1 / frames(env.Attack)
	decay := (1 - env.Sustain) / frames(env.Decay)
	releaseFrames := frames(env.Release)

	count := len(buf) / BytesPerFrame
	for f := range count {
		var sample float64
		for v := range e.voices {
			vo := &e.voices[v]
			if vo.state == voiceOff {
				continue
			}

			switch vo.state {
			case voiceAttack:
				vo.env += rise
				if vo.env >= 1 {
					vo.env = 1
					vo.state = voiceDecay
				}
			case voiceDecay:
				vo.env -= decay
				if vo.env <= env.Sustain {
					vo.env = env.Sustain
					vo.state = voiceSustain
				}
			case voiceRelease:
				if vo.fall == 0 {
					// Released on this frame: fall from wherever the
					// envelope stood, so that a key let go during its
					// attack does not jump.
					vo.fall = max(vo.env, 1e-9) / releaseFrames
				}
				vo.env -= vo.fall
				if vo.env <= 0 {
					*vo = voice{}
					continue
				}
			}

			sample += vo.gain * vo.env * vo.osc.next(e.Timbre, e.Smooth)
		}

		writeFrame(buf[f*BytesPerFrame:], float32(clamp(sample)))
	}

	return count * BytesPerFrame, nil
}

func (e *Engine) apply(c command) {
	if !c.on {
		for i := range e.voices {
			if e.voices[i].state != voiceOff && e.voices[i].key == c.key {
				e.voices[i].state = voiceRelease
				e.voices[i].fall = 0
			}
		}
		return
	}

	inc := e.Tuning.Frequency(c.key) / SampleRate

	// Retrigger the same key rather than stacking a second voice on it,
	// so that a repeated note does not grow louder than a single one.
	for i := range e.voices {
		if e.voices[i].state != voiceOff && e.voices[i].key == c.key {
			e.voices[i].osc.inc = inc
			e.voices[i].gain = c.velocity * headroom
			e.voices[i].state = voiceAttack
			e.voices[i].fall = 0
			return
		}
	}

	slot := -1
	for i := range e.voices {
		if e.voices[i].state == voiceOff {
			slot = i
			break
		}
	}
	if slot < 0 {
		// Steal the quietest voice, which is the one furthest into its
		// release and therefore the least missed.
		slot = 0
		for i := range e.voices {
			if e.voices[i].env*e.voices[i].gain < e.voices[slot].env*e.voices[slot].gain {
				slot = i
			}
		}
	}

	e.voices[slot] = voice{
		key:   c.key,
		osc:   oscillator{inc: inc, lfsr: 1},
		gain:  c.velocity * headroom,
		state: voiceAttack,
	}
}

// clamp keeps a dense chord from wrapping around into noise. Summing
// twenty four voices can exceed one, and a float32 that overflows the
// converter's range does not saturate, it screams.
func clamp(v float64) float64 {
	switch {
	case v > 1:
		return 1
	case v < -1:
		return -1
	}
	return v
}

func writeFrame(buf []byte, v float32) {
	bits := math.Float32bits(v)
	binary.LittleEndian.PutUint32(buf[0:], bits)
	binary.LittleEndian.PutUint32(buf[4:], bits)
}
