package ireal

import (
	"sort"
	"strings"
)

// A Kind says what a token of a chart is.
type Kind uint8

const (
	Unknown Kind = iota // a character the lexer has no rule for, kept as is

	Space       // positions what follows in its cell; also any stray whitespace, see singles
	Comma       // separates two chords without moving them
	EmptyCells  // "XyQ"
	Bar         // "|", or "LZ"
	BarRepeated // "Kcl": a bar line, then the bar before repeated
	DoubleOpen  // "["
	DoubleClose // "]"
	RepeatOpen  // "{"
	RepeatClose // "}"
	FinalBar    // "Z"

	Section       // "*A": a rehearsal mark, see Token.Section
	TimeSignature // "T44", see Token.Time
	Ending        // "N1", see Token.Ending
	Segno         // "S"
	Coda          // "Q"
	Fermata       // "f"
	Comment       // "<D.S. al Coda>", see Token.Comment

	RepeatOne // "x": the bar before, again
	RepeatTwo // "r": the two bars before, again
	NoChord   // "n"
	Slash     // "p": a beat on the chord before
	Spacer    // "Y": vertical room before a line
	PlayerEnd // "U": where the player stops
	Small     // "s": the chords that follow are written small
	Large     // "l": back to large

	Chord     // see Token.Chord
	Alternate // "(C-7)": a chord written small above, see Token.Chord
)

var kindNames = [...]string{
	"Unknown", "Space", "Comma", "EmptyCells", "Bar", "BarRepeated",
	"DoubleOpen", "DoubleClose", "RepeatOpen", "RepeatClose", "FinalBar",
	"Section", "TimeSignature", "Ending", "Segno", "Coda", "Fermata",
	"Comment", "RepeatOne", "RepeatTwo", "NoChord", "Slash", "Spacer",
	"PlayerEnd", "Small", "Large", "Chord", "Alternate",
}

func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "Kind(?)"
}

// A Token is one piece of a chart. Raw is its exact text: the Raw of
// every token, joined in order, is the chart again.
type Token struct {
	Kind Kind
	Pos  int // byte offset in the chart
	Raw  string

	Section byte        // the mark: A to D, i for intro, v for verse
	Time    TimeSig     // for a TimeSignature
	Ending  int         // 1, 2, 3
	Comment CommentText // for a Comment
	Chord   ChordSymbol // for a Chord or an Alternate
}

// A TimeSig is a time signature, as a numerator and a denominator.
type TimeSig struct {
	Beats, Unit int
}

// A CommentText is the text of a comment. The app can move a comment
// down with a leading "*" and two digits; Shift keeps them.
type CommentText struct {
	Text    string
	Shift   int
	Shifted bool
}

// A ChordSymbol is a chord as written, cut into its three parts but not
// read: Root is a letter and its accidental, or "W" for a root left
// invisible that repeats the one before; Quality is what follows, in
// the app's dialect (^7, -7b5, 7alt); Bass is the note after the slash,
// if any.
//
// Custom marks a quality the player typed as free text between two
// stars, "A*m7*", outside the app's list: the text is kept as the
// quality, and it is for the chord symbol parser to make sense of it.
type ChordSymbol struct {
	Root, Quality, Bass string
	Custom              bool
}

// qualities is every chord quality the app writes, longest first so
// that the lexer takes "7sus" before "7". From the app's own list,
// completed from real charts.
var qualities = func() []string {
	q := []string{
		"5", "2", "add9", "+", "o", "h", "sus", "^", "-", "^7", "-7", "7",
		"7sus", "h7", "o7", "^9", "^13", "6", "69", "^7#11", "^9#11",
		"^7#5", "-6", "-69", "-^7", "-^9", "-9", "-11", "-7b5", "h9",
		"-b6", "-#5", "-add9", "9", "7b9", "7#9", "7#11", "7b5", "7#5",
		"9#11", "9b5", "9#5", "7b13", "7#9#5", "7#9b5", "7#9#11",
		"7b9#11", "7b9b5", "7b9#5", "7b9#9", "7b9b13", "7alt", "13",
		"13#11", "13#9", "13b9", "11", "7b9sus", "7susadd3", "9sus",
		"13sus", "7b13sus",
	}
	sort.Slice(q, func(i, j int) bool { return len(q[i]) > len(q[j]) })
	return q
}()

// Lex cuts a chart into tokens, without losing a character.
//
// It never fails: what it has no rule for comes out as an Unknown token
// of one byte, so that a reader can report it and a writer still gets
// the chart back whole.
func Lex(chart string) []Token {
	var out []Token
	for i := 0; i < len(chart); {
		t := lexOne(chart, i)
		out = append(out, t)
		i += len(t.Raw)
	}
	return out
}

