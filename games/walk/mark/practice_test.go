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
	if waiting() != "F7" || p.at != 0 || p.Choruses() != 1 || !p.Landed(0) {
		t.Errorf("after the turnaround: waiting on %s, %d choruses, bar 1 landed %v; want F7 again, one chorus, the marks kept", waiting(), p.Choruses(), p.Landed(0))
	}
	p.Play(Fs2)
	if len(p.landed) != 0 {
		t.Errorf("a note on the next chorus: %d landed; want the marks cleared", len(p.landed))
	}
}

// At palier 0, only the root lands, even on a bar where the chord
// carries on: the F7 of bar 4 waits for an F.
func TestPalier0(t *testing.T) {
	grid, err := grids.Changes(grids.JazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	p := NewPractice(Palier0, Expect(grid, at120(), 1))
	p.Play(F2)
	p.Play(Bb1)
	p.Play(F2)
	if got := chordName(p.Waiting().Chord); got != "F7" || p.Waiting().N != 12 {
		t.Fatalf("after F, B♭, F: waiting on %s, beat %d; want F7, bar 4", got, p.Waiting().N)
	}
	if got := p.Play(A2); got != ChordTone || p.Waiting().N != 12 {
		t.Errorf("A on the F7 of bar 4: %d, waiting on beat %d; want a chord tone, still bar 4", got, p.Waiting().N)
	}
	if p.Play(F2); p.Waiting().N != 16 {
		t.Errorf("F on the F7 of bar 4: waiting on beat %d; want bar 5", p.Waiting().N)
	}
}
