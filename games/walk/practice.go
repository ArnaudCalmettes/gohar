package main

// practice is the phase without tempo (see "La boucle de jeu" in
// docs/walk.md): the chart waits for the player. The cursor stays on
// the next arrival, or beat 1 of a bar where the chord carries on,
// until a note lands it, then moves on; another note is marked and the
// cursor stays. No drums, no metronome: the marks of time disappear,
// those of pitch remain.
type practice struct {
	rules    Rules
	beats    []Beat // one chorus
	arrivals []int  // the beats that wait for a note: arrivals, held bars
	at       int    // the arrival waited for, in arrivals
	landed   map[int]bool
}

func newPractice(r Rules, beats []Beat) *practice {
	p := &practice{rules: r, beats: beats, landed: map[int]bool{}}
	for i, b := range beats {
		if (b.Arrives || b.Holds) && !b.Chord.Silent {
			p.arrivals = append(p.arrivals, i)
		}
	}
	return p
}

// waiting returns the beat waited for.
func (p *practice) waiting() Beat {
	return p.beats[p.arrivals[p.at]]
}

// play marks a note of the bass against the beat waited for, and moves
// on when it lands it, by the same rules as the marker. The chart
// loops: after the last beat, the first one again, its marks cleared.
func (p *practice) play(key int) Pitch {
	b := p.waiting()
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
