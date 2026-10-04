package main

import (
	"math/rand/v2"
	"time"
)

// The coach is what the walker says during a run (see "Le bonhomme" in
// docs/walk.md): the marks in, a phrase now and then out, never a
// judgement. Nothing when the player plays on time: the snaps say it.
//
//   - the last few notes drifting ahead, on average: "tu presses",
//     "détends-toi"; drifting behind: "tu traînes", "ça traîne". A
//     drift, not a fault: a good player rushes within the window on
//     the beat, and the marker calls those notes on time;
//   - a word of praise when it starts to swing, the walker snapping,
//     every few arrivals landed in a row while it swings, and when the
//     player gets back on the beat after a remark.
//
// Each phrase waits a while before another may come: he does not
// harp on. The figures are starting points, to set by playing.
const (
	heard       = 5                     // the last notes weighed
	drift       = 30 * time.Millisecond // their mean offset, ahead or behind: a remark
	back        = drift / 2             // within it again, after a remark: praise
	praiseEvery = 8                     // arrivals landed in a row, snapping
	quiet       = 8                     // beats of silence after a phrase
)

var (
	rushing  = []string{msgRushing, msgRelax}
	dragging = []string{msgDragging, msgItDrags}
	praise   = []string{msgCool, msgYeah, msgKeepItUp, msgSwinging, msgGroovy, msgGreat, msgIDig}
)

type coach struct {
	recent   []time.Duration // the offsets of the last `heard` notes on a beat, oldest first
	until    int             // the first beat he may speak on again
	remarked bool            // he made a remark, and waits for the player to get back on the beat
	snapping bool            // the walker snapped at the last arrival
	streak   int             // arrivals landed in a row while snapping
	last     string          // the last phrase, never said twice in a row
	rng      *rand.Rand
}

func newCoach(rng *rand.Rand) *coach {
	return &coach{rng: rng}
}

// note takes a note played on `beat`, `off` from it (negative when
// early) and marked `t`, and returns the phrase it calls for, empty for
// none. A note between two beats claims none, and does not count.
func (c *coach) note(off time.Duration, t Timing, beat int) string {
	if t == Between {
		return ""
	}
	c.recent = append(c.recent, off)
	if len(c.recent) > heard {
		c.recent = c.recent[1:]
	}
	if len(c.recent) < heard || beat < c.until {
		return ""
	}
	var sum time.Duration
	for _, o := range c.recent {
		sum += o
	}
	mean := sum / heard
	switch {
	case mean <= -drift:
		c.remarked = true
		return c.say(rushing, beat)
	case mean >= drift:
		c.remarked = true
		return c.say(dragging, beat)
	case c.remarked && mean.Abs() <= back:
		c.remarked = false // back on the beat
		return c.say(praise, beat)
	}
	return ""
}

// arrival takes what became of an arrival on `beat`, `landed` or not,
// and whether the walker snaps now: the moment he starts to is worth a
// word, and so is every `praiseEvery` arrival landed in a row after it.
func (c *coach) arrival(beat int, landed, snapping bool) string {
	started := snapping && !c.snapping
	c.snapping = snapping
	switch {
	case !snapping || !landed:
		c.streak = 0
		return ""
	case started:
		c.streak = 0
	default:
		c.streak++
	}
	if (!started && c.streak%praiseEvery != 0) || beat < c.until {
		return ""
	}
	return c.say(praise, beat)
}

// say picks a phrase among `ids`, never the last one said, and keeps
// him quiet for a while.
func (c *coach) say(ids []string, beat int) string {
	id := ids[c.rng.IntN(len(ids))]
	for id == c.last && len(ids) > 1 {
		id = ids[c.rng.IntN(len(ids))]
	}
	c.last = id
	c.until = beat + quiet
	c.recent = c.recent[:0] // the next remark weighs fresh notes
	return id
}
