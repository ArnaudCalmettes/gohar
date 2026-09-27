package ireal

import (
	"fmt"
	"strings"
	"testing"
)

// timed renders a timeline as chords and their lengths in beats: "C^7:4
// D-7:2 G7:2", "n" for no chord.
func timed(chart string) string {
	var out []string
	for _, s := range Structure(Lex(chart)).Timeline().Spans {
		name := s.Chord.Root + s.Chord.Quality
		if s.Chord.Bass != "" {
			name += "/" + s.Chord.Bass
		}
		if s.NoChord {
			name = "n"
		}
		out = append(out, fmt.Sprintf("%s:%g", name, float64(s.Length)/float64(TicksPerBeat)))
	}
	return strings.Join(out, " ")
}

func TestTimeline(t *testing.T) {
	for name, c := range map[string]struct{ chart, want string }{
		"4/4, a cell a beat":         {"[T44C^7XyQ|D-7 G7LZC^7XyQZ", "C^7:4 D-7:2 G7:2 C^7:4"},
		"3/4 in four cells":          {"[T34CXyQ|D- G7LZCXyQZ", "C:3 D-:2 G7:1 C:3"},
		"two cells share the bar":    {"[CD-|EXyQZ", "C:2 D-:2 E:4"},
		"12/8 in dotted quarters":    {"[T12C D-LZEXyQZ", "C:2 D-:2 E:4"},
		"5/4 rounds up, Take Five":   {"[T54Eb- Bb-7LZEb-XyQZ", "Eb-:3 Bb-7:2 Eb-:5"},
		"empty and again go on":      {"[CXyQ|XyQ|CXyQ|DXyQZ", "C:12 D:4"},
		"x repeats":                  {"[C D-LZKcl LZEXyQZ", "C:2 D-:2 C:2 D-:2 E:4"},
		"a slash is the chord again": {"[CpD-p|EXyQZ", "C:2 D-:2 E:4"},
		"W takes the root before":    {"[C^7 W7LZDXyQZ", "C^7:2 C7:2 D:4"},
		"W walks the bass":           {"[C-7 W/BbLZDXyQZ", "C-7:2 C-7/Bb:2 D:4"},
		"W after no chord":           {"[n W/DLZDXyQZ", "n:4 D:4"},
		"silence before the first":   {"[XyQ|CXyQZ", "n:4 C:4"},
		"no chord":                   {"[C nLZDXyQZ", "C:2 n:2 D:4"},
	} {
		if got := timed(c.chart); got != c.want {
			t.Errorf("%s: %s, want %s", name, got, c.want)
		}
	}
}

func TestTimelineBars(t *testing.T) {
	tl := Structure(Lex("[T34CXyQ|DXyQ|T44EXyQZ")).Timeline()
	b := TicksPerBeat
	if want := []Ticks{0, 3 * b, 6 * b}; fmt.Sprint(tl.Bars) != fmt.Sprint(want) {
		t.Errorf("bars start at %v, want %v", tl.Bars, want)
	}
	if tl.Length != 10*b {
		t.Errorf("length %v, want %v", tl.Length, 10*b)
	}
	if tl.Next(len(tl.Spans)-1) != 0 {
		t.Error("the last span is not followed by the first")
	}
}

// The coda is played once, at the end: the chorus loops back before it.
func TestTimelineCoda(t *testing.T) {
	chart := "[CXyQ|SDXyQ|E Q XyQ|F<D.S. al Coda>XyQ][QXyQ|GXyQZ"
	if got, want := timed(chart), "C:4 D:4 E:4 F:4 D:4 E:4 E:4 G:4"; got != want {
		t.Fatalf("%s, want %s", got, want)
	}
	tl := Structure(Lex(chart)).Timeline()
	if tl.Coda != 6 {
		t.Errorf("the coda starts at span %d, want 6", tl.Coda)
	}
	for i, want := range map[int]int{4: 5, 5: 0, 6: 7, 7: -1} {
		if got := tl.Next(i); got != want {
			t.Errorf("after span %d comes %d, want %d", i, got, want)
		}
	}
}
