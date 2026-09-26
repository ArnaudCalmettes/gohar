package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// farEnough is the least number of notes a distractor must differ by,
// at the same tonic. Lydian against ionian differs by one and is kept
// for a later level.
const farEnough = 2

// modes is the first activity: a mode sounds over its tonic, and the
// player names it among three.
//
// # All seven, in the order of the degrees
//
// With `all`, the level above: the seven modes every time, each in the
// place of its degree. Key 6 is then the mode of the sixth degree of
// the parent scale, and the hand learns the degrees while the ear finds
// the colour. That is the skill aimed at: finding the parent scale of a
// mode at once, which starts with knowing the degrees by heart. Close
// distractors come with it, lydian always beside ionian, with no rule
// to pick them.
type modes struct {
	system harmony.System
	all    bool
}

var _ Activity = modes{}

func (a modes) ID() string {
	if a.all {
		return "modes-all"
	}
	return "modes"
}

func (a modes) notion(d harmony.Degree) dex.Notion {
	return dex.ModeOf(a.system, d)
}

func (a modes) pattern(n dex.Notion) harmony.ScalePattern {
	p, _ := a.system.Mode(n.Degree)
	return p
}

// Plan makes the seven modes all come once before any comes twice, so
// that ten questions cover the whole system.
func (a modes) Plan(rng *rand.Rand, n int) []Question {
	var order []harmony.Degree
	for len(order) < n {
		for _, i := range rng.Perm(7) {
			order = append(order, harmony.Degree(i+1))
		}
	}
	questions := make([]Question, n)
	for i, d := range order[:n] {
		questions[i] = a.question(d, harmony.PitchClass(rng.IntN(12)), rng)
	}
	return questions
}

func (a modes) question(mode harmony.Degree, tonic harmony.PitchClass, rng *rand.Rand) Question {
	if a.all {
		q := Question{Tonic: tonic, Answer: int(mode) - 1}
		for d := harmony.Degree(1); d <= 7; d++ {
			q.Choices = append(q.Choices, a.notion(d))
		}
		return q
	}

	var far []harmony.Degree
	for d := harmony.Degree(1); d <= 7; d++ {
		if d != mode && a.distance(mode, d) >= farEnough {
			far = append(far, d)
		}
	}
	rng.Shuffle(len(far), func(i, j int) { far[i], far[j] = far[j], far[i] })

	degrees := []harmony.Degree{mode, far[0], far[1]}
	rng.Shuffle(3, func(i, j int) { degrees[i], degrees[j] = degrees[j], degrees[i] })

	q := Question{Tonic: tonic}
	for i, d := range degrees {
		q.Choices = append(q.Choices, a.notion(d))
		if d == mode {
			q.Answer = i
		}
	}
	return q
}

// distance counts the notes two modes of the system differ by on the
// same tonic.
func (a modes) distance(x, y harmony.Degree) int {
	px, _ := a.system.Mode(x)
	py, _ := a.system.Mode(y)
	return px.At(0).Difference(py.At(0)).Len()
}

// Retry asks the same mode on another tonic: the ear has to find the
// colour, not remember the notes.
func (a modes) Retry(q Question, rng *rand.Rand) Question {
	return a.question(q.Right().Degree, otherTonic(q.Tonic, rng), rng)
}

func (a modes) Prompt(l language, q Question) string {
	return fmt.Sprintf(l.words.which, l.note(q.Tonic))
}

func (a modes) Sound(q Question) []keyboard.Note {
	notes, _ := scaleNotes(q.Tonic, a.pattern(q.Right()), 0)
	return notes
}

func (a modes) Correction(q Question, chosen int) []keyboard.Note {
	right, end := scaleNotes(q.Tonic, a.pattern(q.Right()), 0)
	wrong, _ := scaleNotes(q.Tonic, a.pattern(q.Choices[chosen]), end+gap)
	return append(right, wrong...)
}

func (a modes) Show(q Question, chosen int) display {
	return display{
		tonic:   q.Tonic,
		right:   a.pattern(q.Right()),
		chosen:  a.pattern(q.Choices[chosen]),
		mistake: chosen != q.Answer,
	}
}

// Facts are a recognition and a meeting: the game names the mode as it
// confirms.
func (a modes) Facts(q Question) []dex.Fact {
	return []dex.Fact{
		{Kind: dex.FactNamed, Notion: q.Right(), Tonic: q.Tonic},
		{Kind: dex.FactHeard, Notion: q.Right(), Tonic: q.Tonic},
	}
}
