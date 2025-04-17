//go:build js && wasm

package goharjs

import (
	"syscall/js"

	"github.com/ArnaudCalmettes/gohar"
)

func ScaleToJS(scale gohar.Scale) any {
	return js.ValueOf(map[string]any{
		"root":    PitchClassToJS(scale.Root),
		"pattern": ScalePatternToJS(scale.Pattern),
	})
}

func ScaleSliceToJS(scales []gohar.Scale) any {
	slice := make([]any, 0, 64)
	for _, scale := range scales {
		slice = append(slice, ScaleToJS(scale))
	}
	return js.ValueOf(slice)
}

func PitchSliceToJS(pitches []gohar.Pitch) any {
	slice := make([]any, 0, 12)
	for _, pitch := range pitches {
		slice = append(slice, int(pitch))
	}
	return js.ValueOf(slice)
}

func PitchClassToJS(p gohar.PitchClass) any {
	return js.ValueOf(int(p))
}

func PitchClassFromJS(value js.Value) gohar.PitchClass {
	return gohar.PitchClass(value.Int())
}

func NoteToJS(note gohar.Note) any {
	repr := uint(note.PitchClass)
	repr |= uint(note.Oct+64) << 8
	return js.ValueOf(int(repr))
}

func NoteFromJS(value js.Value) gohar.Note {
	repr := uint(value.Int())
	return gohar.Note{
		PitchClass: gohar.PitchClass(value.Int() & 0xff),
		Oct:        int8(repr&0xff00>>8) - 64,
	}
}

func NoteSliceToJS(notes []gohar.Note) any {
	slice := make([]any, 0, 12)
	for _, note := range notes {
		slice = append(slice, NoteToJS(note))
	}
	return js.ValueOf(slice)
}

func ScalePatternToJS(pattern gohar.ScalePattern) any {
	return js.ValueOf(int(pattern))
}

func ScalePatternFromJS(value js.Value) gohar.ScalePattern {
	return gohar.ScalePattern(value.Int())
}
