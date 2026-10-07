package grids

import "testing"

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
