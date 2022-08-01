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

func TestFieldEquals(t *testing.T) {
	u := user{Name: "Jon", Email: "jon@collections.go", Age: 33}

	if !FieldEquals[user]("Name", u.Name)(0, u) {
		t.Error("user should've matched")
	}
}

func TestFieldMatch(t *testing.T) {
	u := user{Name: "Jon", Email: "jon@collections.go", Age: 33}

	if !FieldMatch[user]("Age", ValueGT(30))(0, u) {
		t.Error("user should've matched")
	}
}

func TestAnd(t *testing.T) {
	i := 10

	if !And[int](ValueGT(9), ValueLT(11))(0, i) {
		t.Error("10 is greater than 9 and lesser than 11")
	}
}

func TestAndValue(t *testing.T) {
	i := 10

	if !AndValue[int](
		11,
		func(v int) Matcher { return ValueGT(-v) },
		func(v int) Matcher { return ValueLT(v) },
	)(0, i) {
		t.Error("10 is greater than -11 and lesser than 11")
	}
}

func TestOr(t *testing.T) {
	i := 11

	if !Or[int](ValueGT(9), ValueLT(10))(0, i) {
		t.Error("11 is greater than 9 and 10")
	}
}

func TestOrValue(t *testing.T) {
	i := 11

	if !OrValue[int](
		10,
		func(v int) Matcher { return ValueGT(v) },
		func(v int) Matcher { return ValueLT(2 * v) },
	)(0, i) {
		t.Error("11 is greater than 10 and lesser than 20")
	}
}
