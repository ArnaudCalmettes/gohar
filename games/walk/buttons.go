package main

import (
	"fmt"
	"time"

	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
)

// The buttons of the game, for a mouse or a finger (see "Les boutons"
// in docs/walk.md): back to the title, left of the title; the settings
// of the grid, the block of its title and tempo; and play or stop, the
// label of the mode, top right. The rest of the screen is the chart,
// which scrolls. In a palier, back alone.
const (
	buttonBack = iota
	buttonSetup
	buttonPlay
)

// backW is the room of the back button, left of the title: a chevron.
const backW = 18

// drawButtons draws the title block and the mode, and keeps the boxes
// of the buttons for pressed.
func (g *game) drawButtons(c screen.Canvas) {
	x := float64(margin + backW)
	_, h := c.Measure(g.title, g.fonts.ui)
	drawChevronLeft(c, margin+backW/3, float32(titleY+h/2))
	back := hit{x: 0, y: 0, w: x, h: tempoY + h + tapPad}

	c.Text(g.title, g.fonts.ui, x, titleY, ink)
	c.Text(g.pace(), g.fonts.ui, x, tempoY, faint)
	if g.back != nil {
		g.hits = append(g.hits[:0], back)
		return
	}
	w := 0.0
	for _, s := range []string{g.title, g.pace()} {
		sw, _ := c.Measure(s, g.fonts.ui)
		w = max(w, sw)
	}
	setup := hit{x: x - menuPad, y: titleY - menuPad, w: w + 2*menuPad, h: tempoY + h - titleY + 2*menuPad}
	if !g.running {
		// Underlined, faint: the block opens the settings.
		c.Rect(float32(setup.x), float32(setup.y+setup.h), float32(setup.w), 1, faint)
	}
	g.hits = append(g.hits[:0], back, setup, g.drawMode(c))
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

// pressed does what a button pressed does, the boxes those of the last
// frame: back, the title or the course; the settings, at rest; play or
// stop, as Space.
func (g *game) pressed(now time.Time) scene.Transition {
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

// drawChevronLeft draws a chevron pointing left, centred on (`x`, `y`):
// the back button.
func drawChevronLeft(c screen.Canvas, x, y float32) {
	const w = 7
	c.Line(x+w/2, y-w, x-w/2, y, 2, ink)
	c.Line(x-w/2, y, x+w/2, y+w, 2, ink)
}
