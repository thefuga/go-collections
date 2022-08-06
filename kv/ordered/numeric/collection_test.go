package numeric

import (
	"reflect"
	"testing"
)

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

func TestDuplicates(t *testing.T) {
	testCases := []struct {
		name       string
		collection Collection[int, int]
		expected   []int
	}{
		{
			"no duplicates",
			Collect(1, 2, 3, 4),
			[]int{},
		},
		{
			"1 appearing twice",
			Collect(1, 2, 1, 3, 4),
			[]int{1},
		},
		{
			"1 and 2 appearing twice",
			Collect(1, 2, 1, 3, 2),
			[]int{1, 2},
		},
		{
			"every element appearing twice",
			Collect(1, 2, 3, 1, 2, 3),
			[]int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.collection.Duplicates(); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}
