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

func TestMedian(t *testing.T) {
	testCases := []struct {
		name           string
		collection     Collection[int, int]
		expectedMedian float64
	}{
		{
			name:           "1 and 2",
			collection:     Collect(1, 2),
			expectedMedian: 1.5,
		},
		{
			name:           "1 through 3",
			collection:     Collect(1, 2, 3),
			expectedMedian: 2.0,
		},
		{
			name:           "1 through 5",
			collection:     Collect(1, 2, 3, 4, 5),
			expectedMedian: 3.0,
		},
		{
			name:           "1 through 6",
			collection:     Collect(1, 2, 3, 4, 5, 6),
			expectedMedian: 3.5,
		},
		{
			name:           "wikepedia example 1",
			collection:     Collect(1, 3, 3, 6, 7, 8, 9),
			expectedMedian: 6.0,
		},
		{
			name:           "wikepedia example 2",
			collection:     Collect(1, 2, 3, 4, 5, 7, 8, 9),
			expectedMedian: 4.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actualMedian := tc.collection.Median(); actualMedian != tc.expectedMedian {
				t.Errorf("expected median to be %f. Got %f", tc.expectedMedian, actualMedian)
			}
		})
	}
}
