// Command forms looks for the form of every chart of iReal Pro
// playlists from the chords alone, and reports which forms it finds.
//
//	forms [-min 4] [-loose] [-tolerance] [-carrure] [-all] [-title regexp] [-form AABA] playlist.html [more.html ...]
//
// It is a probe: nothing in the analysis uses it yet. It tests one
// idea, that the sections of a tune can be found as the passages that
// come back, with no knowledge of any style. It knows no carrure, no
// AABA, no blues, and reads no [A] [B] mark of the chart: those are
// what the report should show, or fail to.
//
// # The method
//
// The chart is unfolded as the app plays one chorus, and each bar
// becomes what sounds in it: its chords, where they start, their root
// and their tetrad. Two passages are the same when bar after bar their
// chords are, at one transposition: a bridge that repeats a motif a
// fourth up is found too.
//
// Every passage of -min bars or more that comes back later is a
// recurrence. The form grows from its first bar: the longest
// recurrence that starts where a section is known to start is placed
// first, and each occurrence placed tells where more sections start.
// Only when none is left does a recurrence start anywhere, and never a
// transposed one (see formOf).
//
// Sections that differ only by their last bars, their turnarounds, are
// the same section: when one occurrence is followed right away by the
// next, as A1 by A2, the chart tells how long the section is, and every
// occurrence covers that length. A stretch that comes back nowhere is
// a section of its own, named with a fresh letter.
//
// # Two settings
//
// -tolerance lets a recurrence of 6 bars or more hold one bar that
// differs, a substitution, and compares chords by family: a 6 and a
// maj7, a m6 and a m7, a 7 and a 7♭5 are the same chord (see family).
//
// -carrure is the one piece of knowledge of a style: sections of 4 and
// 8 bars, the carrure of the standards (Siron, La partition
// intérieure, p. 146 and 147). It only breaks ties: a recurrence
// whose length is within 2 bars of a multiple of 8 is given that
// length, and one shorter than 8 bars may not start off the grid of 8
// bars. The grid starts where most recurrences do, a pickup or an
// introduction put aside (see phase). A section whose length the chart
// tells, 14 bars in Alone Together, is never rounded.
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
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The settings of the probe, from the flags.
var opts struct {
	minRun                    int
	loose, tolerance, carrure bool
}

func main() {
	flag.IntVar(&opts.minRun, "min", 4, "the shortest recurrence kept, in bars")
	flag.BoolVar(&opts.loose, "loose", false, "compare the roots of the chords only, not their tetrads")
	flag.BoolVar(&opts.tolerance, "tolerance", false, "let a recurrence of 6 bars or more hold one bar that differs, and compare chords by family")
	flag.BoolVar(&opts.carrure, "carrure", false, "prefer sections of 4 and 8 bars")
	all := flag.Bool("all", false, "list every chart under its form")
	title := flag.String("title", "", "only the charts whose title matches this regexp, each listed")
	only := flag.String("form", "", `only the charts of this form, "AABA", each listed`)
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: forms [-min 4] [-loose] [-tolerance] [-carrure] [-all] [-title regexp] [-form AABA] <playlist> [more ...]")
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
			f := formOf(barsOf(changes))
			f.title = s.Title
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

// A slot is one chord sounding in a bar: where it starts in the bar,
// or 0 for the chord held from the bar before, and what it is.
type slot struct {
	at      analysis.Ticks
	silent  bool
	root    harmony.PitchClass
	quality harmony.ChordPattern
}

// A bar is what sounds in one bar, in order.
type bar []slot

// barsOf cuts the changes into bars, as the chart plays them.
func barsOf(c analysis.Changes) []bar {
	if len(c.Bars) == 0 {
		return nil
	}
	last := c.Chords[len(c.Chords)-1]
	end := last.Start + last.Length
	out := make([]bar, len(c.Bars))
	for b, from := range c.Bars {
		to := end
		if b+1 < len(c.Bars) {
			to = c.Bars[b+1]
		}
		for _, ch := range c.Chords {
			if ch.Start+ch.Length <= from || ch.Start >= to {
				continue
			}
			s := slot{at: max(ch.Start, from) - from, silent: ch.Silent}
			if !ch.Silent {
				s.root, s.quality = ch.Chord.Root, ch.Chord.Pattern.Tetrad()
			}
			out[b] = append(out[b], s)
		}
	}
	return out
}

// same reports whether bar `b` sounds as bar `a` transposed by `t`
// semitones: the same roots, and the same tetrads, or the same family
// of them with -tolerance, or anything with -loose.
func same(a, b bar, t harmony.Semitones) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.at != y.at || x.silent != y.silent {
			return false
		}
		if x.silent {
			continue
		}
		if x.root.Transpose(t) != y.root {
			return false
		}
		switch {
		case opts.loose:
		case opts.tolerance && family(x.quality) != family(y.quality),
			!opts.tolerance && x.quality != y.quality:
			return false
		}
	}
	return true
}

