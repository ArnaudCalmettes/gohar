package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/scene"
)

// The title screen, in the final position of the opening to come (see
// "L'écran titre" in docs/walk.md): the walker greyed on the left,
// snapping on the band, and the menu in the middle.
const (
	gameTitle = "Walk With Me" // a name: the same in every language

	// titleBPM is the band's tempo on the title screen, livelier than
	// the game's: it invites, it does not teach.
	titleBPM = 160
)

// The items of the menu, in order.
const (
	itemPlay = iota
	itemOptions
	itemQuit
	items
)

var itemPhrase = [items]string{itemPlay: msgMenuPlay, itemOptions: msgMenuOptions, itemQuit: msgMenuQuit}

// title is the title scene: the jam plays, the reference bass walking
// the blues with the ride and the walker's snaps, while the player
// chooses. The keyboard sounds, to play along.
type title struct {
	*app

	walker walker
	menu   menu
}

func newTitle(a *app) *title {
	return &title{app: a, walker: walker{ease: 1, streak: snapStreak}, menu: menu{items: items}} // at his best: he snaps
}

// Enter starts the jam, unless it plays already: back from the
// options, the blues goes on.
func (t *title) Enter() {
	if t.jam == nil {
		t.jam = newJam(t.app)
	}
}

// Leave lets the jam play: the next scene takes it over, or stops it.
func (t *title) Leave() {}

func (t *title) Update() scene.Transition {
	t.drain(nil) // the keys sound, from the callback; nothing to mark
	t.menu.move()
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return scene.Quit
	case confirmed():
		switch t.menu.chosen {
		case itemQuit:
			return scene.Quit
		case itemOptions:
			return scene.Replace(newOptions(t.app, t))
		}
		return scene.Replace(newGame(t.app))
	}
	t.jam.play(time.Now(), full)
	return scene.Stay
}

func (t *title) Draw(dst *ebiten.Image) {
	labels := make([]string, items)
	for i := range items {
		labels[i] = t.lang.T(itemPhrase[i])
	}
	t.drawMenu(dst, &t.walker, t.walker.gait(true), gameTitle, t.fonts.count, labels, t.menu.chosen, t.lang.T(msgTitleKeys))
}
