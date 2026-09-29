// Command corpus runs the analysis over every chart of iReal Pro
// playlists and reports where the tonality it hears differs from the
// one the app declares.
//
//	corpus [-aside set-aside.txt] [-keys keys.txt] [-evidence] playlist.html [more.html ...]
//
// The analysis does not read the app's key (see docs/grilles.md): it
// finds the tonality from the cadences. The app is often wrong, and so
// is the analysis at times; the report is there to tell which. Charts
// are grouped by how the two tonalities relate, the relative, the
// fifth, the fourth, since each family of gaps tends to have one cause,
// and each chart comes with the clues that help decide: how many
// cadences resolve, whether the tune ends on the tonic heard, whether it
// is a blues.
//
// The tonality to compare with is the app's, unless the -keys file
// gives one checked by ear, "My Lucky Star | F", in the app's spelling:
// the app is wrong at times, and what we have checked is the data
// worth judging against. The charts checked are counted apart.
//
// With -evidence, the report ends with the charts where the candidate
// with the most evidence (see analysis.Candidates) is not the tonality
// heard: where a proof weighs more than the others, or where the count
// goes wrong.
//
// Some charts lack what tells a tonality, and are set apart rather than
// judged: a modal tune, plages for half of it or more, or one where no
// cadence resolves on a tonic, and the titles listed in the -aside file
// with their reason, the modal tunes whose colours the chart does not
// write, Speak No Evil or Infant Eyes, and the blues of a form [Blues]
// does not know.
package main

