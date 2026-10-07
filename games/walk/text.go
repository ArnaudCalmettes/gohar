package main

import (
	"embed"
	"io/fs"

	"github.com/ArnaudCalmettes/gohar/games/lang"
	"github.com/ArnaudCalmettes/gohar/games/walk/figure"
	"github.com/ArnaudCalmettes/gohar/games/walk/mark"
)

// The game's phrases, one file per language in locales/ (see
// "L'internationalisation" in docs/architecture.md). The musical words,
// note and chord names, come from naming.
//
//go:embed locales/*.toml
var locales embed.FS

// newLang speaks `tag`, "fr" or "en", English for any other.
func newLang(tag string) (*lang.Lang, error) {
	files, err := fs.Sub(locales, "locales")
	if err != nil {
		return nil, err
	}
	return lang.New(files, tag)
}

// The IDs of the phrases. Every one is in each file of locales/, which
// TestPhrases checks: call them through these names only.
const (
	msgTempo      = "tempo"      // {{.BPM}}
	msgFreeTempo  = "tempo.free" // the chart waits for the root
	msgKeyPlay    = "keys.play"  // the keys of the game, each as it stands
	msgKeyStop    = "keys.stop"
	msgKeyFree    = "keys.free"  // T, to the phase without tempo
	msgKeyTempo   = "keys.tempo" // T, back to the tempo
	msgKeyDemoOn  = "keys.demo.on"
	msgKeyDemoOff = "keys.demo.off"
	msgKeyGrid    = "keys.grid" // the arrows: the grid and the tempo
	msgKeyMenu    = "keys.menu"
	msgNoMIDI     = "status.nomidi"
	msgMark       = "status.mark" // {{.Note}}, {{.Words}}: what a note was

	msgModeTempo = "mode.tempo"
	msgModeFree  = "mode.free"
	msgModeDemo  = "mode.demo"

	msgOnTime  = "timing.ontime"
	msgEarly   = "timing.early"
	msgLate    = "timing.late"
	msgBetween = "timing.between"

	msgRoot      = "pitch.root"
	msgChordTone = "pitch.chordtone"
	msgOutside   = "pitch.outside"

	msgPlayer = "record.player" // who plays, in a recording's heading
	msgDemo   = "record.demo"

	msgMenuLearn   = "menu.learn" // the course for beginners
	msgMenuPlay    = "menu.play"  // play a grid
	msgMenuOptions = "menu.options"
	msgMenuQuit    = "menu.quit"
	msgTitleKeys   = "title.keys" // the keys of the title screen

	msgOptLang      = "options.lang" // {{.Name}}: the language spoken
	msgLangName     = "lang.name"    // each language's own name for itself
	msgOptCalibrate = "options.calibrate"
	msgOptBack      = "options.back"
	msgOptKeys      = "options.keys"

	msgCourseKeys     = "course.keys"     // the chapters and lessons are named in course.go
	msgCoursePractise = "course.practise" // a lesson done: its activity
	msgCourseReview   = "course.review"   // the lesson again
	msgCourseStart    = "course.start"    // the lesson, not done yet
	msgLessonKeys     = "lesson.keys"     // a bubble to read; what a lesson says comes from lessons, by ID
	msgLessonPlay     = "lesson.play"     // notes to play
	msgCheerYes       = "cheer.yes"       // a right answer in a lesson
	msgCheerGood      = "cheer.good"
	msgCheerRight     = "cheer.right"

	// The walker's annoyance, the easter egg of the lessons (gags.go).
	msgTeaseHey       = "tease.hey"
	msgTeaseOnPurpose = "tease.onpurpose"
	msgTeaseFedUp     = "tease.fedup"
	msgTeaseBack      = "tease.back"
	msgTeaseSit       = "tease.sit"
	msgTeaseLotus     = "tease.lotus"
	msgTeaseThanks    = "tease.thanks" // out of the lotus

	// What the walker says (see figure.Coach).
	msgRushing  = "coach.rushing"
	msgRelax    = "coach.relax"
	msgDragging = "coach.dragging"
	msgItDrags  = "coach.itdrags"
	msgCool     = "coach.cool"
	msgYeah     = "coach.yeah"
	msgKeepItUp = "coach.keepitup"
	msgSwinging = "coach.swinging"
	msgGroovy   = "coach.groovy"
	msgGreat    = "coach.great"
	msgIDig     = "coach.idig"

	// The review of a run.
	msgReviewTitle       = "review.title" // {{.Title}}: the grid's
	msgReviewWorked      = "review.worked"
	msgReviewConsolidate = "review.consolidate"
	msgReviewStart       = "review.start" // {{.Landed}}, {{.Expected}}: the beats
	msgReviewMid         = "review.mid"
	msgReviewHeld        = "review.held"
	msgReviewTime        = "review.time"
	msgReviewSteady      = "review.steady"
	msgReviewRushing     = "review.rushing"
	msgReviewDragging    = "review.dragging"
	msgReviewUnsteady    = "review.unsteady"
	msgReviewKeys        = "review.keys"

	// The advice of a review, one priority.
	msgAdviceWork      = "advice.work"   // {{.What}}: a formula or a situation, below
	msgAdviceSlower    = "advice.slower" // {{.BPM}}: the tempo to try
	msgAdviceFaster    = "advice.faster" // {{.BPM}}
	msgAdviceStay      = "advice.stay"
	msgFormulaAnatole  = "formula.anatole"
	msgFormulaThreeSix = "formula.threesix"
	msgFormulaAeolian  = "formula.aeolian"
	msgFormulaTwoFive  = "formula.twofive"
	msgWorkStart       = "work.start"
	msgWorkMid         = "work.mid"
	msgWorkHeld        = "work.held"
)

