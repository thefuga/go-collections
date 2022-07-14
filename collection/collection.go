package collection

import (
	"reflect"
	"sort"

	"github.com/thefuga/go-collections/errors"
)

type Collection[V any] struct {
	keys   []any
	values map[any]V
}

func Collect[T any](items ...T) Collection[T] {
	return CollectSlice(items)
}

func CollectSlice[T any](items []T) Collection[T] {
	collection := makeCollection[T](len(items))

	for key, item := range items {
		collection.Put(key, item)
	}

	return collection
}

func CollectMap[V any](items map[any]V) Collection[V] {
	collection := makeCollection[V](len(items))

	for key, item := range items {
		collection.Put(key, item)
	}

	return collection
}

func makeCollection[V any](capacity int) Collection[V] {
	return Collection[V]{
		keys:   make([]any, 0, capacity),
		values: make(map[any]V, capacity),
	}
}

func Get[T, K, V any](c Collection[V], k K) (T, error) {
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

func Assert[T any](from any) T {
	return from.(T)
}

func AssertE[T any](from any) (T, error) {
	if to, ok := from.(T); ok {
		return to, nil
	}

	return *new(T), errors.NewTypeError[T](&from)
}

func (c *Collection[V]) Put(k any, v V) {
	c.keys = append(c.keys, k)
	c.values[k] = v
}

func (c *Collection[V]) Push(v V) {
	c.Put(len(c.keys), v)
}

func (c *Collection[V]) Pop() V {
	if c.IsEmpty() {
		return *new(V)
	}

	v := c.Last()

	lastKey := c.keys[len(c.keys)-1]
	c.keys = c.keys[:len(c.keys)-1]
	delete(c.values, lastKey)

	return v

}

func (c Collection[V]) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c Collection[V]) Get(k any) (V, error) {
	if item, found := c.values[k]; found {
		return item, nil
	}

	return *new(V), errors.NewKeyNotFoundError(k)
}

func (c Collection[V]) Count() int {
	return len(c.keys)
}

func (c Collection[V]) Each(closure func(k any, v V)) {
	for _, key := range c.keys {
		closure(key, c.values[key])
	}
}

// TODO change to value or closure
func (c Collection[V]) Search(value V) (any, error) {
	for _, v := range c.keys {
		if reflect.DeepEqual(c.values[v], value) {
			return v, nil
		}
	}

	return nil, errors.NewValueNotFoundError()
}

func (c Collection[V]) Keys() []any {
	return c.keys
}

func (c Collection[V]) Sort(closure func(current, next V) bool) Collection[V] {
	sort.Slice(c.keys, func(i, j int) bool {
		return closure(c.values[c.keys[i]], c.values[c.keys[j]])
	})

	return c
}

func (c Collection[V]) Map(closure func(k any, v V) V) Collection[V] {
	mappedValues := make(map[any]V)

	c.Each(func(k any, v V) {
		mappedValues[k] = closure(k, v)
	})

	return CollectMap(mappedValues)
}

func (c Collection[V]) First() V {
	if len(c.keys) == 0 {
		return *new(V)
	}

	return c.values[c.keys[0]]
}

func (c Collection[V]) FirstOrFail(match Matcher) (any, V, error) {
	var (
		found bool
		v     V
		k     any
	)

	c.Each(func(atK any, atV V) {
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

func (c Collection[V]) Last() V {
	if len(c.keys) == 0 {
		return *new(V)
	}

	return c.values[c.keys[len(c.keys)-1]]
}

func (c Collection[V]) ToSlice() []V {
	slice := make([]V, len(c.keys))

	for i, key := range c.keys {
		slice[i] = c.values[key]
	}

	return slice
}
func (c Collection[V]) Combine(v Collection[V]) Collection[V] {
	if c.Count() != v.Count() {
		return c
	}

	keys := c.values
	values := v.values

	combined := Collect[V]()

	for i := 0; i < len(keys); i++ {
		combined.Put(keys[i], values[i])
	}

	return combined
}

func (c Collection[V]) Concat(concatTo Collection[V]) Collection[V] {
	concatenated := CollectMap(c.values)

	concatTo.Each(func(k any, v V) {
		if _, ok := concatenated.values[k]; ok {
			concatenated.Push(v)
		} else {
			concatenated.Put(k, v)
		}
	})

	return concatenated
}

