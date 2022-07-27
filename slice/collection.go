package slice

import (
	"reflect"
	"sort"

	"github.com/thefuga/go-collections/errors"
)

type Collection[V any] []V

func Collect[V any](values ...V) Collection[V] {
	return append(make(Collection[V], 0, len(values)), values...)
}

func (c Collection[V]) Get(i int) (V, error) {
	if i < 0 || c.Count() <= i {
		return *new(V), errors.NewValueNotFoundError()
	}

	return c[i], nil
}

func (c Collection[V]) Push(v V) Collection[V] {
	return append(c, v)
}

func (c Collection[V]) Put(i int, v V) Collection[V] {
	if c.IsEmpty() {
		c = make(Collection[V], i+1)
		c[i] = v
		return c
	}

	if cap(c) < c.Count()+1 {
		newC := make(Collection[V], 0, c.Count()+1)
		c = append(newC, c...)
	}

	return c.put(i, v)
}

func (c Collection[V]) put(i int, v V) Collection[V] {
	c = append(c[:i+1], c[i:]...)
	c[i] = v
	return c
}

func (c *Collection[V]) Pop() (V, error) {
	v, err := c.Last()

	if err != nil {
		return v, err
	}

	*c = (*c)[:c.Count()-1]

	return v, nil
}

func (c Collection[V]) Count() int {
	return len(c)
}

func (c Collection[V]) Capacity() int {
	return cap(c)
}

func (c Collection[V]) IsEmpty() bool {
	return c.Count() == 0
}

func (c Collection[V]) Each(f func(i int, v V)) Collection[V] {
	for i, v := range c {
		f(i, v)
	}

	return c
}

func (c Collection[V]) Tap(f func(Collection[V])) Collection[V] {
	f(c)

	return c
}

func (c Collection[V]) Search(v V) (int, error) {
	for i := range c {
		if reflect.DeepEqual(c[i], v) {
			return i, nil
		}
	}

	return -1, errors.NewValueNotFoundError()
}

func (c Collection[V]) Sort(f func(current, next V) bool) Collection[V] {
	sort.Slice(c, func(i, j int) bool {
		return f(c[i], c[j])
	})

	return c
}

func (c Collection[V]) Map(f func(i int, v V) V) Collection[V] {
	mappedValues := make(Collection[V], 0, c.Count())

	c.Each(func(_ int, v V) {
		mappedValues = mappedValues.Push(v)
	})

	return mappedValues
}

func (c Collection[V]) First() (V, error) {
	v, err := c.Get(0)

	if err != nil {
		return v, errors.NewEmptyCollectionError(err)
	}

	return v, nil
}

func (c Collection[V]) Last() (V, error) {
	v, err := c.Get(c.Count() - 1)

	if err != nil {
		return v, errors.NewEmptyCollectionError(err)
	}

	return v, nil
}
