package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/calibrate"
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
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
	itemPlay      = iota
	itemCalibrate // until a menu of options takes it over
	itemQuit
	items
)

var itemPhrase = [items]string{itemPlay: msgMenuPlay, itemCalibrate: msgMenuCalibrate, itemQuit: msgMenuQuit}

// title is the title scene: the jam plays, the reference bass walking
// the blues with the ride and the walker's snaps, while the player
// chooses. The keyboard sounds, to play along.
type title struct {
	*app

	walker walker
	chosen int
}

func newTitle(a *app) *title {
	return &title{app: a, walker: walker{ease: 1, streak: snapStreak}} // at his best: he snaps
}

// Enter starts the jam, unless it plays already: back from the
// calibration, the blues goes on.
func (t *title) Enter() {
	if t.jam == nil {
		t.jam = newJam(t.app)
	}
}

// Leave lets the jam play: the next scene takes it over, or stops it.
func (t *title) Leave() {}

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
		switch t.chosen {
		case itemQuit:
			return scene.Quit
		case itemCalibrate:
			return t.calibrate()
		}
		return scene.Replace(newGame(t.app))
	}
	t.jam.play(time.Now(), true)
	return scene.Stay
}

// calibrate opens the shared calibration on the jam's clock, its bar on
// the band's wood blocks: the bass goes on under it, the drums and the
// snaps stop, and come back once the measure is steady. Back to this
// very title, the menu where it was.
func (t *title) calibrate() scene.Transition {
	s, err := calibrate.New(calibrate.Config{
		Width: screenWidth,
		Scale: func() float64 { return t.scale },
		Lang:  t.lang.Tag(),
		Clock: t.jam.m,
		MIDI:  t.midi != "",
		Drain: t.drain,
		Tick:  t.band.tick,
		Play:  func(now time.Time, quiet bool) { t.jam.play(now, !quiet) },
		Back:  func() scene.Scene { return t },
	})
	if err != nil {
		log.Println("calibration:", err)
		return scene.Stay
	}
	return scene.Replace(s)
}

func (t *title) Draw(dst *ebiten.Image) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: t.scale}
	now := time.Now()
	beats := t.jam.m.Beats(now)
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
