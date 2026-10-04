package chordpro

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
)

// barsPerLine is how many bars the converter writes on a row, as a Real
// Book page does.
const barsPerLine = 4

// ErrTime is returned for a chart whose time signature changes, which
// the profile does not write yet.
var ErrTime = errors.New("chordpro: a change of time signature, not in the profile yet")

// FromIReal converts an iReal Pro song into the profile, for a private
// corpus: the grid an iReal chart holds is its author's.
//
// The form is written as it is played, once through: the repeats, the
// endings and the jumps of a D.S. or a D.C. are unfolded, a section
// starts at each rehearsal mark, and the coda after {x_coda}. The
// chords are the timeline's, so that the song gives the analysis the
// same changes as the chart (see [Song.Changes]), and the chords the
// app writes small above another become optional chords,
// joined to it with ~. A bar where no chord starts writes the chord
// that goes on, which reads better than a dot.
//
// A chord that cannot be read becomes N.C., and the error says which,
// as in charts/ireal.
func FromIReal(src ireal.Song) (Song, error) {
	s := Song{
		Title: src.Title,
		Key:   src.Key,
		Tempo: src.BPM,
		ABC:   map[string]string{},
	}
	if src.Composer != "" && src.Composer != "Composer Unknown" { // the app's placeholder
		s.Composers = []string{src.Composer}
	}
	if src.Style != "" {
		s.Meta = append(s.Meta, Meta{Name: "style", Value: src.Style})
	}

	chart := ireal.Structure(ireal.Lex(src.Chart))
	played := chart.Unfold()
	if len(played) == 0 {
		return s, errors.New("chordpro: an empty chart")
	}
	first := chart.Measures[played[0].Index].Time
	for _, p := range played[1:] {
		if chart.Measures[p.Index].Time != first {
			return s, ErrTime
		}
	}
	s.Time = Time{Beats: first.Beats, Unit: first.Unit}
	beats := s.Time.pulses()

	// The chords of each bar, beat by beat, from the timeline.
	t := chart.Timeline()
	slots := make([][]Cell, len(t.Bars))
	for b := range slots {
		slots[b] = make([]Cell, beats)
	}
	var errs []error
	sounding := Chord{NoChord: true}
	for _, sp := range t.Spans {
		ch := Chord{NoChord: true}
		if !sp.NoChord {
			r, err := sp.Chord.Read()
			if err != nil {
				errs = append(errs, fmt.Errorf("bar %d: %v: %w", sp.Bar+1, sp.Chord, err))
			} else {
				ch = Chord{Root: r.Root, Pattern: r.Pattern, Bass: r.Bass, HasBass: r.HasBass}
			}
		}
		beat := int((sp.Start - t.Bars[sp.Bar]) / ireal.TicksPerBeat)
		slots[sp.Bar][beat].Chords = []Chord{ch}
	}
	for b := range slots {
		for beat, c := range slots[b] {
			switch {
			case len(c.Chords) > 0:
				sounding = c.Chords[0]
			case beat == 0:
				slots[b][0].Chords = []Chord{sounding}
			}
		}
	}

	// The chords written small above another: optional.
	for b, p := range played {
		for _, e := range p.Events {
			if !e.HasAlternate {
				continue
			}
			r, err := e.Alternate.Read()
			if err != nil {
				errs = append(errs, fmt.Errorf("bar %d: %v: %w", b+1, e.Alternate, err))
				continue
			}
			beat := cellBeat(e.Cell, chart.Measures[p.From].Cells, beats)
			cell := &slots[b][beat]
			if len(cell.Chords) == 0 {
				cell.Chords = []Chord{soundingAt(slots, b, beat)}
			}
			cell.Chords = append(cell.Chords, Chord{Root: r.Root, Pattern: r.Pattern, Bass: r.Bass, HasBass: r.HasBass, Optional: true})
		}
	}

	// Sections at the rehearsal marks and at the coda, four bars a row.
	var g *Grid
	for b, p := range played {
		if mark := chart.Measures[p.Index].Section; mark != 0 || g == nil || p.Coda {
			label := ""
			switch {
			case mark != 0:
				label = string(mark)
			case p.Coda:
				label = "Coda"
			}
			if p.Coda {
				s.Coda = len(s.Body)
			}
			g = &Grid{Label: label, Shape: fmt.Sprintf("%dx%d", barsPerLine, beats)}
			s.Body = append(s.Body, Item{Grid: g})
		}
		n := len(g.Lines)
		if n == 0 || len(g.Lines[n-1].Bars) == barsPerLine {
			g.Lines = append(g.Lines, Line{Barlines: []Barline{"|"}})
			n++
		}
		l := &g.Lines[n-1]
		l.Bars = append(l.Bars, Bar{Cells: shorten(slots[b])})
		l.Barlines = append(l.Barlines, "|")
	}
	for i, it := range s.Body {
		last := &it.Grid.Lines[len(it.Grid.Lines)-1]
		last.Barlines[len(last.Barlines)-1] = "||"
		if i == len(s.Body)-1 {
			last.Barlines[len(last.Barlines)-1] = "|."
		}
	}
	return s, errors.Join(errs...)
}

// cellBeat is the beat a cell of a measure `cells` wide falls on, as
// charts/ireal counts it: rounded up, never past the last beat.
func cellBeat(cell, cells, beats int) int {
	if cells == 0 {
		return 0
	}
	return min((cell*beats+cells-1)/cells, beats-1)
}

// soundingAt returns the chord sounding at a beat of a bar.
func soundingAt(slots [][]Cell, bar, beat int) Chord {
	for b := bar; b >= 0; b-- {
		for i := len(slots[b]) - 1; i >= 0; i-- {
			if b == bar && i > beat {
				continue
			}
			if c := slots[b][i]; len(c.Chords) > 0 {
				return c.Chords[0]
			}
		}
	}
	return Chord{NoChord: true}
}

// shorten writes a bar with the shortcuts musicians use: one cell when
// a chord lasts the bar, two when two chords share it in halves, a cell
// per beat otherwise.
func shorten(beats []Cell) []Cell {
	var at []int
	for i, c := range beats {
		if len(c.Chords) > 0 {
			at = append(at, i)
		}
	}
	half := len(beats) / 2
	switch {
	case len(at) == 1 && at[0] == 0:
		return beats[:1]
	case len(at) == 2 && at[0] == 0 && at[1] == half && len(beats)%2 == 0:
		return []Cell{beats[0], beats[half]}
	}
	return beats
}

// Name returns a file name for a song: its title, lower case, the
// spaces as dashes, ".cho".
func Name(title string) string {
	var sb strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			sb.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			sb.WriteRune('-')
		}
	}
	return sb.String() + ".cho"
}
