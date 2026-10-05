package main

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/synth"
)

// audioFlags are what the command line says of the sound. The browser
// reads them too, from the page's address, and ignores those it has no
// use for (see audio_js.go).
type audioFlags struct {
	sf2        string // a soundfont other than the one embedded
	chip, list bool
	bass       synth.Preset
	split      int
	demo       bool
	device     time.Duration // the desktop's audio buffer
	gain       float64       // the browser's master gain, 0 for its default
}
