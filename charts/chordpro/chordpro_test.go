package chordpro

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// blues is the jazz blues of Walk with me, written in the profile.
const blues = `{title: 12 Bar Blues}
{meta: style Medium Swing}
{meta: chart_author gohar}
{meta: chart_license CC0-1.0}

{start_of_grid label="A" shape="4x4"}
| F7 | Bb7 | F7 | F7 |
| Bb7 | Bb7 | F6 | D7 |
| Gm7 | C7 | F6 D7~(Ab7) | G7 C7~(Gb7) |.
{end_of_grid}
`

// played says the changes as a musician reads a chart: bar and beat,
// counted from 1, then the chord, spelled in F.
func played(t *testing.T, c analysis.Changes) []string {
	t.Helper()
	namer, err := naming.NewNamer(naming.English)
	if err != nil {
		t.Fatal(err)
	}
	namer = namer.WithTonality(analysis.MajorTonalities(f)[0]) // these charts are in F, or in C
	var out []string
	for _, ch := range c.Chords {
		bar := c.Bar(ch.Start)
		beat := int((ch.Start-c.Bars[bar])/analysis.TicksPerBeat) + 1
		name := "N.C."
		if !ch.Silent {
			q, _ := style.Symbol(ch.Chord.Pattern)
			name = asciiSigns.Replace(namer.Name(ch.Chord.Root) + q)
			if ch.Bass != ch.Chord.Root {
				name += "/" + asciiSigns.Replace(namer.Name(ch.Bass))
			}
		}
		out = append(out, fmt.Sprintf("%d.%d %s", bar+1, beat, name))
	}
	return out
}

const f = 5 // the pitch class of F

func changes(t *testing.T, src string) analysis.Changes {
	t.Helper()
	s, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	c, err := s.Changes()
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func check(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("\n got %q\nwant %q", got, want)
	}
}

func TestReadChord(t *testing.T) {
	for in, want := range map[string]string{
		"Bbmaj7":     "Bbmaj7",
		"B♭maj7":     "Bbmaj7",
		"Am7b5":      "Am7b5",
		"D7(b9)":     "D7(b9)",
		"C7alt":      "C7alt",
		"F/A":        "F/A",
		"N.C.":       "N.C.",
		"(Db7)":      "(Db7)",
		"Gm6":        "Gm6",
		"Ebdim7":     "Ebdim7",
		"F6/9":       "F6/9",
		"C7sus4":     "C7sus4",
		"Cm(maj7)":   "Cm(maj7)",
		"F#m7":       "F#m7",
		"Bb7b5(b9)":  "Bb7b5(b9)",
		"G7sus(b9)":  "G7sus4(b9)",
		"D7#5(b9)":   "D7#5(b9)",
		"C7b5(#9)":   "C7b5(#9)",
		"C7(#11,b9)": "C7(b9,#11)",
		"Dm(maj7,9)": "Dm(maj7,9)",
		"F#sus2":     "F#sus2",
		"Dsus2/B":    "Dsus2/B",
	} {
		c, err := ReadChord(in)
		if err != nil {
			t.Errorf("%s: %v", in, err)
			continue
		}
		if got := c.String(); got != want {
			t.Errorf("%s: written %s, want %s", in, got, want)
		}
	}
	for _, in := range []string{"H7", "Cxyz", "(C7", "C/H"} {
		if _, err := ReadChord(in); err == nil {
			t.Errorf("%s: read without an error", in)
		}
	}
}

// The blues, a chord a bar, two in the last ones, the optional
// chromatic dominants left out.
func TestBlues(t *testing.T) {
	check(t, played(t, changes(t, blues)), []string{
		"1.1 F7", "2.1 Bb7", "3.1 F7", "5.1 Bb7", "7.1 F6", "8.1 D7",
		"9.1 Gm7", "10.1 C7", "11.1 F6", "11.3 D7", "12.1 G7", "12.3 C7",
	})
}

// A bar written beat by beat, a dot where nothing changes; two chords
// in one beat, joined with ~; a restrike.
func TestBeats(t *testing.T) {
	check(t, played(t, changes(t, "{sog}\n| F7 . . D7 | Gm7~C7 . / . |\n{eog}\n")), []string{
		"1.1 F7", "1.4 D7", "2.1 Gm7", "2.1 C7",
	})
}

// A repeat with two endings: C D E, then C D F.
func TestRepeatsAndEndings(t *testing.T) {
	check(t, played(t, changes(t, "{sog}\n|: C | D |1 E :|2 F |.\n{eog}\n")), []string{
		"1.1 C", "2.1 D", "3.1 E", "4.1 C", "5.1 D", "6.1 F",
	})
}

