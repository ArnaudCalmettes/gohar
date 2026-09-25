package dex_test

import (
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A notion is designated the way the library designates it, so the same
// shape met in two games is the same entry. This is the whole reason
// there is no free string here.
func TestNotionIdentity(t *testing.T) {
	t.Run("the same mode designated twice is the same value", func(t *testing.T) {
		a := dex.ModeOf(harmony.MelodicMinor, 4)
		b := dex.ModeOf(harmony.MelodicMinor, 4)
		assert.Equal(t, a, b)
	})

	t.Run("two degrees of a system are two notions", func(t *testing.T) {
		assert.NotEqual(t,
			dex.ModeOf(harmony.MelodicMinor, 4),
			dex.ModeOf(harmony.MelodicMinor, 5),
		)
	})

	t.Run("a mode and a tetrad never collide", func(t *testing.T) {
		assert.NotEqual(t,
			dex.ModeOf(harmony.NaturalMajor, 1),
			dex.TetradOf(harmony.ChordMajorSeventh),
		)
	})

	t.Run("a notion keys a map", func(t *testing.T) {
		seen := map[dex.Notion]int{}
		seen[dex.ModeOf(harmony.HarmonicMinor, 5)]++
		seen[dex.ModeOf(harmony.HarmonicMinor, 5)]++
		assert.Equal(t, 2, seen[dex.ModeOf(harmony.HarmonicMinor, 5)])
	})

	t.Run("the zero value designates nothing", func(t *testing.T) {
		assert.True(t, dex.Notion{}.IsZero())
		assert.False(t, dex.ModeOf(harmony.NaturalMajor, 1).IsZero())
	})
}

// The written form is what persistence keeps, so it has to survive a
// round trip and it must not depend on the numeric value of a constant.
func TestNotionTextIsStable(t *testing.T) {
	notions := []dex.Notion{
		dex.ModeOf(harmony.MelodicMinor, 4),
		dex.TetradOf(harmony.ChordDominantSeventh),
		dex.TetrachordOf(harmony.TetrachordHarmonic),
		dex.ProgressionOf(1),
	}

	for _, n := range notions {
		t.Run(n.String()+" reads back as itself", func(t *testing.T) {
			back, err := dex.ParseNotion(n.String())
			require.NoError(t, err)
			assert.Equal(t, n, back)
		})
	}

	t.Run("two notions never share a written form", func(t *testing.T) {
		seen := map[string]dex.Notion{}
		for _, n := range notions {
			previous, clash := seen[n.String()]
			assert.False(t, clash, "%v and %v write the same", previous, n)
			seen[n.String()] = n
		}
	})

	t.Run("nonsense is refused rather than guessed", func(t *testing.T) {
		_, err := dex.ParseNotion("mode:banane")
		assert.Error(t, err)
	})
}

// The dex body is not written yet: New, Apply, Look, Tonics and Visible
// are still TODO, and so are the decisions they rest on (what a wrong
// answer refreshes, whether a choice also counts as a production, when
// a slot really left the colour open).
//
// The tests below are the specification of that body. They are kept
// compiling and skipped rather than deleted, because rewriting them
// from memory later would quietly lose the rules they encode.
func skipUntilBody(t *testing.T) {
	t.Helper()
	t.Skip("dex body not implemented yet")
}

// Nothing ever removes a mark. What passes is time, and that is a
// different thing from a mistake.
func TestMarksAreNeverUndone(t *testing.T) {
	skipUntilBody(t)

	d := dex.New()
	lydianDominant := dex.ModeOf(harmony.MelodicMinor, 4)
	at := time.Date(2026, 9, 22, 21, 0, 0, 0, time.UTC)

	d.Apply(dex.Report{
		Game: "eartrainer", At: at,
		Facts: []dex.Fact{{Kind: dex.FactNamed, Notion: lydianDominant, Correct: true}},
	})

	before, ok := d.Look(lydianDominant)
	require.True(t, ok)
	require.True(t, before.Recognized.Held())

	t.Run("a later report cannot take the mark back", func(t *testing.T) {
		d.Apply(dex.Report{
			Game: "eartrainer", At: at.Add(time.Hour),
			Facts: []dex.Fact{{Kind: dex.FactNamed, Notion: lydianDominant, Correct: false}},
		})

		after, ok := d.Look(lydianDominant)
		require.True(t, ok)
		assert.True(t, after.Recognized.Held(),
			"a wrong answer is not a fact that undoes anything")
	})

	t.Run("a correct answer refreshes instead", func(t *testing.T) {
		later := at.Add(48 * time.Hour)
		d.Apply(dex.Report{
			Game: "eartrainer", At: later,
			Facts: []dex.Fact{{Kind: dex.FactNamed, Notion: lydianDominant, Correct: true}},
		})

		after, ok := d.Look(lydianDominant)
		require.True(t, ok)
		assert.Equal(t, later, after.Recognized.Last)
		assert.Equal(t, before.Recognized.First, after.Recognized.First,
			"the first time never moves")
		assert.Greater(t, after.Recognized.Count, before.Recognized.Count)
	})
}

// The hand has twelve topographies, so production is kept per tonic.
// The summary is derived from the detail, never the other way round.
func TestProductionIsKeptPerTonic(t *testing.T) {
	skipUntilBody(t)

	d := dex.New()
	twoFiveOne := dex.ProgressionOf(1)
	at := time.Date(2026, 9, 22, 21, 0, 0, 0, time.UTC)

	produce := func(tonic harmony.PitchClass, when time.Time) {
		d.Apply(dex.Report{
			Game: "shmup", At: when,
			Facts: []dex.Fact{{
				Kind: dex.FactProduced, Notion: twoFiveOne, Tonic: tonic,
			}},
		})
	}

	produce(0, at)
	produce(10, at.Add(time.Hour))

	e, ok := d.Look(twoFiveOne)
	require.True(t, ok)

	t.Run("each tonic carries its own date and count", func(t *testing.T) {
		assert.Equal(t, at, e.Produced[0].Last)
		assert.Equal(t, at.Add(time.Hour), e.Produced[10].Last)
		assert.Len(t, e.Produced, 2)
	})

	t.Run("the mask is derived for display", func(t *testing.T) {
		want, err := harmony.NewPitchSet(0, 10)
		require.NoError(t, err)
		assert.Equal(t, want, e.Tonics())
	})

	t.Run("a tonic never played is absent, not zero", func(t *testing.T) {
		_, played := e.Produced[6]
		assert.False(t, played)
	})

	t.Run("recognition is not detailed by tonic", func(t *testing.T) {
		d.Apply(dex.Report{
			Game: "eartrainer", At: at,
			Facts: []dex.Fact{{
				Kind: dex.FactNamed, Notion: twoFiveOne, Tonic: 3, Correct: true,
			}},
		})

		e, ok := d.Look(twoFiveOne)
		require.True(t, ok)
		assert.True(t, e.Recognized.Held(),
			"a relative ear recognises a shape whatever the key")
	})
}

// The four marks are independent. Laying a shape under the hands
// without naming it by ear is a real path, and the most common one for
// someone who explores at the keyboard.
func TestMarksAreIndependent(t *testing.T) {
	skipUntilBody(t)

	d := dex.New()
	phrygianDominant := dex.ModeOf(harmony.HarmonicMinor, 5)
	at := time.Date(2026, 9, 22, 21, 0, 0, 0, time.UTC)

	d.Apply(dex.Report{
		Game: "shmup", At: at,
		Facts: []dex.Fact{{
			Kind: dex.FactProduced, Notion: phrygianDominant, Tonic: 4,
		}},
	})

	e, ok := d.Look(phrygianDominant)
	require.True(t, ok)
	assert.True(t, e.Produced[4].Held())
	assert.False(t, e.Recognized.Held(), "producing says nothing about the ear")
	assert.False(t, e.Used.Held())
}

// The freshness fact exists because notions are not independent.
func TestSoundedRefreshesWithoutMarking(t *testing.T) {
	skipUntilBody(t)

	d := dex.New()
	minorSeventh := dex.TetradOf(harmony.ChordMinorSeventh)
	at := time.Date(2026, 9, 22, 21, 0, 0, 0, time.UTC)

	t.Run("a notion only heard in passing opens no mark", func(t *testing.T) {
		d.Apply(dex.Report{
			Game: "shmup", At: at,
			Facts: []dex.Fact{{Kind: dex.FactSounded, Notion: minorSeventh}},
		})

		e, ok := d.Look(minorSeventh)
		if ok {
			assert.False(t, e.Recognized.Held())
			assert.False(t, e.Used.Held())
			assert.Empty(t, e.Produced)
		}
	})

	t.Run("it does refresh what is already known", func(t *testing.T) {
		d.Apply(dex.Report{
			Game: "shmup", At: at,
			Facts: []dex.Fact{{
				Kind: dex.FactProduced, Notion: minorSeventh, Tonic: 2,
			}},
		})

		later := at.Add(72 * time.Hour)
		d.Apply(dex.Report{
			Game: "shmup", At: later,
			Facts: []dex.Fact{{Kind: dex.FactSounded, Notion: minorSeventh}},
		})

		e, ok := d.Look(minorSeventh)
		require.True(t, ok)
		assert.Equal(t, later, e.Produced[2].Last,
			"working a two five one sounds a minor seventh, and the dex has to know")
	})
}

// A discovery is the moment worth showing.
func TestDiscoveryIsReported(t *testing.T) {
	skipUntilBody(t)

	d := dex.New()
	lydianDominant := dex.ModeOf(harmony.MelodicMinor, 4)
	at := time.Date(2026, 9, 22, 21, 0, 0, 0, time.UTC)

	first := d.Apply(dex.Report{
		Game: "shmup", At: at,
		Facts: []dex.Fact{{Kind: dex.FactHeard, Notion: lydianDominant}},
	})

	require.Len(t, first, 1)
	assert.True(t, first[0].Discovery)
	assert.Equal(t, lydianDominant, first[0].Notion)

	second := d.Apply(dex.Report{
		Game: "shmup", At: at.Add(time.Hour),
		Facts: []dex.Fact{{Kind: dex.FactHeard, Notion: lydianDominant}},
	})

	for _, c := range second {
		assert.False(t, c.Discovery, "a wild lydian dominant appears only once")
	}
}

// Where a mark was first earned is what makes a collection tellable.
func TestMarksRememberWhereTheyWereEarned(t *testing.T) {
	skipUntilBody(t)

	d := dex.New()
	n := dex.ModeOf(harmony.MelodicMinor, 4)
	at := time.Date(2026, 9, 22, 21, 0, 0, 0, time.UTC)

	d.Apply(dex.Report{
		Game: "eartrainer", At: at,
		Facts: []dex.Fact{{Kind: dex.FactNamed, Notion: n, Correct: true}},
	})
	d.Apply(dex.Report{
		Game: "shmup", At: at.Add(time.Hour),
		Facts: []dex.Fact{{Kind: dex.FactProduced, Notion: n, Tonic: 10}},
	})

	e, ok := d.Look(n)
	require.True(t, ok)
	assert.Equal(t, dex.GameID("eartrainer"), e.Recognized.FirstIn)
	assert.Equal(t, dex.GameID("shmup"), e.Produced[10].FirstIn)
}

// A notion never crossed does not show.
func TestNeverMetIsNotVisible(t *testing.T) {
	skipUntilBody(t)

	d := dex.New()
	assert.Empty(t, d.Visible(),
		"a beginner facing thirty five empty slots reads everything they do not know")
}
