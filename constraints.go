package collections

import "github.com/thefuga/go-collections/errors"

type (
	UnsignedInteger interface {
		uint8 | uint16 | uint32 | uint64 | uint
	}

	SignedInteger interface {
		int8 | int16 | int32 | int64 | int
	}

	Integer interface {
		UnsignedInteger | SignedInteger
	}

	Float interface {
		float32 | float64
	}

	Number interface {
		Integer | Float
	}

	Relational interface {
		Number | string
	}
)

func Assert[T any](from any) (T, bool) {
	toAny, ok := from.(T)
	return toAny, ok
}

func AssertE[T any](from any) (T, error) {
	if to, ok := from.(T); ok {
		return to, nil
	}

	return *new(T), errors.NewTypeError[T](&from)
}
