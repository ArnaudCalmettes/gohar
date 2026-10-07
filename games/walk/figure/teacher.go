package figure

import (
	"image/color"
	"math"

	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// shoulder is from the neck to each shoulder, sideways, from the front.
const shoulder = 3.0

// DrawTeacher draws the walker as he gives a lesson: standing on (`x`,
// `y`), facing the player, `s` times the size of his sketch, in `col`.
// Both hands in his pockets; when `snap` is true, the right one is out,
// up by his shoulder, snapping, for a right answer.
//
// The sketch is the walker's (see Walker.Draw), seen from the front:
// the same lengths, angles in radians from the downward vertical,
// positive toward the right of the screen.
func DrawTeacher(c screen.Canvas, x, y float32, s float64, snap bool, col color.Color) {
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

		stance = 0.2 // each leg, apart

		// The hands in the pockets: at the sides of the hips, low
		// enough for the arms to hang almost straight, the elbows barely
		// out, relaxed. Nearer the middle, both hands would meet where
		// they should not; higher, the elbows flare.
		pocketX = 2.5
		pocketY = 1.8

		snapUp  = 0.5 // the snapping arm out...
		snapFor = 2.6 // ...the forearm up, the hand by the shoulder

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

	// The shoulders, and a hand in each pocket: the elbow where the
	// arm's two lengths put it, out. elbowOut bends to the left going
	// down, out for the left arm: the right one is mirrored there and
	// back.
	shY := neckY + shoulderDrop*head
	pocket := func(side float64) {
		sx, hx, hy := hipX+side*shoulder, hipX+side*pocketX, hipY+pocketY
		m := -side
		ex, ey := elbowOut(m*(sx-hipX), shY, m*(hx-hipX), hy, upper, fore)
		ex = hipX + m*ex
		c.Line(float32(sx), float32(shY), float32(ex), float32(ey), stroke, col)
		c.Line(float32(ex), float32(ey), float32(hx), float32(hy), stroke, col)
	}
	for _, side := range []float64{-1, 1} {
		c.Line(float32(hipX), float32(neckY), float32(hipX+side*shoulder), float32(shY), stroke, col)
	}
	pocket(-1)
	if !snap {
		pocket(1)
		return
	}

	// The right hand out of its pocket, up by the shoulder, snapping.
	ex, ey := limb(hipX+shoulder, shY, upper, snapUp)
	hx, hy := limb(ex, ey, fore, snapFor)
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

// A Pose is how the teacher holds himself in place: standing, or, worn
// down by a player who does not listen, sitting, then meditating (see
// "Rendre le cours amusant" in docs/debutants.md).
type Pose int

const (
	Standing Pose = iota // DrawTeacher
	Sitting              // DrawSitting
	Lotus                // DrawLotus

	// The frames between standing and sitting, held a moment, as a
	// sprite cycle would have them.
	SittingDown // DrawCrouch
	GettingUp   // DrawGettingUp
)

// Draw draws him in pose `p`, his feet (or his seat) on (`x`, `y`), `s`
// times the size of his sketch, in `col`; standing, snapping when
// `snap`.
func (p Pose) Draw(c screen.Canvas, x, y float32, s float64, snap bool, col color.Color) {
	switch p {
	case Sitting:
		DrawSitting(c, x, y, s, col)
	case Lotus:
		DrawLotus(c, x, y, s, col)
	case SittingDown:
		DrawCrouch(c, x, y, s, col)
	case GettingUp:
		DrawGettingUp(c, x, y, s, col)
	default:
		DrawTeacher(c, x, y, s, snap, col)
	}
}

// The resting poses, at scale 1, in logical units and radians.
const (
	sitLean  = 0.55 // the torso, back from the vertical
	sitThigh = 0.8  // the thigh, up from the ground, the knee raised

	lotusLift  = 4.0  // the hip above the ground, on the crossed legs
	lotusKnee  = 11.0 // each knee, sideways from the middle, just off the ground
	lotusFoot  = 4.0  // each foot, on the ground across the middle
	lotusHandX = 10.0 // the hands on the knees, out...
	lotusHandY = -1.0 // ...and a little above the hip, near enough for the elbows to bend in

	// Crouching, sitting down: the seat low, a little behind the feet,
	// the torso leaning back toward the hand behind him on the ground;
	// the other arm out in front.
	downHipX   = -2.0
	downHipY   = 4.0
	downLean   = -0.5 // toward the front, as walker.Draw has it
	downHand   = 9.0  // the hand on the ground, behind the seat
	downArm    = 1.0  // the front arm...
	downElbow  = 0.5  // ...and its forearm, bent down
	crouchFoot = 2.0  // the feet, ahead of `x`, the second one further

	// Getting up: still seated, the knees up and the feet drawn in, the
	// torso forward, a hand on the ground in front of the feet.
	upLean  = 0.3  // toward the front
	upThigh = 1.1  // up from the ground, steeper than sitting
	upHand  = 12.0 // the hand on the ground, ahead of the seat
)

// Head returns the centre of the head of the teacher in pose `p`, his
// feet (or his seat) on (`x`, `y`), `s` times the size of his sketch,
// and its radius: where the tail of his bubble aims.
func (p Pose) Head(x, y, s float64) (hx, hy, r float64) {
	r = sketchHead * s
	switch p {
	case Sitting:
		up := (sketchTorso + sketchNeckGap*sketchHead) * s
		return x - up*math.Sin(sitLean), y - up*math.Cos(sitLean), r
	case Lotus:
		return x, y - (lotusLift+sketchTorso+sketchNeckGap*sketchHead)*s, r
	case SittingDown:
		up := (sketchTorso + sketchNeckGap*sketchHead) * s
		return x + downHipX*s + up*math.Sin(downLean), y - downHipY*s - up*math.Cos(downLean), r
	case GettingUp:
		up := (sketchTorso + sketchNeckGap*sketchHead) * s
		return x + up*math.Sin(upLean), y - up*math.Cos(upLean), r
	}
	return x, y - sketchHeadAbove*s, r
}

// DrawSitting draws him sitting on the grass, his seat on (`x`, `y`),
// facing right: the torso leaning back, the arms straight behind him
// down to the ground, the knees raised, the feet flat.
func DrawSitting(c screen.Canvas, x, y float32, s float64, col color.Color) {
	c.Scale *= s
	x, y = x/float32(s), y/float32(s)
	line := func(ax, ay, bx, by float64) {
		c.Line(float32(ax), float32(ay), float32(bx), float32(by), sketchStroke, col)
	}
	hipX, hipY := float64(x), float64(y)

	// The legs: the thigh up and forward, the shin down to the ground.
	kx, ky := hipX+sketchThigh*math.Cos(sitThigh), hipY-sketchThigh*math.Sin(sitThigh)
	fx := kx + math.Sqrt(max(sketchShin*sketchShin-(hipY-ky)*(hipY-ky), 0))
	line(hipX, hipY, kx, ky)
	line(kx, ky, fx, hipY)

	// The torso and the head, leaning back.
	nx, ny := hipX-sketchTorso*math.Sin(sitLean), hipY-sketchTorso*math.Cos(sitLean)
	line(hipX, hipY, nx, ny)
	hx, hy, _ := Sitting.Head(hipX, hipY, 1)
	c.Circle(float32(hx), float32(hy), sketchHead, col)

	// The arms, straight from the shoulder down to the ground behind.
	sx, sy := nx+sketchShoulderDrop*sketchHead*math.Sin(sitLean), ny+sketchShoulderDrop*sketchHead*math.Cos(sitLean)
	arm := sketchUpper + sketchFore
	back := math.Sqrt(max(arm*arm-(hipY-sy)*(hipY-sy), 0))
	line(sx, sy, sx-back, hipY)
}

// DrawLotus draws him in the lotus position, facing the player, his
// seat above (`x`, `y`): the legs crossed on the ground, each foot on
// the other side, the back straight, the hands on the knees.
func DrawLotus(c screen.Canvas, x, y float32, s float64, col color.Color) {
	c.Scale *= s
	x, y = x/float32(s), y/float32(s)
	line := func(ax, ay, bx, by float64) {
		c.Line(float32(ax), float32(ay), float32(bx), float32(by), sketchStroke, col)
	}
	hipX, hipY := float64(x), float64(y)-lotusLift
	ground := float64(y)

	// The legs: out to the knees, back across to the feet.
	for _, side := range []float64{-1, 1} {
		kx := hipX + side*lotusKnee
		line(hipX, hipY, kx, ground-1)
		line(kx, ground-1, hipX-side*lotusFoot, ground)
	}

	neckY := hipY - sketchTorso
	line(hipX, hipY, hipX, neckY)
	hx, hy, _ := Lotus.Head(hipX, ground, 1)
	c.Circle(float32(hx), float32(hy), sketchHead, col)

	// The arms, relaxed: the elbows in, by the body, the forearms out
	// to the knees. Elbows out, he would seem to lean on them, anything
	// but calm.
	shY := neckY + sketchShoulderDrop*sketchHead
	for _, side := range []float64{-1, 1} {
		sx := hipX + side*shoulder
		hx, hy := hipX+side*lotusHandX, hipY+lotusHandY
		line(hipX, neckY, sx, shY)
		// elbowOut bends to the left going down, in for the right arm:
		// the left one is mirrored there and back.
		m := side
		ex, ey := elbowOut(m*(sx-hipX), shY, m*(hx-hipX), hy, sketchUpper, sketchFore)
		ex = hipX + m*ex
		line(sx, shY, ex, ey)
		line(ex, ey, hx, hy)
	}
}

// DrawCrouch draws him on his way down to sit, facing right, his feet
// by (`x`, `y`): crouching, the seat low, a hand behind him on the
// ground, the other arm out in front.
func DrawCrouch(c screen.Canvas, x, y float32, s float64, col color.Color) {
	c.Scale *= s
	x, y = x/float32(s), y/float32(s)
	line := func(ax, ay, bx, by float64) {
		c.Line(float32(ax), float32(ay), float32(bx), float32(by), sketchStroke, col)
	}
	limb := func(fx, fy, length, angle float64) (float64, float64) {
		tx, ty := fx+length*math.Sin(angle), fy+length*math.Cos(angle)
		line(fx, fy, tx, ty)
		return tx, ty
	}
	ground := float64(y)
	hipX, hipY := float64(x)+downHipX, ground-downHipY

	// The legs, the knees up in front: elbowOut, from the foot to the
	// seat, bends to the left going up, that is forward.
	for _, f := range []float64{crouchFoot, 2 * crouchFoot} {
		fx := float64(x) + f
		kx, ky := elbowOut(fx, ground, hipX, hipY, sketchShin, sketchThigh)
		line(hipX, hipY, kx, ky)
		line(kx, ky, fx, ground)
	}

	nx, ny := limb(hipX, hipY, sketchTorso, math.Pi-downLean) // up the torso
	hx, hy, _ := SittingDown.Head(float64(x), ground, 1)
	c.Circle(float32(hx), float32(hy), sketchHead, col)
	sx := nx - sketchShoulderDrop*sketchHead*math.Sin(downLean)
	sy := ny + sketchShoulderDrop*sketchHead*math.Cos(downLean)

	// The hand behind, a little stretched to reach the ground: a stick
	// figure's arms are short.
	line(sx, sy, hipX-downHand, ground)
	ex, ey := limb(sx, sy, sketchUpper, downArm)
	limb(ex, ey, sketchFore, downArm+downElbow)
}

// DrawGettingUp draws him on his way up, facing right, his seat on
// (`x`, `y`): still seated, the feet drawn in, the torso forward, a
// hand on the ground in front of him to push on. The other arm is
// behind the body, unseen.
func DrawGettingUp(c screen.Canvas, x, y float32, s float64, col color.Color) {
	c.Scale *= s
	x, y = x/float32(s), y/float32(s)
	line := func(ax, ay, bx, by float64) {
		c.Line(float32(ax), float32(ay), float32(bx), float32(by), sketchStroke, col)
	}
	hipX, hipY := float64(x), float64(y)

	kx, ky := hipX+sketchThigh*math.Cos(upThigh), hipY-sketchThigh*math.Sin(upThigh)
	fx := kx + math.Sqrt(max(sketchShin*sketchShin-(hipY-ky)*(hipY-ky), 0))
	line(hipX, hipY, kx, ky)
	line(kx, ky, fx, hipY)

	nx, ny := hipX+sketchTorso*math.Sin(upLean), hipY-sketchTorso*math.Cos(upLean)
	line(hipX, hipY, nx, ny)
	hx, hy, _ := GettingUp.Head(hipX, hipY, 1)
	c.Circle(float32(hx), float32(hy), sketchHead, col)

	// The hand ahead, a little stretched to reach the ground, as
	// DrawCrouch's behind.
	sx := nx - sketchShoulderDrop*sketchHead*math.Sin(upLean)
	sy := ny + sketchShoulderDrop*sketchHead*math.Cos(upLean)
	line(sx, sy, hipX+upHand, hipY)
}
