package main

import (
	"image/color"
	"math"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/screen"
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

// The walker's sketch, at scale 1, as walker.draw and drawTeacher draw
// him: lengths in logical units, the same in profile and from the front.
const (
	sketchStroke = 1.5
	sketchThigh  = 8.0
	sketchShin   = 8.0
	sketchTorso  = 14.0
	sketchHead   = 5.0 // its radius

	sketchNeckGap      = 1.2 // the centre of the head above the neck, in radii
	sketchShoulderDrop = 0.6 // the shoulders below the neck, in radii
	sketchUpper        = 6.0 // the arm, shoulder to elbow
	sketchFore         = 7.0 // elbow to hand

	// sketchHeadAbove is the centre of his head above his feet, standing
	// straight: where the tail of a bubble aims.
	sketchHeadAbove = sketchThigh + sketchShin + sketchTorso + sketchNeckGap*sketchHead
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
// size of his sketch, in `col`: ink
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
func (w *walker) draw(c screen.Canvas, x, y float32, s float64, g gait, beats float64, now time.Time, col color.Color) {
	// The sketch, at scale 1: lengths in logical units, angles in
	// radians from the downward vertical, positive toward the front.
	const (
		stroke       = sketchStroke
		thigh        = sketchThigh
		shin         = sketchShin
		torso        = sketchTorso
		head         = sketchHead
		neckGap      = sketchNeckGap
		shoulderDrop = sketchShoulderDrop
		upper        = sketchUpper
		fore         = sketchFore

		stride       = 0.4 // each leg, walking
		grooveStride = 0.5
		shuffle      = 0.1  // looking for the tempo, on the spot
		knee         = 0.5  // the crossing leg: the knee forward...
		foot         = -0.9 // ...the foot up behind
		elbow        = 0.3  // the forearm, from the upper arm

		hopLift = 5.0
		lean    = 0.3 // stumbling

		snapUpper  = 0.3 // the snapping arm, up front
		snapFolded = 2.6 // on 1 and 3, the hand up by the chest
		snapOpen   = 1.3 // on 2 and 4, the hand out, snapping
		spark      = 0.6 // the spread of the sparks
		sparkFrom  = 3.0
		sparkTo    = 7.0
	)
	// Drawn at scale 1 on a canvas `s` times larger.
	c.Scale *= s
	x, y = x/float32(s), y/float32(s)

	beat := int(math.Floor(beats))
	strong := beat%2 == 0

	legs := stride
	switch g {
	case searching:
		legs, strong = shuffle, true // no crossing
	case grooving, snapping:
		legs = grooveStride
	}

	lift, tilt := 0.0, 0.0
	if now.Sub(w.hop) < hopTime {
		lift = hopLift
	}
	if now.Sub(w.stumble) < stumbleTime {
		tilt = lean
	}

	hipX, hipY := float64(x), float64(y)-(thigh+shin)-lift
	if strong {
		hipY = float64(y) - (thigh+shin)*math.Cos(legs) - lift
	}
	// limb draws a segment from a point at an angle from the downward
	// vertical, positive toward the front, and returns its end.
	limb := func(fx, fy, length, angle float64) (float64, float64) {
		tx, ty := fx+length*math.Sin(angle), fy+length*math.Cos(angle)
		c.Line(float32(fx), float32(fy), float32(tx), float32(ty), stroke, col)
		return tx, ty
	}

	if strong {
		limb(hipX, hipY, thigh+shin, legs)
		limb(hipX, hipY, thigh+shin, -legs)
	} else {
		limb(hipX, hipY, thigh+shin, 0)
		kx, ky := limb(hipX, hipY, thigh, knee)
		limb(kx, ky, shin, foot)
	}

	neckX := hipX + torso*math.Sin(tilt)
	neckY := hipY - torso*math.Cos(tilt)
	c.Line(float32(hipX), float32(hipY), float32(neckX), float32(neckY), stroke, col)
	c.Circle(float32(neckX+head*math.Sin(tilt)), float32(neckY-neckGap*head), float32(head), col)

	// The arms, each with its elbow: against the legs on the strong
	// beats, down on the weak ones.
	shX, shY := neckX, neckY+shoulderDrop*head
	swing := 0.0
	if strong {
		swing = legs
	}
	ex, ey := limb(shX, shY, upper, -swing)
	limb(ex, ey, fore, -swing+elbow)
	if g != snapping {
		ex, ey = limb(shX, shY, upper, swing)
		limb(ex, ey, fore, swing+elbow)
		return
	}

	// The jazz snap, from the elbow: very bent on 1 and 3, the hand up
	// by the chest; opening on 2 and 4, the hand out in front, where it
	// snaps, sparks at the fingers.
	ex, ey = limb(shX, shY, upper, snapUpper)
	if strong {
		limb(ex, ey, fore, snapFolded)
		return
	}
	hx, hy := limb(ex, ey, fore, snapOpen)
	for _, a := range []float64{-spark, 0, spark} {
		dx, dy := math.Cos(a), math.Sin(a) // away from the body
		c.Line(float32(hx+sparkFrom*dx), float32(hy+sparkFrom*dy), float32(hx+sparkTo*dx), float32(hy+sparkTo*dy), stroke/2, col)
	}
}
