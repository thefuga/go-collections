package collections

import (
	"math"
	"math/rand"
	"reflect"
	"sort"

	"github.com/thefuga/go-collections/errors"
	"github.com/thefuga/go-collections/internal"
)

// Get calls GetE, omitting the error.
func Get[T any](slice []T, i int) T {
	v, _ := GetE(slice, i)
	return v
}

// GetE indexes the slice with i, returning the corresponding value when it exists.
// Should i be negative or greater then the slice's len, a zeroed T value and
// an errors.ValueNotFound error is returned.
// This function is safe to be used with empty slices.
func GetE[T any](slice []T, i int) (T, error) {
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
	return GetE(slice, 0)
}

// Push appends v to the slice.
func Push[T any](slice []T, v T) []T {
	return append(slice, v)
}

// Put sets the position i in the slice with v. Should i be greater than the slice,
// a new slice with capacity to store v at i index is allocated.
// Put preserves the original values, shifting the slice so the new  value can be
// stored without affecting the previous values.
func Put[T any](slice []T, i int, v T) []T {
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
	v, _ := GetE(slice, len(slice)-1)
	return v
}

// LastE simply calls GetE with the slice and len-1 as the target index.
func LastE[T any](slice []T) (T, error) {
	return GetE(slice, len(slice)-1)
}

// LastBy uses LastByE, omitting the error.
func LastBy[T any](slice []T, matcher Matcher[int, T]) T {
	v, _ := LastByE(slice, matcher)
	return v
}

// LastByE returns the last matched element in the slice.
func LastByE[T any](slice []T, matcher Matcher[int, T]) (T, error) {
	for i := len(slice) - 1; i >= 0; i-- {
		if matcher(i, slice[i]) {
			return slice[i], nil
		}
	}

	return *new(T), errors.NewValueNotFoundError()
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
	i, _ := SearchE(slice, v)
	return i
}

