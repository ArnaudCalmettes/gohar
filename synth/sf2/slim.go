package sf2

import (
	"fmt"
	"slices"
)

// A Keep is a preset to keep: its bank and patch, and for a drum kit
// the keys a game strikes. Nil keys keep them all.
type Keep struct {
	Bank, Patch int
	Keys        []int
}

// Slim returns a soundfont holding only the presets `keep` names, the
// instruments they play and the samples those play. The INFO list is
// copied as it is.
//
// For a preset with keys, the zones that cover none of them go, in the
// preset and in its instruments, unless another kept preset plays the
// same instrument in full. A kit of fifty percussions, of which a game
// strikes three, keeps three.
//
// A global zone, the first one of a preset or an instrument when it
// names no instrument or no sample, is always kept: it holds the
// defaults of the other zones. A stereo sample keeps the other side it
// links to.
func Slim(f *File, keep []Keep) (*File, error) {
	presets := make([]int, 0, len(keep)) // the kept presets, by index
	keysOf := map[int][]int{}            // their keys, nil for all
	for _, k := range keep {
		i := slices.IndexFunc(f.Presets[:len(f.Presets)-1], func(p PresetHeader) bool {
			return int(p.Bank) == k.Bank && int(p.Patch) == k.Patch
		})
		if i < 0 {
			return nil, fmt.Errorf("sf2: no preset %d:%d", k.Bank, k.Patch)
		}
		presets = append(presets, i)
		keysOf[i] = k.Keys
	}
	slices.Sort(presets)

	// The preset zones kept, and the keys each instrument is played on.
	presetZones := map[int][]int{}
	instKeys := map[int]*keySet{}
	for _, p := range presets {
		for z := int(f.Presets[p].Bag); z < int(f.Presets[p+1].Bag); z++ {
			gens := f.PresetGens[f.PresetBags[z].Gen:f.PresetBags[z+1].Gen]
			inst, ok := last(gens, genInstrument)
			if !ok { // a global zone
				presetZones[p] = append(presetZones[p], z)
				continue
			}
			if !covers(gens, keysOf[p]) {
				continue
			}
			presetZones[p] = append(presetZones[p], z)
			instKeys[inst] = instKeys[inst].add(keysOf[p])
		}
	}

	// The instrument zones kept, and the samples they play.
	insts := sortedKeys(instKeys)
	instZones := map[int][]int{}
	used := map[int]bool{}
	for _, i := range insts {
		for z := int(f.Instruments[i].Bag); z < int(f.Instruments[i+1].Bag); z++ {
			gens := f.InstGens[f.InstBags[z].Gen:f.InstBags[z+1].Gen]
			s, ok := last(gens, genSampleID)
			if !ok {
				instZones[i] = append(instZones[i], z)
				continue
			}
			if !covers(gens, instKeys[i].keys()) {
				continue
			}
			instZones[i] = append(instZones[i], z)
			used[s] = true
		}
	}
	for changed := true; changed; { // the other side of each stereo pair
		changed = false
		for s := range used {
			h := f.SampleHdrs[s]
			if h.Type&(rightSample|leftSample|linkedSample) != 0 && !used[int(h.Link)] {
				used[int(h.Link)], changed = true, true
			}
		}
	}

	out := &File{Info: f.Info}
	sampleAt := out.copySamples(f, sortedKeys(used))
	instAt := map[int]int{}
	for _, i := range insts {
		instAt[i] = len(out.Instruments)
		out.Instruments = append(out.Instruments, InstHeader{Name: f.Instruments[i].Name, Bag: uint16(len(out.InstBags))})
		for _, z := range instZones[i] {
			out.InstBags = append(out.InstBags, Bag{Gen: uint16(len(out.InstGens)), Mod: uint16(len(out.InstMods))})
			out.InstMods = append(out.InstMods, f.InstMods[f.InstBags[z].Mod:f.InstBags[z+1].Mod]...)
			for _, g := range f.InstGens[f.InstBags[z].Gen:f.InstBags[z+1].Gen] {
				if g.Oper == genSampleID {
					g.Amount = uint16(sampleAt[int(g.Amount)])
				}
				out.InstGens = append(out.InstGens, g)
			}
		}
	}
	out.Instruments = append(out.Instruments, InstHeader{Name: terminal("EOI"), Bag: uint16(len(out.InstBags))})
	out.InstBags = append(out.InstBags, Bag{Gen: uint16(len(out.InstGens)), Mod: uint16(len(out.InstMods))})
	out.InstMods = append(out.InstMods, Mod{})
	out.InstGens = append(out.InstGens, Gen{})

	for _, p := range presets {
		h := f.Presets[p]
		h.Bag = uint16(len(out.PresetBags))
		out.Presets = append(out.Presets, h)
		for _, z := range presetZones[p] {
			out.PresetBags = append(out.PresetBags, Bag{Gen: uint16(len(out.PresetGens)), Mod: uint16(len(out.PresetMods))})
			out.PresetMods = append(out.PresetMods, f.PresetMods[f.PresetBags[z].Mod:f.PresetBags[z+1].Mod]...)
			for _, g := range f.PresetGens[f.PresetBags[z].Gen:f.PresetBags[z+1].Gen] {
				if g.Oper == genInstrument {
					g.Amount = uint16(instAt[int(g.Amount)])
				}
				out.PresetGens = append(out.PresetGens, g)
			}
		}
	}
	out.Presets = append(out.Presets, PresetHeader{Name: terminal("EOP"), Bag: uint16(len(out.PresetBags))})
	out.PresetBags = append(out.PresetBags, Bag{Gen: uint16(len(out.PresetGens)), Mod: uint16(len(out.PresetMods))})
	out.PresetMods = append(out.PresetMods, Mod{})
	out.PresetGens = append(out.PresetGens, Gen{})
	return out, nil
}

