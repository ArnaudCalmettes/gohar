package keyboard

import "testing"

// A key pressed gently still sounds: under the floor, at the floor; a
// key struck harder keeps its own velocity.
func TestVelocityFloor(t *testing.T) {
	for _, tc := range []struct {
		midi uint8
		want float64
	}{
		{1, MinVelocity},
		{20, MinVelocity},
		{64, 64.0 / 127},
		{127, 1},
	} {
		if got := velocityOf(tc.midi); got != tc.want {
			t.Errorf("velocity %d: got %.3f, want %.3f", tc.midi, got, tc.want)
		}
	}
}
