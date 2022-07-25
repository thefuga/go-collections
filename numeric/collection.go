package numeric

import (
	"github.com/thefuga/go-collections/collection"
)

type (
	unsignedInteger interface {
		uint8 | uint16 | uint32 | uint64 | uint
	}

	signedInteger interface {
		int8 | int16 | int32 | int64 | int
	}

	integer interface {
		unsignedInteger | signedInteger
	}

	float interface {
		float32 | float64
	}

	number interface {
		integer | float
	}

	Collection[K comparable, V number] struct {
		collection.Collection[K, V]
	}
)

func Collect[K number](n ...K) Collection[int, K] {
	return Collection[int, K]{collection.Collect(n...)}
}

func (c Collection[K, V]) Average() V {
	return c.Sum() / V(c.Count())
}

func (c Collection[K, V]) Sum() V {
	var sum V

	for _, v := range c.ToSlice() {
		sum += v
	}

	return sum
}

func (c Collection[K, V]) Min() V {
	min := c.First()

	for _, v := range c.ToSlice() {
		if v < min {
			min = v
		}
	}

	return min
}

func (c Collection[K, V]) Max() V {
	max := c.First()

	for _, v := range c.ToSlice() {
		if v > max {
			max = v
		}
	}

	return max
}
