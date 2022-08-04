package collections

import (
	"reflect"
	"sort"

	"github.com/thefuga/go-collections/errors"
)

func Get[T any](i int, slice []T) T {
	v, _ := GetE(i, slice)
	return v
}

func GetE[T any](i int, slice []T) (T, error) {
	if len(slice) == 0 {
		return *new(T), errors.NewEmptyCollectionError(
			errors.NewValueNotFoundError(),
		)
	}

	if i < 0 || len(slice) <= i {
		return *new(T), errors.NewIndexOutOfBoundsError(
			errors.NewValueNotFoundError(),
		)
	}

	return slice[i], nil
}

func First[T any](slice []T) T {
	v, _ := FirstE(slice)
	return v
}

func FirstE[T any](slice []T) (T, error) {
	return GetE(0, slice)
}

func Push[T any](v T, slice []T) []T {
	return append(slice, v)
}

func Put[T any](i int, v T, slice []T) []T {
	if len(slice) == 0 {
		slice = make([]T, i+1)
		slice[i] = v
		return slice
	}

	slice = append(slice, *new(T))
	copy(slice[i+1:], slice[i:])
	slice[i] = v

	return slice
}

func Pop[T any](slice *[]T) T {
	v, _ := PopE(slice)
	return v
}

func PopE[T any](slice *[]T) (T, error) {
	v, err := LastE(*slice)

	if err != nil {
		return v, err
	}

	*slice = (*slice)[:len(*slice)-1]

	return v, nil
}

func Shift[T any](slice *[]T) T {
	v, _ := ShiftE(slice)
	return v
}

func ShiftE[T any](slice *[]T) (T, error) {
	v, err := FirstE(*slice)

	if err != nil {
		return v, err
	}

	*slice = (*slice)[1:]

	return v, nil
}

func Last[T any](slice []T) T {
	v, _ := GetE(len(slice)-1, slice)
	return v
}

func LastE[T any](slice []T) (T, error) {
	return GetE(len(slice)-1, slice)
}

func Each[T any](f func(i int, v T), slice []T) {
	for i, v := range slice {
		f(i, v)
	}
}

func Search[T any](v T, slice []T) int {
	i, _ := SearchE(v, slice)
	return i
}

func SearchE[T any](v T, slice []T) (int, error) {
	for i := range slice {
		if reflect.DeepEqual(slice[i], v) {
			return i, nil
		}
	}

	return -1, errors.NewValueNotFoundError()
}

func Map[T any](f func(i int, v T) T, slice []T) []T {
	mappedValues := make([]T, 0, len(slice))

	Each(func(i int, v T) {
		mappedValues = Push(f(i, v), mappedValues)
	}, slice)

	return mappedValues
}

func Sort[T any](slice []T, f func(current, next T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return f(slice[i], slice[j])
	})
}

func Copy[V any](slice []V) []V {
	copied := make([]V, len(slice))
	copy(slice, copied)
	return copied
}

func Cut[V any](slice *[]V, i int, optionalJ ...int) []V {
	cutted, _ := CutE(slice, i, optionalJ...)
	return cutted
}

func CutE[V any](slice *[]V, i int, optionalJ ...int) ([]V, error) {
	sliceLen := len(*slice)
	i, j := bounds(i, optionalJ...)
	if i > sliceLen || j > sliceLen {
		return nil, errors.NewIndexOutOfBoundsError()
	}

	cutted := make([]V, j-i)
	copy(cutted, (*slice)[i:])

	copy((*slice)[i:], (*slice)[j:])
	for k, n := sliceLen-j+i, sliceLen; k < n; k++ {
		(*slice)[k] = *new(V)
	}
	*slice = (*slice)[:sliceLen-j+i]

	return cutted, nil
}

func DeleteE[V any](slice *[]V, i int, optionalJ ...int) error {
	sliceLen := len(*slice)

	i, j := bounds(i, optionalJ...)
	if i < 0 || i >= sliceLen || j >= sliceLen {
		return errors.NewIndexOutOfBoundsError()
	}

	copy((*slice)[i:], (*slice)[i+1:])
	(*slice)[sliceLen-1] = *new(V)
	(*slice) = (*slice)[:sliceLen-1]

	return nil
}

func ForgetE[V any](slice *[]V, i int, optionalJ ...int) error {
	return DeleteE(slice, i, optionalJ...)
}

func Tally[T comparable](slice []T) map[T]int {
	m := map[T]int{}
	for _, v := range slice {
		m[v]++
	}
	return m
}

// Mode returns the values that appear most often in the slice. Order is not guaranteed.
func Mode[T comparable](slice []T) []T {
	maxCount := 0
	mode := []T{}

	for v, count := range Tally(slice) {
		if count > maxCount {
			maxCount = count
			mode = []T{v}
		} else if count == maxCount {
			mode = append(mode, v)
		}
	}

	return mode
}

func bounds(i int, optionalJ ...int) (int, int) {
	var j int

	if len(optionalJ) > 0 {
		j = optionalJ[0]
	} else {
		j = i
	}

	return i, j
}

func Contains[V any](slice []V, matcher Matcher) bool {
	for i, v := range slice {
		if matcher(i, v) {
			return true
		}
	}

	return false
}
