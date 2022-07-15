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

	Collection[T number] struct {
		collection.Collection[T]
	}
)

func Collect[T number](n ...T) Collection[T] {
	return Collection[T]{collection.Collect(n...)}
}

func (c Collection[T]) Average() T {
	return c.Sum() / T(c.Count())
}

func (c Collection[T]) Sum() T {
	var sum T

	for _, v := range c.ToSlice() {
		sum += v
	}

	return sum
}
