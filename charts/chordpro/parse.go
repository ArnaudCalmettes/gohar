package chordpro

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parse reads a song in gohar's profile. Lines outside the grids and
// the blocks of melody, lyrics among them, are not the profile's: they
// are skipped. Errors say the line, and the song read up to them is
// kept.
func Parse(r io.Reader) (Song, error) {
	s := Song{Time: Time{Beats: 4, Unit: 4}, ABC: map[string]string{}}
	p := parser{song: &s}
	sc := bufio.NewScanner(r)
	var errs []error
	for n := 1; sc.Scan(); n++ {
		if err := p.line(strings.TrimSpace(sc.Text())); err != nil {
			errs = append(errs, fmt.Errorf("line %d: %w", n, err))
		}
	}
	if err := sc.Err(); err != nil {
		errs = append(errs, err)
	}
	if p.grid != nil || p.abc != nil {
		errs = append(errs, errors.New("chordpro: a section is not closed"))
	}
	return s, errors.Join(errs...)
}

// ParseString reads a song from a string.
func ParseString(s string) (Song, error) {
	return Parse(strings.NewReader(s))
}

type parser struct {
	song     *Song
	grid     *Grid            // the grid being read, nil outside one
	abc      *strings.Builder // the block of melody being read
	abcLabel string
}

// The directives of the profile, with their short names.
var aliases = map[string]string{
	"t": "title", "st": "subtitle",
	"sog": "start_of_grid", "eog": "end_of_grid",
}

func (p *parser) line(l string) error {
	if p.abc != nil {
		if name, _ := directive(l); name == "end_of_abc" {
			p.song.ABC[p.abcLabel] = p.abc.String()
			p.abc = nil
			return nil
		}
		p.abc.WriteString(l + "\n")
		return nil
	}
	if l == "" || strings.HasPrefix(l, "#") {
		return nil // a comment, in ChordPro
	}
	name, value := directive(l)
	switch {
	case name != "":
		return p.directive(name, value)
	case p.grid != nil:
		return p.gridLine(l)
	}
	return nil // lyrics, or anything the profile does not read
}

// directive splits "{name: value}" or "{name value}"; an empty name for
// a line that is none.
func directive(l string) (name, value string) {
	inner, ok := strings.CutPrefix(l, "{")
	if !ok {
		return "", ""
	}
	if inner, ok = strings.CutSuffix(inner, "}"); !ok {
		return "", ""
	}
	i := strings.IndexAny(inner, ": ")
	if i < 0 {
		name = inner
	} else {
		name, value = inner[:i], strings.TrimSpace(inner[i+1:])
	}
	name = strings.ToLower(strings.TrimSpace(name))
	if long, ok := aliases[name]; ok {
		name = long
	}
	return name, value
}

func (p *parser) directive(name, value string) error {
	s := p.song
	switch name {
	case "title":
		s.Title = value
	case "composer":
		s.Composers = append(s.Composers, value)
	case "copyright":
		s.Copyright = value
	case "key":
		s.Key = value
	case "time":
		beats, unit, ok := strings.Cut(value, "/")
		b, err1 := strconv.Atoi(beats)
		u, err2 := strconv.Atoi(unit)
		if !ok || err1 != nil || err2 != nil || b <= 0 || u <= 0 {
			return fmt.Errorf("chordpro: time %q", value)
		}
		s.Time = Time{Beats: b, Unit: u}
	case "tempo":
		t, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("chordpro: tempo %q", value)
		}
		s.Tempo = t
	case "meta":
		n, v, _ := strings.Cut(value, " ")
		s.Meta = append(s.Meta, Meta{Name: n, Value: strings.TrimSpace(v)})
	case "start_of_grid":
		if p.grid != nil {
			return errors.New("chordpro: a grid inside a grid")
		}
		attrs := attributes(value)
		p.grid = &Grid{Label: attrs["label"], Shape: attrs["shape"]}
	case "end_of_grid":
		if p.grid == nil {
			return errors.New("chordpro: end_of_grid outside a grid")
		}
		s.Body = append(s.Body, Item{Grid: p.grid})
		p.grid = nil
	case "x_coda":
		if s.Coda != 0 || len(s.Body) == 0 {
			return errors.New("chordpro: x_coda: one coda, after the chorus")
		}
		s.Coda = len(s.Body)
	case "x_play":
		if p.find(value) == nil {
			return fmt.Errorf("chordpro: x_play %q: no grid with that label before", value)
		}
		s.Body = append(s.Body, Item{Recall: value})
	case "start_of_abc":
		p.abc, p.abcLabel = &strings.Builder{}, attributes(value)["label"]
	default:
		s.Meta = append(s.Meta, Meta{Name: name, Value: value})
	}
	return nil
}