// SearchE searches for v in the slice. The evaluation is done using reflect.DeepEqual.
// Should the value be found, it's index is returned. Otherwise, T's zeroed value and
// an instance of errors.ValueNotFoundError is returned.
func SearchE[T any](slice []T, v T) (int, error) {
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
func Map[T any, R any](slice []T, f func(i int, v T) R) []R {
	mappedValues := make([]R, 0, len(slice))

	Each(func(i int, v T) {
		mappedValues = Push(mappedValues, f(i, v))
	}, slice)

	return mappedValues
}

// Reduce reduces the collection to a single value, passing the result of each
// iteration into the subsequent iteration
func Reduce[T, V any](slice []T, f func(carry V, v T, i int) V, carry V) V {
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

// SortBy sorts `slice` based on `f`.
func SortBy[T any, S internal.Relational](slice []T, f func(t T) S) []T {
	sort.Slice(slice, func(i, j int) bool {
		return f(slice[i]) < f(slice[j])
	})
	return slice
}

// SortByDesc sorts desc `slice` based on f.
func SortByDesc[T any, S internal.Relational](slice []T, f func(t T) S) []T {
	sort.Slice(slice, func(i, j int) bool {
		return f(slice[i]) > f(slice[j])
	})
	return slice
}

// Copy returns a copy of the input slice.
func Copy[V any](slice []V) []V {
	copied := make([]V, len(slice))
	copy(copied, slice)
	return copied
}

// Cut uses CutE, omitting the error.
func Cut[V any](slice *[]V, i int, optionalJ ...int) []V {
	cut, _ := CutE(slice, i, optionalJ...)
	return cut
}

// CutE removes and returns the portion of the slice limited by i (included) and j (not included).
// Should either i or j be out of bounds, an instance of errors.IndexOutOfBounds is returned.
func CutE[V any](slice *[]V, i int, optionalJ ...int) ([]V, error) {
	sliceLen := len(*slice)
	i, j := bounds(i, optionalJ...)
	if i > sliceLen || j > sliceLen {
		return nil, errors.NewIndexOutOfBoundsError()
	}

	cut := make([]V, j-i)
	copy(cut, (*slice)[i:])

	copy((*slice)[i:], (*slice)[j:])
	for k, n := sliceLen-j+i, sliceLen; k < n; k++ {
		(*slice)[k] = *new(V)
	}
	*slice = (*slice)[:sliceLen-j+i]

	return cut, nil
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

// Tally counts the occurrence of each element on the slice.
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
func Contains[V any](slice []V, matcher AnyMatcher) bool {
	for i, v := range slice {
		if matcher(i, v) {
			return true
		}
	}

	return false
}

// FirstWhere uses FirstWhereE, omitting the error.
func FirstWhere[V any](slice []V, matcher AnyMatcher) V {
	v, _ := FirstWhereE(slice, matcher)
	return v
}

// FirstWhereE returns the first value matched by the given matcher. Should no value
// match, an instance of errors.ValueNotFoundError is returned.
func FirstWhereE[V any](slice []V, matcher AnyMatcher) (V, error) {
	for i, v := range slice {
		if matcher(i, v) {
			return v, nil
		}
	}

	return *new(V), errors.NewValueNotFoundError()
}

// FirstWhereField uses FirstWhereFieldE, omitting the error.
func FirstWhereField[V any](slice []V, field string, matcher AnyMatcher) V {
	v, _ := FirstWhereFieldE(slice, field, matcher)
	return v
}

// FirstWhereFieldE uses FieldMatcher to match a struct field from elements present
// on S.
// Should either no element match, the field doesn't exist on the struct V, or V is not
// a struct, an instance of errors.ValueNotFoundError is returned.
func FirstWhereFieldE[V any](slice []V, field string, matcher AnyMatcher) (V, error) {
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

// Diff returns a slice containing the elements that appear in the Left slice but not in the Right slice.
func Diff[V comparable](leftSlice, rightSlice []V) []V {
	seen := makeSeenMap(rightSlice)
	diff := []V{}

	for _, v := range leftSlice {
		if _, ok := seen[v]; !ok {
			diff = append(diff, v)
		}

	}

	return diff
}

// Intersect creates a new slice containing the elements present in both left
// and right slices. The given slices are left untoutched.
func Intersect[V comparable](leftSlice, rightSlice []V) []V {
	if len(rightSlice) > len(leftSlice) {
		return intersect(rightSlice, leftSlice)
	}

	return intersect(leftSlice, rightSlice)
}

func intersect[V comparable](leftSlice, rightSlice []V) []V {
	seen := makeSeenMap(rightSlice)
	intersection := []V{}

	for _, v := range leftSlice {
		if _, ok := seen[v]; ok {
			intersection = append(intersection, v)
		}

	}

	return intersection
}

func makeSeenMap[V comparable](slice []V) map[V]struct{} {
	seen := make(map[V]struct{}, len(slice))
	for _, v := range slice {
		seen[v] = struct{}{}
	}

	return seen
}

// Zip merges the values of the given slices at their corresponding indexes
func Zip[V any](slices ...[]V) [][]V {
	if len(slices) == 0 {
		return [][]V{}
	}

	minLen := len(slices[0])
	for _, s := range slices {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	result := make([][]V, minLen)

	for i := range result {
		result[i] = make([]V, len(slices))
		for j := range result[i] {
			result[i][j] = slices[j][i]
		}
	}

	return result
}

// Unique returns all distinct items in the slice
func Unique[V comparable](slice []V) []V {
	unique := []V{}

	seen := map[V]struct{}{}
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			unique = append(unique, v)
			seen[v] = struct{}{}
		}
	}

	return unique
}

// Unique uses the returned value of `f` to return distinct values in `slice`
func UniqueBy[V any, T comparable](slice []V, f func(v V) T) []V {
	unique := []V{}

	seen := map[T]struct{}{}
	for _, v := range slice {
		t := f(v)
		if _, ok := seen[t]; !ok {
			unique = append(unique, v)
			seen[t] = struct{}{}
		}
	}

	return unique
}

// GroupBy groups the slice's items by the return value of `f`
func GroupBy[V any, T comparable](slice []V, f func(v V) T) map[T][]V {
	result := map[T][]V{}

	for _, v := range slice {
		t := f(v)
		if group, ok := result[t]; ok {
			result[t] = append(group, v)
		} else {
			result[t] = []V{v}
		}
	}

	return result
}

// Partition divides the slice into two slices based on the given predicate function.
// It returns a slice of elements that satisfy the predicate and a slice of elements that do not.
func Partition[V any](slice []V, predicate func(v V) bool) ([]V, []V) {
	pass, reject := []V{}, []V{}
	for _, v := range slice {
		if predicate(v) {
			pass = append(pass, v)
		} else {
			reject = append(reject, v)
		}
	}
	return pass, reject
}

// Chunk breaks the slice into multiple, smaller slices of a given size
func Chunk[V any](slice []V, size int) [][]V {
	resultLen := int(math.Ceil(float64(len(slice)) / float64(size)))
	result := make([][]V, resultLen)

	for i := range result {
		result[i] = make([]V, 0, size)
	}

	for i, v := range slice {
		result[i/size] = append(result[i/size], v)
	}

	return result
}

// Reverse reverses the given slice
func Reverse[V any](slice []V) []V {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// SumBy accumulates values returned by `f`
func SumBy[V any, T internal.Number](slice []V, f func(v V) T) T {
	var sum T
	for _, v := range slice {
		sum += f(v)
	}
	return sum
}

// Range returns a slice containing integers in the specified range (i.e. [min, max])
func Range[T internal.Integer](min, max T) []T {
	result := make([]T, max-min+1)
	for i := range result {
		result[i] = T(i) + min
	}
	return result
}

// Interpose adds `sep` between every element in `slice`
func Interpose[V any](slice []V, sep V) []V {
	if len(slice) == 0 {
		return slice
	}

	resultLen := len(slice)*2 - 1
	result := make([]V, resultLen)
	result[resultLen-1] = slice[len(slice)-1]

	for i := 0; i+1 < len(slice); i++ {
		result[i*2] = slice[i]
		result[i*2+1] = sep
	}

	return result
}

// ForPage returns a slice containing the items that would be present on a given page number
func ForPage[V any](slice []V, page, size int) []V {
	lower := internal.Max((page-1)*size, 0)
	upper := internal.Min(lower+size, len(slice))
	return slice[lower:upper]
}

// KeyBy keys the collection by the given key
// If multiple items have the same key, the last one will appear in the new collection
func KeyBy[V any, K comparable](slice []V, f func(v V) K) map[K]V {
	result := map[K]V{}
	for _, v := range slice {
		result[f(v)] = v
	}
	return result
}

// PadRight will fill the slice with `pad` to the specified `size`.
// No padding will take place if the `size` is less than or equal to the length of `slice`
func PadRight[V any](slice []V, size int, pad V) []V {
	if size <= len(slice) {
		return slice
	}

	result := make([]V, size)
	copy(result, slice)

	for i := len(slice); i < size; i++ {
		result[i] = pad
	}
	return result
}

// PadLeft will fill the slice with `pad` from the left to the specified `size`.
// No padding will take place if the `size` is less than or equal to the length of `slice`
func PadLeft[V any](slice []V, size int, pad V) []V {
	if size <= len(slice) {
		return slice
	}

	offset := size - len(slice)

	result := make([]V, size)
	copy(result[offset:], slice)

	for i := 0; i < offset; i++ {
		result[i] = pad
	}
	return result
}

// Pad will fill the slice with `pad` to the specified `size`.
// To pad to the left, specify a negative `size`.
// No padding will take place if the absolute value of `size` is less than or equal to the length of `slice`
func Pad[V any](slice []V, size int, pad V) []V {
	if size >= 0 {
		return PadRight(slice, size, pad)
	}

	return PadLeft(slice, -size, pad)
}

// Prepend adds `value` to the beginning of `slice`
func Prepend[V any](slice []V, value V) []V {
	return append([]V{value}, slice...)
}

// Random uses RandomE, omitting the error.
func Random[V any](slice []V) V {
	v, _ := RandomE(slice)
	return v
}

// RandomE returns a random item from `slice` and errors if `slice` is empty.
func RandomE[V any](slice []V) (V, error) {
	if len(slice) == 0 {
		return *new(V), errors.NewEmptyCollectionError()
	}

	return slice[rand.Intn(len(slice))], nil
}

// Shuffle pseudo-randomizes the order of elements.
func Shuffle[V any](slice []V) []V {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
	return slice
}

// Skip returns `slice` with `skip` elements removed from the beginning
func Skip[V any](slice []V, skip int) []V {
	return slice[internal.Min(skip, len(slice)):]
}

// SkipUntil skips over items from `slice` until `matcher` returns true and
// then returns the remaining items in the slice
func SkipUntil[V any](slice []V, matcher AnyMatcher) []V {
	for i, v := range slice {
		if matcher(i, v) {
			return slice[i:]
		}
	}

	return slice[len(slice):]
}

// SkipWhile skips over items from `slice` while `matcher` returns true and
// then returns the remaining items in the slice
func SkipWhile[V any](slice []V, matcher AnyMatcher) []V {
	return SkipUntil(slice, Not(matcher))
}

// Nth creates a new slice consisting of every n-th element, starting at 0.
func Nth[V any, N internal.Integer](slice []V, n N) []V {
	return NthOffset(slice, n, 0)
}

// NthOffset creates a new slice consisting of every n-th element, starting at the given offset.
func NthOffset[V any, N internal.Integer](slice []V, n N, off N) []V {
	nthLen := N((len(slice) / int(n)))
	nthSlice := make([]V, 0, nthLen)

	for nth := off; nth < N(len(slice)); nth = nth + n {
		nthSlice = append(nthSlice, slice[nth])
	}

	return nthSlice
}

// Sliding returns a "sliding window" view of the items in `slice`. Each window
// will by `step` items apart
func SlidingStep[V any](slice []V, window, step int) [][]V {
	if step < 1 || window < 1 || len(slice) == 0 {
		return nil
	}

	if window >= len(slice) {
		return [][]V{slice}
	}

	result := make([][]V, 0, len(slice)-window+1)
	for i := 0; i <= len(slice)-window; i += step {
		result = append(result, slice[i:i+window])
	}
	return result
}

// Sliding returns a "sliding window" view of the items in `slice`
func Sliding[V any](slice []V, window int) [][]V {
	return SlidingStep(slice, window, 1)
}

// Splice returns a slice of items starting at the specified index,
// and the updated slice with the items removed.
func Splice[V any](slice []V, idx int) ([]V, []V) {
	if idx >= len(slice) {
		return nil, nil
	}
	return slice[idx:], slice[:idx]
}

// SpliceN returns a slice of `slice` starting at the `index` with length `size`,
// and the updated slice with the items removed.
func SpliceN[V any](slice []V, idx, size int) ([]V, []V) {
	if idx < 0 || size < 1 {
		return nil, nil
	}

	if idx >= len(slice) {
		return nil, slice
	}

	if idx+size > len(slice) {
		size = len(slice) - idx
	}

	return Copy(slice[idx : idx+size]), append(slice[:idx], slice[idx+size:]...)
}

// Split breaks `slice` into the given number of groups
func Split[V any](slice []V, numberOfGroups int) [][]V {
	if len(slice) == 0 {
		return nil
	}

	numberOfGroups = internal.Min(len(slice), numberOfGroups)
	groupSize := internal.DivFloor(len(slice), numberOfGroups)
	remain := len(slice) % numberOfGroups

	result := make([][]V, numberOfGroups)
	start := 0
	for i := 0; i < numberOfGroups; i++ {
		size := groupSize
		if i < remain {
			size++
		}
		result[i] = slice[start:internal.Min(start+size, len(slice))]
		start += size
	}
	return result
}

// Take returns a slice with the specified number of items from `slice`.
// You may also pass a negative integer to take the specified number of items from the end of the `slice`.
func Take[V any](slice []V, n int) []V {
	if n < 0 {
		i := len(slice) - internal.Min(-n, len(slice))
		return slice[i:]
	}

	return slice[:internal.Min(n, len(slice))]
}

// TakeWhile returns items in the `slice` until `matcher` returns false
func TakeWhile[V any](slice []V, matcher AnyMatcher) []V {
	for i, v := range slice {
		if !matcher(i, v) {
			return slice[:i]
		}
	}

	return slice
}

// TakeUntil returns items in the `slice` until `matcher` returns true
func TakeUntil[V any](slice []V, matcher AnyMatcher) []V {
	return TakeWhile(slice, Not(matcher))
}

// Times creates a new slice by invoking `f` `n` times:
func Times[V any](n int, f func(i int) V) []V {
	result := make([]V, n)
	for i := range result {
		result[i] = f(i + 1)
	}
	return result
}
