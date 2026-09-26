// Command latency measures how long it takes for a keypress to become a
// sound.
//
// # Why this exists before anything else
//
// The project synthesises its own audio, so the player's keyboard is a
// controller and this program's output is what they hear themselves
// play. Past roughly twenty or thirty milliseconds a musician feels the
// lag and stops trusting the instrument, so the whole design depends on
// a number nobody has measured yet.
//
// # What it can and cannot tell you
//
// Software can measure the part it controls: the delay between the
// keypress and the moment the sound is written into a buffer, plus the
// buffer size that has been asked for. It cannot measure what the
// operating system's mixer and the sound card add afterwards, which on
// some machines is the larger half.
//
// So the printed figure is a floor, not the answer. The answer is your
// ear, which is why the metronome mode is here: latency is nearly
// impossible to judge on an isolated click and obvious against a pulse.
// Turn it on, play along, and see whether you feel yourself dragging.
//
// # What it deliberately leaves out
//
// MIDI. This measures the audio half only. A MIDI keyboard adds its own
// delay on top, and mixing the two into one number would hide which is
// which. Measure this first, then the MIDI path separately.
//
//	1 through 5   ask for a buffer of 10, 20, 30, 50 or 100 ms
//	space         one blip
//	m             metronome on and off
//	q             quit
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	sampleRate = 48000
	channels   = 2

	// blipLength is short enough to feel like a click rather than a
	// note, which is what makes a timing judgement easy.
	blipLength = 30 * time.Millisecond

	// rampLength fades the blip in and out. Without it the waveform
	// starts and stops on a discontinuity, which is heard as a crack
	// and would sit on top of the very thing being measured.
	rampLength = 2 * time.Millisecond

	blipFrequency = 880.0

	metronomePeriod = 500 * time.Millisecond
)

// bufferSizes are the candidates, smallest first. Smaller means less
// delay and more risk of the stream running dry, which is heard as
// crackling.
var bufferSizes = []time.Duration{
	10 * time.Millisecond,
	20 * time.Millisecond,
	30 * time.Millisecond,
	50 * time.Millisecond,
	100 * time.Millisecond,
}

// A voice is one blip being sounded.
//
// A fixed array of these rather than a slice, because Read runs on the
// audio goroutine and must not allocate: a garbage collection pause in
// the middle of filling a buffer is heard as a gap.
type voice struct {
	active    bool
	phase     float64
	remaining int
}

const maxVoices = 8

// A stream turns triggers into samples. It implements io.Reader, which
// is what the audio context consumes.
//
// The game goroutine and the audio goroutine meet here and nowhere
// else. The mutex is held for the few instructions it takes to move a
// trigger across, never for the synthesis itself.
type stream struct {
	mu sync.Mutex

	// pending holds keypresses the audio side has not seen yet, with
	// the moment they happened.
	pending  [maxVoices]time.Time
	nPending int

	// lastDelay is the measured gap between a keypress and the buffer
	// that carried it. Read writes it, the game loop reads it.
	lastDelay time.Duration
	worst     time.Duration

	voices [maxVoices]voice
}

// Trigger records a keypress. Called from the game loop.
func (s *stream) Trigger(at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.nPending < len(s.pending) {
		s.pending[s.nPending] = at
		s.nPending++
	}
}

// Delays returns the last and worst measured delays.
func (s *stream) Delays() (last, worst time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastDelay, s.worst
}

// ResetWorst clears the running maximum, so that a change of buffer
// size can be judged on its own.
func (s *stream) ResetWorst() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.worst = 0
}

// Read fills `buf` with float32 little endian stereo samples.
//
// Runs on the audio goroutine. No allocation, no logging, no waiting on
// anything the game loop might hold for long: whatever happens here
// happens between two moments the sound card is expecting data.
func (s *stream) Read(buf []byte) (int, error) {
	const frameSize = channels * 4

	frames := len(buf) / frameSize
	if frames == 0 {
		return 0, nil
	}

	s.startPending(time.Now())

	for i := range frames {
		var sample float64
		for v := range s.voices {
			sample += s.voices[v].next()
		}
		sample = clip(sample * 0.3)

		bits := math.Float32bits(float32(sample))
		offset := i * frameSize
		binary.LittleEndian.PutUint32(buf[offset:], bits)
		binary.LittleEndian.PutUint32(buf[offset+4:], bits)
	}

	return frames * frameSize, nil
}

