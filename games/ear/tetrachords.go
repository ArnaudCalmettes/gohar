package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// naturalTetrachords are the four shapes the natural system is built
// from: every one of its seven modes is two of them, a gap apart.
var naturalTetrachords = []harmony.Tetrachord{
	harmony.TetrachordMajor,
	harmony.TetrachordMinor,
	harmony.TetrachordPhrygian,
	harmony.TetrachordLydian,
}

// tetrachords is the step before the modes: four notes over a pedal,
// and the player names the shape.
//
// # Always the same choices, in the same order
//
// Four shapes, all four offered every time, in the same places. A mode
// among seven needs distractors picked for the level; a tetrachord
// among four does not, and fixed places let the hand learn where each
// answer lives while the ear does the work.
type tetrachords struct {
	shapes []harmony.Tetrachord
}

var _ Activity = tetrachords{}

func (a tetrachords) ID() string {
	return "tetrachords"
}

// tetrachordPattern lays a tetrachord out as a four note pattern on its tonic.
func tetrachordPattern(t harmony.Tetrachord) harmony.ScalePattern {
	x, y, z := t.Steps()
	set, err := harmony.NewPitchSet(0, harmony.PitchClass(x), harmony.PitchClass(x+y), harmony.PitchClass(x+y+z))
	if err != nil {
		return 0
	}
	p, _ := harmony.NewScalePattern(set, 0)
	return p
}

// Plan makes every shape come once before any comes twice.
func (a tetrachords) Plan(rng *rand.Rand, n int) []Question {
	var order []int
	for len(order) < n {
		order = append(order, rng.Perm(len(a.shapes))...)
	}
	questions := make([]Question, n)
	for i, shape := range order[:n] {
		questions[i] = a.question(shape, harmony.PitchClass(rng.IntN(12)))
	}
	return questions
}

func (a tetrachords) question(shape int, tonic harmony.PitchClass) Question {
	q := Question{Tonic: tonic, Answer: shape}
	for _, t := range a.shapes {
		q.Choices = append(q.Choices, dex.TetrachordOf(t))
	}
	return q
}

// Retry asks the same shape on another tonic.
func (a tetrachords) Retry(q Question, rng *rand.Rand) Question {
	return a.question(q.Answer, otherTonic(q.Tonic, rng))
}

func (a tetrachords) Prompt(l language, q Question) string {
	return fmt.Sprintf(l.words.whichShape, l.note(q.Tonic))
}

func (a tetrachords) Sound(q Question) []keyboard.Note {
	notes, _ := scaleNotes(q.Tonic, tetrachordPattern(q.Right().Tetrachord), 0)
	return notes
}

func (a tetrachords) Correction(q Question, chosen int) []keyboard.Note {
	right, end := scaleNotes(q.Tonic, tetrachordPattern(q.Right().Tetrachord), 0)
	wrong, _ := scaleNotes(q.Tonic, tetrachordPattern(q.Choices[chosen].Tetrachord), end+gap)
	return append(right, wrong...)
}

func (a tetrachords) Show(q Question, chosen int) display {
	return display{
		tonic:   q.Tonic,
		right:   tetrachordPattern(q.Right().Tetrachord),
		chosen:  tetrachordPattern(q.Choices[chosen].Tetrachord),
		mistake: chosen != q.Answer,
	}
}

func (a tetrachords) Facts(q Question) []dex.Fact {
	return []dex.Fact{
		{Kind: dex.FactNamed, Notion: q.Right(), Tonic: q.Tonic},
		{Kind: dex.FactHeard, Notion: q.Right(), Tonic: q.Tonic},
	}
}
