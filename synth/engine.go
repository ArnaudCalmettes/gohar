package synth

import (
	"encoding/binary"
	"math"
	"sync"
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

// attack and release are short ramps, not an expressive envelope.
//
// They exist to avoid the click a square edge makes, which on a short
// note is louder than the note. Anything musical, a real envelope or a
// velocity curve, comes later and belongs above this.
const (
	attack  = 3 * time.Millisecond
	release = 40 * time.Millisecond
)

type voiceState uint8

const (
	voiceOff voiceState = iota
	voiceAttack
	voiceSustain
	voiceRelease
)

const headroom = 1.0 / 6

type voice struct {
	key   int
	phase float64
	step  float64
	gain  float64
	env   float64
	state voiceState
}

type command struct {
	on       bool
	key      int
	velocity float64
	at       time.Time
}

// An Engine turns key presses into samples. It is the whole instrument
// for now: one sine per key, no filter, no effect.
//
// # Two goroutines, one handover
//
// NoteOn and NoteOff are called from wherever the keys come from, a
// MIDI callback or a game loop. Read is called by the audio driver on
// its own goroutine and must never wait, never allocate and never
// block on anything.
//
// So the two sides share a queue of commands and nothing else. The
// mutex is held only while that queue is drained at the top of Read,
// never while samples are computed, and the voices are touched by the
// audio goroutine alone.
//
// The zero Engine is usable and plays at 440.
type Engine struct {
	Tuning Tuning

	mu        sync.Mutex
	commands  [maxCommands]command
	nCommands int
	last      time.Duration
	worst     time.Duration
	histogram Histogram

	voices [maxVoices]voice
}

// NoteOn starts a key, or restarts it if it was already sounding.
//
// at is when the press happened, which the caller knows better than we
// do: a MIDI driver hands over an event that already waited. It is
// recorded so that Delays can say how long the press took to reach the
// samples.
func (e *Engine) NoteOn(key int, velocity float64, at time.Time) {
	e.push(command{on: true, key: key, velocity: velocity, at: at})
}

// NoteOff releases a key. A key not sounding is not an error.
func (e *Engine) NoteOff(key int, at time.Time) {
	e.push(command{key: key, at: at})
}

func (e *Engine) push(c command) {
	e.mu.Lock()
	if e.nCommands < len(e.commands) {
		e.commands[e.nCommands] = c
		e.nCommands++
	}
	e.mu.Unlock()
}

// Delays returns the last and the worst delay between a key event and
// the moment it reached the samples.
//
// This is the part of the latency that belongs to us. What the driver,
// the bus and the card add afterwards is not visible from here, and the
// only honest measurement of the whole remains the ear.
func (e *Engine) Delays() (last, worst time.Duration) {
	e.mu.Lock()
	last, worst = e.last, e.worst
	e.mu.Unlock()
	return
}

// Histogram returns how the delays between a key event and the samples
// were spread since the engine started.
//
// A copy, taken under the lock the audio goroutine holds a few
// microseconds per buffer: cheap enough to read every frame.
func (e *Engine) Histogram() Histogram {
	e.mu.Lock()
	h := e.histogram
	e.mu.Unlock()
	return h
}

func (e *Engine) Read(buf []byte) (int, error) {
	now := time.Now()

	e.mu.Lock()
	for i := range e.nCommands {
		d := now.Sub(e.commands[i].at)
		e.last = d
		if d > e.worst {
			e.worst = d
		}
		e.histogram.add(d)
		e.apply(e.commands[i])
	}
	e.nCommands = 0
	e.mu.Unlock()

	frames := len(buf) / BytesPerFrame
	attackFrames := float64(SampleRate) * attack.Seconds()
	releaseFrames := float64(SampleRate) * release.Seconds()

	for f := range frames {
		var sample float64
		for v := range e.voices {
			vo := &e.voices[v]
			if vo.state == voiceOff {
				continue
			}

			switch vo.state {
			case voiceAttack:
				vo.env += vo.gain / attackFrames
				if vo.env >= vo.gain {
					vo.env = vo.gain
					vo.state = voiceSustain
				}
			case voiceRelease:
				vo.env -= vo.gain / releaseFrames
				if vo.env <= 0 {
					*vo = voice{}
					continue
				}
			}

			sample += vo.env * math.Sin(vo.phase)
			vo.phase += vo.step
			if vo.phase > 2*math.Pi {
				vo.phase -= 2 * math.Pi
			}
		}

		writeFrame(buf[f*BytesPerFrame:], float32(clamp(sample)))
	}

	return frames * BytesPerFrame, nil
}

func (e *Engine) apply(c command) {
	if !c.on {
		for i := range e.voices {
			if e.voices[i].state != voiceOff && e.voices[i].key == c.key {
				e.voices[i].state = voiceRelease
			}
		}
		return
	}

	step := 2 * math.Pi * e.Tuning.Frequency(c.key) / SampleRate

	// Retrigger the same key rather than stacking a second voice on it,
	// so that a repeated note does not grow louder than a single one.
	for i := range e.voices {
		if e.voices[i].state != voiceOff && e.voices[i].key == c.key {
			e.voices[i].step = step
			e.voices[i].gain = c.velocity * headroom
			e.voices[i].state = voiceAttack
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
			if e.voices[i].env < e.voices[slot].env {
				slot = i
			}
		}
	}

	e.voices[slot] = voice{
		key:   c.key,
		step:  step,
		gain:  c.velocity * headroom,
		state: voiceAttack,
	}
}

// clamp keeps a dense chord from wrapping around into noise. Summing
// twenty four sines can exceed one, and a float32 that overflows the
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
