package settings

import (
	"os"
	"path/filepath"
	"testing"
)

type prefs struct {
	Lang string
	BPM  int
}

func TestSaveThenLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "prefs.json")
	if err := Save(path, prefs{"fr", 100}); err != nil {
		t.Fatal(err)
	}
	var got prefs
	if found, err := Load(path, &got); !found || err != nil || got != (prefs{"fr", 100}) {
		t.Errorf("read back %+v, found %v, %v", got, found, err)
	}
}

// The first run: nothing saved, the defaults stay.
func TestLoadMissing(t *testing.T) {
	got := prefs{Lang: "en"}
	if found, err := Load(filepath.Join(t.TempDir(), "none.json"), &got); found || err != nil || got.Lang != "en" {
		t.Errorf("found %v, %v, %+v", found, err, got)
	}
}

// A damaged file is an error, never a fresh start.
func TestLoadDamaged(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path, &prefs{}); err == nil {
		t.Error("a damaged file read without error")
	}
}
