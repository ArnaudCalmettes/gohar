package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/scene"
)

// The tempo a player sets to practise a grid: in steps, between two
// bounds a walking bass can hold.
const (
	bpmStep = 5
	minBPM  = 60
	maxBPM  = 240
)

// The items of the settings of a grid, in order.
const (
	setGrid = iota
	setTempo
	setMode
	setDemo
	setBack
	setItems
)

// setup is the settings of the grid, over the game at rest: the grid,
// its tempo, the mode, the demo. The arrows, or a tap left or right of
// a value's middle, change it at once; Escape, or Back, goes back to
// the chart. The tempo the player sets holds for the session, until
// the grid changes: each grid comes with its own (see grids.Tune).
type setup struct {
	g    *game
	menu menu
}

func newSetup(g *game) *setup { return &setup{g: g, menu: menu{items: setItems}} }

func (s *setup) Enter() {}
func (s *setup) Leave() {}

func (s *setup) Update() scene.Transition {
	g := s.g
	g.drain(nil)
	s.menu.move()
	tp, tapped := g.tapped()
	if tapped {
		s.menu.chosen = tp.item
	}
	step := 0
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyTab):
		return scene.Pop
	case s.menu.chosen == setBack && (confirmed() || tapped):
		return scene.Pop
	case tapped:
		step = tp.side
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
		step = -1
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight), confirmed():
		step = 1
	}
	if step != 0 {
		s.change(step)
	}
	return scene.Stay
}

// change moves the value chosen one `step` along: the grids go round,
// the tempo stops at its bounds, the mode and the demo switch.
func (s *setup) change(step int) {
	g := s.g
	switch s.menu.chosen {
	case setGrid:
		g.load((g.current + step + len(g.tunes)) % len(g.tunes))
	case setTempo:
		g.bpm = min(max(g.bpm+float64(step*bpmStep), minBPM), maxBPM)
	case setMode:
		g.practicing = !g.practicing
	case setDemo:
		g.band.Demo = !g.band.Demo
	}
}

func (s *setup) Draw(dst *ebiten.Image) {
	g := s.g
	labels := make([]string, setItems)
	labels[setGrid] = g.lang.T(msgSetGrid, "Title", g.title)
	labels[setTempo] = g.tempoWords()
	labels[setMode] = g.lang.T(msgSetModeTempo)
	if g.practicing {
		labels[setMode] = g.lang.T(msgSetModeFree)
	}
	labels[setDemo] = g.lang.T(msgSetDemoOff)
	if g.band.Demo {
		labels[setDemo] = g.lang.T(msgSetDemoOn)
	}
	labels[setBack] = g.lang.T(msgOptBack)
	g.drawList(dst, nil, 0, g.lang.T(msgSetTitle), g.fonts.heading, entries(labels, nil), s.menu.chosen, g.lang.T(msgSetKeys), layout{top: menuY - 30, gap: menuGap, apart: backApart})
}
