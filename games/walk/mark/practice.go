package mark

// Practice is the phase without tempo (see "La boucle de jeu" in
// docs/walk.md): the chart waits for the player. The cursor stays on
// the next arrival, or beat 1 of a bar where the chord carries on,
// until a note lands it, then moves on; another note is marked and the
// cursor stays. No drums, no metronome: the marks of time disappear,
// those of pitch remain.
type Practice struct {
	rules    Rules
	beats    []Beat // one chorus
	arrivals []int  // the beats that wait for a note: arrivals, held bars
	at       int    // the arrival waited for, in arrivals
	landed   map[int]bool
}

// NewPractice waits, by the rules `r`, on the arrivals of `beats`, one
// chorus.
func NewPractice(r Rules, beats []Beat) *Practice {
	p := &Practice{rules: r, beats: beats, landed: map[int]bool{}}
	for i, b := range beats {
		if (b.Arrives || b.Holds) && !b.Chord.Silent {
			p.arrivals = append(p.arrivals, i)
		}
	}
	return p
}

// Waiting returns the beat waited for.
func (p *Practice) Waiting() Beat {
	return p.beats[p.arrivals[p.at]]
}

// Arrival is the rank of the arrival waited for, among those of the
// chorus: it moves on when a note lands it.
func (p *Practice) Arrival() int { return p.at }

// Landed tells whether beat `n` of the chorus is landed, this time
// round.
func (p *Practice) Landed(n int) bool { return p.landed[n] }

// Play marks a note of the bass against the beat waited for, and moves
// on when it lands it, by the same rules as the marker. The chart
// loops: after the last beat, the first one again, its marks cleared.
func (p *Practice) Play(key int) Pitch {
	b := p.Waiting()
	pitch := pitchOf(key, b.Chord)
	if p.rules.lands(b, pitch) {
		p.landed[b.N] = true
		p.at++
		if p.at == len(p.arrivals) {
			p.at = 0
			clear(p.landed)
		}
	}
	return pitch
}
