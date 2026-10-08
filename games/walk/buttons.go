package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// The buttons of the game, for a mouse or a finger (see "Les boutons"
// in docs/walk.md), each a framed square with a sign everyone knows:
// back, a chevron, top left; the settings of the grid, a cog, right of
// its title; play, a triangle, or stop, a square, in the middle, the
// main command of the screen. The mode stays top right, a label. The
// rest of the screen is the chart, which scrolls. In a palier, back
// alone.
const (
	buttonBack = iota
	buttonSetup
	buttonPlay
)

// The buttons' squares: `buttonSize` wide, from `buttonY`; play is
// bigger, `playSize`.
const (
	buttonSize = 22
	playSize   = 28
	buttonY    = 14
	buttonGap  = 8 // between a button and the text beside it
)

// button is the box of a square button `size` wide from (`x`, `y`),
// widened for a finger.
func button(x, y, size float64) hit {
	return hit{x: x - tapPad/2, y: y - tapPad/2, w: size + tapPad, h: size + tapPad}
}

// drawFrame draws the frame of a button, in `col`.
func drawFrame(c screen.Canvas, x, y, size float64, col color.Color) {
	c.RoundRect(float32(x), float32(y), float32(size), float32(size), 4, 1.5, col)
}

// drawButtons draws the buttons and the title block, and keeps the
// boxes of the buttons for clicked.
func (g *game) drawButtons(c screen.Canvas) {
	drawBack(c, margin, buttonY)
	back := button(margin, buttonY, buttonSize)

	x := float64(margin + buttonSize + buttonGap)
	c.Text(g.title, g.fonts.ui, x, buttonY, ink)
	c.Text(g.pace(), g.fonts.ui, x, buttonY+16, faint)
	g.drawMode(c)
	if g.back != nil {
		g.hits = append(g.hits[:0], back)
		return
	}

	// The cog, right of the title or the tempo, whichever is wider;
	// faint while the band plays, when the settings do not open.
	w := 0.0
	for _, s := range []string{g.title, g.pace()} {
		sw, _ := c.Measure(s, g.fonts.ui)
		w = max(w, sw)
	}
	cx := x + w + buttonGap
	var col color.Color = ink
	if g.running {
		col = faint
	}
	drawCog(c, cx, buttonY, col)

	px := float64(screenWidth/2 - playSize/2)
	drawPlay(c, px, buttonY-3, g.running)
	g.hits = append(g.hits[:0], back, button(cx, buttonY, buttonSize), button(px, buttonY-3, playSize))
}

// drawBack draws the back button from (`x`, `y`): a chevron pointing
// left.
func drawBack(c screen.Canvas, x, y float64) {
	drawFrame(c, x, y, buttonSize, ink)
	mx, my := float32(x+buttonSize/2), float32(y+buttonSize/2)
	const w = 4
	c.Line(mx+w/2, my-w, mx-w/2, my, 2, ink)
	c.Line(mx-w/2, my, mx+w/2, my+w, 2, ink)
}

// drawCog draws the settings button from (`x`, `y`): a cog, its teeth
// round a disc, a hole in the middle.
func drawCog(c screen.Canvas, x, y float64, col color.Color) {
	drawFrame(c, x, y, buttonSize, col)
	mx, my := float32(x+buttonSize/2), float32(y+buttonSize/2)
	const (
		teeth = 8
		r     = 5.5 // the disc
		tooth = 2.5 // how far a tooth stands out
	)
	for i := range teeth {
		a := float64(i) * 2 * math.Pi / teeth
		dx, dy := float32(math.Cos(a)), float32(math.Sin(a))
		c.Line(mx+dx*(r-1), my+dy*(r-1), mx+dx*(r+tooth), my+dy*(r+tooth), 2.5, col)
	}
	c.Circle(mx, my, r, col)
	c.Circle(mx, my, r/2.2, paper)
}

// drawPlay draws the play button from (`x`, `y`): a triangle pointing
// right, or, while the band plays, the square of stop.
func drawPlay(c screen.Canvas, x, y float64, playing bool) {
	drawFrame(c, x, y, playSize, ink)
	mx, my := float32(x+playSize/2), float32(y+playSize/2)
	if playing {
		const half = 5
		c.Rect(mx-half, my-half, 2*half, 2*half, ink)
		return
	}
	const h = 7 // half the triangle's height
	c.Polygon([]float32{mx - h*0.7, my - h, mx + h, my, mx - h*0.7, my + h}, ink)
}

// pace says the tempo as the game stands: none without tempo.
func (g *game) pace() string {
	if g.practicing {
		return g.lang.T(msgFreeTempo)
	}
	return g.tempoWords()
}

// tempoWords says the tempo: the grid's, or the one the player chose,
// with the grid's beside it.
func (g *game) tempoWords() string {
	bpm := fmt.Sprintf("%.0f", g.bpm)
	if g.bpm != g.tuneBPM {
		return g.lang.T(msgTempoAside, "BPM", bpm, "Grid", fmt.Sprintf("%.0f", g.tuneBPM))
	}
	return g.lang.T(msgTempo, "BPM", bpm)
}

// clicked does what a button clicked does, the boxes those of the last
// frame: back, the title or the course; the settings, at rest; play or
// stop, as Space.
func (g *game) clicked(now time.Time) scene.Transition {
	tp, ok := g.tapped()
	if !ok {
		return scene.Stay
	}
	switch {
	case tp.item == buttonBack && g.back != nil:
		return scene.Replace(g.back)
	case tp.item == buttonBack:
		return scene.Replace(newTitle(g.app))
	case g.back != nil:
	case tp.item == buttonSetup && !g.running:
		return scene.Push(newSetup(g))
	case tp.item == buttonPlay && g.running:
		g.stop(now)
	case tp.item == buttonPlay:
		g.start(now)
	}
	return scene.Stay
}
