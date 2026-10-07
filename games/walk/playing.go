package main

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/games/walk/band"
)

// The walker's phrases in a lesson: a note a beat, or a bar each,
// snapping on 2 and 4.

// Play plays `keys` from now, scheduled ahead as the band is; after a
// wrong note, from the next bar (see Miss).
func (l *lesson) Play(keys []int) { l.playEvery(keys, phraseNote) }

// PlayBars plays `keys` a bar each, as the roots of a line of a grid.
func (l *lesson) PlayBars(keys []int) { l.playEvery(keys, beatsPerBar*phraseNote) }

// playEvery plays `keys`, a note every `step`, as Play says.
func (l *lesson) playEvery(keys []int, step time.Duration) {
	if len(l.corrections) > 0 {
		l.corrections = append(l.corrections, func() { l.play(keys, step, l.correctAt) })
		return
	}
	l.react(func() { l.play(keys, step, time.Now().Add(band.Lookahead)) })
}

// play plays `keys` from `now`, a note every `step`, each held almost to
// the next. A note a bar, he snaps on 2 and 4.
func (l *lesson) play(keys []int, step time.Duration, now time.Time) {
	snaps := step == beatsPerBar*phraseNote
	if snaps {
		for i := range keys {
			bar := now.Add(time.Duration(i) * step)
			l.band.SnapAt(grooveVel, bar.Add(phraseNote))
			l.band.SnapAt(grooveVel, bar.Add(3*phraseNote))
		}
	}
	for i, k := range keys {
		at := now.Add(time.Duration(i) * step)
		off := at.Add(time.Duration(phraseHold * float64(step)))
		if l.bass {
			l.band.Pluck(k, phraseVel, at, off)
			continue
		}
		l.band.Piano.ScheduleOn(k, phraseVel, at)
		l.band.Piano.ScheduleOff(k, off)
	}
	l.phrase, l.phraseAt, l.phraseStep, l.phraseSnap, l.playing = keys, now, step, snaps, true
}

// snapping tells whether the walker's hand is up, snapping, at `now`:
// for a right answer, or on 2 and 4 under his line of a grid.
func (l *lesson) snapping(now time.Time) bool {
	if now.Sub(l.cheerAt) < cheerTime {
		return true
	}
	if !l.playing || !l.phraseSnap || now.Before(l.phraseAt) {
		return false
	}
	t := now.Sub(l.phraseAt)
	beat := int(t / phraseNote)
	return beat%2 == 1 && t%phraseNote < cheerTime
}
