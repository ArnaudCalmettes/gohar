package grids

import (
	"strings"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
)

// DefaultTempo is the tempo of a grid that says neither its tempo nor a
// style the game knows.
const DefaultTempo = 120

// styleTempos are the tempos the game gives the styles of iReal Pro, in
// lower case: a convention of the game, which no source sets, for a
// walking bass to sound right in the style. A grid that writes its
// {tempo} keeps it.
var styleTempos = map[string]int{
	"ballad":          70,
	"slow swing":      90,
	"medium swing":    120,
	"medium up swing": 160,
	"up tempo swing":  200,
}

// tempoOf is the tempo of the song `s`: the one it writes, or that of
// its style, or DefaultTempo.
func tempoOf(s chordpro.Song) int {
	if s.Tempo > 0 {
		return s.Tempo
	}
	for _, m := range s.Meta {
		if m.Name != "style" {
			continue
		}
		if t, ok := styleTempos[strings.ToLower(strings.TrimSpace(m.Value))]; ok {
			return t
		}
	}
	return DefaultTempo
}
