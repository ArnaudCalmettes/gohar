//go:build js && wasm

package convert

import (
	"syscall/js"

	. "github.com/ArnaudCalmettes/gohar"
)

func ScaleToJS(scale Scale) any {
	return js.ValueOf(map[string]any{
		"root":    NoteToJS(scale.Root),
		"pattern": ScalePatternToJS(scale.Pattern),
	})
}

func ScaleSliceToJS(scales []Scale) any {
	slice := make([]any, 0, 64)
	for _, scale := range scales {
		slice = append(slice, ScaleToJS(scale))
	}
	return js.ValueOf(slice)
}

func PitchSliceToJS(pitches []Pitch) any {
	slice := make([]any, 0, 12)
	for _, pitch := range pitches {
		slice = append(slice, int(pitch))
	}
	return js.ValueOf(slice)
}

func NoteToJS(note Note) any {
	repr := int(note.Base) | int(note.Alt)<<8 | int(note.Oct)<<16
	return js.ValueOf(repr)
}

func NoteFromJS(value js.Value) Note {
	repr := value.Int()
	return Note{
		Base: byte(repr & 0xff),
		Alt:  Pitch((repr >> 8) & 0xff),
		Oct:  int8(repr >> 16 & 0xff),
	}
}

func NoteSliceToJS(notes []Note) any {
	slice := make([]any, 0, 12)
	for _, note := range notes {
		slice = append(slice, NoteToJS(note))
	}
	return js.ValueOf(slice)
}

func ScalePatternToJS(pattern ScalePattern) any {
	return js.ValueOf(int(pattern))
}

func ScalePatternFromJS(value js.Value) ScalePattern {
	return ScalePattern(value.Int())
}
