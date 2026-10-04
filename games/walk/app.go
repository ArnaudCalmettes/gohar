package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/lang"
)

// eventBuffer is how many key events wait for the scene on top.
const eventBuffer = 64

// app is what the scenes share, built once by main and handed to each
// scene it builds (see "Les scènes" in docs/architecture.md): the band,
// which plays across the scenes, the keyboard, the language, the fonts,
// and the options of the command line until a settings screen takes
// them over.
type app struct {
	band   *band
	events chan keyboard.Event // the keys, for the scene on top
	midi   string              // the keyboard's name, empty without one
	lang   *lang.Lang
	fonts  fonts
	scale  float64   // the window's, set by layout
	rec    *recorder // nil without -record
	jam    *jam      // the music of the menus, nil while a game plays

	// pair names the keyboard and the output, for the calibration;
	// latency is the offset measured for it, zero until one is.
	pair    string
	latency time.Duration

	// The grids, and the one chosen last, kept for the session.
	tunes   []tune
	current int

	// The options.
	bpm      float64
	choruses int
	untimed  bool // a game starts without tempo
}

// onKey is the MIDI callback: it sounds the key at once and hands the
// event to the scene on top, never waiting. A full channel drops the
// event for the display only: the sound has already gone.
func (a *app) onKey(e keyboard.Event) {
	a.band.key(e)
	select {
	case a.events <- e:
	default:
	}
}

// drain hands each key event waiting to `f`, if not nil: a scene that
// ignores the keys still drains them, so that the next one does not
// receive them late.
func (a *app) drain(f func(keyboard.Event)) {
	for {
		select {
		case e := <-a.events:
			if f != nil {
				f(e)
			}
		default:
			return
		}
	}
}

// layout scales the logical screen to the window, for every scene.
func (a *app) layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	w, h := float64(outsideWidth)*s, float64(outsideHeight)*s
	a.scale = min(w/screenWidth, h/screenHeight)
	return int(screenWidth * a.scale), int(screenHeight * a.scale)
}
