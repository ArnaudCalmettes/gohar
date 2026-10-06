package main

import (
	"image/color"
	"math"

	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// drawTeacher draws the walker as he gives a lesson: standing on (`x`,
// `y`), facing the player, `s` times the size of his sketch, in `col`.
// One hand on the hip, the other free; when `snap` is true, the free
// hand is up by his shoulder, snapping, for a right answer.
//
// The sketch is the walker's (see walker.draw), seen from the front:
// the same lengths, angles in radians from the downward vertical,
// positive toward the right of the screen.
func drawTeacher(c screen.Canvas, x, y float32, s float64, snap bool, col color.Color) {
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

		shoulder = 3.0 // from the neck to each shoulder, sideways

		stance = 0.2 // each leg, apart

		pocketX = 1.0 // the hand in the pocket: against the hip, from the middle
		pocketY = 1.5 // and a little under it, where the leg starts

		hang    = 0.12 // the free arm at rest, straight down, a little out
		snapUp  = 0.5  // the free arm out...
		snapFor = 2.6  // ...the forearm up, the hand by the shoulder

		spark     = 0.6
		sparkFrom = 3.0
		sparkTo   = 7.0
	)
	c.Scale *= s
	x, y = x/float32(s), y/float32(s)

	limb := func(fx, fy, length, angle float64) (float64, float64) {
		tx, ty := fx+length*math.Sin(angle), fy+length*math.Cos(angle)
		c.Line(float32(fx), float32(fy), float32(tx), float32(ty), stroke, col)
		return tx, ty
	}

	hipX, hipY := float64(x), float64(y)-(thigh+shin)*math.Cos(stance)
	limb(hipX, hipY, thigh+shin, -stance)
	limb(hipX, hipY, thigh+shin, stance)

	neckY := hipY - torso
	c.Line(float32(hipX), float32(hipY), float32(hipX), float32(neckY), stroke, col)
	c.Circle(float32(hipX), float32(neckY-neckGap*head), float32(head), col)

	// The hand in the pocket, on the left: the hand on the hip, the
	// elbow where the arm's two lengths put it, as little out as they
	// allow.
	shY := neckY + shoulderDrop*head
	lx := hipX - shoulder
	c.Line(float32(hipX), float32(neckY), float32(lx), float32(shY), stroke, col)
	hx, hy := hipX-pocketX, hipY+pocketY
	ex, ey := elbowOut(lx, shY, hx, hy, upper, fore)
	c.Line(float32(lx), float32(shY), float32(ex), float32(ey), stroke, col)
	c.Line(float32(ex), float32(ey), float32(hx), float32(hy), stroke, col)

	// The free arm, on the right: straight down at rest, up snapping.
	rx := hipX + shoulder
	c.Line(float32(hipX), float32(neckY), float32(rx), float32(shY), stroke, col)
	if !snap {
		limb(rx, shY, upper+fore, hang)
		return
	}
	ex, ey = limb(rx, shY, upper, snapUp)
	hx, hy = limb(ex, ey, fore, snapFor)
	for _, a := range []float64{-spark, 0, spark} {
		dx, dy := math.Sin(math.Pi/2+a), -math.Cos(math.Pi/2+a) // out, away from the head
		c.Line(float32(hx+sparkFrom*dx), float32(hy+sparkFrom*dy), float32(hx+sparkTo*dx), float32(hy+sparkTo*dy), stroke/2, col)
	}
}

// elbowOut places the elbow of an arm from the shoulder (`sx`, `sy`) to
// the hand (`hx`, `hy`), `upper` then `fore` long: on the left of the
// line from shoulder to hand, the side away from the body for the left
// arm. A hand out of reach gets the arm straight.
func elbowOut(sx, sy, hx, hy, upper, fore float64) (float64, float64) {
	dx, dy := hx-sx, hy-sy
	d := math.Hypot(dx, dy)
	if d >= upper+fore || d == 0 {
		return sx + dx*upper/(upper+fore), sy + dy*upper/(upper+fore)
	}
	a := (upper*upper - fore*fore + d*d) / (2 * d) // along the line
	h := math.Sqrt(max(upper*upper-a*a, 0))        // across it
	ux, uy := dx/d, dy/d
	return sx + a*ux - h*uy, sy + a*uy + h*ux // (-uy, ux): to the left, going down
}
