package ireal

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

func TestChanges(t *testing.T) {
	const c, cs, d, g harmony.PitchClass = 0, 1, 2, 7
	changes, err := Structure(Lex("[C^7 A7/C#LZD-7 nLZG7XyQZ")).Changes()
	if err != nil {
		t.Fatal(err)
	}
	type want struct {
		root, bass harmony.PitchClass
		silent     bool
		beats      int
	}
	for i, w := range []want{{c, c, false, 2}, {9, cs, false, 2}, {d, d, false, 2}, {0, 0, true, 2}, {g, g, false, 4}} {
		got := changes.Chords[i]
		if got.Silent != w.silent || (!w.silent && (got.Chord.Root != w.root || got.Bass != w.bass)) ||
			got.Length != Ticks(w.beats)*TicksPerBeat {
			t.Errorf("change %d: %+v, want %+v", i, got, w)
		}
	}
	if !changes.Chords[1].Inverted() || changes.Chords[0].Inverted() {
		t.Error("the bass does not tell inversions")
	}
	if !changes.Loops || changes.Next(4) != 0 {
		t.Error("a chart does not loop")
	}
}

// A chord no one can read is a silence, and the error says where.
func TestChangesUnread(t *testing.T) {
	changes, err := Structure(Lex("[C^7XyQ|D*7us*XyQZ")).Changes()
	if err == nil {
		t.Fatal("no error for an unreadable chord")
	}
	if len(changes.Chords) != 2 || changes.Chords[0].Silent || !changes.Chords[1].Silent {
		t.Errorf("%+v", changes.Chords)
	}
}
