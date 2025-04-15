package gohar

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	AltSharp       = "♯"
	AltFlat        = "♭"
	AltNatural     = "♮"
	AltDoubleSharp = "𝄪"
	AltDoubleFlat  = "𝄫"
)

var (
	noteRegexp        = regexp.MustCompile(`^([a-gA-G])(#|b|##|bb|♭|♯|𝄫|𝄪|♮)?([+-]?\d)?$`)
	toUnicodeReplacer = strings.NewReplacer(
		"bb", AltDoubleFlat,
		"##", AltDoubleSharp,
		"n", AltNatural,
		"b", AltFlat,
		"#", AltSharp,
	)
	_ = toUnicodeReplacer

	ErrCannotParseNote   = errors.New("cannot parse note")
	ErrUnknownAlteration = errors.New("unknown alteration")
)

func ParseNote(input string) (Note, error) {
	match := noteRegexp.FindStringSubmatch(input)
	var n Note
	if len(match) == 0 {
		return n, fmt.Errorf("%w: %q", ErrCannotParseNote, input)
	}
	base := strings.ToUpper(match[1])[0]
	alt, _ := ParseAlteration(match[2])
	oct, _ := strconv.Atoi(match[3])
	pc, err := NewPitchClassFromChar(base, alt)
	return Note{pc, int8(oct)}, err
}

func ParsePitch(input string) (Pitch, error) {
	if n, err := ParseNote(input); err != nil {
		return 0, err
	} else {
		return n.Pitch(), nil
	}
}

func ParseAlteration(alt string) (Pitch, error) {
	switch alt {
	case "bb", AltDoubleFlat:
		return -2, nil
	case "b", AltFlat:
		return -1, nil
	case "", "n", AltNatural:
		return 0, nil
	case "#", AltSharp:
		return 1, nil
	case "##", AltDoubleSharp:
		return 2, nil
	default:
		return 0, ErrUnknownAlteration
	}
}
