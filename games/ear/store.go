package main

import (
	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/games/settings"
)

// The collection is kept by settings, as JSON; the dex does not care
// where, it only knows its form.

// defaultDexPath is where the collection lives unless told otherwise.
func defaultDexPath() string {
	return settings.Path("dex.json")
}

// loadDex reads the collection, or starts an empty one when there is
// none yet.
func loadDex(path string) (*dex.Dex, error) {
	d := dex.New()
	if _, err := settings.Load(path, d); err != nil {
		return nil, err
	}
	return d, nil
}

// saveDex writes the collection, a crash halfway leaving the previous
// one intact.
func saveDex(path string, d *dex.Dex) error {
	return settings.Save(path, d)
}
