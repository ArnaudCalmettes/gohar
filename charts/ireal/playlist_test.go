package ireal

import (
	"net/url"
	"strings"
	"testing"
)

// link makes an irealb:// URL from raw fields, scrambling the charts:
// songs are given as their fields, with the chart unscrambled and
// marked by a leading "@".
func link(name string, songs ...[]string) string {
	var parts []string
	for _, fields := range songs {
		f := make([]string, len(fields))
		for i, v := range fields {
			if strings.HasPrefix(v, "@") {
				v = chartPrefix + unscramble(v[1:])
			}
			f[i] = v
		}
		parts = append(parts, strings.Join(f, "="))
	}
	body := strings.Join(parts, songSeparator)
	if name != "" {
		body += songSeparator + name
	}
	return scheme + url.PathEscape(body)
}

const chart = "{*AT44C^7XyQ|A-7XyQ|D-7XyQ|G7XyQ}"

func TestParsePlaylist(t *testing.T) {
	p, err := Parse(link("Mine",
		[]string{"First", "Me", "", "Medium Swing", "C", "0", "@" + chart, "Jazz-Medium Swing", "120", "3"},
		[]string{"Second", "Me", "", "Bossa Nova", "A-", "2", "@" + chart, "", "0", "1"},
	))
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Mine" || len(p.Songs) != 2 {
		t.Fatalf("%q with %d songs", p.Name, len(p.Songs))
	}
	want := Song{"First", "Me", "Medium Swing", "C", 0, chart, "Jazz-Medium Swing", 120, 3}
	if p.Songs[0] != want {
		t.Errorf("got  %+v\nwant %+v", p.Songs[0], want)
	}
	if s := p.Songs[1]; s.Key != "A-" || s.Transpose != 2 || s.CompStyle != "" || s.Repeats != 1 {
		t.Errorf("second song %+v", s)
	}
}

// A single song has no playlist name after it.
func TestParseSingleSong(t *testing.T) {
	p, err := Parse(link("", []string{"Alone", "Me", "", "Waltz", "F", "0", "@" + chart, "", "0", "0"}))
	if err != nil || p.Name != "" || len(p.Songs) != 1 || p.Songs[0].Title != "Alone" {
		t.Fatalf("%+v, %v", p, err)
	}
}

// Older writers dropped fields; each layout is read around the chart.
func TestParseOlderLayouts(t *testing.T) {
	for name, fields := range map[string][]string{
		"no empty third field": {"T", "Me", "Waltz", "F", "0", "@" + chart, "Jazz-Waltz", "90", "2"},
		"no transposition":     {"T", "Me", "Waltz", "F", "@" + chart, "Jazz-Waltz", "90", "2"},
		"no comp style":        {"T", "Me", "Waltz", "F", "0", "@" + chart, "90", "2"},
		"neither":              {"T", "Me", "Waltz", "F", "@" + chart, "90", "2"},
	} {
		p, err := Parse(link("", fields))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		s := p.Songs[0]
		if s.Style != "Waltz" || s.Key != "F" || s.Chart != chart || s.BPM != 90 || s.Repeats != 2 {
			t.Errorf("%s: %+v", name, s)
		}
	}
}

// A plus sign is an augmented chord, not an encoded space.
func TestParseKeepsPlus(t *testing.T) {
	p, err := Parse(link("", []string{"T", "Me", "", "Waltz", "C", "0", "@C+XyQZ", "", "0", "0"}))
	if err != nil || p.Songs[0].Chart != "C+XyQZ" {
		t.Fatalf("%+v, %v", p, err)
	}
}

// Whatever comes before the scheme is ignored: an href copied whole.
func TestParseFromHref(t *testing.T) {
	href := `<a href="` + link("", []string{"T", "Me", "", "Waltz", "C", "0", "@" + chart, "", "0", "0"}) + `">`
	if _, err := Parse(strings.TrimSuffix(href, `">`)); err != nil {
		t.Fatal(err)
	}
	if _, err := Parse("irealbook://T=Me=Waltz=C=n=CXyQZ"); err != ErrScheme {
		t.Errorf("irealbook:// accepted: %v", err)
	}
}
