package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/settings"
	"github.com/ArnaudCalmettes/gohar/games/walk/lessons"
)

// courseFile keeps the lessons done, between two runs.
const courseFile = "course.json"

// The course's lists. The chapters sit under the big heading, as the
// title's menu does; a chapter's lessons, up to five, under its name in
// a smaller hand, from higher up. Back stands apart below either.
var (
	chaptersLayout = layout{top: menuY, gap: menuGap, apart: backApart}
	lessonsLayout  = layout{top: 120, gap: 32, apart: backApart}
)

const backApart = 18

// The phrases of the course name its chapters and lessons by ID, as
// lessons.Course gives them: "chapter.1", "lesson.1.1".
func chapterPhrase(id string) string { return "chapter." + id }
func lessonPhrase(id string) string  { return "lesson." + id }

// init lists the course's phrases with the others, for TestPhrases: the
// names of its chapters and lessons, and what the lessons say.
func init() {
	for _, c := range lessons.Course {
		phrases = append(phrases, chapterPhrase(c.ID))
		for _, l := range c.Lessons {
			phrases = append(phrases, lessonPhrase(l.ID))
		}
	}
	phrases = append(phrases, lessons.Phrases()...)
}

// course is the scene of the course for beginners, over the lighter
// jam: first the chapters open, then the lessons of the one chosen, a
// green check on what is done. A lesson written plays in a lesson
// scene (see lesson.go), which comes back here.
type course struct {
	*app

	back   *title // the title to go back to, the jam still playing
	walker walker
	menu   menu

	done    lessons.Progress
	open    []lessons.Chapter // what the player sees, a prefix of lessons.Course
	chapter int               // the chapter open, or -1 for the list of chapters
}

func newCourse(a *app, back *title) *course {
	done := lessons.Progress{}
	if _, err := settings.Load(settings.Path(courseFile), &done); err != nil {
		log.Println("course:", err) // the course starts over, the file stays as it is
	}
	c := &course{app: a, back: back, walker: walker{ease: grooveFrom}, done: done, open: lessons.Open(lessons.Course, done)}
	c.showChapters(0)
	return c
}

// showChapters lists the chapters, `chosen` chosen.
func (c *course) showChapters(chosen int) {
	c.chapter = -1
	c.menu = menu{chosen: chosen, items: len(c.open) + 1} // and Back
}

// showLessons lists the lessons of chapter `i`, `chosen` chosen.
func (c *course) showLessons(i, chosen int) {
	c.chapter = i
	c.menu = menu{chosen: chosen, items: len(c.open[i].Lessons) + 1}
}

// finish marks lesson `id` done, saves it, and opens what follows: the
// next lesson appears, chosen, in the chapter's list.
func (c *course) finish(id string) {
	c.done[id] = true
	if err := settings.Save(settings.Path(courseFile), c.done); err != nil {
		log.Println("course:", err)
	}
	c.open = lessons.Open(lessons.Course, c.done)
	chosen := c.menu.chosen
	if chosen+1 < len(c.open[c.chapter].Lessons) {
		chosen++
	}
	c.showLessons(c.chapter, chosen)
}

// Enter starts the jam again on the way back from a lesson, which stops
// it, without the count-in: here the walker grooves, he does not snap.
func (c *course) Enter() {
	if c.jam == nil {
		c.jam = newJam(c.app, false)
	}
}

func (c *course) Leave() {}

func (c *course) Update() scene.Transition {
	c.drain(nil)
	c.menu.move()
	back := c.menu.chosen == c.menu.items-1
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape), confirmed() && back:
		if c.chapter < 0 {
			return scene.Replace(c.back)
		}
		c.showChapters(c.chapter)
	case confirmed() && c.chapter < 0:
		c.showLessons(c.menu.chosen, 0)
	case confirmed():
		id := c.open[c.chapter].Lessons[c.menu.chosen].ID
		if steps := lessons.Script(id); steps != nil {
			return scene.Replace(newLesson(c.app, c, id, steps))
		}
	}
	c.jam.play(time.Now(), light)
	return scene.Stay
}

func (c *course) Draw(dst *ebiten.Image) {
	var labels []string
	var done []bool
	if c.chapter < 0 {
		for i, ch := range c.open {
			labels = append(labels, c.lang.T(chapterPhrase(ch.ID)))
			done = append(done, c.done.Finished(lessons.Course[i]))
		}
		labels = append(labels, c.lang.T(msgOptBack))
		c.drawList(dst, &c.walker, c.walker.gait(true), c.lang.T(msgMenuLearn), c.fonts.count, labels, done, c.menu.chosen, c.lang.T(msgCourseKeys), chaptersLayout)
		return
	}
	for _, l := range c.open[c.chapter].Lessons {
		labels = append(labels, c.lang.T(lessonPhrase(l.ID)))
		done = append(done, c.done[l.ID])
	}
	labels = append(labels, c.lang.T(msgOptBack))
	heading := c.lang.T(chapterPhrase(c.open[c.chapter].ID)) // longer than a menu's: in the chords' hand
	c.drawList(dst, &c.walker, c.walker.gait(true), heading, c.fonts.chord, labels, done, c.menu.chosen, c.lang.T(msgCourseKeys), lessonsLayout)
}
