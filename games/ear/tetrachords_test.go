package main

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

var naturalShapes = tetrachords{shapes: naturalTetrachords}

func TestTetrachordsCoverTheirShapes(t *testing.T) {
	questions := naturalShapes.Plan(rand.New(rand.NewPCG(1, 2)), 10)
	seen := map[dex.Notion]bool{}
	for i, q := range questions {
		seen[q.Right()] = true
		if len(q.Choices) != 4 {
			t.Errorf("question %d offers %d choices, want all 4", i, len(q.Choices))
		}
		if !slices.Equal(q.Choices, questions[0].Choices) {
			t.Errorf("question %d moves the choices around", i)
		}
	}
	if len(seen) != 4 {
		t.Errorf("%d shapes asked in 10 questions, want all 4", len(seen))
	}
}

// The major tetrachord on its tonic is do ré mi fa.
func TestTetrachordPattern(t *testing.T) {
	var got []harmony.Semitones
	for _, n := range tetrachordPattern(harmony.TetrachordMajor).Offsets() {
		got = append(got, n)
	}
	if !slices.Equal(got, []harmony.Semitones{0, 2, 4, 5}) {
		t.Errorf("major tetrachord %v, want 0 2 4 5", got)
	}
}

// Four notes, no octave added: a tetrachord is not a scale.
func TestTetrachordSoundsFourNotesOverThePedal(t *testing.T) {
	q := naturalShapes.question(0, 0)
	if notes := naturalShapes.Sound(q); len(notes) != 5 {
		t.Errorf("%d notes, want the pedal and four", len(notes))
	}
}
