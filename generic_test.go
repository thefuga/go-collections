package collections

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		description string
		sut         []int
		i           int
		v           int
		err         error
	}{
		{
			"calling Get with empty slice",
			[]int{},
			0,
			0,
			fmt.Errorf("value not found: empty collection"),
		},
		{
			"calling Get with out of bounds index",
			[]int{1},
			2,
			0,
			fmt.Errorf("value not found: index out of bounds"),
		},
		{
			"calling Get with slice with values",
			[]int{1, 2, 3, 4},
			0,
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			getV := Get(tc.i, tc.sut)
			getEV, err := GetE(tc.i, tc.sut)

			if tc.v != getV || tc.v != getEV {
				t.Errorf("expected get values to be %d and %d. got %d", getV, getEV, tc.v)
			}

			if tc.err != nil {
				if err == nil || err.Error() != tc.err.Error() {
					t.Errorf("expect error to be '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestPush(t *testing.T) {
	expectedPushed := []int{1}
	var sut []int

	if pushed := Push(1, sut); !reflect.DeepEqual(pushed, expectedPushed) {
		t.Errorf("expected slice after push to be %v. got %v", expectedPushed, pushed)
	}
}

func TestPut(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		i           int
		v           string
		expectation []string
	}{
		{
			"putting at 0 on empty slice",
			[]string{},
			0,
			"foo",
			[]string{"foo"},
		},
		{
			"putting at 1 on an empty slice",
			[]string{},
			1,
			"foo",
			[]string{"", "foo"},
		},
		{
			"putting in the middle of a slice",
			[]string{"foo", "baz", "foo", "bar", "baz"},
			1,
			"bar",
			[]string{"foo", "bar", "baz", "foo", "bar", "baz"},
		},
		{
			"prepending to a slice",
			[]string{"bar", "baz"},
			0,
			"foo",
			[]string{"foo", "bar", "baz"},
		},
		{
			"appending to a slice",
			[]string{"foo", "bar"},
			2,
			"baz",
			[]string{"foo", "bar", "baz"},
		},
		{
			"appending to a high cap slice",
			append(make([]string, 0, 10), "foo", "baz"),
			1,
			"bar",
			append(make([]string, 0, 10), "foo", "bar", "baz"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := Put(tc.i, tc.v, tc.sut)

			if !reflect.DeepEqual(actual, tc.expectation) {
				t.Errorf("expected %v. got %v", tc.expectation, actual)
			}

			if len(actual) != len(tc.expectation) {
				t.Errorf("expected %d. got %d", len(tc.expectation), len(actual))
			}
		})
	}
}

func TestPop(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		v           string
		count       int
		capacity    int
	}{
		{
			"popping an empty slice",
			[]string{},
			"",
			0,
			0,
		},
		{
			"popping a slice with items",
			[]string{"foo", "bar"},
			"bar",
			1,
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualV := Pop(&tc.sut)

			if actualV != tc.v {
				t.Errorf("expected %s. got %s", tc.v, actualV)
			}

			if len(tc.sut) != tc.count {
				t.Errorf("expected count after poping to be %d. got %d", tc.count, len(tc.sut))
			}

			if cap(tc.sut) != tc.capacity {
				t.Errorf("expected capacity after poping to be %d. got %d", tc.capacity, cap(tc.sut))
			}
		})
	}
}

func TestPopE(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		v           string
		err         error
		count       int
		capacity    int
	}{
		{
			"popping an empty slice",
			[]string{},
			"",
			fmt.Errorf("value not found: empty collection"),
			0,
			0,
		},
		{
			"popping a slice with items",
			[]string{"foo", "bar"},
			"bar",
			nil,
			1,
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualV, actualErr := PopE(&tc.sut)

			if actualV != tc.v {
				t.Errorf("expected %s. got %s", tc.v, actualV)
			}

			if actualErr != nil && tc.err != nil {
				if actualErr.Error() != tc.err.Error() {
					t.Errorf("expected '%s'. got '%s'", tc.err.Error(), actualErr.Error())
				}
			}

			if len(tc.sut) != tc.count {
				t.Errorf("expected count after poping to be %d. got %d", tc.count, len(tc.sut))
			}

			if cap(tc.sut) != tc.capacity {
				t.Errorf("expected capacity after poping to be %d. got %d", tc.capacity, cap(tc.sut))
			}
		})
	}
}

