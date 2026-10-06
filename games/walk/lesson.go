package main

import (
	"image/color"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/walk/lessons"
)

// The lesson's screen looks like the game's (see "L'écran" in
// docs/debutants.md): the walker in his place on the left, his bubble
// above the keyboard on the right, the keyboard below where the game
// has it, reduced to three octaves to read better, its white keys named
// in both notations.
const (
	lessonLow  = 48 // C3
	lessonHigh = 84 // C6

	// The bubble: a block of lines right of the walker, and a stroke
	// from under it down to his head.
	bubbleX     = 190
	bubbleTop   = 60
	bubbleW     = screenWidth - margin - bubbleX
	bubbleLineH = 20

	// labelBottom puts the two names of a white key above its bottom
	// edge, the French one over the letter.
	labelBottom = 3

	// The walker's phrases: one note after the other, held almost to
	// the next.
	phraseNote = 500 * time.Millisecond
	phraseHold = 0.9
	phraseVel  = 0.7

	// The sign of a right answer: the walker snaps, and says so, for a
	// moment.
	cheerTime = 400 * time.Millisecond
	cheerVel  = 0.9
)

// cheers are what the walker says to a right answer, one at random,
// never the same twice in a row.
var cheers = []string{msgCheerYes, msgCheerGood, msgCheerRight}

var (
	// shownTint lights the keys the walker shows, and the answer after
	// a miss: the demo's blue, lighter, under what the hands play.
	shownTint = color.RGBA{0xbf, 0xdc, 0xf2, 0xff}

	// The names of the white keys, from C, in both notations.
	solfege = []string{"do", "ré", "mi", "fa", "sol", "la", "si"}
	letters = []string{"C", "D", "E", "F", "G", "A", "B"}
	whitePC = []int{0, 2, 4, 5, 7, 9, 11}
)

// lesson is the scene of a lesson: the walker teaches, through the
// lessons.Stage it lends to the lesson's Runner. Any key of the
// computer's keyboard says a bubble is read, but R, which plays the step
// again, and Esc, which goes back to the course without marking the
// lesson done.
type lesson struct {
	*app

	back   *course
	id     string
	runner *lessons.Runner
	piano  screen.Piano
	down   [screen.MIDIKeys]bool // what the hands hold

	bubble []string // the walker's phrase, in lines
	shown  []int    // the pitch classes lit

	// The walker's phrase playing: its keys, from `phraseAt`; `playing`
	// until the runner is told it ended.
	phrase   []int
	phraseAt time.Time
	playing  bool

	upAt   time.Time // when the last key went up
	missAt time.Time // the last casserole

	// The walker's reactions to a key, his cheer, a phrase played again,
	// wait for it to go up, or for a beat after it went down: answering
	// a key still held cuts the player short. The casserole does not
	// wait: it is the sound of the wrong note itself.
	held    int       // the keys down
	pressAt time.Time // when the last one went down
	pending []func()  // the reactions waiting, in order

	cheerAt time.Time // the last right answer
	cheer   int       // the last cheer said, in cheers

	walker walker  // for his walkout, in profile
	out    walkout // the easter egg (see walkout.go)

	turned [screen.MIDIKeys]bool // the keys that turned a page, silent, until they go up
}

func newLesson(a *app, back *course, id string, steps []lessons.Step) *lesson {
	p := newPiano()
	p.Low, p.High = lessonLow, lessonHigh
	l := &lesson{app: a, back: back, id: id, piano: p}
	l.runner = lessons.NewRunner(l, steps)
	return l
}

// Enter stops the jam: a lesson of chapter 1 plays without tempo, and
// the walker's phrases must be heard alone.
func (l *lesson) Enter() {
	if l.jam != nil {
		l.jam.stop(time.Now())
		l.jam = nil
	}
	l.runner.Start()
}

// Leave gives the keys their sound back.
func (l *lesson) Leave() { l.band.silent.Store(false) }

// The lessons.Stage the lesson's runner plays on.

func (l *lesson) Say(phrase string) {
	l.bubble = l.wrap(l.lang.T(phrase), bubbleW)
}

func (l *lesson) Light(pcs []int) { l.shown = pcs }

// react does `f` at once, or once the key that brought it about is up,
// or a beat after it went down.
func (l *lesson) react(f func()) {
	if l.held == 0 {
		f()
		return
	}
	l.pending = append(l.pending, f)
}

// Play plays `keys` from now, scheduled ahead as the band is, or a beat
// after a casserole that just sounded: the phrase played again after a
// wrong note must not cover it.
func (l *lesson) Play(keys []int) { l.react(func() { l.play(keys) }) }

func (l *lesson) play(keys []int) {
	now := time.Now().Add(lookahead)
	if after := l.missAt.Add(phraseNote); now.Before(after) {
		now = after
	}
	for i, k := range keys {
		at := now.Add(time.Duration(i) * phraseNote)
		l.band.piano.ScheduleOn(k, phraseVel, at)
		l.band.piano.ScheduleOff(k, at.Add(time.Duration(phraseHold*float64(phraseNote))))
	}
	l.phrase, l.phraseAt, l.playing = keys, now, true
}

