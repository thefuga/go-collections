package collections

import (
	"github.com/thefuga/go-collections/errors"
	"github.com/thefuga/go-collections/internal"
)

// Sum sums all the values stored on the numeric slice and returns the result.
func Sum[T internal.Number](slice []T) T {
	var sum T

	for _, v := range slice {
		sum += v
	}

	return sum
}

// AverageE calculates the average value of the slice. Should the slice be empty,
// an instance of errors.EmptyCollectionError is returned.
func AverageE[T internal.Number](slice []T) (T, error) {
	if len(slice) == 0 {
		return *new(T), errors.NewEmptyCollectionError()
	}

	return Sum(slice) / T(len(slice)), nil
}

// Average uses AverageE, omitting the error.
func Average[T internal.Number](slice []T) T {
	avg, _ := AverageE(slice)
	return avg
}

// MinE returns the minimal value stored on the numeric slice. Should the slice be
// empty, an error is returned.
func MinE[T internal.Number](slice []T) (T, error) {
	min, err := FirstE(slice)

	if err != nil {
		return min, err
	}

	for _, v := range slice {
		if v < min {
			min = v
		}
	}

	return min, nil
}

// Min uses MinE, omitting the error.
func Min[T internal.Number](slice []T) T {
	min, _ := MinE(slice)
	return min
}

// MaxE returns the maximum value stored on the numeric slice. Should the slice be
// empty, an error is returned.
func MaxE[T internal.Number](slice []T) (T, error) {
	max, err := FirstE(slice)

	if err != nil {
		return max, err
	}

	for _, v := range slice {
		if v > max {
			max = v
		}
	}

	return max, nil
}

// Max uses MaxE, omitting the error.
func Max[T internal.Number](slice []T) T {
	max, _ := MaxE(slice)
	return max
}

// Median calculates and returns the median value of the slice.
func Median[T internal.Number](slice []T) float64 {
	Sort(slice, Asc[T]())

	halfway := int(len(slice) / 2)
	if len(slice)%2 == 0 {
		return float64(slice[halfway]+slice[halfway-1]) / 2.0
	}

	return float64(slice[halfway])
}
