package main

import (
	"math/rand/v2"
	"testing"
)

// The grids of the game, as written: their titles and their length in
// bars, the repeat of Autumn Leaves played.
func TestGrids(t *testing.T) {
	tunes, err := readTunes()
	if err != nil {
		t.Fatal(err)
	}
	want := []struct {
		title string
		bars  int
	}{
		{"12 Bar Blues", 12},
		{"Tune Up", 16},
		{"Autumn Leaves", 32},
	}
	if len(tunes) != len(want) {
		t.Fatalf("%d grids, want %d", len(tunes), len(want))
	}
	for i, w := range want {
		if tunes[i].title != w.title || len(tunes[i].grid.Bars) != w.bars {
			t.Errorf("grid %d: %q in %d bars, want %q in %d", i, tunes[i].title, len(tunes[i].grid.Bars), w.title, w.bars)
		}
	}
}

// On Tune Up and Autumn Leaves as on the blues, the reference line lands
// every arrival: two chords a bar, II-V-I in major and in minor.
func TestWalkLandsEveryGrid(t *testing.T) {
	tunes, err := readTunes()
	if err != nil {
		t.Fatal(err)
	}
	for _, tu := range tunes[1:] {
		beats := Expect(tu.grid, at120(), 2)
		for seed := range uint64(20) {
			line := Walk(beats, rand.New(rand.NewPCG(seed, 1)))
			m := at120()
			k := NewMarker(FirstPalier, m, beats)
			for n, key := range line {
				k.Play(Note{key, m.At(n)})
			}
			for _, bm := range k.Close(m.At(len(beats))) {
				if bm.Kind != Landed {
					p := m.Position(bm.Beat)
					t.Errorf("%s, seed %d: bar %d, beat %d marked %d, want landed", tu.title, seed, p.Bar, p.Beat, bm.Kind)
				}
			}
		}
	}
}
