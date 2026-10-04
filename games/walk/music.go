package main

import (
	"math/rand/v2"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// A jam is the music of the menus: the reference bass walking the blues
// without a break, from the title to the calibration and back, with the
// drums and the snaps or without. It lives in the app, not in a scene:
// a scene on top schedules it at each tick, and the next one takes over
// on the same beat. The game stops it, and counts in on its own.
type jam struct {
	band  *band
	m     tempo.Metronome
	swing tempo.Swing
	bass  *bassist
	next  int // the next beat to schedule; negative in the count-in
}

// newJam starts the count-in a lookahead from now: two bars of the
// walker's snaps on 2 and 4, alone, then the band comes in, at
// `titleBPM`.
func newJam(a *app) *jam {
	now := time.Now()
	beat := time.Duration(float64(time.Minute) / titleBPM)
	m := tempo.NewMetronome(now.Add(lookahead+countIn*beat), titleBPM, perBar)
	return &jam{band: a.band, m: m, swing: tempo.NewSwing(m, swingRatio), bass: newBassist(a), next: -countIn}
}

// play queues the beats due within the lookahead of `now`: the snaps of
// the count-in, then the bass, with the drums and the snaps when
// `drums`, discreet without them.
func (j *jam) play(now time.Time, drums bool) {
	_, end := j.m.Due(now, now.Add(lookahead))
	for ; j.next < end; j.next++ {
		at := j.m.At(j.next)
		switch {
		case j.next < 0:
			if j.m.Position(j.next).Beat%2 == 0 { // 2 and 4
				j.band.snapAt(at)
			}
		case drums:
			j.band.beat(j.m.Position(j.next), j.bass.key(j.next), at, j.swing.AtBeats(float64(j.next)+0.5), true)
		default:
			j.band.walk(j.bass.key(j.next), underVel, at)
		}
	}
}

// stop releases the bass.
func (j *jam) stop(now time.Time) {
	j.band.stop(now)
}

// A bassist walks the chart chorus after chorus, a new line drawn at
// each one. Each chorus is walked on its own: the line does not carry
// over the double bar yet.
type bassist struct {
	beats []Beat // one chorus
	line  []int  // the chorus playing: a key per beat
	rng   *rand.Rand
}

func newBassist(a *app) *bassist {
	seed := uint64(time.Now().UnixNano())
	return &bassist{
		beats: Expect(a.grid, tempo.NewMetronome(time.Time{}, a.bpm, perBar), 1), // only the numbering
		rng:   rand.New(rand.NewPCG(seed, seed>>32|1)),
	}
}

// key returns the key of beat `n`, counted from 0 at the first chorus.
func (b *bassist) key(n int) int {
	i := n % len(b.beats)
	if i == 0 || b.line == nil {
		b.line = Walk(b.beats, b.rng)
	}
	return b.line[i]
}
