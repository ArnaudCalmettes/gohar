package calibrate

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/games/lang"
	"github.com/ArnaudCalmettes/gohar/games/scene"
	"github.com/ArnaudCalmettes/gohar/games/screen"
	"github.com/ArnaudCalmettes/gohar/games/settings"
	"github.com/ArnaudCalmettes/gohar/games/tempo"
)

// Config is what a game lends the scene: its screen, its keyboard, its
// sound, and where to go back.
type Config struct {
	Width float64        // the logical screen's
	Scale func() float64 // the window's scale, as the game's layout set it
	Lang  string         // "fr" or "en"

	// Clock is the game's metronome, when music plays already: the bar
	// falls on its beats, in its tempo, and the music goes on under it.
	// Without one, the zero value, the scene beats on its own, at BPM,
	// or `defaultBPM` for 0.
	Clock tempo.Metronome
	BPM   float64

	// MIDI tells whether a keyboard listens; Drain hands each key event
	// waiting to its argument, as the game receives them.
	MIDI  bool
	Drain func(func(keyboard.Event))

	// Tick schedules the sound of one beat at `at`: "ta", or "TI" when
	// `accent`, the fourth beat.
	Tick func(at time.Time, accent bool)

	// Play, if not nil, schedules the game's music on Clock at each
	// tick, up to a lookahead from `now`: `quiet` while the bar beats,
	// the music discreet under it, and no longer once the measure is
	// steady. The scene left, the next one takes the music over.
	Play func(now time.Time, quiet bool)

	// Key names the pair of keyboard and output measured (see Key): the
	// offset is saved for it once steady, unless it is empty. Done, if
	// not nil, hands the offset to the game at once.
	Key  string
	Done func(offset time.Duration)

	// Back builds the scene to go back to.
	Back func() scene.Scene
}

// The pace of the bar, unless the game sets it, and how far ahead it is
// scheduled (see "Le son" in docs/walk.md). Only the fourth beat is
// pressed: a lively pace stays easy.
const (
	defaultBPM = 120
	perBar     = 4
	lookahead  = 100 * time.Millisecond
)

// The layout, in logical units: four circles in a row across the
// middle, the text below.
const (
	circleY   = 140 // their centres
	circleGap = 80  // from one centre to the next
	radius    = 22
	pop       = 10 // how much the beat sounding grows, fading over the beat
	titleY    = 30
	askY      = 220
	tapY      = 250
	measureY  = 270
	savedY    = 290
	keysY     = 330
	bigSize   = 28 // the numbers in the circles
	textSize  = 12
)

var (
	paper  = color.White
	ink    = color.Black
	faint  = color.RGBA{0x88, 0x88, 0x88, 0xff}
	beatIn = color.RGBA{0x8c, 0xc8, 0x5a, 0xff} // "ta", green as in the original
	tiIn   = color.RGBA{0xf0, 0x7c, 0xb8, 0xff} // "TI", pink
)

// Scene beats "ta, ta, ta, TI" until the taps on the TI are steady,
// then shows the offset measured. Not saved yet: the next step keeps it
// for the pair of keyboard and audio output.
type Scene struct {
	cfg   Config
	lang  *lang.Lang
	big   *screen.Font
	small *screen.Font

	m      tempo.Metronome
	meter  *Meter
	next   int // the next beat to schedule
	last   time.Duration
	tapped bool
	result time.Duration
	done   bool
	saved  bool // the offset is in the file
}

func New(cfg Config) (*Scene, error) {
	l, err := newLang(cfg.Lang)
	if err != nil {
		return nil, err
	}
	big, err := screen.NewFont(bigSize, screen.GoRegular)
	if err != nil {
		return nil, err
	}
	small, err := screen.NewFont(textSize, screen.GoRegular)
	if err != nil {
		return nil, err
	}
	return &Scene{cfg: cfg, lang: l, big: big, small: small}, nil
}

// Enter starts the bar on the first downbeat of the clock not yet
// scheduled, or a lookahead from now on a clock of its own.
func (s *Scene) Enter() {
	now := time.Now()
	s.m = s.cfg.Clock
	if s.m.PerBar() == 0 {
		bpm := s.cfg.BPM
		if bpm == 0 {
			bpm = defaultBPM
		}
		s.m = tempo.NewMetronome(now.Add(lookahead), bpm, perBar)
	}
	_, end := s.m.Due(now, now.Add(lookahead))
	per := s.m.PerBar()
	s.next = (max(end, 0) + per - 1) / per * per
	s.meter = NewMeter(s.m)
}

