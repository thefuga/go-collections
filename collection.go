package collections

import (
	"sort"

	"github.com/thefuga/go-collections/errors"
)

func Sum[T Number](slice []T) T {
	var sum T

	for _, v := range slice {
		sum += v
	}

	return sum
}

func AverageE[T Number](slice []T) (T, error) {
	if len(slice) == 0 {
		return *new(T), errors.NewEmptyCollectionError()
	}

	return Sum(slice) / T(len(slice)), nil
}

func Average[T Number](slice []T) T {
	avg, _ := AverageE(slice)
	return avg
}

func First[T Number](slice []T) (T, error) {
	if len(slice) == 0 {
		return *new(T), errors.NewEmptyCollectionError()
	}

	return slice[0], nil
}

func MinE[T Number](slice []T) (T, error) {
	min, err := First(slice)

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

func Min[T Number](slice []T) T {
	min, _ := MinE(slice)
	return min
}

func MaxE[T Number](slice []T) (T, error) {
	max, err := First(slice)

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

func Max[T Number](slice []T) T {
	max, _ := MaxE(slice)
	return max
}

func Median[T Number](slice []T) float64 {
	Sort(slice, Asc[T]())

	halfway := int(len(slice) / 2)
	if len(slice)%2 == 0 {
		return float64(slice[halfway]+slice[halfway-1]) / 2.0
	}

	return float64(slice[halfway])
}

func Sort[T Number](slice []T, f func(current, next T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return f(slice[i], slice[j])
	})
}
