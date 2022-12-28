package collections

import "reflect"

// Matcher is used by matchers on functions that  must compare keys and values from
// a collection.
// It is used as a functional option. To learn more, see: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Matcher func(key any, value any) bool

// KeyEquals builds a matcher to compare the given key to the key passed by the matcher caller.
func KeyEquals(key any) Matcher {
	return func(collectionKey any, _ any) bool {
		return key == collectionKey
	}
}

// ValueEquals builds a matcher to compare the given value (with reflect.DeepEqual)
// to the value passed by the matcher caller.
func ValueEquals(value any) Matcher {
	return func(_ any, collectionValue any) bool {
		return reflect.DeepEqual(value, collectionValue)
	}
}

// ValueDiffers builds a matcher to compare the given value (with reflect.DeepEqual)
// to the value passed by the matcher caller. It has the opposite behavior from ValueEquals
func ValueDiffers(value any) Matcher {
	return func(_ any, collectionValue any) bool {
		return !reflect.DeepEqual(value, collectionValue)
	}
}

// ValueGT compares the given numeric value to check if it is greater than the value
// given by the matcher's caller.
func ValueGT[T Number](value T) Matcher {
	return func(_ any, collectionValue any) bool {
		if cast, ok := collectionValue.(T); ok {
			return value < cast
		}
		return false
	}
}

// ValueLT compares the given numeric value to check if it is lesser than the value
// given by the matcher's caller.
func ValueLT[T Number](value T) Matcher {
	return func(_ any, collectionValue any) bool {
		if cast, ok := collectionValue.(T); ok {
			return value > cast
		}
		return false
	}
}

// FieldEquals uses FieldMatch composed with ValueEquals as the matcher.
func FieldEquals[V any](field string, value any) Matcher {
	return FieldMatch[V](field, ValueEquals(value))
}

// FieldMatch will attempt to retrieve the value corresponding to the given struct
// field name. V must be a struct, otherwise calls to the matcher will always return false.
// The retrieved value will be used to supply the given matcher.
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

// And combines all the given matchers into a single matcher which returns true
// when all matchers return true.
func And[V any](matchers ...Matcher) Matcher {
	return func(i any, collectionValue any) bool {
		match := true

		for _, matcher := range matchers {
			match = match && matcher(i, collectionValue)
		}

		return match
	}
}

// AndValue is similar to And, but it receives matchers wrapped by a function which
// will receive v. It is useful to compare build matchers dynamically at the execution time
// rather than at the function's call time (i.e. the composed matchers won't be called until
// the higher order matcher is called).
func AndValue[V any](v V, matchers ...func(V) Matcher) Matcher {
	return func(i any, collectionValue any) bool {
		match := true

		for _, matcher := range matchers {
			match = match && matcher(v)(i, collectionValue)
		}

		return match
	}
}

// Or combines all the given matchers into a single matcher which returns true
// when at least one of the given matcher returns true.
func Or[V any](matchers ...Matcher) Matcher {
	return func(i any, collectionValue any) bool {
		match := true

		for _, matcher := range matchers {
			match = match || matcher(i, collectionValue)
		}

		return match
	}
}

// OrValue is similar to Or, but it receives matchers wrapped by a function which
// will receive v. It is useful to compare build matchers dynamically at the execution time
// rather than at the function's call time (i.e. the composed matchers won't be called until
// the higher order matcher is called).
func OrValue[V any](v V, matchers ...func(V) Matcher) Matcher {
	return func(i any, collectionValue any) bool {
		match := true

		for _, matcher := range matchers {
			match = match || matcher(v)(i, collectionValue)
		}

		return match
	}
}

// Asc can be used as a Sort param to order collections in ascending order. It only
// works on slices holding Relational values.
func Asc[T Relational]() func(T, T) bool {
	return func(current, next T) bool {
		return current < next
	}
}

// Desc can be used as a Sort param to order collections in descending order. It only
// works on slices holding Relational values.
func Desc[T Relational]() func(T, T) bool {
	return func(current, next T) bool {
		return current > next
	}
}
