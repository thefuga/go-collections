package collections

import "testing"

func TestKeyEquals(t *testing.T) {
	if !KeyEquals("foo")("foo", nil) {
		t.Error("both keys are equal")
	}

	if KeyEquals("foo")("bar", nil) {
		t.Error("keys are different")
	}
}

func TestValueEquals(t *testing.T) {
	if !ValueEquals("foo")(nil, "foo") {
		t.Error("both values are equal")
	}

	if ValueEquals("foo")(nil, "bar") {
		t.Error("values are different")
	}
}

func TestValueDiffers(t *testing.T) {
	if !ValueDiffers("foo")(nil, "bar") {
		t.Error("values are different")
	}

	if ValueDiffers("foo")(nil, "foo") {
		t.Error("values are equal")
	}
}

func TestAsc(t *testing.T) {
	if !Asc[int]()(1, 2) {
		t.Error("1 is lesser than 2")
	}

	if Asc[int]()(2, 1) {
		t.Error("2 is greater than 1")
	}
}

func TestDesc(t *testing.T) {
	if !Desc[int]()(2, 1) {
		t.Error("2 is greater than 1")
	}

	if Desc[int]()(1, 2) {
		t.Error("1 is lesser than 2")
	}
}

func TestAssert(t *testing.T) {
	var concreteType string
	defer func() {
		if assertionErr := recover(); assertionErr != nil {
			t.Error("Unexpected error casting value.")
		}
	}()

	underlyingValue := "generic value"
	var genericType any = underlyingValue

	concreteType, _ = Assert[string](genericType)

	if concreteType != underlyingValue {
		t.Error("Expected concreteType to have the value of underlyingValue")
	}
}

func TestAssertE(t *testing.T) {
	underlyingValue := "generic value"
	var genericType any = underlyingValue

	concreteType, assertionErr := AssertE[string](genericType)

	if assertionErr != nil {
		t.Error("Unexpected error casting value.")
	}

	if concreteType != underlyingValue {
		t.Error("Expected concreteType to have the value of underlyingValue")
	}

	zeroValue, assertionErr := AssertE[int](genericType)

	if zeroValue != 0 {
		t.Error("Cast value should be zeroed when an invalid type is given")
	}

	if assertionErr.Error() != "interface conversion: interface {} is string, not int" {
		t.Error("Trying to get a value with the wrong type parameter must return a type error!")
	}
}
