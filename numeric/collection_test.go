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

func TestMin(t *testing.T) {
	collection := Collect(3, 2, 1, 0, -1, -2, -3)
	expectedMin := -3

	if actualMin := collection.Min(); actualMin != expectedMin {
		t.Errorf("expected min to be %d. Got %d", expectedMin, actualMin)
	}
}

func TestMinFloat(t *testing.T) {
	collection := Collect(3.3, 2.2, 1.1, 0.0, -1.1, -2.2, -3.3)
	expectedMin := -3.3

	if actualMin := collection.Min(); actualMin != expectedMin {
		t.Errorf("expected min to be %f. Got %f", expectedMin, actualMin)
	}
}
