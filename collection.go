package main

import (
	"reflect"
	"sort"
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

func Assert[T any](from any) T {
	return from.(T)
}

func AssertE[T any](from any) (T, error) {
	if to, ok := from.(T); ok {
		return to, nil
	}

	return *new(T), NewTypeError[T](&from)
}

func (c *Collection[V]) Put(k any, v V) {
	c.keys = append(c.keys, k)
	c.values[k] = v
}

func (c *Collection[V]) Push(v V) {
	c.Put(len(c.keys), v)
}

func (c Collection[V]) Get(k any) (V, error) {
	if item, found := c.values[k]; found {
		return item, nil
	}

	return *new(V), NewKeyNotFoundError(k)
}

func (c Collection[V]) Count() int {
	return len(c.keys)
}

func (c Collection[V]) Each(closure func(k any, v V)) {
	for _, key := range c.keys {
		closure(key, c.values[key])
	}
}

func (c Collection[V]) Search(value V) (any, error) {
	for _, v := range c.keys {
		if reflect.DeepEqual(c.values[v], value) {
			return v, nil
		}
	}

	return nil, NewValueNotFoundError()
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

func (c Collection[V]) Last() V {
	if len(c.keys) == 0 {
		return *new(V)
	}

	return c.values[c.keys[len(c.keys)-1]]
}
