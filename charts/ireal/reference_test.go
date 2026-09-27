package ireal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// A fiche is the analysis of a tune as a book prints it, transcribed by
// hand into testdata/fiches: the oracle the analysis of charts will be
// held to. See docs/grilles.md.
//
// Bars are numbered as in the book, from 1. Chords are written in the
// app's spelling so that the same reader reads both; "%" is a bar where
// the chord before goes on.
type fiche struct {
	Title  string `json:"title"` // as in the app's playlists
	Source string `json:"source"`
	From   int    `json:"from"` // the played bar the book's bar 1 is, from 1

	Bars    []string `json:"bars"`
	Degrees []string `json:"degrees"` // as printed, one per chord

	Structure string `json:"structure"`
	Length    int    `json:"length"` // bars in the form
	Rhythm    string `json:"rhythm"` // chords per bar: "1", "1-2"
	Key       string `json:"key"`    // in the app's spelling: Eb, A-

	Modulations []struct {
		Key  string `json:"key"`
		Bars [2]int `json:"bars"`
	} `json:"modulations"`
	Borrowings []struct {
		Bars   [2]int `json:"bars"`
		Chords string `json:"chords"`
		Scale  string `json:"scale"` // "*" when the book names none
	} `json:"borrowings"`
	Cadences []struct {
		Name string `json:"name"`
		Bars [2]int `json:"bars"`
	} `json:"cadences"`

	// Toward holds the book's brackets: a two five starting at the
	// first bar that resolves on the second.
	Toward [][2]int `json:"toward"`

	Notes []string `json:"notes"`
}

func fiches(t *testing.T) map[string]fiche {
	files, _ := filepath.Glob("testdata/fiches/*.json")
	if len(files) == 0 {
		t.Fatal("no fiche in testdata/fiches")
	}
	out := map[string]fiche{}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		var f fiche
		if err := json.Unmarshal(data, &f); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		out[filepath.Base(file)] = f
	}
	return out
}

// chords returns the chord symbols written in one bar of a fiche.
func chords(bar string) []ChordSymbol {
	var out []ChordSymbol
	for _, tok := range Lex(bar) {
		if tok.Kind == Chord {
			out = append(out, tok.Chord)
		}
	}
	return out
}

// A fiche must hold together before anything is compared with it.
func TestFichesHoldTogether(t *testing.T) {
	for name, f := range fiches(t) {
		if len(f.Degrees) != len(f.Bars) {
			t.Errorf("%s: %d bars, %d bars of degrees", name, len(f.Bars), len(f.Degrees))
			continue
		}
		for i, bar := range f.Bars {
			n := len(chords(bar))
			if bar == "%" {
				n = 1
			}
			if d := len(strings.Fields(f.Degrees[i])); d != n || (bar == "%") != (f.Degrees[i] == "%") {
				t.Errorf("%s, bar %d: %q under %q", name, i+1, f.Degrees[i], bar)
			}
			for _, c := range chords(bar) {
				if _, err := c.Read(); err != nil {
					t.Errorf("%s, bar %d: %v", name, i+1, err)
				}
			}
		}
	}
}

// localSongs reads the charts in testdata/local by title, the first
// one when a title comes twice.
func localSongs(t *testing.T) map[string]Song {
	files, _ := filepath.Glob("testdata/local/*.html")
	href := regexp.MustCompile(`href="(irealb://[^"]*)"`)
	out := map[string]Song{}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range href.FindAllSubmatch(data, -1) {
			p, err := Parse(string(m[1]))
			if err != nil {
				continue
			}
			for _, s := range p.Songs {
				if _, ok := out[s.Title]; !ok {
					out[s.Title] = s
				}
			}
		}
	}
	return out
}

// The chart in the app and the book's grid must be the same tune bar
// for bar, compared by root and tetrad: the book and the typist choose
// extensions each their own way (D♭7 for D♭7♯11). What the chart and
// the book disagree on beyond the chords, the key given by the app or
// the length played, is logged: it is what an analysis has to cope
// with, not an error.
func TestFichesAgainstCharts(t *testing.T) {
	songs := localSongs(t)
	if len(songs) == 0 {
		t.Skip("no charts in testdata/local")
	}
	for name, f := range fiches(t) {
		s, ok := songs[f.Title]
		if !ok {
			t.Logf("%s: %q is not in testdata/local", name, f.Title)
			continue
		}
		tl := Structure(Lex(s.Chart)).Timeline()
		if s.Key != f.Key {
			t.Logf("%s: the app says %s, the book %s", name, s.Key, f.Key)
		}
		if len(tl.Bars) != f.Length {
			t.Logf("%s: %d bars played, a form of %d in the book", name, len(tl.Bars), f.Length)
		}
		for i, bar := range f.Bars {
			played := f.From - 1 + i
			var got []ChordSymbol
			for _, sp := range tl.Spans {
				if sp.Bar == played && !sp.NoChord {
					got = append(got, sp.Chord)
				}
			}
			want := chords(bar)
			if len(got) != len(want) {
				t.Errorf("%s, bar %d: the chart has %v, the book %q", name, i+1, got, bar)
				continue
			}
			for j := range want {
				if !sameTetrad(got[j], want[j]) {
					t.Errorf("%s, bar %d: the chart has %v, the book %v", name, i+1, got[j], want[j])
				}
			}
		}
	}
}

func sameTetrad(a, b ChordSymbol) bool {
	ra, err := a.Read()
	if err != nil {
		return false
	}
	rb, err := b.Read()
	if err != nil {
		return false
	}
	return ra.Root.Class() == rb.Root.Class() && ra.Pattern.Tetrad() == rb.Pattern.Tetrad()
}
