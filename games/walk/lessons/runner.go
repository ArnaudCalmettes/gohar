package lessons

import "slices"

// A Runner plays a lesson's steps on a Stage, one after the other. The
// game feeds it what happens, the keys, the end of the walker's
// phrases, the bubbles read, and it moves on when a step is over.
//
// A step the player ends at the keyboard does not hand over at once:
// the Runner waits, Waiting, until the game calls Resume, once the keys
// are up and a beat has gone by. Asked again straight away, the player
// feels rushed.
type Runner struct {
	stage   Stage
	steps   []Step
	at      int   // the step playing, len(steps) once over
	held    []int // the keys down, in the order they went down
	waiting bool  // the step before is over, the next not begun
}

// NewRunner prepares `steps` on `stage`. Start begins the first.
func NewRunner(stage Stage, steps []Step) *Runner {
	return &Runner{stage: stage, steps: steps}
}

// Start begins the first step.
func (r *Runner) Start() {
	r.at = 0
	r.begin()
}

// Done tells whether every step is over.
func (r *Runner) Done() bool { return r.at >= len(r.steps) }

// Waiting tells whether the Runner waits for Resume to begin the next
// step. Meanwhile the keys sound, and count for nothing.
func (r *Runner) Waiting() bool { return r.waiting && !r.Done() }

// Held tells whether a key is down.
func (r *Runner) Held() bool { return len(r.held) > 0 }

// Resume begins the step the Runner waits for.
func (r *Runner) Resume() {
	if r.waiting {
		r.waiting = false
		r.begin()
	}
}

// Again plays the step playing over: the bubble said again, the phrase
// heard again (the R key of docs/debutants.md).
func (r *Runner) Again() {
	r.waiting = false
	r.begin()
}

// NoteOn tells that MIDI key `key` went down.
func (r *Runner) NoteOn(key int) {
	r.held = append(r.held, key)
	if r.waiting {
		return
	}
	if r.next(func(st Step) bool { return st.pressed(r.stage, r.held) }) {
		r.waiting = true
	}
}

// NoteOff tells that MIDI key `key` went up.
func (r *Runner) NoteOff(key int) {
	r.held = slices.DeleteFunc(r.held, func(k int) bool { return k == key })
	if !r.waiting && r.next(func(st Step) bool { return st.released(r.stage, r.held) }) {
		r.begin()
	}
}

// PhraseEnded tells that the walker's phrase is over.
func (r *Runner) PhraseEnded() {
	if !r.waiting && r.next(func(st Step) bool { return st.phraseEnded(r.stage) }) {
		r.begin()
	}
}

// Read tells that the bubble is read: the player pressed a key of the
// computer's keyboard, or the game judged its reading time over.
func (r *Runner) Read() {
	if !r.waiting && r.next(func(st Step) bool { return st.read() }) {
		r.begin()
	}
}

// next moves past the step playing when `over` says it is over, and
// tells whether it did. The caller begins the next one, or waits.
func (r *Runner) next(over func(Step) bool) bool {
	if r.Done() || !over(r.steps[r.at]) {
		return false
	}
	r.at++
	return true
}

// begin begins the step playing, if any.
func (r *Runner) begin() {
	if !r.Done() {
		r.steps[r.at].begin(r.stage)
	}
}
