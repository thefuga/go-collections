package ordered

import (
	"reflect"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/errors"
	"github.com/thefuga/go-collections/internal"
	"github.com/thefuga/go-collections/kv"
	"github.com/thefuga/go-collections/slice"
)

// Collection is composed by a slice collection of keys and a kv collection of key-value
// pairs. This ensures the order of the values, allowing for deterministic execution of
// iterative methods - which is not possible when iterating over maps.
type Collection[K comparable, V any] struct {
	keys   slice.Collection[K]
	values kv.Collection[K, V]
}

// Collect returns the result of CollectSlice passing the given items.
func Collect[T any](items ...T) Collection[int, T] {
	return CollectSlice(items)
}

// CollectSlice makes a new collection, given each item a numeric sequential key corresponding
// to the item index.
func CollectSlice[T any](items []T) Collection[int, T] {
	collection := makeCollection[int, T](len(items))

	for key, item := range items {
		collection.Put(key, item)
	}

	return collection
}

// CollectMap returns the given map as a Collection. The order of the keys is
// randomized, but preserved for future method calls.
// Should you need an specific order, immediately call Sort on the returned collection
func CollectMap[K comparable, V any](items map[K]V) Collection[K, V] {
	collection := Collection[K, V]{
		keys:   make(slice.Collection[K], 0, len(items)),
		values: items,
	}

	for key := range items {
		collection.keys = append(collection.keys, key)
	}

	return collection
}

func makeCollection[K comparable, V any](capacity int) Collection[K, V] {
	return Collection[K, V]{
		keys:   make([]K, 0, capacity),
		values: make(map[K]V, capacity),
	}
}

// Get attempts to get the item corresponding to k in c. Should the key exist,
// It's value will be converted - when possible - to T. This is useful when working
// with collections where the type of the values are unknown (e.g. Collection[string, any]).
func Get[K comparable, T, V any](c Collection[K, V], k K) (T, error) {
	var (
		genericValue any
		getErr       error
	)

	genericValue, getErr = c.GetE(k)

	if getErr != nil {
		return *new(T), getErr
	}

	return internal.AssertE[T](genericValue)
}

// TODO
func CountBy[T comparable, K comparable, V any](c Collection[K, V], f func(v V) T) map[T]int {
	return kv.CountBy(c.values, f)

}

// Put inserts v in the key represented by k. If k already exists on the map, it
// value is overridden. The inserted item will be at the last position of the keys
// slice (i.e. it will be the last element on iterations, unless the collection is sorted).
func (c *Collection[K, V]) Put(k K, v V) Collection[K, V] {
	if _, ok := c.values[k]; !ok {
		c.keys = append(c.keys, k)
	}

	c.values[k] = v

	return *c
}

// Push appends v to the end of the collection. The value will be associated with
// a key corresponding to the count of the collection at the time of pushing.
func (c *Collection[K, V]) Push(v V) Collection[K, V] {
	if castK, ok := internal.Assert[K](c.Count()); ok {
		c.Put(castK, v)
	}

	return *c
}

// Pop uses PopE, omitting the error.
func (c *Collection[K, V]) Pop() V {
	v, _ := c.PopE()
	return v
}

// Pop takes the last element of the collection, deletes and returns it.
func (c *Collection[K, V]) PopE() (V, error) {
	if c.IsEmpty() {
		return *new(V), errors.NewEmptyCollectionError()
	}

	v := c.Last()

	lastKey := c.keys[len(c.keys)-1]
	c.keys = c.keys[:len(c.keys)-1]
	delete(c.values, lastKey)

	return v, nil
}

// IsEmpty checks if the collection is empty.
func (c Collection[K, V]) IsEmpty() bool { return c.keys.IsEmpty() }

// Get calls Get in the underlying values map.
func (c *Collection[K, V]) Get(k K) V { return c.values.Get(k) }

// GetE calls GetE in the underlying values map.
func (c *Collection[K, V]) GetE(k K) (V, error) { return c.values.GetE(k) }

