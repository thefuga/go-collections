package collection

import (
	"reflect"
	"sort"

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

	genericValue, getErr = c.Get(k)

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
	if c.IsEmpty() {
		return *new(V)
	}

	v := c.Last()

	lastKey := c.keys[len(c.keys)-1]
	c.keys = c.keys[:len(c.keys)-1]
	delete(c.values, lastKey)

	return v

}

func (c Collection[K, V]) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *Collection[K, V]) Get(k K) (V, error) {
	item, found := c.values[k]

	if !found {
		return *new(V), errors.NewKeyNotFoundError(k)
	}

	return item, nil
}

func (c Collection[K, V]) Count() int {
	return len(c.keys)
}

func (c Collection[K, V]) Each(closure func(k K, v V)) Collection[K, V] {
	for _, key := range c.keys {
		closure(key, c.values[key])
	}

	return c
}

func (c Collection[K, V]) Tap(closure func(Collection[K, V])) Collection[K, V] {
	closure(c)

	return c
}

func (c Collection[K, V]) Search(value V) (K, error) {
	for _, v := range c.keys {
		if reflect.DeepEqual(c.values[v], value) {
			return v, nil
		}
	}

	return *new(K), errors.NewValueNotFoundError()
}

func (c Collection[K, V]) Keys() []K {
	return c.keys
}

func (c Collection[K, V]) Sort(closure func(current, next V) bool) Collection[K, V] {
	sort.Slice(c.keys, func(i, j int) bool {
		return closure(c.values[c.keys[i]], c.values[c.keys[j]])
	})

	return c
}

func (c Collection[K, V]) Map(closure func(k K, v V) V) Collection[K, V] {
	mappedValues := make(map[K]V)

	c.Each(func(k K, v V) {
		mappedValues[k] = closure(k, v)
	})

	return CollectMap(mappedValues)
}

func (c Collection[K, V]) Only(keys []K) Collection[K, V] {
	onlyValues := make(map[K]V)

	for _, key := range keys {
		value, err := c.Get(key)
		if err == nil {
			onlyValues[key] = value
		}
	}

	return CollectMap(onlyValues)
}

func (c Collection[K, V]) First() V {
	if len(c.keys) == 0 {
		return *new(V)
	}

	return c.values[c.keys[0]]
}

func (c Collection[K, V]) FirstOrFail(match Matcher) (any, V, error) {
	var (
		found bool
		v     V
		k     any
	)

	c.Each(func(atK K, atV V) {
		if !match(atK, atV) {
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

func (c Collection[K, V]) Last() V {
	if len(c.keys) == 0 {
		return *new(V)
	}

	return c.values[c.keys[len(c.keys)-1]]
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

func (c Collection[K, V]) Contains(match Matcher) bool {
	var contains bool

	c.Each(func(k K, v V) {
		if match(k, v) {
			contains = true
			return
		}
	})

	return contains
}

func (c Collection[K, V]) Every(match Matcher) bool {
	contains := true

	c.Each(func(k K, v V) {
		if !match(k, v) {
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

func (c Collection[K, V]) Merge(other Collection[K, V]) (Collection[K, V], error) {
	other.Each(func(k K, v V) {
		c.Put(k, v)
	})

	return c, nil
}
