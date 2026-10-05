// Package sf2 reads and writes SoundFont 2 files, table by table, and
// cuts a soundfont down to the instruments a game plays (see Slim).
//
// It does not play anything: the synth plays soundfonts through
// meltysynth, which reads SF2 but does not write it. The layout is the
// one of the SoundFont 2.04 specification: a RIFF file holding three
// lists, INFO (names and versions), sdta (the samples, 16-bit PCM, with
// an optional sm24 chunk for 24-bit), and pdta (nine tables of records
// pointing at one another by index).
package sf2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// A File is a soundfont as its tables hold it. Each table keeps the
// terminal record the specification asks for ("EOP", "EOI", "EOS" and
// the zero records): the zones of record i run from its index to the
// index of record i+1.
type File struct {
	Info      []byte // the body of the INFO list, kept as read
	Samples   []byte // smpl: 16-bit little-endian PCM, every sample end to end
	Samples24 []byte // sm24: the low byte of 24-bit samples, or nil

	Presets     []PresetHeader
	PresetBags  []Bag
	PresetMods  []Mod
	PresetGens  []Gen
	Instruments []InstHeader
	InstBags    []Bag
	InstMods    []Mod
	InstGens    []Gen
	SampleHdrs  []SampleHeader
}

// A PresetHeader is a phdr record: a General MIDI program, bank and
// patch, and its first zone.
type PresetHeader struct {
	Name       [20]byte
	Patch      uint16
	Bank       uint16
	Bag        uint16 // its first zone in PresetBags
	Library    uint32
	Genre      uint32
	Morphology uint32
}

// A Bag is a zone: its first generator and its first modulator.
type Bag struct {
	Gen, Mod uint16
}

// A Mod is a modulator, copied as it is.
type Mod struct {
	Src, Dest uint16
	Amount    int16
	AmountSrc uint16
	Transform uint16
}

// A Gen is a generator: a parameter of a zone and its value. The value
// is a union: a signed amount, or a range of two bytes, low then high.
type Gen struct {
	Oper   uint16
	Amount uint16
}

// An InstHeader is an inst record: an instrument and its first zone.
type InstHeader struct {
	Name [20]byte
	Bag  uint16 // its first zone in InstBags
}

// A SampleHeader is a shdr record: where a sample lies in Samples, in
// sample points, and how it loops and sounds.
type SampleHeader struct {
	Name            [20]byte
	Start, End      uint32
	LoopStart       uint32
	LoopEnd         uint32
	Rate            uint32
	OriginalPitch   uint8
	PitchCorrection int8
	Link            uint16 // the other side of a stereo pair
	Type            uint16
}

// The generators Slim reads; the others are copied as they are.
const (
	genInstrument = 41 // a preset zone's instrument
	genKeyRange   = 43 // a zone's keys, low and high byte
	genSampleID   = 53 // an instrument zone's sample
)

// The sample types that point at another sample through Link.
const (
	rightSample  = 2
	leftSample   = 4
	linkedSample = 8
)

// The zero points the specification asks after each sample, so that a
// player interpolating past the end reads silence.
const padding = 46

// name reads the name of a record, up to its first zero byte.
func name(b [20]byte) string {
	if i := bytes.IndexByte(b[:], 0); i >= 0 {
		return string(b[:i])
	}
	return string(b[:])
}

// Read reads a soundfont, its three lists in the order the
// specification gives them.
func Read(r io.Reader) (*File, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(data) < 12 || string(data[:4]) != "RIFF" || string(data[8:12]) != "sfbk" {
		return nil, errors.New("sf2: not a soundfont")
	}
	f := &File{}
	body := data[12:]
	for len(body) >= 8 {
		id, size := string(body[:4]), int(binary.LittleEndian.Uint32(body[4:8]))
		if 8+size > len(body) {
			return nil, fmt.Errorf("sf2: chunk %q runs past the end", id)
		}
		chunk := body[8 : 8+size]
		body = body[8+size+size%2:]
		if id != "LIST" || len(chunk) < 4 {
			continue
		}
		switch string(chunk[:4]) {
		case "INFO":
			f.Info = append([]byte(nil), chunk[4:]...)
		case "sdta":
			err = f.readSamples(chunk[4:])
		case "pdta":
			err = f.readTables(chunk[4:])
		}
		if err != nil {
			return nil, err
		}
	}
	if f.Presets == nil || f.SampleHdrs == nil {
		return nil, errors.New("sf2: no pdta list")
	}
	return f, nil
}

