package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

var formulaNames = [formulas]string{"", "anatole", "III-VI-II-V-I", "aeolian", "II-V-I"}

// formulasRead writes the formulas of a grid of the game as a chart
// reads: each chord, then the formula it is part of, between brackets.
func formulasRead(t *testing.T, name string) string {
	t.Helper()
	tune, err := readTune(name)
	if err != nil {
		t.Fatal(err)
	}
	f := formulasOf(tune.grid)
	var out []string
	for i, ch := range tune.grid.Chords {
		s := symbolOf(tune.written[i]).String()
		if k := f[ch.Start]; k != noFormula {
			s += "[" + formulaNames[k] + "]"
		}
		out = append(out, s)
	}
	return strings.Join(out, " ")
}

// The blues ends on two anatoles, F6 D7 Gm7 C7 then F6 D7 G7 C7; Tune Up
// is three II-V-I, the first played again.
func TestFormulas(t *testing.T) {
	for name, want := range map[string]string{
		jazzBlues:     "F7 B♭7 F7 B♭7 F6[anatole] D7[anatole] Gm7[anatole] C7[anatole] F6[anatole] D7[anatole] G7[anatole] C7[anatole]",
		"tune-up.cho": "Em7[II-V-I] A7[II-V-I] Dmaj7[II-V-I] Dm7[II-V-I] G7[II-V-I] Cmaj7[II-V-I] Cm7[II-V-I] F7[II-V-I] B♭maj7[II-V-I] Gm7 Em7[II-V-I] A7[II-V-I] Dmaj7[II-V-I]",
	} {
		if got := formulasRead(t, name); got != want {
			t.Errorf("%s:\n got %s\nwant %s", name, got, want)
		}
	}
}

// run plays two choruses of the blues at 120, every arrival landed but
// the beats `missed` (bar and beat, from 1, in both choruses), every
// note off the beat by `off`, give or take `spread`.
func run(t *testing.T, off, spread time.Duration, missed ...[2]int) *summary {
	t.Helper()
	grid, err := readGrid(jazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	beats := Expect(grid, tempo.NewMetronome(time.Time{}, 120, perBar), 2)
	chorus := len(beats) / 2
	miss := map[int]bool{}
	for _, m := range missed {
		n := (m[0]-1)*perBar + m[1] - 1
		miss[n], miss[chorus+n] = true, true
	}
	s := newSummary(12, chorus, formulasOf(grid))
	for n, b := range beats {
		k := Landed
		if miss[n] {
			k = Missed
		}
		s.beat(n, b, k)
		d := off + spread
		if n%2 == 1 {
			d = off - spread
		}
		s.note(NoteMark{Off: d, Timing: OnTime})
	}
	return s
}

func said(a advice) string {
	switch a.kind {
	case workFormula:
		return "work the " + formulaNames[a.what] + " at this tempo"
	case workSit:
		return fmt.Sprintf("work situation %d at this tempo", a.sit)
	case slower:
		return fmt.Sprintf("slow down to %.0f", a.bpm)
	case faster:
		return fmt.Sprintf("speed up to %.0f", a.bpm)
	case stay:
		return "stay at this tempo"
	}
	return "nothing"
}

func TestAdvice(t *testing.T) {
	d7, gm7, g7 := [2]int{8, 1}, [2]int{9, 1}, [2]int{12, 1}
	for _, tc := range []struct {
		name        string
		off, spread time.Duration
		bpm         float64
		missed      [][2]int
		feel        feel
		want        string
	}{
		{"D7, Gm7 and G7 missed, in time", 2 * time.Millisecond, 20 * time.Millisecond, 120, [][2]int{d7, gm7, g7}, feelSteady, "work the anatole at this tempo"},
		{"the same, rushing", -40 * time.Millisecond, 10 * time.Millisecond, 120, [][2]int{d7, gm7, g7}, feelRushing, "slow down to 110"},
		{"the same, scattered, fast", 0, 50 * time.Millisecond, 160, [][2]int{d7, gm7, g7}, feelUnsteady, "slow down to 155"},
		{"everything landed, dragging", 30 * time.Millisecond, 0, 120, nil, feelDragging, "stay at this tempo"},
		{"everything landed, in time", 0, 10 * time.Millisecond, 120, nil, feelSteady, "speed up to 130"},
		{"everything landed, in time, at 140", 0, 10 * time.Millisecond, 140, nil, feelSteady, "speed up to 145"},
	} {
		s := run(t, tc.off, tc.spread, tc.missed...)
		if got := s.feel(); got != tc.feel {
			t.Errorf("%s: feel %d, want %d", tc.name, got, tc.feel)
		}
		if got := said(advise(s, tc.bpm)); got != tc.want {
			t.Errorf("%s: %s, want %s", tc.name, got, tc.want)
		}
	}
}
