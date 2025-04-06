package keyboard

import (
	. "github.com/ArnaudCalmettes/gohar"
)

type Keyboard struct {
	Keys []Key
}

type Key struct {
	Pitch
	Flags KeyFlag
}

type KeyFlag uint8

const (
	keyFlagBlack KeyFlag = 1 << iota
	keyFlagPressed
)

func (k Key) isBlack() bool {
	return k.Flags&keyFlagBlack != 0
}

func (k Key) IsPressed() bool {
	return k.Flags&keyFlagPressed != 0
}

func New(lowest, highest Pitch) *Keyboard {
	lowest, highest, ambitus := adjustAmbitus(lowest, highest)
	keys := make([]Key, 0, int(ambitus)+1)
	for pitch := lowest; pitch <= highest; pitch++ {
		var flag KeyFlag
		if isBlackKey(pitch) {
			flag |= keyFlagBlack
		}
		keys = append(keys, Key{Pitch: pitch, Flags: flag})
	}
	return &Keyboard{Keys: keys}
}

func (k *Keyboard) Press(pitches ...Pitch) {
	lowest := k.Keys[0].Pitch
	for _, pitch := range pitches {
		i := int(pitch - lowest)
		if i >= 0 && i < len(k.Keys) {
			k.Keys[pitch-lowest].Flags |= keyFlagPressed
		}
	}
}

// Adjust boundaries so the leftmost and rightmost keys are white
// and the keyboard is at least one octave wide.
func adjustAmbitus(low, high Pitch) (lowest, highest, ambitus Pitch) {
	lowest, highest = low, high
	if ambitus = high - low; ambitus < 0 {
		lowest, highest, ambitus = highest, lowest, -ambitus
	}
	if isBlackKey(lowest) {
		lowest--
		ambitus++
	}
	if isBlackKey(highest) {
		highest++
		ambitus++
	}
	if ambitus < PitchDiffOctave {
		highest = lowest + PitchDiffOctave
		ambitus = PitchDiffOctave
	}
	return
}

func isBlackKey(pitch Pitch) bool {
	switch pitch.Normalize() {
	case PitchAFlat, PitchBFlat, PitchDFlat, PitchEFlat, PitchGFlat:
		return true
	default:
		return false
	}
}
