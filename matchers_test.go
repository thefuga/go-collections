package collections

import (
	"testing"

	"github.com/thefuga/go-collections/internal"
)

func TestKeyEquals(t *testing.T) {
	if !KeyEquals("foo")("foo", nil) {
		t.Error("both keys are equal")
	}

	if KeyEquals("foo")("bar", nil) {
		t.Error("keys are different")
	}
}

func TestValueEquals(t *testing.T) {
	if !ValueEquals[int]("foo")(1, "foo") {
		t.Error("both values are equal")
	}

	if ValueEquals[int]("foo")(1, "bar") {
		t.Error("values are different")
	}
}

func TestValueDeepEquals(t *testing.T) {
	type user struct {
		name string
	}
	bob1 := user{"bob"}
	bob2 := user{"bob"}
	alice := user{"alice"}

	if !ValueDeepEquals[int](bob1)(1, bob2) {
		t.Error("both values are equal")
	}

	if ValueDeepEquals[int](bob1)(1, alice) {
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

	concreteType, _ = internal.Assert[string](genericType)

	if concreteType != underlyingValue {
		t.Error("Expected concreteType to have the value of underlyingValue")
	}
}

func TestAssertE(t *testing.T) {
	underlyingValue := "generic value"
	var genericType any = underlyingValue

	concreteType, assertionErr := internal.AssertE[string](genericType)

	if assertionErr != nil {
		t.Error("Unexpected error casting value.")
	}

	if concreteType != underlyingValue {
		t.Error("Expected concreteType to have the value of underlyingValue")
	}

	zeroValue, assertionErr := internal.AssertE[int](genericType)

	if zeroValue != 0 {
		t.Error("Cast value should be zeroed when an invalid type is given")
	}

	if assertionErr.Error() != "interface conversion: interface {} is string, not int" {
		t.Error("Trying to get a value with the wrong type parameter must return a type error!")
	}
}

func TestFieldEquals(t *testing.T) {
	u := user{Name: "Jon", Email: "jon@collections.go", Age: 33}

	if !FieldEquals[user]("Name", u.Name)(0, u) {
		t.Error("user should've matched")
	}
}

func TestFieldMatch(t *testing.T) {
	u := user{Name: "Jon", Email: "jon@collections.go", Age: 33}

	if !FieldMatch[user]("Age", ValueCastGT(30))(0, u) {
		t.Error("user should've matched")
	}
}

func TestNot(t *testing.T) {
	matcher := ValueEquals[int](1)
	notMatcher := Not(matcher)

	if notMatcher(0, 1) == matcher(0, 1) {
		t.Error("got the same result using Not.")
	}
}

func TestAnd(t *testing.T) {
	i := 10

	if !And[int](ValueCastGT(9), ValueCastLT(11))(0, i) {
		t.Error("10 is greater than 9 and lesser than 11")
	}
}

func TestAndValue(t *testing.T) {
	i := 10

	if !AndValue[int, int](
		11,
		func(v int) Matcher[int, int] { return ValueGT[int](-v) },
		func(v int) Matcher[int, int] { return ValueLT[int](v) },
	)(0, i) {
		t.Error("10 is greater than -11 and lesser than 11")
	}
}

func TestOr(t *testing.T) {
	i := 11

	if !Or[int](ValueGT[int](9), ValueLT[int](10))(0, i) {
		t.Error("11 is greater than 9 and 10")
	}
}

func TestOrValue(t *testing.T) {
	i := 11

	if !OrValue[int](
		10,
		func(v int) Matcher[int, int] { return ValueGT[int](v) },
		func(v int) Matcher[int, int] { return ValueLT[int](2 * v) },
	)(0, i) {
		t.Error("11 is greater than 10 and lesser than 20")
	}
}
