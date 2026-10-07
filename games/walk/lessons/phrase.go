package lessons

import "slices"

// The steps that ask for phrases, several notes in a row: Repeat, Hang,
// its phrase left hanging, and PlayLine, a line of a grid.

// Repeat is "Répéter après lui": the walker plays a phrase, then the
// player plays it back, note after note, without a mistake. A wrong
// note sounds the casserole, and the walker plays the phrase again, to
// be played back again from its start: "do sol do" is not found among
// the notes of a scale.
//
// Without Keys, the walker only says the phrase, "C, E, C", and the
// player plays it from the words. A wrong note then shows every key of
// the phrase, lit until it is played through, from its start again.
type Repeat struct {
	Phrase string
	Keys   []int    // what the walker plays, MIDI numbers; nil to say it only
	Want   []Target // a target for each note

	at int // the note expected next, -1 while the walker plays
}

func (st *Repeat) said() string { return st.Phrase }

func (st *Repeat) begin(s Stage) {
	s.Say(st.Phrase)
	s.Light(nil)
	if st.Keys == nil {
		st.at = 0
		return
	}
	st.at = -1
	s.Play(st.Keys)
}

func (st *Repeat) phraseEnded(Stage) bool {
	if st.Keys != nil {
		st.at = 0
	}
	return false
}

func (st *Repeat) pressed(s Stage, held []int) bool {
	if st.at < 0 {
		return false // he is still playing: the keys sound, nothing counts
	}
	if st.Want[st.at].Judge(held) != Hit {
		if st.Keys == nil {
			s.Miss(st.hint(), st.at)
			st.at = 0
			return false
		}
		s.Miss(nil, st.at)
		st.at = -1
		s.Play(st.Keys)
		return false
	}
	if st.Keys != nil {
		s.Light(nil) // what a Show lit, out at the first note
	}
	st.at++
	if st.at < len(st.Want) {
		return false
	}
	if st.Keys == nil {
		s.Light(nil) // the phrase's keys, lit after a miss
	}
	s.Hit()
	return true
}

// hint is every pitch class the phrase wants, once each, in order.
func (st *Repeat) hint() []int {
	var pcs []int
	for _, w := range st.Want {
		for _, p := range w.Hint() {
			if !slices.Contains(pcs, p) {
				pcs = append(pcs, p)
			}
		}
	}
	return pcs
}

func (*Repeat) released(Stage, []int) bool { return false }
func (*Repeat) read() bool                 { return false }

// Hang is a phrase left hanging, the gag of lesson 1.2 (see
// docs/debutants/chapitre-1.md): said, and played back from the words
// as a Repeat without Keys is, but its last note is held by the game
// once the player lets it go. The walker frets, and the note that
// resolves it, the first one of its pitch class above, blinks; played,
// it lets go of the held note, and he thanks the player.
type Hang struct {
	Phrase  string   // what is asked, said
	Want    []Target // a target for each note, the last one left hanging
	Humpf   string   // said while it hangs
	Resolve Note     // what resolves it
	Thanks  string   // said once resolved

	at   int // the note expected next
	held int // the key left hanging, 0 before
}

func (st *Hang) said() string { return st.Phrase }

func (st *Hang) alsoSays() []string { return []string{st.Humpf, st.Thanks} }

func (st *Hang) begin(s Stage) {
	if st.held != 0 { // played over, R: nothing left hanging
		s.Release()
		s.Blink(0)
	}
	st.at, st.held = 0, 0
	s.Say(st.Phrase)
	s.Light(nil)
}

func (st *Hang) pressed(s Stage, held []int) bool {
	if st.held == 0 {
		r := Repeat{Want: st.Want}
		if st.Want[st.at].Judge(held) != Hit {
			s.Miss(r.hint(), st.at)
			st.at = 0
			return false
		}
		if st.at++; st.at < len(st.Want) {
			return false
		}
		s.Light(nil) // the phrase's keys, lit after a miss
		st.held = held[len(held)-1]
		s.Hold(st.held)
		s.Say(st.Humpf)
		up := (int(st.Resolve) - pc(st.held) + 12) % 12
		if up == 0 {
			up = 12
		}
		s.Blink(st.held + up)
		return false
	}
	if st.Resolve.Judge(held) != Hit {
		s.Miss(nil, 0)
		return false
	}
	s.Blink(0)
	s.Release()
	s.Say(st.Thanks)
	return true
}

func (*Hang) released(Stage, []int) bool { return false }
func (*Hang) phraseEnded(Stage) bool     { return false }
func (*Hang) read() bool                 { return false }

// PlayLine plays a line of a grid on the bass, the roots of its chords,
// a note a bar. With Keys, the walker plays it, the bar he plays shaded,
// over when his phrase ends. Without, the player plays `Want`, a target
// a bar, the line waiting for him, without tempo, as the first palier
// does: a wrong root sounds the casserole and shows the right one, and
// the line stays on its bar.
type PlayLine struct {
	Phrase string
	Chords []string
	Keys   []int    // the walker's roots, MIDI numbers; nil for the player's turn
	Want   []Target // the player's, a target a bar

	at int // the bar played next
}

func (st *PlayLine) said() string { return st.Phrase }

func (st *PlayLine) begin(s Stage) {
	s.Chords(nil)
	s.Bass(true)
	s.Say(st.Phrase)
	s.Light(nil)
	if st.Keys != nil {
		s.Line(st.Chords, -1)
		s.PlayBars(st.Keys)
		return
	}
	st.at = 0
	s.Line(st.Chords, 0)
}

func (st *PlayLine) pressed(s Stage, held []int) bool {
	if st.Keys != nil {
		return false // he plays: the keys sound, nothing counts
	}
	if st.Want[st.at].Judge(held) != Hit {
		s.Miss(st.Want[st.at].Hint(), 0) // a bar of its own: shown on the next one
		return false
	}
	s.Light(nil)
	if st.at++; st.at < len(st.Want) {
		s.Line(st.Chords, st.at)
		return false
	}
	s.Line(st.Chords, -1)
	s.Hit()
	return true
}

func (st *PlayLine) phraseEnded(Stage) bool  { return st.Keys != nil }
func (*PlayLine) released(Stage, []int) bool { return false }
func (*PlayLine) read() bool                 { return false }
