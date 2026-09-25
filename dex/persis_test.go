package dex_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ArnaudCalmettes/gohar/dex"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A collection written down and read back is the same collection, marks,
// tonics, games and silhouettes included.
func TestDexRoundTrip(t *testing.T) {
	at := time.Date(2026, 9, 25, 21, 0, 0, 0, time.UTC)
	lydianDominant := dex.ModeOf(harmony.MelodicMinor, 4)
	phrygianNatural6 := dex.ModeOf(harmony.MelodicMinor, 2)

	d := dex.New()
	d.Apply(dex.Report{Game: "ear", At: at, Facts: []dex.Fact{
		{Kind: dex.FactHeard, Notion: lydianDominant},
		{Kind: dex.FactNamed, Notion: lydianDominant},
		{Kind: dex.FactSounded, Notion: phrygianNatural6},
	}})
	d.Apply(dex.Report{Game: "shmup", At: at.Add(time.Hour), Facts: []dex.Fact{
		{Kind: dex.FactChosen, Notion: lydianDominant, Tonic: 7},
		{Kind: dex.FactProduced, Notion: lydianDominant, Tonic: 0},
	}})

	data, err := json.Marshal(d)
	require.NoError(t, err)

	back := dex.New()
	require.NoError(t, json.Unmarshal(data, back))

	for _, n := range []dex.Notion{lydianDominant, phrygianNatural6} {
		want, _ := d.Look(n)
		got, ok := back.Look(n)
		require.True(t, ok, n.String())
		assert.Equal(t, want, got, n.String())
	}
	assert.Equal(t, d.Visible(), back.Visible())
}

// Reading is strict: what does not parse exactly is refused.
func TestDexReadIsStrict(t *testing.T) {
	cases := map[string]string{
		"another version":      `{"version":2,"entries":{}}`,
		"an unknown notion":    `{"version":1,"entries":{"mode:banane":{"met":{"count":1}}}}`,
		"a tonic past B":       `{"version":1,"entries":{"tetrad:000091":{"produced":{"12":{"count":1}}}}}`,
		"a mark with no count": `{"version":1,"entries":{"tetrad:000091":{"met":{"count":0}}}}`,
		"an empty entry":       `{"version":1,"entries":{"tetrad:000091":{}}}`,
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, json.Unmarshal([]byte(data), dex.New()))
		})
	}
}