// Leave has nothing to release: the ticks are strokes, they end alone,
// and the music is the next scene's.
func (s *Scene) Leave() {}

func (s *Scene) Update() scene.Transition {
	// Once the measure is steady, any key goes back, on the computer's
	// keyboard or on the MIDI one; the tap that steadied it is a tick
	// behind, and does not count.
	finished := s.done
	back := inpututil.IsKeyJustPressed(ebiten.KeyEscape) || (finished && len(inpututil.AppendJustPressedKeys(nil)) > 0)
	s.cfg.Drain(func(e keyboard.Event) {
		if !e.Down {
			return
		}
		if finished {
			back = true
			return
		}
		if s.done {
			return
		}
		if off, ok := s.meter.Tap(e.At); ok {
			s.last, s.tapped = off, true
		}
		if mean, ok := s.meter.Mean(); ok {
			s.result, s.done = mean, true
			s.keep()
		}
	})
	if back {
		return scene.Replace(s.cfg.Back())
	}
	now := time.Now()
	_, end := s.m.Due(now, now.Add(lookahead))
	for ; s.next < end && !s.done; s.next++ { // the bar stops with the measure
		s.cfg.Tick(s.m.At(s.next), s.next%s.m.PerBar() == s.m.PerBar()-1)
	}
	if s.cfg.Play != nil {
		s.cfg.Play(now, !s.done)
	}
	return scene.Stay
}

func (s *Scene) Draw(dst *ebiten.Image) {
	dst.Fill(paper)
	c := screen.Canvas{Dst: dst, Scale: s.cfg.Scale()}
	mid := s.cfg.Width / 2
	c.Centred(s.lang.T(msgTitle), s.small, mid, titleY, ink)

	// The circles, the beat sounding grown, shrinking back over the
	// beat: the eye follows the bar as the ear does.
	x := s.m.Beats(time.Now())
	sounding, frac := -1, 0.0
	if x >= 0 && !s.done {
		sounding, frac = int(x)%perBar, x-math.Floor(x)
	}
	for i := range perBar {
		cx := mid + (float64(i)-float64(perBar-1)/2)*circleGap
		r := float64(radius)
		if i == sounding {
			r += pop * (1 - frac)
		}
		var col color.Color = beatIn
		if i == perBar-1 {
			col = tiIn
		}
		c.Circle(float32(cx), circleY, float32(r), col)
		label := fmt.Sprint(i + 1)
		_, h := c.Measure(label, s.big)
		c.Centred(label, s.big, cx, circleY-h/2, ink)
	}

	ask := s.lang.T(msgAsk)
	switch {
	case s.done:
		ask = s.lang.T(msgDone)
	case !s.cfg.MIDI:
		ask = s.lang.T(msgNoMIDI)
	}
	c.Centred(ask, s.small, mid, askY, ink)
	if s.tapped {
		c.Centred(s.lang.T(msgTap, "Ms", ms(s.last)), s.small, mid, tapY, faint)
	}
	switch {
	case s.done:
		c.Centred(s.lang.T(msgSteady, "Ms", ms(s.result)), s.small, mid, measureY, ink)
		if s.saved {
			c.Centred(s.lang.T(msgSaved), s.small, mid, savedY, faint)
		}
		return // the keys are said above
	case s.meter.Taps() > 0:
		mean, _ := s.meter.Mean()
		c.Centred(s.lang.T(msgMeasure, "Ms", ms(mean), "Taps", s.meter.Taps(), "Of", window), s.small, mid, measureY, faint)
	}
	c.Centred(s.lang.T(msgKeys), s.small, mid, keysY, faint)
}

// keep saves the offset measured for the pair, and hands it to the
// game. A file that does not write loses nothing but this measure: the
// game still applies it until it quits.
func (s *Scene) keep() {
	if s.cfg.Key != "" {
		if err := Save(settings.Path(File), s.cfg.Key, s.result); err != nil {
			log.Println("calibration:", err)
		} else {
			s.saved = true
		}
	}
	if s.cfg.Done != nil {
		s.cfg.Done(s.result)
	}
}

// ms writes `d` in signed milliseconds: +40 late, -30 early.
func ms(d time.Duration) string {
	return fmt.Sprintf("%+d", d.Milliseconds())
}