// startPending moves the keypresses waiting at the door into voices,
// and records how long they waited.
func (s *stream) startPending(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.nPending {
		delay := now.Sub(s.pending[i])
		s.lastDelay = delay
		if delay > s.worst {
			s.worst = delay
		}

		for v := range s.voices {
			if !s.voices[v].active {
				s.voices[v] = voice{
					active:    true,
					remaining: int(blipLength.Seconds() * sampleRate),
				}
				break
			}
		}
	}
	s.nPending = 0
}

// next advances the voice by one sample and returns its contribution.
func (v *voice) next() float64 {
	if !v.active {
		return 0
	}

	total := int(blipLength.Seconds() * sampleRate)
	ramp := int(rampLength.Seconds() * sampleRate)
	elapsed := total - v.remaining

	gain := 1.0
	switch {
	case elapsed < ramp:
		gain = float64(elapsed) / float64(ramp)
	case v.remaining < ramp:
		gain = float64(v.remaining) / float64(ramp)
	}

	out := math.Sin(v.phase) * gain
	v.phase += 2 * math.Pi * blipFrequency / sampleRate
	if v.phase > 2*math.Pi {
		v.phase -= 2 * math.Pi
	}

	v.remaining--
	if v.remaining <= 0 {
		v.active = false
	}
	return out
}

func clip(x float64) float64 {
	switch {
	case x > 1:
		return 1
	case x < -1:
		return -1
	}
	return x
}

type game struct {
	stream *stream
	player *audio.Player

	bufferIndex int

	metronome     bool
	lastMetronome time.Time
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.stream.Trigger(time.Now())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.metronome = !g.metronome
		g.lastMetronome = time.Now()
	}

	for i := range bufferSizes {
		key := ebiten.Key1 + ebiten.Key(i)
		if inpututil.IsKeyJustPressed(key) {
			g.bufferIndex = i
			g.player.SetBufferSize(bufferSizes[i])
			g.stream.ResetWorst()
		}
	}

	if g.metronome {
		if now := time.Now(); now.Sub(g.lastMetronome) >= metronomePeriod {
			g.lastMetronome = now
			g.stream.Trigger(now)
		}
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	last, worst := g.stream.Delays()

	metronome := "off"
	if g.metronome {
		metronome = "on"
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"buffer asked for : %v   (keys 1-5)\n"+
			"queueing delay   : %5.1f ms   worst %5.1f ms\n"+
			"floor estimate   : %5.1f ms   = buffer + queueing\n"+
			"metronome        : %s   (m)\n"+
			"\n"+
			"space to blip, q to quit\n"+
			"\n"+
			"The figure is a floor. Your sound card adds more and\n"+
			"cannot be measured from here. Turn the metronome on and\n"+
			"tap along: lag is obvious against a pulse and nearly\n"+
			"invisible on a lone click.\n"+
			"\n"+
			"Go down the buffer sizes until you hear crackling, then\n"+
			"come back up one. That is the machine's limit.",
		bufferSizes[g.bufferIndex],
		float64(last.Microseconds())/1000,
		float64(worst.Microseconds())/1000,
		float64(bufferSizes[g.bufferIndex].Microseconds())/1000+
			float64(worst.Microseconds())/1000,
		metronome,
	))
}

func (g *game) Layout(int, int) (int, int) { return 460, 260 }

func main() {
	s := &stream{}

	context := audio.NewContext(sampleRate)
	player, err := context.NewPlayerF32(s)
	if err != nil {
		log.Fatal(err)
	}

	player.SetBufferSize(bufferSizes[0])
	player.Play()

	ebiten.SetWindowSize(920, 520)
	ebiten.SetWindowTitle("latence audio")

	g := &game{stream: s, player: player}
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
