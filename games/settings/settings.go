// Package settings keeps what the games remember between two runs, as
// JSON files in one folder of the user's configuration: the dex of ear,
// the calibrated latency shared by every game (games/calibrate), and
// soon the language (see "Les scènes" in docs/architecture.md).
//
// The desktop keeps files; the browser will keep local storage behind
// the same two calls.
package settings

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// Path is where the file `name` lives: gohar's folder in the user's
// configuration, or the working directory when there is none.
func Path(name string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return name
	}
	return filepath.Join(dir, "gohar", name)
}

// Load reads the JSON file at `path` into `v`, and says whether there
// was one. A missing file leaves `v` as it is, the first run. A file
// that exists and does not read is an error, never a fresh start:
// silently replacing what was saved would lose it.
func Load(path string, v any) (bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(data, v)
}

// Save writes `v` to `path` as JSON, through a temporary file, so that a
// crash halfway leaves the previous one intact.
func Save(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), dirMode); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, fileMode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

const (
	dirMode  = 0o755
	fileMode = 0o644
)