// find returns the grid labelled `label` read so far, nil for none.
func (p *parser) find(label string) *Grid {
	for _, it := range p.song.Body {
		if it.Grid != nil && it.Grid.Label == label {
			return it.Grid
		}
	}
	return nil
}

// attributes reads `label="A" shape="4x4"`. A bare value, as in older
// files, is the shape.
func attributes(s string) map[string]string {
	attrs := map[string]string{}
	for s = strings.TrimSpace(s); s != ""; s = strings.TrimSpace(s) {
		name, rest, ok := strings.Cut(s, "=")
		if !ok || strings.ContainsAny(name, " \"") {
			attrs["shape"] = s
			break
		}
		rest = strings.TrimSpace(rest)
		if v, after, ok := strings.Cut(strings.TrimPrefix(rest, `"`), `"`); strings.HasPrefix(rest, `"`) && ok {
			attrs[strings.TrimSpace(name)], s = v, after
			continue
		}
		v, after, _ := strings.Cut(rest, " ")
		attrs[strings.TrimSpace(name)], s = v, after
	}
	return attrs
}

// gridLine reads a row of bars: bar lines and cells, separated by
// spaces. What follows the last bar line is a comment, in ChordPro.
func (p *parser) gridLine(l string) error {
	fields := strings.Fields(l)
	if len(fields) == 0 || !isBarline(fields[0]) {
		return fmt.Errorf("chordpro: grid line %q: no bar line to start", l)
	}
	line := Line{Barlines: []Barline{Barline(fields[0])}}
	var cells []string
	for _, f := range fields[1:] {
		if !isBarline(f) {
			cells = append(cells, f)
			continue
		}
		bar, err := readBar(cells)
		if err != nil {
			return err
		}
		line.Bars = append(line.Bars, bar)
		line.Barlines = append(line.Barlines, Barline(f))
		cells = nil
	}
	if len(line.Bars) > 0 {
		p.grid.Lines = append(p.grid.Lines, line)
	}
	return nil
}

// isBarline tells a bar line: an optional ":" closing a repeat, the bar,
// then a second bar, a final dot, a ":" opening a repeat, or the number
// of an ending, and the ">" ChordPro allows after it.
func isBarline(f string) bool {
	f = strings.TrimPrefix(f, ":")
	rest, ok := strings.CutPrefix(f, "|")
	if !ok {
		return false
	}
	rest = strings.TrimSuffix(rest, ">")
	switch rest {
	case "", "|", ".", ":":
		return true
	}
	_, err := strconv.Atoi(rest)
	return err == nil
}

// readBar reads the cells of one bar.
func readBar(cells []string) (Bar, error) {
	if len(cells) == 1 {
		switch cells[0] {
		case "%":
			return Bar{Repeat: 1}, nil
		case "%%":
			return Bar{Repeat: 2}, nil
		}
	}
	if len(cells) == 0 {
		return Bar{}, errors.New("chordpro: an empty bar")
	}
	var b Bar
	for _, c := range cells {
		var cell Cell
		switch c {
		case ".":
		case "/":
			cell.Restrike = true
		default:
			for _, sym := range strings.Split(c, "~") {
				ch, err := ReadChord(sym)
				if err != nil {
					return b, err
				}
				cell.Chords = append(cell.Chords, ch)
			}
		}
		b.Cells = append(b.Cells, cell)
	}
	return b, nil
}
