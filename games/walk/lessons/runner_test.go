package lessons

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// The tests name keys as a musician does, C4 for MIDI 60, and black
// keys by their sharps.
var names = []string{"C", "C♯", "D", "D♯", "E", "F", "F♯", "G", "G♯", "A", "A♯", "B"}

// key is the MIDI number of `name`, "C4" or "F♯3".
func key(name string) int {
	i := len(name) - 1
	octave := int(name[i] - '0')
	return slices.Index(names, name[:i]) + 12*(octave+1)
}

func keys(names ...string) []int {
	var out []int
	for _, n := range names {
		out = append(out, key(n))
	}
	return out
}

func keyName(k int) string { return fmt.Sprintf("%s%d", names[pc(k)], k/12-1) }

// stage writes down what the lesson does, in words.
type stage struct{ log []string }

func (s *stage) Say(phrase string) { s.log = append(s.log, "says "+phrase) }

func (s *stage) Add(phrase string) { s.log = append(s.log, "adds "+phrase) }

func (s *stage) Light(pcs []int) {
	if pcs == nil {
		s.log = append(s.log, "lights out")
		return
	}
	s.log = append(s.log, "lights "+pcNames(pcs))
}

func (s *stage) Play(ks []int) {
	var ns []string
	for _, k := range ks {
		ns = append(ns, keyName(k))
	}
	s.log = append(s.log, "plays "+strings.Join(ns, " "))
}

func (s *stage) Miss(pcs []int, _ int) { s.log = append(s.log, "casserole, shows "+pcNames(pcs)) }

func (s *stage) Hit() { s.log = append(s.log, "right") }

func (s *stage) Hold(k int) { s.log = append(s.log, "holds "+keyName(k)) }

func (s *stage) Release() { s.log = append(s.log, "lets go") }

func (s *stage) Blink(k int) {
	if k == 0 {
		s.log = append(s.log, "blinks no more")
		return
	}
	s.log = append(s.log, "blinks "+keyName(k))
}

func (s *stage) PlayBars(ks []int) {
	s.Play(ks)
	s.log[len(s.log)-1] += ", a bar each"
}

func (s *stage) Chords(cs []string) { s.log = append(s.log, "writes "+strings.Join(cs, " ")) }

func (s *stage) Line(cs []string, bar int) {
	s.log = append(s.log, fmt.Sprintf("line %s, bar %d", strings.Join(cs, " | "), bar))
}

func (s *stage) Bass(on bool) { s.log = append(s.log, fmt.Sprint("bass ", on)) }

func (s *stage) Tease(a Annoyance) {
	s.log = append(s.log, [...]string{Hey: "Hé…", OnPurpose: "on purpose?", FedUp: "walks out"}[a])
}

func pcNames(pcs []int) string {
	var ns []string
	for _, p := range pcs {
		ns = append(ns, names[p])
	}
	return strings.Join(ns, " ")
}

// take returns what happened since the last take.
func (s *stage) take() []string {
	out := s.log
	s.log = nil
	return out
}

