package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type Collection[K comparable, V any] struct {
	keys     []K
	keyValue map[K]V
}

func Collect[T any](items ...T) Collection[int, T] {
	return CollectSlice(items)
}

func CollectSlice[T any](items []T) Collection[int, T] {
	collection := Collection[int, T]{
		keys:     make([]int, len(items)),
		keyValue: make(map[int]T, len(items)),
	}

	for key, item := range items {
		collection.keys[key] = key
		collection.keyValue[key] = item
	}

	return collection
}

func CollectMap[K comparable, V any](items map[K]V) Collection[K, V] {
	collection := Collection[K, V]{
		keys:     make([]K, len(items)),
		keyValue: make(map[K]V, len(items)),
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
		collection.keyValue[k] = items[k]
	}

	return collection
}

func (c Collection[K, V]) Get(k K) (V, error) {
	if item, found := c.keyValue[k]; found {
		return item, nil
	}

	return *new(V), fmt.Errorf("Item not found")
}

func (c Collection[K, V]) Count() int {
	return len(c.keys)
}

func (c Collection[K, V]) Each(closure func(k K, v V)) {
	for _, key := range c.keys {
		closure(key, c.keyValue[key])
	}
}

func (c Collection[K, V]) Search(value V) (K, error) {
	for _, v := range c.keys {
		if reflect.DeepEqual(c.keyValue[v], value) {
			return v, nil
		}
	}

	return *new(K), errors.New("Value not found")
}

func (c Collection[K, V]) Keys() []K {
	return c.keys
}