// A recurrence is a passage of `length` bars starting at bar `a` that
// sounds again from bar `b`, transposed by `shift` semitones.
type recurrence struct {
	a, b, length int
	shift        harmony.Semitones
}

// recurrences finds every passage of -min bars or more that comes
// back, at every distance and transposition. A passage that repeats
// over and over (A A A) is cut into as many recurrences as the
// distance allows, each of them no longer than the distance: two
// occurrences never overlap.
//
// With -tolerance, two runs of bars that match, parted by one bar that
// does not, make one recurrence when it is 6 bars long or more.
func recurrences(bars []bar) []recurrence {
	var out []recurrence
	n := len(bars)
	for lag := opts.minRun; lag < n; lag++ {
		for t := harmony.Semitones(0); t < 12; t++ {
			var runs [][2]int // [start, end) of the bars that match
			run := 0
			for i := 0; i+lag <= n; i++ {
				if i+lag < n && same(bars[i], bars[i+lag], t) {
					run++
					continue
				}
				if run > 0 {
					runs = append(runs, [2]int{i - run, i})
				}
				run = 0
			}
			spans := runs
			if opts.tolerance {
				for k := 0; k+1 < len(runs); k++ {
					if from, to := runs[k][0], runs[k+1][1]; runs[k+1][0] == runs[k][1]+1 && to-from >= 6 {
						spans = append(spans, [2]int{from, to})
					}
				}
			}
			for _, sp := range spans {
				for start := sp[0]; sp[1]-sp[0] >= opts.minRun && start+opts.minRun <= sp[1]; start += lag {
					out = append(out, recurrence{a: start, b: start + lag, length: min(lag, sp[1]-start), shift: t})
				}
			}
		}
	}
	// The longest first; at equal length, untransposed first, then the
	// earliest.
	slices.SortStableFunc(out, func(x, y recurrence) int {
		switch {
		case x.length != y.length:
			return y.length - x.length
		case (x.shift == 0) != (y.shift == 0):
			if x.shift == 0 {
				return -1
			}
			return 1
		}
		return x.a - y.a
	})
	return out
}

// An occurrence is one place a section sounds: from which bar, over
// how many bars the recurrence matched, as which section and at which
// transposition from the section's first occurrence.
type occurrence struct {
	start, length int
	label         int
	shift         harmony.Semitones
}

// A section is one stretch of the form, as the report writes it.
type section struct {
	label string
	bars  int
	shift harmony.Semitones
}

// A form is what the probe finds in one chart.
type form struct {
	title    string
	bars     int
	sections []section
}

// Letters returns the form as its letters alone, "AABA".
func (f form) Letters() string {
	var b strings.Builder
	for _, s := range f.sections {
		b.WriteString(s.label)
	}
	return b.String()
}

