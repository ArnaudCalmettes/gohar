package ireal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// The charts in testdata/local are exports from the app that are not
// ours to publish: the directory is ignored by git, and this test only
// runs where someone put some there. It checks what a smoke test can:
// every playlist reads, every chart lexes whole, with no unknown token.
func TestLocalCharts(t *testing.T) {
	files, _ := filepath.Glob("testdata/local/*.html")
	if len(files) == 0 {
		t.Skip("no charts in testdata/local")
	}
	href := regexp.MustCompile(`href="(irealb://[^"]*)"`)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range href.FindAllSubmatch(data, -1) {
			p, err := Parse(string(m[1]))
			if err != nil {
				t.Errorf("%s: %v", file, err)
				continue
			}
			unknown, measures := 0, 0
			var odds []string
			for _, s := range p.Songs {
				ts := Lex(s.Chart)
				ms := Structure(ts).Measures
				measures += len(ms)
				for _, m := range ms {
					if !m.filled() {
						t.Errorf("%s, %q: an empty measure", file, s.Title)
					}
				}
				if u := (Chart{ms}).Unfold(); len(u) < len(ms) || len(u) >= 16*len(ms)+64 {
					t.Errorf("%s, %q: %d measures played for %d written", file, s.Title, len(u), len(ms))
				}
				for _, odd := range oddities(ms) {
					odd = fmt.Sprintf("%s, %q: %s", filepath.Base(file), s.Title, odd)
					odds = append(odds, odd)
				}
				if join(ts) != s.Chart {
					t.Errorf("%s: round trip lost %q", file, s.Title)
				}
				for _, tok := range ts {
					if tok.Kind == Unknown {
						unknown++
						if unknown <= 10 {
							t.Errorf("%s, %q: unknown %q in %q", file, s.Title, tok.Raw, around(s.Chart, tok.Pos))
						}
					}
				}
			}
			t.Logf("%s: %q, %d songs, %d measures, %d unknown tokens, %d oddities", filepath.Base(file), p.Name, len(p.Songs), measures, unknown, len(odds))
			for _, odd := range odds {
				t.Log("  " + odd)
			}
		}
	}
}

// oddities are what a chart may well hold, written that way by whoever
// typed it, but what a reader of the form will have to cope with: a
// repeat never closed, endings out of order. They are logged, not
// failed: the charts are not ours to fix.
func oddities(ms []Measure) []string {
	var problems []string
	if len(ms) == 0 {
		return []string{"no measure"}
	}
	open, ending := false, 0
	for i, m := range ms {
		if m.Open == RepeatOpen {
			open, ending = true, 0
		}
		if m.Ending > 0 {
			if m.Ending != ending+1 {
				problems = append(problems, fmt.Sprintf("measure %d: ending %d after ending %d", i+1, m.Ending, ending))
			}
			ending = m.Ending
		}
		if m.Close == RepeatClose {
			open = false
		}
	}
	if open {
		problems = append(problems, "a repeat is never closed")
	}
	return problems
}

func around(s string, i int) string {
	return s[max(0, i-12):min(len(s), i+12)]
}
