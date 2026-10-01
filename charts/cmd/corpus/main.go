// Command corpus runs the analysis over every chart of iReal Pro
// playlists and reports where the tonality it hears differs from the
// one the app declares.
//
//	corpus [-aside set-aside.txt] [-keys keys.txt] [-forms] [-phrases] playlist.html [more.html ...]
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
// The report ends with the form: the sections the analysis finds from
// the chords (see analysis.Sections), compared with the rehearsal marks
// of the chart, [A] [B], which the analysis does not read. As the key,
// the marks are often missing or loose. With -forms, the charts that
// disagree are listed, the form marked above the form found.
//
// Then the phrases: where the analysis comes to rest (see
// analysis.Phrases), compared with where the sections found conclude
// (see analysis.Conclusions), the rule of Siron that is to replace ours.
// With -phrases, the charts where they part are listed.
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
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

func main() {
	asideFile := flag.String("aside", "", "titles to set aside, one per line, with their reason after \" | \"")
	keysFile := flag.String("keys", "", "tonalities checked by ear, one title per line, the key after \" | \" as the app spells it (F, A-)")
	withForms := flag.Bool("forms", false, "also list the charts whose form found disagrees with the marks of the chart")
	withPhrases := flag.Bool("phrases", false, "also list the charts where the phrases do not rest where the sections conclude")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: corpus [-aside file] [-keys file] [-forms] [-phrases] <playlist> [more ...]")
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
	fmt.Print(report(songs, aside, keys, *withForms, *withPhrases))
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

	// The sections the analysis finds from the chords, and those the
	// chart marks (see forms).
	found, marked []analysis.Section

	// The bars where the phrases come to rest, and those where the
	// sections found conclude, -1 for a section that concludes on none
	// (see phrases).
	rests, conclusions []int
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
	chart := ireal.Structure(ireal.Lex(s.Chart))
	changes, err := chart.Changes()
	r.marked = chart.MarkedSections()
	if err == nil {
		r.found = analysis.Sections(changes)
	}
	if err != nil || len(changes.Chords) == 0 {
		r.declared = nil
		return r
	}
	h := analysis.Hear(changes, nil)
	blocks, phrases := h.Blocks, h.Phrases
	for _, s := range h.Sensed {
		if s.Resolves != nil && !s.Across {
			r.cadences++
		}
	}
	r.heard = h.Heard
	for _, p := range phrases {
		if p.Tonic != nil && p.Arrives >= 0 {
			r.rests = append(r.rests, changes.Bar(changes.Chords[p.Arrives].Start))
		}
	}
	for _, k := range analysis.Conclusions(changes, r.found, blocks) {
		bar := -1
		if k.Arrives >= 0 {
			bar = changes.Bar(changes.Chords[k.Arrives].Start)
		}
		r.conclusions = append(r.conclusions, bar)
	}
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

func report(songs []ireal.Song, aside, keys map[string]string, withForms, withPhrases bool) string {
	groups := map[string][]reading{}
	var all []reading
	for _, s := range songs {
		r := analyse(s, aside, keys)
		rel := relation(r)
		groups[rel] = append(groups[rel], r)
		all = append(all, r)
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
	b.WriteString(forms(all, withForms))
	b.WriteString(phrasing(all, withPhrases))
	return b.String()
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

// name writes tonalities by their tonic, spelled with the fewest
// accidentals in their key signature (see naming.TonicSpelling), and
// "m" for a minor one: "F♯m", "D♭".
func name(ts []harmony.Tonality) string {
	if ts == nil {
		return "-"
	}
	n := naming.English.Name(naming.TonicSpelling(ts[0]), naming.Signs)
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
