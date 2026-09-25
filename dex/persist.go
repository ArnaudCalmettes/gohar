package dex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// formatVersion is written with every dex. Reading refuses any other
// value rather than guessing: a collection is the player's, and a
// misread one moves marks silently.
const formatVersion = 1

// The written form of a dex. Kept apart from the types the package
// works with, so that renaming a field in Go never renames it on disk.
//
// Notions are keyed by [Notion.String], tonics by their number, and a
// mark never earned is left out rather than written as zeros.
type (
	dexJSON struct {
		Version int                  `json:"version"`
		Entries map[string]entryJSON `json:"entries"`
	}

	entryJSON struct {
		Met        *markJSON           `json:"met,omitempty"`
		Recognized *markJSON           `json:"recognized,omitempty"`
		Produced   map[string]markJSON `json:"produced,omitempty"`
		Used       *markJSON           `json:"used,omitempty"`
		Overheard  *markJSON           `json:"overheard,omitempty"`
	}

	markJSON struct {
		First   time.Time `json:"first"`
		Last    time.Time `json:"last"`
		Count   int       `json:"count"`
		FirstIn GameID    `json:"firstIn,omitempty"`
	}
)

// MarshalJSON writes the collection.
//
// Where it goes is the caller's business: a file on a desktop, local
// storage in a browser. The dex only knows the form.
func (d *Dex) MarshalJSON() ([]byte, error) {
	out := dexJSON{Version: formatVersion, Entries: map[string]entryJSON{}}
	for n, e := range d.entries {
		ej := entryJSON{
			Met:        markOut(e.Met),
			Recognized: markOut(e.Recognized),
			Used:       markOut(e.Used),
			Overheard:  markOut(e.Overheard),
		}
		if len(e.Produced) > 0 {
			ej.Produced = make(map[string]markJSON, len(e.Produced))
			for tonic, m := range e.Produced {
				ej.Produced[strconv.Itoa(int(tonic))] = *markOut(m)
			}
		}
		out.Entries[n.String()] = ej
	}
	return json.Marshal(out)
}

// UnmarshalJSON reads back what MarshalJSON wrote, into an empty dex or
// over an existing one, whose entries it replaces.
//
// Strict, for the same reason as [ParseNotion]: an unknown notion, a
// tonic outside the twelve or a mark with no count is an error, never
// a guess.
func (d *Dex) UnmarshalJSON(data []byte) error {
	var in dexJSON
	if err := json.Unmarshal(data, &in); err != nil {
		return fmt.Errorf("dex: %w", err)
	}
	if in.Version != formatVersion {
		return fmt.Errorf("dex: format version %d, expected %d", in.Version, formatVersion)
	}

	entries := make(map[Notion]*Entry, len(in.Entries))
	for key, ej := range in.Entries {
		n, err := ParseNotion(key)
		if err != nil {
			return err
		}
		if n.IsZero() {
			return fmt.Errorf("dex: an entry designates no notion")
		}

		e := &Entry{Notion: n}
		for _, m := range []struct {
			in  *markJSON
			out *Mark
		}{
			{ej.Met, &e.Met},
			{ej.Recognized, &e.Recognized},
			{ej.Used, &e.Used},
			{ej.Overheard, &e.Overheard},
		} {
			if m.in == nil {
				continue
			}
			if *m.out, err = markIn(*m.in, key); err != nil {
				return err
			}
		}

		if len(ej.Produced) > 0 {
			e.Produced = make(map[harmony.PitchClass]Mark, len(ej.Produced))
			for t, mj := range ej.Produced {
				v, err := strconv.ParseUint(t, 10, 8)
				tonic := harmony.PitchClass(v)
				if err != nil || !tonic.IsValid() {
					return fmt.Errorf("dex: %s: %q is not a tonic", key, t)
				}
				if e.Produced[tonic], err = markIn(mj, key); err != nil {
					return err
				}
			}
		}

		// An entry with nothing in it was never written by MarshalJSON,
		// which only knows entries some fact opened.
		if !e.held() && !e.Overheard.Held() {
			return fmt.Errorf("dex: %s holds no mark", key)
		}
		entries[n] = e
	}

	d.entries = entries
	return nil
}

func markOut(m Mark) *markJSON {
	if !m.Held() {
		return nil
	}
	return &markJSON{First: m.First, Last: m.Last, Count: m.Count, FirstIn: m.FirstIn}
}

func markIn(m markJSON, key string) (Mark, error) {
	if m.Count <= 0 {
		return Mark{}, fmt.Errorf("dex: %s: a written mark has no count", key)
	}
	if m.Last.Before(m.First) {
		return Mark{}, fmt.Errorf("dex: %s: a mark ends before it starts", key)
	}
	return Mark{First: m.First, Last: m.Last, Count: m.Count, FirstIn: m.FirstIn}, nil
}