// Count returns the number of elements stored on the collection
func (c Collection[K, V]) Count() int { return c.keys.Count() }

// Each iterates over the underlying keys slice, passing the k and corresponding
// value to f. The order in which the key-value pairs are passed to f is always the
// same considering the same collection (i.e. it executes deterministically).
func (c Collection[K, V]) Each(f func(k K, v V)) Collection[K, V] {
	c.keys.Each(func(_ int, k K) {
		f(k, c.values[k])
	})

	return c
}

// Tap passes the collection to f and returns the collection.
func (c Collection[K, V]) Tap(f func(Collection[K, V])) Collection[K, V] {
	f(c)
	return c
}

// Search uses SearchE, omitting the error.
func (c Collection[K, V]) Search(v V) K {
	k, _ := c.SearchE(v)
	return k
}

// SearchE searches for v in the collection. The evaluation is done using reflect.DeepEqual.
// Should the value be found, it's key is returned. Otherwise, T's zeroed value and
// an instance of errors.ValueNotFoundError is returned.
// SearchE iterates on the keys collection, which guarantees the order of executions (i.e. even if multiple values are present, the same key will be returned on multiple calls).
func (c Collection[K, V]) SearchE(v V) (K, error) {
	for _, k := range c.keys {
		if reflect.DeepEqual(c.values[k], v) {
			return k, nil
		}
	}

	return *new(K), errors.NewValueNotFoundError()
}

// Keys returns the keys stored in the collection. Order is guaranteed.
func (c Collection[K, V]) Keys() slice.Collection[K] { return c.keys }

// Sort sorts the collection keys and returns the ordered collection. The underlying map
// is not affected.
func (c Collection[K, V]) Sort(f func(current, next V) bool) Collection[K, V] {
	c.keys.Sort(func(i, j K) bool {
		return f(c.Get(i), c.Get(j))
	})

	return c
}

// Map applies f to each element of the map and builds a new map with f's returned
// value. The built map is returned.
// Map iterates over the underlying keys collection, which guarantees ordering.
func (c Collection[K, V]) Map(f func(k K, v V) V) Collection[K, V] {
	mappedValues := makeCollection[K, V](c.Count())

	c.Each(func(k K, v V) {
		mappedValues.Put(k, f(k, v))
	})

	return mappedValues
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

	return CollectMap(onlyValues)
}

// First calls FirstE, omitting the error.
func (c Collection[K, V]) First() V {
	first, _ := c.FirstE()
	return first
}

// FirstE returns the value associated to the first key in the collection. It returns
// an error should the collection be empty.
func (c Collection[K, V]) FirstE() (V, error) {
	first, err := c.keys.FirstE()

	if err != nil {
		return *new(V), err
	}

	return c.values[first], nil
}

// FirstOrFail returns the key-value pair of the first occurrence matched by f.
// Should no entry match, an instance of NewValueNotFoundError is returned.
func (c Collection[K, V]) FirstOrFail(f collections.AnyMatcher) (K, V, error) {
	var (
		found bool
		v     V
		k     K
	)

	// Optimize this to return on the first occurrence. Add more test scenarios!
	c.Each(func(atK K, atV V) {
		if !f(atK, atV) {
			return
		}

		found = true
		k = atK
		v = atV
	})

	if !found {
		return *new(K), *new(V), errors.NewValueNotFoundError()
	}

	return k, v, nil
}

// Last calls LastE, omitting the error.
func (c Collection[K, V]) Last() V { v, _ := c.LastE(); return v }

// LastE returns the value associated with the last element in the keys collection.
// Should the collection be empty, an error is returned.
func (c Collection[K, V]) LastE() (V, error) {
	last, err := c.keys.LastE()

	if err != nil {
		return *new(V), err
	}

	return c.values[last], nil
}

// ToSlice makes a new slice containing all values from the collection. The values
// are obtained through the keys collection, which guarantees ordering.
func (c Collection[K, V]) ToSlice() []V {
	slice := make([]V, len(c.keys))

	for i, key := range c.keys {
		slice[i] = c.values[key]
	}

	return slice
}

