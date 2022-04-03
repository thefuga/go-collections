package main

import (
	"errors"
	"reflect"
)

type Collection[K comparable, V any] map[K]V

func Collect[T any](items ...T) Collection[int, T] {
	return CollectSlice(items)
}

func CollectSlice[T any](items []T) Collection[int, T] {
	collection := make(Collection[int, T])

	for key, item := range items {
		collection[key] = item
	}

	return collection
}

func CollectMap[K comparable, V any](items Collection[K, V]) Collection[K, V] {
	return items
}

func (c Collection[K, V]) Count() int {
	return len(c)
}

func (c Collection[K, V]) Each(closure func(k K, v V)) {
	for k, v := range c {
		closure(k, v)
	}
}

func (c Collection[K, V]) Search(value V) (K, error) {
	for k, v := range c {
		if reflect.DeepEqual(v, value) {
			return k, nil
		}
	}

	return *new(K), errors.New("Value not found")
}

func (c Collection[K, V]) Keys() []K {
	keys := make([]K, c.Count())
	i := 0

	for k, _ := range c {
		keys[i] = k
		i++
	}

	return keys
}
