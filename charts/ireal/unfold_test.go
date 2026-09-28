package ireal

import (
	"strings"
	"testing"
)

// played unfolds a chart and names each bar by its first chord, or by
// "-" when nothing new sounds in it.
func played(chart string) string {
	var names []string
	for _, b := range Structure(Lex(chart)).Unfold() {
		name := "-"
		if len(b.Events) > 0 {
			name = b.Events[0].Chord.Root
		}
		names = append(names, name)
	}
	return strings.Join(names, " ")
}

func TestUnfold(t *testing.T) {
	for name, c := range map[string]struct{ chart, want string }{
		"plain":                         {"[CXyQ|DXyQ|EXyQZ", "C D E"},
		"repeat":                        {"{CXyQ|DXyQ}EXyQZ", "C D C D E"},
		"repeat 3x":                     {"{CXyQ|DXyQ<3x>}EXyQZ", "C D C D C D E"},
		"ending after the sign":         {"{CXyQ|N1DXyQ}XyQXyQ LZN2EXyQ|FXyQZ", "C D C E F"},
		"endings inside":                {"{CXyQ|N1DXyQ|N2EXyQ}FXyQZ", "C D C E F"},
		"endings out of order":          {"{CXyQ|N1DXyQ}XyQ|N1EXyQZ", "C D C E"},
		"closed, never opened":          {"[CXyQ|DXyQ}EXyQZ", "C D C D E"},
		"opened, never closed":          {"{CXyQ|DXyQ]EXyQZ", "C D E"},
		"x repeats a bar":               {"[CXyQKcl LZDXyQZ", "C C D"},
		"r repeats two":                 {"[CXyQ|DXyQ|XyQr|XyQLZEXyQZ", "C D C D E"},
		"empty goes on":                 {"[CXyQ|XyQ|DXyQZ", "C - D"},
		"D.C. al Fine":                  {"[CXyQ|D<Fine>XyQ|EXyQ<D.C. al Fine>Z", "C D E C D"},
		"D.S. al Coda":                  {"[CXyQ|SDXyQ|E Q XyQ|F<D.S. al Coda>XyQ][QGXyQZ", "C D E F D E G"},
		"D.S. takes the last ending":    {"[SCXyQ{DXyQ|N1EXyQ}XyQ LZN2F<Fine>XyQ<D.S. al Fine>|GXyQZ", "C D E D F C D F"},
		"a direction after other words": {"[CXyQ|D<Fine>XyQ|EXyQ<After solos, D.C. al Fine>Z", "C D E C D"},
		// I Get A Kick Out Of You: the fine is in the first ending.
		"D.C. al 1st ending": {"{CXyQ|N1D<Fine>XyQ}XyQ LZN2EXyQ|FXyQ<D.C. al 1st ending>Z", "C D C E F C D"},
		// Cherokee: the fine is in the second ending.
		"D.C. al 2nd ending": {"{CXyQ|N1DXyQ}XyQ LZN2E<Fine>XyQ|FXyQ<D.C. al 2nd ending>Z", "C D C E F C E"},
		// I Believe In You: the third ending is written after the D.C.
		"D.C. al 3rd ending":            {"{CXyQ|N1DXyQ}XyQ LZN2EXyQ|FXyQ<D.C. al 3rd End.>|N3GXyQZ", "C D C E F C G"},
		"a bare D.C. plays to the end":  {"[CXyQ|DXyQ<D.C.>|EXyQZ", "C D C D E"},
		"a D.C. on cue is not followed": {"[CXyQ|DXyQ<D.C. on cue>|EXyQZ", "C D E"},
	} {
		if got := played(c.chart); got != c.want {
			t.Errorf("%s: %s, want %s", name, got, c.want)
		}
	}
}

// A chart that jumps forever is still cut short.
func TestUnfoldIsBounded(t *testing.T) {
	ms := []Measure{{Close: RepeatClose, Cells: 4}}
	ms[0].Marks = []Mark{{Kind: Comment, Comment: CommentText{Text: "99x"}}}
	if n := len(Chart{ms}.Unfold()); n > 99 || n == 0 {
		t.Errorf("%d bars", n)
	}
}