// ToSlice collection simply returns ToSlice as a slice.Collection type
func (c Collection[K, V]) ToSliceCollection() slice.Collection[V] {
	return c.ToSlice()
}

// Combine calls CombineE, omitting the error.
func (c Collection[K, V]) Combine(v Collection[K, V]) Collection[K, V] {
	combined, _ := c.CombineE(v)
	return combined
}

// CombineE makes a new collection using the receiver values as keys on the new collection,
// and the given collection values as values on the new collection.
// The values on the keys collection must be of the same type of it's keys, otherwise an
// error is returned.
func (c Collection[K, V]) CombineE(values Collection[K, V]) (Collection[K, V], error) {
	combined := makeCollection[K, V](c.Count())

	for i := 0; i < c.Count(); i++ {
		k, err := internal.AssertE[K](c.values[c.keys[i]])
		if err != nil {
			return combined, err
		}

		v, err := internal.AssertE[V](values.values[values.keys[i]])
		if err != nil {
			return combined, err
		}

		combined.Put(k, v)
	}

	return combined, nil
}

// Concat appends the given collection to the receiving collection. Should keys
// on concatTo already exist on the base collection, the value will be pushed (when possible)
// to the end of the collection. In case the value cannot be pushed, it is discarded.
// Concat ensures all keys and values on the base collection are preserved.
func (c Collection[K, V]) Concat(concatTo Collection[K, V]) Collection[K, V] {
	concatenated := CollectMap(c.values)

	concatTo.Each(func(k K, v V) {
		if _, ok := concatenated.values[k]; ok {
			concatenated.Push(v)
		} else {
			concatenated.Put(k, v)
		}
	})

	return concatenated
}

// Contains checks if any values on the Collection match f.
func (c Collection[K, V]) Contains(f collections.AnyMatcher) bool {
	_, _, err := c.FirstOrFail(f)
	return err == nil
}

// Every checks if every value on the Collection match f.
func (c Collection[K, V]) Every(f collections.AnyMatcher) bool {
	contains := true

	c.Each(func(k K, v V) {
		if !f(k, v) {
			contains = false
			return
		}
	})

	return contains
}

// Flip makes a new collection, flipping the values with the keys.
// The type of the keys and values must be the same (e.g. Collection[string, string]).
// In case the types are not compatible (e.g. Collection[string, struct{}]), the entries
// won't be flipped.
func (c Collection[K, V]) Flip() Collection[K, V] {
	flippedKeys := make([]K, 0, c.Count())
	flippedValues := make(map[K]V, c.Count())

	c.Each(func(k K, v V) {
		castKey, keyOk := internal.Assert[K](v)
		castValue, valueOk := internal.Assert[V](k)

		if keyOk && valueOk {
			flippedKeys = append(flippedKeys, castKey)
			flippedValues[castKey] = castValue
		} else {
			flippedKeys = append(flippedKeys, k)
			flippedValues[k] = v
		}
	})

	c.keys = flippedKeys
	c.values = flippedValues

	return c
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
	values := make(map[K]V)

	c.Each(func(k K, v V) {
		if f(k, c.values[k]) {
			values[k] = c.values[k]
		}
	})

	return CollectMap(values)
}

// Reject makes a new collection containing only the kill value pairs not matched
// by f.
func (c Collection[K, V]) Reject(f func(k K, v V) bool) Collection[K, V] {
	return c.Filter(func(k K, v V) bool {
		return !f(k, v)
	})
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

// ForgetE checks if the key exist before attempting to delete it. Should the key
// not exist, an instance of errors.KeyNotFoundError is returned. The original collection
// is always returned.
func (c Collection[K, V]) ForgetE(k K) (Collection[K, V], error) {
	if err := c.keys.ForgetE(c.keys.Search(k)); err != nil {
		return c, err
	}

	c.values.Forget(k)

	return c, nil
}
