package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ArnaudCalmettes/gohar/dex"
)

// The desktop store: one JSON file. The browser will have its own, in
// local storage; the dex does not care, it only knows its form.

// defaultDexPath is where the collection lives unless told otherwise.
func defaultDexPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "dex.json"
	}
	return filepath.Join(dir, "gohar", "dex.json")
}

// loadDex reads the collection, or starts an empty one when there is
// none yet. A file that exists and does not read is an error, never a
// fresh start: silently replacing a collection would lose it.
func loadDex(path string) (*dex.Dex, error) {
	d := dex.New()
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return d, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, d); err != nil {
		return nil, err
	}
	return d, nil
}

// saveDex writes the collection through a temporary file, so that a
// crash halfway leaves the previous one intact.
func saveDex(path string, d *dex.Dex) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
