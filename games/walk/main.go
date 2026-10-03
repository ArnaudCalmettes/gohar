// Command walk is Walk With Me, the game that teaches the left hand to
// replace the bassist (docs/walk.md).
//
// For now it only counts: two bars of count-in, then the bars, beat by
// beat, the way the game will schedule them, a window ahead.
//
// Without a soundfont, a tick marks every beat and a noise burst stands
// for the snap on 2 and 4: two instruments on the same beats, so that
// a slip of the mixer or the clocks would be heard as a flam. With
// the soundfont, the first palier as the player will hear it: finger
// snaps on 2 and 4 count in, then a double bass plays the roots of a
// blues in F.
//
// The soundfont is GeneralUser GS, which `make sounds` downloads into
// the user's cache directory; -sf2 names another one.
//
//	go run ./walk -bars 12
//	go run ./walk -list
//	go run ./walk -sf2 other.sf2 -bass 0:32
//	go run ./walk -chip
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ArnaudCalmettes/gohar/synth"
)

// lookahead is how far ahead the sound schedules. Ebiten's loop ticks
// every 16 ms, too coarse to play a beat on time: notes are queued in
// advance and the synth applies them to the sample (see "Le son" in
// docs/walk.md).
const lookahead = 100 * time.Millisecond

const (
	tickKey  = 81 // a high triangle, soft
	noiseKey = 84 // the noise is pitched by the key: high, a short hiss
	hold     = 30 * time.Millisecond
)

// snapWAV is the finger snap, from newagesoup on Freesound, CC0 (see
// sounds/CREDITS.md). Its attack falls 3 ms into the file, well under
// the spread of the humanizing to come.
//
//go:embed sounds/snap.wav
var snapWAV []byte

// bluesInF is the roots of the plainest twelve-bar blues, one per bar,
// in the register of the double bass: F2, B♭1, C2.
var bluesInF = []int{41, 34, 41, 41, 34, 34, 41, 41, 36, 34, 41, 36}

// A beat is what one beat plays.
type beat interface {
	play(n int, p Position, at time.Time) string
}

func main() {
	bpm := flag.Float64("bpm", 120, "tempo, in beats per minute")
	bars := flag.Int("bars", 2, "bars to play after the count-in")
	sf2 := flag.String("sf2", defaultSoundFont(), "the soundfont to play the bass with")
	chip := flag.Bool("chip", false, "play 8-bit sounds rather than the soundfont")
	list := flag.Bool("list", false, "list the presets of the soundfont, and quit")
	bassPreset := flag.String("bass", "0:32", "bank:patch of the bass in -sf2")
	flag.Parse()

	var mix synth.Mixer
	var b beat
	if _, err := os.Stat(*sf2); err != nil && !*chip {
		log.Printf("no soundfont at %s: 8-bit sounds instead (run make sounds)", *sf2)
		*chip = true
	}
	if *chip {
		b = newChip(&mix)
	} else {
		sf := load(*sf2)
		if *list {
			for _, p := range sf.Presets() {
				fmt.Printf("%3d:%-3d %s\n", p.Bank, p.Patch, p.Name)
			}
			return
		}
		b = newBand(&mix, sf, preset(*bassPreset))
	}

	dev, err := synth.Open(&mix, synth.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	const perBar, countIn = 4, 2 * 4
	start := time.Now().Add(lookahead + countIn*time.Duration(float64(time.Minute) / *bpm))
	m := NewMetronome(start, *bpm, perBar)

	next, last := -countIn, *bars*perBar
	for next < last {
		now := time.Now()
		_, end := m.Due(now, now.Add(lookahead))
		for ; next < end && next < last; next++ {
			p := m.Position(next)
			mark := b.play(next, p, m.At(next))
			// Printed when scheduled, so a window ahead of the beat.
			fmt.Printf("%d.%d %s (in %v)\n", p.Bar, p.Beat, mark, time.Until(m.At(next)).Round(time.Millisecond))
		}
		time.Sleep(lookahead / 4)
	}
	time.Sleep(time.Second) // let the last beat ring
}

// chip is the sound without a soundfont: a tick and a noise burst.
type chip struct{ tick, noise *synth.Engine }

func newChip(mix *synth.Mixer) chip {
	tick, err := synth.NewEngine("triangle", true)
	if err != nil {
		log.Fatal(err)
	}
	noise, err := synth.NewEngine("noise", true)
	if err != nil {
		log.Fatal(err)
	}
	mix.Add(tick, 0.5)
	mix.Add(noise, 1)
	return chip{tick, noise}
}

func (c chip) play(_ int, p Position, at time.Time) string {
	c.tick.ScheduleOn(tickKey, 0.8, at)
	c.tick.ScheduleOff(tickKey, at.Add(hold))
	if p.Beat%2 == 0 {
		c.noise.ScheduleOn(noiseKey, 0.8, at)
		c.noise.ScheduleOff(noiseKey, at.Add(hold))
		return "tick snap"
	}
	return "tick"
}

// band is the sound with a soundfont: the roots on the bass, the snap
// from its own recording. The count-in is snaps alone.
type band struct {
	bass *synth.Sampler
	snap *synth.Clip
	root int // the note the bass holds, 0 for none
}

func newBand(mix *synth.Mixer, sf *synth.SoundFont, bass synth.Preset) *band {
	b, err := synth.NewSampler(sf, bass)
	if err != nil {
		log.Fatal(err)
	}
	s, err := synth.LoadClip(bytes.NewReader(snapWAV))
	if err != nil {
		log.Fatal(err)
	}
	// The balance of the renders validated by ear, bass to snap as 4 to
	// 0.6, a little lower overall. meltysynth plays at half scale, hence
	// the bass's large gain.
	mix.Add(b, 3)
	mix.Add(s, 0.45)
	return &band{bass: b, snap: s}
}

func (b *band) play(n int, p Position, at time.Time) string {
	mark := "count"
	if n >= 0 {
		root := bluesInF[(p.Bar-1)%len(bluesInF)]
		// Legato, as Siskind asks: each note holds until the next one,
		// released on the same date it is replaced.
		if b.root != 0 {
			b.bass.ScheduleOff(b.root, at)
		}
		b.bass.ScheduleOn(root, 0.8, at)
		b.root = root
		mark = fmt.Sprintf("bass %d", root)
	}
	// A clip rings out on its own: no release needed.
	if p.Beat%2 == 0 {
		b.snap.ScheduleOn(0, 1, at)
		mark += " snap"
	}
	return mark
}

// defaultSoundFont is where `make sounds` puts GeneralUser GS.
func defaultSoundFont() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "GeneralUser-GS.sf2"
	}
	return filepath.Join(dir, "gohar", "GeneralUser-GS.sf2")
}

func load(path string) *synth.SoundFont {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sf, err := synth.LoadSoundFont(f)
	if err != nil {
		log.Fatal(err)
	}
	return sf
}

func preset(s string) synth.Preset {
	var p synth.Preset
	if _, err := fmt.Sscanf(s, "%d:%d", &p.Bank, &p.Patch); err != nil {
		log.Fatalf("preset %q: want bank:patch, as 0:32", s)
	}
	return p
}
