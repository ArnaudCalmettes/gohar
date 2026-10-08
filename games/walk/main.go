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
// The sounds come from GeneralUser GS, slimmed down to what the game
// plays and embedded (sounds/walk.sf2, made by `make slim`); -chip
// plays the 8-bit sounds of ear instead.
//
//	go run ./walk                  # the title screen, then play: the space bar starts
//	go run ./walk -demo            # the band plays the bass itself
//	go run ./walk -practice        # without tempo, T switches back
//	go run ./walk -record line.txt # write down the notes heard
//	go run ./walk -lang en         # in English; the session's language by default
//	go run ./walk -list            # the presets of the soundfont
//	go run ./walk -sf2 ~/.cache/gohar/GeneralUser-GS.sf2 -bass 0:33
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ArnaudCalmettes/gohar/games/calibrate"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/lang"
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/settings"
	"github.com/ArnaudCalmettes/gohar/games/walk/band"
	"github.com/ArnaudCalmettes/gohar/games/walk/grids"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
	"github.com/ArnaudCalmettes/gohar/synth"
)

func main() {
	// The options saved are the flags' defaults: a flag given wins.
	saved, err := loadPrefs(lang.System())
	if err != nil {
		fmt.Fprintln(os.Stderr, "options:", err)
	}
	bpm := flag.Float64("bpm", 0, "tempo of every grid, in beats per minute; 0 for each grid's own")
	choruses := flag.Int("choruses", 2, "choruses to play after the count-in")
	demo := flag.Bool("demo", false, "the band walks the bass itself, the reference line")
	practicing := flag.Bool("practice", false, "start without tempo: the chart waits for the roots")
	sf2 := flag.String("sf2", "", "the soundfont of the band: the bass, the piano, the drums; the one embedded if empty")
	chip := flag.Bool("chip", false, "play 8-bit sounds rather than the soundfont")
	list := flag.Bool("list", false, "list the presets of the soundfont, and quit")
	bassPreset := flag.String("bass", "0:32", "bank:patch of the bass in the soundfont")
	split := flag.Int("split", mark.FirstPalier.Split, "the lowest key of the right hand")
	useMIDI := flag.Bool("midi", true, "listen to a MIDI keyboard if there is one")
	port := flag.String("port", "", "an input's number or part of its name; the first one if empty")
	device := flag.Duration("device", synth.DefaultBuffer, "device buffer")
	gain := flag.Float64("gain", 0, "in the browser, the master gain of the band; 0 for its default")
	record := flag.String("record", "", "write down the notes the marker hears, a line each, in this file")
	language := flag.String("lang", saved.Lang, "the language of the game: fr or en")
	flag.Parse()

	bd, closeAudio, err := band.Open(band.Flags{
		SF2: *sf2, Chip: *chip, List: *list,
		Bass: preset(*bassPreset), Split: *split, Demo: *demo, Device: *device, Gain: *gain,
	})
	if err != nil {
		log.Fatal(err)
	}
	if bd == nil {
		return // -list has listed the presets
	}
	defer closeAudio()

	tunes, err := grids.All()
	if err != nil {
		log.Fatal(err)
	}
	fs, err := newFonts()
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

	l, err := newLang(*language)
	if err != nil {
		log.Fatal(err)
	}
	a := &app{
		band:     bd,
		events:   make(chan keyboard.Event, eventBuffer),
		midi:     midiName,
		lang:     l,
		fonts:    fs,
		tunes:    tunes,
		bpm:      *bpm,
		choruses: *choruses,
		untimed:  *practicing,
	}
	if *record != "" {
		a.rec, err = newRecorder(*record)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := a.rec.Close(); err != nil {
				fmt.Fprintln(os.Stderr, "record:", err)
			}
		}()
	}
	if keys != nil {
		if err := keys.Listen(a.onKey); err != nil {
			fmt.Fprintln(os.Stderr, "MIDI:", err)
			a.midi = ""
		}
	}
	// Without a keyboard, the walker plays and the player listens: the
	// demo is on, and D turns it off at rest.
	if a.midi == "" {
		bd.Demo = true
	}
	// oto opens the default output and does not name it: the buffer is
	// what the game knows of it, and what changes its latency.
	a.pair = calibrate.Key(a.midi, fmt.Sprintf("default output, %v", *device))
	if d, ok, err := calibrate.Load(settings.Path(calibrate.File), a.pair); err != nil {
		fmt.Fprintln(os.Stderr, "calibration:", err)
	} else if ok {
		a.latency = d
	}

	ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)
	ebiten.SetWindowTitle(gameTitle)
	d := scene.New(newTitle(a), a.layout)
	err = ebiten.RunGame(d)
	d.Close() // the scenes release what they hold, before the recorder closes
	if err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}

func preset(s string) synth.Preset {
	var p synth.Preset
	if _, err := fmt.Sscanf(s, "%d:%d", &p.Bank, &p.Patch); err != nil {
		log.Fatalf("preset %q: want bank:patch, as 0:32", s)
	}
	return p
}
