package slice

import (
	"github.com/thefuga/go-collections"
)

type Collection[V any] []V

func Collect[V any](values ...V) Collection[V] {
	return append(make(Collection[V], 0, len(values)), values...)
}

func (c Collection[V]) Get(i int) V { return collections.Get(i, c) }

func (c Collection[V]) GetE(i int) (V, error) { return collections.GetE(i, c) }

func (c Collection[V]) Push(v V) Collection[V] { return append(c, v) }

func (c Collection[V]) Put(i int, v V) Collection[V] { return collections.Put(i, v, c) }

func (c *Collection[V]) Pop() V { return collections.Pop((*[]V)(c)) }

func (c *Collection[V]) PopE() (V, error) { return collections.PopE((*[]V)(c)) }

func (c Collection[V]) Count() int { return len(c) }

func (c Collection[V]) Capacity() int { return cap(c) }

func (c Collection[V]) IsEmpty() bool { return c.Count() == 0 }

func (c Collection[V]) Search(v V) int { return collections.Search(v, c) }

func (c Collection[V]) SearchE(v V) (int, error) { return collections.SearchE(v, c) }

func (c Collection[V]) Map(f func(i int, v V) V) Collection[V] { return collections.Map(f, c) }

func (c Collection[V]) First() V { return collections.First(c) }

func (c Collection[V]) FirstE() (V, error) { return collections.FirstE(c) }

func (c Collection[V]) Last() V { return collections.Last(c) }

func (c Collection[V]) LastE() (V, error) { return collections.LastE(c) }

func (c Collection[V]) Each(f func(i int, v V)) Collection[V] {
	collections.Each(f, c)
	return c
}

func (c Collection[V]) Sort(f func(current, next V) bool) Collection[V] {
	collections.Sort(c, f)
	return c
}

func (c Collection[V]) Tap(f func(Collection[V])) Collection[V] {
	f(c)
	return c
}

func (c *Collection[V]) ForgetE(i int) error {
	return collections.ForgetE((*[]V)(c), i)
}
