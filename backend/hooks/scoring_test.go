package hooks

import (
	"testing"
)

func TestCalculateMatchPoints(t *testing.T) {
	tests := []struct {
		name       string
		phase      string
		actHome    int
		actAway    int
		predHome   int
		predAway   int
		wantPoints int
	}{
		// --- GROUP PHASE TESTS (+2 points for result, +1 bonus for exact score) ---
		{
			name:    "Group Stage: Exact Home Win",
			phase:   "group",
			actHome: 2, actAway: 1,
			predHome: 2, predAway: 1,
			wantPoints: 3, // 2 (result) + 1 (exact)
		},
		{
			name:    "Group Stage: Correct Result Only (Home Win)",
			phase:   "group",
			actHome: 3, actAway: 1,
			predHome: 1, predAway: 0,
			wantPoints: 2, // 2 (result) + 0 (exact)
		},
		{
			name:    "Group Stage: Exact Draw",
			phase:   "group",
			actHome: 1, actAway: 1,
			predHome: 1, predAway: 1,
			wantPoints: 3, // 2 (result) + 1 (exact)
		},
		{
			name:    "Group Stage: Correct Result Only (Draw)",
			phase:   "group",
			actHome: 2, actAway: 2,
			predHome: 0, predAway: 0,
			wantPoints: 2, // 2 (result) + 0 (exact)
		},
		{
			name:    "Group Stage: Completely Incorrect",
			phase:   "group",
			actHome: 0, actAway: 2,
			predHome: 1, predAway: 0,
			wantPoints: 0,
		},

		// --- ROUND OF 32 / 16 TESTS (+3 points for result, +2 bonus for exact score) ---
		{
			name:    "R32: Exact Away Win",
			phase:   "r32",
			actHome: 1, actAway: 3,
			predHome: 1, predAway: 3,
			wantPoints: 5, // 3 (result) + 2 (exact)
		},
		{
			name:    "R16: Correct Result Only (Away Win)",
			phase:   "r16",
			actHome: 0, actAway: 2,
			predHome: 1, predAway: 3,
			wantPoints: 3, // 3 (result) + 0 (exact)
		},

		// --- QUARTER / SEMI-FINAL TESTS (+5 points for result, +3 bonus for exact score) ---
		{
			name:    "QF: Exact Home Win",
			phase:   "qf",
			actHome: 1, actAway: 0,
			predHome: 1, predAway: 0,
			wantPoints: 8, // 5 (result) + 3 (exact)
		},
		{
			name:    "SF: Correct Result Only (Home Win)",
			phase:   "sf",
			actHome: 4, actAway: 2,
			predHome: 2, predAway: 1,
			wantPoints: 5, // 5 (result) + 0 (exact)
		},

		// --- FINAL / THIRD PLACE TESTS (+8 points for result, +5 bonus for exact score) ---
		{
			name:    "Final: Exact Draw",
			phase:   "final",
			actHome: 3, actAway: 3,
			predHome: 3, predAway: 3,
			wantPoints: 13, // 8 (result) + 5 (exact)
		},
		{
			name:    "Final: Correct Result Only (Draw)",
			phase:   "final",
			actHome: 2, actAway: 2,
			predHome: 1, predAway: 1,
			wantPoints: 8, // 8 (result) + 0 (exact)
		},
		{
			name:    "Final: Completely Incorrect",
			phase:   "final",
			actHome: 1, actAway: 0,
			predHome: 0, predAway: 2,
			wantPoints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateMatchPoints(tt.phase, tt.actHome, tt.actAway, tt.predHome, tt.predAway)
			if got != tt.wantPoints {
				t.Errorf("calculateMatchPoints(%s, actual: %d-%d, predicted: %d-%d) = %d; want %d",
					tt.phase, tt.actHome, tt.actAway, tt.predHome, tt.predAway, got, tt.wantPoints,
				)
			}
		})
	}
}
