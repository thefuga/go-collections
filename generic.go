package collections

import (
	"reflect"
	"sort"

	"github.com/thefuga/go-collections/errors"
)

// Get calls GetE, omitting the error.
func Get[T any](i int, slice []T) T {
	v, _ := GetE(i, slice)
	return v
}

// GetE indexes the slice with i, returning the corresponding value when it exists.
// Should i be negative or greater then the slice's len, a zeroed T value and
// an errors.ValueNotFound error is returned.
// This function is safe to be used with empty slices.
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

// First calls FirstE, omitting the error.
func First[T any](slice []T) T {
	v, _ := FirstE(slice)
	return v
}

// FirstE simply calls GetE with the slice and 0 the target index.
func FirstE[T any](slice []T) (T, error) {
	return GetE(0, slice)
}

// Push appends v to the slice.
func Push[T any](v T, slice []T) []T {
	return append(slice, v)
}

// Put sets the position i in the slice with v. Should i be greater than the slice,
// a new slice with capacity to store v at i index is allocated.
// Put preserves the original values, shifting the slice so the new  value can be
// stored without affecting the previous values.
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

// Pop uses PopE, omitting the error.
func Pop[T any](slice *[]T) T {
	v, _ := PopE(slice)
	return v
}

// PopE takes the last element of the slice, deletes and returns it.
func PopE[T any](slice *[]T) (T, error) {
	v, err := LastE(*slice)

	if err != nil {
		return v, err
	}

	*slice = (*slice)[:len(*slice)-1]

	return v, nil
}

// Shift uses ShiftE, omitting the error.
func Shift[T any](slice *[]T) T {
	v, _ := ShiftE(slice)
	return v
}

// ShiftE (aka pop front) is equivalent to PopE, but deletes and returns the first
// slice element.
func ShiftE[T any](slice *[]T) (T, error) {
	v, err := FirstE(*slice)

	if err != nil {
		return v, err
	}

	*slice = (*slice)[1:]

	return v, nil
}

// Last uses LastE, omitting the error.
func Last[T any](slice []T) T {
	v, _ := GetE(len(slice)-1, slice)
	return v
}

// LastE simply calls GetE with the slice and len-1 as the target index.
func LastE[T any](slice []T) (T, error) {
	return GetE(len(slice)-1, slice)
}

// Each ia a typical for loop. The current index and values are passed to the closure
// on each iteration.
func Each[T any](f func(i int, v T), slice []T) {
	for i, v := range slice {
		f(i, v)
	}
}

// Search uses SearchE, omitting the error.
func Search[T any](v T, slice []T) int {
	i, _ := SearchE(v, slice)
	return i
}

// SearchE searches for v in the slice. The evaluation is done using reflect.DeepEqual.
// Should the value be found, it's index is returned. Otherwise, T's zeroed value and
// an instance of errors.ValueNotFoundError is returned.
func SearchE[T any](v T, slice []T) (int, error) {
	for i := range slice {
		if reflect.DeepEqual(slice[i], v) {
			return i, nil
		}
	}

	return -1, errors.NewValueNotFoundError()
}

// Map applies f to each element of the slice and builds a new slice with f's returned
// value. The built slice is returned.
// The mapped slice has the same order as the input slice.
func Map[T any](f func(i int, v T) T, slice []T) []T {
	mappedValues := make([]T, 0, len(slice))

	Each(func(i int, v T) {
		mappedValues = Push(f(i, v), mappedValues)
	}, slice)

	return mappedValues
}

// Reduce reduces the collection to a single value, passing the result of each
// iteration into the subsequent iteration
func Reduce[T, V any](f func(carry V, v T, i int) V, carry V, slice []T) V {
	Each(func(i int, v T) { carry = f(carry, v, i) }, slice)

	return carry
}

// Sort sorts the slice based on f. It can be used used with Asc or Desc functions
// or with a custom closure.
func Sort[T any](slice []T, f func(current, next T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return f(slice[i], slice[j])
	})
}

// Copy returns a copy of the input slice.
func Copy[V any](slice []V) []V {
	copied := make([]V, len(slice))
	copy(slice, copied)
	return copied
}

// Cut uses CutE, omitting the error.
func Cut[V any](slice *[]V, i int, optionalJ ...int) []V {
	cutted, _ := CutE(slice, i, optionalJ...)
	return cutted
}

// CutE removes and returns the portion of the slice limited by i (included) and j (not included).
// Should either i or j be out of bounds, an instance of errors.IndexOutOfBounds is returned.
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

// DeleteE deletes the element corresponding to i from the slice. Every element on the
// right of i will be re-indexed.
// Should either i be out of bounds, an instance of errors.IndexOutOfBounds is returned.
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

// ForgetE is an alias to DeleteE.
func ForgetE[V any](slice *[]V, i int, optionalJ ...int) error {
	return DeleteE(slice, i, optionalJ...)
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

// TODO
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

// Contains checks if the slice holds at least one value matching the given matcher.
func Contains[V any](slice []V, matcher Matcher) bool {
	for i, v := range slice {
		if matcher(i, v) {
			return true
		}
	}

	return false
}

// FirstWhere uses FirstWhereE, omitting the error.
func FirstWhere[V any](slice []V, matcher Matcher) V {
	v, _ := FirstWhereE(slice, matcher)
	return v
}

// FirstWhereE returns the first value matched by the given matcher. Should no value
// match, an instance of errors.ValueNotFoundError is returned.
func FirstWhereE[V any](slice []V, matcher Matcher) (V, error) {
	for i, v := range slice {
		if matcher(i, v) {
			return v, nil
		}
	}

	return *new(V), errors.NewValueNotFoundError()
}

// FirstWhereField uses FirstWhereFieldE, omitting the error.
func FirstWhereField[V any](slice []V, field string, matcher Matcher) V {
	v, _ := FirstWhereFieldE(slice, field, matcher)
	return v
}

// FirstWhereFieldE uses FieldMatcher to match a struct field from elements present
// on S.
// Should either no element match, the field doesn't exist on the struct V, or V is not
// a struct, an instance of errors.ValueNotFoundError is returned.
func FirstWhereFieldE[V any](slice []V, field string, matcher Matcher) (V, error) {
	for i, v := range slice {
		if FieldMatch[V](field, matcher)(i, v) {
			return v, nil
		}
	}

	return *new(V), errors.NewValueNotFoundError()
}

func Duplicates[V comparable](slice []V) []V {
	seen := make(map[V]uint8)
	duplicates := []V{}

	for _, n := range slice {
		switch seen[n] {
		case 0:
			seen[n] = 1
		case 1:
			seen[n] = 2
			duplicates = append(duplicates, n)
		}
	}

	return duplicates
}

func Diff[V comparable](slice1, slice2 []V) []V {
	difference := []V{}

	for _, v := range slice1 {
		if !Contains(slice2, func(_ any, v2 any) bool {
			return v == v2
		}) {
			difference = append(difference, v)
		}
	}

	return difference
}
