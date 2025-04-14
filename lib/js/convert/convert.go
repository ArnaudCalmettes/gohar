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

func PitchClassToJS(p PitchClass) any {
	return js.ValueOf(p.Serialize())
}

func PitchClassFromJS(value js.Value) PitchClass {
	return DeserializePitchClass(uint16(value.Int()))
}

func NoteToJS(note Note) any {
	repr := uint(note.PitchClass.Serialize())
	repr |= uint(note.Oct+64) << 16
	return js.ValueOf(int(repr))
}

func NoteFromJS(value js.Value) Note {
	repr := uint(value.Int())
	return Note{
		PitchClass: DeserializePitchClass(uint16(value.Int())),
		Oct:        int8(repr&0xff0000>>16) - 64,
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
