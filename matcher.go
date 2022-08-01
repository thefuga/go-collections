package collections

import "reflect"

type Matcher func(key any, value any) bool

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

func ValueGT[T Number](value T) Matcher {
	return func(_ any, collectionValue any) bool {
		if cast, ok := collectionValue.(T); ok {
			return value < cast
		}
		return false
	}
}

func ValueLT[T Number](value T) Matcher {
	return func(_ any, collectionValue any) bool {
		if cast, ok := collectionValue.(T); ok {
			return value > cast
		}
		return false
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
