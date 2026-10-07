package main

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/settings"
	"github.com/ArnaudCalmettes/gohar/games/walk/band"
	"github.com/ArnaudCalmettes/gohar/games/walk/figure"
	"github.com/ArnaudCalmettes/gohar/games/walk/lessons"
)

// courseFile keeps the lessons done, between two runs.
const courseFile = "course.json"

// The course's lists. The chapters sit under the big heading, as the
// title's menu does; a chapter's lessons, up to five, under its name in
// a smaller hand, from higher up. Back stands apart below either.
var (
	chaptersLayout = layout{top: menuY, gap: menuGap, apart: backApart, left: true}
	lessonsLayout  = layout{heading: 48, top: 96, gap: 27, branchGap: 22, apart: backApart, left: true}
)

const backApart = 14

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
// green check on what is done. The lessons are a tree: the one chosen
// shows its branches under it, the lesson itself and, once done, its
// activity; Enter or the right arrow goes down to them, Escape or the
// left arrow back up. A lesson with no other branch than itself, not
// done yet, starts at once. A lesson or an activity plays in a lesson
// scene (see lesson.go), which comes back here.
type course struct {
	*app

	back   *title // the title to go back to, the jam still playing
	walker figure.Walker
	menu   menu

	done    lessons.Progress
	open    []lessons.Chapter // what the player sees, a prefix of lessons.Course
	chapter int               // the chapter open, or -1 for the list of chapters
	inside  bool              // down in the branches of the lesson chosen
	branch  menu              // the branch chosen there
	rng     *rand.Rand        // the order of the activities
}

// The branches of a lesson: the lesson, and its activity once done.
const (
	branchLesson = iota
	branchPractise
)

func newCourse(a *app, back *title) *course {
	done := lessons.Progress{}
	if _, err := settings.Load(settings.Path(courseFile), &done); err != nil {
		log.Println("course:", err) // the course starts over, the file stays as it is
	}
	seed := uint64(time.Now().UnixNano())
	c := &course{
		app: a, back: back, walker: figure.NewWalker(figure.Grooving),
		done: done, open: lessons.Open(lessons.Course, done),
		rng: rand.New(rand.NewPCG(seed, seed>>32|1)),
	}
	c.showChapters(0)
	return c
}

// showChapters lists the chapters, `chosen` chosen.
func (c *course) showChapters(chosen int) {
	c.chapter, c.inside = -1, false
	c.menu = menu{chosen: chosen, items: len(c.open) + 1} // and Back
}

// showLessons lists the lessons of chapter `i`, `chosen` chosen.
func (c *course) showLessons(i, chosen int) {
	c.chapter, c.inside = i, false
	c.menu = menu{chosen: chosen, items: len(c.open[i].Lessons) + 1}
}

// goDown opens the branches of the lesson chosen, `chosen` chosen.
func (c *course) goDown(chosen int) {
	c.inside = true
	c.branch = menu{chosen: chosen, items: len(c.branches(c.id()))}
}

// id is the ID of the lesson chosen in the chapter's list.
func (c *course) id() string { return c.open[c.chapter].Lessons[c.menu.chosen].ID }

// branches lists the branches of lesson `id`: the lesson, and its
// activity once the lesson is done, if written.
func (c *course) branches(id string) []int {
	if c.done[id] && lessons.Activity(id, c.rng) != nil {
		return []int{branchLesson, branchPractise}
	}
	return []int{branchLesson}
}

// play starts lesson `id`, or its activity when `activity`, if written.
func (c *course) play(id string, activity bool) scene.Transition {
	if activity {
		if steps := lessons.Activity(id, c.rng); steps != nil {
			return scene.Replace(newLesson(c.app, c, steps, func() {}))
		}
		return scene.Stay
	}
	if steps := lessons.Script(id); steps != nil {
		return scene.Replace(newLesson(c.app, c, steps, func() { c.finish(id) }))
	}
	return scene.Stay
}

