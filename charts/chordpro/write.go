package chordpro

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Write writes `s` in gohar's profile: the metadata, then the body, a
// grid per section and a recall where one plays again, the chords in
// ASCII. The cells are written as they are, shortcuts included.
func Write(w io.Writer, s Song) error {
	b := bufio.NewWriter(w)
	line := func(format string, args ...any) { fmt.Fprintf(b, format+"\n", args...) }

	if s.Title != "" {
		line("{title: %s}", s.Title)
	}
	for _, c := range s.Composers {
		line("{composer: %s}", c)
	}
	if s.Copyright != "" {
		line("{copyright: %s}", s.Copyright)
	}
	if s.Key != "" {
		line("{key: %s}", s.Key)
	}
	if s.Time != (Time{}) && s.Time != (Time{Beats: 4, Unit: 4}) {
		line("{time: %d/%d}", s.Time.Beats, s.Time.Unit)
	}
	if s.Tempo != 0 {
		line("{tempo: %d}", s.Tempo)
	}
	for _, m := range s.Meta {
		line("{meta: %s %s}", m.Name, m.Value)
	}
	for i, it := range s.Body {
		b.WriteString("\n")
		if s.Coda != 0 && i == s.Coda {
			line("{x_coda}")
		}
		if it.Grid == nil {
			line("{x_play: %s}", it.Recall)
			continue
		}
		g := it.Grid
		open := "{start_of_grid"
		if g.Label != "" {
			open += fmt.Sprintf(" label=%q", g.Label)
		}
		if g.Shape != "" {
			open += fmt.Sprintf(" shape=%q", g.Shape)
		}
		line("%s}", open)
		for _, l := range g.Lines {
			line("%s", l.String())
		}
		line("{end_of_grid}")
		if abc, ok := s.ABC[g.Label]; ok {
			line("{start_of_abc label=%q}", g.Label)
			b.WriteString(abc)
			line("\n{end_of_abc}")
		}
	}
	return b.Flush()
}

// String writes a row of the grid: "| F7 | Bb7 | F7 . . D7 | % |".
func (l Line) String() string {
	var sb strings.Builder
	for i, bar := range l.Bars {
		sb.WriteString(string(l.Barlines[i]) + " " + bar.String() + " ")
	}
	sb.WriteString(string(l.Barlines[len(l.Bars)]))
	return sb.String()
}

// String writes the cells of a bar, or its repeat sign.
func (b Bar) String() string {
	switch b.Repeat {
	case 1:
		return "%"
	case 2:
		return "%%"
	}
	cells := make([]string, len(b.Cells))
	for i, c := range b.Cells {
		cells[i] = c.String()
	}
	return strings.Join(cells, " ")
}

// String writes a cell: its chords joined with ~, a dot for a hold, a
// slash for a restrike.
func (c Cell) String() string {
	switch {
	case c.Restrike:
		return "/"
	case len(c.Chords) == 0:
		return "."
	}
	chords := make([]string, len(c.Chords))
	for i, ch := range c.Chords {
		chords[i] = ch.String()
	}
	return strings.Join(chords, "~")
}
