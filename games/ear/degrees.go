package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// degrees is functional ear training: a pedal on the tonic in the same
// octave, one note of the scale beside it, and the player answers its
// degree.
//
// It is the exercise of the beginner at the keyboard: the left hand
// holds the tonic, the right index falls at random in the octave above,
// and the player sings the scale up to that note to find its number.
//
// # What it teaches the dex
//
// A degree is not a notion: it is the interval from the tonic that the
// dex keeps, an elementary notion. A right answer names that interval
// and marks it heard.
//
// # The tonic stays
//
// One tonic for the whole series, drawn at random: the tonal memory is
// built over the series, the pedal recalling the tonic at each
// question. A moving tonic is the level above.
//
// # The correction walks home
//
// After the answer, right or wrong, the scale walks from the note to
// the nearest tonic, as the singing does: up from the tonic for a note
// of the lower tetrachord, up to the octave for one of the upper. The
// pedal goes down to the bass for it: in the same octave, a walk that
// starts on the tonic would start on a unison with it.
type degrees struct {
	scale harmony.ScalePattern
	asked []harmony.Degree // the degrees a question may land on; the answers are always all seven
}

var _ Activity = degrees{}

func (a degrees) ID() string {
	return "degrees"
}

// interval is the notion behind degree `d`: the interval from the
// tonic up to it.
func (a degrees) interval(d harmony.Degree) dex.Notion {
	n, _ := a.scale.Offset(d)
	return dex.IntervalOf(harmony.Interval{Degrees: harmony.Degrees(d - 1), Semitones: n})
}

// Plan draws one tonic for the series, and makes every asked degree
// come once before any comes twice.
func (a degrees) Plan(rng *rand.Rand, n int) []Question {
	tonic := harmony.PitchClass(rng.IntN(12))
	var order []harmony.Degree
	for len(order) < n {
		for _, i := range rng.Perm(len(a.asked)) {
			order = append(order, a.asked[i])
		}
	}
	questions := make([]Question, n)
	for i, d := range order[:n] {
		questions[i] = a.question(d, tonic)
	}
	return questions
}

// question offers the seven degrees in their order, whatever is asked:
// key 3 is always the third degree.
func (a degrees) question(d harmony.Degree, tonic harmony.PitchClass) Question {
	q := Question{Tonic: tonic, Answer: int(d) - 1}
	for c := harmony.Degree(1); c <= 7; c++ {
		q.Choices = append(q.Choices, a.interval(c))
	}
	return q
}

// Retry asks the same degree in the same key: the tonic stays for the
// series.
func (a degrees) Retry(q Question, rng *rand.Rand) Question {
	return a.question(a.degree(q.Answer), q.Tonic)
}

func (a degrees) degree(choice int) harmony.Degree {
	return harmony.Degree(choice + 1)
}

// Prompt names the scale rather than the tonic: the degree is read in
// a scale, and the player sings that one.
func (a degrees) Prompt(l language, q Question) string {
	name, ok := l.namer.ScaleName(q.Tonic, a.scale)
	if !ok {
		name = l.note(q.Tonic)
	}
	return fmt.Sprintf(l.words.whichDegree, name)
}

// Sound is the pedal and the one note, in the octave above the tonic.
func (a degrees) Sound(q Question) []keyboard.Note {
	n, _ := a.scale.Offset(a.degree(q.Answer))
	notes, _ := pathNotes(q.Tonic, []int{int(n)}, pedalBeside, 0)
	return notes
}

// path is the walk from degree `d` to the nearest tonic, in the order
// it is played: up from the tonic to a note of the lower tetrachord,
// up from a note of the upper one to the octave.
func (a degrees) path(d harmony.Degree) []int {
	var offsets []int
	for degree, n := range a.scale.Offsets() {
		if (d <= 4 && degree <= d) || (d > 4 && degree >= d) {
			offsets = append(offsets, int(n))
		}
	}
	if d > 4 {
		offsets = append(offsets, 12)
	}
	return offsets
}

// Correction walks home from the right note, and after a mistake from
// the chosen one too.
func (a degrees) Correction(q Question, chosen int) ([]keyboard.Note, []cue) {
	right, end := pathNotes(q.Tonic, a.path(a.degree(q.Answer)), pedalBelow, 0)
	if chosen == q.Answer {
		return right, []cue{{0, q.Answer}}
	}
	wrong, _ := pathNotes(q.Tonic, a.path(a.degree(chosen)), pedalBelow, end+gap)
	return append(right, wrong...), []cue{{0, q.Answer}, {end + gap, chosen}}
}

// Show lights the walks the correction plays.
func (a degrees) Show(q Question, chosen int) display {
	return display{
		tonic:   q.Tonic,
		right:   a.pathPattern(a.degree(q.Answer)),
		chosen:  a.pathPattern(a.degree(chosen)),
		key:     a.scale,
		mistake: chosen != q.Answer,
	}
}

// pathPattern is the walk as a set on the tonic, the octave folded
// onto it.
func (a degrees) pathPattern(d harmony.Degree) harmony.ScalePattern {
	classes := []harmony.PitchClass{0}
	for _, n := range a.path(d) {
		classes = append(classes, harmony.PitchClass(n%12))
	}
	set, err := harmony.NewPitchSet(classes...)
	if err != nil {
		return 0
	}
	p, _ := harmony.NewScalePattern(set, 0)
	return p
}

func (a degrees) Facts(q Question) []dex.Fact {
	return []dex.Fact{
		{Kind: dex.FactNamed, Notion: q.Right(), Tonic: q.Tonic},
		{Kind: dex.FactHeard, Notion: q.Right(), Tonic: q.Tonic},
	}
}
