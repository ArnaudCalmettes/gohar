package lang

import (
	"os"
	"slices"
	"testing"
)

func load(t *testing.T, want string) *Lang {
	t.Helper()
	l, err := New(os.DirFS("testdata"), want)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func TestSays(t *testing.T) {
	fr := load(t, "fr-FR")
	for _, tc := range []struct{ got, want string }{
		{fr.T("hello", "Name", "Arnaud"), "Bonjour Arnaud"},
		{fr.N("landed", 1), "1 arrivée posée"},
		{fr.N("landed", 3), "3 arrivées posées"},
		{fr.T("onlyEnglish"), "Only in English"},
		{fr.T("nowhere"), "nowhere"},
	} {
		if tc.got != tc.want {
			t.Errorf("%q, want %q", tc.got, tc.want)
		}
	}
	if fr.Tag() != "fr" {
		t.Errorf("fr-FR speaks %q, want fr", fr.Tag())
	}
}

// A language without a file speaks English.
func TestFallsBack(t *testing.T) {
	if de := load(t, "de"); de.Tag() != "en" || de.T("hello", "Name", "Ada") != "Hello Ada" {
		t.Errorf("German speaks %q: %q", de.Tag(), de.T("hello", "Name", "Ada"))
	}
}

func TestMissing(t *testing.T) {
	ids := []string{"hello", "landed", "onlyEnglish", "nowhere"}
	if got := load(t, "fr").Missing(ids); !slices.Equal(got, []string{"onlyEnglish", "nowhere"}) {
		t.Errorf("missing in French: %v", got)
	}
	if got := load(t, "en").Missing(ids); !slices.Equal(got, []string{"nowhere"}) {
		t.Errorf("missing in English: %v", got)
	}
}

func TestSystem(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "fr_FR.UTF-8")
	if got := System(); got != "fr-FR" {
		t.Errorf("LANG=fr_FR.UTF-8 gives %q", got)
	}
	t.Setenv("LANG", "C")
	if got := System(); got != Fallback {
		t.Errorf("LANG=C gives %q", got)
	}
}
