package calibrate

import (
	"embed"
	"io/fs"

	"github.com/ArnaudCalmettes/gohar/games/lang"
)

// The scene's own phrases, one file per language in locales/: shared by
// the games, it speaks for itself.
//
//go:embed locales/*.toml
var locales embed.FS

func newLang(tag string) (*lang.Lang, error) {
	files, err := fs.Sub(locales, "locales")
	if err != nil {
		return nil, err
	}
	return lang.New(files, tag)
}

// The IDs of the phrases, all in each file of locales/ (TestPhrases).
const (
	msgTitle   = "title"
	msgAsk     = "ask"     // press on the fourth beat
	msgNoMIDI  = "nomidi"  // no keyboard to calibrate
	msgTap     = "tap"     // {{.Ms}}: the last tap
	msgMeasure = "measure" // {{.Ms}}, {{.Taps}}, {{.Of}}: so far
	msgSteady  = "steady"  // {{.Ms}}: the result
	msgSaved   = "saved"   // kept for this keyboard
	msgKeys    = "keys"
	msgDone    = "done" // finished: any key goes back
)

var phrases = []string{msgTitle, msgAsk, msgNoMIDI, msgTap, msgMeasure, msgSteady, msgSaved, msgKeys, msgDone}
