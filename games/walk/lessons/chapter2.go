package lessons

// pulseBPM is the tempo of chapter 2: 90 at least. Slower, the lesson
// would fall into the opposite trap: playing slowly is harder than
// playing fast, each beat left alone too long (see
// docs/debutants/chapitre-2.md).
const pulseBPM = 90

// theMeter is lesson 2.1, "La mesure" (see
// docs/debutants/chapitre-2.md): the band comes in, the hi-hat on every
// beat, the walker counts the bar, and the player plays every beat,
// then the 1 only. A miss never stops the pulse.
func theMeter() []Step {
	return []Step{
		Band{Phrase: "l2.1.listen", BPM: pulseBPM},
		Band{Phrase: "l2.1.count", BPM: pulseBPM, Count: true},
		Say{Phrase: "l2.1.one"},
		&Groove{Phrase: "l2.1.every", Again: "l2.1.again", On: []int{1, 2, 3, 4}, Bars: 2},
		&Groove{Phrase: "l2.1.first", Again: "l2.1.again", On: []int{1}, Bars: 2},
		Say{Phrase: "l2.1.meter"},
		Band{Phrase: "l2.1.bravo"}, // the band stops
	}
}

// bluesInC is the blues of the game in C, as the palier 0 plays it, a
// bar a string, two chords in the turnaround's bars; bluesRoots, the
// root on the 1 of each bar. Its first line is the activity of lesson
// 2.1; the whole of it, the closing activity's.
var (
	bluesInC = []string{
		"C7", "F7", "C7", RepeatBar,
		"F7", RepeatBar, "C6", "A7",
		"Dm7", "G7", "C6 A7", "D7 G7",
	}
	bluesRoots = []Target{Do, Fa, Do, Do, Fa, Fa, Do, La, Re, Sol, Do, Re}
)

// meterActivity is the activity of lesson 2.1: the first line of the
// blues in C, the root of each bar on its 1, slowly, with every help,
// the hi-hat on every beat and the walker counting. The whole blues
// waits for the activity that closes the chapter.
func meterActivity() []Step {
	const line = 4
	return []Step{
		Band{Phrase: "a2.1.hello", BPM: pulseBPM, Count: true},
		&Groove{Phrase: "a2.1.go", Again: "l2.1.again", On: []int{1}, Bars: line, Want: bluesRoots[:line], Chords: bluesInC[:line]},
		Band{Phrase: "a2.1.bravo"},
	}
}
