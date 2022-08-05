package numeric

import (
	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/kv/ordered"
)

// Collection is an ordered collection which ensures all values are numbers. This
// allows for methods that are only applicable to numbers, such as sums and so on.
type Collection[K comparable, V collections.Number] struct {
	ordered.Collection[K, V]
}

// Collect returns an ordered numeric collection containing the given values
func Collect[V collections.Number](n ...V) Collection[int, V] {
	return Collection[int, V]{ordered.Collect(n...)}
}

// Average calculates the average value stored on the collection. Should the slice be empty,
// 0 is returned.
func (c Collection[K, V]) Average() V {
	return collections.Average(c.ToSlice())
}

// Sum sums all the values stored on the collection and returns the result.
func (c Collection[K, V]) Sum() V {
	return collections.Sum(c.ToSlice())
}

// Min returns the minimal value stored on collection. Should the slice be
// empty, 0 is returned.
func (c Collection[K, V]) Min() V {
	return collections.Min(c.ToSlice())
}

// Max returns the maximum value stored on the collection. Should the slice be
// empty, 0 is returned.
func (c Collection[K, V]) Max() V {
	return collections.Max(c.ToSlice())
}

// Median calculates and returns the median value of the collection.
func (c Collection[K, V]) Median() float64 {
	return collections.Median(c.ToSlice())
}

// Duplicates returns duplicate values from the collection
func (c Collection[K, V]) Duplicates() []V {
	return collections.Duplicates(c.ToSlice())
}
