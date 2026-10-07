package figure

import (
	"image/color"
	"math"
	"strings"

	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// The walker's bubbles, in the game and in the lessons, drawn the way a
// comic strip drawn in strokes does: the phrase in the chart's hand, a
// line under it, and a tail from that line toward his head. No frame.
const (
	BubbleGap = 3 // between the phrase and its line
	tailGap   = 4 // between the tail's end and his head
)

// A Bubble is where the walker's phrase is written, and in what: its
// lines from (`Left`, `Top`), `LineH` apart, under them a line `Width`
// long, all in `Font` and `Ink`.
type Bubble struct {
	Font         *screen.Font
	Ink          color.Color
	Left, Top    float64
	Width, LineH float64
}

// Draw draws `lines` in the bubble, and the tail toward the walker's
// head, its centre on (`hx`, `hy`), `r` its radius (see Pose.Head). The
// tail leaves the line at the point nearest his head, straight above
// it when he stands under the line, and aims at the centre of his head,
// wherever he walks; it stops a few pixels short of it. A bubble of
// several lines may end below his head: the tail then goes up to it.
func (b Bubble) Draw(c screen.Canvas, lines []string, hx, hy, r float64) {
	ly := b.Top
	for _, line := range lines {
		c.Text(line, b.Font, b.Left, ly, b.Ink)
		ly += b.LineH
	}
	under := ly + BubbleGap
	c.Line(float32(b.Left), float32(under), float32(b.Left+b.Width), float32(under), 1, b.Ink)

	fx := min(max(hx, b.Left), b.Left+b.Width)
	dx, dy := hx-fx, hy-under
	d := math.Hypot(dx, dy)
	stop := r + tailGap
	if d <= stop {
		return // his head against the line: no room for a tail
	}
	k := 1 - stop/d
	c.Line(float32(fx), float32(under), float32(fx+dx*k), float32(under+dy*k), 1, b.Ink)
}

// Wrap cuts `s` into lines no wider than the bubble, between words.
func (b Bubble) Wrap(s string) []string {
	measure := func(s string) float64 {
		w, _ := screen.Canvas{Scale: 1}.Measure(s, b.Font)
		return w
	}
	var lines []string
	line := ""
	for _, word := range strings.Fields(s) {
		next := word
		if line != "" {
			next = line + " " + word
		}
		if line != "" && measure(next) > b.Width {
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
