package numeric

import "testing"

func TestAverageInts(t *testing.T) {
	collection := Collect(1, 2, 3, 4, 5)
	expectedAvg := 3

	if actualAvg := collection.Average(); actualAvg != expectedAvg {
		t.Errorf("Expected average to be %d. Got %d", expectedAvg, actualAvg)
	}
}

func TestAverageFloats(t *testing.T) {
	collection := Collect(1.1, 2.2, 3.3, 4.4, 5.5)
	expectedAvg := 3.3

	if actualAvg := collection.Average(); actualAvg != expectedAvg {
		t.Errorf("Expected average to be %f. Got %f", expectedAvg, actualAvg)
	}
}

func TestSum(t *testing.T) {
	collection := Collect(1, 2, 3, 4, 5)
	expectedSum := 15

	if actualSum := collection.Sum(); actualSum != expectedSum {
		t.Errorf("expected sum to be %d. Got %d", expectedSum, actualSum)
	}
}

func TestMax(t *testing.T) {
	collection := Collect(-3, -2, -1, 0, 1, 2, 3)
	expectedMax := 3

	if actualMax := collection.Max(); actualMax != expectedMax {
		t.Errorf("expected min to be %d. Got %d", expectedMax, actualMax)
	}
}

func TestMaxFloat(t *testing.T) {
	collection := Collect(-3.3, -2.2, -1.1, 0.0, 1.1, 2.2, 3.3)
	expectedMax := 3.3

	if actualMax := collection.Max(); actualMax != expectedMax {
		t.Errorf("expected min to be %f. Got %f", expectedMax, actualMax)
	}
}
