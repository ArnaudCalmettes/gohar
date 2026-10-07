package grids

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// The grids of the game, as written: their titles and their length in
// bars, the repeat of Autumn Leaves played, the A of Satin Doll
// recalled.
func TestGrids(t *testing.T) {
	tunes, err := All()
	if err != nil {
		t.Fatal(err)
	}
	want := []struct {
		title string
		bars  int
	}{
		{"12 Bar Blues", 12},
		{"Satin Doll", 32},
		{"Tune Up", 16},
		{"Autumn Leaves", 32},
	}
	if len(tunes) != len(want) {
		t.Fatalf("%d grids, want %d", len(tunes), len(want))
	}
	for i, w := range want {
		if tunes[i].Title != w.title || len(tunes[i].Grid.Bars) != w.bars {
			t.Errorf("grid %d: %q in %d bars, want %q in %d", i, tunes[i].Title, len(tunes[i].Grid.Bars), w.title, w.bars)
		}
	}
}

// The blues in F, heard in F and moved to C, then to D♭: spelled letter
// by letter, G♭7 rather than F♯7.
func TestInKey(t *testing.T) {
	blues, err := Read(JazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	inC, err := InKey(blues, naming.SpelledNote{Letter: naming.LetterC})
	if err != nil {
		t.Fatal(err)
	}
	written := func(t Tune) string {
		var out []string
		for _, w := range t.Written {
			out = append(out, w.String())
		}
		return strings.Join(out, " ")
	}
	if got, want := written(inC), "C7 F7 C7 F7 C6 A7 Dm7 G7 C6 A7 D7 G7"; got != want {
		t.Errorf("in C: %s, want %s", got, want)
	}
	inDb, err := InKey(blues, naming.SpelledNote{Letter: naming.LetterD, Accidental: naming.FlatSign})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := written(inDb), "Db7 Gb7 Db7 Gb7 Db6 Bb7 Ebm7 Ab7 Db6 Bb7 Eb7 Ab7"; got != want {
		t.Errorf("in D♭: %s, want %s", got, want)
	}
}
