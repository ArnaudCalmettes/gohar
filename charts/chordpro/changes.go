package chordpro

import (
	"errors"
	"fmt"

	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// Played returns the bars of the song in the order they sound, once
// through, the coda included: the body in order, each recall playing
// its grid again, each repeat taken once, the first ending on the first
// time, the second on the second. A % plays the bar before it again, %%
// the two bars before.
func (s Song) Played() ([]Bar, error) {
	bars, _, _, err := s.played()
	return bars, err
}

// A Rehearsal is a rehearsal mark as the chart writes it: the label of
// a grid, "A", "B", at the bar it starts from, an index in the bars
// played. A grid played again by {x_play} is marked again.
type Rehearsal struct {
	Label string
	Bar   int
}

// Rehearsals returns the rehearsal marks of the song, in the order they
// sound, as written: a chart's own reading of its form, which nothing
// checks (see analysis.Sections for the form found in the chords). A
// grid without a label has none.
func (s Song) Rehearsals() ([]Rehearsal, error) {
	_, _, marks, err := s.played()
	return marks, err
}

// played returns the bars as Played does, the first bar of the coda,
// -1 when there is none, and the rehearsal marks.
func (s Song) played() ([]Bar, int, []Rehearsal, error) {
	var out []Bar
	var marks []Rehearsal
	coda := -1
	for i, it := range s.Body {
		if s.Coda != 0 && i == s.Coda {
			coda = len(out)
		}
		g := it.Grid
		if g == nil {
			if g = s.grid(it.Recall); g == nil {
				return out, coda, marks, fmt.Errorf("chordpro: x_play %q: no such grid", it.Recall)
			}
		}
		if g.Label != "" {
			marks = append(marks, Rehearsal{Label: g.Label, Bar: len(out)})
		}
		bars, err := unfold(g)
		if err != nil {
			return out, coda, marks, fmt.Errorf("chordpro: grid %q: %w", g.Label, err)
		}
		for _, b := range bars {
			switch b.Repeat {
			case 0:
				out = append(out, b)
			default:
				if len(out) < b.Repeat {
					return out, coda, marks, errors.New("chordpro: a repeat sign with nothing before it")
				}
				out = append(out, out[len(out)-b.Repeat:]...)
			}
		}
	}
	return out, coda, marks, nil
}

func (s Song) grid(label string) *Grid {
	for _, it := range s.Body {
		if it.Grid != nil && it.Grid.Label == label {
			return it.Grid
		}
	}
	return nil
}

// A written bar, with the bar lines that open and close it.
type written struct {
	bar         Bar
	open, close Barline
}

// unfold plays the repeats and endings of one grid. A repeat goes back
// to the bar line that opened it, or to the start of the grid; an
// ending lasts from its bar line to the one that closes the repeat,
// the second to a double or final bar line, or the end of the grid.
func unfold(g *Grid) ([]Bar, error) {
	var bars []written
	for _, l := range g.Lines {
		for i, b := range l.Bars {
			bars = append(bars, written{bar: b, open: l.Barlines[i], close: l.Barlines[i+1]})
		}
	}
	var out []Bar
	start, pass, ending := 0, 1, 0
	for i, guard := 0, 0; i < len(bars); i++ {
		if guard++; guard > 4*len(bars) {
			return out, errors.New("repeats that never end")
		}
		w := bars[i]
		_, opens, e := w.open.parts()
		if opens && i != start {
			start, pass = i, 1
		}
		if e != 0 {
			ending = e
		}
		closes, _, _ := w.close.parts()
		if ending == 0 || ending == pass {
			out = append(out, w.bar)
			if closes && pass == 1 {
				pass, ending, i = 2, 0, start-1
				continue
			}
		}
		switch {
		case closes && pass == 2 && ending == 0:
			pass, start = 1, i+1
		case closes:
			ending = 0 // the first ending is over
		case w.close.double() && ending != 0:
			pass, ending, start = 1, 0, i+1
		}
	}
	return out, nil
}

// Changes times the chords of the song for the analysis: each chord
// with its bass, when it starts and how long it sounds, looping as a
// chart is played. A chord written again the same goes on, as in
// charts/ireal, and so do the optional chords: they are the player's
// to add, the tune does not need them.
//
// A coda starts a change of its own, even on a chord held into it, and
// the tune ends on its last chord, as in charts/ireal.
//
// # Cells to beats
//
// A bar of as many cells as beats has a cell per beat. A bar of fewer
// cells shares the bar evenly, rounded up to the next beat, as charts/ireal
// does: | F7 | lasts the bar, | Gm7 C7 | is two halves, the shortcut
// musicians write. Chords joined with ~ share their cell evenly.
func (s Song) Changes() (analysis.Changes, error) {
	c, _, err := s.Spelled()
	return c, err
}

// Spelled returns the changes as Changes does, and beside them the
// chords as the grid writes them, one for each change: the B♭ of
// Bbmaj7, never an A♯. A game shows them so; the analysis spells by the
// tonality it hears, and a grid that modulates needs more than one.
func (s Song) Spelled() (analysis.Changes, []Chord, error) {
	bars, coda, _, err := s.played()
	if err != nil {
		return analysis.Changes{}, nil, err
	}
	var written []Chord
	beats := s.Time.pulses()
	barLength := analysis.Ticks(beats) * analysis.TicksPerBeat
	c := analysis.Changes{Loops: true, Bars: make([]analysis.Ticks, len(bars))}
	var last *Chord
	add := func(ch Chord, at analysis.Ticks) {
		if last != nil && last.same(ch) {
			return
		}
		if n := len(c.Chords); n > 0 && c.Chords[n-1].Start == at {
			c.Chords = c.Chords[:n-1] // two on one beat: the later sounds
			written = written[:n-1]
		}
		if n := len(c.Chords); n > 0 {
			c.Chords[n-1].Length = at - c.Chords[n-1].Start
		}
		change := analysis.Change{Start: at, Silent: ch.NoChord}
		if !ch.NoChord {
			change.Chord = ch.Harmony()
			change.Bass = change.Chord.Root
			if ch.HasBass {
				change.Bass = ch.Bass.Class()
			}
		}
		c.Chords = append(c.Chords, change)
		written = append(written, ch)
		chord := ch
		last = &chord
	}
	for i, b := range bars {
		bar := analysis.Ticks(i) * barLength
		c.Bars[i] = bar
		if i == coda && last != nil {
			held := *last
			last = nil
			add(held, bar) // the coda starts its own change
			c.Coda = len(c.Chords) - 1
		}
		n := len(b.Cells)
		for j, cell := range b.Cells {
			if cell.held() {
				continue
			}
			at := bar + analysis.Ticks(min((j*beats+n-1)/n, beats-1))*analysis.TicksPerBeat
			if n == beats {
				at = bar + analysis.Ticks(j)*analysis.TicksPerBeat
			}
			var chords []Chord
			for _, ch := range cell.Chords {
				if !ch.Optional {
					chords = append(chords, ch)
				}
			}
			for k, ch := range chords {
				add(ch, at+analysis.Ticks(k)*analysis.TicksPerBeat/analysis.Ticks(len(chords)))
			}
		}
	}
	if len(c.Chords) == 0 {
		return c, nil, errors.New("chordpro: no chord")
	}
	end := analysis.Ticks(len(bars)) * barLength
	c.Chords[len(c.Chords)-1].Length = end - c.Chords[len(c.Chords)-1].Start
	if c.Chords[0].Start > 0 {
		// The silence before the first chord, as charts/ireal keeps it.
		c.Chords = append([]analysis.Change{{Start: 0, Length: c.Chords[0].Start, Silent: true}}, c.Chords...)
		written = append([]Chord{{NoChord: true}}, written...)
		if c.Coda > 0 {
			c.Coda++
		}
	}
	if c.Coda > 0 {
		c.End = len(c.Chords) - 1
	}
	return c, written, nil
}