import (
	"bufio"
	"flag"
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
	asideFile := flag.String("aside", "", "titles to set aside, one per line, with their reason after \" | \"")
	withEvidence := flag.Bool("evidence", false, "also list the charts where the evidence counted leads to another tonality than the one heard")
	keysFile := flag.String("keys", "", "tonalities checked by ear, one title per line, the key after \" | \" as the app spells it (F, A-)")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: corpus [-aside file] <playlist> [more ...]")
		os.Exit(2)
	}
	aside, err := readList(*asideFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	keys, err := readList(*keysFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var songs []ireal.Song
	for _, path := range flag.Args() {
		s, err := read(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		songs = append(songs, s...)
	}
	fmt.Print(report(songs, aside, keys, *withEvidence))
}

// readList reads a list of titles, each with a value after " | ": the
// reason a chart is set aside, or its tonality checked by ear. Blank
// lines and lines starting with # are skipped.
func readList(path string) (map[string]string, error) {
	out := map[string]string{}
	if path == "" {
		return out, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		title, reason, _ := strings.Cut(line, " | ")
		out[strings.TrimSpace(title)] = strings.TrimSpace(reason)
	}
	return out, sc.Err()
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
	cadences                    int    // blocks that resolve
	endsOnTonic, blues, picardy bool   // the last chord is the tonic heard
	plages                      int    // modal plages
	modal                       bool   // a modal tune, plages for half of it
	aside                       string // why the chart is set aside, listed
	checked                     bool   // the tonality compared with was checked by ear

	// The candidates with the most evidence, and all the candidates
	// with their count, "Cm 9, E♭ 4" (see analysis.Candidates).
	leaders  [][]harmony.Tonality
	evidence string
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
	modal    = "set aside: modal, or no cadence resolving on a tonic"
	listed   = "set aside: listed, with the reason"
	checked  = "same tonality, checked by ear"
)

var order = []string{relative, fifth, fourth, parallel, other, nothing, unread, modal, listed, checked, same}

func relation(r reading) string {
	switch {
	case r.declared == nil:
		return unread
	case r.aside != "":
		return listed
	case r.modal, r.cadences == 0 && !r.blues:
		return modal
	case r.heard == nil:
		return nothing
	}
	d, h := r.declared[0].Tonic(), r.heard[0].Tonic()
	dm, hm := minor(r.declared), minor(r.heard)
	up := (int(h) - int(d) + 12) % 12
	switch {
	case up == 0 && dm == hm && r.checked:
		return checked
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

func analyse(s ireal.Song, aside, keys map[string]string) reading {
	r := reading{title: s.Title, aside: aside[strings.TrimSpace(s.Title)]}
	if key, ok := keys[strings.TrimSpace(s.Title)]; ok {
		s.Key, r.checked = key, true
	}
	if d, ok := s.DeclaredTonalities(); ok {
		r.declared = d
	}
	changes, err := ireal.Structure(ireal.Lex(s.Chart)).Changes()
	if err != nil || len(changes.Chords) == 0 {
		r.declared = nil
		return r
	}
	blocks := analysis.Blocks(changes, analysis.Approaches(changes))
	phrases := analysis.Phrases(changes, blocks)
	for _, s := range analysis.Sense(changes, blocks, phrases, analysis.Tune(changes, phrases)) {
		if s.Resolves != nil && !s.Across {
			r.cadences++
		}
	}
	r.heard = analysis.Tune(changes, phrases)
	var counts []string
	cands := analysis.Candidates(changes, blocks, phrases)
	for _, cd := range cands {
		if len(cd.Evidence) == len(cands[0].Evidence) {
			r.leaders = append(r.leaders, cd.Tonic)
		}
		counts = append(counts, fmt.Sprintf("%s %d", name(cd.Tonic), len(cd.Evidence)))
	}
	r.evidence = strings.Join(counts, ", ")
	r.picardy = analysis.Picardy(changes, phrases)
	_, r.blues = analysis.Blues(changes)
	plages := analysis.Modal(changes, blocks)
	r.plages, r.modal = len(plages), analysis.IsModal(changes, plages)
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

func report(songs []ireal.Song, aside, keys map[string]string, withEvidence bool) string {
	groups := map[string][]reading{}
	for _, s := range songs {
		r := analyse(s, aside, keys)
		rel := relation(r)
		groups[rel] = append(groups[rel], r)
	}
	var b strings.Builder
	judged := len(songs) - len(groups[modal]) - len(groups[listed])
	fmt.Fprintf(&b, "%d charts, %d judged\n\n", len(songs), judged)
	for _, rel := range order {
		fmt.Fprintf(&b, "%5d  %s\n", len(groups[rel]), rel)
	}
	for _, rel := range order {
		rs := groups[rel]
		if rel == same || rel == checked || len(rs) == 0 {
			continue
		}
		slices.SortFunc(rs, func(a, b reading) int { return strings.Compare(a.title, b.title) })
		fmt.Fprintf(&b, "\n%s (%d)\n\n", rel, len(rs))
		fmt.Fprintf(&b, "  %s %s %s %s\n", pad("chart", 44), pad("declared", 9), pad("heard", 9), "clues")
		for _, r := range rs {
			fmt.Fprintf(&b, "  %s %s %s %s\n", pad(cut(r.title, 44), 44), pad(name(r.declared), 9), pad(name(r.heard), 9), clues(r))
		}
	}
	if withEvidence {
		b.WriteString(elsewhere(groups))
	}
	return b.String()
}

// elsewhere lists the charts judged where the evidence gathered does
// not lead to the tonality heard: the candidate with the most evidence,
// or those tied for it, is another. Either a proof is more decisive than
// the others there, or the count is wrong: the charts to look at before
// the count decides anything. How often each agrees with the tonality
// compared with comes first.
func elsewhere(groups map[string][]reading) string {
	var rs []reading
	heardAgrees, leaderAgrees, judged := 0, 0, 0
	for rel, g := range groups {
		if rel == modal || rel == listed {
			continue
		}
		for _, r := range g {
			judged++
			if r.declared != nil && sameTonality(r.heard, r.declared) {
				heardAgrees++
			}
			if r.declared != nil && len(r.leaders) == 1 && sameTonality(r.leaders[0], r.declared) {
				leaderAgrees++
			}
			if !slices.ContainsFunc(r.leaders, func(t []harmony.Tonality) bool { return sameTonality(t, r.heard) }) {
				rs = append(rs, r)
			}
		}
	}
	slices.SortFunc(rs, func(a, b reading) int { return strings.Compare(a.title, b.title) })
	var b strings.Builder
	fmt.Fprintf(&b, "\nthe evidence counted: of %d charts judged, %d agree with the tonality heard, %d with the most evidence, ties excluded\n", judged, heardAgrees, leaderAgrees)
	fmt.Fprintf(&b, "\nthe evidence leads elsewhere (%d)\n\n", len(rs))
	fmt.Fprintf(&b, "  %s %s %s %s\n", pad("chart", 44), pad("declared", 9), pad("heard", 9), "evidence")
	for _, r := range rs {
		fmt.Fprintf(&b, "  %s %s %s %s\n", pad(cut(r.title, 44), 44), pad(name(r.declared), 9), pad(name(r.heard), 9), r.evidence)
	}
	return b.String()
}

// sameTonality reports whether two sets of tonalities name the same
// tonic in the same mode, as name writes them.
func sameTonality(a, b []harmony.Tonality) bool {
	return a != nil && b != nil && name(a) == name(b)
}

// clues lists what helps tell whether the app or the analysis is
// wrong: few cadences resolving on a tonic mean little to hear a
// tonality from, a tune that
// does not end on the tonic heard leaves it to its last cadence, and a
// video game chart is less reliable than a standard.
func clues(r reading) string {
	var out []string
	if r.aside != "" {
		out = append(out, r.aside)
	}
	if r.checked {
		out = append(out, "checked by ear")
	}
	out = append(out, fmt.Sprintf("%d cadences", r.cadences))
	if r.heard != nil && !r.endsOnTonic {
		out = append(out, "ends elsewhere")
	}
	if r.blues {
		out = append(out, "blues")
	}
	if r.plages > 0 {
		out = append(out, fmt.Sprintf("%d modal plage%s", r.plages, plural(r.plages)))
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

func plural(n int) string {
	if n > 1 {
		return "s"
	}
	return ""
}
