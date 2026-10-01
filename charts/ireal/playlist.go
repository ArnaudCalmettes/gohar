package ireal

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	scheme = "irealb://"

	// chartPrefix opens every chart, before the scrambled part.
	chartPrefix = "1r34LbKcu7"

	// songSeparator separates the songs of a playlist. The empty third
	// field of a song only ever makes "==", so it cannot be mistaken
	// for it.
	songSeparator = "==="
)

// A Playlist is what one URL holds: a name, empty for a single song,
// and its songs.
type Playlist struct {
	Name  string
	Songs []Song
}

// A Song is one chart with what the app keeps about it.
type Song struct {
	Title    string
	Composer string
	Style    string // the style shown on the chart: Medium Swing, Bossa Nova
	Key      string // as the app writes it: Eb, A- for A minor

	// Transpose is how many semitones the player transposed the song
	// by in the app.
	Transpose int

	// Chart is the chart unscrambled, abbreviations included: see
	// [Lex].
	Chart string

	CompStyle string // the accompaniment style: Jazz-Medium Swing, Pop-Funk
	BPM       int    // zero when the app's default applies
	Repeats   int
}

// ErrScheme is returned for a URL that is not irealb://. The older
// irealbook:// is not read yet: no real sample of it has been seen.
var ErrScheme = errors.New("ireal: not an irealb:// URL")

// Parse reads an irealb:// URL, as the app shares it or as it appears
// in the href of an exported page. Anything before the scheme is
// ignored.
func Parse(link string) (Playlist, error) {
	i := strings.Index(link, scheme)
	if i < 0 {
		return Playlist{}, ErrScheme
	}
	// PathUnescape and not QueryUnescape: a plus sign is an augmented
	// chord, not a space.
	decoded, err := url.PathUnescape(link[i+len(scheme):])
	if err != nil {
		return Playlist{}, fmt.Errorf("ireal: %w", err)
	}

	var p Playlist
	parts := strings.Split(decoded, songSeparator)
	if len(parts) > 1 {
		p.Name = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}
	for i, part := range parts {
		s, err := parseSong(part)
		if err != nil {
			return Playlist{}, fmt.Errorf("ireal: song %d: %w", i+1, err)
		}
		p.Songs = append(p.Songs, s)
	}
	return p, nil
}

// parseSong reads the fields of one song around its chart, which is
// the only field that can be found for sure.
//
// The app writes ten fields today. Older writers dropped the empty
// third one, the transposition or the accompaniment style, so the
// fields are read from the chart outwards rather than by position:
// before it, the key, preceded by the transposition when there is one;
// after it, the accompaniment style when three fields remain, then the
// tempo and the repeats.
func parseSong(s string) (Song, error) {
	// The fields always found first, before those that may be dropped.
	const (
		titleField = iota
		composerField
		firstOptionalField
	)
	f := strings.Split(s, "=")
	c := -1
	for i, field := range f {
		if strings.HasPrefix(field, chartPrefix) {
			c = i
			break
		}
	}
	if c < firstOptionalField {
		return Song{}, errors.New("no chart")
	}

	song := Song{
		Title:    f[titleField],
		Composer: f[composerField],
		Chart:    unscramble(f[c][len(chartPrefix):]),
	}

	before := f[firstOptionalField:c]
	if n := len(before); n > 0 {
		if t, err := strconv.Atoi(before[n-1]); err == nil {
			song.Transpose = t
			before = before[:n-1]
		}
	}
	// What remains is [empty] style key, or style key, or key alone.
	if n := len(before); n > 0 {
		song.Key = before[n-1]
		if n > 1 {
			song.Style = before[n-2]
		}
	}

	after := f[c+1:]
	if len(after) == 3 {
		song.CompStyle, after = after[0], after[1:]
	}
	ints := []*int{&song.BPM, &song.Repeats}
	for i, v := range after {
		if i >= len(ints) || v == "" {
			continue
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return Song{}, fmt.Errorf("field %q after the chart: %w", v, err)
		}
		*ints[i] = n
	}
	return song, nil
}
