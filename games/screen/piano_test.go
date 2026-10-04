package screen

import (
	"image/color"
	"testing"
)

var (
	orange = color.RGBA{0xe0, 0xa0, 0x40, 0xff}
	blue   = color.RGBA{0x6c, 0xb4, 0xe8, 0xff}
)

// Three octaves, C3 to C6: 22 white keys over 220 units, 10 each.
func threeOctaves() *Piano {
	return &Piano{Low: 48, High: 84, X: 0, Y: 0, W: 220, H: 70}
}

func TestKeyRect(t *testing.T) {
	p := threeOctaves()
	if x, _, w, _ := p.KeyRect(50); x != 10 || w != 10 { // D3, the second white key
		t.Errorf("D3 at %v, %v wide", x, w)
	}
	if x, _, w, h := p.KeyRect(49); x != 7 || w != 6 || h != 70*blackLength { // C♯3, across C3 and D3
		t.Errorf("C♯3 at %v, %v wide, %v long", x, w, h)
	}
}

// The player first: a key both play glows the player's colour, and a
// released one fades.
func TestUpdate(t *testing.T) {
	p := threeOctaves()
	var player, demo [MIDIKeys]bool
	player[60], demo[60], demo[64] = true, true, true
	p.Update(Lit{&player, orange}, Lit{&demo, blue})
	if p.color[60] != orange || p.color[64] != blue {
		t.Errorf("C4 lit %v, E4 lit %v", p.color[60], p.color[64])
	}
	demo[64] = false
	p.Update(Lit{&player, orange}, Lit{&demo, blue})
	if p.glow[64] != glowDecay {
		t.Errorf("E4 released keeps %v of its light, want %v", p.glow[64], glowDecay)
	}
}

// A low F, F1, below the range, lights F3, dimmed.
func TestFold(t *testing.T) {
	p := threeOctaves()
	var down [MIDIKeys]bool
	down[29] = true
	p.Update(Lit{&down, orange})
	levels, _ := p.levels()
	if levels[53] != ghost {
		t.Errorf("F1 lights F3 at %v, want %v", levels[53], ghost)
	}
}
