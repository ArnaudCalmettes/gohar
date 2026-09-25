package analysis_test

import (
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// No fake clock anywhere in this file. The caller drives, so a test
// that wants to skip four seconds passes a time four seconds later.
// That is the payoff for not putting a Clock port in the engine.
var t0 = time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

func at(ms int) time.Time { return t0.Add(time.Duration(ms) * time.Millisecond) }

func engine(t *testing.T, c analysis.Config) *analysis.Engine {
	t.Helper()
	e, err := analysis.NewEngine(newRecognizer(t), c)
	require.NoError(t, err)
	return e
}

func play(e *analysis.Engine, ms int, pitches ...harmony.Pitch) {
	for _, p := range pitches {
		e.NoteOn(p, at(ms))
	}
	e.Advance(at(ms))
}

// The gather window is what turns a rolled voicing into a chord. Notes
// played one after another, and released, still form one reading.
func TestArpeggioFormsAChord(t *testing.T) {
	e := engine(t, analysis.DefaultConfig)

	t.Run("notes rolled inside the window read as one chord", func(t *testing.T) {
		e.NoteOn(60, at(0))
		e.NoteOff(60, at(80))
		e.NoteOn(64, at(100))
		e.NoteOff(64, at(180))
		e.NoteOn(67, at(200))
		e.Advance(at(250))

		got, ok := e.Reading()
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(0), got.Root)
		assert.Equal(t, harmony.ChordMajorTriad, got.Tetrad)
	})

	t.Run("a note played after the window has closed does not join it", func(t *testing.T) {
		e.Reset()

		e.NoteOn(60, at(0))
		e.NoteOff(60, at(50))
		e.Advance(at(1000))
		e.NoteOn(64, at(1000))
		e.NoteOn(67, at(1000))
		e.NoteOn(71, at(1000))
		e.Advance(at(1050))

		got, ok := e.Reading()
		require.True(t, ok)

		sounding := harmony.Chord{Root: got.Root, Pattern: got.Pattern}.Set()
		assert.False(t, sounding.Contains(0),
			"the C fell out of the window and must not be in the reading",
		)
		assert.Equal(t, harmony.PitchClass(4), got.Root,
			"three notes left, and they are an E minor triad")
	})
}

// The hysteresis exists so the display does not flicker. A passing note
// that briefly tips the score must not change what the player sees.
func TestHysteresis(t *testing.T) {
	t.Run("a reading survives a note too brief to settle", func(t *testing.T) {
		e := engine(t, analysis.DefaultConfig)
		play(e, 0, 60, 64, 67)

		before, ok := e.Reading()
		require.True(t, ok)

		e.NoteOn(69, at(100))
		e.NoteOff(69, at(140))
		e.Advance(at(180))

		after, ok := e.Reading()
		require.True(t, ok)
		assert.Equal(t, before.Root, after.Root,
			"forty milliseconds is a passing note, not a chord change",
		)
	})

	t.Run("a reading changes once the new one has outlasted Settle", func(t *testing.T) {
		e := engine(t, analysis.DefaultConfig)
		play(e, 0, 60, 64, 67)

		e.NoteOff(60, at(100))
		e.NoteOff(64, at(100))
		e.NoteOff(67, at(100))
		play(e, 600, 62, 65, 69)
		e.Advance(at(900))

		got, ok := e.Reading()
		require.True(t, ok)
		assert.Equal(t, harmony.PitchClass(2), got.Root)
	})

	t.Run("a still valid reading is not traded for an equal one", func(t *testing.T) {
		e := engine(t, analysis.DefaultConfig)
		play(e, 0, 60, 63, 66, 69)

		before, ok := e.Reading()
		require.True(t, ok)

		for ms := 200; ms <= 3000; ms += 200 {
			e.Advance(at(ms))
			after, ok := e.Reading()
			require.True(t, ok)
			assert.Equal(t, before.Root, after.Root,
				"four roots explain a diminished seventh chord, and the display "+
					"must not walk between them while the chord is held",
			)
		}
	})
}

// The two time constants must not be one. A key that moves as fast as a
// chord is not a key.
func TestTonalityIsSlowerThanChords(t *testing.T) {
	e := engine(t, analysis.DefaultConfig)

	t.Run("nothing is claimed before enough has been heard", func(t *testing.T) {
		play(e, 0, 60, 64, 67)
		assert.True(t, e.Tonality().IsZero(),
			"one chord is not a key, and the zero value says so",
		)
	})

	t.Run("chords keep moving while the key stays put", func(t *testing.T) {
		e.Reset()

		progression := [][]harmony.Pitch{
			{60, 64, 67}, {65, 69, 72}, {55, 59, 62, 65}, {60, 64, 67},
			{62, 65, 69}, {55, 59, 62, 65}, {60, 64, 67},
		}
		for i, chord := range progression {
			play(e, i*800, chord...)
			e.Advance(at(i*800 + 400))
		}

		var keys []harmony.Tonality
		for i := range progression {
			e.Advance(at(i*800 + 700))
			keys = append(keys, e.Tonality())
		}

		last := keys[len(keys)-1]
		require.False(t, last.IsZero(), "seven chords is enough to hear a key")
		assert.Equal(t, harmony.PitchClass(0), last.Tonic())
	})

	t.Run("a single foreign note does not move the key", func(t *testing.T) {
		before := e.Tonality()
		require.False(t, before.IsZero())

		e.NoteOn(61, at(10_000))
		e.NoteOff(61, at(10_120))
		e.Advance(at(10_400))

		assert.Equal(t, before, e.Tonality(),
			"an accidental is not a modulation",
		)
	})
}

func TestEngineToleratesDeviceMisbehaviour(t *testing.T) {
	e := engine(t, analysis.DefaultConfig)

	t.Run("a note off for a pitch never on is ignored", func(t *testing.T) {
		assert.NotPanics(t, func() {
			e.NoteOff(60, at(0))
			e.Advance(at(10))
		})
	})

	t.Run("an out of order timestamp does not lose the note", func(t *testing.T) {
		e.Reset()

		e.NoteOn(60, at(200))
		e.NoteOn(64, at(180))
		e.NoteOn(67, at(190))
		e.Advance(at(300))

		got, ok := e.Reading()
		require.True(t, ok)
		assert.Equal(t, harmony.ChordMajorTriad, got.Tetrad)
	})
}

func TestAdvanceIsIdempotent(t *testing.T) {
	e := engine(t, analysis.DefaultConfig)
	play(e, 0, 60, 64, 67)

	e.Advance(at(200))
	first, ok := e.Reading()
	require.True(t, ok)

	e.Advance(at(200))
	second, ok := e.Reading()
	require.True(t, ok)

	assert.Equal(t, first, second,
		"a caller may advance twice on one frame without consequence")
}

func TestResetClearsEverything(t *testing.T) {
	e := engine(t, analysis.DefaultConfig)
	play(e, 0, 60, 64, 67)

	require.NotPanics(t, e.Reset)

	_, ok := e.Reading()
	assert.False(t, ok)
	assert.True(t, e.Tonality().IsZero())
	assert.True(t, e.Snapshot().IsEmpty())
}
