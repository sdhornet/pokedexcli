// ABOUTME: Tests for calculateCatchChance, the catch-probability clamp in command_catch.go.
// ABOUTME: Verifies the chance stays within [10, 90] so a catch is never impossible or guaranteed.

package main

import (
	"testing"
)

// calculateCatchChance clamps the catch probability to the range [10, 90].
// Because of that, for any base experience a long run of attempts must produce
// both at least one catch and at least one escape: the floor keeps high
// base-experience pokemon catchable, the ceiling keeps low base-experience
// pokemon escapable. With this many trials a false failure is astronomically
// unlikely (worst case ~0.9^trials).
func TestCalculateCatchChanceStaysBounded(t *testing.T) {
	const trials = 10000

	cases := []struct {
		name           string
		baseExperience int
	}{
		{
			name:           "floor clamp for very high base experience",
			baseExperience: 10000, // raw chance 100 - 5000 = -4900, clamped up to 10
		},
		{
			name:           "ceiling clamp for zero base experience",
			baseExperience: 0, // raw chance 100, clamped down to 90
		},
		{
			name:           "mid range needs no clamping",
			baseExperience: 40, // raw chance 100 - 20 = 80, already inside [10, 90]
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			catches := 0
			escapes := 0
			for i := 0; i < trials; i++ {
				if calculateCatchChance(c.baseExperience) {
					catches++
				} else {
					escapes++
				}
			}

			if catches == 0 {
				t.Errorf("baseExperience %d: got 0 catches in %d trials; catch chance should never be 0%% (floor clamp at 10)", c.baseExperience, trials)
			}
			if escapes == 0 {
				t.Errorf("baseExperience %d: got 0 escapes in %d trials; catch chance should never be 100%% (ceiling clamp at 90)", c.baseExperience, trials)
			}
		})
	}
}
