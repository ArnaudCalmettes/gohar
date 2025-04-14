package gohar

// A Pitch corresponds to a distinct vibration frequency.
// Pitches are represented as a signed integer value in the range [-127; 128]. This range
// is more than enough to model a piano keyboard, whose range is contained within [-50; 50].
//
// Pitch unit is the semitone or half-step.
// An octave corresponds to 12 semitones.
//
// By convention, pitch 0 corresponds to middle C on a piano keyboard.
type Pitch int8

// Add transposes the pitch up or down given interval (pitch difference) value.
func (p Pitch) Add(pitchDiff Pitch) Pitch {
	return p + pitchDiff
}

// Normalize transposes the pitch to the reference [0; 11] octave.
func (p Pitch) Normalize() Pitch {
	p %= Pitch(PitchDiffOctave)
	if p < 0 {
		return p.Add(PitchDiffOctave)
	}
	return p
}

// GetOctave returns the pitch's octave.
func (p Pitch) GetOctave() int8 {
	if p >= 0 {
		return int8(p) / 12
	} else {
		return int8(p+1)/12 - 1
	}
}

// AtOctave transposes the pitch to given octave.
func (p Pitch) AtOctave(oct int8) Pitch {
	return p.Normalize() + Pitch(12*oct)
}

const (
	PitchC            Pitch = 0
	PitchBSharp       Pitch = 0
	PitchDDoubleFlat  Pitch = 0
	PitchBDoubleSharp Pitch = 1
	PitchCSharp       Pitch = 1
	PitchDFlat        Pitch = 1
	PitchCDoubleSharp Pitch = 2
	PitchD            Pitch = 2
	PitchEDoubleFlat  Pitch = 2
	PitchDSharp       Pitch = 3
	PitchEFlat        Pitch = 3
	PitchDDoubleSharp Pitch = 4
	PitchE            Pitch = 4
	PitchFFlat        Pitch = 4
	PitchFDoubleFlat  Pitch = 4
	PitchESharp       Pitch = 5
	PitchF            Pitch = 5
	PitchGDoubleFlat  Pitch = 5
	PitchEDoubleSharp Pitch = 6
	PitchFSharp       Pitch = 6
	PitchGFlat        Pitch = 6
	PitchFDoubleSharp Pitch = 7
	PitchG            Pitch = 7
	PitchADoubleFlat  Pitch = 7
	PitchGSharp       Pitch = 8
	PitchAFlat        Pitch = 8
	PitchGDoubleSharp Pitch = 9
	PitchA            Pitch = 9
	PitchBDoubleFlat  Pitch = 9
	PitchASharp       Pitch = 10
	PitchBFlat        Pitch = 10
	PitchCDoubleFlat  Pitch = 10
	PitchADoubleSharp Pitch = 11
	PitchB            Pitch = 11
	PitchCFlat        Pitch = 11
)

const (
	PitchDiffUnisson           Pitch = 0
	PitchDiffPerfectUnisson    Pitch = 0
	PitchDiffHalfStep          Pitch = 1
	PitchDiffSemitone          Pitch = 1
	PitchDiffMinorSecond       Pitch = 1
	PitchDiffFullStep          Pitch = 2
	PitchDiffTone              Pitch = 2
	PitchDiffMajorSecond       Pitch = 2
	PitchDiffDiminishedThird   Pitch = 2
	PitchDiffAugmentedSecond   Pitch = 3
	PitchDiffMinorThird        Pitch = 3
	PitchDiffMajorThird        Pitch = 4
	PitchDiffDiminishedFourth  Pitch = 4
	PitchDiffFourth            Pitch = 5
	PitchDiffPerfectFourth     Pitch = 5
	PitchDiffAugmentedFourth   Pitch = 6
	PitchDiffDiminishedFifth   Pitch = 6
	PitchDiffFifth             Pitch = 7
	PitchDiffPerfectFifth      Pitch = 7
	PitchDiffAugmentedFifth    Pitch = 8
	PitchDiffMinorSixth        Pitch = 8
	PitchDiffMajorSixth        Pitch = 9
	PitchDiffDiminishedSeventh Pitch = 9
	PitchDiffMinorSeventh      Pitch = 10
	PitchDiffMajorSeventh      Pitch = 11
	PitchDiffOctave            Pitch = 12
	PitchDiffPerfectOctave     Pitch = 12
	PitchDiffMinorNinth        Pitch = 13
	PitchDiffMajorNinth        Pitch = 14
	PitchDiffAugmentedNinth    Pitch = 15
	PitchDiffMinorTenth        Pitch = 15
	PitchDiffMajorTenth        Pitch = 16
	PitchDiffEleventh          Pitch = 17
	PitchDiffPerfectEleventh   Pitch = 17
	PitchDiffAugmentedEleventh Pitch = 18
	PitchDiffMinorThirteenth   Pitch = 20
	PitchDiffMajorThirteenth   Pitch = 21
	PitchDiffMajorFourteenth   Pitch = 23
)
