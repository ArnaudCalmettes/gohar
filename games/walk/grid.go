package main

import (
	"embed"
	"fmt"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The grids of the game, in gohar's profile of ChordPro (see
// docs/formats.md): chords alone, written for gohar, under a free
// licence. Neither melodies nor grids that are someone else's.
//
//go:embed grids/*.cho
var grids embed.FS

// The grids the game offers, in the order the player goes round them:
// the blues first, then II-V-I in major, then in major and minor.
const jazzBlues = "blues.cho"

var gridFiles = []string{jazzBlues, "tune-up.cho", "autumn-leaves.cho"}

// A tune is a grid to play: its title, its changes, and beside each
// change the chord as the grid writes it, for the screen.
type tune struct {
	title   string
	grid    analysis.Changes
	written []chordpro.Chord
}

// readTune reads the grid `name` of the game.
func readTune(name string) (tune, error) {
	f, err := grids.Open("grids/" + name)
	if err != nil {
		return tune{}, err
	}
	defer f.Close()
	s, err := chordpro.Parse(f)
	if err != nil {
		return tune{}, fmt.Errorf("walk: %s: %w", name, err)
	}
	c, written, err := s.Spelled()
	if err != nil {
		return tune{}, fmt.Errorf("walk: %s: %w", name, err)
	}
	return tune{title: s.Title, grid: c, written: written}, nil
}

// readGrid reads the changes of the grid `name`.
func readGrid(name string) (analysis.Changes, error) {
	t, err := readTune(name)
	return t.grid, err
}

// readTunes reads every grid of the game, in order.
func readTunes() ([]tune, error) {
	tunes := make([]tune, len(gridFiles))
	for i, name := range gridFiles {
		t, err := readTune(name)
		if err != nil {
			return nil, err
		}
		tunes[i] = t
	}
	return tunes, nil
}
