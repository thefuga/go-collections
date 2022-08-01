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

func FieldEquals[V any](field string, value any) Matcher {
	return FieldMatch[V](field, ValueEquals(value))
}

func FieldMatch[V any](field string, matcher Matcher) Matcher {
	return func(_, v any) bool {
		cast, ok := v.(V)
		if !ok {
			return false
		}

		fieldVal := reflect.ValueOf(&cast).Elem()

		for fieldNum := 0; fieldNum < fieldVal.NumField(); fieldNum++ {
			if fieldName := fieldVal.Type().Field(fieldNum).Name; fieldName == field {
				if matcher(0, fieldVal.Field(fieldNum).Interface()) {
					return true
				}
			}
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
