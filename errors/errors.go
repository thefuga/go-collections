package errors

import "fmt"

type KeyNotFoundError string

func NewKeyNotFoundError(k any) KeyNotFoundError {
	return KeyNotFoundError(
		fmt.Sprintf("Key '%v' wasn't found in the collection!", k),
	)
}

func (err KeyNotFoundError) Error() string {
	return string(err)
}

type ValueNotFoundError string

func NewValueNotFoundError() ValueNotFoundError {
	return ValueNotFoundError("Value wasn't found in the collection!")
}

func (err ValueNotFoundError) Error() string {
	return string(err)
}

type TypeError string

func NewTypeError[T any](from *any) TypeError {
	actual := getTypeString(from)
	expected := fmt.Sprintf("%T", *new(T))

	return TypeError(fmt.Sprintf(
		"interface conversion: interface {} is %s, not %s",
		actual,
		expected,
	))
}

func getTypeString(from *any) string {
	switch t := (*from).(type) {
	default:
		return fmt.Sprintf("%T", t)
	}
}

func (err TypeError) Error() string {
	return string(err)
}
