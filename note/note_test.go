package note_test

import (
	"testing"

	"github.com/heucuva/comparison"

	"github.com/gotracker/playback/period"
)

// testPeriod defines a sampler period that follows the Amiga-style approach of note
// definition. Useful in calculating resampling.
type testPeriod = period.Amiga

func periodCompareTest(t *testing.T, lhs, rhs testPeriod, expected comparison.Spaceship) {
	t.Helper()

	if lhs.Compare(rhs) != expected {
		t.Fatalf("%v <=> %v was not %v", lhs, rhs, expected)
	}
}

func TestPeriodCompare(t *testing.T) {
	lhs1 := testPeriod(1)
	rhs1 := testPeriod(1)
	periodCompareTest(t, lhs1, rhs1, comparison.SpaceshipEqual)

	lhs2 := testPeriod(1)
	rhs2 := testPeriod(2)
	periodCompareTest(t, lhs2, rhs2, comparison.SpaceshipLeftGreater)

	lhs3 := testPeriod(2)
	rhs3 := testPeriod(1)
	periodCompareTest(t, lhs3, rhs3, comparison.SpaceshipRightGreater)
}
