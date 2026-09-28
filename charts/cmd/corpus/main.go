// Command corpus runs the analysis over every chart of iReal Pro
// playlists and reports where the tonality it hears differs from the
// one the app declares.
//
//	corpus playlist.html [more.html ...]
//
// The analysis does not read the app's key (see docs/grilles.md): it
// finds the tonality from the cadences. The app is often wrong, and so
// is the analysis at times; the report is there to tell which. Charts
// are grouped by how the two tonalities relate, the relative, the
// fifth, the fourth, since each family of gaps tends to have one cause,
// and each chart comes with the clues that help decide: how many
// cadences resolve, whether the tune ends on the tonic heard, whether it
// is a blues.
package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: corpus <playlist> [more ...]")
		os.Exit(2)
	}
	var songs []ireal.Song
	for _, path := range os.Args[1:] {
		s, err := read(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		songs = append(songs, s...)
	}
	fmt.Print(report(songs))
}

var link = regexp.MustCompile(`irealb://[^"\s<]*`)

// read returns the songs of a playlist file, the first of each title.
func read(path string) ([]ireal.Song, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []ireal.Song
	seen := map[string]bool{}
	for _, l := range link.FindAllString(string(data), -1) {
		p, err := ireal.Parse(l)
		if err != nil {
			continue
		}
		for _, s := range p.Songs {
			if !seen[s.Title] {
				seen[s.Title] = true
				out = append(out, s)
			}
		}
	}
	return out, nil
}

// A reading is what the report says of one chart.
type reading struct {
	title                       string
	declared, heard             []harmony.Tonality
	cadences                    int  // blocks that resolve
	endsOnTonic, blues, picardy bool // the last chord is the tonic heard
}

// The relations between the tonality declared and the one heard, in
// the order the report lists them.
const (
	same     = "same tonality"
	parallel = "same tonic, other mode (C and Cm)"
	relative = "relative (C and Am)"
	fifth    = "heard a fifth above the declared (C declared, G heard)"
	fourth   = "heard a fourth above the declared (C declared, F heard)"
	other    = "other"
	nothing  = "no tonality heard"
	unread   = "no key declared, or chords unread"
)

var order = []string{relative, fifth, fourth, parallel, other, nothing, unread, same}

func relation(r reading) string {
	switch {
	case r.declared == nil:
		return unread
	case r.heard == nil:
		return nothing
	}
	d, h := r.declared[0].Tonic(), r.heard[0].Tonic()
	dm, hm := minor(r.declared), minor(r.heard)
	up := (int(h) - int(d) + 12) % 12
	switch {
	case up == 0 && dm == hm:
		return same
	case up == 0:
		return parallel
	case !dm && hm && up == 9, dm && !hm && up == 3:
		return relative
	case up == 7:
		return fifth
	case up == 5:
		return fourth
	}
	return other
}

func minor(ts []harmony.Tonality) bool {
	return analysis.ModesOf(ts[:1]) == analysis.Minor
}

func analyse(s ireal.Song) reading {
	r := reading{title: s.Title}
	if d, ok := s.DeclaredTonalities(); ok {
		r.declared = d
	}
	changes, err := ireal.Structure(ireal.Lex(s.Chart)).Changes()
	if err != nil || len(changes.Chords) == 0 {
		r.declared = nil
		return r
	}
	blocks := analysis.Blocks(changes, analysis.Approaches(changes))
	for _, b := range blocks {
		if b.Target >= 0 {
			r.cadences++
		}
	}
	sensed := analysis.Sense(changes, blocks)
	r.heard = analysis.Tune(changes, sensed)
	r.picardy = analysis.Picardy(changes, sensed)
	_, r.blues = analysis.Blues(changes)
	if r.heard != nil {
		for i := len(changes.Chords) - 1; i >= 0; i-- {
			if ch := changes.Chords[i]; !ch.Silent {
				r.endsOnTonic = ch.Chord.Root == r.heard[0].Tonic()
				break
			}
		}
	}
	return r
}

func report(songs []ireal.Song) string {
	groups := map[string][]reading{}
	for _, s := range songs {
		r := analyse(s)
		rel := relation(r)
		groups[rel] = append(groups[rel], r)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d charts\n\n", len(songs))
	for _, rel := range order {
		fmt.Fprintf(&b, "%5d  %s\n", len(groups[rel]), rel)
	}
	for _, rel := range order {
		rs := groups[rel]
		if rel == same || len(rs) == 0 {
			continue
		}
		slices.SortFunc(rs, func(a, b reading) int { return strings.Compare(a.title, b.title) })
		fmt.Fprintf(&b, "\n%s (%d)\n\n", rel, len(rs))
		fmt.Fprintf(&b, "  %s %s %s %s\n", pad("chart", 44), pad("declared", 9), pad("heard", 9), "clues")
		for _, r := range rs {
			fmt.Fprintf(&b, "  %s %s %s %s\n", pad(cut(r.title, 44), 44), pad(name(r.declared), 9), pad(name(r.heard), 9), clues(r))
		}
	}
	return b.String()
}

// clues lists what helps tell whether the app or the analysis is
// wrong: few cadences mean little to hear a tonality from, a tune that
// does not end on the tonic heard leaves it to its last cadence, and a
// video game chart is less reliable than a standard.
func clues(r reading) string {
	var out []string
	out = append(out, fmt.Sprintf("%d cadences", r.cadences))
	if r.heard != nil && !r.endsOnTonic {
		out = append(out, "ends elsewhere")
	}
	if r.blues {
		out = append(out, "blues")
	}
	if r.picardy {
		out = append(out, "Picardy third")
	}
	if strings.Contains(r.title, "(vgls)") {
		out = append(out, "video game")
	}
	return strings.Join(out, ", ")
}

var flats = [12]string{"C", "D♭", "D", "E♭", "E", "F", "G♭", "G", "A♭", "A", "B♭", "B"}

func name(ts []harmony.Tonality) string {
	if ts == nil {
		return "-"
	}
	n := flats[ts[0].Tonic()]
	if minor(ts) {
		n += "m"
	}
	return n
}

// pad fills a string with spaces to a width, counting runes: "♭" is one.
func pad(s string, w int) string {
	return s + strings.Repeat(" ", max(0, w-utf8.RuneCountInString(s)))
}

// cut shortens a title to a width, counting runes.
func cut(s string, w int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) > w {
		return string(r[:w-1]) + "…"
	}
	return string(r)
}
