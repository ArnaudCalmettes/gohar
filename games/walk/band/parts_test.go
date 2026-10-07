package band

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// described says the strokes as a musician would: the sound, where in the
// beat, and how hard when it leans.
func described(Strokes []Stroke) []string {
	names := map[sound]string{hatSound: "hi-hat", rideSound: "ride", snapSound: "snap", bassSound: "bass"}
	var out []string
	for _, s := range Strokes {
		where := "on the beat"
		if s.at != 0 {
			where = "on the swung and"
		}
		out = append(out, fmt.Sprintf("%s %s %.1f", names[s.sound], where, s.vel))
	}
	slices.Sort(out)
	return out
}

func on(bar, beat int) tempo.Position { return tempo.Position{Bar: bar, Beat: beat} }

// The swing ride: "ding, ding-da", leaning on 2 and 4, with the hi-hat
// and the snaps there, and the bass on every beat.
func TestFullBand(t *testing.T) {
	const f2 = 41 // F2, the root of F7 on bar 1
	for _, c := range []struct {
		name string
		Cue  Cue
		want []string
	}{
		{"beat 1", Cue{Pos: on(1, 1), Key: f2, Snap: true}, []string{
			"bass on the beat 0.8",
			"ride on the beat 0.4",
		}},
		{"beat 2", Cue{Pos: on(1, 2), Key: f2, Snap: true}, []string{
			"bass on the beat 0.8",
			"hi-hat on the beat 0.6",
			"ride on the beat 0.7",
			"ride on the swung and 0.5",
			"snap on the beat 1.0",
		}},
		{"beat 4, the walker not snapping, the player on the bass", Cue{Pos: on(3, 4)}, []string{
			"hi-hat on the beat 0.6",
			"ride on the beat 0.7",
			"ride on the swung and 0.5",
		}},
	} {
		if got := described(Full.Strokes(c.Cue)); !slices.Equal(got, c.want) {
			t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.want)
		}
	}
}

// The game's count-in: "1, 3, 1, 2, 3, 4" on the hi-hat.
func TestCountInOnTheHiHat(t *testing.T) {
	var got []string
	for _, p := range []tempo.Position{on(-1, 1), on(-1, 2), on(-1, 3), on(-1, 4), on(0, 1), on(0, 2), on(0, 3), on(0, 4)} {
		if len(HatCountIn.Strokes(Cue{Pos: p})) > 0 {
			got = append(got, fmt.Sprint(p.Beat))
		}
	}
	if want := []string{"1", "3", "1", "2", "3", "4"}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// The options: the bass and the hi-hat on 2 and 4, nothing else.
func TestLight(t *testing.T) {
	const bb2 = 46 // B♭2
	if got, want := described(Light.Strokes(Cue{Pos: on(2, 1), Key: bb2, Snap: true})), []string{"bass on the beat 0.8"}; !slices.Equal(got, want) {
		t.Errorf("beat 1: got %q, want %q", got, want)
	}
	if got, want := described(Light.Strokes(Cue{Pos: on(2, 2), Key: bb2, Snap: true})), []string{"bass on the beat 0.8", "hi-hat on the beat 0.6"}; !slices.Equal(got, want) {
		t.Errorf("beat 2: got %q, want %q", got, want)
	}
}

// The title's count-in: the snaps alone, on 2 and 4.
func TestIntroSnaps(t *testing.T) {
	if got := SnapCountIn.Strokes(Cue{Pos: on(-1, 1)}); len(got) != 0 {
		t.Errorf("beat 1: got %v, want nothing", got)
	}
	if got, want := described(SnapCountIn.Strokes(Cue{Pos: on(-1, 2)})), []string{"snap on the beat 1.0"}; !slices.Equal(got, want) {
		t.Errorf("beat 2: got %q, want %q", got, want)
	}
}
