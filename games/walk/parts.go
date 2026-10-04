package main

import "github.com/ArnaudCalmettes/gohar/games/tempo"

// The band's parts and arrangements (see "L'orchestre" in docs/walk.md).
//
// A part is what one player of the band does on a beat: the ride, the
// hi-hat, the snaps, the bass. It is a pure function, from what the beat
// is to the strokes it calls for, so that a rhythm reads, and is tested,
// as a musician would say it. An arrangement is the parts playing
// together: the title's, the options', the count-in's. The band only
// plays the strokes (see band.play).

// A sound is what a stroke strikes. The band knows the keys behind each,
// which differ with the soundfont and the chip sounds.
type sound int

const (
	hatSound  sound = iota // the hi-hat, closed by the foot
	rideSound              // the ride cymbal
	snapSound              // the walker's finger snap
	bassSound              // the reference bass, legato: its key is the stroke's
)

// A stroke is one sound, at a velocity from 0 to 1, `at` beats after
// the start of its beat: 0 on the beat, `swingRatio` on the swung "and".
type stroke struct {
	sound sound
	key   int // the bass's note; unused for the others
	vel   float64
	at    float64
}

// A cue is what a part knows of the beat it plays on.
type cue struct {
	pos  tempo.Position // bar and beat, the count-in in bars 0 and below
	key  int            // the note the bass line has on it, 0 for none
	snap bool           // the walker snaps his fingers
}

// backbeat tells the beats a jazz drummer leans on: 2 and 4.
func (c cue) backbeat() bool { return c.pos.Beat%2 == 0 }

// A part is one player's strokes on a beat.
type part func(c cue) []stroke

// An arrangement is the parts that play together.
type arrangement []part

// strokes gathers what every part plays on `c`.
func (a arrangement) strokes(c cue) []stroke {
	var all []stroke
	for _, p := range a {
		all = append(all, p(c)...)
	}
	return all
}

// swingRatio places the ride's "and", the triplet swing of the method
// books: two thirds of the beat, the long-short of the eighths.
const swingRatio = 2.0 / 3

// The velocities of the parts, 0 to 1, set by ear.
const (
	countVel    = 0.7 // the hi-hat of the count-in
	backbeatVel = 0.7 // the ride on 2 and 4
	beatVel     = 0.4 // the ride on 1 and 3
	andVel      = 0.5 // the ride on the "and" of 2 and 4
	hatVel      = 0.6 // the hi-hat on 2 and 4
	endVel      = 0.7 // the ride on the demo's last note
	snapVel     = 1
	bassVel     = 0.8  // the reference bass
	underVel    = 0.55 // the same, discreet, under the calibration
)

// The parts.
var (
	// countInHat counts the game in on the hi-hat, two bars: 1 and 3 in
	// the first, every beat in the second, "1, 3, 1, 2, 3, 4", the time
	// to move the hands from the space bar to the keyboard.
	countInHat part = func(c cue) []stroke {
		if c.pos.Bar < 0 && c.backbeat() {
			return nil
		}
		return []stroke{{sound: hatSound, vel: countVel}}
	}

	// countInSnaps counts the title in with the walker's snaps alone, on
	// 2 and 4.
	countInSnaps part = func(c cue) []stroke {
		if !c.backbeat() {
			return nil
		}
		return []stroke{{sound: snapSound, vel: snapVel}}
	}

	// ride is the jazz ride pattern: every beat, leaning on 2 and 4, and
	// the swung "and" after 2 and 4: "ding, ding-da, ding, ding-da".
	ride part = func(c cue) []stroke {
		if !c.backbeat() {
			return []stroke{{sound: rideSound, vel: beatVel}}
		}
		return []stroke{{sound: rideSound, vel: backbeatVel}, {sound: rideSound, vel: andVel, at: swingRatio}}
	}

	// hatOnTwoAndFour closes the hi-hat with the foot on 2 and 4.
	hatOnTwoAndFour part = func(c cue) []stroke {
		if !c.backbeat() {
			return nil
		}
		return []stroke{{sound: hatSound, vel: hatVel}}
	}

	// snaps are the walker's on 2 and 4, when he snaps: the juice.
	snaps part = func(c cue) []stroke {
		if !c.snap || !c.backbeat() {
			return nil
		}
		return []stroke{{sound: snapSound, vel: snapVel}}
	}

	// bass walks the line, when the band has one: none in the count-in,
	// none when the player walks it.
	bass part = walkingAt(bassVel)

	// softBass is the same, discreet.
	softBass part = walkingAt(underVel)

	// endRide strikes the ride once, on the demo's last note.
	endRide part = func(cue) []stroke {
		return []stroke{{sound: rideSound, vel: endVel}}
	}
)

// walkingAt is the bass's part at `vel`.
func walkingAt(vel float64) part {
	return func(c cue) []stroke {
		if c.key == 0 {
			return nil
		}
		return []stroke{{sound: bassSound, key: c.key, vel: vel}}
	}
}

// The arrangements.
var (
	hatCountIn  = arrangement{countInHat}                         // the game's count-in
	snapCountIn = arrangement{countInSnaps}                       // the title's count-in
	full        = arrangement{ride, hatOnTwoAndFour, snaps, bass} // the game, the title
	light       = arrangement{hatOnTwoAndFour, bass}              // the options
	bassOnly    = arrangement{softBass}                           // under the calibration's bar
	ending      = arrangement{endRide, bass}                      // the demo's last note
)
