// Command ear trains the ear: something sounds over its tonic, the
// player names it. A menu offers the activities, all open: the four
// tetrachords, then the seven modes of the natural system.
//
// The tracer bullet of the project. It crosses every layer, synth,
// keyboard, harmony, naming in two languages and the dex, with as
// little game as will do. What it is for is the list of what does not
// hold; the rules are in docs/OREILLE.md.
//
//	1 to 7, click  pick an activity, answer
//	r              listen again
//	space          next, then the same activity again
//	enter          back to the menu, at the end of a series
//	h              delay figures, for the test under load
//	l              switch between French and English
//	n              switch between signs and words, si♭ or si bémol
//	escape         quit
//
// -timbre picks the sound. The console timbres are smoothed unless
// -authentic asks for their raw aliasing.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
	"github.com/ArnaudCalmettes/gohar/synth"
)

func main() {
	lang := flag.String("lang", "fr", "fr or en")
	notation := flag.String("notation", "signs", "signs or words: si\u266d or si bémol")
	useMIDI := flag.Bool("midi", true, "listen to a MIDI keyboard if there is one")
	port := flag.String("port", "", "an input's number or part of its name; the first one if empty")
	device := flag.Duration("device", synth.DefaultBuffer, "device buffer")
	player := flag.Duration("player", 0, "player buffer, zero means the same")
	a4 := flag.Float64("a4", 440, "diapason in hertz")
	timbre := flag.String("timbre", "sine", "sine, pulse12, pulse25, square, triangle or noise")
	authentic := flag.Bool("authentic", false, "keep the consoles' aliasing and stepped triangle instead of smoothing them")
	dexPath := flag.String("dex", defaultDexPath(), "where the collection is kept")
	seed := flag.Uint64("seed", 0, "random seed, zero for a new one each run")
	flag.Parse()

	l, err := newLanguage(*lang)
	if err != nil {
		log.Fatal(err)
	}
	otherCode := "en"
	if *lang == "en" {
		otherCode = "fr"
	}
	other, err := newLanguage(otherCode)
	if err != nil {
		log.Fatal(err)
	}
	switch *notation {
	case "signs":
	case "words":
		l, other = l.withNotation(naming.Words), other.withNotation(naming.Words)
	default:
		log.Fatalf("unknown notation %q, want signs or words", *notation)
	}

	d, err := loadDex(*dexPath)
	if err != nil {
		log.Fatalf("reading %s: %v", *dexPath, err)
	}

	engine, err := synth.NewEngine(*timbre, !*authentic)
	if err != nil {
		log.Fatal(err)
	}
	engine.Tuning = synth.Tuning{A4: *a4}
	out, err := synth.Open(engine, synth.Options{Device: *device, Player: *player})
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	s := *seed
	if s == 0 {
		s = uint64(time.Now().UnixNano())
	}
	rng := rand.New(rand.NewPCG(s, s>>32|1))

	// The game must run without a keyboard: in a browser there may be
	// none, and answering is done with the mouse anyway.
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

	g, err := newGame(l, other, engine, d, *dexPath, rng, midiName)
	if err != nil {
		log.Fatal(err)
	}
	if keys != nil {
		if err := keys.Listen(g.onKey); err != nil {
			fmt.Fprintln(os.Stderr, "MIDI:", err)
			g.midi = ""
		}
	}

	ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)
	ebiten.SetWindowTitle("gohar")
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
	if err := out.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "audio:", err)
	}
}
