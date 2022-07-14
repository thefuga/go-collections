package collection

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