func expect(t *testing.T, s *stage, want ...string) {
	t.Helper()
	if got := s.take(); !slices.Equal(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

// The opening of lesson 1.1, in short: the pairs shown, then pressed
// together, then a C found, and a motif repeated.
func TestRunner(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{
		Say{Phrase: "hello"},
		Show{Phrase: "pairs", PCs: Pair},
		&Ask{Phrase: "press a pair", Want: Pair},
		&Ask{Phrase: "find C", Want: Do},
		&Repeat{Phrase: "after me", Keys: keys("C4", "D4", "C4"), Want: []Target{Do, Re, Do}},
	})
	r.Start()
	expect(t, s, "says hello")

	r.NoteOn(key("C4")) // a key during a bubble: it sounds, nothing counts
	r.NoteOff(key("C4"))
	expect(t, s)
	r.Read()
	expect(t, s, "says pairs", "lights C♯ D♯")
	r.Read()
	expect(t, s, "says press a pair")

	// The pair, one key after the other: on the way, then the answer.
	r.NoteOn(key("D♯4"))
	expect(t, s)
	r.NoteOn(key("C♯4"))
	expect(t, s, "lights out", "right")

	// The next step waits for the game, the keys up and a beat gone:
	// until then a key counts for nothing.
	if !r.Waiting() {
		t.Fatal("not waiting after the pair")
	}
	r.NoteOff(key("D♯4"))
	r.NoteOff(key("C♯4"))
	r.NoteOn(key("G4"))
	r.NoteOff(key("G4"))
	expect(t, s)
	r.Resume()
	expect(t, s, "says find C")

	// A D for a C: the casserole, the Cs shown, and the lesson waits.
	r.NoteOn(key("D3"))
	r.NoteOff(key("D3"))
	expect(t, s, "casserole, shows C")
	r.NoteOn(key("C5")) // any C
	expect(t, s, "lights out", "right")
	r.NoteOff(key("C5"))
	r.Resume()
	expect(t, s, "says after me", "lights out", "plays C4 D4 C4")

	// The motif: nothing counts while he plays, then note after note,
	// the C held under the left hand.
	r.NoteOn(key("C4"))
	r.NoteOff(key("C4"))
	expect(t, s)
	r.PhraseEnded()
	r.NoteOn(key("C3"))
	expect(t, s, "lights out")

	// A wrong note: the casserole, and he plays the motif again, to be
	// played back from its start.
	r.NoteOn(key("E4"))
	expect(t, s, "casserole, shows ", "plays C4 D4 C4")
	r.NoteOff(key("E4"))
	r.NoteOff(key("C3"))
	r.PhraseEnded()
	for _, n := range []string{"C3", "D4"} {
		r.NoteOn(key(n))
		r.NoteOff(key(n))
	}
	expect(t, s, "lights out", "lights out")
	if r.Done() {
		t.Fatal("done before the motif's last C")
	}
	r.NoteOn(key("C3"))
	expect(t, s, "lights out", "right")
	if !r.Done() {
		t.Error("not done after the motif")
	}
}

// A scale holds "do sol do", but it is not "do sol do": the motif is
// played back note for note, or not at all.
func TestRepeatExactly(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Repeat{Phrase: "after me", Keys: keys("C4", "G4", "C5"), Want: []Target{Do, Sol, Do}}})
	r.Start()
	r.PhraseEnded()
	for _, n := range []string{"B3", "C4", "D4", "E4", "F4", "G4", "A4", "B4", "C5"} {
		r.NoteOn(key(n))
		r.NoteOff(key(n))
		r.PhraseEnded() // each wrong note has him play it again
	}
	if r.Done() {
		t.Error("a scale passed for C G C")
	}
}

// Said rather than shown, "C, E, C" is played from the words. A wrong
// note shows the keys of the whole motif, lit until it is played
// through, from its start again.
func TestRepeatSaid(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Repeat{Phrase: "C, E, C", Want: []Target{Do, Mi, Do}}})
	r.Start()
	expect(t, s, "says C, E, C", "lights out")
	for _, n := range []string{"C4", "F4"} {
		r.NoteOn(key(n))
		r.NoteOff(key(n))
	}
	expect(t, s, "casserole, shows C E")
	for _, n := range []string{"C4", "E4", "C5"} {
		r.NoteOn(key(n))
		r.NoteOff(key(n))
	}
	expect(t, s, "lights out", "right")
	if !r.Done() {
		t.Error("not over after C E C")
	}
}

// "Joue C D E F G A B": the B hangs, held by the game, until the C
// above it, blinking, resolves it. A wrong note on the way starts the
// phrase over; once it hangs, a wrong note is only a casserole.
func TestHang(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Hang{
		Phrase: "play C D E F G A B", Want: []Target{Do, Re, Mi, Fa, Sol, La, Si},
		Humpf: "humpf", Resolve: Do, Thanks: "thanks",
	}})
	r.Start()
	s.take()
	for _, n := range []string{"C4", "D4", "E4", "F4", "G4", "A4", "B4"} {
		r.NoteOn(key(n))
		r.NoteOff(key(n))
	}
	expect(t, s, "lights out", "holds B4", "says humpf", "blinks C5")
	r.NoteOn(key("D5"))
	r.NoteOff(key("D5"))
	expect(t, s, "casserole, shows ")
	r.NoteOn(key("C5"))
	expect(t, s, "blinks no more", "lets go", "says thanks")
	r.NoteOff(key("C5"))
	if !r.Done() {
		t.Error("not over once resolved")
	}
}

// "Et mi♯ ?" Two misses, a "no" each, nothing shown; at the third,
// the answer given and F shown, to be played.
func TestGuess(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Guess{Phrase: "and E♯?", Want: MiSharp, Nope: "no", Answer: "on F", Chances: 3}})
	r.Start()
	s.take()
	for _, n := range []string{"E4", "F♯4"} {
		r.NoteOn(key(n))
		r.NoteOff(key(n))
		expect(t, s, "casserole, shows ", "says no")
	}
	r.NoteOn(key("G4"))
	r.NoteOff(key("G4"))
	expect(t, s, "casserole, shows F", "says on F")
	r.NoteOn(key("F4"))
	expect(t, s, "lights out", "right")
	r.NoteOff(key("F4"))
	if !r.Done() {
		t.Error("not over once F is played")
	}
}

