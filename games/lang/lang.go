// Package lang translates the phrases of the games: menus, status
// lines, marks, the walker's bubbles (see "L'internationalisation" in
// docs/architecture.md). The musical words, notes, chords, modes, are
// naming's business and never pass through here.
//
// It wraps go-i18n: each game embeds its message files, one TOML file
// per language, and asks for a phrase by its ID. English is the
// fallback: a phrase missing in French is said in English, and Missing
// lists those, for a test to catch.
package lang

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Fallback is the language said when a phrase or a language is missing.
const Fallback = "en"

// A Lang says the phrases of one game in one language.
type Lang struct {
	tag    string
	bundle *i18n.Bundle
	loc    *i18n.Localizer
}

// New loads every TOML file at the root of `files`, one per language,
// named after it as goi18n writes them (active.fr.toml), and speaks
// `want`, a language tag such as "fr" or "fr-FR". A language without a
// file speaks the fallback.
func New(files fs.FS, want string) (*Lang, error) {
	b := i18n.NewBundle(language.MustParse(Fallback))
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	names, err := fs.Glob(files, "*.toml")
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if _, err := b.LoadMessageFileFS(files, path.Clean(name)); err != nil {
			return nil, err
		}
	}
	tag := Fallback
	if t, err := language.Parse(want); err == nil {
		if m, _, conf := language.NewMatcher(b.LanguageTags()).Match(t); conf != language.No {
			base, _ := m.Base()
			tag = base.String()
		}
	}
	return &Lang{tag: tag, bundle: b, loc: i18n.NewLocalizer(b, tag)}, nil
}

// Tag is the language spoken: "fr", "en".
func (l *Lang) Tag() string { return l.tag }

// T says the phrase `id`. `args` fill its blanks, as name and value
// pairs: T("tempo", "BPM", 100) for "{{.BPM}} à la noire". A phrase
// missing everywhere says its ID, so that the gap shows on screen.
func (l *Lang) T(id string, args ...any) string {
	return l.say(id, nil, args)
}

// N says the phrase `id` for a count of `n`, in the plural form the
// language gives it: "1 arrivée posée", "3 arrivées posées". The count
// fills the blank {{.Count}}.
func (l *Lang) N(id string, n int, args ...any) string {
	return l.say(id, n, append(args, "Count", n))
}

func (l *Lang) say(id string, count any, args []any) string {
	data := map[string]any{}
	for i := 0; i+1 < len(args); i += 2 {
		if k, ok := args[i].(string); ok {
			data[k] = args[i+1]
		}
	}
	s, err := l.loc.Localize(&i18n.LocalizeConfig{MessageID: id, TemplateData: data, PluralCount: count})
	if s == "" && err != nil {
		return id
	}
	return s
}

// Missing returns the IDs among `ids` that the language spoken does not
// say itself, falling back or not at all. A test calls it for every
// language of a game, with every ID the game uses.
func (l *Lang) Missing(ids []string) []string {
	var out []string
	for _, id := range ids {
		_, tag, err := l.loc.LocalizeWithTag(&i18n.LocalizeConfig{MessageID: id})
		var missing *i18n.MessageNotFoundErr
		if errors.As(err, &missing) || tag.String() != l.tag {
			out = append(out, id)
		}
	}
	return out
}

// System returns the language of the user's session, from the usual
// variables, the fallback when none says.
func System() string {
	for _, v := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		s := os.Getenv(v)
		if s == "" || s == "C" || s == "POSIX" {
			continue
		}
		s, _, _ = strings.Cut(s, ".") // fr_FR.UTF-8
		return strings.ReplaceAll(s, "_", "-")
	}
	return Fallback
}
