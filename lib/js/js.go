// package js implements Javascript bindings.

//go:build js && wasm

package js

import (
	"fmt"
	"syscall/js"

	"github.com/ArnaudCalmettes/gohar"
	"github.com/ArnaudCalmettes/gohar/lib/js/convert"
)

// ImportGoharBindings imports bindings to the current javascript environment
// under the "gohar" namespace. This function is intended to be run at
// initialization time from within a wasm binary.
func ImportGoharBindings() {
	js.Global().Set("gohar", js.ValueOf(map[string]any{
		"isLoaded":            js.ValueOf(true),
		"setLocale":           js.FuncOf(SetLocale),
		"noteName":            js.FuncOf(NoteName),
		"scalePatternName":    js.FuncOf(ScalePatternName),
		"scalePatternPitches": js.FuncOf(ScalePatternPitches),
		"scalePatterns": js.ValueOf([]any{
			int(gohar.ScalePatternMajor),
			int(gohar.ScalePatternMelodicMinor),
			int(gohar.ScalePatternHarmonicMinor),
			int(gohar.ScalePatternHarmonicMajor),
			int(gohar.ScalePatternDoubleHarmonicMajor),
		}),
	}))
}

// SetLocale sets gohar's locale.
// Supported locales: "en", "fr".
//
// TypeScript signature:
//
//	function setLocale(locale: string) => void
func SetLocale(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		panic(fmt.Errorf("setLocale: expected 1 arg, got %d", len(args)))
	}
	switch args[0].String() {
	case "fr":
		gohar.CurrentLocale = &gohar.LocaleFrench
	case "en":
		gohar.CurrentLocale = &gohar.LocaleEnglish
	}
	return js.Null()
}

// NoteName returns a note's name in the current locale.
//
// Typescript signature:
//
//	function noteName(note: number) => string
func NoteName(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		panic(fmt.Errorf("noteName: expected 1 arg, got %d", len(args)))
	}
	note := convert.NoteFromJS(args[0])
	result, err := gohar.NoteName(note)
	if err != nil {
		panic(fmt.Errorf("noteName: %w", err))
	}
	return js.ValueOf(result)
}

// ScalePatternName returns the name of a ScalePattern in the current locale.
//
// Typescript signature:
//
//	function scalePatternName(pattern: number) => string
func ScalePatternName(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		panic(fmt.Errorf("scalePatternName: expected 1 arg, got %d", len(args)))
	}
	pattern := convert.ScalePatternFromJS(args[0])
	result, err := gohar.ScalePatternName(pattern)
	if err != nil {
		panic(fmt.Errorf("scalePatternName: %w", err))
	}
	return js.ValueOf(result)
}

// ScalePatternPitches instanciates a scale pattern and returns the corresponding pitches.
// If rootPitch is provided, start the scale on this pitch instead of the default C (0).
//
// Typescript signature:
//
//	function scalePatternPitches(pattern: number, rootPitch?: number) => number[]
func ScalePatternPitches(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		panic(fmt.Errorf("scalePatternPitches: expected at least 1 arg, got %d", len(args)))
	}
	pattern := convert.ScalePatternFromJS(args[0])
	var root gohar.Pitch
	if len(args) >= 2 {
		root = gohar.Pitch(args[1].Int())
	}
	pitches, err := pattern.IntoPitches(make([]gohar.Pitch, 0, 12), root)
	if err != nil {
		panic(fmt.Errorf("scalePatternPitches: %w", err))
	}
	return convert.PitchSliceToJS(pitches)
}
