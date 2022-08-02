package kv

import (
	"fmt"
	"reflect"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/errors"
	"github.com/thefuga/go-collections/slice"
)

type Collection[K comparable, V any] map[K]V

func Collect[T any](items ...T) Collection[int, T] {
	return CollectSlice(items)
}

func CollectSlice[T any](items []T) Collection[int, T] {
	collection := make(Collection[int, T], len(items))

	for key, item := range items {
		collection.Put(key, item)
	}

	return collection
}

func CollectMap[K comparable, V any](items map[K]V) Collection[K, V] {
	return items
}

func Combine[K comparable, V any](
	keys slice.Collection[K], values slice.Collection[V],
) Collection[K, V] {
	if keys.Count() != values.Count() {
		return nil
	}

	count := keys.Count()

	collection := make(Collection[K, V], keys.Count())

	for i := 0; i < count; i++ {
		collection.Put(keys.Get(i), values.Get(i))
	}

	return collection
}

func CombineE[K comparable, V any](
	keys slice.Collection[K], values slice.Collection[V],
) (Collection[K, V], error) {
	if keys.Count() != values.Count() {
		return nil, fmt.Errorf("") // TODO new  error
	}

	count := keys.Count()

	collection := make(Collection[K, V], keys.Count())

	for i := 0; i < count; i++ {
		collection.Put(keys.Get(i), values.Get(i))
	}

	return collection, nil
}

func CountBy[T comparable, K comparable, V any](c Collection[K, V], f func(v V) T) map[T]int {
	count := map[T]int{}

	c.Each(func(_ K, v V) {
		count[f(v)]++
	})

	return count
}

func (c Collection[K, V]) Each(f func(k K, v V)) Collection[K, V] {
	for key, value := range c {
		f(key, value)
	}

	return c
}

func (c *Collection[K, V]) Get(k K) V {
	v, _ := c.GetE(k)
	return v
}

func (c Collection[K, V]) GetE(k K) (V, error) {
	item, found := c[k]

	if !found {
		return *new(V), errors.NewKeyNotFoundError(k)
	}

	return item, nil
}

func (c Collection[K, V]) Search(v V) K {
	k, _ := c.SearchE(v)
	return k
}

func (c Collection[K, V]) SearchE(v V) (K, error) {
	for k, collectionV := range c {
		if reflect.DeepEqual(collectionV, v) {
			return k, nil
		}
	}

	return *new(K), errors.NewValueNotFoundError()
}

func (c Collection[K, V]) Map(f func(k K, v V) V) Collection[K, V] {
	mapped := make(Collection[K, V], c.Count())

	c.Each(func(k K, v V) {
		mapped.Put(k, f(k, v))
	})

	return mapped
}

func (c Collection[K, V]) Count() int {
	return len(c)
}

func (c Collection[K, V]) Put(k K, v V) Collection[K, V] {
	c[k] = v

	return c
}

func (c Collection[K, V]) IsEmpty() bool { return len(c) == 0 }

func (c Collection[K, V]) Keys() slice.Collection[K] {
	keys := make(slice.Collection[K], 0, c.Count())

	c.Each(func(k K, _ V) {
		keys.Push(k)
	})

	return keys
}

func (c Collection[K, V]) Values() slice.Collection[V] {
	values := make(slice.Collection[V], 0, c.Count())

	c.Each(func(_ K, v V) {
		values.Push(v)
	})

	return values
}

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

func (c Collection[K, V]) Tap(f func(Collection[K, V])) Collection[K, V] {
	f(c)
	return c
}

func (c Collection[K, V]) ToSlice() []V {
	slice := make([]V, c.Count())

	for _, value := range c {
		slice = append(slice, value)
	}

	return slice
}

func (c Collection[K, V]) ToSliceCollection() slice.Collection[V] {
	return c.ToSlice()
}

func (c Collection[K, V]) Concat(concatTo Collection[K, V]) Collection[K, V] {
	concatenated := make(Collection[K, V], c.Count())

	concatTo.Each(func(k K, v V) {
		concatenated.Put(k, v)
	})

	return concatenated
}

// TODO
func (c Collection[K, V]) Contains(f collections.Matcher) bool {
	return true
}

func (c Collection[K, V]) Every(f collections.Matcher) bool {
	contains := true

	// TODO change this to return early. Add more test scenarios.
	c.Each(func(k K, v V) {
		if !f(k, v) {
			contains = false
			return
		}
	})

	return contains
}

func (c Collection[K, V]) Flip() Collection[K, V] {
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

	return flippedValues
}

//TODO check if makes sense to have Merge and Concat
func (c Collection[K, V]) Merge(other Collection[K, V]) Collection[K, V] {
	other.Each(func(k K, v V) {
		c.Put(k, v)
	})

	return c
}

func (c Collection[K, V]) Filter(f func(k K, v V) bool) Collection[K, V] {
	filtered := make(map[K]V)

	c.Each(func(k K, v V) {
		if f(k, v) {
			filtered[k] = c[k]
		}
	})

	return filtered
}

func (c Collection[K, V]) Reject(f func(k K, v V) bool) Collection[K, V] {
	return c.Filter(func(k K, v V) bool {
		return !f(k, v)
	})
}

func (c Collection[K, V]) Forget(k K) Collection[K, V] {
	delete(c, k)
	return c
}

func (c Collection[K, V]) When(
	execute bool, f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	if !execute {
		return c
	}

	return f(c)
}

func (c Collection[K, V]) WhenEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.When(c.IsEmpty(), f)
}

func (c Collection[K, V]) WhenNotEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.When(!c.IsEmpty(), f)
}

func (c Collection[K, V]) Unless(
	execute bool, f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.When(!execute, f)
}

func (c Collection[K, V]) UnlessEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.WhenNotEmpty(f)
}

func (c Collection[K, V]) UnlessNotEmpty(
	f func(collection Collection[K, V]) Collection[K, V],
) Collection[K, V] {
	return c.WhenEmpty(f)
}
