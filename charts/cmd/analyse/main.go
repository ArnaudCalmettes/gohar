// Command analyse prints a chart of an iReal Pro playlist with what the
// analysis sees in it, bar by bar, in the terminal.
//
//	analyse [-key heard|declared|F|A-] [-legend] [-smells] [-maj7 maj|natural|delta] [-minus] [-halfdim] [-dimsign] [-slash69] playlist.html "tenderly"
//
// The playlist is an export of the app: an HTML file holding irealb://
// links, or a text file holding one. The title is matched without case,
// on any part of it; when several songs match, they are listed.
//
// # The tonality of the tune
//
// The degrees are counted in the tonality of the tune, and the chords
// alone cannot always tell it (see docs/grilles.md, "Le plafond des
// grilles seules"), but the analysis judges as an analyst does: by
// default the chart is analysed in the tonality it concludes (see
// analysis.Tune), and the heading says when the app declares another.
//
// -key declared analyses it in the key the app declares instead, the
// tonality heard when it declares none, and -key F or -key A- in the
// one the reader forces, in the app's spelling. A tune that opens on
// the relative of its key, Lullaby Of Birdland on F minor in A flat,
// is not counted from there: its first bars tonicise F minor, which
// the tune may install later.
//
// # The spelling of the chords
//
// The chords are not written as the app writes them: their roots and
// basses are respelled so that their movements read plainly and their
// names follow the degrees, and the heading counts the smells, the
// roots and basses written with a double accidental or as E♯, F♭, B♯,
// C♭, and the enharmonies tolerated, when run with -smells (see
// respell).
//
// Their qualities are written by naming (see naming.ChordStyle): Cmaj7,
// Cm7, Cm7♭5, Cdim7 by default; -maj7 natural or -maj7 delta writes C♮7
// or CΔ7, -minus C-7, -halfdim Cø, -dimsign C°7, -slash69 C6/9 rather
// than C6⁄9.
//
// It is the test bench of the analysis (see docs/grilles.md), and grows
// with it: for now each chord says how it prepares the next one.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

const barsPerRow = 4

func main() {
	withLegend := flag.Bool("legend", false, "explain the marks before the chart")
	withSmells := flag.Bool("smells", false, "list the spellings that smell and the enharmonies tolerated")
	key := flag.String("key", "heard", `the tonality to analyse the tune in: "heard" by the analysis, "declared" by the app, or a key as the app spells it (F, A-)`)
	maj7 := flag.String("maj7", "maj", `how to write the major seventh: "maj" (Cmaj7), "natural" (C♮7) or "delta" (CΔ7)`)
	flag.BoolVar(&style.Minus, "minus", false, "write a minor chord C-7 rather than Cm7")
	flag.BoolVar(&style.HalfDiminishedSign, "halfdim", false, "write a half-diminished chord Cø rather than Cm7♭5")
	flag.BoolVar(&style.DiminishedSign, "dimsign", false, "write a diminished chord C°7 rather than Cdim7")
	flag.BoolVar(&style.Slash, "slash69", false, "write a sixth chord with its ninth C6/9 rather than C6⁄9")
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: analyse [-key heard|declared|F|A-] [-legend] [-smells] [-maj7 maj|natural|delta] [-minus] [-halfdim] [-dimsign] [-slash69] <playlist> <title>")
		os.Exit(2)
	}
	sevenths := map[string]naming.MajorSeventhSign{"natural": naming.NaturalSeventh, "delta": naming.DeltaSeventh, "maj": naming.MajSeventh}
	sign, ok := sevenths[*maj7]
	if !ok {
		fmt.Fprintf(os.Stderr, "-maj7 %s: not maj, natural or delta\n", *maj7)
		os.Exit(2)
	}
	style.MajorSeventh = sign
	if _, ok := (ireal.Song{Key: *key}).DeclaredTonalities(); !ok && *key != "declared" && *key != "heard" {
		fmt.Fprintf(os.Stderr, "-key %s: not a key, nor declared or heard\n", *key)
		os.Exit(2)
	}
	song, err := find(flag.Arg(0), flag.Arg(1))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(render(song, *key, *withLegend, *withSmells))
}