// copySamples copies the samples `kept` of `f`, in their order, each
// followed by its padding of silence, and returns where each one went.
func (out *File) copySamples(f *File, kept []int) map[int]int {
	at := map[int]int{}
	for n, s := range kept {
		at[s] = n
	}
	for _, s := range kept {
		h := f.SampleHdrs[s]
		start := uint32(len(out.Samples) / 2)
		out.Samples = append(out.Samples, f.Samples[2*h.Start:2*h.End]...)
		out.Samples = append(out.Samples, make([]byte, 2*padding)...)
		if f.Samples24 != nil {
			out.Samples24 = append(out.Samples24, f.Samples24[h.Start:h.End]...)
			out.Samples24 = append(out.Samples24, make([]byte, padding)...)
		}
		shift := func(x uint32) uint32 { return x - h.Start + start }
		h.Start, h.End, h.LoopStart, h.LoopEnd = start, shift(h.End), shift(h.LoopStart), shift(h.LoopEnd)
		if h.Type&(rightSample|leftSample|linkedSample) != 0 {
			h.Link = uint16(at[int(h.Link)])
		}
		out.SampleHdrs = append(out.SampleHdrs, h)
	}
	out.SampleHdrs = append(out.SampleHdrs, SampleHeader{Name: terminal("EOS")})
	return at
}

// last returns the value of the zone's last generator when it is `oper`:
// the instrument of a preset zone, the sample of an instrument zone. A
// zone without it is a global zone.
func last(gens []Gen, oper uint16) (int, bool) {
	if len(gens) == 0 || gens[len(gens)-1].Oper != oper {
		return 0, false
	}
	return int(gens[len(gens)-1].Amount), true
}

// covers tells whether a zone sounds on one of `keys`, nil for any. A
// zone without a key range covers every key.
func covers(gens []Gen, keys []int) bool {
	if keys == nil {
		return true
	}
	lo, hi := 0, 127
	for _, g := range gens {
		if g.Oper == genKeyRange {
			lo, hi = int(g.Amount&0xff), int(g.Amount>>8)
		}
	}
	return slices.ContainsFunc(keys, func(k int) bool { return k >= lo && k <= hi })
}

// A keySet is the keys an instrument is played on: nil for all of them.
type keySet struct {
	all  bool
	some map[int]bool
}

// add widens the set by `keys`, nil for all. A nil set is a new one.
func (s *keySet) add(keys []int) *keySet {
	if s == nil {
		s = &keySet{some: map[int]bool{}}
	}
	if keys == nil {
		s.all = true
	}
	for _, k := range keys {
		s.some[k] = true
	}
	return s
}

// keys returns the keys of the set, nil for all.
func (s *keySet) keys() []int {
	if s.all {
		return nil
	}
	return sortedKeys(s.some)
}

func sortedKeys[V any](m map[int]V) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

// terminal names a terminal record.
func terminal(s string) [20]byte {
	var b [20]byte
	copy(b[:], s)
	return b
}
