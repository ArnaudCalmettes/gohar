package main

import (
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/chart"
	"github.com/ArnaudCalmettes/gohar/games/walk/lessons"
)

// What the walker writes in a lesson, chords or a line of a grid, and
// the bass the player plays it on (see lessons.Stage).

func (l *lesson) Chords(chords []string) { l.chords = symbols(chords) }

func (l *lesson) Line(chords []string, bar int) {
	l.line, l.lineBar = nil, bar
	for _, ch := range chords {
		if ch == lessons.RepeatBar {
			l.line = append(l.line, nil) // drawn as the repeat sign
			continue
		}
		// Two chords in a bar, "C6 A7", share it: the second on beat 3.
		var cells []chart.Cell
		for i, s := range symbols(strings.Fields(ch)) {
			cells = append(cells, chart.Cell{Symbol: s, Beat: i * 2})
		}
		l.line = append(l.line, cells)
	}
}

func (l *lesson) Bass(on bool) {
	l.bass = on
	l.band.BassHand.Store(on)
	l.piano.Low, l.piano.High = lessonLow, lessonHigh
	if on {
		l.piano.Low, l.piano.High = bassLow, bassHigh
	}
}

// symbols reads `chords`, in ChordPro, as a chart draws them. A lesson
// writes them in its script: one it cannot read is a bug, and logged.
func symbols(chords []string) []chart.Symbol {
	var out []chart.Symbol
	for _, s := range chords {
		c, err := chordpro.ReadChord(s)
		if err != nil {
			log.Println("lesson:", err)
			continue
		}
		out = append(out, chart.SymbolOf(c))
	}
	return out
}

// drawWritten draws what the walker wrote: chords side by side, or a
// line of a grid, its bar played shaded; while he plays the line, the
// bar of the note he plays. A grid longer than a row shows the row of
// the bar played, the first one before.
func (l *lesson) drawWritten(c screen.Canvas, now time.Time) {
	pen := l.pen()
	x := float64(bubbleX)
	for _, s := range l.chords {
		x = pen.DrawSymbol(c, s, x, writtenY) + chordsGap
	}
	bar := l.lineBar
	if l.playing && !now.Before(l.phraseAt) {
		bar = int(now.Sub(l.phraseAt) / l.phraseStep)
	}
	row := max(bar, 0) / barsPerRow
	for i, cells := range l.line {
		if i/barsPerRow != row {
			continue
		}
		var fill color.Color
		if i == bar {
			fill = pale
		}
		x := float32(bubbleX + i%barsPerRow*barW)
		pen.DrawBar(c, cells, x, writtenY, barW, lineBarH, fill, false, i == len(l.line)-1 || i%barsPerRow == barsPerRow-1)
	}
}
