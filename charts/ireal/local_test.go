package ireal

import (
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
			unknown := 0
			for _, s := range p.Songs {
				ts := Lex(s.Chart)
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
			t.Logf("%s: %q, %d songs, %d unknown tokens", filepath.Base(file), p.Name, len(p.Songs), unknown)
		}
	}
}

func around(s string, i int) string {
	return s[max(0, i-12):min(len(s), i+12)]
}
