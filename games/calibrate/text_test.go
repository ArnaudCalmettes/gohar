package calibrate

import "testing"

// Every phrase the scene says is in every language it speaks.
func TestPhrases(t *testing.T) {
	for _, tag := range []string{"fr", "en"} {
		l, err := newLang(tag)
		if err != nil {
			t.Fatal(err)
		}
		if missing := l.Missing(phrases); len(missing) > 0 {
			t.Errorf("%s lacks %v", tag, missing)
		}
	}
}
