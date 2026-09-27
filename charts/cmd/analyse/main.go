// Command analyse prints a chart of an iReal Pro playlist with what the
// analysis sees in it, bar by bar, in the terminal.
//
//	analyse playlist.html "tenderly"
//
// The playlist is an export of the app: an HTML file holding irealb://
// links, or a text file holding one. The title is matched without case,
// on any part of it; when several songs match, they are listed.
//
// It is the test bench of the analysis (see docs/grilles.md), and grows
// with it: for now each chord says how it prepares the next one.
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

const barsPerRow = 4

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: analyse <playlist> <title>")
		os.Exit(2)
	}
	song, err := find(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(render(song))
}

var link = regexp.MustCompile(`irealb://[^"\s<]*`)

// find reads the playlists of a file and picks the song whose title
// holds the words asked for, or the one whose title is exactly them.
func find(path, title string) (ireal.Song, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ireal.Song{}, err
	}
	want := strings.ToLower(title)
	var found []ireal.Song
	for _, l := range link.FindAllString(string(data), -1) {
		p, err := ireal.Parse(l)
		if err != nil {
			continue
		}
		for _, s := range p.Songs {
			t := strings.ToLower(s.Title)
			if t == want {
				return s, nil
			}
			if strings.Contains(t, want) {
				found = append(found, s)
			}
		}
	}
	switch len(found) {
	case 0:
		return ireal.Song{}, fmt.Errorf("no song in %s has %q in its title", path, title)
	case 1:
		return found[0], nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d songs match %q:\n", len(found), title)
	for _, s := range found {
		fmt.Fprintf(&b, "  %s\n", s.Title)
	}
	return ireal.Song{}, fmt.Errorf("%s", strings.TrimRight(b.String(), "\n"))
}

