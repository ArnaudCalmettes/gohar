package main

import (
	"fmt"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// words is the interface text of one language. The musical words come
// from naming; these are only the game's own sentences, too few to
// deserve a package.
type words struct {
	start      string
	question   string // question number, total
	retry      string
	which      string // tonic name
	replay     string
	answerKeys string
	right      string
	wrong      string // right mode name
	next       string
	end        string
	discovered string // mode name
	recognized string // mode name
	nothingNew string
	again      string
	noMIDI     string
	keys       string
}

var french = words{
	start:      "Cliquer ou appuyer sur Espace pour commencer",
	question:   "Question %d / %d",
	retry:      "(reprise)",
	which:      "Quel est ce mode ? Tonique : %s",
	replay:     "R : réécouter",
	answerKeys: "1, 2, 3 ou clic pour répondre",
	right:      "Oui !",
	wrong:      "C'était %s. Écoute les deux.",
	next:       "Espace : suivante",
	end:        "Fin de la série",
	discovered: "Découverte : %s",
	recognized: "Reconnu pour la première fois : %s",
	nothingNew: "Rien de nouveau cette fois, mais tout est rafraîchi.",
	again:      "Espace : une autre série",
	noMIDI:     "sans clavier MIDI",
	keys:       "L : English   N : signes ou mots",
}

var english = words{
	start:      "Click or press Space to start",
	question:   "Question %d / %d",
	retry:      "(retry)",
	which:      "Which mode is this? Tonic: %s",
	replay:     "R: listen again",
	answerKeys: "1, 2, 3 or click to answer",
	right:      "Yes!",
	wrong:      "It was %s. Listen to both.",
	next:       "Space: next",
	end:        "End of the series",
	discovered: "Discovered: %s",
	recognized: "Recognised for the first time: %s",
	nothingNew: "Nothing new this time, but everything is refreshed.",
	again:      "Space: another series",
	noMIDI:     "no MIDI keyboard",
	keys:       "L: français   N: signs or words",
}

// language bundles what the game needs to speak one language.
type language struct {
	words  words
	locale naming.Locale
	namer  *naming.Namer
}

func newLanguage(code string) (language, error) {
	l := language{words: french, locale: naming.French}
	switch code {
	case "fr":
	case "en":
		l.words, l.locale = english, naming.English
	default:
		return language{}, fmt.Errorf("unknown language %q, want fr or en", code)
	}
	var err error
	l.namer, err = naming.NewNamer(l.locale)
	return l, err
}

// withNotation returns the same language writing accidentals as signs
// or words.
func (l language) withNotation(n naming.Notation) language {
	l.namer = l.namer.WithNotation(n)
	return l
}

// mode is the systematic name. The alternatives are not shown yet: the
// seven natural modes have none.
func (l language) mode(d harmony.Degree) string {
	m, ok := naming.Lookup(system, d)
	if !ok {
		return "?"
	}
	return l.namer.ModeNameOf(m)
}

func (l language) note(c harmony.PitchClass) string {
	return l.namer.Name(c)
}

// change turns what the dex reports into a line worth showing, or ""
// for what is not worth it here.
func (l language) change(c dex.Change) string {
	if c.Notion.Kind != dex.KindMode || c.Notion.System != system {
		return ""
	}
	name := l.mode(c.Notion.Degree)
	switch {
	case c.Discovery:
		return fmt.Sprintf(l.words.discovered, name)
	case c.Marked == dex.FactNamed:
		return fmt.Sprintf(l.words.recognized, name)
	}
	return ""
}
