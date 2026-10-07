// Package grids holds the grids of Walk with me, in gohar's profile of
// ChordPro (see docs/formats.md): chords alone, written for gohar,
// under a free licence. Neither melodies nor grids that are someone
// else's.
package grids

import (
	"embed"
	"fmt"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

//go:embed *.cho
var files embed.FS

// JazzBlues is the first grid, and the one tests play most.
const JazzBlues = "blues.cho"

// SatinDollA is the A of Satin Doll alone, for the lessons.
const SatinDollA = "satin-doll-a.cho"

// Files are the grids the game offers, in the order the player goes
// round them: the blues first, then II-V in many keys, to vary the
// paths, then II-V-I in major, then in major and minor.
var Files = []string{JazzBlues, "satin-doll.cho", "tune-up.cho", "autumn-leaves.cho"}

// A Tune is a grid to play: its title, its changes, and beside each
// change the chord as the grid writes it, for the screen.
type Tune struct {
	Title   string
	Grid    analysis.Changes
	Written []chordpro.Chord
}

// Read reads the grid `name`.
func Read(name string) (Tune, error) {
	f, err := files.Open(name)
	if err != nil {
		return Tune{}, err
	}
	defer f.Close()
	s, err := chordpro.Parse(f)
	if err != nil {
		return Tune{}, fmt.Errorf("walk: %s: %w", name, err)
	}
	c, written, err := s.Spelled()
	if err != nil {
		return Tune{}, fmt.Errorf("walk: %s: %w", name, err)
	}
	return Tune{Title: s.Title, Grid: c, Written: written}, nil
}

// Changes reads the changes of the grid `name`.
func Changes(name string) (analysis.Changes, error) {
	t, err := Read(name)
	return t.Grid, err
}

// All reads every grid, in order.
func All() ([]Tune, error) {
	tunes := make([]Tune, len(Files))
	for i, name := range Files {
		t, err := Read(name)
		if err != nil {
			return nil, err
		}
		tunes[i] = t
	}
	return tunes, nil
}
