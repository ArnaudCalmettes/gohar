package synth

import (
	"math"
	"testing"
	"time"
)

var timbres = []Timbre{Sine, Pulse12, Pulse25, Square, Triangle, Noise}

// The audio goroutine never allocates, whatever it plays.
func TestReadDoesNotAllocate(t *testing.T) {
	buf := make([]byte, 256*BytesPerFrame)
	for _, timbre := range timbres {
		for _, smooth := range []bool{false, true} {
			e := &Engine{Timbre: timbre, Smooth: smooth, Envelope: ChipEnvelope}
			allocs := testing.AllocsPerRun(50, func() {
				e.NoteOn(60, 1, time.Now())
				e.NoteOn(64, 1, time.Now())
				_, _ = e.Read(buf)
				e.NoteOff(60, time.Now())
				_, _ = e.Read(buf)
			})
			if allocs != 0 {
				t.Errorf("timbre %d, smooth %v: %v allocations per run", timbre, smooth, allocs)
			}
		}
	}
}

// sweep runs an oscillator over whole periods and returns its samples.
func sweep(timbre Timbre, smooth bool, inc float64, n int) []float64 {
	o := oscillator{inc: inc, lfsr: 1}
	out := make([]float64, n)
	for i := range out {
		out[i] = o.next(timbre, smooth)
	}
	return out
}

// The console triangle climbs in sixteen steps; the smooth one does not.
func TestTriangleSteps(t *testing.T) {
	levels := func(samples []float64) int {
		seen := map[float64]bool{}
		for _, v := range samples {
			seen[v] = true
		}
		return len(seen)
	}
	if n := levels(sweep(Triangle, false, 1.0/320, 3200)); n != 16 {
		t.Errorf("authentic triangle takes %d levels, want 16", n)
	}
	if n := levels(sweep(Triangle, true, 1.0/320, 3200)); n <= 16 {
		t.Errorf("smooth triangle takes only %d levels", n)
	}
}

// A pulse keeps no offset, whatever its width, and the smooth one only
// differs from the raw one around its edges.
func TestPulse(t *testing.T) {
	for _, timbre := range []Timbre{Pulse12, Pulse25, Square} {
		samples := sweep(timbre, false, 1.0/400, 4000)
		var sum float64
		for _, v := range samples {
			sum += v
		}
		if mean := sum / float64(len(samples)); math.Abs(mean) > 1e-6 {
			t.Errorf("timbre %d has a mean of %v", timbre, mean)
		}

		raw, smooth := sweep(timbre, false, 1.0/400, 400), sweep(timbre, true, 1.0/400, 400)
		differ := 0
		for i := range raw {
			if raw[i] != smooth[i] {
				differ++
			}
		}
		if differ == 0 || differ > 8 {
			t.Errorf("timbre %d: smoothing touched %d samples of a period, want a handful", timbre, differ)
		}
	}
}

// The noise register cycles through its states without getting stuck.
func TestNoiseMoves(t *testing.T) {
	samples := sweep(Noise, false, 1.0/100, 4800)
	ups := 0
	for _, v := range samples {
		if v > 0 {
			ups++
		}
	}
	if ups == 0 || ups == len(samples) {
		t.Errorf("noise stuck at one value (%d ups of %d)", ups, len(samples))
	}
}

// The envelope settles on the sustain level, and a released voice ends.
func TestEnvelope(t *testing.T) {
	e := &Engine{Envelope: ChipEnvelope}
	buf := make([]byte, 480*BytesPerFrame) // 10 ms

	e.NoteOn(69, 1, time.Now())
	for range 30 { // 300 ms, past attack and decay
		_, _ = e.Read(buf)
	}
	if v := e.voices[0]; v.state != voiceSustain || math.Abs(v.env-ChipEnvelope.Sustain) > 1e-9 {
		t.Errorf("after the decay: state %d, level %v, want sustain at %v", v.state, v.env, ChipEnvelope.Sustain)
	}

	e.NoteOff(69, time.Now())
	for range 10 { // 100 ms, past the release
		_, _ = e.Read(buf)
	}
	if e.voices[0].state != voiceOff {
		t.Errorf("a released voice is still sounding: state %d", e.voices[0].state)
	}
}