var phrases = []string{
	msgTempo, msgFreeTempo, msgNoMIDI, msgMark,
	msgKeyPlay, msgKeyStop, msgKeyFree, msgKeyTempo, msgKeyDemoOn, msgKeyDemoOff, msgKeyGrid, msgKeyMenu,
	msgModeTempo, msgModeFree, msgModeDemo,
	msgOnTime, msgEarly, msgLate, msgBetween,
	msgRoot, msgChordTone, msgOutside,
	msgPlayer, msgDemo,
	msgMenuLearn, msgMenuPlay, msgMenuOptions, msgMenuQuit, msgTitleKeys,
	msgOptLang, msgLangName, msgOptCalibrate, msgOptBack, msgOptKeys,
	msgCourseKeys, msgCoursePractise, msgCourseReview, msgCourseStart, msgLessonKeys, msgLessonPlay, msgCheerYes, msgCheerGood, msgCheerRight,
	msgTeaseHey, msgTeaseOnPurpose, msgTeaseFedUp, msgTeaseBack, msgTeaseSit, msgTeaseLotus, msgTeaseThanks,
	msgRushing, msgRelax, msgDragging, msgItDrags,
	msgCool, msgYeah, msgKeepItUp, msgSwinging, msgGroovy, msgGreat, msgIDig,
	msgReviewTitle, msgReviewWorked, msgReviewConsolidate,
	msgReviewStart, msgReviewMid, msgReviewHeld,
	msgReviewTime, msgReviewSteady, msgReviewRushing, msgReviewDragging, msgReviewUnsteady, msgReviewKeys,
	msgAdviceWork, msgAdviceSlower, msgAdviceFaster, msgAdviceStay,
	msgFormulaAnatole, msgFormulaThreeSix, msgFormulaAeolian, msgFormulaTwoFive,
	msgWorkStart, msgWorkMid, msgWorkHeld,
}

// coachLines are what the coach says.
var coachLines = figure.Lines{
	Rushing:  []string{msgRushing, msgRelax},
	Dragging: []string{msgDragging, msgItDrags},
	Praise:   []string{msgCool, msgYeah, msgKeepItUp, msgSwinging, msgGroovy, msgGreat, msgIDig},
}

var (
	timingPhrase = map[mark.Timing]string{mark.OnTime: msgOnTime, mark.Early: msgEarly, mark.Late: msgLate, mark.Between: msgBetween}
	pitchPhrase  = map[mark.Pitch]string{mark.Root: msgRoot, mark.ChordTone: msgChordTone, mark.Outside: msgOutside}
)
