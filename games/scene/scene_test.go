package scene

import (
	"errors"
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// fake is a scene that writes down what the director does to it, and
// asks for `next` on its next Update.
type fake struct {
	name string
	log  *[]string
	next Transition
}

func (f *fake) Enter() { *f.log = append(*f.log, f.name+" enters") }
func (f *fake) Leave() { *f.log = append(*f.log, f.name+" leaves") }

func (f *fake) Update() Transition {
	*f.log = append(*f.log, f.name+" updates")
	t := f.next
	f.next = Stay
	return t
}

func (f *fake) Draw(*ebiten.Image) { *f.log = append(*f.log, f.name+" draws") }

func noLayout(w, h int) (int, int) { return w, h }

// run ticks `d` once, draws it, and returns what happened.
func run(t *testing.T, d *Director, log *[]string) ([]string, error) {
	t.Helper()
	*log = nil
	err := d.Update()
	d.Draw(nil)
	return *log, err
}

func check(t *testing.T, step string, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("%s:\n got %q\nwant %q", step, got, want)
	}
}

// The title, the game over it by Replace, the settings pushed over the
// game and popped, then a quit: who enters, leaves, updates and draws.
func TestTitleGameSettings(t *testing.T) {
	var log []string
	title := &fake{name: "title", log: &log}
	game := &fake{name: "game", log: &log}
	settings := &fake{name: "settings", log: &log}
	d := New(title, noLayout)

	got, _ := run(t, d, &log)
	check(t, "first tick", got, []string{"title enters", "title updates", "title draws"})

	title.next = Replace(game)
	got, _ = run(t, d, &log)
	check(t, "title to game", got, []string{"title updates", "title leaves", "game enters", "game draws"})

	game.next = Push(settings)
	got, _ = run(t, d, &log)
	check(t, "settings opened", got, []string{"game updates", "settings enters", "game draws", "settings draws"})

	got, _ = run(t, d, &log)
	check(t, "settings on top", got, []string{"settings updates", "game draws", "settings draws"})

	settings.next = Pop
	got, _ = run(t, d, &log)
	check(t, "settings closed", got, []string{"settings updates", "settings leaves", "game draws"})

	game.next = Quit
	got, err := run(t, d, &log)
	check(t, "quit", got, []string{"game updates", "game leaves"})
	if !errors.Is(err, ebiten.Termination) {
		t.Errorf("quit: got %v, want ebiten.Termination", err)
	}
}

// Popping the last scene ends the game too.
func TestPopTheLast(t *testing.T) {
	var log []string
	title := &fake{name: "title", log: &log, next: Pop}
	d := New(title, noLayout)
	if _, err := run(t, d, &log); !errors.Is(err, ebiten.Termination) {
		t.Errorf("got %v, want ebiten.Termination", err)
	}
}

// Closing the window leaves the scenes still on the stack, from the top
// down; a director never started leaves nothing.
func TestClose(t *testing.T) {
	var log []string
	game := &fake{name: "game", log: &log}
	settings := &fake{name: "settings", log: &log}
	d := New(game, noLayout)
	d.Close()
	check(t, "never started", log, nil)

	d = New(game, noLayout)
	game.next = Push(settings)
	run(t, d, &log)
	log = nil
	d.Close()
	check(t, "closed", log, []string{"settings leaves", "game leaves"})
}
