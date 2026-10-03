package main

import (
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// A Beat is what the grid says about one beat: the first brick of the
// game (see "Les briques" in docs/walk.md). What the player should play
// on it is the marker's business, palier by palier; this is only the
// chart, laid out in time.
//
// The scale of the moment, which the walking bass will need, comes
// later, from the analysis.
type Beat struct {
	N        int // the metronome's numbering: 0 is bar 1, beat 1
	Position Position
	Strong   bool

	// Chord is the change sounding at the start of the beat, Next the
	// one after it. Arrives says the chord starts on this beat: where
	// the root is expected, from the first palier on.
	Chord   analysis.Change
	Next    analysis.Change
	Arrives bool

	// Holds is beat 1 of a bar where no chord arrives: the chord of the
	// bar before carries on, like two tied whole notes on a chart. A
	// bassist still plays a note there, and the marker expects one.
	Holds bool

	// PerBar is how many changes start in the bar, the harmonic rhythm:
	// it picks which of Siskind's rules applies (one, two or four
	// chords a bar).
	PerBar int
}

// Expect lays `choruses` choruses of `c` out beat by beat, for the
// metronome `m`. A chart without bars is cut into bars of the
// metronome's length. Changes that start off the beat count from the
// beat they start in.
//
// The coda is not played yet: every chorus is the same.
func Expect(c analysis.Changes, m Metronome, choruses int) []Beat {
	if len(c.Chords) == 0 {
		return nil
	}
	last := c.Chords[len(c.Chords)-1]
	chorusBeats := int((last.Start + last.Length + analysis.TicksPerBeat - 1) / analysis.TicksPerBeat)
	barBeats := m.perBar

	beats := make([]Beat, 0, choruses*chorusBeats)
	i := 0 // the change sounding
	for chorus := range choruses {
		i = 0
		for b := range chorusBeats {
			from := analysis.Ticks(b) * analysis.TicksPerBeat
			to := from + analysis.TicksPerBeat
			arrives := false
			for j := i; j < len(c.Chords) && c.Chords[j].Start < to; j++ {
				if c.Chords[j].Start >= from {
					arrives = true
				}
				i = j
			}
			next := c.Next(i)
			if next < 0 {
				next = i
			}
			n := chorus*chorusBeats + b
			p := m.Position(n)
			beats = append(beats, Beat{
				N:        n,
				Position: p,
				Strong:   m.Strong(p),
				Chord:    c.Chords[i],
				Next:     c.Chords[next],
				Arrives:  arrives,
				Holds:    p.Beat == 1 && !arrives,
				PerBar:   changesInBar(c, b/barBeats, barBeats),
			})
		}
	}
	return beats
}

// changesInBar counts the changes that start in bar `bar` of the chart,
// counted from 0: its bars when it has them, else bars of `barBeats`.
func changesInBar(c analysis.Changes, bar, barBeats int) int {
	var from, to analysis.Ticks
	if bar < len(c.Bars) {
		from = c.Bars[bar]
		to = analysis.Ticks(1 << 62)
		if bar+1 < len(c.Bars) {
			to = c.Bars[bar+1]
		}
	} else {
		from = analysis.Ticks(bar*barBeats) * analysis.TicksPerBeat
		to = from + analysis.Ticks(barBeats)*analysis.TicksPerBeat
	}
	n := 0
	for _, ch := range c.Chords {
		if ch.Start >= from && ch.Start < to {
			n++
		}
	}
	return n
}

// bassKey places a pitch class in the register of the double bass, as a
// MIDI key: from G♯1 to G2, so that F, B♭ and C land on F2, B♭1 and C2,
// within the four strings Siskind asks for (p. 116, 204).
func bassKey(pc harmony.PitchClass) int {
	k := 36 + int(pc)
	if k > 43 {
		k -= 12
	}
	return k
}
