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
	menu       string
	activities map[string]string // activity ID → its name in the menu
	question   string            // question number, total
	retry      string
	which      string // tonic name
	whichShape string // tonic name
	tetrachord string // qualifier from naming: majeur, phrygien
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
	menu:       "Choisir une activité : chiffre ou clic",
	activities: map[string]string{"tetrachords": "Tétracordes", "modes": "Modes", "modes-all": "Les sept modes"},
	question:   "Question %d / %d",
	retry:      "(reprise)",
	which:      "Quel est ce mode ? Tonique : %s",
	whichShape: "Quel est ce tétracorde ? Tonique : %s",
	tetrachord: "le tétracorde %s",
	replay:     "R : réécouter",
	answerKeys: "Chiffre ou clic pour répondre",
	right:      "Oui !",
	wrong:      "C'était %s. Écoute les deux.",
	next:       "Espace : suivant",
	end:        "Fin de la série",
	discovered: "Découverte : %s",
	recognized: "Reconnu pour la première fois : %s",
	nothingNew: "Rien de nouveau cette fois, mais tout est rafraîchi.",
	again:      "Espace : une autre série   Entrée : menu",
	noMIDI:     "sans clavier MIDI",
	keys:       "L : English   N : signes ou mots",
}

var english = words{
	menu:       "Pick an activity: digit or click",
	activities: map[string]string{"tetrachords": "Tetrachords", "modes": "Modes", "modes-all": "All seven modes"},
	question:   "Question %d / %d",
	retry:      "(retry)",
	which:      "Which mode is this? Tonic: %s",
	whichShape: "Which tetrachord is this? Tonic: %s",
	tetrachord: "the %s tetrachord",
	replay:     "R: listen again",
	answerKeys: "Digit or click to answer",
	right:      "Yes!",
	wrong:      "It was %s. Listen to both.",
	next:       "Space: next",
	end:        "End of the series",
	discovered: "Discovered: %s",
	recognized: "Recognised for the first time: %s",
	nothingNew: "Nothing new this time, but everything is refreshed.",
	again:      "Space: another series   Enter: menu",
	noMIDI:     "no MIDI keyboard",
	keys:       "L: Français   N: signs or words",
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

// notion names what an activity offers, or what the dex reports, in
// this language and notation, as it reads in a sentence. For a mode,
// the systematic name; the alternatives are not shown yet.
func (l language) notion(n dex.Notion) string {
	switch n.Kind {
	case dex.KindMode:
		if m, ok := naming.Lookup(n.System, n.Degree); ok {
			return l.namer.ModeNameOf(m)
		}
	case dex.KindTetrachord:
		return fmt.Sprintf(l.words.tetrachord, l.locale.TetrachordName(n.Tetrachord))
	}
	return n.String()
}

// label names a choice on its button, where the question already says
// what kind of thing it is: « phrygien » rather than « le tétracorde
// phrygien », which would not fit four abreast.
func (l language) label(n dex.Notion) string {
	if n.Kind == dex.KindTetrachord {
		return l.locale.TetrachordName(n.Tetrachord)
	}
	return l.notion(n)
}

// activity names an activity in the menu, falling back on its ID.
func (l language) activity(a Activity) string {
	if name, ok := l.words.activities[a.ID()]; ok {
		return name
	}
	return a.ID()
}

func (l language) note(c harmony.PitchClass) string {
	return l.namer.Name(c)
}

// change turns what the dex reports into a line worth showing, or ""
// for what is not worth it here.
func (l language) change(c dex.Change) string {
	name := l.notion(c.Notion)
	switch {
	case c.Discovery:
		return fmt.Sprintf(l.words.discovered, name)
	case c.Marked == dex.FactNamed:
		return fmt.Sprintf(l.words.recognized, name)
	}
	return ""
}
