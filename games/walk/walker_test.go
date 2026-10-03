package main

import (
	"testing"
	"time"
)

// He walks, grooves once the roots land, snaps after eight in a row,
// and comes back down with the misses, slowly: one miss does not undo
// a groove.
func TestWalkerGait(t *testing.T) {
	var w walker
	now := time.Now()
	if g := w.gait(false); g != searching {
		t.Errorf("at rest: gait %d, want searching", g)
	}
	if g := w.gait(true); g != walking {
		t.Errorf("first beats: gait %d, want walking", g)
	}
	for range 4 {
		w.mark(Landed, now)
	}
	if g := w.gait(true); g != grooving {
		t.Errorf("four roots landed: gait %d (ease %.2f), want grooving", g, w.ease)
	}
	for range 4 {
		w.mark(Landed, now)
	}
	if g := w.gait(true); g != snapping {
		t.Errorf("eight roots landed: gait %d (ease %.2f), want snapping", g, w.ease)
	}
	w.mark(Missed, now)
	if g := w.gait(true); g != grooving || !w.stumble.IsZero() {
		t.Errorf("one miss after eight: gait %d, stumbled %v; want grooving, no stumble", g, !w.stumble.IsZero())
	}
	for range 3 {
		w.mark(Missed, now)
	}
	if g := w.gait(true); g != walking || w.stumble.IsZero() {
		t.Errorf("four misses: gait %d, stumbled %v; want walking, a stumble", g, !w.stumble.IsZero())
	}
}
