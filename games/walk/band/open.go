package band

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/synth"
)

// Flags are what the command line says of the sound. The browser
// reads them too, from the page's address, and ignores those it has no
// use for (see open_js.go).
type Flags struct {
	SF2        string // a soundfont other than the one embedded
	Chip, List bool
	Bass       synth.Preset
	Split      int
	Demo       bool
	Device     time.Duration // the desktop's audio buffer
	Gain       float64       // the browser's master gain, 0 for its default
}
