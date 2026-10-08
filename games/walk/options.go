package main

import (
	"log"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/calibrate"
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/walk/band"
	"github.com/ArnaudCalmettes/gohar/games/walk/figure"
)

// languages are those the game speaks, in the order the options go
// round them.
var languages = []string{"fr", "en"}

// The items of the options, in order.
const (
	optLang = iota
	optCalibrate
	optBack
	optItems
)

// options is the options scene, over the jam lighter: the bass and the
// hi-hat on 2 and 4 only, the walker grooving without his snaps. The
// left and right arrows change a value, which takes effect at once; the
// options are saved on the way back to the title.
type options struct {
	*app

	back   *title // the title to go back to, the jam still playing
	walker figure.Walker
	menu   menu
}

func newOptions(a *app, back *title) *options {
	return &options{app: a, back: back, walker: figure.NewWalker(figure.Grooving), menu: menu{items: optItems}}
}

func (o *options) Enter() {}

// Leave saves the options.
func (o *options) Leave() {
	if err := (prefs{Lang: o.lang.Tag()}).save(); err != nil {
		log.Println("options:", err)
	}
}

func (o *options) Update() scene.Transition {
	o.drain(nil)
	o.menu.move()
	// A tap on a value changes it, down left of its middle, up right of
	// it, as the arrows; on the others, it opens them.
	tp, tapped := o.tapped()
	if tapped {
		o.menu.chosen = tp.item
	}
	value := o.menu.chosen == optLang
	step := 0
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return scene.Replace(o.back)
	case tapped && value:
		step = tp.side
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
		step = -1
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight):
		step = 1
	case confirmed() || tapped:
		switch o.menu.chosen {
		case optCalibrate:
			return o.calibrate()
		case optBack:
			return scene.Replace(o.back)
		}
	}
	if step != 0 {
		o.change(step)
	}
	o.jam.Play(time.Now(), band.Light)
	return scene.Stay
}

// change moves the value chosen one `step` along.
func (o *options) change(step int) {
	switch o.menu.chosen {
	case optLang:
		i := slices.Index(languages, o.lang.Tag())
		tag := languages[(i+step+len(languages))%len(languages)]
		l, err := newLang(tag)
		if err != nil {
			log.Println("options:", err)
			return
		}
		o.lang = l
	}
}

// calibrate opens the shared calibration on the jam's clock: the bass
// alone under its bar, then the lighter jam again once the measure is
// steady, and back to these options.
func (o *options) calibrate() scene.Transition {
	s, err := calibrate.New(calibrate.Config{
		Width: screenWidth,
		Scale: func() float64 { return o.scale },
		Lang:  o.lang.Tag(),
		Clock: o.jam.Metronome(),
		MIDI:  o.midi != "",
		Drain: o.drain,
		Tick:  o.band.Tick,
		Play: func(now time.Time, quiet bool) {
			if quiet {
				o.jam.Play(now, band.BassOnly)
				return
			}
			o.jam.Play(now, band.Light)
		},
		Key:  o.pair,
		Done: func(d time.Duration) { o.latency = d },
		Back: func() scene.Scene { return o },
	})
	if err != nil {
		log.Println("calibration:", err)
		return scene.Stay
	}
	return scene.Replace(s)
}

func (o *options) Draw(dst *ebiten.Image) {
	labels := make([]string, optItems)
	labels[optLang] = o.lang.T(msgOptLang, "Name", o.lang.T(msgLangName))
	labels[optCalibrate] = o.lang.T(msgOptCalibrate)
	labels[optBack] = o.lang.T(msgOptBack)
	o.drawMenu(dst, &o.walker, o.walker.Gait(true), o.lang.T(msgMenuOptions), o.fonts.count, labels, o.menu.chosen, o.lang.T(msgOptKeys))
}
