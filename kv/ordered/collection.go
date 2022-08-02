package ordered

import (
	"reflect"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/errors"
	"github.com/thefuga/go-collections/kv"
	"github.com/thefuga/go-collections/slice"
)

type Collection[K comparable, V any] struct {
	keys   slice.Collection[K]
	values kv.Collection[K, V]
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

func Get[K comparable, T, V any](c Collection[K, V], k K) (T, error) {
	var (
		genericValue any
		getErr       error
	)

	genericValue, getErr = c.GetE(k)

	if getErr != nil {
		return *new(T), getErr
	}

	return collections.AssertE[T](genericValue)
}

func CountBy[T comparable, K comparable, V any](c Collection[K, V], f func(v V) T) map[T]int {
	return kv.CountBy(c.values, f)

}

func (c *Collection[K, V]) Put(k K, v V) Collection[K, V] {
	c.keys = append(c.keys, k)
	c.values[k] = v

	return *c
}

func (c *Collection[K, V]) Push(v V) Collection[K, V] {
	if castK, ok := collections.Assert[K](c.Count()); ok {
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

func (c Collection[K, V]) IsEmpty() bool { return c.keys.IsEmpty() }

func (c *Collection[K, V]) Get(k K) V { return c.values.Get(k) }

func (c *Collection[K, V]) GetE(k K) (V, error) { return c.values.GetE(k) }

func (c Collection[K, V]) Count() int { return c.keys.Count() }

func (c Collection[K, V]) Each(f func(k K, v V)) Collection[K, V] {
	c.keys.Each(func(_ int, k K) {
		f(k, c.values[k])
	})

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

func (c Collection[K, V]) Keys() slice.Collection[K] { return c.keys }

func (c Collection[K, V]) Sort(f func(current, next V) bool) Collection[K, V] {
	c.keys.Sort(func(i, j K) bool {
		return f(c.Get(i), c.Get(j))
	})

	return c
}

func (c Collection[K, V]) Map(f func(k K, v V) V) Collection[K, V] {
	mappedValues := makeCollection[K, V](c.Count())

	c.Each(func(k K, v V) {
		mappedValues.Put(k, f(k, v))
	})

	return mappedValues
}

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

func (c Collection[K, V]) First() V {
	first, _ := c.FirstE()
	return first
}

func (c Collection[K, V]) FirstE() (V, error) {
	first, err := c.keys.FirstE()

	if err != nil {
		return *new(V), err
	}

	return c.values[first], nil
}

func (c Collection[K, V]) FirstOrFail(f collections.Matcher) (K, V, error) {
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

func (c Collection[K, V]) Last() V { v, _ := c.LastE(); return v }

func (c Collection[K, V]) LastE() (V, error) {
	last, err := c.keys.LastE()

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

func (c Collection[K, V]) ToSliceCollection() slice.Collection[V] {
	return c.ToSlice()
}

// Combine doesn't preserve order and keys and values must be of the same type
func (c Collection[K, V]) Combine(v Collection[K, K]) Collection[K, V] {
	combined := makeCollection[K, V](c.Count())

	for i := 0; i < c.Count(); i++ {
		k, _ := collections.Assert[K](c.values[c.keys[i]])
		v, _ := collections.Assert[V](v.values[v.keys[i]])

		combined.Put(k, v)
	}

	return combined
}

// Combine doesn't preserve order and keys and values must be of the same type
func (c Collection[K, V]) CombineE(v Collection[K, K]) (Collection[K, V], error) {
	combined := makeCollection[K, V](c.Count())

	for i := 0; i < c.Count(); i++ {
		k, err := collections.AssertE[K](c.values[c.keys[i]])
		if err != nil {
			return combined, err
		}

		v, err := collections.AssertE[V](v.values[v.keys[i]])
		if err != nil {
			return combined, err
		}

		combined.Put(k, v)
	}

	return combined, nil
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
	_, _, err := c.FirstOrFail(f)
	return err == nil
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
		castKey, keyOk := collections.Assert[K](v)
		castValue, valueOk := collections.Assert[V](k)

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

func (c Collection[K, V]) ForgetE(k K) (Collection[K, V], error) {
	if err := c.keys.ForgetE(c.keys.Search(k)); err != nil {
		return c, err
	}

	c.values.Forget(k)

	return c, nil
}
