// slice package provides a custom slice collection type, functions and methods
// related to slices.
package slice

import (
	"github.com/thefuga/go-collections"
)

// Collection is a custom generic slice type, specially useful to pipe generic slice methods.
// Most methods of Collection will be delegated to the generics functions using the
// receiver as the argument slice.
type Collection[V any] []V

// Collect a new Collection from the given values. It is essentially a syntactic sugar
// []V{}.
func Collect[V any](values ...V) Collection[V] {
	return values
}

// Contains passes the collection and the given params to the generic Contains function.
func (c Collection[V]) Contains(matcher collections.AnyMatcher) bool {
	return collections.Contains(c, matcher)
}

// Get passes the collection and the given params to the generic Get function.
func (c Collection[V]) Get(i int) V { return collections.Get(i, c) }

// GetE passes the collection and the given params to the generic GetE function.
func (c Collection[V]) GetE(i int) (V, error) { return collections.GetE(i, c) }

// Push passes the collection and the given params to the generic Push function.
func (c Collection[V]) Push(v V) Collection[V] { return append(c, v) }

// Put passes the collection and the given params to the generic Put function.
func (c Collection[V]) Put(i int, v V) Collection[V] { return collections.Put(i, v, c) }

// Pop passes the collection and the given params to the generic Pop function.
func (c *Collection[V]) Pop() V { return collections.Pop((*[]V)(c)) }

// PopE passes the collection and the given params to the generic PopE function.
func (c *Collection[V]) PopE() (V, error) { return collections.PopE((*[]V)(c)) }

// Count returns the length of the collection.
func (c Collection[V]) Count() int { return len(c) }

// Capacity returns the capacity of the collection.
func (c Collection[V]) Capacity() int { return cap(c) }

// IsEmpty checks if the collection is empty.
func (c Collection[V]) IsEmpty() bool { return c.Count() == 0 }

// Search passes the collection and the given params to the generic Search function.
func (c Collection[V]) Search(v V) int { return collections.Search(v, c) }

// SearchE passes the collection and the given params to the generic SearchE function.
func (c Collection[V]) SearchE(v V) (int, error) { return collections.SearchE(v, c) }

// Map passes the collection and the given params to the generic Map function.
func (c Collection[V]) Map(f func(i int, v V) V) Collection[V] { return collections.Map(f, c) }

// First passes the collection and the given params to the generic First function.
func (c Collection[V]) First() V { return collections.First(c) }

// FirstE passes the collection and the given params to the generic FirstE function.
func (c Collection[V]) FirstE() (V, error) { return collections.FirstE(c) }

// Last passes the collection and the given params to the generic Last function.
func (c Collection[V]) Last() V { return collections.Last(c) }

// LastE passes the collection and the given params to the generic LastE function.
func (c Collection[V]) LastE() (V, error) { return collections.LastE(c) }

// Each passes the collection and the given params to the generic Each function and
// returns the collection.
func (c Collection[V]) Each(f func(i int, v V)) Collection[V] {
	collections.Each(f, c)
	return c
}

// Sort passes the collection and the given params to the generic Sort function and
// returns the collection.
func (c Collection[V]) Sort(f func(current, next V) bool) Collection[V] {
	collections.Sort(c, f)
	return c
}

// Tap passes the collection to f and returns the collection.
func (c Collection[V]) Tap(f func(Collection[V])) Collection[V] {
	f(c)
	return c
}

// ForgetE passes the collection and the given params to the generic ForgetE function.
func (c *Collection[V]) ForgetE(i int) error {
	return collections.ForgetE((*[]V)(c), i)
}

// ToSlice simply returns the collection as a slice.
func (c Collection[V]) ToSlice() []V { return c }
