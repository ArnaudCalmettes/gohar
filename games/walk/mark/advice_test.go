package mark

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/games/walk/grids"
)

var formulaNames = [Formulas]string{"", "anatole", "III-VI-II-V-I", "aeolian", "II-V-I"}

// formulasRead writes the formulas of a grid of the game as a chart
// reads: each chord, then the formula it is part of, between brackets.
func formulasRead(t *testing.T, name string) string {
	t.Helper()
	tune, err := grids.Read(name)
	if err != nil {
		t.Fatal(err)
	}
	f := FormulasOf(tune.Grid)
	var out []string
	for i, ch := range tune.Grid.Chords {
		s := tune.Written[i].String()
		if len(s) > 1 && s[1] == 'b' { // a flat root, as a musician writes it
			s = s[:1] + "♭" + s[2:]
		}
		if k := f[ch.Start]; k != NoFormula {
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
		grids.JazzBlues: "F7 B♭7 F7 B♭7 F6[anatole] D7[anatole] Gm7[anatole] C7[anatole] F6[anatole] D7[anatole] G7[anatole] C7[anatole]",
		"tune-up.cho":   "Em7[II-V-I] A7[II-V-I] Dmaj7[II-V-I] Dm7[II-V-I] G7[II-V-I] Cmaj7[II-V-I] Cm7[II-V-I] F7[II-V-I] B♭maj7[II-V-I] Gm7 Em7[II-V-I] A7[II-V-I] Dmaj7[II-V-I]",
	} {
		if got := formulasRead(t, name); got != want {
			t.Errorf("%s:\n got %s\nwant %s", name, got, want)
		}
	}
}

// run plays two choruses of the blues at 120, every arrival landed but
// the beats `missed` (bar and beat, from 1, in both choruses), every
// note off the beat by `off`, give or take `spread`.
func run(t *testing.T, off, spread time.Duration, missed ...[2]int) *Summary {
	t.Helper()
	grid, err := grids.Changes(grids.JazzBlues)
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
	s := NewSummary(12, chorus, FormulasOf(grid))
	for n, b := range beats {
		k := Landed
		if miss[n] {
			k = Missed
		}
		s.Beat(n, b, k)
		d := off + spread
		if n%2 == 1 {
			d = off - spread
		}
		s.Note(NoteMark{Off: d, Timing: OnTime})
	}
	return s
}

func said(a Advice) string {
	switch a.Kind {
	case WorkFormula:
		return "work the " + formulaNames[a.What] + " at this tempo"
	case WorkSituation:
		return fmt.Sprintf("work situation %d at this tempo", a.Situation)
	case Slower:
		return fmt.Sprintf("slow down to %.0f", a.BPM)
	case Faster:
		return fmt.Sprintf("speed up to %.0f", a.BPM)
	case Stay:
		return "stay at this tempo"
	}
	return "nothing"
}

// The bounds of the options' tempo, as the game sets them.
const slowest, fastest = 60, 240

func TestAdvice(t *testing.T) {
	d7, gm7, g7 := [2]int{8, 1}, [2]int{9, 1}, [2]int{12, 1}
	for _, tc := range []struct {
		name        string
		off, spread time.Duration
		BPM         float64
		missed      [][2]int
		Feel        Feel
		want        string
	}{
		{"D7, Gm7 and G7 missed, in time", 2 * time.Millisecond, 20 * time.Millisecond, 120, [][2]int{d7, gm7, g7}, FeelSteady, "work the anatole at this tempo"},
		{"the same, rushing", -40 * time.Millisecond, 10 * time.Millisecond, 120, [][2]int{d7, gm7, g7}, FeelRushing, "slow down to 110"},
		{"the same, scattered, fast", 0, 50 * time.Millisecond, 160, [][2]int{d7, gm7, g7}, FeelUnsteady, "slow down to 155"},
		{"everything landed, dragging", 30 * time.Millisecond, 0, 120, nil, FeelDragging, "stay at this tempo"},
		{"everything landed, in time", 0, 10 * time.Millisecond, 120, nil, FeelSteady, "speed up to 130"},
		{"everything landed, in time, at 140", 0, 10 * time.Millisecond, 140, nil, FeelSteady, "speed up to 145"},
	} {
		s := run(t, tc.off, tc.spread, tc.missed...)
		if got := s.Feel(); got != tc.Feel {
			t.Errorf("%s: feel %d, want %d", tc.name, got, tc.Feel)
		}
		if got := said(Advise(s, tc.BPM, slowest, fastest)); got != tc.want {
			t.Errorf("%s: %s, want %s", tc.name, got, tc.want)
		}
	}
}
