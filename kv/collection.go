package kv

import (
	"reflect"
	"sort"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/errors"
)

type Collection[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

func Collect[T any](items ...T) Collection[int, T] {
	return CollectSlice(items)
}

func CollectSlice[T any](items []T) Collection[int, T] {
	collection := makeCollection[int, T](len(items))

	for key, item := range items {
		collection.Put(key, item)
	}

	return collection
}

func CollectMap[K comparable, V any](items map[K]V) Collection[K, V] {
	collection := makeCollection[K, V](len(items))

	for key, item := range items {
		collection.Put(key, item)
	}

	return collection
}

func makeCollection[K comparable, V any](capacity int) Collection[K, V] {
	return Collection[K, V]{
		keys:   make([]K, 0, capacity),
		values: make(map[K]V, capacity),
	}
}

func Get[K comparable, T, V any](c Collection[K, V], k K) (T, error) {
	var (
		genericValue any
		getErr       error
	)

	genericValue, getErr = c.GetE(k)

	if getErr != nil {
		return *new(T), getErr
	}

	return AssertE[T](genericValue)
}

func Assert[T any](from any) (T, bool) {
	toAny, ok := from.(T)
	return toAny, ok
}

func AssertE[T any](from any) (T, error) {
	if to, ok := from.(T); ok {
		return to, nil
	}

	return *new(T), errors.NewTypeError[T](&from)
}

func CountBy[T comparable, K comparable, V any](c Collection[K, V], f func(v V) T) map[T]int {
	count := map[T]int{}

	c.Each(func(_ K, v V) {
		count[f(v)]++
	})

	return count
}

func (c *Collection[K, V]) Put(k K, v V) Collection[K, V] {
	c.keys = append(c.keys, k)
	c.values[k] = v

	return *c
}

func (c *Collection[K, V]) Push(v V) Collection[K, V] {
	if castK, ok := Assert[K](c.Count()); ok {
		c.Put(castK, v)
	}

	return *c
}

func (c *Collection[K, V]) Pop() V {
	v, _ := c.PopE()
	return v
}

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

func (c Collection[K, V]) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *Collection[K, V]) Get(k K) V {
	v, _ := c.GetE(k)
	return v
}

func (c *Collection[K, V]) GetE(k K) (V, error) {
	item, found := c.values[k]

	if !found {
		return *new(V), errors.NewKeyNotFoundError(k)
	}

	return item, nil
}

func (c Collection[K, V]) Count() int {
	return len(c.keys)
}

func (c Collection[K, V]) Each(f func(k K, v V)) Collection[K, V] {
	for _, key := range c.keys {
		f(key, c.values[key])
	}

	return c
}

func (c Collection[K, V]) Tap(f func(Collection[K, V])) Collection[K, V] {
	f(c)
	return c
}

func (c Collection[K, V]) Search(v V) K {
	k, _ := c.SearchE(v)
	return k
}

func (c Collection[K, V]) SearchE(v V) (K, error) {
	for _, k := range c.keys {
		if reflect.DeepEqual(c.values[k], v) {
			return k, nil
		}
	}

	return *new(K), errors.NewValueNotFoundError()
}

func (c Collection[K, V]) Keys() []K { return c.keys }

func (c Collection[K, V]) Sort(f func(current, next V) bool) Collection[K, V] {
	sort.Slice(c.keys, func(i, j int) bool {
		return f(c.values[c.keys[i]], c.values[c.keys[j]])
	})

	return c
}

func (c Collection[K, V]) Map(f func(k K, v V) V) Collection[K, V] {
	mappedValues := make(map[K]V)

	c.Each(func(k K, v V) {
		mappedValues[k] = f(k, v)
	})

	return CollectMap(mappedValues)
}

func (c Collection[K, V]) Only(keys []K) Collection[K, V] {
	onlyValues := make(map[K]V)

	for _, key := range keys {
		value, err := c.GetE(key)
		if err == nil {
			onlyValues[key] = value
		}
	}

	return CollectMap(onlyValues)
}

func (c Collection[K, V]) First() V {
	first, _ := c.FirstE()
	return first
}

func (c Collection[K, V]) FirstE() (V, error) {
	first, err := collections.FirstE(c.keys)

	if err != nil {
		return *new(V), err
	}

	return c.values[first], nil
}

func (c Collection[K, V]) FirstOrFail(f collections.Matcher) (any, V, error) {
	var (
		found bool
		v     V
		k     any
	)

	c.Each(func(atK K, atV V) {
		if !f(atK, atV) {
			return
		}

		found = true
		k = atK
		v = atV
	})

	if !found {
		return nil, *new(V), errors.NewValueNotFoundError()
	}

	return k, v, nil
}

func (c Collection[K, V]) Last() V { v, _ := c.LastE(); return v }

func (c Collection[K, V]) LastE() (V, error) {
	last, err := collections.LastE(c.keys)

	if err != nil {
		return *new(V), err
	}

	return c.values[last], nil
}

func (c Collection[K, V]) ToSlice() []V {
	slice := make([]V, len(c.keys))

	for i, key := range c.keys {
		slice[i] = c.values[key]
	}

	return slice
}

// Combine doesn't preserve order and keys and values must be of the same type
func (c Collection[K, V]) Combine(v Collection[K, K]) Collection[K, V] {
	combined := makeCollection[K, V](c.Count())

	for i := 0; i < c.Count(); i++ {
		k, _ := Assert[K](c.values[c.keys[i]])
		v, _ := Assert[V](v.values[v.keys[i]])

		combined.Put(k, v)
	}

	return combined
}

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

func (c Collection[K, V]) Contains(f collections.Matcher) bool {
	var contains bool

	c.Each(func(k K, v V) {
		if f(k, v) {
			contains = true
			return
		}
	})

	return contains
}

func (c Collection[K, V]) Every(f collections.Matcher) bool {
	contains := true

	c.Each(func(k K, v V) {
		if !f(k, v) {
			contains = false
			return
		}
	})

	return contains
}

func (c Collection[K, V]) Flip() Collection[K, V] {
	flippedKeys := make([]K, 0, c.Count())
	flippedValues := make(map[K]V, c.Count())

	c.Each(func(k K, v V) {
		castKey, keyOk := Assert[K](v)
		castValue, valueOk := Assert[V](k)

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

func (c Collection[K, V]) Merge(other Collection[K, V]) Collection[K, V] {
	other.Each(func(k K, v V) {
		c.Put(k, v)
	})

	return c
}

func (c Collection[K, V]) Filter(f func(k K, v V) bool) Collection[K, V] {
	values := make(map[K]V)

	c.Each(func(k K, v V) {
		if f(k, c.values[k]) {
			values[k] = c.values[k]
		}
	})

	return CollectMap(values)
}

func (c Collection[K, V]) Reject(f func(k K, v V) bool) Collection[K, V] {
	return c.Filter(func(k K, v V) bool {
		return !f(k, v)
	})
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
