package calibrate

import (
	"path/filepath"
	"testing"
	"time"
)

// Two keyboards on the same output keep an offset each; a pair never
// measured has none.
func TestSaveThenLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	piano := Key("Digital Piano", "default output, 10ms")
	pads := Key("MPK mini", "default output, 10ms")
	if err := Save(path, piano, 42*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if err := Save(path, pads, -8*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		key  string
		want time.Duration
		ok   bool
	}{
		{piano, 42 * time.Millisecond, true},
		{pads, -8 * time.Millisecond, true},
		{Key("Digital Piano", "default output, 40ms"), 0, false},
	} {
		got, ok, err := Load(path, c.key)
		if err != nil || got != c.want || ok != c.ok {
			t.Errorf("%s: got %v, %v, %v; want %v, %v", c.key, got, ok, err, c.want, c.ok)
		}
	}
}

// A first run, without a file, has no offset and no error.
func TestNoFileYet(t *testing.T) {
	_, ok, err := Load(filepath.Join(t.TempDir(), File), Key("Digital Piano", "default output, 10ms"))
	if ok || err != nil {
		t.Errorf("got %v, %v; want no offset, no error", ok, err)
	}
}
