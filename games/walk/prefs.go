package main

import "github.com/ArnaudCalmettes/gohar/games/settings"

// prefsFile keeps what the options set, between two runs.
const prefsFile = "walk.json"

// prefs are the options a player sets and finds again: the language
// and the game's tempo. They are the flags' defaults, so that a flag
// given on the command line still wins.
type prefs struct {
	Lang string
	BPM  float64
}

// defaultBPM is the game's tempo before any is set.
const defaultBPM = 120

// loadPrefs reads the options saved, the defaults for any missing: the
// session's language, `defaultBPM`.
func loadPrefs(system string) (prefs, error) {
	p := prefs{Lang: system, BPM: defaultBPM}
	_, err := settings.Load(settings.Path(prefsFile), &p)
	return p, err
}

func (p prefs) save() error {
	return settings.Save(settings.Path(prefsFile), p)
}
