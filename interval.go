package gohar

type Interval struct {
	ScaleDiff int8
	PitchDiff Pitch
}

var (
	IntUnisson           = Interval{}
	IntMinorSecond       = Interval{1, 1}
	IntMajorSecond       = Interval{1, 2}
	IntAugmentedSecond   = Interval{1, 3}
	IntMinorThird        = Interval{2, 3}
	IntMajorThird        = Interval{2, 4}
	IntDiminishedFourth  = Interval{3, 4}
	IntPerfectFourth     = Interval{3, 5}
	IntAugmentedFourth   = Interval{3, 6}
	IntDiminishedFifth   = Interval{4, 6}
	IntPerfectFifth      = Interval{4, 7}
	IntAugmentedFifth    = Interval{4, 8}
	IntMinorSixth        = Interval{5, 8}
	IntMajorSixth        = Interval{5, 9}
	IntDiminishedSeventh = Interval{6, 9}
	IntAugmentedSixth    = Interval{5, 10}
	IntMinorSeventh      = Interval{6, 10}
	IntMajorSeventh      = Interval{6, 11}
	IntOctave            = Interval{7, 12}
	IntMinorNinth        = Interval{8, 13}
	IntMajorNinth        = Interval{8, 14}
	IntAugmentedNinth    = Interval{8, 15}
	IntMinorTenth        = Interval{9, 15}
	IntMajorTenth        = Interval{9, 16}
	IntPerfectEleventh   = Interval{10, 17}
	IntAugmentedEleventh = Interval{10, 18}
	IntMinorThirteenth   = Interval{12, 20}
	IntMajorThirteenth   = Interval{12, 21}
	IntMajorFourteenth   = Interval{13, 23}
)

func (i Interval) Down() Interval {
	return Interval{-i.ScaleDiff, -i.PitchDiff}
}
