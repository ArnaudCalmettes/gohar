package main

import (
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// The expectations read as a chart: chord names, bars and beats.

var t0 = time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)

// At 120, a beat lasts half a second: easy to check by hand.
func at120() tempo.Metronome { return tempo.NewMetronome(t0, 120, 4) }

var noteNames = [12]string{"C", "D♭", "D", "E♭", "E", "F", "F♯", "G", "A♭", "A", "B♭", "B"}

func chordName(ch analysis.Change) string {
	switch ch.Chord.Pattern {
	case harmony.ChordDominantSeventh:
		return noteNames[ch.Chord.Root] + "7"
	case harmony.ChordMinorSeventh:
		return noteNames[ch.Chord.Root] + "m7"
	case harmony.ChordMajorSeventh:
		return noteNames[ch.Chord.Root] + "maj7"
	case harmony.ChordMajorSixth:
		return noteNames[ch.Chord.Root] + "6"
	}
	return noteNames[ch.Chord.Root] + "?"
}

// beatAt finds bar `bar`, beat `beat` of chorus 1.
func beatAt(t *testing.T, beats []Beat, bar, beat int) Beat {
	t.Helper()
	for _, b := range beats {
		if b.Position == (tempo.Position{Bar: bar, Beat: beat}) {
			return b
		}
	}
	t.Fatalf("no bar %d, beat %d", bar, beat)
	return Beat{}
}

// The jazz blues of the first milestone, in F:
// F7 | B♭7 | F7 | F7 | B♭7 | B♭7 | F6 | D7 | Gm7 | C7 | F6 D7 | G7 C7.
func TestExpectJazzBlues(t *testing.T) {
	grid, err := readGrid(jazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	beats := Expect(grid, at120(), 2)
	if len(beats) != 2*48 {
		t.Fatalf("%d beats, want two choruses of twelve bars", len(beats))
	}
	for _, tc := range []struct {
		bar, beat int
		chord     string
		arrives   bool
		strong    bool
		next      string
		perBar    int
	}{
		{1, 1, "F7", true, true, "B♭7", 1},
		{1, 2, "F7", false, false, "B♭7", 1},
		{2, 1, "B♭7", true, true, "F7", 1}, // the IV in bar 2
		{8, 1, "D7", true, true, "Gm7", 1}, // the VI7
		{9, 1, "Gm7", true, true, "C7", 1},
		{11, 1, "F6", true, true, "D7", 2},
		{11, 3, "D7", true, true, "G7", 2}, // two chords a bar
		{12, 3, "C7", true, true, "F7", 2},
		{12, 4, "C7", false, false, "F7", 2}, // the turnaround leads back
		{13, 1, "F7", true, true, "B♭7", 1},  // the second chorus
	} {
		b := beatAt(t, beats, tc.bar, tc.beat)
		if got := chordName(b.Chord); got != tc.chord {
			t.Errorf("bar %d, beat %d: %s, want %s", tc.bar, tc.beat, got, tc.chord)
		}
		if b.Arrives != tc.arrives {
			t.Errorf("bar %d, beat %d: arrives %v, want %v", tc.bar, tc.beat, b.Arrives, tc.arrives)
		}
		if b.Strong != tc.strong {
			t.Errorf("bar %d, beat %d: strong %v, want %v", tc.bar, tc.beat, b.Strong, tc.strong)
		}
		if got := chordName(b.Next); got != tc.next {
			t.Errorf("bar %d, beat %d: next %s, want %s", tc.bar, tc.beat, got, tc.next)
		}
		if b.PerBar != tc.perBar {
			t.Errorf("bar %d, beat %d: %d changes in the bar, want %d", tc.bar, tc.beat, b.PerBar, tc.perBar)
		}
	}
}

// A II-V in one bar, then the I for a bar: two chords a bar, the second
// arriving on beat 3, which Siskind's second rule is about.
func TestExpectTwoChordsABar(t *testing.T) {
	d, g, c := harmony.PitchClass(2), harmony.PitchClass(7), harmony.PitchClass(0)
	half, bar := 2*analysis.TicksPerBeat, 4*analysis.TicksPerBeat
	changes := analysis.Changes{
		Bars: []analysis.Ticks{0, bar},
		Chords: []analysis.Change{
			{Chord: harmony.Chord{Root: d, Pattern: harmony.ChordMinorSeventh}, Bass: d, Start: 0, Length: half},
			{Chord: harmony.Chord{Root: g, Pattern: harmony.ChordDominantSeventh}, Bass: g, Start: half, Length: half},
			{Chord: harmony.Chord{Root: c, Pattern: harmony.ChordMajorSeventh}, Bass: c, Start: bar, Length: bar},
		},
	}
	beats := Expect(changes, at120(), 1)
	for _, tc := range []struct {
		bar, beat int
		chord     string
		arrives   bool
		perBar    int
	}{
		{1, 1, "Dm7", true, 2},
		{1, 2, "Dm7", false, 2},
		{1, 3, "G7", true, 2},
		{1, 4, "G7", false, 2},
		{2, 1, "Cmaj7", true, 1},
	} {
		b := beatAt(t, beats, tc.bar, tc.beat)
		if got := chordName(b.Chord); got != tc.chord || b.Arrives != tc.arrives || b.PerBar != tc.perBar {
			t.Errorf("bar %d, beat %d: %s arrives %v, %d a bar; want %s arrives %v, %d a bar",
				tc.bar, tc.beat, got, b.Arrives, b.PerBar, tc.chord, tc.arrives, tc.perBar)
		}
	}
	if last := beatAt(t, beats, 2, 4); chordName(last.Next) != "Cmaj7" {
		t.Errorf("at the end of a chart that does not loop, next is %s, want the chord itself", chordName(last.Next))
	}
}

// Bar 4 carries on the F7 of bar 3, as two tied whole notes would: its
// beat 1 holds, and only its beat 1.
func TestExpectHolds(t *testing.T) {
	grid, err := readGrid(jazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	beats := Expect(grid, at120(), 1)
	for _, tc := range []struct {
		bar, beat int
		holds     bool
	}{
		{1, 1, false}, // F7 arrives
		{4, 1, true},  // F7 carries on from bar 3
		{4, 2, false},
		{6, 1, true}, // B♭7 carries on from bar 5
		{11, 3, false},
	} {
		if b := beatAt(t, beats, tc.bar, tc.beat); b.Holds != tc.holds {
			t.Errorf("bar %d, beat %d: holds %v, want %v", tc.bar, tc.beat, b.Holds, tc.holds)
		}
	}
}

// On the blues, beat 3 of bar 11 in the second chorus: bar 11, chorus
// 2, counted from 1 as a musician counts them.
func TestBarOf(t *testing.T) {
	chorus := 12 * perBar
	bar, n := barOf(chorus+10*perBar+2, chorus)
	if bar+1 != 11 || n+1 != 2 {
		t.Errorf("bar %d of chorus %d, want bar 11 of chorus 2", bar+1, n+1)
	}
}
