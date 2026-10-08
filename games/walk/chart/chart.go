package chart

import (
	"image/color"

	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/grids"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// A Cell is one chord written in a bar, at the beat it starts on.
type Cell struct {
	Symbol Symbol
	Beat   int // in the bar, from 0
}

// Bars writes the chords of `t` bar by bar, as the grid spells them:
// one tonality is not enough to spell a grid that modulates, the C of
// Tune Up's Cmaj7 would be a B♯ in D. A bar where no chord starts is
// left empty, and drawn as the repeat sign.
//
// Respelling by the tonalities the analysis hears, zone by zone, as
// charts/cmd/analyse does, is for later (see docs/chantiers.md).
func Bars(t grids.Tune) [][]Cell {
	bars := make([][]Cell, len(t.Grid.Bars))
	for i, ch := range t.Grid.Chords {
		bar := t.Grid.Bar(ch.Start)
		if bar < 0 {
			continue
		}
		beat := int((ch.Start - t.Grid.Bars[bar]) / analysis.TicksPerBeat)
		bars[bar] = append(bars[bar], Cell{Symbol: SymbolOf(t.Written[i]), Beat: beat})
	}
	return bars
}

// Small tells whether a bar of `bars` holds two chords: then all of
// them are written small, one size for the whole chart.
func Small(bars [][]Cell) bool {
	for _, cells := range bars {
		if len(cells) > 1 {
			return true
		}
	}
	return false
}

// Within a bar: a chord from the start of its beat, and how far an
// exponent rises above the line's top.
const (
	chordDX, chordDY = 8, 14
	raisedDY         = 3
)

// A Pen is what the chart is written with: the chart's hand in two
// sizes, each with the smaller one of its exponents, and the ink.
type Pen struct {
	Chord, ChordRaised *screen.Font
	Small, SmallRaised *screen.Font // two chords in a bar
	Ink, Paper         color.Color  // the paper fills a rehearsal mark
}

// DrawBar draws the bar `cells` from (`x`, `y`), `w` wide and `h`
// high: its `fill`, if not nil, its barline, the closing one too when
// it ends its row (`last`), and its chords, or the repeat sign for a
// bar where none starts. Two chords sharing the bar are written small,
// and all of them when `small` (see Small).
func (p Pen) DrawBar(c screen.Canvas, cells []Cell, x, y, w, h float32, fill color.Color, small, last bool) {
	if fill != nil {
		c.Rect(x, y, w, h, fill)
	}
	c.Line(x, y, x, y+h, 1, p.Ink)
	if last {
		c.Line(x+w, y, x+w, y+h, 1, p.Ink)
	}
	dy := float64(y) + min(chordDY, float64(h)/2-10) // about the middle of a shrunk bar
	font, raised := p.Chord, p.ChordRaised
	if small || len(cells) > 1 {
		font, raised = p.Small, p.SmallRaised
	}
	if len(cells) == 0 {
		c.Centred("%", font, float64(x+w/2), dy, p.Ink)
	}
	for _, cl := range cells {
		cx := float64(x) + chordDX + float64(cl.Beat)*float64(w)/mark.BeatsPerBar
		p.drawSymbol(c, cl.Symbol, font, raised, cx, dy)
	}
}

// rehearsalPad is the room around a rehearsal mark, inside its box.
const rehearsalPad = 3

// DrawRehearsal draws the rehearsal mark `label` framed, in the small
// hand, its box ending on the barline at `x`, from the top `y` of the
// bar: left of the bar it marks, as a Real Book prints it.
func (p Pen) DrawRehearsal(c screen.Canvas, label string, x, y float32) {
	w, h := c.Measure(label, p.Small)
	bw, bh := float32(w)+2*rehearsalPad, float32(h)+2*rehearsalPad
	bx := x - bw - rehearsalPad
	c.Rect(bx, y, bw, bh, p.Ink)
	c.Rect(bx+1, y+1, bw-2, bh-2, p.Paper)
	c.Text(label, p.Small, float64(bx)+rehearsalPad, float64(y)+rehearsalPad, p.Ink)
}

// DrawSymbol draws the symbol `s` alone from (`x`, `y`), in the large
// hand, and returns where it ends: a chord shown outside a grid.
func (p Pen) DrawSymbol(c screen.Canvas, s Symbol, x, y float64) float64 {
	return p.drawSymbol(c, s, p.Chord, p.ChordRaised, x, y)
}

// drawSymbol draws the symbol `s` from (`x`, `y`): its line in `font`,
// its raised run in the smaller `raised`, a little above, as an
// exponent. It returns where the symbol ends.
func (p Pen) drawSymbol(c screen.Canvas, s Symbol, font, raised *screen.Font, x, y float64) float64 {
	c.Text(s.Line, font, x, y, p.Ink)
	w, _ := c.Measure(s.Line, font)
	x += w
	if s.Raised != "" {
		c.Text(s.Raised, raised, x, y-raisedDY, p.Ink)
		w, _ = c.Measure(s.Raised, raised)
		x += w
	}
	c.Text(s.Bass, font, x, y, p.Ink)
	w, _ = c.Measure(s.Bass, font)
	return x + w
}

// The marks of the arrivals: landed, missed, doubled.
var (
	LandedInk  = color.RGBA{0x2a, 0x9d, 0x5a, 0xff}
	MissedInk  = color.RGBA{0xc0, 0x39, 0x2b, 0xff}
	DoubledInk = color.RGBA{0xe0, 0x8a, 0x1e, 0xff}
)

// DrawMark draws the mark of an arrival centred on `x`, `y`: a green
// dot when landed, a red cross when missed, two orange strokes when
// doubled.
func DrawMark(c screen.Canvas, k mark.BeatKind, x, y float32) {
	const (
		r      = 3.5 // the dot, the arms of the cross
		apart  = 2   // the two strokes, from the centre
		half   = 4   // their half length
		stroke = 1.5
	)
	switch k {
	case mark.Landed:
		c.Circle(x, y, r, LandedInk)
	case mark.Missed:
		c.Line(x-r, y-r, x+r, y+r, stroke, MissedInk)
		c.Line(x-r, y+r, x+r, y-r, stroke, MissedInk)
	case mark.Doubled:
		c.Line(x-apart, y-half, x-apart, y+half, stroke, DoubledInk)
		c.Line(x+apart, y-half, x+apart, y+half, stroke, DoubledInk)
	}
}
