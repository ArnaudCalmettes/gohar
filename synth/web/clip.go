//go:build js && wasm

package web

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/ArnaudCalmettes/gohar/synth"
)

// A clip is a recorded sound, a synth.Instrument as the desktop's Clip
// is: decoded once by the browser, and started at its date by Web Audio
// itself, on the audio thread. The key is ignored, as on the desktop:
// a clip has one sound; the velocity scales it.
type clip struct {
	s      *Synth
	buffer js.Value // the decoded AudioBuffer
	gain   float64
}

// Clip decodes `wav` and plays it at `gain`.
func (s *Synth) Clip(wav []byte, gain float64) (synth.Instrument, error) {
	buf := js.Global().Get("Uint8Array").New(len(wav))
	js.CopyBytesToJS(buf, wav)
	decoded, err := await(s.ctx.Call("decodeAudioData", buf.Get("buffer")))
	if err != nil {
		return nil, fmt.Errorf("web: decoding a clip: %w", err)
	}
	return &clip{s: s, buffer: decoded, gain: gain}, nil
}

func (c *clip) NoteOn(_ int, velocity float64, _ time.Time) { c.start(velocity, 0) }

func (c *clip) ScheduleOn(_ int, velocity float64, at time.Time) {
	c.start(velocity, max(time.Until(at).Seconds(), 0))
}

// NoteOff and ScheduleOff do nothing: a clip plays to its end.
func (c *clip) NoteOff(int, time.Time)              {}
func (c *clip) ScheduleOff(int, time.Time)          {}
func (c *clip) Read([]byte) (int, error)            { return 0, nil }
func (c *clip) Delays() (last, worst time.Duration) { return 0, 0 }
func (c *clip) Histogram() synth.Histogram          { return synth.Histogram{} }

// start plays the clip `in` seconds from now, at `velocity`: a source
// and a gain of its own, which the browser drops once played.
func (c *clip) start(velocity, in float64) {
	ctx := c.s.ctx
	g := ctx.Call("createGain")
	g.Get("gain").Set("value", c.gain*velocity)
	g.Call("connect", ctx.Get("destination"))
	src := ctx.Call("createBufferSource")
	src.Set("buffer", c.buffer)
	src.Call("connect", g)
	src.Call("start", ctx.Get("currentTime").Float()+in)
}
