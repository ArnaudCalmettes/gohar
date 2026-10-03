package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A recorder writes down what the marker hears, a note per line, to
// read a line away from the game: the player's, or the demo's, to
// compare them. Practice is not recorded: it has no time.
//
//	chorus  bar.beat  offset  chord  note   mark
//	1       9.4       +12ms   Gm7    D2     sur le temps, note de l'accord
//
// The bar is counted within the chorus. The note carries its octave,
// MIDI's: C4 is middle C, the double bass sounds E1 to G2.
type recorder struct {
	f *os.File
	w *bufio.Writer
}

func newRecorder(path string) (*recorder, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &recorder{f: f, w: bufio.NewWriter(f)}, nil
}

// run opens a run with a heading line: when, at which tempo, who plays.
func (r *recorder) run(at time.Time, bpm float64, demo bool) {
	who := "joueur"
	if demo {
		who = "démo"
	}
	fmt.Fprintf(r.w, "\n# %s, %g bpm, %s\n", at.Format("2006-01-02 15:04:05"), bpm, who)
}

// note writes one line.
func (r *recorder) note(chorus int, p Position, off time.Duration, chord, note, mark string) {
	fmt.Fprintf(r.w, "%d\t%d.%d\t%+dms\t%s\t%s\t%s\n", chorus, p.Bar, p.Beat, off.Milliseconds(), chord, note, mark)
}

// flush writes out what is buffered: at the end of each run, so that a
// crash loses one run at most.
func (r *recorder) flush() error {
	return r.w.Flush()
}

func (r *recorder) Close() error {
	err := r.w.Flush()
	if cerr := r.f.Close(); err == nil {
		err = cerr
	}
	return err
}

// octaveName appends the octave of `key` to `name`, as MIDI counts
// them: 28 is E1, 43 is G2.
func octaveName(name string, key int) string {
	pc := harmony.PitchClass(key % octave)
	o := key/octave - 1
	// B♯ and C♭ spell across the octave line: C♭2 sounds as B1.
	switch {
	case pc == 11 && name[0] == 'C':
		o++
	case pc == 0 && name[0] == 'B':
		o--
	}
	return fmt.Sprintf("%s%d", name, o)
}
