package synth

import (
	"fmt"
	"io"
	"time"

	"github.com/ebitengine/oto/v3"
)

// Options are the buffer sizes of the audio path.
//
// # Why they are named here and nowhere else
//
// Four queues separate a decision from a sound, and every one of them
// has a default calibrated for playing a file rather than for playing
// an instrument. oto's player ring defaults to half a second, and its
// ALSA driver to a period of 1024 frames repeated twice, which is
// 42.7 ms at this sample rate. Both are sane for a backing track and
// ruinous for a keyboard.
//
// Hence the rule this package exists to enforce: no buffer is ever left
// at its default. That is also why the audio device lives here and not
// in the games. A game that opened oto itself would have to remember
// the rule, and one day one of them would not.
//
// Measured on 2026-09-25, a Focusrite Scarlett Solo over USB holds 5 ms
// on both without a glitch, for a floor around 12.5 ms and no audible
// flam. Those are one machine's figures, not a promise: the numbers a
// player gets are the player's, which is why calibration stays on the
// programme.
type Options struct {
	// Device is what ALSA hands the card. Zero means DefaultBuffer.
	Device time.Duration

	// Player is the ring oto keeps above the device, filled the moment
	// playback starts. Zero means the same as Device.
	Player time.Duration
}

// DefaultBuffer is the starting point, not a recommendation. Raise it
// when a machine crackles, and expect to have to.
const DefaultBuffer = 10 * time.Millisecond

// A Device is an open audio output. One per process: oto refuses a
// second context, and so would the sound card.
type Device struct {
	ctx    *oto.Context
	player *oto.Player
	opts   Options
}

// Open starts playing `src` and returns once the device is ready.
//
// `src` is read from the audio goroutine, so it must not allocate, must
// not block, and must always fill the buffer it is given. An Engine
// satisfies all three.
func Open(src io.Reader, opts Options) (*Device, error) {
	if opts.Device == 0 {
		opts.Device = DefaultBuffer
	}
	if opts.Player == 0 {
		opts.Player = opts.Device
	}

	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   SampleRate,
		ChannelCount: ChannelCount,
		Format:       oto.FormatFloat32LE,
		BufferSize:   opts.Device,
	})
	if err != nil {
		return nil, fmt.Errorf("synth: opening the audio device: %w", err)
	}
	<-ready

	p := ctx.NewPlayer(src)
	p.SetBufferSize(sizeOf(opts.Player))
	p.Play()

	return &Device{ctx: ctx, player: p, opts: opts}, nil
}

// Queued reports how much audio is already waiting to be played.
//
// The part of the latency an engine cannot see: a sample it writes now
// comes out after everything already in the queue. Useful to check that
// a buffer setting took, and to notice when it drifts.
func (d *Device) Queued() time.Duration {
	return spanOf(d.player.BufferedSize())
}

// Floor is the best case between writing a sample and hearing it: what
// is queued plus one device buffer.
//
// A floor and not a measurement. The bus and the converter add their
// own, invisible from here. The only honest reading of the whole path
// is a flam heard against a real sound.
func (d *Device) Floor() time.Duration {
	return d.Queued() + d.opts.Device
}

// Err reports a failure that happened on the audio goroutine, where
// there was nobody to return it to.
func (d *Device) Err() error {
	return d.player.Err()
}

func (d *Device) Close() error {
	return d.player.Close()
}

func sizeOf(d time.Duration) int {
	frames := int(int64(d) * SampleRate / int64(time.Second))
	return frames * BytesPerFrame
}

func spanOf(n int) time.Duration {
	frames := n / BytesPerFrame
	return time.Duration(frames) * time.Second / SampleRate
}
