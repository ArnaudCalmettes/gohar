package band

import (
	"math/rand/v2"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/games/walk/bass"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// A Jam is the music of the menus: the reference bass walking the blues
// without a break, from the title to the options and the calibration
// and back, each screen with its own arrangement. It lives in the app,
// not in a scene: a scene on top schedules it at each tick, and the
// next one takes over on the same beat. The game stops it, and counts
// in on its own.
type Jam struct {
	band *Band
	m    tempo.Metronome
	bass *bassist
	next int // the next beat to schedule; negative in the count-in
}

// NewJam starts the count-in a lookahead from now: two bars of the
// walker's snaps on 2 and 4, alone, then `b` comes in, walking `grid`
// at `bpm`. Without `snapIn`, the band comes in at once: on a screen
// where the walker does not snap, the count-in's snaps would come from
// nowhere.
func NewJam(b *Band, grid analysis.Changes, bpm float64, snapIn bool) *Jam {
	now := time.Now()
	in := 0
	if snapIn {
		in = CountIn
	}
	beat := time.Duration(float64(time.Minute) / bpm)
	m := tempo.NewMetronome(now.Add(Lookahead+time.Duration(in)*beat), bpm, mark.BeatsPerBar)
	return &Jam{band: b, m: m, bass: newBassist(grid), next: -in}
}

// Metronome is the jam's clock: where the music is, for the screens
// that move on its beats.
func (j *Jam) Metronome() tempo.Metronome { return j.m }

// Play queues the beats due within the lookahead of `now`: the snaps of
// the count-in, then the bass line in arrangement `a`, the scene's.
// The walker of the menus always snaps: the parts that want him, play.
func (j *Jam) Play(now time.Time, a Arrangement) {
	_, end := j.m.Due(now, now.Add(Lookahead))
	for ; j.next < end; j.next++ {
		c := Cue{Pos: j.m.Position(j.next), Snap: true}
		if j.next < 0 {
			j.band.Play(SnapCountIn.Strokes(c), j.next, j.m)
			continue
		}
		c.Key = j.bass.key(j.next)
		j.band.Play(a.Strokes(c), j.next, j.m)
	}
}

// Stop releases the bass.
func (j *Jam) Stop(now time.Time) {
	j.band.Stop(now)
}

// A bassist walks the chart chorus after chorus, a new line drawn at
// each one. Each chorus is walked on its own: the line does not carry
// over the double bar yet.
type bassist struct {
	beats []mark.Beat // one chorus
	line  []int       // the chorus playing: a key per beat
	rng   *rand.Rand
}

func newBassist(grid analysis.Changes) *bassist {
	seed := uint64(time.Now().UnixNano())
	return &bassist{
		beats: mark.Expect(grid, tempo.NewMetronome(time.Time{}, 120, mark.BeatsPerBar), 1), // the metronome only numbers the beats
		rng:   rand.New(rand.NewPCG(seed, seed>>32|1)),
	}
}

// key returns the key of beat `n`, counted from 0 at the first chorus.
func (b *bassist) key(n int) int {
	i := n % len(b.beats)
	if i == 0 || b.line == nil {
		b.line = bass.Walk(b.beats, b.rng)
	}
	return b.line[i]
}
