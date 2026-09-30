// Command forms reads the form of every chart of iReal Pro playlists
// from the chords alone, and reports which forms it finds.
//
//	forms [-min 4] [-carrure=false] [-tolerance] [-loose] [-all] [-title regexp] [-form AABA] playlist.html [more.html ...]
//
// It is the bench of analysis.Sections, as corpus is the bench of the
// tonality: the form of each chart as the analysis finds it, grouped by
// its letters, without reading the [A] [B] marks of the chart. See
// analysis.Sections for the method and the settings.
//
// A section is written with its letter and its length in bars, and
// its transposition from the first time it sounds when there is one:
// "A8 A8 B8 A8", "A4 A4+5". The form is the letters alone, "AABA", and
// the report groups the charts by form.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The settings of the form, from the flags.
var opts = analysis.DefaultSections

func main() {
	flag.IntVar(&opts.MinRun, "min", opts.MinRun, "the shortest recurrence kept, in bars")
	flag.BoolVar(&opts.Carrure, "carrure", opts.Carrure, "prefer sections of 4 and 8 bars")
	flag.BoolVar(&opts.Tolerance, "tolerance", opts.Tolerance, "let a recurrence of 6 bars or more hold one bar that differs, and compare chords by family")
	flag.BoolVar(&opts.Loose, "loose", opts.Loose, "compare the roots of the chords only, not their tetrads")
	all := flag.Bool("all", false, "list every chart under its form")
	title := flag.String("title", "", "only the charts whose title matches this regexp, each listed")
	only := flag.String("form", "", `only the charts of this form, "AABA", each listed`)
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: forms [-min 4] [-carrure=false] [-tolerance] [-loose] [-all] [-title regexp] [-form AABA] <playlist> [more ...]")
		os.Exit(2)
	}
	var filter *regexp.Regexp
	if *title != "" {
		var err error
		if filter, err = regexp.Compile("(?i)" + *title); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		*all = true
	}
	if *only != "" {
		*all = true
	}
	var found []form
	for _, path := range flag.Args() {
		songs, err := read(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, s := range songs {
			if filter != nil && !filter.MatchString(s.Title) {
				continue
			}
			changes, err := ireal.Structure(ireal.Lex(s.Chart)).Changes()
			if err != nil || len(changes.Chords) == 0 {
				continue
			}
			f := form{title: s.Title, bars: len(changes.Bars), sections: analysis.SectionsWith(changes, opts)}
			if *only != "" && f.Letters() != *only {
				continue
			}
			found = append(found, f)
		}
	}
	fmt.Print(report(found, *all))
}

var link = regexp.MustCompile(`irealb://[^"\s<]*`)

// read returns the songs of a playlist file, the first of each title,
// as corpus does.
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

// A form is what the analysis finds in one chart.
type form struct {
	title    string
	bars     int
	sections []analysis.Section
}

// Letters returns the form as its letters alone, "AABA".
func (f form) Letters() string {
	var b strings.Builder
	for _, s := range f.sections {
		b.WriteString(s.Label)
	}
	return b.String()
}

// Lengths returns the form with the lengths and transpositions of the
// sections, "A8 A8 B8 A8".
func (f form) Lengths() string {
	var parts []string
	for _, s := range f.sections {
		p := fmt.Sprintf("%s%d", s.Label, s.Bars)
		if s.Shift != 0 {
			p += fmt.Sprintf("+%d", s.Shift)
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, " ")
}

// report groups the charts by form, the most frequent first.
func report(found []form, all bool) string {
	var b strings.Builder
	compared := "root and tetrad"
	switch {
	case opts.Loose:
		compared = "root only"
	case opts.Tolerance:
		compared = "root and family of tetrads, one bar that differs allowed"
	}
	carrure := ""
	if opts.Carrure {
		carrure = ", sections of 4 and 8 bars preferred"
	}
	fmt.Fprintf(&b, "Forms of %d charts, from recurrences of %d bars or more (chords compared on %s%s)\n\n", len(found), opts.MinRun, compared, carrure)

	byForm := map[string][]form{}
	for _, f := range found {
		byForm[f.Letters()] = append(byForm[f.Letters()], f)
	}
	letters := make([]string, 0, len(byForm))
	for l := range byForm {
		letters = append(letters, l)
	}
	slices.SortFunc(letters, func(x, y string) int {
		if d := len(byForm[y]) - len(byForm[x]); d != 0 {
			return d
		}
		return strings.Compare(x, y)
	})

	for _, l := range letters {
		fs := byForm[l]
		fmt.Fprintf(&b, "%-12s %4d  %5.1f%%   %s\n", l, len(fs), 100*float64(len(fs))/float64(len(found)), variants(fs))
		if all {
			for _, f := range fs {
				fmt.Fprintf(&b, "    %-40s %3d bars   %s\n", f.title, f.bars, f.Lengths())
			}
			continue
		}
		for i, f := range fs {
			if i == 3 {
				break
			}
			fmt.Fprintf(&b, "    %s\n", f.title)
		}
	}
	return b.String()
}

// variants lists the three most frequent ways a form is laid out, with
// their counts: "A8 A8 B8 A8 ×301, A8 A8 B8 A10 ×12".
func variants(fs []form) string {
	count := map[string]int{}
	for _, f := range fs {
		count[f.Lengths()]++
	}
	keys := make([]string, 0, len(count))
	for k := range count {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(x, y string) int {
		if d := count[y] - count[x]; d != 0 {
			return d
		}
		return strings.Compare(x, y)
	})
	var parts []string
	for i, k := range keys {
		if i == 3 {
			parts = append(parts, "…")
			break
		}
		parts = append(parts, fmt.Sprintf("%s ×%d", k, count[k]))
	}
	return strings.Join(parts, ", ")
}