// An AABA: the A recalled, and % repeating the bar before.
func TestRecall(t *testing.T) {
	src := `{start_of_grid label="A"}
| C | Am7 | Dm7 | G7 |.
{end_of_grid}
{x_play: A}
{start_of_grid label="B"}
| E7 | % | A7 | % |.
{end_of_grid}
{x_play: A}
`
	s, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	bars, err := s.Played()
	if err != nil {
		t.Fatal(err)
	}
	if len(bars) != 16 {
		t.Errorf("%d bars played, want 16", len(bars))
	}
	got := played(t, changes(t, src))
	if got[len(got)-1] != "16.1 G7" || !slices.Contains(got, "9.1 E7") || !slices.Contains(got, "11.1 A7") {
		t.Errorf("got %q", got)
	}
}

// Written again, the song reads the same, and writes the same.
func TestWriteThenRead(t *testing.T) {
	s, err := ParseString(blues)
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err := Write(&b, s); err != nil {
		t.Fatal(err)
	}
	again, err := ParseString(b.String())
	if err != nil {
		t.Fatalf("%v\n%s", err, b.String())
	}
	var b2 bytes.Buffer
	Write(&b2, again)
	if b.String() != b2.String() {
		t.Errorf("written twice, not the same:\n%s\n%s", b.String(), b2.String())
	}
	c1, _ := s.Changes()
	c2, _ := again.Changes()
	check(t, played(t, c2), played(t, c1))
	if !strings.Contains(b.String(), "| F6 D7~(Ab7) |") {
		t.Errorf("the optional chord is lost:\n%s", b.String())
	}
}

// walkBlues is the iReal link of the same blues.
const walkBlues = "irealb://12%20Bar%20Blues=Composer%20Unknown==Medium%20Swing=F=5=1r34LbKcu7X7D%7CQQ%7CBb7QyX%2C7bB%7CQyX7bBQ%7CyX7F%7CQyX7F%7CQyX%7CF6XyyX7F%5ByQ%7CG%2D7XyQ%7CC7XyQ%7CF6%20D7%28Ab7%29LZG7%20C7%28Gb7%29%20%5D%20=Jazz%2DMedium%20Swing=100=30"

// Converted from iReal, written and read again, a chart gives the
// analysis the same changes as the iReal chart itself.
func TestFromIRealSameChanges(t *testing.T) {
	p, err := ireal.Parse(walkBlues)
	if err != nil {
		t.Fatal(err)
	}
	src := p.Songs[0]
	want, err := ireal.Structure(ireal.Lex(src.Chart)).Changes()
	if err != nil {
		t.Fatal(err)
	}
	s, err := FromIReal(src)
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err := Write(&b, s); err != nil {
		t.Fatal(err)
	}
	got := changes(t, b.String())
	check(t, played(t, got), played(t, want))
	if !strings.Contains(b.String(), "(Ab7)") {
		t.Errorf("the chord written small is lost:\n%s", b.String())
	}
}

// A D.S. al Coda, converted: the jump unfolded, the coda after
// {x_coda}, starting a change of its own on the chord held into it,
// as in the iReal chart.
func TestFromIRealCoda(t *testing.T) {
	src := ireal.Song{Title: "Coda", Chart: "[CXyQ|SDXyQ|E Q XyQ|F<D.S. al Coda>XyQ][QXyQ|GXyQZ"}
	want, err := ireal.Structure(ireal.Lex(src.Chart)).Changes()
	if err != nil {
		t.Fatal(err)
	}
	s, err := FromIReal(src)
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err := Write(&b, s); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "{x_coda}") {
		t.Errorf("no coda written:\n%s", b.String())
	}
	got := changes(t, b.String())
	check(t, played(t, got), played(t, want))
	if got.Coda != want.Coda || got.End != want.End {
		t.Errorf("coda %d, end %d; want %d, %d\n%s", got.Coda, got.End, want.Coda, want.End, b.String())
	}
}

// The chords as the grid writes them, beside the changes: the C of
// Cmaj7 and the B♭ of Bbmaj7 in a tune heard in D.
func TestSpelled(t *testing.T) {
	s, err := ParseString("{sog}\n| Em7 | A7 | Dmaj7 | Cmaj7 | Bbmaj7 |\n{eog}\n")
	if err != nil {
		t.Fatal(err)
	}
	c, written, err := s.Spelled()
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, ch := range written {
		got = append(got, ch.String())
	}
	check(t, got, []string{"Em7", "A7", "Dmaj7", "Cmaj7", "Bbmaj7"})
	if len(written) != len(c.Chords) {
		t.Errorf("%d chords written for %d changes", len(written), len(c.Chords))
	}
}