// The first line of the blues, played on the bass: a root a bar, the
// line waiting for the player; a wrong root, the casserole, and the
// line stays on its bar.
func TestPlayLine(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&PlayLine{
		Phrase: "your turn", Chords: []string{"C7", "F7", "C7", "C7"},
		Want: []Target{Do, Fa, Do, Do},
	}})
	r.Start()
	expect(t, s, "writes ", "bass true", "says your turn", "lights out", "line C7 | F7 | C7 | C7, bar 0")
	r.NoteOn(key("C2"))
	r.NoteOff(key("C2"))
	expect(t, s, "lights out", "line C7 | F7 | C7 | C7, bar 1")
	r.NoteOn(key("G2"))
	r.NoteOff(key("G2"))
	expect(t, s, "casserole, shows F")
	for _, n := range []string{"F2", "C2", "C2"} {
		r.NoteOn(key(n))
		r.NoteOff(key(n))
	}
	expect(t, s,
		"lights out", "line C7 | F7 | C7 | C7, bar 2",
		"lights out", "line C7 | F7 | C7 | C7, bar 3",
		"lights out", "line C7 | F7 | C7 | C7, bar -1", "right")
	if !r.Done() {
		t.Error("not over after the fourth bar")
	}
}

// A long phrase missed: the walker adds his reminder to take it slowly,
// once, however many misses follow. A short one goes without.
func TestRemindSlowly(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Repeat{Phrase: "C, D, E, F, G", Want: []Target{Do, Re, Mi, Fa, Sol}}})
	r.Start()
	s.take()
	for range 2 {
		r.NoteOn(key("C4"))
		r.NoteOff(key("C4"))
		r.NoteOn(key("E4"))
		r.NoteOff(key("E4"))
	}
	expect(t, s, "adds "+slowly, "casserole, shows C D E F G", "casserole, shows C D E F G")

	r = NewRunner(s, []Step{&Repeat{Phrase: "C, E, C", Want: []Target{Do, Mi, Do}}})
	r.Start()
	s.take()
	r.NoteOn(key("D4"))
	r.NoteOff(key("D4"))
	expect(t, s, "casserole, shows C E")
}

func TestGroup(t *testing.T) {
	for _, tc := range []struct {
		group Group
		held  []string
		want  Verdict
	}{
		{Pair, []string{"C♯4", "D♯4"}, Hit},
		{Pair, []string{"D♯2", "C♯2"}, Hit},
		{Pair, []string{"C♯4"}, Partial},
		{Pair, []string{"C♯4", "D♯5"}, Miss}, // two octaves apart
		{Pair, []string{"D♯4", "F♯4"}, Miss}, // across two groups
		{Pair, []string{"C4"}, Miss},
		{Trio, []string{"F♯3", "G♯3", "A♯3"}, Hit},
		{Trio, []string{"G♯3", "A♯3"}, Partial},
		{Trio, []string{"F♯3", "A♯3"}, Partial}, // the middle one to come
		{Trio, []string{"F♯3", "G♯3", "A♯3", "C♯4"}, Miss},
	} {
		if got := tc.group.Judge(keys(tc.held...)); got != tc.want {
			t.Errorf("%s on %v: got %d, want %d", pcNames(tc.group), tc.held, got, tc.want)
		}
	}
}

// Ten Cs for a pair: nobody misses that in good faith. The walker gets
// annoyed, then walks out; back, he counts from nothing again.
func TestTease(t *testing.T) {
	s := &stage{}
	r := NewRunner(s, []Step{&Ask{Phrase: "press a pair", Want: Pair, Teasing: true}})
	r.Start()
	s.take()
	var said []string
	for range 10 {
		r.NoteOn(key("C4"))
		r.NoteOff(key("C4"))
		for _, e := range s.take() {
			if e != "casserole, shows C♯ D♯" {
				said = append(said, e)
			}
		}
	}
	if want := []string{"Hé…", "on purpose?", "walks out"}; !slices.Equal(said, want) {
		t.Errorf("got %q, want %q", said, want)
	}
	r.Again() // he is back
	for range missesHey - 1 {
		r.NoteOn(key("C4"))
		r.NoteOff(key("C4"))
	}
	for _, e := range s.take() {
		if e == "Hé…" {
			t.Error("annoyed again too soon")
		}
	}
}