// singles are the tokens of one fixed character.
//
// A line break or a tab can end up inside a chart: one song of a real
// playlist ends in "Z" followed by a newline and four spaces, saved as
// such by the app. They mean nothing, and are spaces like the others.
var singles = map[byte]Kind{
	' ': Space, '\n': Space, '\r': Space, '\t': Space, ',': Comma, '|': Bar, '[': DoubleOpen, ']': DoubleClose,
	'{': RepeatOpen, '}': RepeatClose, 'Z': FinalBar, 'S': Segno,
	'Q': Coda, 'f': Fermata, 'x': RepeatOne, 'r': RepeatTwo,
	'n': NoChord, 'p': Slash, 'Y': Spacer, 'U': PlayerEnd, 's': Small,
	'l': Large,
}

func lexOne(s string, i int) Token {
	rest := s[i:]
	tok := func(k Kind, n int) Token { return Token{Kind: k, Pos: i, Raw: rest[:n]} }

	switch {
	case strings.HasPrefix(rest, "XyQ"):
		return tok(EmptyCells, 3)
	case strings.HasPrefix(rest, "LZ"):
		return tok(Bar, 2)
	case strings.HasPrefix(rest, "Kcl"):
		return tok(BarRepeated, 3)
	}

	switch c := rest[0]; {
	case c == '*' && len(rest) > 1:
		t := tok(Section, 2)
		t.Section = rest[1]
		return t

	case c == 'T' && len(rest) > 2 && isDigit(rest[1]) && isDigit(rest[2]):
		t := tok(TimeSignature, 3)
		t.Time = TimeSig{int(rest[1] - '0'), int(rest[2] - '0')}
		// Two digits leave no room for twelve: T12 is 12/8.
		if t.Time == (TimeSig{1, 2}) {
			t.Time = TimeSig{12, 8}
		}
		return t

	case c == 'N' && len(rest) > 1 && isDigit(rest[1]):
		t := tok(Ending, 2)
		t.Ending = int(rest[1] - '0')
		return t

	case c == '<':
		end := strings.IndexByte(rest, '>')
		if end < 0 {
			return tok(Unknown, 1)
		}
		t := tok(Comment, end+1)
		t.Comment = parseComment(rest[1:end])
		return t

	case c == '(':
		end := strings.IndexByte(rest, ')')
		if end < 0 {
			return tok(Unknown, 1)
		}
		sym, n := lexChord(rest[1:end])
		if n != end-1 {
			return tok(Unknown, 1)
		}
		t := tok(Alternate, end+1)
		t.Chord = sym
		return t

	case isRoot(c):
		sym, n := lexChord(rest)
		t := tok(Chord, n)
		t.Chord = sym
		return t
	}

	if k, ok := singles[rest[0]]; ok {
		return tok(k, 1)
	}
	return tok(Unknown, 1)
}

// lexChord reads a chord symbol at the start of s and returns it with
// its length: a root, the longest quality the app knows, and a bass.
func lexChord(s string) (ChordSymbol, int) {
	var c ChordSymbol
	n := note(s)
	if n == 0 {
		return c, 0
	}
	c.Root = s[:n]
	if q, ok := customQuality(s[n:]); ok {
		c.Quality, c.Custom = q, true
		n += len(q) + 2
	} else {
		n += c.knownQuality(s[n:])
	}
	if n < len(s) && s[n] == '/' {
		if b := note(s[n+1:]); b > 0 && s[n+1] != 'W' {
			c.Bass = s[n+1 : n+1+b]
			n += 1 + b
		}
	}
	return c, n
}

// knownQuality takes the longest quality of the app's list at the start
// of s, and returns its length.
func (c *ChordSymbol) knownQuality(s string) int {
	for _, q := range qualities {
		if strings.HasPrefix(s, q) {
			c.Quality = q
			return len(q)
		}
	}
	return 0
}

// customQuality reads a free quality between two stars, "*m7*". Short,
// and with nothing in it that belongs to the chart around it, so that
// a rehearsal mark after a chord is not taken for one.
func customQuality(s string) (string, bool) {
	if len(s) < 3 || s[0] != '*' {
		return "", false
	}
	end := strings.IndexByte(s[1:], '*')
	if end <= 0 || end > 12 || strings.ContainsAny(s[1:1+end], " |[]{}<>,") {
		return "", false
	}
	return s[1 : 1+end], true
}

// note returns the length of a note at the start of s: a letter from A
// to G, or W, and an optional accidental. Zero when there is none.
func note(s string) int {
	if s == "" || !isRoot(s[0]) {
		return 0
	}
	if len(s) > 1 && s[0] != 'W' && (s[1] == '#' || s[1] == 'b') {
		return 2
	}
	return 1
}

func isRoot(c byte) bool  { return c >= 'A' && c <= 'G' || c == 'W' }
func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// parseComment reads the optional shift before the text of a comment.
func parseComment(s string) CommentText {
	if len(s) >= 3 && s[0] == '*' && isDigit(s[1]) && isDigit(s[2]) {
		return CommentText{Text: s[3:], Shift: int(s[1]-'0')*10 + int(s[2]-'0'), Shifted: true}
	}
	return CommentText{Text: s}
}
