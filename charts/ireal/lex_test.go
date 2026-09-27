package ireal

import (
	"strings"
	"testing"
)

func join(ts []Token) string {
	var b strings.Builder
	for _, t := range ts {
		b.WriteString(t.Raw)
	}
	return b.String()
}

func kinds(ts []Token) []Kind {
	var k []Kind
	for _, t := range ts {
		k = append(k, t.Kind)
	}
	return k
}

// A chart with a bit of everything comes back whole, with no unknown
// token.
func TestLexLosesNothing(t *testing.T) {
	c := "[T44*ASC^7XyQKcl LZA-7(C-7) D7b9/F#LZ x <D.S. al Coda>]{N1G7sus,W/B}N2 n pp f Y*iQ<*48Fine>r|XyQ sC7alt lUZ "
	ts := Lex(c)
	if got := join(ts); got != c {
		t.Fatalf("round trip\ngot  %q\nwant %q", got, c)
	}
	for _, tok := range ts {
		if tok.Kind == Unknown {
			t.Errorf("unknown %q at %d", tok.Raw, tok.Pos)
		}
	}
}

func TestLexTokens(t *testing.T) {
	ts := Lex("[T34*BC-7b5/GbXyQKclLZ(E-)<*69Coda>N2")
	want := []Kind{DoubleOpen, TimeSignature, Section, Chord, EmptyCells, BarRepeated, Bar, Alternate, Comment, Ending}
	if got := kinds(ts); len(got) != len(want) {
		t.Fatalf("got %v\nwant %v", got, want)
	}
	for i, k := range want {
		if ts[i].Kind != k {
			t.Errorf("token %d: %v, want %v", i, ts[i].Kind, k)
		}
	}
	if ts[1].Time != (TimeSig{3, 4}) || ts[2].Section != 'B' || ts[9].Ending != 2 {
		t.Errorf("details %+v %+v %+v", ts[1], ts[2], ts[9])
	}
	if c := ts[3].Chord; c != (ChordSymbol{"C", "-7b5", "Gb", false}) {
		t.Errorf("chord %+v", c)
	}
	if c := ts[8].Comment; c.Text != "Coda" || c.Shift != 69 || !c.Shifted {
		t.Errorf("comment %+v", c)
	}
	if ts[7].Chord.Root != "E" || ts[7].Chord.Quality != "-" {
		t.Errorf("alternate %+v", ts[7].Chord)
	}
}

// Two digits leave no room for twelve: T12 is 12/8.
func TestLexTwelveEight(t *testing.T) {
	if got := Lex("T12")[0].Time; got != (TimeSig{12, 8}) {
		t.Errorf("T12 is %v", got)
	}
}

// The longest quality wins: 7sus is not 7 then a small s, 7alt is not
// 7 then an a.
func TestLexLongestQuality(t *testing.T) {
	for in, want := range map[string]ChordSymbol{
		"C7sus":     {"C", "7sus", "", false},
		"Bb7alt":    {"Bb", "7alt", "", false},
		"F#-7b5/C":  {"F#", "-7b5", "C", false},
		"Eb^7#11":   {"Eb", "^7#11", "", false},
		"G":         {"G", "", "", false},
		"W/B":       {"W", "", "B", false},
		"A7susadd3": {"A", "7susadd3", "", false},
		"A*m7*":     {"A", "m7", "", true},
		"Gb*o^7*/C": {"Gb", "o^7", "C", true},
	} {
		ts := Lex(in)
		if len(ts) != 1 || ts[0].Chord != want {
			t.Errorf("%s: %v %+v, want %+v", in, kinds(ts), ts[0].Chord, want)
		}
	}
}

// After a chord, s and l are sizes, not part of it.
func TestLexSizeAfterChord(t *testing.T) {
	if got := kinds(Lex("C7sC7l")); len(got) != 4 || got[1] != Small || got[3] != Large {
		t.Errorf("%v", got)
	}
}

// A rehearsal mark right after a chord is not a free quality.
func TestLexSectionAfterChord(t *testing.T) {
	if got := kinds(Lex("C*AD")); len(got) != 3 || got[1] != Section {
		t.Errorf("%v", got)
	}
}

// A line break saved in a chart is whitespace like a space.
func TestLexStrayWhitespace(t *testing.T) {
	c := "A-7XyQZ \n    "
	ts := Lex(c)
	if join(ts) != c {
		t.Fatal("round trip")
	}
	for _, tok := range ts[3:] {
		if tok.Kind != Space {
			t.Errorf("%q is %v", tok.Raw, tok.Kind)
		}
	}
}
