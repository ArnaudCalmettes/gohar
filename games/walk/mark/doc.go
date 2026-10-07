// Package mark takes note of what the player plays in Walk with me:
// the beats a grid expects (Expect), the mark of each note and each
// arrival against them (Marker), the chart that waits without tempo
// (Practice), and what a run leaves (Summary, Advise). It never grades
// the player: a mark is what happened, said as it is (see "Les briques"
// and "Le bilan" in docs/walk.md).
//
// It knows neither the screen nor the sound: its tests play the grids
// of the game, note by note, on a clock of their own.
package mark
