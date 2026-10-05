package sf2

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

// A soundfont made for the test: three samples (a low tone, a high
// tone, a click), a melodic instrument on the first, and a kit striking
// each of them on its own key, played by three presets.
func testFile() *File {
	f := &File{Info: infoChunk()}
	var hdrs []SampleHeader
	for n, hz := range []float64{110, 440, 3000} {
		start := uint32(len(f.Samples) / 2)
		for i := range 4000 {
			v := int16(8000 * math.Sin(2*math.Pi*hz*float64(i)/22050) * math.Exp(-float64(i)/2000))
			f.Samples = append(f.Samples, byte(v), byte(uint16(v)>>8))
		}
		end := uint32(len(f.Samples) / 2)
		f.Samples = append(f.Samples, make([]byte, 2*padding)...)
		hdrs = append(hdrs, SampleHeader{Name: terminal("s" + string(rune('0'+n))), Start: start, End: end,
			LoopStart: start + 1000, LoopEnd: end - 100, Rate: 22050, OriginalPitch: 60, Type: 1})
	}
	f.SampleHdrs = append(hdrs, SampleHeader{Name: terminal("EOS")})

	// Instrument 0: the low tone everywhere. Instrument 1: the kit,
	// a global zone, then key 44 → s0, 51 → s1, 60 → s2.
	f.Instruments = []InstHeader{{Name: terminal("tone"), Bag: 0}, {Name: terminal("kit"), Bag: 1}, {Name: terminal("EOI"), Bag: 5}}
	f.InstBags = []Bag{{0, 0}, {1, 0}, {2, 0}, {4, 0}, {6, 0}, {8, 0}}
	keys := func(lo, hi int) Gen { return Gen{Oper: genKeyRange, Amount: uint16(lo | hi<<8)} }
	f.InstGens = []Gen{
		{Oper: genSampleID, Amount: 0},
		{Oper: 17, Amount: 0}, // a global zone: the pan
		keys(44, 44), {Oper: genSampleID, Amount: 0},
		keys(51, 51), {Oper: genSampleID, Amount: 1},
		keys(60, 60), {Oper: genSampleID, Amount: 2},
		{},
	}
	f.InstMods = []Mod{{}}

	// Presets 0:0 the tone, 0:1 the tone again, 128:0 the kit.
	f.Presets = []PresetHeader{
		{Name: terminal("Tone"), Bank: 0, Patch: 0, Bag: 0},
		{Name: terminal("Other"), Bank: 0, Patch: 1, Bag: 1},
		{Name: terminal("Kit"), Bank: 128, Patch: 0, Bag: 2},
		{Name: terminal("EOP"), Bag: 3},
	}
	f.PresetBags = []Bag{{0, 0}, {1, 0}, {2, 0}, {3, 0}}
	f.PresetGens = []Gen{{Oper: genInstrument, Amount: 0}, {Oper: genInstrument, Amount: 0}, {Oper: genInstrument, Amount: 1}, {}}
	f.PresetMods = []Mod{{}}
	return f
}

// infoChunk is the least INFO list a player reads: the version, the
// sound engine, the name.
func infoChunk() []byte {
	var b bytes.Buffer
	chunk(&b, "ifil", []byte{2, 0, 1, 0})
	chunk(&b, "isng", []byte("EMU8000\x00"))
	chunk(&b, "INAM", []byte("test\x00\x00"))
	return b.Bytes()
}

func encode(t *testing.T, f *File) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := Write(&b, f); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// Read then written, a soundfont is the same file, byte for byte.
func TestReadWrite(t *testing.T) {
	data := encode(t, testFile())
	f, err := Read(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if again := encode(t, f); !bytes.Equal(again, data) {
		t.Errorf("written again, %d bytes for %d", len(again), len(data))
	}
}

// render plays `key` on the preset `bank`:`patch` of the soundfont
// `data`, for a quarter of a second, and returns the left channel.
func render(t *testing.T, data []byte, bank, patch, key int) []float32 {
	t.Helper()
	sf, err := meltysynth.NewSoundFont(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	s, err := meltysynth.NewSynthesizer(sf, meltysynth.NewSynthesizerSettings(44100))
	if err != nil {
		t.Fatal(err)
	}
	ch := int32(0)
	if bank == 128 {
		ch = 9 // the drums' channel
	} else {
		s.ProcessMidiMessage(ch, 0xb0, 0, int32(bank))
	}
	s.ProcessMidiMessage(ch, 0xc0, int32(patch), 0)
	s.NoteOn(ch, int32(key), 100)
	left, right := make([]float32, 11025), make([]float32, 11025)
	s.Render(left, right)
	return left
}

// A played is a preset kept, and the keys the test plays on it.
type played struct {
	keep Keep
	keys []int
}

// slimmed cuts `data` down to the presets of `play`, and checks every
// key it plays sounds as it did, sample for sample.
func slimmed(t *testing.T, data []byte, play []played) []byte {
	t.Helper()
	f, err := Read(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	var keep []Keep
	for _, p := range play {
		keep = append(keep, p.keep)
	}
	s, err := Slim(f, keep)
	if err != nil {
		t.Fatal(err)
	}
	out := encode(t, s)
	for _, p := range play {
		k := p.keep
		for _, key := range p.keys {
			before, after := render(t, data, k.Bank, k.Patch, key), render(t, out, k.Bank, k.Patch, key)
			if !slices.Equal(before, after) {
				t.Errorf("preset %d:%d, key %d: not the same sound once slimmed", k.Bank, k.Patch, key)
			}
			if !slices.ContainsFunc(after, func(x float32) bool { return x != 0 }) {
				t.Errorf("preset %d:%d, key %d: silent", k.Bank, k.Patch, key)
			}
		}
	}
	return out
}

// The tone and two keys of the kit kept: the other preset goes, and the
// click with it; what is left sounds the same.
func TestSlim(t *testing.T) {
	data := encode(t, testFile())
	tone, kit := Keep{Bank: 0, Patch: 0}, Keep{Bank: 128, Patch: 0, Keys: []int{44, 51}}
	out := slimmed(t, data, []played{{tone, []int{48, 60, 72}}, {kit, []int{44, 51}}})
	f, err := Read(bytes.NewReader(out))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, p := range f.Presets {
		names = append(names, name(p.Name))
	}
	if !slices.Equal(names, []string{"Tone", "Kit", "EOP"}) {
		t.Errorf("presets %q", names)
	}
	if n := len(f.SampleHdrs) - 1; n != 2 {
		t.Errorf("%d samples kept, want 2: the low tone and the high one", n)
	}
	if len(out) >= len(data) {
		t.Errorf("%d bytes, no smaller than %d", len(out), len(data))
	}
}

// On GeneralUser GS, where make sounds puts it: what Walk with me
// plays, the double bass, the piano and four keys of the Jazz kit,
// sounds the same once slimmed.
func TestGeneralUser(t *testing.T) {
	dir, err := os.UserCacheDir()
	if err != nil {
		t.Skip(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "gohar", "GeneralUser-GS.sf2"))
	if err != nil {
		t.Skip("no GeneralUser GS: make sounds")
	}
	bass, piano, kit := Keep{Bank: 0, Patch: 32}, Keep{Bank: 0, Patch: 0}, Keep{Bank: 128, Patch: 32, Keys: []int{44, 51, 76, 77}}
	out := slimmed(t, data, []played{
		{bass, []int{28, 40, 52}}, {piano, []int{48, 60, 72, 84}}, {kit, []int{44, 51, 76, 77}},
	})
	t.Logf("%d bytes, from %d", len(out), len(data))
}
