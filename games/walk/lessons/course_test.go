package lessons

import (
	"slices"
	"testing"
)

// ids lists what a player sees: each chapter, then its lessons.
func ids(chapters []Chapter) []string {
	var out []string
	for _, c := range chapters {
		out = append(out, "chapter "+c.ID)
		for _, l := range c.Lessons {
			out = append(out, l.ID)
		}
	}
	return out
}

func TestOpen(t *testing.T) {
	for _, tc := range []struct {
		name string
		done []string
		want []string
	}{
		{"a new player sees the first lesson", nil, []string{"chapter 1", "1.1"}},
		{"a lesson done opens the next", []string{"1.1"}, []string{"chapter 1", "1.1", "1.2"}},
		{"the closing activity follows the last lesson", []string{"1.1", "1.2", "1.3", "1.4"},
			[]string{"chapter 1", "1.1", "1.2", "1.3", "1.4", "1.end"}},
		{"a chapter done opens the next chapter", []string{"1.1", "1.2", "1.3", "1.4", "1.end"},
			[]string{"chapter 1", "1.1", "1.2", "1.3", "1.4", "1.end", "chapter 2", "2.1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			done := Progress{}
			for _, id := range tc.done {
				done[id] = true
			}
			if got := ids(Open(Course, done)); !slices.Equal(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFinished(t *testing.T) {
	done := Progress{"1.1": true, "1.2": true, "1.3": true, "1.4": true}
	if done.Finished(Course[0]) {
		t.Error("chapter 1 finished before its closing activity")
	}
	done["1.end"] = true
	if !done.Finished(Course[0]) {
		t.Error("chapter 1 not finished with all its lessons done")
	}
}
