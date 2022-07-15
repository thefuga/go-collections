package collection

import "reflect"

type (
	Matcher func(key any, value any) bool

	Relational interface {
		uint8 | uint16 | uint32 | uint64 | uint |
			int8 | int16 | int32 | int64 | int |
			float32 | float64 | string
	}
)

func KeyEquals(key any) Matcher {
	return func(collectionKey any, _ any) bool {
		return key == collectionKey
	}
}

func ValueEquals(value any) Matcher {
	return func(_ any, collectionValue any) bool {
		return reflect.DeepEqual(value, collectionValue)
	}
}

func ValueDiffers(value any) Matcher {
	return func(_ any, collectionValue any) bool {
		return !reflect.DeepEqual(value, collectionValue)
	}
}

func Asc[T Relational]() func(T, T) bool {
	return func(current, next T) bool {
		return current < next
	}
}

func Desc[T Relational]() func(T, T) bool {
	return func(current, next T) bool {
		return current > next
	}
}
