package main

import (
	"math/rand/v2"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// A jam is the music of the menus: the reference bass walking the blues
// without a break, from the title to the options and the calibration
// and back, each screen with its own arrangement. It lives in the app,
// not in a scene: a scene on top schedules it at each tick, and the
// next one takes over on the same beat. The game stops it, and counts
// in on its own.
type jam struct {
	band *band
	m    tempo.Metronome
	bass *bassist
	next int // the next beat to schedule; negative in the count-in
}

// newJam starts the count-in a lookahead from now: two bars of the
// walker's snaps on 2 and 4, alone, then the band comes in, at
// `titleBPM`. Without `snapIn`, the band comes in at once: on a screen
// where the walker does not snap, the count-in's snaps would come from
// nowhere.
func newJam(a *app, snapIn bool) *jam {
	now := time.Now()
	in := 0
	if snapIn {
		in = countIn
	}
	beat := time.Duration(float64(time.Minute) / titleBPM)
	m := tempo.NewMetronome(now.Add(lookahead+time.Duration(in)*beat), titleBPM, perBar)
	return &jam{band: a.band, m: m, bass: newBassist(a), next: -in}
}

// play queues the beats due within the lookahead of `now`: the snaps of
// the count-in, then the bass line in arrangement `a`, the scene's.
// The walker of the menus always snaps: the parts that want him, play.
func (j *jam) play(now time.Time, a arrangement) {
	_, end := j.m.Due(now, now.Add(lookahead))
	for ; j.next < end; j.next++ {
		c := cue{pos: j.m.Position(j.next), snap: true}
		if j.next < 0 {
			j.band.play(snapCountIn.strokes(c), j.next, j.m)
			continue
		}
		c.key = j.bass.key(j.next)
		j.band.play(a.strokes(c), j.next, j.m)
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
		beats: Expect(a.tunes[0].grid, tempo.NewMetronome(time.Time{}, a.bpm, perBar), 1), // the blues; the metronome only numbers the beats
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
