// kv package provides a custom map collection type, functions and methods
// related to maps.
// The types and methods from kv don't guarantee order. See kv/ordered for that.
package kv

import (
	"reflect"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/errors"
	"github.com/thefuga/go-collections/slice"
)

// Collection is a custom generic map type specially useful to pipe generric methos.
type Collection[K comparable, V any] map[K]V

// Collect returns the result of CollectSlice passing the given items.
func Collect[T any](items ...T) Collection[int, T] {
	return CollectSlice(items)
}

// CollectSlice makes a map[int]T, given each item a numeric sequencial key corresponding
// to the item index.
func CollectSlice[T any](items []T) Collection[int, T] {
	collection := make(Collection[int, T], len(items))

	for key, item := range items {
		collection.Put(key, item)
	}

	return collection
}

// CollectMap returns the given map as a Collection.
func CollectMap[K comparable, V any](items map[K]V) Collection[K, V] {
	return items
}

// Combine calls CombineE, omitting the error.
func Combine[K comparable, V any](
	keys slice.Collection[K], values slice.Collection[V],
) Collection[K, V] {
	c, _ := CombineE(keys, values)
	return c
}

// CombineE uses the first slice as keys and the second as values to build a map.
// The order is preserved.
// Should the lenght of keys and values slices be different, an instance of errors.KeysValuesLegthMismatch
// is returned.
func CombineE[K comparable, V any](
	keys slice.Collection[K], values slice.Collection[V],
) (Collection[K, V], error) {
	if keys.Count() != values.Count() {
		return nil, errors.NewKeysValuesLengthMismatch()
	}

	count := keys.Count()

	collection := make(Collection[K, V], keys.Count())

	for i := 0; i < count; i++ {
		collection.Put(keys.Get(i), values.Get(i))
	}

	return collection, nil
}

// TODO
func CountBy[T comparable, K comparable, V any](c Collection[K, V], f func(v V) T) map[T]int {
	count := map[T]int{}

	c.Each(func(_ K, v V) {
		count[f(v)]++
	})

	return count
}

// Each ia a typical for loop. The current key and values are passed to the closure
// on each iteration. Order is not guaranteed on each execution.
func (c Collection[K, V]) Each(f func(k K, v V)) Collection[K, V] {
	for key, value := range c {
		f(key, value)
	}

	return c
}

// Get calls GetE, omitting the error.
func (c *Collection[K, V]) Get(k K) V {
	v, _ := c.GetE(k)
	return v
}

// GetE indexes the map with k, returning the corresponding value when it exists.
// Should k not exist, a zeroed T value and an instance of errors.KeyNotFound
// error is returned.
// This function is safe to be used with empty maps.
func (c Collection[K, V]) GetE(k K) (V, error) {
	item, found := c[k]

	if !found {
		return *new(V), errors.NewKeyNotFoundError(k)
	}

	return item, nil
}

// Search uses SearchE, omitting the error.
func (c Collection[K, V]) Search(v V) K {
	k, _ := c.SearchE(v)
	return k
}

// SearchE searches for v in the collection. The evaluation is done using reflect.DeepEqual.
// Should the value be found, it's key is returned. Otherwise, T's zeroed value and
// an instance of errors.ValueNotFoundError is returned. Considering the lack of ordering
// on maps, should multiple matching values exist, multiple calls to SearchE may return
// multiple different keys.
func (c Collection[K, V]) SearchE(v V) (K, error) {
	for k, collectionV := range c {
		if reflect.DeepEqual(collectionV, v) {
			return k, nil
		}
	}

	return *new(K), errors.NewValueNotFoundError()
}

// Map applies f to each element of the map and builds a new map with f's returned
// value. The built map is returned.
// Order is not guaranteed.
func (c Collection[K, V]) Map(f func(k K, v V) V) Collection[K, V] {
	mapped := make(Collection[K, V], c.Count())

	c.Each(func(k K, v V) {
		mapped.Put(k, f(k, v))
	})

	return mapped
}

// Count returns the number of elements stored on the collection
func (c Collection[K, V]) Count() int {
	return len(c)
}

// Put inserts v in the key represented by k. If k already exists on the map, it's
// value is overridden.
func (c Collection[K, V]) Put(k K, v V) Collection[K, V] {
	c[k] = v

	return c
}

// IsEmpty checks if the collection is empty.
func (c Collection[K, V]) IsEmpty() bool { return len(c) == 0 }

// Keys returns a new slice.Collection containing all of the maps keys. Ordering is not
// guaranteed.
func (c Collection[K, V]) Keys() slice.Collection[K] {
	keys := make(slice.Collection[K], 0, c.Count())

	c.Each(func(k K, _ V) {
		keys = keys.Push(k)
	})

	return keys
}

// KeysValues returns two slices containing all of the keys and values of the map.
// Indexes from both slice will correspond to the former key-value pair (i.e. calling
// Combine(collection.KeysValues()) returns a new collection equivalent to the one
// that originated the new collection).
func (c Collection[K, V]) KeysValues() (slice.Collection[K], slice.Collection[V]) {
	keys := make(slice.Collection[K], 0, c.Count())
	values := make(slice.Collection[V], 0, c.Count())

	c.Each(func(k K, v V) {
		keys = keys.Push(k)
		values = values.Push(v)
	})

	return keys, values
}

