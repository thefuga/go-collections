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

	if cap(slice) < len(slice)+1 {
		newSlice := make([]T, 0, len(slice)+1)
		slice = append(newSlice, slice...)
	}

	slice = append(slice[:i+1], slice[i:]...)
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

	Each(func(_ int, v T) {
		mappedValues = Push(v, mappedValues)
	}, slice)

	return mappedValues
}

func Sort[T any](slice []T, f func(current, next T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return f(slice[i], slice[j])
	})
}


	}

}