// render lays the chart out four bars to a row, as a lead sheet does,
// with the analysis on lines of its own above the chords: over each
// chord, how it prepares the next, and above, the blocks with the tonality
// they announce and, when it differs, the block read as En Harmonie
// brackets it. Under the chords, their degrees on the installed tonic,
// and under the degrees the sensed tonic where it changes, from the
// tonality the app gives the tune.
func render(s ireal.Song) string {
	chart := ireal.Structure(ireal.Lex(s.Chart))
	tl := chart.Timeline()
	changes, err := tl.Changes()
	kinds := analysis.Approaches(changes)
	passing := analysis.PassingChords(changes)
	blocks := analysis.Blocks(changes, kinds)
	home, _ := s.HomeTonalities()
	sensed := analysis.Sense(changes, blocks, home)
	degrees := analysis.Degrees(changes, passing, sensed)
	bracket := analysis.Bracketed(changes, blocks, passing, sensed)

	// Each bar is two lines that line up word for word: the analysis
	// above, the chords below. A bar holds "%" when the chord before
	// goes on.
	cells := make([]cell, len(tl.Bars))
	started := make([]bool, len(tl.Bars))
	words := make([]word, len(tl.Spans)) // where each chord sits in its bar
	for i, sp := range tl.Spans {
		name := "N.C."
		if !sp.NoChord {
			name = symbol(sp.Chord)
		}
		off, w := cells[sp.Bar].add(walk(passing[i])+label(kinds[i]), name, degrees[i].String(), heard(sensed, i))
		words[i] = word{sp.Bar, off, w}
		started[sp.Bar] = true
	}
	width := 0
	for i := range cells {
		if !started[i] {
			cells[i].add("", "%", "", "")
		}
		width = max(width, cells[i].width())
	}

	// The structure: a section or the coda starts a row of its own, as
	// on a lead sheet, and its mark sits over the bar number.
	marks := make([]string, len(tl.Bars))
	for bar, p := range chart.Unfold() {
		if sec := chart.Measures[p.Index].Section; sec != 0 {
			marks[bar] = "[" + string(sec) + "]"
		}
	}
	if tl.Coda > 0 {
		marks[tl.Spans[tl.Coda].Bar] = "coda"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s, %s (%s), %d bars played\n", s.Title, s.Key, s.Style, len(tl.Bars))
	fmt.Fprint(&b, legend)
	for start := 0; start < len(cells); {
		end := start + 1
		for end < len(cells) && end-start < barsPerRow && marks[end] == "" {
			end++
		}
		fmt.Fprintf(&b, "\n%s", drawBlocks(blocks, tl, words, bracketed(blocks, degrees, bracket), start, end, width))
		fmt.Fprintf(&b, "\n%4s  ", marks[start])
		for _, c := range cells[start:end] {
			fmt.Fprintf(&b, " %s  ", pad(c.top, width))
		}
		fmt.Fprintf(&b, "\n%4d |", start+1)
		for _, c := range cells[start:end] {
			fmt.Fprintf(&b, " %s |", pad(c.bottom, width))
		}
		fmt.Fprintf(&b, "\n%4s  ", "")
		for _, c := range cells[start:end] {
			fmt.Fprintf(&b, " %s  ", pad(c.below, width))
		}
		var tonics strings.Builder
		for _, c := range cells[start:end] {
			fmt.Fprintf(&tonics, " %s  ", pad(c.tonic, width))
		}
		if strings.TrimSpace(tonics.String()) != "" {
			fmt.Fprintf(&b, "\n%4s  %s", "", strings.TrimRight(tonics.String(), " "))
		}
		fmt.Fprintln(&b)
		start = end
	}
	if err != nil {
		fmt.Fprintf(&b, "\nunread: %v\n", err)
	}
	return b.String()
}

// A cell is one bar on its four lines: the analysis over the chords,
// the degrees under them, and the sensed tonic under the degrees.
type cell struct{ top, bottom, below, tonic string }

// add puts a word on each line, all padded to the widest so that what
// follows stays aligned, and returns where it starts in the bar and how
// wide it is.
func (c *cell) add(top, bottom, below, tonic string) (int, int) {
	w := max(utf8.RuneCountInString(top), utf8.RuneCountInString(bottom),
		utf8.RuneCountInString(below), utf8.RuneCountInString(tonic))
	if c.bottom != "" {
		c.top, c.bottom, c.below, c.tonic = c.top+" ", c.bottom+" ", c.below+" ", c.tonic+" "
	}
	off := utf8.RuneCountInString(c.bottom)
	c.top += pad(top, w)
	c.bottom += pad(bottom, w)
	c.below += pad(below, w)
	c.tonic += pad(tonic, w)
	return off, w
}

// heard names the sensed tonic at a change when it differs from the
// change before: its ground, then in brackets the local tonic a
// cadence has just tonicised. Nothing when nothing changed.
func heard(sensed []analysis.Sensed, i int) string {
	name := func(s analysis.Sensed) string {
		if s.Ground == nil {
			return "?"
		}
		n := short(flats[s.Ground[0].Tonic()], s.Ground)
		if s.Local != nil {
			n += " (" + short(flats[s.Local[0].Tonic()], s.Local) + ")"
		}
		return n
	}
	if i > 0 && name(sensed[i]) == name(sensed[i-1]) {
		return ""
	}
	return name(sensed[i])
}

// A word is where a chord sits: its bar, its offset in the bar, its
// width.
type word struct{ bar, off, width int }

// Columns of a row: the bar number, then each bar with a space before
// and " |" after.
const margin = 6

func column(bar, first, width int) int {
	return margin + (bar-first)*(width+3) + 1
}

// drawBlocks draws the line of blocks over a row of bars: from the
// first chord of each block to its V, the tonality it announces then a
// rule. A block cut by the end of a row goes on in the next one, and a
// block across the loop, its two at the end of the chorus and its V at
// the start, is drawn in both places. A name that would run into the
// next block's is cut short.
func drawBlocks(blocks []analysis.Block, tl ireal.Timeline, words []word, readings []string, first, end, width int) string {
	line := []rune(strings.Repeat(" ", column(end, first, width)))
	type name struct {
		at          int
		text, short []rune
	}
	var names []name
	for n, bl := range blocks {
		start := bl.Five
		if bl.Sus >= 0 {
			start = bl.Sus
		}
		if bl.Two >= 0 {
			start = bl.Two
		}
		a, z := words[start], words[bl.Five]
		wraps := a.bar > z.bar
		covers := func(bar int) bool {
			if wraps {
				return bar >= a.bar || bar <= z.bar
			}
			return bar >= a.bar && bar <= z.bar
		}
		from, to := -1, -1
		for bar := first; bar < end; bar++ {
			if !covers(bar) {
				continue
			}
			if from < 0 {
				from = column(bar, first, width)
			}
			to = column(bar, first, width) + width
		}
		if from < 0 {
			continue
		}
		if a.bar >= first && a.bar < end {
			from = column(a.bar, first, width) + a.off
			long, short := announced(bl, tl)
			if readings[n] != "" {
				long += " : " + readings[n]
			}
			names = append(names, name{from, []rune(long + " "), []rune(short + " ")})
		}
		if z.bar >= first && z.bar < end {
			to = column(z.bar, first, width) + z.off + z.width
		}
		for i := from; i < to && i < len(line); i++ {
			line[i] = '─'
		}
	}
	for i, n := range names {
		stop := len(line)
		if i+1 < len(names) {
			stop = names[i+1].at
		}
		text := n.text
		if n.at+len(text) > stop {
			text = n.short
		}
		for j, r := range text {
			if n.at+j >= stop {
				break
			}
			line[n.at+j] = r
		}
	}
	return strings.TrimRight(string(line), " ")
}

// bracketed gives for each block its chords as En Harmonie brackets
// them, in the tonality it announces ("II V"), when its degrees differ
// from those on the installed tonic; "" when they do not, as when a two
// five announces the minor of the tonic (Fm7♭5 B♭7 in E flat).
func bracketed(blocks []analysis.Block, ground, bracket []analysis.Degree) []string {
	out := make([]string, len(blocks))
	for n, bl := range blocks {
		var ds []string
		differs := false
		for _, i := range []int{bl.Two, bl.Sus, bl.Five} {
			if i >= 0 {
				ds = append(ds, bracket[i].String())
				differs = differs || bracket[i].Number != ground[i].Number ||
					bracket[i].Accidental != ground[i].Accidental
			}
		}
		if differs {
			out[n] = strings.Join(ds, " ")
		}
	}
	return out
}

// announced names the tonalities a block announces: their tonic as the
// chart spells the target, or in flats, then their scales: "E♭" for E
// flat major, "Fm harm" for F harmonic minor, "D♭ M/m mel" when major
// and melodic minor both remain, "Gm nat/harm/mel" when the three
// minors do; "…" when the block does not resolve. The short name,
// for when the next block leaves no room, keeps the mode alone: "Gm",
// "D♭(m)" when both remain.
func announced(bl analysis.Block, tl ireal.Timeline) (string, string) {
	if len(bl.Announced) == 0 {
		return "?", "?"
	}
	tonic := flats[bl.Announced[0].Tonic()]
	if bl.Target >= 0 && !tl.Spans[bl.Target].NoChord {
		tonic = note(tl.Spans[bl.Target].Chord.Root)
	}
	name, short := scales(tonic, bl.Announced)
	if bl.Target < 0 {
		name, short = name+"…", short+"…"
	}
	return name, short
}

// short is the short name of tonalities on a tonic: "Gm", "D♭(m)".
func short(tonic string, ts []harmony.Tonality) string {
	_, s := scales(tonic, ts)
	return s
}

// scales names tonalities on a tonic, long and short: see [announced].
func scales(tonic string, ts []harmony.Tonality) (string, string) {
	major := false
	var minors []string
	for _, t := range ts {
		switch sc, _ := harmony.NamedScaleOf(t.Pattern()); sc {
		case harmony.NamedMajor, harmony.NamedHarmonicMajor:
			major = true
		case harmony.NamedNaturalMinor:
			minors = append(minors, "nat")
		case harmony.NamedHarmonicMinor:
			minors = append(minors, "harm")
		case harmony.NamedMelodicMinor:
			minors = append(minors, "mel")
		}
	}
	var name, short string
	switch {
	case len(minors) == 0:
		name, short = tonic, tonic
	case !major && len(minors) == 1 && minors[0] == "nat":
		name, short = tonic+"m", tonic+"m"
	case !major:
		name, short = tonic+"m "+strings.Join(minors, "/"), tonic+"m"
	default:
		name, short = tonic+" M/m "+strings.Join(minors, "/"), tonic+"(m)"
	}
	return name, short
}

var flats = [12]string{"C", "D♭", "D", "E♭", "E", "F", "G♭", "G", "A♭", "A", "B♭", "B"}

func (c cell) width() int {
	return max(utf8.RuneCountInString(c.top), utf8.RuneCountInString(c.bottom),
		utf8.RuneCountInString(c.below), utf8.RuneCountInString(c.tonic))
}

func pad(s string, w int) string {
	return s + strings.Repeat(" ", max(0, w-utf8.RuneCountInString(s)))
}

// legend explains the marks over the chords, one to a line.
const legend = `
over each chord:
  V→    prepares the next as its dominant
  ♭II→  as its chromatic dominant
  °→    as its diminished chord, a dominant without its root
  sus→  as its suspension
  II→   as the two of a dominant
  ↗ ↘   a passing chord, its bass walking up or down

above, the blocks, [II] [sus4] V, and the tonality each one announces:
  Fm harm ─     a two five announcing F harmonic minor, whatever it lands on
  Fm harm : II V ─
                and its chords read in F minor, as En Harmonie brackets
                them, when that differs from the degrees below
  D♭ M/m mel… ─ one that does not resolve, nor decide between D♭ major
                and D♭ melodic minor

under each chord, its degree on the tonic the ear has installed there,
IIIm7♭5 VI7 for Gm7♭5 C7 in E♭; ♭VII7, ♯Vdim7, I/3

under the degrees, the tonic the ear senses, where it changes:
  E♭            installed, from the app's key or the first cadence
  E♭ (Fm)       and F minor, which a cadence has just tonicised
`

// walk marks a passing chord by the way its bass goes.
func walk(dir int) string {
	switch dir {
	case 1:
		return "↗"
	case -1:
		return "↘"
	}
	return ""
}

// label says how a chord prepares the next.
func label(k harmony.ApproachKind) string {
	var out []string
	for _, l := range []struct {
		kind harmony.ApproachKind
		text string
	}{
		{harmony.DominantApproach, "V"},
		{harmony.ChromaticApproach, "♭II"},
		{harmony.DiminishedApproach, "°"},
		{harmony.SuspensionApproach, "sus"},
		{harmony.TwoApproach, "II"},
	} {
		if k.Has(l.kind) {
			out = append(out, l.text)
		}
	}
	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, "+") + "→"
}

// spelling writes the app's qualities the way a lead sheet does, until
// naming renders chord symbols (see docs/chantiers.md).
var spelling = strings.NewReplacer(
	"^", "maj", "-", "m", "h7", "m7♭5", "h9", "m9♭5", "h", "m7♭5",
	"o", "dim", "b", "♭", "#", "♯",
)

func symbol(c ireal.ChordSymbol) string {
	s := note(c.Root) + spelling.Replace(c.Quality)
	if c.Bass != "" {
		s += "/" + note(c.Bass)
	}
	return s
}

func note(n string) string {
	return strings.NewReplacer("b", "♭", "#", "♯").Replace(n)
}