// tuneOf returns the tonality to analyse a song in, as -key asks, and
// why: the one the app declares, the one the analysis hears, or the one
// the reader forces.
func tuneOf(s ireal.Song, key string, heard []harmony.Tonality) ([]harmony.Tonality, string) {
	switch key {
	case "heard":
		return heard, "as heard"
	case "declared":
		t, ok := s.DeclaredTonalities()
		switch {
		case !ok:
			return heard, "as heard, the app declaring none"
		}
		return t, "as the app declares"
	}
	t, _ := (ireal.Song{Key: key}).DeclaredTonalities()
	return t, "as asked"
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
func render(s ireal.Song, key string, withLegend, withSmells bool) string {
	chart := ireal.Structure(ireal.Lex(s.Chart))
	tl := chart.Timeline()
	changes, err := tl.Changes()
	passing := analysis.PassingChords(changes)
	var why string
	h := analysis.Hear(changes, func(heard []harmony.Tonality) []harmony.Tonality {
		tune, reason := tuneOf(s, key, heard)
		why = reason
		return tune
	})
	kinds, blocks, phrases, tune, sensed := h.Kinds, h.Blocks, h.Phrases, h.Tune, h.Sensed
	plages := analysis.Modal(changes, blocks)
	degrees := analysis.Degrees(changes, passing, sensed)
	bracket := analysis.Bracketed(changes, blocks, passing, sensed)
	pedals := analysis.Pedals(changes, analysis.Grounds(changes, sensed), analysis.Sections(changes))
	spelled := respell(tl, changes, sensed, blocks, degrees, pedals)
	nm := spelled.names
	held := pedalRules(pedals, tl)
	formulas := cellRules(analysis.InTonality(changes, analysis.Cells(changes), sensed))
	steps := make([]string, len(changes.Chords)) // how a II-V, or a V, follows the one before it
	for _, l := range analysis.Links(changes, kinds) {
		steps[l.To] = l.Step.String() + " "
	}

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
		degree, tonic := degrees[i].String(), heard(nm, sensed, i)
		if inPlage(plages, i) {
			degree, tonic = "modal", ""
		}
		off, w := cells[sp.Bar].add(steps[i]+walk(passing[i])+label(kinds[i]), name, degree, tonic)
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
	read := analysis.ReadTune(changes, phrases)
	fmt.Fprint(&b, heardIn(nm, s, read, plages))
	if tune != nil {
		fmt.Fprintf(&b, "analysed in %s, %s\n", nm.of(tune, -1), why)
	}
	if withSmells {
		fmt.Fprint(&b, smells(spelled, tl))
	}
	fmt.Fprint(&b, reading(nm, read, tl))
	fmt.Fprint(&b, form(nm, changes, tl, blocks))
	fmt.Fprint(&b, areas(nm, changes, tl, blocks, sensed))
	if withLegend {
		fmt.Fprint(&b, legend)
	} else {
		fmt.Fprint(&b, "\nrun with -legend to explain the notation\n")
	}
	for start := 0; start < len(cells); {
		end := start + 1
		for end < len(cells) && end-start < barsPerRow && marks[end] == "" {
			end++
		}
		fmt.Fprintf(&b, "\n%s", drawBlocks(nm, blocks, tl, words, bracketed(blocks, degrees, bracket), start, end, width))
		fmt.Fprintf(&b, "\n%4s  ", marks[start])
		for _, c := range cells[start:end] {
			fmt.Fprintf(&b, " %s  ", pad(c.top, width))
		}
		fmt.Fprintf(&b, "\n%4d |", start+1)
		for _, c := range cells[start:end] {
			fmt.Fprintf(&b, " %s |", pad(c.bottom, width))
		}
		if line := drawRules(held, words, start, end, width); line != "" {
			fmt.Fprintf(&b, "\n%s", line)
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
		if line := drawRules(formulas, words, start, end, width); line != "" {
			fmt.Fprintf(&b, "\n%s", line)
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

// heardIn says in which tonality the analysis hears the tune, and when
// it differs, what the app declares: the analysis does not read the
// app's key, which is often wrong.
func heardIn(nm namer, s ireal.Song, r analysis.TuneReading, plages []analysis.Plage) string {
	var modal string
	if n := len(plages); n > 0 {
		modal = fmt.Sprintf(", %d modal plage%s, the chart not saying which mode", n, plural(n))
	}
	if r.Tonality == nil {
		return "no tonality heard" + modal + "\n"
	}
	name := nm.of(r.Tonality, -1)
	line := "heard in " + name + modal
	if r.Blues {
		line += ", a blues"
	}
	if r.First != nil && nm.of(r.First, -1) != name {
		line += ", setting out from " + nm.of(r.First, -1)
	}
	// The declared key is spelled as a key signature would be, B♭ and
	// not the A♯ a zone of the tune may have chosen.
	if declared, ok := s.DeclaredTonalities(); ok {
		if d := (namer{}).of(declared, -1); d != (namer{}).of(r.Tonality, -1) {
			line += ", where the app declares " + d
		}
	}
	return line + "\n"
}

// form writes the sections the analysis finds from the chords, the
// marks of the chart left aside, and how each ends: the tonic its
// conclusive cadence lands on, at which bar, strong or weak, and the
// chords of the turnaround that goes on from it (see
// analysis.Conclusion).
//
//	the form it finds: A8 A8 B8 A8
//	  A   bars 1-8     C at bar 7, strong, then A7 Dm7 G7
func form(nm namer, changes analysis.Changes, tl ireal.Timeline, blocks []analysis.Block) string {
	sections := analysis.Sections(changes)
	if len(sections) == 0 {
		return ""
	}
	var b, lengths strings.Builder
	for i, s := range sections {
		if i > 0 {
			lengths.WriteString(" ")
		}
		fmt.Fprintf(&lengths, "%s%d", s.Label, s.Bars)
		if s.Shift != 0 {
			fmt.Fprintf(&lengths, "+%d", s.Shift)
		}
	}
	fmt.Fprintf(&b, "the form it finds: %s\n", lengths.String())
	for i, k := range analysis.Conclusions(changes, sections, blocks) {
		s := sections[i]
		line := "no conclusive cadence"
		if k.Half >= 0 {
			line = fmt.Sprintf("half cadence on %s, bar %d, the V of %s", symbol(tl.Spans[k.Half].Chord), tl.Spans[k.Half].Bar+1, nm.of(k.Tonic, -1))
		}
		if k.Arrives >= 0 {
			line = fmt.Sprintf("%s at bar %d", nm.of(k.Tonic, -1), tl.Spans[k.Arrives].Bar+1)
			if k.Strong {
				line += ", strong"
			} else {
				line += ", weak"
			}
			if k.Loop >= 0 {
				// The turnaround goes from the tonic the section lands on
				// back to the start: its chords, to the end of the section.
				var loop []string
				for i := k.Loop; i < len(tl.Spans) && tl.Spans[i].Bar < s.From+s.Bars; i++ {
					if !tl.Spans[i].NoChord {
						loop = append(loop, symbol(tl.Spans[i].Chord))
					}
				}
				line += ", then " + strings.Join(loop, " ")
			}
		}
		fmt.Fprintf(&b, "  %-3s %-12s %s\n", s.Label, fmt.Sprintf("bars %d-%d", s.From+1, s.From+s.Bars), line)
	}
	return b.String()
}

// reading tells what the tonality heard rests on (see
// analysis.ReadTune): the first chord, the last one, and when they part,
// how long each of their tonics is heard.
//
//	how the tonality is heard:
//	  first chord    Cm6, bar 1: Cm, its first cadence resolves on it
//	  last chord     E♭6, bar 35: E♭, the turnaround left out
//	  heard longest  Cm 20 bars, E♭ 12 bars
func reading(nm namer, r analysis.TuneReading, tl ireal.Timeline) string {
	if r.Tonality == nil {
		return ""
	}
	name := func(t []harmony.Tonality) string { return nm.of(t, -1) }
	chord := func(i int) string {
		return fmt.Sprintf("%s, bar %d", symbol(tl.Spans[i].Chord), tl.Spans[i].Bar+1)
	}
	var b strings.Builder
	b.WriteString("how the tonality is heard:\n")
	if r.Blues {
		fmt.Fprintf(&b, "  a blues, in %s by its form\n", name(r.Tonality))
		return b.String()
	}
	switch {
	case r.Opens < 0:
	case r.First == nil && r.Cadence == nil:
		fmt.Fprintf(&b, "  first chord    %s: not the tonic, no cadence resolves\n", chord(r.Opens))
	case r.First == nil:
		fmt.Fprintf(&b, "  first chord    %s: not the tonic, the first cadence goes to %s\n", chord(r.Opens), name(r.Cadence))
	default:
		fmt.Fprintf(&b, "  first chord    %s: %s, its first cadence resolves on it\n", chord(r.Opens), name(r.First))
	}
	if r.Stops < 0 {
		fmt.Fprintf(&b, "  last chord     none concludes: %s, where the first phrase rests\n", name(r.Tonality))
		return b.String()
	}
	if r.Picardy {
		fmt.Fprintf(&b, "  last chord     %s: a picardy third, the tune in %s, the turnaround left out\n", chord(r.Stops), name(r.Last))
	} else {
		fmt.Fprintf(&b, "  last chord     %s: %s, the turnaround left out\n", chord(r.Stops), name(r.Last))
	}
	if r.Heard != nil {
		f, l := r.First[0].Tonic(), r.Last[0].Tonic()
		fmt.Fprintf(&b, "  heard longest  %s %.4g bars, %s %.4g bars\n", name(r.First), r.Heard[f], name(r.Last), r.Heard[l])
	}
	return b.String()
}

// areas lists the stretches heard around another tonic than the first,
// with what Siron weighs to tell a true modulation from a transitory
// one (see analysis.TonalArea):
//
//	the tonal areas it hears, and what would make each a true modulation:
//	  C      bars 5-8      4 bars, leaving D (2 steps), the first tonality: transitory
func areas(nm namer, changes analysis.Changes, tl ireal.Timeline, blocks []analysis.Block, sensed []analysis.Sensed) string {
	found := analysis.TonalAreas(changes, blocks, sensed, analysis.Sections(changes))
	if len(found) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("the tonal areas it hears, and what would make each a true modulation:\n")
	for _, a := range found {
		last := changes.Chords[a.To]
		bars := fmt.Sprintf("bars %d-%d", tl.Spans[a.From].Bar+1, changes.Bar(last.Start+last.Length-1)+1)
		line := fmt.Sprintf("%g bars, leaving %s (%s)", a.Bars, nm.of(a.Leaves, -1), distance(a))
		if a.First {
			line += ", the first tonality"
		}
		if a.Opens {
			line += ", opening a section"
		}
		if a.Closes {
			line += ", closing a section"
		}
		if a.True {
			line += ": true modulation"
		} else {
			line += ": transitory"
		}
		fmt.Fprintf(&b, "  %-6s %-12s %s\n", nm.of(a.Tonic, a.From), bars, line)
	}
	return b.String()
}

// distance writes how far an area is from the tonic it leaves, on the
// cycle of fifths: "1 step", "2 steps", or for a relative, which shares
// its key signature, "relative minor", "relative major".
func distance(a analysis.TonalArea) string {
	switch {
	case a.Distance == 0 && analysis.ModesOf(a.Leaves) == analysis.Minor:
		return "relative minor"
	case a.Distance == 0:
		return "relative major"
	case a.Distance == 1:
		return "1 step"
	}
	return fmt.Sprintf("%d steps", a.Distance)
}

// heard names the sensed tonic at a change, where it changes: the
// ground, when it changes or comes back after a region; in parentheses
// the region a cadence opens, the ground around it left unsaid; and in
// square brackets the tonic a dominant makes of the chord, for that
// chord only, when it is not the region's. In bars 7 and 8 of 'Round
// Midnight, heard in E♭m, A♭7 opens a region of D♭ and G♭7 B7 B♭7♯5
// read "(D♭) [G♭] [B]".
func heard(nm namer, sensed []analysis.Sensed, i int) string {
	name := func(t []harmony.Tonality) string { return nm.of(t, -1) }
	same := func(a, b []harmony.Tonality) bool {
		return (a == nil) == (b == nil) && (a == nil || name(a) == name(b))
	}
	s := sensed[i]
	var before analysis.Sensed
	if i > 0 {
		before = sensed[i-1]
	}
	var parts []string
	switch {
	case s.Ground == nil:
		if i == 0 || before.Ground != nil {
			parts = append(parts, "?")
		}
	case i == 0 || !same(s.Ground, before.Ground) || s.Region == nil && before.Region != nil:
		parts = append(parts, name(s.Ground))
	}
	if s.Region != nil && (i == 0 || !same(s.Region, before.Region)) {
		parts = append(parts, "("+nm.of(s.Region, i)+")")
	}
	if s.Tonicised != nil && (s.Region == nil || s.Tonicised[0].Tonic() != s.Region[0].Tonic()) {
		parts = append(parts, "["+nm.of(s.Tonicised, i)+"]")
	}
	return strings.Join(parts, " ")
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

// ruleEnd returns where a rule that ends on word `i` stops: at the end
// of its bar, padding included, unless another chord follows in the
// bar, where it stops a space before.
func ruleEnd(words []word, i, first, width int) int {
	z := words[i]
	if i+1 < len(words) && words[i+1].bar == z.bar {
		return column(z.bar, first, width) + words[i+1].off - 1
	}
	return column(z.bar, first, width) + width
}

// A rule is a name over a stretch of changes, ruled from its first to
// its last: a cell, a pedal.
type rule struct {
	from, to int
	text     string
}

// cellRules names the cells: "anatole", "anatole ♭II…".
func cellRules(formulas []analysis.Cell) []rule {
	var out []rule
	for _, cl := range formulas {
		text := cl.Kind.String()
		if cl.Substituted {
			text += " ♭II"
		}
		if !cl.Resolves {
			text += "…"
		}
		out = append(out, rule{cl.From, cl.To, text})
	}
	return out
}

// pedalRules names the pedals as En Harmonie writes them, "B♭ ped.",
// the bass spelled as the chords over it show it.
func pedalRules(pedals []analysis.Pedal, tl ireal.Timeline) []rule {
	var out []rule
	for _, p := range pedals {
		bass := ""
		for k := p.From; k <= p.To && bass == ""; k++ {
			bass = tl.Spans[k].Chord.Bass
		}
		out = append(out, rule{p.From, p.To, note(bass) + " ped."})
	}
	return out
}

// drawRules writes rules on a line of a row, each named at its first
// chord and ruled to the end of its last: "anatole ─────". A rule that
// goes on from the row before is drawn without its name.
func drawRules(rules []rule, words []word, first, end, width int) string {
	line := []rune(strings.Repeat(" ", column(end, first, width)))
	for _, r := range rules {
		a, z := words[r.from], words[r.to]
		if z.bar < first || a.bar >= end {
			continue
		}
		from, to := column(first, first, width), column(end, first, width)
		if a.bar >= first {
			from = column(a.bar, first, width) + a.off
		}
		if z.bar < end {
			to = ruleEnd(words, r.to, first, width)
		}
		name := []rune(r.text + " ")
		if a.bar < first {
			name = nil
		}
		for i := from; i < to && i < len(line); i++ {
			line[i] = '─'
			if k := i - from; k < len(name) {
				line[i] = name[k]
			}
		}
	}
	return strings.TrimRight(string(line), " ")
}

// drawBlocks draws the line of blocks over a row of bars: from the
// first chord of each block to its V, the tonality it announces then a
// rule. A block cut by the end of a row goes on in the next one, and a
// block across the loop, its two at the end of the chorus and its V at
// the start, is drawn in both places. A name that would run into the
// next block's is cut short.
func drawBlocks(nm namer, blocks []analysis.Block, tl ireal.Timeline, words []word, readings []string, first, end, width int) string {
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
			long, short := announced(nm, bl, tl)
			if readings[n] != "" {
				long += " : " + readings[n]
			}
			names = append(names, name{from, []rune(long + " "), []rune(short + " ")})
		}
		if z.bar >= first && z.bar < end {
			to = ruleEnd(words, bl.Five, first, width)
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

// announced names the tonalities a block announces: their tonic spelled
// as the chords under it are (see namer.tonic), then their scales: "E♭" for E
// flat major, "Fm harm" for F harmonic minor, "D♭ M/m mel" when major
// and melodic minor both remain, "Gm nat/harm/mel" when the three
// minors do; "…" when the block does not resolve. The short name,
// for when the next block leaves no room, keeps the mode alone: "Gm",
// "D♭(m)" when both remain.
func announced(nm namer, bl analysis.Block, tl ireal.Timeline) (string, string) {
	if len(bl.Announced) == 0 {
		return "?", "?"
	}
	at := bl.Five
	if bl.Target >= 0 {
		at = bl.Target
	}
	tonic := nm.tonic(bl.Announced, at)
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

func (c cell) width() int {
	return max(utf8.RuneCountInString(c.top), utf8.RuneCountInString(c.bottom),
		utf8.RuneCountInString(c.below), utf8.RuneCountInString(c.tonic))
}

func pad(s string, w int) string {
	return s + strings.Repeat(" ", max(0, w-utf8.RuneCountInString(s)))
}

// legend explains the marks over the chords, one to a line.
const legend = `
over each chord, how it prepares the next:
  V→               as its dominant
  ♭II→             as its chromatic dominant
  °→               as its diminished chord, a dominant without its root
  sus→             as its suspension
  II→              as the two of a dominant, a fifth above it, or a
                   half tone above a chromatic dominant, as the two of
                   its tritone twin: Fm7 E7 (the two of B♭7)
  IV→              as its subdominant, in a plagal cadence (IV, ♭VII7
                   or IIm7♭5)
  ↗ ↘              a passing chord, its bass walking up or down
  ½ step 5th II→   a II-V following the one before it by a half tone,
                   a tone, or down the cycle of fifths
  5th V→  ½ V→     a dominant following the one it resolves, down the
                   cycle of fifths, or a half tone below (chromatic)

above, the blocks ([II] [sus4] V) and the tonality each announces:
  Fm harm ──       F harmonic minor
  D♭ M/m mel… ──   D♭ major or melodic minor, not resolved (deceptive)
  Fm harm : II V   the block read in that tonality, as En Harmonie
                   brackets it, when it differs from the degrees below

under the chords, the pedals, a bass held under chords on other roots:
  B♭ ped. ───      over Fm7/B♭ B♭7 E♭maj7/B♭

under each chord, its degree in the tonality of the passage:
  IIIm7♭5 VI7      Gm7♭5 C7 in E♭ major
  ♭VII7 ♯Vdim7 I/3 a borrowed chord, a passing chord, an inversion
  modal            a chord held four bars or more with no cadence: a
                   modal plage, whose mode the chart does not say

under the degrees, the tonic the ear hears, where it changes:
  E♭               E♭ is the tonic
  (Fm)             a cadence opens a region in F minor, a transitory
                   modulation, inside the tonic named before
  [Fm]             a dominant makes of this chord a I of F minor, for
                   this chord only: a tonicisation

under the tonics, the cells, formulas heard as one:
  anatole ───      I VI II V, in any colour
  III-VI-II-V-I ─  the same, the III standing for the I, then the I
  anatole ♭II ─    with some X7 standing for their tritone twin
  anatole… ─       its V not going to its I, the anatole being cyclic
  ♭VI-♭VII-I ─     the aeolian cadence, on a minor I, or borrowed on a
                   major one (the Mario Cadence)
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
		{harmony.PlagalApproach, "IV"},
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

// style is how chord symbols are written, set once from the flags.
var style naming.ChordStyle

// spelling writes the app's qualities the way a lead sheet does, for
// the few that naming cannot render: a quality typed by hand that the
// reader does not know.
var spelling = strings.NewReplacer(
	"^", "maj", "-", "m", "h7", "m7♭5", "h9", "m9♭5", "h", "m7♭5",
	"o", "dim", "b", "♭", "#", "♯",
)

func symbol(c ireal.ChordSymbol) string {
	return note(c.Root) + quality(c) + slash(c)
}

// quality writes the quality of a chord symbol in the chosen style. The
// root, which respell may have rewritten in signs, does not matter here:
// the quality is read on C.
func quality(c ireal.ChordSymbol) string {
	if r, err := (ireal.ChordSymbol{Root: "C", Quality: c.Quality, Custom: c.Custom}).Read(); err == nil {
		if q, ok := style.Symbol(r.Pattern); ok {
			return q
		}
	}
	return spelling.Replace(c.Quality)
}

// slash writes the bass of a chord symbol, "/E", or nothing.
func slash(c ireal.ChordSymbol) string {
	if c.Bass == "" {
		return ""
	}
	return "/" + note(c.Bass)
}

func note(n string) string {
	return strings.NewReplacer("b", "♭", "#", "♯").Replace(n)
}

// inPlage reports whether change i is in a modal plage.
func inPlage(plages []analysis.Plage, i int) bool {
	for _, p := range plages {
		if p.From <= i && i <= p.To {
			return true
		}
	}
	return false
}

func plural(n int) string {
	if n > 1 {
		return "s"
	}
	return ""
}
