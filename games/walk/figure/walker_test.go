package figure

import (
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// He walks, grooves once the roots land, snaps after eight in a row,
// and comes back down with the misses, slowly: one miss does not undo
// a groove.
func TestWalkerGait(t *testing.T) {
	var w Walker
	now := time.Now()
	if g := w.Gait(false); g != Searching {
		t.Errorf("at rest: gait %d, want Searching", g)
	}
	if g := w.Gait(true); g != Walking {
		t.Errorf("first beats: gait %d, want Walking", g)
	}
	for range 4 {
		w.Mark(mark.Landed, now)
	}
	if g := w.Gait(true); g != Grooving {
		t.Errorf("four roots landed: gait %d (ease %.2f), want Grooving", g, w.ease)
	}
	for range 4 {
		w.Mark(mark.Landed, now)
	}
	if g := w.Gait(true); g != Snapping {
		t.Errorf("eight roots landed: gait %d (ease %.2f), want Snapping", g, w.ease)
	}
	w.Mark(mark.Missed, now)
	if g := w.Gait(true); g != Grooving || !w.stumble.IsZero() {
		t.Errorf("one miss after eight: gait %d, stumbled %v; want Grooving, no stumble", g, !w.stumble.IsZero())
	}
	for range 3 {
		w.Mark(mark.Missed, now)
	}
	if g := w.Gait(true); g != Walking || w.stumble.IsZero() {
		t.Errorf("four misses: gait %d, stumbled %v; want Walking, a stumble", g, !w.stumble.IsZero())
	}
}
