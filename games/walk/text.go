package main

import (
	"embed"
	"io/fs"

	"github.com/ArnaudCalmettes/gohar/games/lang"
)

// The game's phrases, one file per language in locales/ (see
// "L'internationalisation" in docs/architecture.md). The musical words,
// note and chord names, come from naming.
//
//go:embed locales/*.toml
var locales embed.FS

// newLang speaks `tag`, "fr" or "en", English for any other.
func newLang(tag string) (*lang.Lang, error) {
	files, err := fs.Sub(locales, "locales")
	if err != nil {
		return nil, err
	}
	return lang.New(files, tag)
}

// The IDs of the phrases. Every one is in each file of locales/, which
// TestPhrases checks: call them through these names only.
const (
	msgTempo     = "tempo"       // {{.BPM}}
	msgFreeTempo = "tempo.free"  // the chart waits for the root
	msgKeys      = "status.keys" // the keys of the game
	msgNoMIDI    = "status.nomidi"
	msgMark      = "status.mark" // {{.Note}}, {{.Words}}: what a note was

	msgModeTempo = "mode.tempo"
	msgModeFree  = "mode.free"
	msgModeDemo  = "mode.demo"
	msgModeHint  = "mode.hint" // {{.Mode}}, at rest

	msgOnTime  = "timing.ontime"
	msgEarly   = "timing.early"
	msgLate    = "timing.late"
	msgBetween = "timing.between"

	msgRoot      = "pitch.root"
	msgChordTone = "pitch.chordtone"
	msgOutside   = "pitch.outside"

	msgPlayer = "record.player" // who plays, in a recording's heading
	msgDemo   = "record.demo"
)

var phrases = []string{
	msgTempo, msgFreeTempo, msgKeys, msgNoMIDI, msgMark,
	msgModeTempo, msgModeFree, msgModeDemo, msgModeHint,
	msgOnTime, msgEarly, msgLate, msgBetween,
	msgRoot, msgChordTone, msgOutside,
	msgPlayer, msgDemo,
}

var (
	timingPhrase = map[Timing]string{OnTime: msgOnTime, Early: msgEarly, Late: msgLate, Between: msgBetween}
	pitchPhrase  = map[Pitch]string{Root: msgRoot, ChordTone: msgChordTone, Outside: msgOutside}
)
