// Command walk is Walk With Me, the game that teaches the left hand to
// replace the bassist (docs/walk.md).
//
// The first palier: the chart of a jazz blues, two bars of hi-hat to
// count in, then the ride while the player plays the roots at the MIDI
// keyboard, and the arrivals marked on the chart. The finger snaps on 2
// and 4 come when the walker snaps, once it has been rolling. The
// left hand sounds like a double bass, the right like a piano, split at
// sol2 (G3). Without tempo, the chart waits for each root instead.
//
// The sounds come from GeneralUser GS, which `make sounds` downloads
// into the user's cache directory; without it, the chip sounds of ear.
//
//	go run ./walk                  # play, the space bar starts
//	go run ./walk -demo            # the band plays the bass itself
//	go run ./walk -practice        # without tempo, T switches back
//	go run ./walk -record line.txt # write down the notes heard
//	go run ./walk -list            # the presets of the soundfont
//	go run ./walk -sf2 other.sf2 -bass 0:33
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/synth"
)

func main() {
	bpm := flag.Float64("bpm", 100, "tempo, in beats per minute")
	choruses := flag.Int("choruses", 2, "choruses to play after the count-in")
	demo := flag.Bool("demo", false, "the band walks the bass itself, the reference line")
	practicing := flag.Bool("practice", false, "start without tempo: the chart waits for the roots")
	sf2 := flag.String("sf2", defaultSoundFont(), "the soundfont to play the bass and the piano with")
	chip := flag.Bool("chip", false, "play 8-bit sounds rather than the soundfont")
	list := flag.Bool("list", false, "list the presets of the soundfont, and quit")
	bassPreset := flag.String("bass", "0:32", "bank:patch of the bass in the soundfont")
	split := flag.Int("split", FirstPalier.Split, "the lowest key of the right hand")
	useMIDI := flag.Bool("midi", true, "listen to a MIDI keyboard if there is one")
	port := flag.String("port", "", "an input's number or part of its name; the first one if empty")
	device := flag.Duration("device", synth.DefaultBuffer, "device buffer")
	record := flag.String("record", "", "write down the notes the marker hears, a line each, in this file")
	flag.Parse()

	if _, err := os.Stat(*sf2); err != nil && !*chip {
		log.Printf("no soundfont at %s: 8-bit sounds instead (run make sounds)", *sf2)
		*chip = true
	}
	var sf *synth.SoundFont
	if !*chip {
		sf = load(*sf2)
		if *list {
			for _, p := range sf.Presets() {
				fmt.Printf("%3d:%-3d %s\n", p.Bank, p.Patch, p.Name)
			}
			return
		}
	}

	var mix synth.Mixer
	bd, err := newBand(&mix, sf, preset(*bassPreset), *split, *demo)
	if err != nil {
		log.Fatal(err)
	}
	out, err := synth.Open(&mix, synth.Options{Device: *device})
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	grid, err := readGrid(jazzBlues)
	if err != nil {
		log.Fatal(err)
	}

	// The game runs without a keyboard: the demo needs none.
	var midiName string
	var keys keyboard.Source
	if *useMIDI {
		defer keyboard.Shutdown()
		keys, err = keyboard.OpenMIDI(*port)
		if err != nil {
			fmt.Fprintln(os.Stderr, "MIDI:", err)
			keys = nil
		} else {
			defer keys.Close()
			midiName = keys.Name()
		}
	}

	g, err := newGame("12 Bar Blues", grid, *bpm, *choruses, bd, midiName, *practicing)
	if err != nil {
		log.Fatal(err)
	}
	if *record != "" {
		g.rec, err = newRecorder(*record)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := g.rec.Close(); err != nil {
				fmt.Fprintln(os.Stderr, "record:", err)
			}
		}()
	}
	if keys != nil {
		if err := keys.Listen(g.onKey); err != nil {
			fmt.Fprintln(os.Stderr, "MIDI:", err)
			g.midi = ""
		}
	}

	ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)
	ebiten.SetWindowTitle("Walk With Me")
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
	if err := out.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "audio:", err)
	}
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
