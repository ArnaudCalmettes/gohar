package lessons

import (
	"math/rand/v2"
	"slices"
)

// The scripts of the lessons written so far, by ID. Each call builds
// fresh steps: Ask and Repeat keep where they are.
var scripts = map[string]func() []Step{
	"1.1":   firstSteps,
	"1.2":   cMajorScale,
	"1.3":   sharps,
	"1.4":   flats,
	"1.end": firstGrids,
}

// Script returns the steps of lesson `id`, or nil while it is not
// written.
func Script(id string) []Step {
	if s, ok := scripts[id]; ok {
		return s()
	}
	return nil
}

// Phrases returns the phrases every script and every activity says,
// for the game to check that it can say them all.
func Phrases() []string {
	var out []string
	add := func(steps []Step) {
		for _, st := range steps {
			if p := st.said(); p != "" {
				out = append(out, p)
			}
			if o, ok := st.(interface{ alsoSays() []string }); ok {
				out = append(out, o.alsoSays()...)
			}
		}
	}
	for _, s := range scripts {
		add(s())
	}
	for _, a := range activities {
		add(a(rand.New(rand.NewPCG(1, 2)))) // any order says the same phrases
	}
	out = append(out, slowly)                                 // said at a miss in a long phrase
	for _, n := range slices.Concat(sevenNotes, fiveSharps) { // drawn by a script, not all said each time
		out = append(out, findLetter[n], findSolfege[n])
	}
	for _, n := range fiveFlats {
		out = append(out, findLetterFlat[n], findSolfegeFlat[n])
	}
	return out
}
