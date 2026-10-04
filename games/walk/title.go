package main

import (
	"image/color"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// The title screen, in the final position of the opening to come (see
// "L'écran titre" in docs/walk.md): the walker greyed on the left,
// snapping on the band, and the menu in the middle.
const (
	gameTitle = "Walk With Me" // a name: the same in every language

	titleNameY = 70  // the top of the name
	menuY      = 170 // the top of the first item
	menuGap    = 36  // from one item to the next
	menuPad    = 6   // around the item chosen

	// titleBPM is the band's tempo on the title screen, livelier than
	// the game's: it invites, it does not teach.
	titleBPM = 160
)

// The items of the menu, in order.
const (
	itemPlay = iota
	itemQuit
	items
)

var itemPhrase = [items]string{itemPlay: msgMenuPlay, itemQuit: msgMenuQuit}

// title is the title scene: the reference bass walks a blues chorus
// after chorus, with the ride and the walker's snaps, while the player
// chooses. The keyboard sounds, to play along.
type title struct {
	*app

	m      tempo.Metronome
	swing  tempo.Swing
	beats  []Beat // one chorus
	line   []int  // the chorus playing: a key per beat
	next   int    // the next beat to schedule, counted from the first chorus
	rng    *rand.Rand
	walker walker
	chosen int
}

func newTitle(a *app) *title {
	return &title{app: a, walker: walker{ease: 1, streak: snapStreak}} // at his best: he snaps
}

// Enter starts the count-in a lookahead from now: two bars of the
// walker's snaps on 2 and 4, alone, then the band comes in.
func (t *title) Enter() {
	now := time.Now()
	beat := time.Duration(float64(time.Minute) / titleBPM)
	t.m = tempo.NewMetronome(now.Add(lookahead+countIn*beat), titleBPM, perBar)
	t.next = -countIn
	t.swing = tempo.NewSwing(t.m, swingRatio)
	t.beats = Expect(t.grid, t.m, 1)
	seed := uint64(now.UnixNano())
	t.rng = rand.New(rand.NewPCG(seed, seed>>32|1))
}

// Leave releases the bass: the game counts in on its own.
func (t *title) Leave() {
	t.band.stop(time.Now())
}

func (t *title) Update() scene.Transition {
	t.drain(nil) // the keys sound, from the callback; nothing to mark
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return scene.Quit
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp):
		t.chosen = (t.chosen + items - 1) % items
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown):
		t.chosen = (t.chosen + 1) % items
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter), inpututil.IsKeyJustPressed(ebiten.KeySpace):
		if t.chosen == itemQuit {
			return scene.Quit
		}
		return scene.Replace(newGame(t.app))
	}
	t.schedule(time.Now())
	return scene.Stay
}

// schedule queues the beats due within the lookahead of `now`: the
// snaps of the count-in, then the band, a new line drawn at each chorus.
// Each chorus is walked on its own: the line does not carry over the
// double bar yet.
func (t *title) schedule(now time.Time) {
	_, end := t.m.Due(now, now.Add(lookahead))
	for ; t.next < end; t.next++ {
		if t.next < 0 {
			if t.m.Position(t.next).Beat%2 == 0 { // 2 and 4
				t.band.snapAt(t.m.At(t.next))
			}
			continue
		}
		n := t.next % len(t.beats)
		if n == 0 {
			t.line = Walk(t.beats, t.rng)
		}
		snap := t.walker.gait(true) == snapping
		t.band.beat(t.m.Position(t.next), t.line[n], t.m.At(t.next), t.swing.AtBeats(float64(t.next)+0.5), snap)
	}
}

func (t *title) Draw(dst *ebiten.Image) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: t.scale}
	now := time.Now()
	beats := t.m.Beats(now)
	t.walker.draw(c, walkerX, walkerY, walkerScale, t.walker.gait(true), beats, now, faint) // snapping from the count-in

	const mid = screenWidth / 2
	c.Centred(gameTitle, t.fonts.count, mid, titleNameY, ink)
	for i := range items {
		label := t.lang.T(itemPhrase[i])
		y := float64(menuY + i*menuGap)
		var col color.Color = faint
		if i == t.chosen {
			w, h := c.Measure(label, t.fonts.chord)
			c.Rect(float32(mid-w/2-menuPad), float32(y-menuPad), float32(w+2*menuPad), float32(h+2*menuPad), pale)
			col = ink
		}
		c.Centred(label, t.fonts.chord, mid, y, col)
	}
	c.Text(t.lang.T(msgTitleKeys), t.fonts.ui, margin, statusY, faint)
}