// finish marks lesson `id` done, saves it, and opens what follows: its
// activity, chosen under it, not to be missed; or, without one, the
// next lesson, chosen in the chapter's list.
func (c *course) finish(id string) {
	c.done[id] = true
	if err := settings.Save(settings.Path(courseFile), c.done); err != nil {
		log.Println("course:", err)
	}
	c.open = lessons.Open(lessons.Course, c.done)
	chosen := c.menu.chosen
	if len(c.branches(id)) > 1 {
		c.showLessons(c.chapter, chosen)
		c.goDown(branchPractise)
		return
	}
	if chosen+1 < len(c.open[c.chapter].Lessons) {
		chosen++
	}
	c.showLessons(c.chapter, chosen)
}

// Enter starts the jam again on the way back from a lesson, which stops
// it, without the count-in: here the walker grooves, he does not snap.
func (c *course) Enter() {
	if c.jam == nil {
		c.jam = c.newJam(false)
	}
}

func (c *course) Leave() {}

func (c *course) Update() scene.Transition {
	c.drain(nil)
	c.jam.Play(time.Now(), band.Light)
	right := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight)
	up := inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft)
	down := confirmed() || right
	switch {
	case c.chapter < 0:
		c.menu.move()
		back := c.menu.chosen == c.menu.items-1
		switch {
		case up, down && back:
			return scene.Replace(c.back)
		case down:
			c.showLessons(c.menu.chosen, 0)
		}
	case c.inside:
		c.branch.move()
		switch {
		case up:
			c.inside = false
		case down:
			return c.play(c.id(), c.branches(c.id())[c.branch.chosen] == branchPractise)
		}
	default:
		c.menu.move()
		back := c.menu.chosen == c.menu.items-1
		switch {
		case up, down && back:
			c.showChapters(c.chapter)
		case down && len(c.branches(c.id())) == 1:
			return c.play(c.id(), false) // a single branch: no tree to go down
		case down:
			c.goDown(branchLesson)
		}
	}
	return scene.Stay
}

func (c *course) Draw(dst *ebiten.Image) {
	var items []entry
	if c.chapter < 0 {
		for i, ch := range c.open {
			items = append(items, entry{label: c.lang.T(chapterPhrase(ch.ID)), done: c.done.Finished(lessons.Course[i])})
		}
		items = append(items, entry{label: c.lang.T(msgOptBack)})
		c.drawList(dst, &c.walker, c.walker.Gait(true), c.lang.T(msgMenuLearn), c.fonts.count, items, c.menu.chosen, c.lang.T(msgCourseKeys), chaptersLayout)
		return
	}
	// The lessons, the one chosen with its branches under it.
	chosen := c.menu.chosen
	for i, l := range c.open[c.chapter].Lessons {
		if i == c.menu.chosen {
			chosen = len(items)
		}
		items = append(items, entry{label: c.lang.T(lessonPhrase(l.ID)), done: c.done[l.ID]})
		if i != c.menu.chosen {
			continue
		}
		for j, b := range c.branches(l.ID) {
			if c.inside && j == c.branch.chosen {
				chosen = len(items)
			}
			items = append(items, entry{label: c.lang.T(c.branchPhrase(l.ID, b)), branch: true})
		}
	}
	if c.menu.chosen == c.menu.items-1 {
		chosen = len(items)
	}
	items = append(items, entry{label: c.lang.T(msgOptBack)})
	heading := c.lang.T(chapterPhrase(c.open[c.chapter].ID))
	c.drawList(dst, &c.walker, c.walker.Gait(true), heading, c.fonts.heading, items, chosen, c.lang.T(msgCourseKeys), lessonsLayout)
}

// branchPhrase names branch `b` of lesson `id`.
func (c *course) branchPhrase(id string, b int) string {
	switch {
	case b == branchPractise:
		return msgCoursePractise
	case c.done[id]:
		return msgCourseReview
	}
	return msgCourseStart
}
