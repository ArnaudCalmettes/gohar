package main

import (
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// Keys in the register of the double bass, named as a bassist reads
// them.
const (
	F2  = 41
	Fs2 = 42 // F♯2
	A2  = 45
	C2  = 36
	Bb1 = 34 // B♭1
)

// runOf marks the jazz blues in F at 120: a beat is 500 ms, so the
// first palier is on time within 62 ms and loose within 166 ms.
func runOf(t *testing.T, r Rules) (*Marker, tempo.Metronome) {
	t.Helper()
	grid, err := readGrid(jazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	m := at120()
	return NewMarker(r, m, Expect(grid, m, 1)), m
}

// at is bar `bar`, beat `beat`, moved by `ms` milliseconds.
func at(m tempo.Metronome, bar, beat, ms int) time.Time {
	return m.At((bar-1)*4 + beat - 1).Add(time.Duration(ms) * time.Millisecond)
}

func TestMarkerNotes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		key    int
		ms     int
		timing Timing
		pitch  Pitch
	}{
		{"F on bar 1, on time", F2, 10, OnTime, Root},
		{"F 100 ms late", F2, 100, Late, Root},
		{"F 120 ms early", F2, -120, Early, Root},
		{"A, the third of F7", A2, 0, OnTime, ChordTone},
		{"F♯, not in F7", Fs2, 0, OnTime, Outside},
	} {
		k, m := runOf(t, FirstPalier)
		mark, ok := k.Play(Note{tc.key, at(m, 1, 1, tc.ms)})
		if !ok || mark.Beat != 0 || mark.Timing != tc.timing || mark.Pitch != tc.pitch {
			t.Errorf("%s: %+v, want beat 0, timing %d, pitch %d", tc.name, mark, tc.timing, tc.pitch)
		}
	}
}

// What becomes of bar 1, beat 1, where F7 arrives.
func TestMarkerArrival(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rules Rules
		notes []struct{ key, ms int } // ms from bar 1, beat 1
		want  BeatKind
	}{
		{"F, on time", FirstPalier, []struct{ key, ms int }{{F2, 0}}, Landed},
		{"F, 150 ms late", FirstPalier, []struct{ key, ms int }{{F2, 150}}, Landed},
		{"nothing", FirstPalier, nil, Missed},
		{"F, 200 ms late: between two beats", FirstPalier, []struct{ key, ms int }{{F2, 200}}, Missed},
		{"A instead of F, for a beginner", FirstPalier, []struct{ key, ms int }{{A2, 0}}, Missed},
		{"A instead of F, inversions allowed", withInversions(), []struct{ key, ms int }{{A2, 0}}, Landed},
		{"F♯ instead of F", FirstPalier, []struct{ key, ms int }{{Fs2, 0}}, Missed},
		{"F and C together: which one was meant?", FirstPalier, []struct{ key, ms int }{{F2, 0}, {C2, 10}}, Doubled},
		{"F twice", FirstPalier, []struct{ key, ms int }{{F2, -20}, {F2 + 12, 20}}, Doubled},
	} {
		k, m := runOf(t, tc.rules)
		for _, n := range tc.notes {
			k.Play(Note{n.key, at(m, 1, 1, n.ms)})
		}
		marks := k.Close(at(m, 1, 2, 0))
		if len(marks) != 1 || marks[0] != (BeatMark{0, tc.want}) {
			t.Errorf("%s: %+v, want bar 1, beat 1 marked %d", tc.name, marks, tc.want)
		}
	}
}

func withInversions() Rules {
	r := FirstPalier
	r.Inversions = true
	return r
}

// Between two arrivals the player may repeat the root or stay silent:
// no mark, unless two notes claim the same beat.
func TestMarkerBetweenArrivals(t *testing.T) {
	k, m := runOf(t, FirstPalier)
	k.Play(Note{F2, at(m, 1, 1, 0)})
	k.Play(Note{F2, at(m, 1, 2, 0)}) // the root again: fine
	k.Play(Note{A2, at(m, 1, 3, 0)}) // the third, off an arrival: fine
	k.Play(Note{F2, at(m, 1, 4, 0)})
	k.Play(Note{C2, at(m, 1, 4, 30)}) // two notes on beat 4
	marks := k.Close(at(m, 2, 1, 0))
	want := []BeatMark{{0, Landed}, {3, Doubled}}
	if len(marks) != len(want) || marks[0] != want[0] || marks[1] != want[1] {
		t.Errorf("bar 1: %+v, want F7 landed, beat 4 doubled", marks)
	}
}

// Bar 11 holds two chords, F6 then D7 on beat 3: two arrivals.
func TestMarkerTwoArrivalsInABar(t *testing.T) {
	k, m := runOf(t, FirstPalier)
	k.Close(at(m, 11, 1, -200)) // the first ten bars, unplayed
	k.Play(Note{F2, at(m, 11, 1, 0)})
	k.Play(Note{Bb1, at(m, 11, 3, 0)}) // B♭ under D7: wrong root
	marks := k.Close(at(m, 12, 1, 0))
	want := []BeatMark{{40, Landed}, {42, Missed}}
	if len(marks) != len(want) || marks[0] != want[0] || marks[1] != want[1] {
		t.Errorf("bar 11: %+v, want F6 landed, D7 missed", marks)
	}
}

// A beat is only marked once its loose window has closed: a note may
// still come late.
func TestMarkerWaitsForTheWindow(t *testing.T) {
	k, m := runOf(t, FirstPalier)
	if marks := k.Close(at(m, 1, 1, 160)); len(marks) != 0 {
		t.Errorf("160 ms after bar 1, beat 1: %+v, want nothing yet", marks)
	}
	k.Play(Note{F2, at(m, 1, 1, 165)})
	if marks := k.Close(at(m, 1, 1, 170)); len(marks) != 1 || marks[0].Kind != Landed {
		t.Errorf("170 ms after: %+v, want F7 landed", marks)
	}
}

func TestMarkerIgnores(t *testing.T) {
	k, m := runOf(t, FirstPalier)
	if _, ok := k.Play(Note{60, at(m, 1, 1, 0)}); ok {
		t.Error("a middle C, in the right hand's zone, was marked")
	}
	if _, ok := k.Play(Note{F2, at(m, 0, 4, 0)}); ok {
		t.Error("a note in the count-in was marked")
	}
}

// Beat 1 of bar 4, where F7 carries on: a note is expected, any note of
// the tetrad.
func TestMarkerHeldBar(t *testing.T) {
	for _, tc := range []struct {
		name string
		key  int
		want BeatKind
	}{
		{"F, the root", F2, Landed},
		{"C, the fifth", C2, Landed},
		{"A, the third", A2, Landed},
		{"E♭, the seventh", 39, Landed},
		{"F♯, not in F7", Fs2, Missed},
		{"nothing", 0, Missed},
	} {
		k, m := runOf(t, FirstPalier)
		k.Close(at(m, 4, 1, -250)) // the first three bars, unplayed
		if tc.key != 0 {
			k.Play(Note{tc.key, at(m, 4, 1, 0)})
		}
		marks := k.Close(at(m, 4, 2, 0))
		if len(marks) != 1 || marks[0] != (BeatMark{12, tc.want}) {
			t.Errorf("%s on bar 4, beat 1: %+v, want %d", tc.name, marks, tc.want)
		}
	}
}
