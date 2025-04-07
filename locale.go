package gohar

import (
	"errors"
	"fmt"
)

type Locale struct {
	NoteNames  []string
	ScaleNames map[ScalePattern]string
}

var (
	LocaleFrench = Locale{
		NoteNames: []string{"do", "ré", "mi", "fa", "sol", "la", "si"},
		ScaleNames: map[ScalePattern]string{
			ScalePatternMajor:               "majeur",
			ScalePatternMelodicMinor:        "mineur mélodique",
			ScalePatternHarmonicMinor:       "mineur harmonique",
			ScalePatternHarmonicMajor:       "majeur harmonique",
			ScalePatternDoubleHarmonicMajor: "majeur double harmonique",
		},
	}

	LocaleEnglish = Locale{
		NoteNames: []string{"c", "d", "e", "f", "g", "a", "b"},
		ScaleNames: map[ScalePattern]string{
			ScalePatternMajor:               "major",
			ScalePatternMelodicMinor:        "melodic minor",
			ScalePatternHarmonicMinor:       "harmonic minor",
			ScalePatternHarmonicMajor:       "harmonic major",
			ScalePatternDoubleHarmonicMajor: "double harmonic major",
		},
	}

	CurrentLocale = &LocaleEnglish
)

func (loc *Locale) NoteName(note Note) (string, error) {
	if err := CheckNoteIsPrintable(note); err != nil {
		return "", err
	}
	s := loc.basename(note.Base) + altToString(note.Alt)
	if note.Oct != 0 {
		s += fmt.Sprintf("%d", note.Oct)
	}
	return s, nil
}

func (loc *Locale) ScalePatternName(pattern ScalePattern) (string, error) {
	if name, ok := loc.ScaleNames[pattern]; ok {
		return name, nil
	}
	return "", ErrUnknownScalePattern
}

func (loc *Locale) ScaleName(scale Scale) (string, error) {
	note, noteErr := loc.NoteName(scale.Root)
	name, nameErr := loc.ScalePatternName(scale.Pattern)
	return note + " " + name, errors.Join(noteErr, nameErr)
}

func (loc *Locale) basename(b byte) string {
	idx := int(b) - int('C')
	if idx < 0 {
		idx += 7
	}
	return loc.NoteNames[idx]
}

var ErrLocaleNotSet = errors.New("gohar.CurrentLocale is not set")

func NoteName(note Note) (string, error) {
	if CurrentLocale != nil {
		return CurrentLocale.NoteName(note)
	}
	return "", ErrLocaleNotSet
}

func ScalePatternName(pattern ScalePattern) (string, error) {
	if CurrentLocale != nil {
		return CurrentLocale.ScalePatternName(pattern)
	}
	return "", ErrLocaleNotSet
}

func ScaleName(scale Scale) (string, error) {
	if CurrentLocale != nil {
		return CurrentLocale.ScaleName(scale)
	}
	return "", ErrLocaleNotSet
}