// Only returns a new collection containing only the key-value pairs from the keys slice.
func (c Collection[K, V]) Only(keys []K) Collection[K, V] {
	onlyValues := make(map[K]V, len(keys))

	for _, key := range keys {
		value, err := c.GetE(key)
		if err == nil {
			onlyValues[key] = value
		}
	}

	return onlyValues
}

// Tap passes the collection to f and returns the collection.
func (c Collection[K, V]) Tap(f func(Collection[K, V])) Collection[K, V] {
	f(c)
	return c
}

// Values returns a new slice.Collection containing all of the map's values. Ordering is not
// guaranteed.
func (c Collection[K, V]) Values() slice.Collection[V] {
	values := make(slice.Collection[V], 0, c.Count())

	c.Each(func(_ K, v V) {
		values = values.Push(v)
	})

	return values
}

// ToSlice is an alias to Collection.Values.
func (c Collection[K, V]) ToSlice() []V { return c.Values() }

// Copy returns a new collection equivalent to the copied.
func (c Collection[K, V]) Copy() Collection[K, V] {
	return Combine(c.KeysValues())
}

// Concat puts the key-value pairs from concatTo on the original collection. Should two
// keys equal, the original value will be preserved. To override the values, see Collection.Merge.
func (c Collection[K, V]) Concat(concatTo Collection[K, V]) Collection[K, V] {
	concatTo.Each(func(k K, v V) {
		if _, ok := c[k]; ok {
			return
		}
		c.Put(k, v)
	})

	return c
}

// Contains checks if any values on the Collection match f.
func (c Collection[K, V]) Contains(f collections.Matcher) bool {
	return c.Values().Contains(f)
}

// Every checks if every value on the Collection match f.
func (c Collection[K, V]) Every(f collections.Matcher) bool {
	for k, v := range c {
		if !f(k, v) {
			return false
		}
	}

	return true
}

// Flip calls Collection.FlipE, omitting the error.
func (c Collection[K, V]) Flip() Collection[K, V] {
	v, _ := c.FlipE()
	return v
}

// Flip makes a new collection, flipping the values with the keys.
func (c Collection[K, V]) FlipE() (Collection[K, V], error) {
	if _, err := collections.AssertE[V](*new(K)); err != nil {
		return c, err

	}

	flippedValues := make(map[K]V, c.Count())

	c.Each(func(k K, v V) {
		castKey, keyOk := collections.Assert[K](v)
		castValue, valueOk := collections.Assert[V](k)

		if keyOk && valueOk {
			flippedValues[castKey] = castValue
		} else {
			flippedValues[k] = v
		}
	})

	return flippedValues, nil
}

// Merge works similarly to Concat, but overrides conflicting keys.
func (c Collection[K, V]) Merge(other Collection[K, V]) Collection[K, V] {
	other.Each(func(k K, v V) {
		c.Put(k, v)
	})

	return c
}

// Filter makes a new collection containing only the key-value pairs matched by f.
func (c Collection[K, V]) Filter(f func(k K, v V) bool) Collection[K, V] {
	filtered := make(map[K]V)

	c.Each(func(k K, v V) {
		if f(k, v) {
			filtered[k] = c[k]
		}
	})

	return filtered
}

// Reject makes a new collection containing only the kill value pairs not matched
// by f.
func (c Collection[K, V]) Reject(f func(k K, v V) bool) Collection[K, V] {
	return c.Filter(func(k K, v V) bool {
		return !f(k, v)
	})
}

// Forget deletes the given key and returns the collection.
func (c Collection[K, V]) Forget(k K) Collection[K, V] {
	delete(c, k)
	return c
}

// ForgetE checks if the key exist before attempting to delete it. Should the key
// not exist, an instance of errors.KeyNotFoundError is returned. The original collection
// is always returned.
func (c Collection[K, V]) ForgetE(k K) (Collection[K, V], error) {
	if _, ok := c[k]; !ok {
		return c, errors.NewKeyNotFoundError(k)
	}

	delete(c, k)
	return c, nil
}

// When calls f with the collection when execute is true. Usually, execute will be
// the result of a function call.
// E.g.: c.When(!c.IsEmpty(), func(c Collection){ fmt.Printf("%v", c)}) prints the collection
// when the collection is not empty.
func (c Collection[K, V]) When(
	execute bool, f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	if !execute {
		return c
	}

	return f(c)
}

// WhenEmpty calls f with the collection when the collection is empty.
func (c Collection[K, V]) WhenEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.When(c.IsEmpty(), f)
}

// WhenNotEmpty calls f with the collection when the collection is not empty.
func (c Collection[K, V]) WhenNotEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.When(!c.IsEmpty(), f)
}

// Unless calls f with the collection when execute is false. Usually, execute will be
// the result of a function call.
// E.g.: c.Unless(c.IsEmpty(), func(c Collection){ fmt.Printf("%v", c)}) prints the collection
// when the collection is not empty.
func (c Collection[K, V]) Unless(
	execute bool, f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.When(!execute, f)
}

// UnlessEmpty calls f with the collection when the collection is not empty.
func (c Collection[K, V]) UnlessEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.WhenNotEmpty(f)
}

// UnlessNotEmpty calls f with the collection when the collection is empty.
func (c Collection[K, V]) UnlessNotEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.WhenEmpty(f)
}
