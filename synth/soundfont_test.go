package synth

import (
	"os"
	"testing"
	"time"
)

// No soundfont ships with the repository: GOHAR_SF2 names one to test
// with, any General MIDI one will do. Without it, these tests skip.
func loadTestSoundFont(t *testing.T) *SoundFont {
	t.Helper()
	path := os.Getenv("GOHAR_SF2")
	if path == "" {
		t.Skip("GOHAR_SF2 names no soundfont")
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sf, err := LoadSoundFont(f)
	if err != nil {
		t.Fatal(err)
	}
	return sf
}

func TestMidiVelocity(t *testing.T) {
	for _, tc := range []struct {
		v    float64
		want int32
	}{{0, 1}, {0.5, 64}, {1, 127}, {2, 127}} {
		if got := midiVelocity(tc.v); got != tc.want {
			t.Errorf("velocity %v: %d, want %d", tc.v, got, tc.want)
		}
	}
}

func TestSamplerUnknownPreset(t *testing.T) {
	sf := loadTestSoundFont(t)
	if _, err := NewSampler(sf, Preset{Bank: 99, Patch: 99}); err == nil {
		t.Error("bank 99, patch 99 accepted")
	}
}

// A note dated 10 ms ahead starts within a block of frame 480. A
// recorded attack may start quietly, so the first sound may come a
// little later still: the bound is loose on that side.
func TestSamplerScheduledNote(t *testing.T) {
	sf := loadTestSoundFont(t)
	p := sf.Presets()[0]
	s, err := NewSampler(sf, p)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 2048*BytesPerFrame)
	s.ScheduleOn(48, 1, t0.Add(10*time.Millisecond))
	s.render(s, buf, t0)

	first := -1
	for f := range len(buf) / BytesPerFrame {
		if samples(buf[f*BytesPerFrame : (f+1)*BytesPerFrame])[0] != 0 {
			first = f
			break
		}
	}
	if first < 480 || first > 480+samplerBlock+480 {
		t.Errorf("%s: first sound on frame %d, want from 480, within a block and an attack", p.Name, first)
	}
}

func TestSamplerDoesNotAllocate(t *testing.T) {
	sf := loadTestSoundFont(t)
	s, err := NewSampler(sf, sf.Presets()[0])
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 256*BytesPerFrame)
	allocs := testing.AllocsPerRun(50, func() {
		now := time.Now()
		s.NoteOn(48, 0.8, now)
		s.ScheduleOff(48, now.Add(2*time.Millisecond))
		_, _ = s.Read(buf)
	})
	if allocs != 0 {
		t.Errorf("%v allocations per run", allocs)
	}
}
