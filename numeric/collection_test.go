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