// Hit snaps, has the walker raise his free hand, snapping, and says a
// cheer in his bubble, until the next step says its own phrase.
func (l *lesson) Hit() { l.react(l.hit) }

func (l *lesson) hit() {
	now := time.Now()
	l.band.snap.ScheduleOn(noiseKey, cheerVel, now)
	l.band.snap.ScheduleOff(noiseKey, now.Add(hold)) // a clip ignores it
	l.cheerAt = now
	l.cheer = (l.cheer + 1 + rand.IntN(len(cheers)-1)) % len(cheers)
	l.bubble = []string{l.lang.T(cheers[l.cheer])}
}

// Miss sounds the casserole at once, with the wrong note, and shows
// the right keys.
func (l *lesson) Miss(pcs []int) {
	l.missAt = time.Now()
	l.band.strike(casserole, casseroleVel, l.missAt)
	l.shown = pcs
}

func (l *lesson) Update() scene.Transition {
	now := time.Now()
	l.drain(func(e keyboard.Event) {
		if l.turned[e.Key] {
			l.turned[e.Key] = e.Down // up at last: a key like the others
			return
		}
		if e.Down && l.teaching() && l.runner.Reading() {
			// A bubble to read: the key turns the page, unheard (see
			// band.silent) and unlit, as a key of the computer would.
			l.turned[e.Key] = true
			l.runner.Read()
			return
		}
		if l.down[e.Key] == e.Down {
			return // a repeat, or a key held across scenes
		}
		l.down[e.Key] = e.Down
		if !e.Down {
			l.held--
			l.upAt = now
			l.runner.NoteOff(e.Key) // always: the runner knows what is up
			return
		}
		l.held++
		l.pressAt = now
		if !l.teaching() {
			l.walkoutKey(now) // the key sounds, and counts for nothing
			return
		}
		l.runner.NoteOn(e.Key)
	})
	l.walkoutTick(now)
	if len(l.pending) > 0 && (l.held == 0 || now.Sub(l.pressAt) >= phraseNote) {
		for _, f := range l.pending {
			f()
		}
		l.pending = nil
	}
	// A step played through hands over once the keys are up, and a beat
	// of the walker's phrases later: asked again at once, the player
	// would feel rushed.
	if l.teaching() && l.runner.Waiting() && !l.runner.Held() && now.Sub(l.upAt) >= phraseNote {
		l.runner.Resume()
	}
	if l.playing && now.After(l.phraseAt.Add(time.Duration(len(l.phrase))*phraseNote)) {
		l.playing = false
		l.runner.PhraseEnded()
	}
	keys := inpututil.AppendJustPressedKeys(nil)
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return scene.Replace(l.back)
	case !l.teaching():
		if len(keys) > 0 {
			l.walkoutKey(now)
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyR):
		l.runner.Again()
	case len(keys) > 0:
		l.runner.Read()
	}
	// The keys turn the page in silence while a bubble waits to be read;
	// set here, before the next key comes.
	l.band.silent.Store(l.teaching() && l.runner.Reading())
	if l.runner.Done() {
		l.back.finish(l.id)
		return scene.Replace(l.back)
	}
	var walker [screen.MIDIKeys]bool // the note of his phrase sounding now
	if l.playing && !now.Before(l.phraseAt) {
		if i := int(now.Sub(l.phraseAt) / phraseNote); i < len(l.phrase) {
			walker[l.phrase[i]] = true
		}
	}
	l.piano.Update(screen.Lit{Down: &l.down, Color: playerLit}, screen.Lit{Down: &walker, Color: demoLit})
	return scene.Stay
}

func (l *lesson) Draw(dst *ebiten.Image) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: l.scale}
	now := time.Now()
	if !l.drawWalkout(dst, c, now) {
		drawTeacher(c, walkerX, walkerY, walkerScale, now.Sub(l.cheerAt) < cheerTime, ink)
	}
	if len(l.bubble) > 0 {
		l.drawSpeech(c, l.bubble, bubbleX, bubbleTop, bubbleW, bubbleLineH, l.walkerAt(now))
	}

	l.piano.Draw(c, func(k int) screen.Dress {
		if slices.Contains(l.shown, k%12) {
			return screen.Dress{Base: shownTint}
		}
		return screen.Dress{}
	})
	l.drawNames(c)
	keys := msgLessonKeys
	if l.teaching() && !l.runner.Reading() {
		keys = msgLessonPlay // the walker waits for notes, not for a page turned
	}
	c.Text(l.lang.T(keys), l.fonts.ui, margin, statusY, faint)
}

// drawNames writes the names of the white keys on them, "do" over "C".
func (l *lesson) drawNames(c screen.Canvas) {
	for k := l.piano.Low; k <= l.piano.High; k++ {
		i := slices.Index(whitePC, k%12)
		if i < 0 {
			continue
		}
		x, y, w, h := l.piano.KeyRect(k)
		mid := float64(x + w/2)
		_, lh := c.Measure(letters[i], l.fonts.ui)
		bottom := float64(y+h) - labelBottom - lh
		c.Centred(letters[i], l.fonts.ui, mid, bottom, faint)
		c.Centred(solfege[i], l.fonts.ui, mid, bottom-lh, faint)
	}
}
