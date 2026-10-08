package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// review is the summary of a run, shown once the player has played it
// through, over the game (see "Le bilan" in docs/walk.md): on the left,
// the situations that worked and those to consolidate, and how the notes
// sat on the beat; on the right, the whole chart, each bar tinted as
// the situations are. Any key, on the computer or the piano, goes back
// to the game.
type review struct {
	g *game
	s *mark.Summary
}

func newReview(g *game) *review { return &review{g: g, s: g.summary} }

func (r *review) Enter() {}
func (r *review) Leave() {}

func (r *review) Update() scene.Transition {
	done := len(inpututil.AppendJustPressedKeys(nil)) > 0 ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 // a tap anywhere, as any key
	r.g.drain(func(e keyboard.Event) { done = done || e.Down })
	if done {
		return scene.Pop
	}
	return scene.Stay
}

// The review's layout, in logical units, the chart's row shrunk to fit
// the longest grids.
const (
	reviewLine  = 15 // between two lines of text
	reviewGap   = 10 // between two blocks
	swatch      = 9  // the square of colour beside a heading
	reviewChart = adviceY - reviewGap - chartY
	adviceY     = statusY - reviewLine - reviewGap // the advice, above the keys
)

// The words of a review: the feel, the formulas and the situations to
// work, as the advice names them.
var (
	feelPhrase    = [...]string{mark.FeelSteady: msgReviewSteady, mark.FeelRushing: msgReviewRushing, mark.FeelDragging: msgReviewDragging, mark.FeelUnsteady: msgReviewUnsteady}
	formulaPhrase = [mark.Formulas]string{mark.Anatole: msgFormulaAnatole, mark.ThreeSix: msgFormulaThreeSix, mark.Aeolian: msgFormulaAeolian, mark.TwoFive: msgFormulaTwoFive}
	workPhrase    = [mark.Situations]string{mark.BarStart: msgWorkStart, mark.MidBar: msgWorkMid, mark.Held: msgWorkHeld}
)

// adviceWords says the advice `a`, empty for none.
func adviceWords(g *game, a mark.Advice) string {
	bpm := fmt.Sprintf("%.0f", a.BPM)
	switch a.Kind {
	case mark.WorkFormula:
		return g.lang.T(msgAdviceWork, "What", g.lang.T(formulaPhrase[a.What]))
	case mark.WorkSituation:
		return g.lang.T(msgAdviceWork, "What", g.lang.T(workPhrase[a.Situation]))
	case mark.Slower:
		return g.lang.T(msgAdviceSlower, "BPM", bpm)
	case mark.Faster:
		return g.lang.T(msgAdviceFaster, "BPM", bpm)
	case mark.Stay:
		return g.lang.T(msgAdviceStay)
	}
	return ""
}

// situationPhrase names each situation.
var situationPhrase = [mark.Situations]string{mark.BarStart: msgReviewStart, mark.MidBar: msgReviewMid, mark.Held: msgReviewHeld}

func (r *review) Draw(dst *ebiten.Image) {
	g := r.g
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: g.scale}
	c.Text(g.lang.T(msgReviewTitle, "Title", g.title), g.fonts.ui, margin, titleY, ink)

	y := float64(chartY)
	y = r.drawSituations(c, y, true, msgReviewWorked, workedTint)
	y = r.drawSituations(c, y, false, msgReviewConsolidate, consolidateTint)
	if f := r.s.Feel(); f != mark.FeelNone {
		c.Text(g.lang.T(msgReviewTime), g.fonts.bubble, margin, y, ink)
		y += reviewLine + 3
		c.Text(g.lang.T(feelPhrase[f]), g.fonts.ui, margin, y, ink)
	}

	r.drawMap(c)
	if a := adviceWords(g, mark.Advise(r.s, g.bpm, minBPM, maxBPM)); a != "" {
		c.Text(a, g.fonts.bubble, margin, adviceY, ink)
	}
	c.Text(g.lang.T(msgReviewKeys), g.fonts.ui, margin, statusY, faint)
}

// drawSituations draws, from `y`, the situations that worked, or those
// to consolidate, under their heading and its swatch; nothing when
// there is none. It returns where the next block starts.
func (r *review) drawSituations(c screen.Canvas, y float64, worked bool, heading string, tint color.Color) float64 {
	g := r.g
	var lines []string
	for sit, t := range r.s.BySituation {
		if t.Expected == 0 || t.Worked() != worked {
			continue
		}
		lines = append(lines, g.lang.T(situationPhrase[sit], "Landed", t.Landed, "Expected", t.Expected))
	}
	if len(lines) == 0 {
		return y
	}
	c.Rect(margin, float32(y)+3, swatch, swatch, tint)
	c.Text(g.lang.T(heading), g.fonts.bubble, margin+swatch+5, y, ink)
	y += reviewLine + 3
	for _, l := range lines {
		c.Text(l, g.fonts.ui, margin, y, ink)
		y += reviewLine
	}
	return y + reviewGap
}

// drawMap draws the whole chart, its rows shrunk when the grid is long,
// each bar tinted: worked, to consolidate, or left blank when it
// expected nothing; the rehearsal marks as in the game.
func (r *review) drawMap(c screen.Canvas) {
	g := r.g
	rows := (len(g.bars) + barsPerRow - 1) / barsPerRow
	h := min(float32(rowH), float32(reviewChart)/float32(rows))
	pen := g.pen()
	for i, cells := range g.bars {
		x := float32(chartX + i%barsPerRow*barW)
		y := float32(chartY) + float32(i/barsPerRow)*h
		var fill color.Color
		switch t := r.s.ByBar[i]; {
		case t.Worked():
			fill = workedTint
		case t.Expected > 0:
			fill = consolidateTint
		}
		pen.DrawBar(c, cells, x, y, barW, h-4, fill, h < rowH || g.small, i%barsPerRow == barsPerRow-1 || i == len(g.bars)-1)
		if label, ok := g.rehearsals[i]; ok {
			pen.DrawRehearsal(c, label, x, y)
		}
	}
}
