package main

import (
	"image/color"
	"math"
	"time"
)

// The stick figure of the first milestone (see "Le bonhomme" in
// docs/walk.md): drawn with strokes by the engine, so that the walk can
// be judged against the game before any drawing is ordered. His step
// follows the metronome, his gait the player's ease, with inertia: one
// missed root does not trip him, a run of them does.

// A gait is the walker's state, from the lowest to the highest.
type gait int

const (
	searching gait = iota // looking for the tempo: count-in, practice, at rest
	walking               // steady steps, the default
	grooving              // bouncing on the beats, the head follows
	snapping              // it has been rolling for a while: snaps on 2 and 4
)

const (
	// easeWeight is how much one mark moves the ease: about the last
	// eight arrivals count.
	easeWeight = 0.25

	grooveFrom = 0.6
	snapFrom   = 0.85
	snapStreak = 8 // arrivals landed in a row before he snaps

	hopTime     = 250 * time.Millisecond
	stumbleTime = 400 * time.Millisecond
)

type walker struct {
	ease   float64 // 0 to 1, a moving average of landed against missed
	streak int     // arrivals landed in a row

	hop, stumble time.Time // when the last reaction started
}

// mark takes note of what became of an arrival: a landed one lifts the
// ease and makes him hop, a missed or doubled one lowers it, and he
// stumbles only when the ease has fallen below the walk.
func (w *walker) mark(k BeatKind, now time.Time) {
	if k == Landed {
		w.ease += easeWeight * (1 - w.ease)
		w.streak++
		w.hop = now
		return
	}
	w.ease -= easeWeight * w.ease
	w.streak = 0
	if w.ease < grooveFrom {
		w.stumble = now
	}
}

// gait returns the state to draw. `tempo` is false at rest, during the
// count-in and in practice: he looks for the tempo.
func (w *walker) gait(tempo bool) gait {
	switch {
	case !tempo:
		return searching
	case w.ease >= snapFrom && w.streak >= snapStreak:
		return snapping
	case w.ease >= grooveFrom:
		return grooving
	}
	return walking
}

// draw draws him standing on (`x`, `y`), facing right, `s` times the
// size of the sketch his lengths are written for, in `col`: ink
// for the player, grey when the demo plays. `beats` is where
// the music is, in beats. Two frames, as a sprite cycle would have,
// never in between:
//   - on the strong beats, 1 and 3, the legs apart, one forward and one
//     back, and the arms swinging against them;
//   - on the weak beats, 2 and 4, the legs crossing: one straight under
//     the hip, the other with the knee bent and the foot lifted behind,
//     and the arms down.
//
// The hip sits lower with the legs apart, which gives the bounce for
// free. The reactions are frames held for their length: up for a hop,
// leaning for a stumble.
func (w *walker) draw(c canvas, x, y float32, s float64, g gait, beats float64, now time.Time, col color.Color) {
	var (
		stroke = float32(1.5 * s)
		thigh  = 8.0 * s
		shin   = 8.0 * s
		upper  = 6.0 * s // the arm, shoulder to elbow
		fore   = 7.0 * s // elbow to hand
		torso  = 14.0 * s
		head   = 5.0 * s
	)
	beat := int(math.Floor(beats))
	strong := beat%2 == 0

	stride := 0.4 // radians, each leg from the vertical
	switch g {
	case searching:
		stride, strong = 0.1, true // a shuffle on the spot, no crossing
	case grooving, snapping:
		stride = 0.5
	}

	lift, lean := 0.0, 0.0
	if now.Sub(w.hop) < hopTime {
		lift = 5 * s
	}
	if now.Sub(w.stumble) < stumbleTime {
		lean = 0.3
	}

	hipX, hipY := float64(x), float64(y)-(thigh+shin)-lift
	if strong {
		hipY = float64(y) - (thigh+shin)*math.Cos(stride) - lift
	}
	// limb draws a segment from a point at an angle from the downward
	// vertical, positive toward the front, and returns its end.
	limb := func(fx, fy, length, angle float64) (float64, float64) {
		tx, ty := fx+length*math.Sin(angle), fy+length*math.Cos(angle)
		c.line(float32(fx), float32(fy), float32(tx), float32(ty), stroke, col)
		return tx, ty
	}

	if strong {
		limb(hipX, hipY, thigh+shin, stride)
		limb(hipX, hipY, thigh+shin, -stride)
	} else {
		limb(hipX, hipY, thigh+shin, 0)
		kx, ky := limb(hipX, hipY, thigh, 0.5) // the knee forward
		limb(kx, ky, shin, -0.9)               // the foot up behind
	}

	neckX := hipX + torso*math.Sin(lean)
	neckY := hipY - torso*math.Cos(lean)
	c.line(float32(hipX), float32(hipY), float32(neckX), float32(neckY), stroke, col)
	c.circle(float32(neckX+head*math.Sin(lean)), float32(neckY-1.2*head), float32(head), col)

	// The arms, each with its elbow: against the legs on the strong
	// beats, down on the weak ones.
	shX, shY := neckX, neckY+0.6*head
	swing := 0.0
	if strong {
		swing = stride
	}
	ex, ey := limb(shX, shY, upper, -swing)
	limb(ex, ey, fore, -swing+0.3)
	if g != snapping {
		ex, ey = limb(shX, shY, upper, swing)
		limb(ex, ey, fore, swing+0.3)
		return
	}

	// The jazz snap, from the elbow: very bent on 1 and 3, the hand up
	// by the chest; opening on 2 and 4, the hand out in front, where it
	// snaps, sparks at the fingers.
	ex, ey = limb(shX, shY, upper, 0.3)
	if strong {
		limb(ex, ey, fore, 2.6) // folded, the hand up and forward
		return
	}
	hx, hy := limb(ex, ey, fore, 1.3) // open, the hand forward
	for _, a := range []float64{-0.6, 0, 0.6} {
		dx, dy := math.Cos(a), math.Sin(a) // away from the body
		c.line(float32(hx+3*s*dx), float32(hy+3*s*dy), float32(hx+7*s*dx), float32(hy+7*s*dy), stroke/2, col)
	}
}
