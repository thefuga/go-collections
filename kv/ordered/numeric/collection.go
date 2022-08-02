package numeric

import (
	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/kv/ordered"
)

type Collection[K comparable, V collections.Number] struct {
	ordered.Collection[K, V]
}

func Collect[K collections.Number](n ...K) Collection[int, K] {
	return Collection[int, K]{ordered.Collect(n...)}
}

func (c Collection[K, V]) Average() V {
	return collections.Average(c.ToSlice())
}

func (c Collection[K, V]) Sum() V {
	return collections.Sum(c.ToSlice())
}

func (c Collection[K, V]) Min() V {
	return collections.Min(c.ToSlice())
}

func (c Collection[K, V]) Max() V {
	return collections.Max(c.ToSlice())
}

func (c Collection[K, V]) Median() float64 {
	return collections.Median(c.ToSlice())
}
