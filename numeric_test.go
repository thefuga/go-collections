package collections

import (
	"testing"
)

func TestSum(t *testing.T) {
	testCases := []struct {
		description string
		sut         []int
		sum         int
	}{
		{
			"calling sum with empty slice",
			[]int{},
			0,
		},
		{
			"calling sum with slice with values",
			[]int{1, 2, 3, 4},
			10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			sum := Sum(tc.sut)

			if sum != tc.sum {
				t.Errorf("expected sum to be %d. got %d", tc.sum, sum)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	testCases := []struct {
		description string
		sut         []int
		avg         int
	}{
		{
			"calling Average with empty slice",
			[]int{},
			0,
		},
		{
			"calling Average with slice with values",
			[]int{2, 2, 2, 2},
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if avg := Average(tc.sut); avg != tc.avg {
				t.Errorf("expected average to be %d. got %d", tc.avg, avg)
			}
		})
	}
}

func TestMin(t *testing.T) {
	testCases := []struct {
		description string
		sut         []int
		min         int
	}{
		{
			"calling min with empty slice",
			[]int{},
			0,
		},
		{
			"calling min with slice with values",
			[]int{4, 1, 3, 2},
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if min := Min(tc.sut); min != tc.min {
				t.Errorf("expected min value to be %d. got %d", tc.min, min)
			}
		})
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		description string
		sut         []int
		max         int
	}{
		{
			"calling max with empty slice",
			[]int{},
			0,
		},
		{
			"calling max  with slice with values",
			[]int{1, 4, 3, 2},
			4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if max := Max(tc.sut); max != tc.max {
				t.Errorf("expected max value to be %d. got %d", tc.max, max)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	testCases := []struct {
		description string
		sut         []int
		median      float64
	}{
		{
			"1 and 2",
			[]int{1, 2},
			1.5,
		},
		{
			"1 through 3",
			[]int{1, 2, 3},
			2.0,
		},
		{
			"1 through 5",
			[]int{1, 2, 3, 4, 5},
			3.0,
		},
		{
			"1 through 6",
			[]int{1, 2, 3, 4, 5, 6},
			3.5,
		},
		{
			"wikepedia example 1",
			[]int{1, 3, 3, 6, 7, 8, 9},
			6.0,
		},
		{
			"wikepedia example 2",
			[]int{1, 2, 3, 4, 5, 7, 8, 9},
			4.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if median := Median(tc.sut); median != tc.median {
				t.Errorf("expected median to be %f. Got %f", tc.median, median)
			}
		})
	}
}