func TestEach(t *testing.T) {
	var eachResult []any
	sut := []any{1, "foo", 1.1}

	Each(func(_ int, v any) {
		eachResult = append(eachResult, v)
	}, sut)

	if !reflect.DeepEqual(eachResult, sut) {
		t.Errorf("expected visited values to be %v. got %v", sut, eachResult)
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		description string
		sut         []any
		input       any
		i           int
	}{
		{
			"searching with an empty slice",
			[]any{},
			"foo",
			-1,
		},
		{
			"searching an unexisting element",
			[]any{1, "foo", 1.0},
			"bar",
			-1,
		},
		{
			"searching an existing element",
			[]any{1, "foo", 1.0},
			"foo",
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			i := Search(tc.input, tc.sut)

			if i != tc.i {
				t.Errorf("expected resulting index to be %d. got %d", tc.i, i)
			}
		})
	}
}

func TestSearchE(t *testing.T) {
	testCases := []struct {
		description string
		sut         []any
		input       any
		i           int
		err         error
	}{
		{
			"searching with an empty slice",
			[]any{},
			"foo",
			-1,
			fmt.Errorf("value not found"),
		},
		{
			"searching an unexisting element",
			[]any{1, "foo", 1.0},
			"bar",
			-1,
			fmt.Errorf("value not found"),
		},
		{
			"searching an existing element",
			[]any{1, "foo", 1.0},
			"foo",
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			i, err := SearchE(tc.input, tc.sut)

			if i != tc.i {
				t.Errorf("expected resulting index to be %d. got %d", tc.i, i)
			}

			if err != nil {
				if err.Error() != tc.err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestSort(t *testing.T) {
	sut := []int{3, 2, 4, 1}
	sorted := []int{1, 2, 3, 4}

	Sort(sut, Asc[int]())

	if !reflect.DeepEqual(sut, sorted) {
		t.Errorf("expected sorted slice to be %v. got %v", sorted, sut)
	}
}

func TestMap(t *testing.T) {
	testCases := []struct {
		description      string
		sut              []int
		mappedCollection []int
	}{
		{
			"mapping an empty slice",
			[]int{},
			[]int{},
		},
		{
			"mapping a slice with values",
			[]int{1, 2, 3},
			[]int{2, 3, 4},
		},
	}

	f := func(_ int, v int) int {
		return v + 1
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mappedCollection := Map(f, tc.sut)

			if !reflect.DeepEqual(mappedCollection, tc.mappedCollection) {
				t.Errorf(
					"expected mapped collection to be %v. got %v",
					tc.mappedCollection,
					mappedCollection,
				)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	testCases := []struct {
		description string
		input       []int
		fn          func(int, int, int) int
		initial     int
		expected    int
	}{
		{
			"reducing with sum",
			[]int{1, 2, 3, 4, 5},
			func(carry int, n int, _ int) int {
				return carry + n
			},
			0,
			15,
		},
		{
			"reducing with subtraction",
			[]int{1, 2, 3, 4, 5},
			func(carry int, n int, _ int) int {
				return carry - n
			},
			0,
			-15,
		},
		{
			"reducing with multiplication",
			[]int{1, 2, 3, 4, 5},
			func(carry int, n int, _ int) int {
				return carry * n
			},
			1,
			120,
		},
		{
			"reducing with a fixed value",
			[]int{1, 2, 3, 4, 5},
			func(carry int, n int, _ int) int {
				return 2
			},
			1,
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := Reduce(tc.fn, tc.initial, tc.input)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %d. got %d", tc.expected, actual)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		v           string
	}{
		{
			"calling First with an empty slice",
			[]string{},
			"",
		},
		{
			"calling first with a slice with values",
			[]string{"foo", "bar", "baz"},
			"foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v := First(tc.sut)

			if v != tc.v {
				t.Errorf("expected returned value to be '%s', got '%s'", tc.v, v)
			}
		})
	}
}

func TestFirstE(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		v           string
		err         error
	}{
		{
			"calling First with an empty slice",
			[]string{},
			"",
			fmt.Errorf("value not found: empty collection"),
		},
		{
			"calling first with a slice with values",
			[]string{"foo", "bar", "baz"},
			"foo",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v, err := FirstE(tc.sut)

			if v != tc.v {
				t.Errorf("expected returned value to be '%s', got '%s'", tc.v, v)
			}

			if err != nil && tc.err != nil {
				if err.Error() != tc.err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestLast(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		v           string
	}{
		{
			"calling Last with an empty slice",
			[]string{},
			"",
		},
		{
			"calling last with a slice with values",
			[]string{"foo", "bar", "baz"},
			"baz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v := Last(tc.sut)

			if v != tc.v {
				t.Errorf("expected returned value to be '%s', got '%s'", tc.v, v)
			}
		})
	}
}

func TestLastE(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		v           string
		err         error
	}{
		{
			"calling Last with an empty slice",
			[]string{},
			"",
			fmt.Errorf("value not found: empty collection"),
		},
		{
			"calling last with a slice with values",
			[]string{"foo", "bar", "baz"},
			"baz",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v, err := LastE(tc.sut)

			if v != tc.v {
				t.Errorf("expected returned value to be '%s', got '%s'", tc.v, v)
			}

			if err != nil && tc.err != nil {
				if err.Error() != tc.err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestCopy(t *testing.T) {
	from := []string{"foo", "bar", "baz"}

	if to := Copy(from); !reflect.DeepEqual(from, to) {
		t.Errorf("expected copied slice to be %v. got %v", from, to)
	}
}

func TestCut(t *testing.T) {
	testCases := []struct {
		description string
		from        []string
		expected    []string
		remaining   []string
		begin       int
		end         int
	}{
		{
			"cutting an invalid interval",
			[]string{"foo", "bar", "baz"},
			nil,
			[]string{"foo", "bar", "baz"},
			4,
			5,
		},
		{
			"cutting a valid interval",
			[]string{"foo", "bar", "baz"},
			[]string{"bar", "baz"},
			[]string{"foo"},
			1,
			3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := Cut(&tc.from, tc.begin, tc.end)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected cutted slice to be %v. got %v", tc.expected, actual)
			}

			if !reflect.DeepEqual(tc.remaining, tc.from) {
				t.Errorf("expected remaining slice to be %v. got %v", tc.remaining, tc.from)
			}
		})
	}
}

func TestCutE(t *testing.T) {
	testCases := []struct {
		description string
		from        []string
		expected    []string
		remaining   []string
		begin       int
		end         int
		err         error
	}{
		{
			"cutting an invalid interval",
			[]string{"foo", "bar", "baz"},
			nil,
			[]string{"foo", "bar", "baz"},
			4,
			5,
			fmt.Errorf("index out of bounds"),
		},
		{
			"cutting a valid interval",
			[]string{"foo", "bar", "baz"},
			[]string{"bar", "baz"},
			[]string{"foo"},
			1,
			3,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := CutE(&tc.from, tc.begin, tc.end)

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected cutted slice to be %v. got %v", tc.expected, actual)
			}

			if !reflect.DeepEqual(tc.remaining, tc.from) {
				t.Errorf("expected remaining slice to be %v. got %v", tc.remaining, tc.from)
			}

			if tc.err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestForgetE(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		expected    []string
		i           int
		err         error
	}{
		{
			"deleting an unexisting key",
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "bar", "baz"},
			3,
			fmt.Errorf("index out of bounds"),
		},
		{
			"deleting a valid key",
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "baz"},
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := ForgetE(&tc.sut, tc.i)
			if !reflect.DeepEqual(tc.sut, tc.expected) {
				t.Errorf(
					"expected slice after deletting the key to be %v. got %v",
					tc.expected,
					tc.sut,
				)
			}

			if tc.err != nil || err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestDeleteE(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		expected    []string
		i           int
		err         error
	}{
		{
			"deleting an unexisting key",
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "bar", "baz"},
			3,
			fmt.Errorf("index out of bounds"),
		},
		{
			"deleting a valid key",
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "baz"},
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := DeleteE(&tc.sut, tc.i)
			if !reflect.DeepEqual(tc.sut, tc.expected) {
				t.Errorf(
					"expected slice after deletting the key to be %v. got %v",
					tc.expected,
					tc.sut,
				)
			}

			if tc.err != nil || err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}

}

func TestShift(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		remaining   []string
		length      int
		v           string
	}{
		{
			"shifting an empty slice",
			[]string{},
			[]string{},
			0,
			"",
		},
		{
			"shifting a slice with values",
			[]string{"foo", "bar", "baz"},
			[]string{"bar", "baz"},
			2,
			"foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v := Shift(&tc.sut)

			if v != tc.v {
				t.Errorf("expected returned value to be '%s'. got %s", tc.v, v)
			}

			if length := len(tc.sut); length != tc.length {
				t.Errorf("expected sut lenght to be %d. got %d", tc.length, length)
			}
		})
	}
}

func TestShiftE(t *testing.T) {
	testCases := []struct {
		description string
		sut         []string
		remaining   []string
		length      int
		v           string
		err         error
	}{
		{
			"shifting an empty slice",
			[]string{},
			[]string{},
			0,
			"",
			fmt.Errorf("value not found: empty collection"),
		},
		{
			"shifting a slice with values",
			[]string{"foo", "bar", "baz"},
			[]string{"bar", "baz"},
			2,
			"foo",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v, err := ShiftE(&tc.sut)

			if v != tc.v {
				t.Errorf("expected returned value to be '%s'. got %s", tc.v, v)
			}

			if length := len(tc.sut); length != tc.length {
				t.Errorf("expected sut lenght to be %d. got %d", tc.length, length)
			}

			if tc.err != nil || err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestTally(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		expected map[int]int
	}{
		{"empty", []int{}, map[int]int{}},
		{"1,2,3", []int{1, 2, 3}, map[int]int{1: 1, 2: 1, 3: 1}},
		{"1,2,1", []int{1, 2, 1}, map[int]int{1: 2, 2: 1}},
		{"1,3,3,7", []int{1, 3, 3, 7}, map[int]int{1: 1, 3: 2, 7: 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := Tally(tc.slice); !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected '%v'. Got '%v'", tc.expected, actual)
			}
		})
	}
}

func TestMode(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{"empty", []int{}, []int{}},
		{"1,2,3", []int{1, 2, 3}, []int{1, 2, 3}},
		{"1,2,1", []int{1, 2, 1}, []int{1}},
		{"1,3,3,1", []int{1, 3, 3, 1}, []int{1, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := Mode(tc.slice); !reflect.DeepEqual(Tally(actual), Tally(tc.expected)) {
				t.Errorf("expected '%v'. Got '%v'", tc.expected, actual)
			}
		})
	}
}

func TestFirstWhereField(t *testing.T) {
	users := []user{
		{Name: "Jon", Email: "jon@collections.go", Age: 33},
		{Name: "Jane", Email: "jane@collections.go", Age: 27},
		{Name: "Alice", Email: "alice@collections.go", Age: 40},
		{Name: "Bob", Email: "bob@collections.go", Age: 22},
		{Name: "Eve", Email: "eve@collections.go", Age: 30},
	}
	testCases := []struct {
		description string
		sut         []user
		field       string
		matcher     Matcher
		user        user
	}{
		{
			"slice contains a matching object",
			users,
			"Name",
			ValueEquals("Alice"),
			user{Name: "Alice", Email: "alice@collections.go", Age: 40},
		},
		{
			"criteria doesn't match slice elements",
			users,
			"Age",
			ValueGT(50),
			user{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if user := FirstWhereField(tc.sut, tc.field, tc.matcher); user != tc.user {
				t.Errorf("expected user to be %v. got %v", tc.user, user)
			}
		})
	}
}

func TestFirstWhereFieldE(t *testing.T) {
	users := []user{
		{Name: "Jon", Email: "jon@collections.go", Age: 33},
		{Name: "Jane", Email: "jane@collections.go", Age: 27},
		{Name: "Alice", Email: "alice@collections.go", Age: 40},
		{Name: "Bob", Email: "bob@collections.go", Age: 22},
		{Name: "Eve", Email: "eve@collections.go", Age: 30},
	}
	testCases := []struct {
		description string
		sut         []user
		field       string
		matcher     Matcher
		user        user
		err         error
	}{
		{
			"slice contains a matching object",
			users,
			"Name",
			ValueEquals("Alice"),
			user{Name: "Alice", Email: "alice@collections.go", Age: 40},
			nil,
		},
		{
			"slice contains an object matching custom matcher",
			users,
			"Age",
			ValueGT(30),
			user{Name: "Jon", Email: "jon@collections.go", Age: 33},
			nil,
		},
		{
			"criteria doesn't match slice elements",
			users,
			"Age",
			ValueGT(50),
			user{},
			fmt.Errorf("value not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			user, err := FirstWhereFieldE(tc.sut, tc.field, tc.matcher)

			if user != tc.user {
				t.Errorf("expected user to be %v. got %v", tc.user, user)
			}

			if tc.err != nil || err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got %s", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestFirstWhereE(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	testCases := []struct {
		description string
		sut         []int
		matcher     Matcher
		v           int
		err         error
	}{
		{
			"slice contains a matching value",
			slice,
			ValueGT(3),
			4,
			nil,
		},
		{
			"slice does not contain matching values",
			slice,
			ValueGT(5),
			*new(int),
			fmt.Errorf("value not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v, err := FirstWhereE(tc.sut, tc.matcher)

			if v != tc.v {
				t.Errorf("expected returned value to be %d. got %d", tc.v, v)
			}

			if tc.err != nil || err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got %s", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestFirstWhereWithComposedMatchers(t *testing.T) {
	users := []user{
		{Name: "Jon", Email: "jon@collections.go", Age: 33},
		{Name: "Jane", Email: "jane@collections.go", Age: 27},
		{Name: "Alice", Email: "alice@collections.go", Age: 40},
		{Name: "Bob", Email: "bob@collections.go", Age: 22},
		{Name: "Eve", Email: "eve@collections.go", Age: 30},
	}
	testCases := []struct {
		description string
		sut         []user
		field       string
		matcher     Matcher
		user        user
		err         error
	}{
		{
			"slice contains a matching object composed with 'AndValue'",
			users,
			"Name",
			AndValue(user{Name: "Alice", Age: 40}, usernameMatch, ageMatch),
			user{Name: "Alice", Email: "alice@collections.go", Age: 40},
			nil,
		},
		{
			"slice contains a matching object composed with 'OrValue'",
			users,
			"Name",
			OrValue(user{Name: "Alice", Age: 33}, usernameMatch, ageMatch),
			user{Name: "Jon", Email: "jon@collections.go", Age: 33},
			nil,
		},
		{
			"slice does not contain matching objects",
			users,
			"Name",
			AndValue(user{Name: "Alice", Age: 33}, usernameMatch, ageMatch),
			user{},
			fmt.Errorf("value not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			user, err := FirstWhereE(tc.sut, tc.matcher)

			if user != tc.user {
				t.Errorf("expected user to be %v. got %v", tc.user, user)
			}

			if tc.err != nil || err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got %s", tc.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestFirstWhere(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	testCases := []struct {
		description string
		sut         []int
		matcher     Matcher
		v           int
	}{
		{
			"slice contains a matching value",
			slice,
			ValueGT(3),
			4,
		},
		{
			"slice does not contain matching values",
			slice,
			ValueGT(5),
			*new(int),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if v := FirstWhere(tc.sut, tc.matcher); v != tc.v {
				t.Errorf("expected returned value to be %d. got %d", tc.v, v)
			}
		})
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		description string
		slice       []int
		matcher     Matcher
		contains    bool
	}{
		{
			"collection contains at least one matching value",
			[]int{1, 2, 3, 4},
			ValueEquals(3),
			true,
		},
		{
			"collection does not contain matching values",
			[]int{1, 2, 3, 4},
			ValueEquals(5),
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if contains := Contains(tc.slice, tc.matcher); contains != tc.contains {
				t.Errorf("Contains result should be  %v. got %v", tc.contains, contains)
			}
		})
	}
}

func TestDuplicates(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		expected []string
	}{
		{
			"no duplicates",
			[]string{"1", "2", "3", "4"},
			[]string{},
		},
		{
			"1 appearing twice",
			[]string{"1", "2", "1", "3", "4"},
			[]string{"1"},
		},
		{
			"1 and 2 appearing twice",
			[]string{"1", "2", "1", "3", "2"},
			[]string{"1", "2"},
		},
		{
			"every element appearing twice",
			[]string{"1", "2", "3", "1", "2", "3"},
			[]string{"1", "2", "3"},
		},
		{
			"every element appearing thrice",
			[]string{"1", "2", "3", "1", "2", "3", "1", "2", "3"},
			[]string{"1", "2", "3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Duplicates(tc.slice); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

type TestCase[T any] struct {
	Name          string
	ReceiverSlice []T
	DiffSlice     []T
	Expected      []T
}

func TestDiffWithInteger(t *testing.T) {

	integerCases := []TestCase[int]{
		{
			"ordered values",
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3},
			[]int{4, 5},
		},
		{
			"unordored values",
			[]int{5, 4, 3, 2, 1},
			[]int{3, 2, 1},
			[]int{5, 4},
		},
		{
			"values repeats",
			[]int{1, 2, 2, 2, 2},
			[]int{2},
			[]int{1},
		},
	}

	for _, tc := range integerCases {
		t.Run(tc.Name, func(t *testing.T) {
			if got := Diff(tc.ReceiverSlice, tc.DiffSlice); !reflect.DeepEqual(got, tc.Expected) {
				t.Errorf("Expected '%v' got '%v' instead", tc.Expected, got)
			}
		})
	}

}

func TestDiffWithString(t *testing.T) {
	stringCases := []TestCase[string]{
		{
			"small case values",
			[]string{"foo", "bar"},
			[]string{"foo"},
			[]string{"bar"},
		},
		{
			"upper case values",
			[]string{"FOO", "BAR"},
			[]string{"foo"},
			[]string{"FOO", "BAR"},
		},
	}

	for _, tc := range stringCases {
		t.Run(tc.Name, func(t *testing.T) {
			if got := Diff(tc.ReceiverSlice, tc.DiffSlice); !reflect.DeepEqual(got, tc.Expected) {
				t.Errorf("Expected '%v' got '%v' instead", tc.Expected, got)
			}
		})
	}
}
