package figure

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// The phrases of the tests, said as the French locale says them.
var (
	rushing  = []string{"Tu presses", "Détends-toi"}
	dragging = []string{"Tu traînes", "Ça traîne"}
	praise   = []string{"Cool !", "Yeah !", "Continue comme ça !", "Ça swingue !", "Groovy !", "Super !", "Ça joue !"}
)

func testCoach() *Coach {
	return NewCoach(rand.New(rand.NewPCG(1, 2)), Lines{Rushing: rushing, Dragging: dragging, Praise: praise})
}

// play hands `offs`, in milliseconds (negative when early), to `c`, a
// note per beat from `from`, all on time for the marker, and returns
// the phrases said, by beat.
func play(c *Coach, from int, offs ...int) map[int]string {
	said := map[int]string{}
	for i, o := range offs {
		if id := c.Note(time.Duration(o)*time.Millisecond, mark.OnTime, from+i); id != "" {
			said[from+i] = id
		}
	}
	return said
}

// Within the window on the beat, but drifting ahead, 40 ms on average:
// "tu presses" or "détends-toi", once, then silence for a while.
func TestRushing(t *testing.T) {
	c := testCoach()
	said := play(c, 0, -30, -50, -35, -45, -40, -40, -40)
	if len(said) != 1 || !slices.Contains(rushing, said[4]) {
		t.Errorf("got %v, want a remark on rushing at the fifth note, and no more", said)
	}
}

// Drifting behind: "tu traînes" or "ça traîne".
func TestDragging(t *testing.T) {
	c := testCoach()
	said := play(c, 0, 30, 40, 50, 35, 45)
	if !slices.Contains(dragging, said[4]) {
		t.Errorf("got %v, want a remark on dragging", said)
	}
}

// The spread of a player on time, -30 to +20 ms around the beat: he
// says nothing, the snaps say it.
func TestOnTimeIsSilent(t *testing.T) {
	c := testCoach()
	if said := play(c, 0, -30, 20, -10, 5, -25, 15, 0, -20, 10, -5); len(said) > 0 {
		t.Errorf("got %v, want nothing", said)
	}
}

// A note between two beats does not count.
func TestBetweenDoesNotCount(t *testing.T) {
	c := testCoach()
	for i := range heard {
		if id := c.Note(-200*time.Millisecond, mark.Between, i); id != "" {
			t.Fatalf("got %q for a note between two beats", id)
		}
	}
}

// After a remark, five notes back around the beat: a word of praise.
func TestBackOnTheBeat(t *testing.T) {
	c := testCoach()
	play(c, 0, -50, -50, -50, -50, -50)
	said := play(c, 8, -10, 5, 0, -5, 10)
	if !slices.Contains(praise, said[12]) {
		t.Errorf("got %v, want praise at the fifth note back on the beat", said)
	}
}

// The walker starts to snap: a word of praise; then one every eight
// arrivals landed in a row; a missed one starts the count again.
func TestPraiseWhileItSwings(t *testing.T) {
	c := testCoach()
	if id := c.Arrival(0, true, false); id != "" {
		t.Errorf("walking: got %q, want nothing", id)
	}
	if id := c.Arrival(4, true, true); !slices.Contains(praise, id) {
		t.Errorf("starts to snap: got %q, want praise", id)
	}
	beat := 4
	arrive := func(landed bool) string { beat += 4; return c.Arrival(beat, landed, true) }
	for i := 1; i < praiseEvery; i++ {
		if id := arrive(true); id != "" {
			t.Fatalf("arrival %d after the start: got %q, want nothing", i, id)
		}
	}
	if id := arrive(true); !slices.Contains(praise, id) {
		t.Errorf("eighth arrival landed: got %q, want praise", id)
	}
	arrive(false)
	for i := 1; i < praiseEvery; i++ {
		if id := arrive(true); id != "" {
			t.Fatalf("arrival %d after a miss: got %q, want nothing", i, id)
		}
	}
}

// Never the same word twice in a row.
func TestNeverTwiceInARow(t *testing.T) {
	c := testCoach()
	last := ""
	for i := range 50 {
		id := c.say(praise, i*quiet)
		if id == last {
			t.Fatalf("%q twice in a row", id)
		}
		last = id
	}
}
