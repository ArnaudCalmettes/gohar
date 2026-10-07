package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/band"
	"github.com/ArnaudCalmettes/gohar/games/walk/figure"
	"github.com/ArnaudCalmettes/gohar/games/walk/lessons"
)

// The easter egg of the lessons (see "Rendre le cours amusant" in
// docs/debutants.md): ten casseroles in a row on a request nobody misses
// in good faith, and the walker, fed up, walks out of the screen. GAME
// OVER. A key, and he walks back; the lesson goes on where it was,
// nothing lost.

// The GAME OVER jingle, on the piano, high: four tritones going down
// a semitone at a time, the last one shaken in a tremolo.
const (
	jingleTop     = 83 // B5, the top note of the first tritone
	jingleSteps   = 4
	jingleBPM     = 140 // a tritone a beat
	jingleNote    = time.Minute / jingleBPM
	jingleTremolo = 3 * jingleNote        // the last tritone
	jingleShake   = 45 * time.Millisecond // each of its notes, in turn
	jingleVel     = 0.7
	tritone       = 6 // semitones
)

// The stages of a walkout.
type walkoutStage int

const (
	inPlace   walkoutStage = iota // teaching
	leaving                       // walking out to the right
	gone                          // GAME OVER
	returning                     // walking back from the right
	sulking                       // back in place, his last word said
)

const (
	walkTime  = 2500 * time.Millisecond // out, or back
	walkSteps = 2.5                     // beats a second

	offScreen = screenWidth + 40 // where he is out of sight, his feet

	// The box he is drawn in to come back facing left, mirrored: around
	// his feet, wide enough for his stride.
	mirrorW = 90
	mirrorH = 160
)

// A walkout is where the walker is in his walkout, and since when.
type walkout struct {
	stage  walkoutStage
	since  time.Time
	mirror *ebiten.Image // the box for walking back, mirrored
}

// Tease answers misses in a row on a teasing request, once the key is
// up: a word in the bubble, then the walkout.
func (l *lesson) Tease(a lessons.Annoyance) {
	l.react(func() {
		switch a {
		case lessons.Hey:
			l.Say(msgTeaseHey)
		case lessons.OnPurpose:
			l.Say(msgTeaseOnPurpose)
		case lessons.FedUp:
			l.Say(msgTeaseFedUp)
			l.out = walkout{stage: leaving, since: time.Now(), mirror: l.out.mirror}
			l.shown = nil
		}
	})
}

// teaching tells whether the walker is there to teach: the lesson's
// runner hears the keys only then.
func (l *lesson) teaching() bool { return l.out.stage == inPlace }

// walkoutKey moves the walkout on at a key, of either keyboard: GAME
// OVER brings him back; back, the step he left starts again.
func (l *lesson) walkoutKey(now time.Time) {
	switch l.out.stage {
	case gone:
		l.out.stage, l.out.since = returning, now
	case sulking:
		l.out.stage = inPlace
		l.runner.Again()
	}
}

// walkoutTick moves the walkout on with time.
func (l *lesson) walkoutTick(now time.Time) {
	t := now.Sub(l.out.since)
	switch {
	case l.out.stage == leaving && t >= walkTime:
		l.out.stage, l.out.since = gone, now
		l.bubble = nil
		l.jingle(now)
	case l.out.stage == returning && t >= walkTime:
		l.out.stage, l.out.since = sulking, now
		l.Say(msgTeaseBack)
	}
}

// jingle plays the GAME OVER jingle from `now`: the first three
// tritones struck together, the last one in a tremolo, its two notes in
// turn.
func (l *lesson) jingle(now time.Time) {
	p := l.band.Piano
	at := now.Add(band.Lookahead)
	hold := time.Duration(phraseHold * float64(jingleNote.Nanoseconds())) // not a constant: rounded, not refused
	for i := range jingleSteps - 1 {
		hi := jingleTop - i
		for _, k := range []int{hi - tritone, hi} {
			p.ScheduleOn(k, jingleVel, at)
			p.ScheduleOff(k, at.Add(hold))
		}
		at = at.Add(jingleNote)
	}
	hi := jingleTop - (jingleSteps - 1)
	for i := range int(jingleTremolo / jingleShake) {
		k := hi - tritone*(1-i%2) // the low note first
		p.ScheduleOn(k, jingleVel, at)
		p.ScheduleOff(k, at.Add(jingleShake))
		at = at.Add(jingleShake)
	}
}

// walkerAt is where the walker's feet are now, across the screen.
func (l *lesson) walkerAt(now time.Time) float64 {
	t := float64(now.Sub(l.out.since))
	switch l.out.stage {
	case leaving:
		return walkerX + (offScreen-walkerX)*min(t/float64(walkTime), 1)
	case gone:
		return offScreen
	case returning:
		return offScreen - (offScreen-walkerX)*min(t/float64(walkTime), 1)
	}
	return walkerX
}

// drawWalkout draws the walker as the walkout has him: walking out in
// profile, GAME OVER while he is gone, coming back mirrored, facing
// left. It returns false while he is in place, for the teacher's pose.
func (l *lesson) drawWalkout(dst *ebiten.Image, c screen.Canvas, now time.Time) bool {
	t := now.Sub(l.out.since).Seconds()
	switch l.out.stage {
	case leaving:
		l.walker.Draw(c, float32(l.walkerAt(now)), walkerY, walkerScale, figure.Walking, t*walkSteps, now, ink)
	case gone:
		c.Centred(gameOver, l.fonts.count, screenWidth/2, headingY, ink)
	case returning:
		l.drawMirrored(dst, l.walkerAt(now), t*walkSteps, now)
	default:
		return false
	}
	return true
}

// gameOver is the same in every language: the screen of every game.
const gameOver = "GAME OVER"

// drawMirrored draws the walker facing left, his feet at `x`: drawn
// facing right in a box of his own, then flipped onto `dst`.
func (l *lesson) drawMirrored(dst *ebiten.Image, x, beats float64, now time.Time) {
	w, h := int(mirrorW*l.scale), int(mirrorH*l.scale)
	if l.out.mirror == nil || l.out.mirror.Bounds().Dx() != w || l.out.mirror.Bounds().Dy() != h {
		l.out.mirror = ebiten.NewImage(w, h)
	}
	l.out.mirror.Clear()
	box := screen.Canvas{Dst: l.out.mirror, Scale: l.scale}
	l.walker.Draw(box, mirrorW/2, mirrorH-5, walkerScale, figure.Walking, beats, now, ink)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate((x+mirrorW/2)*l.scale, (walkerY-mirrorH+5)*l.scale)
	dst.DrawImage(l.out.mirror, op)
}
