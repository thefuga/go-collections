// errors holds custom errors common to all collections types.
package errors

import "fmt"

type KeyNotFoundError error

func NewKeyNotFoundError(k any, cause ...error) KeyNotFoundError {
	return wrap("key '%v' not found", []any{k}, cause)
}

type ValueNotFoundError error

func NewValueNotFoundError(cause ...error) ValueNotFoundError {
	return wrap("value not found", nil, cause)
}

type TypeError error

func NewTypeError[T any](from *any, cause ...error) TypeError {
	actual := getTypeString(from)
	expected := fmt.Sprintf("%T", *new(T))

	return wrap(
		"interface conversion: interface {} is %s, not %s",
		[]any{actual, expected},
		cause,
	)
}

func getTypeString(from *any) string {
	switch t := (*from).(type) {
	default:
		return fmt.Sprintf("%T", t)
	}
}

type EmptyCollectionError error

func NewEmptyCollectionError(cause ...error) error {
	return wrap("empty collection", nil, cause)
}

type IntexOutOfBoundsError error

func NewIndexOutOfBoundsError(cause ...error) error {
	return wrap("index out of bounds", nil, cause)
}

type KeysValuesLegthMismatch error

func NewKeysValuesLengthMismatch(cause ...error) error {
	return wrap("keys and values don't have the same length", nil, cause)
}

func wrap(format string, args []any, cause []error) error {
	msg := fmt.Sprintf(format, args...)

	if len(cause) > 0 {
		return fmt.Errorf("%w: %s", cause[0], msg)
	}

	return fmt.Errorf(msg)
}
