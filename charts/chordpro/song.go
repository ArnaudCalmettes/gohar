// Package chordpro reads and writes gohar's profile of ChordPro: the
// grid of a jazz tune, its form and its chords, in an open text format
// other software reads too (see docs/formats.md).
//
// The profile is a subset of ChordPro and a few extensions: the
// metadata, one grid section per section of the tune, the bars and
// beats of the grid with its repeats and endings, a recall of a section
// with {x_play}, the coda after {x_coda}, the chords as musicians write
// them, and the optional chords between parentheses. The ABC blocks of
// the melody are kept as they are, for a reader to come.
package chordpro

// A Song is one file: its metadata and its body, in the order written.
type Song struct {
	Title     string
	Composers []string
	Copyright string // the composition's
	Key       string // as written: an annotation, never read by the analysis
	Time      Time
	Tempo     int // 0 when unsaid

	// Meta holds the other metadata, {meta: name value} and the
	// directives the profile does not name, in the order written:
	// chart_author and chart_license among them.
	Meta []Meta

	// Body is the tune in the order it is played: each grid plays where
	// it is defined, each recall plays a grid defined before it.
	Body []Item

	// Coda is where the coda starts in the body, after {x_coda}: the
	// chorus goes round before it, and it is played once, on the last
	// chorus, to conclude. 0 when the tune has none.
	Coda int

	// ABC holds the blocks of melody as written, by label.
	ABC map[string]string
}

// A Meta is one piece of metadata.
type Meta struct {
	Name, Value string
}

// Time is the time signature, 4/4 by default.
type Time struct {
	Beats, Unit int
}

// pulses is how many beats a bar counts, the beat being the pulse:
// compound meters (6/8, 9/8, 12/8) count in dotted quarters, as in
// charts/ireal.
func (t Time) pulses() int {
	if t.Unit == 8 && t.Beats > 3 && t.Beats%3 == 0 {
		return t.Beats / 3
	}
	return t.Beats
}

// An Item of the body is a grid, or the recall of one.
type Item struct {
	Grid   *Grid  // nil for a recall
	Recall string // the label of the grid played again
}

// A Grid is one section of the tune, line by line as written.
type Grid struct {
	Label string
	Shape string // as written, kept for the writer
	Lines []Line
}

// A Line is a row of the grid: its bars, and the bar lines around them,
// one more than the bars.
type Line struct {
	Bars     []Bar
	Barlines []Barline
}

// A Barline is the line between two bars, as written: "|", "||", "|.",
// "|:", ":|", ":|:", and the endings "|1", ":|2".
type Barline string

// An ending reads a bar line: whether it closes a repeat, opens one,
// and the ending it opens, 0 for none.
func (b Barline) parts() (closes, opens bool, ending int) {
	s := string(b)
	if len(s) > 1 && s[0] == ':' {
		closes, s = true, s[1:]
	}
	s = s[1:] // the bar itself
	switch {
	case s == ":":
		opens = true
	case len(s) > 0 && s[0] >= '1' && s[0] <= '9':
		ending = int(s[0] - '0')
	}
	return closes, opens, ending
}

// double tells a double or final bar line, where an ending stops.
func (b Barline) double() bool {
	return b == "||" || b == "|."
}

// A Bar is the cells of one measure, as written, or the repeat of the
// bars before it.
type Bar struct {
	Cells  []Cell
	Repeat int // 1 for %, 2 for %%: the cells are empty
}

// A Cell is what a bar holds at one place: a chord or several, joined
// with ~ in the same beat, or the chord before going on.
type Cell struct {
	Chords   []Chord // empty for a hold
	Restrike bool    // "/": the chord before, struck again
}

// held tells a cell where no new chord starts: a hold, a restrike, or
// optional chords only, which the player may add but the tune does not
// need.
func (c Cell) held() bool {
	for _, ch := range c.Chords {
		if !ch.Optional {
			return false
		}
	}
	return true
}
