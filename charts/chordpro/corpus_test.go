package chordpro

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

var link = regexp.MustCompile(`irealb://[^"\s<]*`)

// Every chart of the private corpus, in charts/ireal/testdata/local,
// converted, written and read again, gives the analysis the same
// changes as the iReal chart itself: the profile loses nothing on the
// way. Skipped where the corpus is not, as the tests of charts/ireal.
func TestCorpusSameChanges(t *testing.T) {
	files, _ := filepath.Glob("../ireal/testdata/local/*.html")
	if len(files) == 0 {
		t.Skip("no charts in ../ireal/testdata/local")
	}
	const shown = 25
	var songs, same, skipped int
	var differ []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range link.FindAllString(string(data), -1) {
			p, err := ireal.Parse(l)
			if err != nil {
				continue
			}
			for _, src := range p.Songs {
				songs++
				want, _ := ireal.Structure(ireal.Lex(src.Chart)).Changes()
				s, err := FromIReal(src)
				if errors.Is(err, ErrTime) {
					skipped++
					continue
				}
				var b bytes.Buffer
				if err := Write(&b, s); err != nil {
					t.Fatal(err)
				}
				again, err := ParseString(b.String())
				if err != nil {
					differ = append(differ, src.Title+": "+err.Error())
					continue
				}
				got, err := again.Changes()
				if err != nil {
					differ = append(differ, src.Title+": "+err.Error())
					continue
				}
				// charts/ireal lets a chord go on when it is written the
				// same, Bb^7 after Bb^7; the profile when it is the same
				// chord, Bbmaj7 after Bb^ too. Merged the same, they
				// must agree.
				g, w := played(t, merged(got)), played(t, merged(want))
				switch {
				case !slices.Equal(g, w):
					i := 0
					for i < min(len(g), len(w)) && g[i] == w[i] {
						i++
					}
					at := func(s []string) string {
						if i < len(s) {
							return s[i]
						}
						return "the end"
					}
					differ = append(differ, src.Title+": "+at(g)+", want "+at(w))
				case got.Coda != 0 && want.Coda != 0 && got.Chords[got.Coda].Start != want.Chords[want.Coda].Start,
					(got.Coda == 0) != (want.Coda == 0):
					differ = append(differ, src.Title+": the coda differs")
				default:
					same++
				}
			}
		}
	}
	t.Logf("%d songs, %d the same, %d skipped for a change of time signature", songs, same, skipped)
	for i, d := range differ {
		if i == shown {
			t.Errorf("… and %d more", len(differ)-shown)
			break
		}
		t.Error(d)
	}
}

// merged lets a chord go on when the next one is the same chord.
func merged(c analysis.Changes) analysis.Changes {
	var out []analysis.Change
	for _, ch := range c.Chords {
		if n := len(out); n > 0 && out[n-1].Silent == ch.Silent && out[n-1].Chord == ch.Chord && out[n-1].Bass == ch.Bass {
			out[n-1].Length += ch.Length
			continue
		}
		out = append(out, ch)
	}
	c.Chords = out
	return c
}
