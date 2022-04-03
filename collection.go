package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type Collection[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

func Collect[T any](items ...T) Collection[int, T] {
	return CollectSlice(items)
}

func CollectSlice[T any](items []T) Collection[int, T] {
	collection := Collection[int, T]{
		keys:   make([]int, len(items)),
		values: make(map[int]T, len(items)),
	}

	for key, item := range items {
		collection.keys[key] = key
		collection.values[key] = item
	}

	return collection
}

func CollectMap[K comparable, V any](items map[K]V) Collection[K, V] {
	collection := Collection[K, V]{
		keys:   make([]K, len(items)),
		values: make(map[K]V, len(items)),
	}

	keys := make([]K, len(items))
	i := 0

	for k, _ := range items {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool { return reflect.DeepEqual(i, j) })

	for i, k := range keys {
		collection.keys[i] = k
		collection.values[k] = items[k]
	}

	return collection
}

func (c Collection[K, V]) Get(k K) (V, error) {
	if item, found := c.values[k]; found {
		return item, nil
	}

	return *new(V), fmt.Errorf("Item not found")
}

func (c Collection[K, V]) Count() int {
	return len(c.keys)
}

func (c Collection[K, V]) Each(closure func(k K, v V)) {
	for _, key := range c.keys {
		closure(key, c.values[key])
	}
}

func (c Collection[K, V]) Search(value V) (K, error) {
	for _, v := range c.keys {
		if reflect.DeepEqual(c.values[v], value) {
			return v, nil
		}
	}

	return *new(K), errors.New("Value not found")
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

func (c Collection[K, V]) First() V {
	if len(c.keys) == 0 {
		return *new(V)
	}

	return c.values[c.keys[0]]
}

func (c Collection[K, V]) Last() V {
	if len(c.keys) == 0 {
		return *new(V)
	}

	return c.values[c.keys[len(c.keys)-1]]
}
