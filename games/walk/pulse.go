package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/band"
	"github.com/ArnaudCalmettes/gohar/games/walk/chart"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// The pulse of a lesson, and the count of the bar that follows it (see
// lessons.Stage): chapter 2, where the band comes in.

// The count of the bar: its digits, `countGap` apart, centred under the
// bubble, where the lessons write; the marks of the beats `countMarkDY`
// under their tops.
const (
	countGap    = 70
	countMarkDY = 80
)

// lastOf returns the beat marked at position `beat` of the bar within
// the last bar before `x`, the pulse's beats now, if any.
func (l *lesson) lastOf(beat int, x float64) (int, bool) {
	for n := range l.beatMarks {
		if n%l.pulse.Metronome().PerBar() == beat-1 && float64(n) > x-float64(l.pulse.Metronome().PerBar()) && float64(n) <= x {
			return n, true
		}
	}
	return 0, false
}

func (l *lesson) Pulse(bpm float64) {
	switch {
	case bpm == 0:
		l.pulse = nil // the beats scheduled ring out
	case l.pulse == nil || l.pulse.BPM() != bpm:
		l.pulse, l.beat = band.NewPulse(l.band, bpm), -1
		clear(l.beatMarks)
	}
}

func (l *lesson) Timing() (int, float64) {
	if l.pulse == nil {
		return 0, 0
	}
	m := l.pulse.Metronome()
	n, off := m.Nearest(l.strikeAt)
	return n, float64(off) / float64(m.Beat())
}

func (l *lesson) MarkBeat(n int, landed bool) {
	if l.beatMarks == nil {
		l.beatMarks = map[int]bool{}
	}
	l.beatMarks[n] = landed
}

func (l *lesson) Count(on bool) { l.counting = on && l.pulse != nil }

// keepPulse queues the beats of the pulse due soon, if it plays, and
// tells the runner each beat as it begins, for the steps in rhythm.
func (l *lesson) keepPulse(now time.Time) {
	if l.pulse == nil {
		return
	}
	l.pulse.Play(now)
	n := int(math.Floor(l.pulse.Metronome().Beats(now)))
	for l.beat < n && l.pulse != nil {
		l.beat++
		if l.beat >= 0 && l.teaching() {
			l.runner.Beat(l.beat)
		}
	}
}

// drawCount draws the count of the bar, "1 2 3 4", the beat playing in
// ink, the others faint, all alike, as the hi-hat plays them. Under
// each digit, the mark of its beat in the last bar, landed or missed,
// as the game's chart marks them: discreet, and gone a bar later.
func (l *lesson) drawCount(c screen.Canvas, now time.Time) {
	if !l.counting || l.pulse == nil {
		return
	}
	m := l.pulse.Metronome()
	x := m.Beats(now)
	playing := 0 // before the first beat, none
	if x >= 0 {
		playing = m.Position(int(math.Floor(x))).Beat
	}
	mid := float64(bubbleX + bubbleW/2)
	for beat := 1; beat <= m.PerBar(); beat++ {
		var col color.Color = faint
		if beat == playing {
			col = ink
		}
		bx := mid + float64(beat-1)*countGap - float64(m.PerBar()-1)*countGap/2
		c.Centred(fmt.Sprint(beat), l.fonts.count, bx, writtenY, col)
		if n, ok := l.lastOf(beat, x); ok {
			kind := mark.Missed
			if l.beatMarks[n] {
				kind = mark.Landed
			}
			chart.DrawMark(c, kind, float32(bx), writtenY+countMarkDY)
		}
	}
}
