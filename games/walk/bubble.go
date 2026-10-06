package main

import (
	"math"
	"strings"

	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// The walker's bubbles, in the game and in the lessons, drawn the way a
// comic strip drawn in strokes does: the phrase in the chart's hand, a
// line under it, and a tail from that line toward his head. No frame.
const (
	bubbleGap = 3 // between the phrase and its line
	tailGap   = 4 // between the tail's end and his head
)

// drawSpeech draws `lines` from (`left`, `top`), `lineH` apart, the line
// under them `width` long, and the tail toward the walker, his feet at
// `x`. The tail leaves the line at the point nearest his head, straight
// above it when he stands under the line, and aims at the centre of his
// head, wherever he walks; it stops a few pixels short of it.
func (a *app) drawSpeech(c screen.Canvas, lines []string, left, top, width, lineH, x float64) {
	y := top
	for _, line := range lines {
		c.Text(line, a.fonts.bubble, left, y, ink)
		y += lineH
	}
	under := y + bubbleGap
	c.Line(float32(left), float32(under), float32(left+width), float32(under), 1, ink)

	hx, hy := x, float64(walkerY-sketchHeadAbove*walkerScale)
	fx := min(max(hx, left), left+width)
	dx, dy := hx-fx, hy-under
	d := math.Hypot(dx, dy)
	stop := float64(sketchHead*walkerScale + tailGap)
	if dy <= 0 || d <= stop {
		return // his head level with the line: no tail
	}
	k := 1 - stop/d
	c.Line(float32(fx), float32(under), float32(fx+dx*k), float32(under+dy*k), 1, ink)
}

// wrap cuts `s` into lines no wider than `width` in the bubble's hand,
// between words.
func (a *app) wrap(s string, width float64) []string {
	measure := func(s string) float64 {
		w, _ := screen.Canvas{Scale: 1}.Measure(s, a.fonts.bubble)
		return w
	}
	var lines []string
	line := ""
	for _, word := range strings.Fields(s) {
		next := word
		if line != "" {
			next = line + " " + word
		}
		if line != "" && measure(next) > width {
			lines = append(lines, line)
			next = word
		}
		line = next
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}
