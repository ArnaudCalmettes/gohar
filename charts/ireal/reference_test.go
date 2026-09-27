package ireal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// A fiche is the analysis of a tune as a book prints it, transcribed by
// hand into testdata/fiches: the oracle the analysis of charts will be
// held to. See docs/grilles.md. Below, "the book" is the one a fiche
// names in Source, En Harmonie for every fiche so far.
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

// The book's brackets: a two five resolving on a later bar. The five,
// the last chord before that bar, must approach its first chord, and
// the two, the chord before the five, must approach the five.
func TestFichesBrackets(t *testing.T) {
	read := func(c ChordSymbol) harmony.Chord {
		r, err := c.Read()
		if err != nil {
			t.Fatal(err)
		}
		return r.Chord()
	}
	for name, f := range fiches(t) {
		for _, br := range f.Toward {
			var before []ChordSymbol
			for _, bar := range f.Bars[br[0]-1 : br[1]-1] {
				before = append(before, chords(bar)...)
			}
			target := chords(f.Bars[br[1]-1])
			if len(before) < 2 || len(target) == 0 {
				t.Errorf("%s, bracket %v: no two five to read", name, br)
				continue
			}
			two, five := before[len(before)-2], before[len(before)-1]
			if k := harmony.ApproachOf(read(five), read(target[0])); k == harmony.NoApproach {
				t.Errorf("%s, bracket %v: %v does not approach %v", name, br, five, target[0])
			}
			if k := harmony.ApproachOf(read(two), read(five)); !k.Has(harmony.TwoApproach) {
				t.Errorf("%s, bracket %v: %v is not the two of %v", name, br, two, five)
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

// fromFiche builds the changes of a fiche's bars, four beats to a bar
// shared between its chords, "%" letting the chord before go on. The
// changes loop when the fiche holds the whole form.
func fromFiche(t *testing.T, f fiche) analysis.Changes {
	c := analysis.Changes{Loops: len(f.Bars) == f.Length}
	bar := 4 * TicksPerBeat
	for i, b := range f.Bars {
		start := Ticks(i) * bar
		c.Bars = append(c.Bars, start)
		if b == "%" {
			c.Chords[len(c.Chords)-1].Length += bar
			continue
		}
		cs := chords(b)
		for j, sym := range cs {
			r, err := sym.Read()
			if err != nil {
				t.Fatal(err)
			}
			ch := analysis.Change{Chord: r.Chord(), Bass: r.Chord().Root, Length: bar / Ticks(len(cs))}
			ch.Start = start + Ticks(j)*ch.Length
			c.Chords = append(c.Chords, ch)
		}
	}
	return c
}

// The degrees the analysis reads bracketed, against those the book
// prints, bracketed too. They need not all agree: the book gives one
// reading where the analysis may give another that is also true, a two
// five that does not resolve where the book hears a borrowed chord. So
// the test only logs where they part, and how often they agree: a
// number to watch as the analysis grows.
//
// Two counts: the degrees, numeral and inversion, and the notation as
// printed. The book writes a borrowed quality once and then leaves it
// out (A♭7 is IV7 in bar 2 of Tenderly, IV in bar 4); the analysis
// writes it each time.
func TestFichesDegrees(t *testing.T) {
	for name, f := range fiches(t) {
		home, ok := Song{Key: f.Key}.HomeTonalities()
		if !ok {
			t.Fatalf("%s: key %q", name, f.Key)
		}
		c := fromFiche(t, f)
		blocks := analysis.Blocks(c, analysis.Approaches(c))
		got := analysis.Bracketed(c, blocks, analysis.PassingChords(c), analysis.Sense(c, blocks, home))
		var want []string
		for _, d := range f.Degrees {
			if d != "%" {
				want = append(want, strings.Fields(d)...)
			}
		}
		if len(want) != len(got) {
			t.Fatalf("%s: %d degrees in the book, %d chords", name, len(want), len(got))
		}
		degrees, notation := 0, 0
		var parts []string
		for i := range want {
			if got[i].String() == want[i] {
				notation++
			}
			if numeral(got[i].String()) == numeral(want[i]) {
				degrees++
				continue
			}
			bar := 1 + int(c.Chords[i].Start/(4*TicksPerBeat))
			parts = append(parts, fmt.Sprintf("bar %d: %s for %s", bar, got[i], want[i]))
		}
		t.Logf("%s: %d degrees of %d as the book, %d written the same", name, degrees, len(want), notation)
		for _, p := range parts {
			t.Log("  " + p)
		}
	}
}

var degreeParts = regexp.MustCompile(`^([♭♯]?[IV]+)[^/]*(/\d)?$`)

// numeral keeps of a degree its numeral, with its accidental, and its
// inversion: "♭VII/3" of "♭VIIm7/3".
func numeral(d string) string {
	m := degreeParts.FindStringSubmatch(d)
	if m == nil {
		return d
	}
	return m[1] + m[2]
}
