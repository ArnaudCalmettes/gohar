package bass

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/games/walk/grids"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// walkBlues generates two choruses of the jazz blues in F, at 120, for
// each of a handful of seeds: the rules draw among their paths, and
// every path must hold.
func walkBlues(t *testing.T, test func(seed uint64, beats []mark.Beat, line []int)) {
	t.Helper()
	grid, err := grids.Changes(grids.JazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	beats := mark.Expect(grid, at120(), 2)
	for seed := range uint64(20) {
		test(seed, beats, Walk(beats, rand.New(rand.NewPCG(seed, 1))))
	}
}

// The line the marker hears, note for note on the beat, lands every
// arrival and every held bar: the generator and the marker share their
// rules.
func TestWalkLandsEverything(t *testing.T) {
	walkBlues(t, func(seed uint64, beats []mark.Beat, line []int) {
		m := at120()
		k := mark.NewMarker(mark.FirstPalier, m, beats)
		for n, key := range line {
			k.Play(mark.Note{Key: key, At: m.At(n)})
		}
		for _, bm := range k.Close(m.At(len(beats))) {
			if bm.Kind != mark.Landed {
				p := m.Position(bm.Beat)
				t.Errorf("seed %d: bar %d, beat %d marked %d, want landed", seed, p.Bar, p.Beat, bm.Kind)
			}
		}
	})
}

// From E1 to C3, an octave on the root down to C1 aside, and never the
// same key twice in a row (Siskind, p. 116).
func TestWalkRegisterAndRepeats(t *testing.T) {
	walkBlues(t, func(seed uint64, beats []mark.Beat, line []int) {
		for n, key := range line {
			eighth := n > 0 && abs(key-line[n-1]) == octave
			if !eighth && key < lowE1 || key < lowC1 || key > highC3 {
				t.Errorf("seed %d, beat %d: key %d outside the register", seed, n, key)
			}
			if n > 0 && key == line[n-1] {
				t.Errorf("seed %d, beat %d: key %d repeated", seed, n, key)
			}
		}
	})
}

// A fifth at most from one note to the next, or an octave from a root
// to the same root: no wide leap in the middle of a line.
func TestWalkLeaps(t *testing.T) {
	walkBlues(t, func(seed uint64, beats []mark.Beat, line []int) {
		for n := 1; n < len(line); n++ {
			leap := abs(line[n] - line[n-1])
			root := harmony.PitchClass(line[n]%12) == beats[n].Chord.Bass
			if leap > maxLeap && !(leap == octave && root) {
				p := at120().Position(n)
				t.Errorf("seed %d: bar %d, beat %d leaps %d semitones", seed, p.Bar, p.Beat, leap)
			}
		}
	})
}

// After a walk up, F7 followed by F7 carries on climbing, F G A B♭,
// two times in three rather than taking another path. B♭7 followed by
// B♭7, after a walk up to A2, never climbs: B♭ C D E♭ would pass C3.
func TestWalkCarriesOnTheRun(t *testing.T) {
	grid, err := grids.Changes(grids.JazzBlues)
	if err != nil {
		t.Fatal(err)
	}
	beats := mark.Expect(grid, at120(), 1)
	f7, bb7 := beats[:mark.BeatsPerBar], beats[mark.BeatsPerBar:2*mark.BeatsPerBar]
	climbs := 0
	const draws = 300
	for seed := range uint64(draws) {
		rng := rand.New(rand.NewPCG(seed, 1))
		if p, _ := choose(paths(f7, f7[0].Chord.Bass), f7, 28, +1, nil, rng); p.dir == +1 { // from E1
			climbs++
		}
		if p, keys := choose(paths(bb7, bb7[0].Chord.Bass), bb7, 45, +1, nil, rng); p.dir == +1 { // from A2
			t.Errorf("seed %d: after a walk up to A2, %v climbs past C3", seed, keys)
		}
	}
	if climbs < draws/2 || climbs > draws*5/6 {
		t.Errorf("after a walk up from E1, %d climbs in %d draws, want about two in three", climbs, draws)
	}
}

// A bar over a root seldom replays the bar played over that root
// before: only when the register leaves nothing else, a C at the top
// that can only walk down. A demo does not replay one walk down.
func TestWalkVaries(t *testing.T) {
	bars, replays := 0, 0
	walkBlues(t, func(seed uint64, beats []mark.Beat, line []int) {
		last := map[harmony.PitchClass][]int{}
		for bar := 0; bar+mark.BeatsPerBar <= len(line); bar += mark.BeatsPerBar {
			root, keys := beats[bar].Chord.Bass, line[bar:bar+mark.BeatsPerBar]
			if slices.Equal(keys, last[root]) {
				replays++
			}
			last[root] = keys
			bars++
		}
	})
	if replays*50 > bars {
		t.Errorf("%d bars replayed in %d, want under one in fifty", replays, bars)
	}
}

// Bar 9, Gm7 then C7 a fifth below: the fourth beat leads into C, a
// half step or a whole step away, whichever path was drawn.
func TestWalkLeadsIntoTheNextRoot(t *testing.T) {
	walkBlues(t, func(seed uint64, beats []mark.Beat, line []int) {
		c := harmony.PitchClass(0)
		lead := harmony.PitchClass(line[35] % 12) // bar 9, beat 4
		if d := min(lead.Up(c), c.Up(lead)); d < 1 || d > 2 {
			t.Errorf("seed %d: bar 9 ends on %s, %d semitones from C", seed, noteNames[lead], d)
		}
	})
}

// The blues loops and says nothing of its end: after the turnaround,
// G7 C7, the line lands on F, a fifth away at most. A chart that says
// where it ends stops there.
func TestWalkEndsOnTheTonic(t *testing.T) {
	walkBlues(t, func(seed uint64, beats []mark.Beat, line []int) {
		k := Ending(beats, line, 0)
		if k%12 != 5 || abs(k-line[len(line)-1]) > maxLeap {
			t.Errorf("seed %d: after %d, the line ends on %d, want an F nearby", seed, line[len(line)-1], k)
		}
		if k := Ending(beats, line, 3); k != 0 {
			t.Errorf("seed %d: a chart with an end gets a last note, %d", seed, k)
		}
	})
}

// At 120, a beat lasts half a second, as in mark's tests.
func at120() tempo.Metronome { return tempo.NewMetronome(time.Time{}, 120, mark.BeatsPerBar) }

var noteNames = [12]string{"C", "D♭", "D", "E♭", "E", "F", "F♯", "G", "A♭", "A", "B♭", "B"}

// On every grid as on the blues, the reference line lands every
// arrival: two chords a bar, II-V-I in major and in minor, II-V a half
// tone apart.
func TestWalkLandsEveryGrid(t *testing.T) {
	tunes, err := grids.All()
	if err != nil {
		t.Fatal(err)
	}
	for _, tu := range tunes[1:] {
		beats := mark.Expect(tu.Grid, at120(), 2)
		for seed := range uint64(20) {
			line := Walk(beats, rand.New(rand.NewPCG(seed, 1)))
			m := at120()
			k := mark.NewMarker(mark.FirstPalier, m, beats)
			for n, key := range line {
				k.Play(mark.Note{Key: key, At: m.At(n)})
			}
			for _, bm := range k.Close(m.At(len(beats))) {
				if bm.Kind != mark.Landed {
					p := m.Position(bm.Beat)
					t.Errorf("%s, seed %d: bar %d, beat %d marked %d, want landed", tu.Title, seed, p.Bar, p.Beat, bm.Kind)
				}
			}
		}
	}
}
