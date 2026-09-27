package analysis_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

func TestChangesOrder(t *testing.T) {
	five := make([]analysis.Change, 5)
	for name, tc := range map[string]struct {
		changes    analysis.Changes
		next, prev []int // for each change
	}{
		"a chart loops": {
			analysis.Changes{Chords: five, Loops: true},
			[]int{1, 2, 3, 4, 0}, []int{4, 0, 1, 2, 3},
		},
		"a chart with a coda from the fourth change": {
			analysis.Changes{Chords: five, Loops: true, Coda: 3},
			[]int{1, 2, 0, 4, -1}, []int{2, 0, 1, 2, 3},
		},
		"a performance only grows": {
			analysis.Changes{Chords: five},
			[]int{1, 2, 3, 4, -1}, []int{-1, 0, 1, 2, 3},
		},
	} {
		for i := range five {
			if got := tc.changes.Next(i); got != tc.next[i] {
				t.Errorf("%s: after %d comes %d, want %d", name, i, got, tc.next[i])
			}
			if got := tc.changes.Prev(i); got != tc.prev[i] {
				t.Errorf("%s: before %d comes %d, want %d", name, i, got, tc.prev[i])
			}
		}
	}
}

func TestChangeInverted(t *testing.T) {
	c := harmony.Chord{Root: 0, Pattern: harmony.ChordMajorTriad}
	if (analysis.Change{Chord: c, Bass: 0}).Inverted() {
		t.Error("C over C is inverted")
	}
	if !(analysis.Change{Chord: c, Bass: 4}).Inverted() {
		t.Error("C over E is not inverted")
	}
	if (analysis.Change{Chord: c, Bass: 4, Silent: true}).Inverted() {
		t.Error("a silence is inverted")
	}
}

func TestChangesBar(t *testing.T) {
	b := analysis.TicksPerBeat
	c := analysis.Changes{Bars: []analysis.Ticks{0, 4 * b, 7 * b}}
	for at, want := range map[analysis.Ticks]int{0: 0, 3 * b: 0, 4 * b: 1, 6 * b: 1, 7 * b: 2, 20 * b: 2} {
		if got := c.Bar(at); got != want {
			t.Errorf("tick %d in bar %d, want %d", at, got, want)
		}
	}
	if got := (analysis.Changes{}).Bar(0); got != -1 {
		t.Errorf("no bars, yet bar %d", got)
	}
}
