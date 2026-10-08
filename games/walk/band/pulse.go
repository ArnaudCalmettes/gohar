package band

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/games/tempo"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// A Pulse is the band of the lessons of chapter 2: the hi-hat on every
// beat, alone, each beat sounding the same, for the player to count
// them without a doubt (see docs/debutants/chapitre-2.md). It starts a
// lookahead from now, without a count-in, and plays on from step to
// step: a lesson never breaks the pulse.
type Pulse struct {
	band *Band
	bpm  float64
	m    tempo.Metronome
	next int // the next beat to schedule
}

// NewPulse starts the hi-hat of `b` at `bpm`.
func NewPulse(b *Band, bpm float64) *Pulse {
	m := tempo.NewMetronome(time.Now().Add(Lookahead), bpm, mark.BeatsPerBar)
	return &Pulse{band: b, bpm: bpm, m: m}
}

// Metronome is the pulse's clock: the count on the screen follows it,
// and the notes of the player are judged against it.
func (p *Pulse) Metronome() tempo.Metronome { return p.m }

// BPM is the pulse's tempo.
func (p *Pulse) BPM() float64 { return p.bpm }

// Play queues the beats due within the lookahead of `now`.
func (p *Pulse) Play(now time.Time) {
	_, end := p.m.Due(now, now.Add(Lookahead))
	for ; p.next < end; p.next++ {
		p.band.Play(Hat.Strokes(Cue{Pos: p.m.Position(p.next)}), p.next, p.m)
	}
}
