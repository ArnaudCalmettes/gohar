// Command otolatency measures the audio path with nothing on top of it.
//
// # Why a second probe
//
// The games/latency probe drives Ebitengine's audio package, which
// creates its oto context with the sample rate and nothing else. The
// ALSA driver then falls back on its own default, a period of 1024
// frames repeated twice, which is 42.7 ms of device buffer at 48 kHz
// that no call on the player can shorten.
//
// This one creates the oto context itself and sets BufferSize, which is
// the only way to move that floor. No window, no engine, no game loop:
// if the delay is still there, it belongs to oto, ALSA or the card, and
// not to anything we wrote.
//
// # The two buffers, and why both have to be named
//
// There are two queues between a decision and a sound, and forgetting
// either one is how a probe lies.
//
// The device buffer is what ALSA hands the card, set once through
// NewContextOptions.BufferSize. The player buffer is a ring oto keeps
// above it, filled eagerly the moment Play is called. Its default is
// half a second, which is sane for a music file and absurd for an
// instrument: a blip asked for after startup lands behind half a second
// of silence already queued. So this program always sets it, and never
// lets the default through.
//
// # How to read the figures
//
// queued is how much audio was already waiting when the blip was asked
// for, taken from the player. It is the part the earlier probe could
// not see. delay is the old figure, from the request to the moment the
// first sample is written, and it stays small even when nothing is
// audible for a second, which is exactly how the earlier probe missed
// the problem.
//
// The estimate is queued plus the device buffer. It is still a floor:
// the card and the USB bus add their own and cannot be seen from here.
//
// The measurement that counts needs no program. Hit the Enter key hard.
// Its mechanical click reaches the ear through the air at once, the blip
// arrives later, and the flam between the two is the real latency. One
// sound means it is under the ear's resolution, two means it is not.
package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

const (
	sampleRate    = 48000
	channelCount  = 2
	bytesPerFrame = channelCount * 4

	blipFrequency = 880.0
	blipLength    = 30 * time.Millisecond

	// rampLength avoids the click a square edge makes at both ends of
	// the blip, which would be louder than the blip and would blur the
	// very attack we are listening for.
	rampLength = 2 * time.Millisecond

	// maxVoices is how many blips may overlap. A dropped blip is better
	// than an allocation in the audio goroutine.
	maxVoices = 16
)

// bytes converts a duration of audio into the size it occupies.
func bytes(d time.Duration) int {
	n := int(int64(d) * sampleRate / int64(time.Second))
	return n * bytesPerFrame
}

// span converts a size back into the duration it will take to play.
func span(n int) time.Duration {
	frames := n / bytesPerFrame
	return time.Duration(frames) * time.Second / sampleRate
}

// A voice is one blip being written out.
type voice struct {
	remaining int
	phase     float64
}

// stream is the source oto reads from. Read runs on the audio
// goroutine and allocates nothing; Trigger runs on the main goroutine.
// The mutex covers the handover between the two and nothing else, so
// the voices themselves stay private to Read.
type stream struct {
	mu       sync.Mutex
	pending  [maxVoices]time.Time
	nPending int
	last     time.Duration
	worst    time.Duration

	voices [maxVoices]voice
}

// Trigger asks for a blip and records when it was asked for.
func (s *stream) Trigger(at time.Time) {
	s.mu.Lock()
	if s.nPending < len(s.pending) {
		s.pending[s.nPending] = at
		s.nPending++
	}
	s.mu.Unlock()
}

// Delays returns the last and the worst write delay seen.
func (s *stream) Delays() (last, worst time.Duration) {
	s.mu.Lock()
	last, worst = s.last, s.worst
	s.mu.Unlock()
	return
}

func (s *stream) Read(buf []byte) (int, error) {
	now := time.Now()

	s.mu.Lock()
	for i := range s.nPending {
		d := now.Sub(s.pending[i])
		s.last = d
		if d > s.worst {
			s.worst = d
		}
		s.wake()
	}
	s.nPending = 0
	s.mu.Unlock()

	frames := len(buf) / bytesPerFrame
	total := bytes(blipLength) / bytesPerFrame
	ramp := bytes(rampLength) / bytesPerFrame
	step := 2 * math.Pi * blipFrequency / sampleRate

	for f := range frames {
		var sample float64
		for v := range s.voices {
			vo := &s.voices[v]
			if vo.remaining == 0 {
				continue
			}
			done := total - vo.remaining
			sample += 0.25 * math.Sin(vo.phase) * envelope(done, vo.remaining, ramp)
			vo.phase += step
			vo.remaining--
		}

		write(buf[f*bytesPerFrame:], float32(sample))
	}

	return frames * bytesPerFrame, nil
}

// wake starts a voice on the first free slot, dropping the blip when
// every slot is busy.
func (s *stream) wake() {
	total := bytes(blipLength) / bytesPerFrame
	for i := range s.voices {
		if s.voices[i].remaining == 0 {
			s.voices[i] = voice{remaining: total}
			return
		}
	}
}

// envelope is a linear fade in and out, flat in between.
func envelope(done, remaining, ramp int) float64 {
	switch {
	case done < ramp:
		return float64(done) / float64(ramp)
	case remaining < ramp:
		return float64(remaining) / float64(ramp)
	}
	return 1
}

func write(buf []byte, v float32) {
	bits := math.Float32bits(v)
	binary.LittleEndian.PutUint32(buf[0:], bits)
	binary.LittleEndian.PutUint32(buf[4:], bits)
}

func main() {
	device := flag.Duration("device", 20*time.Millisecond,
		"buffer handed to ALSA, the one Ebitengine cannot set")
	player := flag.Duration("player", 0,
		"ring buffer above the device; zero means one device buffer, "+
			"never oto's half second default")
	bpm := flag.Float64("bpm", 0,
		"metronome tempo, zero for no metronome")
	flag.Parse()

	if *player == 0 {
		*player = *device
	}

	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channelCount,
		Format:       oto.FormatFloat32LE,
		BufferSize:   *device,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "oto:", err)
		os.Exit(1)
	}
	<-ready

	s := &stream{}
	p := ctx.NewPlayer(s)
	p.SetBufferSize(bytes(*player))
	p.Play()

	fmt.Printf("device buffer : %v\n", *device)
	fmt.Printf("player buffer : %v\n", *player)
	fmt.Println()
	fmt.Println("Enter to blip, q then Enter to quit.")
	fmt.Println("Hit the key hard and listen for a flam between the click")
	fmt.Println("and the blip. The flam is the latency; the figures are a floor.")
	fmt.Println()

	if *bpm > 0 {
		period := time.Duration(float64(time.Minute) / *bpm)
		go func() {
			for t := range time.Tick(period) {
				s.Trigger(t)
			}
		}()
		fmt.Printf("metronome at %.0f bpm\n\n", *bpm)
	}

	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		if strings.TrimSpace(in.Text()) == "q" {
			return
		}

		queued := span(p.BufferedSize())
		s.Trigger(time.Now())

		// Let the audio goroutine pick the blip up, so the delay
		// printed belongs to it rather than to the one before.
		time.Sleep(*device + 5*time.Millisecond)

		last, worst := s.Delays()
		fmt.Printf("queued %5.1f   delay %5.1f   worst %5.1f   estimate %5.1f ms\n",
			ms(queued), ms(last), ms(worst), ms(queued+*device))
	}
}

func ms(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}
