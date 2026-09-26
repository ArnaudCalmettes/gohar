package main

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

var major = degrees{scale: harmony.ScaleMajor, asked: []harmony.Degree{1, 2, 3, 4, 5, 6, 7}}

// One tonic for the series, seven answers in the order of the degrees.
func TestDegreesKeepTheTonic(t *testing.T) {
	questions := major.Plan(rand.New(rand.NewPCG(1, 2)), 10)
	for i, q := range questions {
		if q.Tonic != questions[0].Tonic {
			t.Errorf("question %d moves the tonic", i)
		}
		if len(q.Choices) != 7 || q.Choices[2].Interval != harmony.IntMajorThird {
			t.Errorf("question %d: %v, want the seven degrees in order", i, q.Choices)
		}
	}
}

// Only the asked degrees come up, the answers staying seven.
func TestDegreesAskOnlyWhatTheLevelAsks(t *testing.T) {
	few := degrees{scale: harmony.ScaleMajor, asked: []harmony.Degree{1, 3, 5}}
	for i, q := range few.Plan(rand.New(rand.NewPCG(1, 2)), 10) {
		if d := few.degree(q.Answer); d != 1 && d != 3 && d != 5 {
			t.Errorf("question %d asks degree %d", i, d)
		}
	}
}

// The walk home: up from the tonic in the lower tetrachord, up to the
// octave in the upper one.
func TestDegreePath(t *testing.T) {
	for _, c := range []struct {
		d    harmony.Degree
		want []int
	}{
		{1, []int{0}},
		{3, []int{0, 2, 4}},
		{4, []int{0, 2, 4, 5}},
		{5, []int{7, 9, 11, 12}},
		{7, []int{11, 12}},
	} {
		if got := major.path(c.d); !slices.Equal(got, c.want) {
			t.Errorf("degree %d: %v, want %v", c.d, got, c.want)
		}
	}
}