// Lengths returns the form with the lengths and transpositions of the
// sections, "A8 A8 B8 A8".
func (f form) Lengths() string {
	var parts []string
	for _, s := range f.sections {
		p := fmt.Sprintf("%s%d", s.label, s.bars)
		if s.shift != 0 {
			p += fmt.Sprintf("+%d", s.shift)
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, " ")
}

// formOf finds the sections of the bars. The recurrences are placed
// one by one, growing the form from its first bar:
//
//   - the best recurrence that starts on a boundary already known, the
//     first bar being one, is placed first;
//   - when none is left, the best untransposed recurrence anywhere, for
//     a tune whose opening comes back nowhere;
//   - a transposed recurrence must start on a boundary: short chains of
//     dominants look alike a fifth apart, and would cut the form at
//     random.
//
// The best is the longest, then the untransposed, then the earliest
// (see [recurrences]). Each occurrence placed adds its start to the
// boundaries, and its end once the length of its section is known (see
// [lengths]).
func formOf(bars []bar) form {
	n := len(bars)
	recs := recurrences(bars)
	grid := phase(recs)
	done := make([]bool, len(recs))
	var placed []occurrence
	labels := 0
	for {
		bounds := boundaries(placed, n)
		i := pick(recs, done, placed, bounds, grid, true)
		if i < 0 {
			i = pick(recs, done, placed, bounds, grid, false)
		}
		if i < 0 {
			break
		}
		done[i] = true
		placed, labels = place(placed, recs[i], labels)
	}
	return sectionsOf(placed, n)
}

// pick returns the first recurrence not `done` that can be placed: one
// that starts on a boundary when `anchored`, an untransposed one
// otherwise. A recurrence that no longer fits is marked done: what is
// placed only grows, so it never will. So is, with -carrure, one
// shorter than 8 bars that starts off the grid of 8 bars that starts
// at bar `grid`.
func pick(recs []recurrence, done []bool, placed []occurrence, bounds map[int]bool, grid int, anchored bool) int {
	for i, r := range recs {
		if done[i] {
			continue
		}
		if fit(placed, r) < opts.minRun || opts.carrure && r.length < 8 && ((r.a-grid)%8 != 0 || (r.b-grid)%8 != 0) {
			done[i] = true
			continue
		}
		if anchored && (bounds[r.a] || bounds[r.b]) || !anchored && r.shift == 0 {
			return i
		}
	}
	return -1
}

// phase returns the bar the grid of 8 bars starts from with -carrure,
// between 0 and 7: the one most recurrences start on, weighed by their
// length. A tune with a pickup bar, or an introduction of 4, has its
// grid shifted as much.
func phase(recs []recurrence) int {
	var weight [8]int
	for _, r := range recs {
		if r.a%8 == r.b%8 {
			weight[r.a%8] += r.length
		}
	}
	best := 0
	for p, w := range weight {
		if w > weight[best] {
			best = p
		}
	}
	return best
}

// startsAt returns the occurrence placed at bar `start`, or -1.
func startsAt(placed []occurrence, start int) int {
	return slices.IndexFunc(placed, func(o occurrence) bool { return o.start == start })
}

// room returns how many bars from `start` are free, up to `length`: 0
// when `start` falls inside an occurrence already placed.
func room(placed []occurrence, start, length int) int {
	for _, o := range placed {
		switch {
		case o.start == start:
		case o.start < start && start < o.start+o.length:
			return 0
		case start < o.start && o.start < start+length:
			length = o.start - start
		}
	}
	return length
}

// fit returns how many bars of `r` can be placed: both occurrences
// must fit, a new one in free bars, or on the start of one placed.
func fit(placed []occurrence, r recurrence) int {
	ia, ib := startsAt(placed, r.a), startsAt(placed, r.b)
	if ia >= 0 && ib >= 0 && placed[ia].label == placed[ib].label {
		return 0 // nothing new
	}
	length := r.length
	if ia < 0 {
		length = min(length, room(placed, r.a, length))
	}
	if ib < 0 {
		length = min(length, room(placed, r.b, length))
	}
	return length
}

// place adds the occurrences of `r` to `placed`, under a new label or
// the label of the occurrence already there, and returns them with the
// number of labels used.
func place(placed []occurrence, r recurrence, labels int) ([]occurrence, int) {
	length := fit(placed, r)
	ia, ib := startsAt(placed, r.a), startsAt(placed, r.b)
	switch {
	case ia < 0 && ib < 0:
		placed = append(placed,
			occurrence{start: r.a, length: length, label: labels},
			occurrence{start: r.b, length: length, label: labels, shift: r.shift})
		labels++
	case ib < 0:
		a := placed[ia]
		placed = append(placed, occurrence{start: r.b, length: length, label: a.label, shift: mod12(a.shift + r.shift)})
	case ia < 0:
		b := placed[ib]
		placed = append(placed, occurrence{start: r.a, length: length, label: b.label, shift: mod12(b.shift - r.shift)})
	default:
		// Two sections found apart are one: relabel the second, its
		// transpositions measured from the first.
		from, delta := placed[ib].label, mod12(placed[ia].shift+r.shift-placed[ib].shift)
		for i := range placed {
			if placed[i].label == from {
				placed[i].label, placed[i].shift = placed[ia].label, mod12(placed[i].shift+delta)
			}
		}
	}
	slices.SortFunc(placed, func(x, y occurrence) int { return x.start - y.start })
	return placed, labels
}

// span returns how many bars there are from occurrence `i` to the next
// one, or to the end of a chart of `n` bars.
func span(placed []occurrence, i, n int) int {
	if i+1 < len(placed) {
		return placed[i+1].start - placed[i].start
	}
	return n - placed[i].start
}

// lengths returns the length of each section whose length the chart
// tells: an occurrence followed by the next one within half its
// matched length again, as A1 by A2, whose last bars differ, spans its
// whole section. The span found most often wins, the shortest at equal
// counts.
func lengths(placed []occurrence, n int) map[int]int {
	counts := map[int]map[int]int{}
	for i, o := range placed {
		s := span(placed, i, n)
		if s*2 > o.length*3 {
			continue
		}
		if counts[o.label] == nil {
			counts[o.label] = map[int]int{}
		}
		counts[o.label][s]++
	}
	out := map[int]int{}
	for label, c := range counts {
		best := 0
		for s, k := range c {
			if best == 0 || k > c[best] || k == c[best] && s < best {
				best = s
			}
		}
		out[label] = best
	}
	return out
}

// extent returns how many bars occurrence `i` covers: the length of
// its section when the chart tells it and it fits before the next
// occurrence, what matched otherwise.
//
// With -carrure, a matched length within 2 bars of a multiple of 8
// becomes that multiple, when it fits. The length the chart tells is
// never rounded.
func extent(placed []occurrence, i, n int, known map[int]int) int {
	o := placed[i]
	s := span(placed, i, n)
	if k, ok := known[o.label]; ok && k >= o.length && k <= s {
		return k
	}
	l := o.length
	if opts.carrure {
		if c := (l + 4) / 8 * 8; c > 0 && c-l <= 2 && l-c <= 2 && c <= s {
			l = c
		}
	}
	return l
}

// boundaries returns the bars where a section is known to start: the
// first, the start of every occurrence, and its end when the length of
// its section is known.
func boundaries(placed []occurrence, n int) map[int]bool {
	out := map[int]bool{0: true}
	known := lengths(placed, n)
	for i, o := range placed {
		out[o.start] = true
		if _, ok := known[o.label]; ok {
			if end := o.start + extent(placed, i, n, known); end < n {
				out[end] = true
			}
		}
	}
	return out
}

// sectionsOf turns the occurrences into the sections of a chart of `n`
// bars. An occurrence covers its section when its length is known (see
// [lengths]), and no more: what follows is another section, the bridge
// after an A' whose last bars prepare it (Billy Boy), unless the tune
// ends there: the last A of All The Things You Are is 12 bars long.
// Otherwise it runs to the next one when what follows is no longer
// than half of it, a turnaround that differs. A longer stretch, and
// whatever comes before the first occurrence, is a section of its own.
func sectionsOf(placed []occurrence, n int) form {
	type raw struct {
		label int // -1 for a stretch that comes back nowhere
		bars  int
		shift harmony.Semitones
	}
	var raws []raw
	stretch := func(from, to int) {
		if to <= from {
			return
		}
		// With -carrure, a long stretch between sections found is cut
		// in sections of 8 bars: a bridge and a last section, the B and
		// the C of AABC, come back nowhere and are still two.
		if l := to - from; opts.carrure && len(placed) > 0 && l > 8 && l%8 == 0 {
			for k := 0; k < l; k += 8 {
				raws = append(raws, raw{label: -1, bars: 8})
			}
			return
		}
		raws = append(raws, raw{label: -1, bars: to - from})
	}
	known := lengths(placed, n)
	pos := 0
	for i, o := range placed {
		stretch(pos, o.start)
		end := o.start + extent(placed, i, n, known)
		_, told := known[o.label]
		if next := o.start + span(placed, i, n); (!told || next == n) && (next-end)*2 <= end-o.start {
			end = next
		}
		raws = append(raws, raw{label: o.label, bars: end - o.start, shift: o.shift})
		pos = end
	}
	stretch(pos, n)

	// Letters in the order the sections first sound, transpositions
	// from that first time.
	f := form{bars: n}
	letter := map[int]string{}
	first := map[int]harmony.Semitones{}
	next := 'A'
	for _, r := range raws {
		s := section{bars: r.bars}
		switch l, ok := letter[r.label]; {
		case r.label < 0:
			s.label = string(next)
			next++
		case ok:
			s.label, s.shift = l, mod12(r.shift-first[r.label])
		default:
			letter[r.label], first[r.label] = string(next), r.shift
			s.label = string(next)
			next++
		}
		f.sections = append(f.sections, s)
	}
	return f
}

// report groups the charts by form, the most frequent first.
func report(found []form, all bool) string {
	var b strings.Builder
	compared := "root and tetrad"
	switch {
	case opts.loose:
		compared = "root only"
	case opts.tolerance:
		compared = "root and family of tetrads, one bar that differs allowed"
	}
	carrure := ""
	if opts.carrure {
		carrure = ", sections of 4 and 8 bars preferred"
	}
	fmt.Fprintf(&b, "Forms of %d charts, from recurrences of %d bars or more (chords compared on %s%s)\n\n", len(found), opts.minRun, compared, carrure)

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

// mod12 brings a transposition between 0 and 11.
func mod12(t harmony.Semitones) harmony.Semitones {
	return (t%12 + 12) % 12
}

// family returns the chord that stands for `p` with -tolerance: the
// major tonic chords are one (triad, 6, maj7), the minor ones another
// (triad, m6, m7, m(maj7)), and the dominants a third, whatever their
// fifth.
func family(p harmony.ChordPattern) harmony.ChordPattern {
	switch p {
	case harmony.ChordMajorTriad, harmony.ChordMajorSixth, harmony.ChordMajorSeventh, harmony.ChordMajorSeventhNo5:
		return harmony.ChordMajorSeventh
	case harmony.ChordMinorTriad, harmony.ChordMinorSixth, harmony.ChordMinorSeventh, harmony.ChordMinorSeventhNo5,
		harmony.ChordMinorMajorSeventh, harmony.ChordMinorMajorSeventhNo5:
		return harmony.ChordMinorSeventh
	case harmony.ChordDominantSeventh, harmony.ChordDominantSeventhNo5, harmony.ChordDominantSeventhFlat5, harmony.ChordDominantSeventhSharp5:
		return harmony.ChordDominantSeventh
	}
	return p
}
