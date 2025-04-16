package abc

import (
	"fmt"
	"strings"

	"github.com/ArnaudCalmettes/gohar"
)

func ScaleToABC(root gohar.PitchClass, pattern gohar.ScalePattern) string {
	var w strings.Builder
	rootBase := root.Base()
	for pc := range pattern.PitchClasses(root) {
		var oct int8
		if pc.Base() < rootBase {
			oct = 1
		}
		fmt.Fprintf(&w, "%s ", NoteToABC(gohar.Note{PitchClass: pc, Oct: oct}))
	}
	return w.String()
}

func NoteToABC(n gohar.Note) string {
	name := PitchClassToABC(n.PitchClass)
	switch n.Oct {
	case 0:
		return strings.ToUpper(name)
	case 1:
		return strings.ToLower(name)
	}
	if n.Oct < 0 {
		return name + strings.Repeat(",", int(-n.Oct))
	} else {
		return name + strings.Repeat("'", int(n.Oct))
	}
}

func PitchClassToABC(pc gohar.PitchClass) string {
	return abcAlt(pc.Alt()) + gohar.LocaleEnglish.NoteNames[pc.Base()]
}

func abcAlt(alt gohar.Pitch) string {
	switch alt {
	case -2:
		return "__"
	case -1:
		return "_"
	case +1:
		return "^"
	case +2:
		return "^^"
	}
	return ""
}