// subchunks calls `fn` on each chunk of a list's body.
func subchunks(body []byte, fn func(id string, data []byte) error) error {
	for len(body) >= 8 {
		id, size := string(body[:4]), int(binary.LittleEndian.Uint32(body[4:8]))
		if 8+size > len(body) {
			return fmt.Errorf("sf2: chunk %q runs past the end", id)
		}
		if err := fn(id, body[8:8+size]); err != nil {
			return err
		}
		body = body[8+size+size%2:]
	}
	return nil
}

func (f *File) readSamples(body []byte) error {
	return subchunks(body, func(id string, data []byte) error {
		switch id {
		case "smpl":
			f.Samples = append([]byte(nil), data...)
		case "sm24":
			f.Samples24 = append([]byte(nil), data...)
		}
		return nil
	})
}

// records reads a table of fixed-size records into `out`, a pointer to
// a slice.
func records[T any](id string, data []byte, out *[]T) error {
	var zero T
	size := binary.Size(zero)
	if len(data)%size != 0 {
		return fmt.Errorf("sf2: %s is %d bytes, not a whole number of %d-byte records", id, len(data), size)
	}
	*out = make([]T, len(data)/size)
	return binary.Read(bytes.NewReader(data), binary.LittleEndian, *out)
}

func (f *File) readTables(body []byte) error {
	return subchunks(body, func(id string, data []byte) error {
		switch id {
		case "phdr":
			return records(id, data, &f.Presets)
		case "pbag":
			return records(id, data, &f.PresetBags)
		case "pmod":
			return records(id, data, &f.PresetMods)
		case "pgen":
			return records(id, data, &f.PresetGens)
		case "inst":
			return records(id, data, &f.Instruments)
		case "ibag":
			return records(id, data, &f.InstBags)
		case "imod":
			return records(id, data, &f.InstMods)
		case "igen":
			return records(id, data, &f.InstGens)
		case "shdr":
			return records(id, data, &f.SampleHdrs)
		}
		return nil
	})
}

// Write writes the soundfont, its lists in the specification's order.
func Write(w io.Writer, f *File) error {
	var sdta, pdta bytes.Buffer
	sdta.WriteString("sdta")
	chunk(&sdta, "smpl", f.Samples)
	if f.Samples24 != nil {
		chunk(&sdta, "sm24", f.Samples24)
	}
	pdta.WriteString("pdta")
	for _, t := range []struct {
		id   string
		data any
	}{
		{"phdr", f.Presets}, {"pbag", f.PresetBags}, {"pmod", f.PresetMods}, {"pgen", f.PresetGens},
		{"inst", f.Instruments}, {"ibag", f.InstBags}, {"imod", f.InstMods}, {"igen", f.InstGens},
		{"shdr", f.SampleHdrs},
	} {
		var b bytes.Buffer
		if err := binary.Write(&b, binary.LittleEndian, t.data); err != nil {
			return fmt.Errorf("sf2: writing %s: %w", t.id, err)
		}
		chunk(&pdta, t.id, b.Bytes())
	}
	var riff bytes.Buffer
	riff.WriteString("sfbk")
	chunk(&riff, "LIST", append([]byte("INFO"), f.Info...))
	chunk(&riff, "LIST", sdta.Bytes())
	chunk(&riff, "LIST", pdta.Bytes())
	var out bytes.Buffer
	chunk(&out, "RIFF", riff.Bytes())
	_, err := w.Write(out.Bytes())
	return err
}

// chunk writes a chunk: its id, its size, its data, padded to an even
// size as RIFF asks.
func chunk(b *bytes.Buffer, id string, data []byte) {
	b.WriteString(id)
	binary.Write(b, binary.LittleEndian, uint32(len(data)))
	b.Write(data)
	if len(data)%2 == 1 {
		b.WriteByte(0)
	}
}
