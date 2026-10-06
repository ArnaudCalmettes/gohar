// Package lessons holds the course for complete beginners (see
// docs/debutants.md): its chapters, their lessons, and those a player
// may open. It knows nothing of Walk with me, which lists the lessons
// and, with the lesson engine to come, plays them.
package lessons

// A Lesson is a step of a chapter, short enough for one session. Its
// ID, "1.1", names it in the player's progress and in the game's
// phrases.
type Lesson struct {
	ID string

	// Closing marks the activity that closes a chapter and brings the
	// player back to the game: palier 0 for chapter 1, palier 1 for
	// chapter 2.
	Closing bool
}

// A Chapter has an objective, a knowledge that unlocks a way of
// playing, and the lessons that lead to it.
type Chapter struct {
	ID      string
	Lessons []Lesson
}

// Course is the ramp up as written so far, in docs/debutants/. Chapter
// 3 has no closing activity yet.
var Course = []Chapter{
	{ID: "1", Lessons: []Lesson{{ID: "1.1"}, {ID: "1.2"}, {ID: "1.3"}, {ID: "1.4"}, {ID: "1.end", Closing: true}}},
	{ID: "2", Lessons: []Lesson{{ID: "2.1"}, {ID: "2.2"}, {ID: "2.3"}, {ID: "2.4"}, {ID: "2.end", Closing: true}}},
	{ID: "3", Lessons: []Lesson{{ID: "3.1"}, {ID: "3.2"}, {ID: "3.3"}, {ID: "3.4"}}},
}

// Progress is the lessons a player has done, by ID.
type Progress map[string]bool

// Open returns `course` as the player sees it: the lessons done, and the
// first one not done yet, ready whenever he wants it. The rest stays
// hidden, and so does a chapter with nothing to show. Chapters open in
// order, so the result is a prefix of `course`, chapter by chapter.
func Open(course []Chapter, done Progress) []Chapter {
	var out []Chapter
	next := true // the first lesson not done is still to show
	for _, c := range course {
		var shown []Lesson
		for _, l := range c.Lessons {
			switch {
			case done[l.ID]:
				shown = append(shown, l)
			case next:
				shown = append(shown, l)
				next = false
			}
		}
		if len(shown) == 0 {
			break
		}
		out = append(out, Chapter{ID: c.ID, Lessons: shown})
	}
	return out
}

// Finished tells whether every lesson of `c` is done.
func (p Progress) Finished(c Chapter) bool {
	for _, l := range c.Lessons {
		if !p[l.ID] {
			return false
		}
	}
	return true
}
