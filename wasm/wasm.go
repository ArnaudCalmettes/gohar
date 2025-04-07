//go:build js && wasm

package wasm

import (
	"fmt"
	"syscall/js"

	"github.com/ArnaudCalmettes/gohar"
	"github.com/ArnaudCalmettes/gohar/lib/js/convert"
)

func ExportJSFuncs() {
	js.Global().Set("gohar", js.ValueOf(map[string]any{
		"isLoaded":            js.ValueOf(true),
		"setLocale":           js.FuncOf(SetLocale),
		"noteName":            js.FuncOf(NoteName),
		"scaleName":           js.FuncOf(ScaleName),
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

func ScaleName(_ js.Value, args []js.Value) any {
	if len(args) != 2 {
		panic(fmt.Errorf("scaleName: expected 2 args, got %d", len(args)))
	}
	note := convert.NoteFromJS(args[0])
	pattern := convert.ScalePatternFromJS(args[1])
	result, err := gohar.ScaleName(gohar.Scale{Root: note, Pattern: pattern})
	if err != nil {
		panic(fmt.Errorf("scaleName: %w", err))
	}
	return js.ValueOf(result)
}

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
