package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
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

// setup is the settings of the grid, a panel over the game at rest:
// the grid, its tempo, the mode, the demo. The arrows, or a tap left or
// right of a value's middle, change it at once; Escape, Close, or a tap
// beside the panel, goes back to the chart. The tempo the player sets holds for the session, until
// the grid changes: each grid comes with its own (see grids.Tune).
type setup struct {
	g     *game
	menu  menu
	panel hit // the panel's box, as last drawn: a tap outside closes it
}

func newSetup(g *game) *setup { return &setup{g: g, menu: menu{items: setItems}} }

func (s *setup) Enter() {}
func (s *setup) Leave() {}

func (s *setup) Update() scene.Transition {
	g := s.g
	g.drain(nil)
	s.menu.move()
	if x, y, ok := g.pressed(); ok && !s.panel.has(x, y) {
		return scene.Pop // a tap beside the panel closes it
	}
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

// The panel of the settings: its items `panelGap` apart, `panelPad`
// inside its frame, the heading and the keys inside it too.
const (
	panelGap = 30
	panelPad = 18
)

// Draw draws the settings as a panel over the game, which shows under
// a veil: the director draws the stack from the bottom (see
// scene.Director), so the game has drawn itself already.
func (s *setup) Draw(dst *ebiten.Image) {
	g := s.g
	c := screen.Canvas{Dst: dst, Scale: g.scale}
	veil := color.NRGBAModel.Convert(paper).(color.NRGBA)
	veil.A = 0xb0
	c.Rect(0, 0, screenWidth, screenHeight, veil)

	labels := s.labels()
	heading, keys := g.lang.T(msgSetTitle), g.lang.T(msgSetKeys)
	w, _ := c.Measure(heading, g.fonts.chord)
	kw, kh := c.Measure(keys, g.fonts.ui)
	w = max(w, kw)
	for _, l := range labels {
		lw, _ := c.Measure(l, g.fonts.chordSmall)
		w = max(w, lw)
	}
	pw := w + 2*panelPad
	ph := float64(panelPad+panelGap+len(labels)*panelGap+panelPad/2) + kh + panelPad
	px, py := (screenWidth-pw)/2, (screenHeight-ph)/2
	s.panel = hit{x: px, y: py, w: pw, h: ph}
	c.Rect(float32(px), float32(py), float32(pw), float32(ph), paper)
	c.RoundRect(float32(px), float32(py), float32(pw), float32(ph), 6, 1.5, ink)

	mid := px + pw/2
	y := py + panelPad
	c.Centred(heading, g.fonts.chord, mid, y, ink)
	y += panelGap + 6
	g.hits = g.hits[:0]
	for i, l := range labels {
		lw, lh := c.Measure(l, g.fonts.chordSmall)
		box := hit{x: px + panelPad/2, y: y - menuPad, w: pw - panelPad, h: lh + 2*menuPad}
		g.hits = append(g.hits, box)
		var col color.Color = faint
		if i == s.menu.chosen {
			c.Rect(float32(mid-lw/2-menuPad), float32(y-menuPad), float32(lw+2*menuPad), float32(lh+2*menuPad), pale)
			col = ink
		}
		c.Centred(l, g.fonts.chordSmall, mid, y, col)
		y += panelGap
	}
	c.Centred(keys, g.fonts.ui, mid, py+ph-panelPad/2-kh, faint)
}

// labels says each item as it stands.
func (s *setup) labels() []string {
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
	labels[setBack] = g.lang.T(msgSetClose)
	return labels
}
