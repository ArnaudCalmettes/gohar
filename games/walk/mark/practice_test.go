package mark

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/games/walk/grids"
)

// Without tempo, the chart waits on F7 until an F is played, then on
// B♭7, and loops after the turnaround.
func TestPractice(t *testing.T) {
	grid, err := grids.Changes(grids.JazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	p := NewPractice(FirstPalier, Expect(grid, at120(), 1))
	waiting := func() string { return chordName(p.Waiting().Chord) }

	if got := p.Play(A2); got != ChordTone || waiting() != "F7" {
		t.Errorf("A on F7: %d, waiting on %s; want a chord tone, still F7", got, waiting())
	}
	if got := p.Play(Fs2); got != Outside || waiting() != "F7" {
		t.Errorf("F♯ on F7: %d, waiting on %s; want outside, still F7", got, waiting())
	}
	if got := p.Play(F2 + 12); got != Root || waiting() != "B♭7" {
		t.Errorf("F, an octave up, on F7: %d, waiting on %s; want the root, then B♭7", got, waiting())
	}
	if !p.landed[0] {
		t.Error("bar 1, beat 1 not marked landed")
	}

	// Through the chorus, root after root, back to the top.
	for range len(p.arrivals) - 1 {
		p.Play(C2 + int(p.Waiting().Chord.Bass))
	}
	if waiting() != "F7" || p.at != 0 || len(p.landed) != 0 {
		t.Errorf("after the turnaround: waiting on %s, %d landed; want F7 again, marks cleared", waiting(), len(p.landed))
	}
}
