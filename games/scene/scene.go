// Package scene is the director of a game's screens: a stack of scenes
// that Ebitengine sees as one ebiten.Game (see "Les scènes" in
// docs/architecture.md).
//
// A scene says what it wants, through the Transition its Update
// returns; it never sees the stack. Whatever the scenes share, the
// sound, the fonts, the language, they get from the game when they are
// built: the director knows nothing of it.
package scene

import "github.com/hajimehoshi/ebiten/v2"

// A Scene is one screen with its own logic.
//
//   - Enter is called when it reaches the stack, before its first
//     Update: what it schedules starts there, on time.
//   - Leave is called when it leaves the stack, replaced, popped, or
//     when the game ends: what it holds, a note, a file, is released
//     there. A scene covered by a Push is not left: it waits, intact.
//   - Update runs on the top scene only, every tick.
//   - Draw runs on every scene of the stack, from the bottom up: a scene
//     pushed over another dims it itself, if it wants to.
type Scene interface {
	Enter()
	Leave()
	Update() Transition
	Draw(dst *ebiten.Image)
}

type op int

const (
	stay op = iota
	replace
	push
	pop
	quit
)

// A Transition is what a scene asks of the director at the end of its
// Update. The zero value is Stay.
type Transition struct {
	op op
	to Scene
}

var (
	Stay = Transition{}         // carry on
	Pop  = Transition{op: pop}  // leave, back to the scene below
	Quit = Transition{op: quit} // leave every scene and end the game
)

// Replace leaves the top scene for `s`: from the title to the game.
func Replace(s Scene) Transition { return Transition{op: replace, to: s} }

// Push puts `s` over the top scene, which waits: the settings opened
// during a game.
func Push(s Scene) Transition { return Transition{op: push, to: s} }

// A Director holds the stack and plays the ebiten.Game.
type Director struct {
	stack   []Scene
	layout  func(outsideWidth, outsideHeight int) (int, int)
	started bool
}

// New starts the stack with `first`, entered on the first Update rather
// than now: the window may take a while to open, and a scene that
// schedules its sound on entering would start late. `layout` is the
// game's Layout.
func New(first Scene, layout func(outsideWidth, outsideHeight int) (int, int)) *Director {
	return &Director{stack: []Scene{first}, layout: layout}
}

// Update runs the top scene and applies its transition. An empty stack
// ends the game.
func (d *Director) Update() error {
	if !d.started {
		d.started = true
		if top := d.top(); top != nil {
			top.Enter()
		}
	}
	top := d.top()
	if top == nil {
		return ebiten.Termination
	}
	d.apply(top.Update())
	if len(d.stack) == 0 {
		return ebiten.Termination
	}
	return nil
}

func (d *Director) apply(t Transition) {
	switch t.op {
	case replace:
		d.top().Leave()
		d.stack[len(d.stack)-1] = t.to
		t.to.Enter()
	case push:
		d.stack = append(d.stack, t.to)
		t.to.Enter()
	case pop:
		d.top().Leave()
		d.stack = d.stack[:len(d.stack)-1]
	case quit:
		d.Close()
	}
}

// Draw draws the stack from the bottom up.
func (d *Director) Draw(dst *ebiten.Image) {
	for _, s := range d.stack {
		s.Draw(dst)
	}
}

func (d *Director) Layout(outsideWidth, outsideHeight int) (int, int) {
	return d.layout(outsideWidth, outsideHeight)
}

// Close leaves every scene, from the top down, and empties the stack.
// The game calls it once RunGame returns: closing the window ends the
// game without a Quit, and the scenes must still release what they
// hold.
func (d *Director) Close() {
	if d.started {
		for i := len(d.stack) - 1; i >= 0; i-- {
			d.stack[i].Leave()
		}
	}
	d.stack = nil
}

func (d *Director) top() Scene {
	if len(d.stack) == 0 {
		return nil
	}
	return d.stack[len(d.stack)-1]
}
