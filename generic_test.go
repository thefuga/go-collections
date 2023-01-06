package collections

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/thefuga/go-collections/errors"
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
				t.Errorf("expected count after popping to be %d. got %d", tc.count, len(tc.sut))
			}

			if cap(tc.sut) != tc.capacity {
				t.Errorf("expected capacity after popping to be %d. got %d", tc.capacity, cap(tc.sut))
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
				t.Errorf("expected count after popping to be %d. got %d", tc.count, len(tc.sut))
			}

			if cap(tc.sut) != tc.capacity {
				t.Errorf("expected capacity after popping to be %d. got %d", tc.capacity, cap(tc.sut))
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
			"searching a nonexisting element",
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
			"searching a nonexisting element",
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

func TestMappingToADifferentType(t *testing.T) {
	slice := []int{1, 2, 3}
	mapper := func(_ int, n int) string {
		return strconv.Itoa(n)
	}
	expected := []string{"1", "2", "3"}

	if got := Map(mapper, slice); !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected '%v'. Got '%v'", expected, got)
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

	if !reflect.DeepEqual(from, []string{"foo", "bar", "baz"}) {
		t.Errorf("copy must not change source")
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
				t.Errorf("expected cut slice to be %v. got %v", tc.expected, actual)
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
				t.Errorf("expected cut slice to be %v. got %v", tc.expected, actual)
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
			"deleting a nonexisting key",
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
					"expected slice after deleting the key to be %v. got %v",
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
			"deleting a nonexisting key",
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
					"expected slice after deleting the key to be %v. got %v",
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
				t.Errorf("expected sut length to be %d. got %d", tc.length, length)
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
				t.Errorf("expected sut length to be %d. got %d", tc.length, length)
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
		matcher     AnyMatcher
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
		matcher     AnyMatcher
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
		matcher     AnyMatcher
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
		matcher     AnyMatcher
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
		matcher     AnyMatcher
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
		matcher     AnyMatcher
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

func TestZip(t *testing.T) {
	testCases := []struct {
		name     string
		slices   [][]int
		expected [][]int
	}{
		{
			"no slices",
			[][]int{},
			[][]int{},
		},
		{
			"same length",
			[][]int{{1, 3, 5}, {2, 4, 6}},
			[][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			"first slice shorter",
			[][]int{{1, 3}, {2, 4, 6}},
			[][]int{{1, 2}, {3, 4}},
		},
		{
			"second slice shorter",
			[][]int{{1, 3, 5}, {2, 4}},
			[][]int{{1, 2}, {3, 4}},
		},
		{
			"first slice empty",
			[][]int{{}, {1, 2, 3}},
			[][]int{},
		},
		{
			"second slice empty",
			[][]int{{1, 2, 3}, {}},
			[][]int{},
		},
		{
			"multiple slices",
			[][]int{{1}, {2}, {3}, {4}, {5}},
			[][]int{{1, 2, 3, 4, 5}},
		},
		{
			"one slice",
			[][]int{{1, 2, 3, 4, 5}},
			[][]int{{1}, {2}, {3}, {4}, {5}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Zip(tc.slices...); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestZippingTwiceReturnsTheOriginalInput(t *testing.T) {
	input := [][]int{{1, 2}, {3, 4}}
	if got := Zip(Zip(input...)...); !reflect.DeepEqual(input, got) {
		t.Errorf("Expected '%v'. Got '%v'", input, got)
	}
}

func TestUnique(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			"no duplicate",
			[]int{1, 2, 3},
			[]int{1, 2, 3},
		},
		{
			"first value duplicate",
			[]int{1, 1, 2, 3},
			[]int{1, 2, 3},
		},
		{
			"duplicate in the end",
			[]int{1, 2, 3, 1},
			[]int{1, 2, 3},
		},
		{
			"every element is duplicate",
			[]int{1, 1, 2, 2, 3, 3},
			[]int{1, 2, 3},
		},
		{
			"every element is duplicate in reverse order",
			[]int{3, 3, 2, 2, 1, 1},
			[]int{3, 2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Unique(tc.input); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestUniqueBy(t *testing.T) {
	type person struct {
		firstName string
		lastName  string
	}

	bobRoss := person{"bob", "ross"}
	alyssaHacker := person{"alyssa", "hacker"}
	bobMartin := person{"bob", "martin"}

	testCases := []struct {
		name     string
		input    []person
		f        func(p person) string
		expected []person
	}{
		{
			"first name",
			[]person{bobRoss, alyssaHacker, bobMartin},
			func(p person) string { return p.firstName },
			[]person{bobRoss, alyssaHacker},
		},
		{
			"last name",
			[]person{bobRoss, alyssaHacker, bobMartin},
			func(p person) string { return p.lastName },
			[]person{bobRoss, alyssaHacker, bobMartin},
		},
		{
			"fixed value",
			[]person{bobRoss, alyssaHacker, bobMartin},
			func(_ person) string { return "foo" },
			[]person{bobRoss},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := UniqueBy(tc.input, tc.f); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestGroupByPredicate(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		f        func(i int) bool
		expected map[bool][]int
	}{
		{
			"values greater than 2",
			[]int{1, 2, 3, 4},
			func(i int) bool { return i > 2 },
			map[bool][]int{false: {1, 2}, true: {3, 4}},
		},
		{
			"group evens",
			[]int{1, 2, 3, 4},
			func(i int) bool { return i%2 == 0 },
			map[bool][]int{true: {2, 4}, false: {1, 3}},
		},
		{
			"group odds",
			[]int{1, 2, 3, 4},
			func(i int) bool { return i%2 == 1 },
			map[bool][]int{false: {2, 4}, true: {1, 3}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := GroupBy(tc.input, tc.f); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestGroupByProperty(t *testing.T) {
	type person struct {
		firstName string
		lastName  string
	}
	firstName := func(p person) string { return p.firstName }

	bobRoss := person{"bob", "ross"}
	alyssaHacker := person{"alyssa", "hacker"}
	bobMartin := person{"bob", "martin"}

	people := []person{bobRoss, alyssaHacker, bobMartin}

	expected := map[string][]person{
		"bob":    {bobRoss, bobMartin},
		"alyssa": {alyssaHacker},
	}

	if got := GroupBy(people, firstName); !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected '%v'. Got '%v'", expected, got)
	}
}

func TestPartition(t *testing.T) {
	testCases := []struct {
		name          string
		input         []int
		f             func(i int) bool
		expectedLeft  []int
		expectedRight []int
	}{
		{
			name:          "values greater less than or equal to 2",
			input:         []int{1, 2, 3, 4},
			f:             func(i int) bool { return i <= 2 },
			expectedLeft:  []int{1, 2},
			expectedRight: []int{3, 4},
		},
		{
			name:          "even and odds",
			input:         []int{1, 2, 3, 4},
			f:             func(i int) bool { return i%2 == 0 },
			expectedLeft:  []int{2, 4},
			expectedRight: []int{1, 3},
		},
		{
			name:          "all true",
			input:         []int{1, 2, 3, 4},
			f:             func(_ int) bool { return true },
			expectedLeft:  []int{1, 2, 3, 4},
			expectedRight: []int{},
		},
		{
			name:          "all false",
			input:         []int{1, 2, 3, 4},
			f:             func(_ int) bool { return false },
			expectedLeft:  []int{},
			expectedRight: []int{1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			left, right := Partition(tc.input, tc.f)
			if !reflect.DeepEqual(left, tc.expectedLeft) {
				t.Errorf("Expected left '%v'. Got '%v'", tc.expectedLeft, left)
			}
			if !reflect.DeepEqual(left, tc.expectedLeft) {
				t.Errorf("Expected right '%v'. Got '%v'", tc.expectedRight, right)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		size     int
		expected [][]int
	}{
		{
			name:     "no elements left",
			input:    []int{1, 2, 3, 4, 5, 6},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:     "elements left",
			input:    []int{1, 2, 3, 4, 5},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:     "one big chunk",
			input:    []int{1, 2, 3, 4, 5, 6},
			size:     6,
			expected: [][]int{{1, 2, 3, 4, 5, 6}},
		},
		{
			name:     "chunks of one",
			input:    []int{1, 2, 3, 4, 5, 6},
			size:     1,
			expected: [][]int{{1}, {2}, {3}, {4}, {5}, {6}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Chunk(tc.input, tc.size); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "even number of elements",
			input:    []int{1, 2, 3, 4},
			expected: []int{4, 3, 2, 1},
		},
		{
			name:     "odd number of elements",
			input:    []int{1, 2, 3},
			expected: []int{3, 2, 1},
		},
		{
			name:     "empty",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "one element",
			input:    []int{1},
			expected: []int{1},
		},
		{
			name:     "two elements",
			input:    []int{1, 2},
			expected: []int{2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Reverse(tc.input); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestSumByInt(t *testing.T) {
	type item struct {
		a int
		b int
	}

	i1 := item{1, 2}
	i2 := item{3, 4}
	i3 := item{5, 6}

	testCases := []struct {
		name     string
		input    []item
		f        func(i item) int
		expected int
	}{
		{
			name:     "sum by property a",
			input:    []item{i1, i2, i3},
			f:        func(i item) int { return i.a },
			expected: i1.a + i2.a + i3.a,
		},
		{
			name:     "sum by property b",
			input:    []item{i1, i2, i3},
			f:        func(i item) int { return i.b },
			expected: i1.b + i2.b + i3.b,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SumBy(tc.input, tc.f); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%d'. Got '%d'", tc.expected, got)
			}
		})
	}
}

func TestSumByFloat(t *testing.T) {
	type item struct {
		a float64
		b float64
	}

	i1 := item{1.0, 2.0}
	i2 := item{3.0, 4.0}
	i3 := item{5.0, 6.0}

	testCases := []struct {
		name     string
		input    []item
		f        func(i item) float64
		expected float64
	}{
		{
			name:     "sum by property a",
			input:    []item{i1, i2, i3},
			f:        func(i item) float64 { return i.a },
			expected: i1.a + i2.a + i3.a,
		},
		{
			name:     "sum by property b",
			input:    []item{i1, i2, i3},
			f:        func(i item) float64 { return i.b },
			expected: i1.b + i2.b + i3.b,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SumBy(tc.input, tc.f); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%f'. Got '%f'", tc.expected, got)
			}
		})
	}
}

func TestRange(t *testing.T) {
	testCases := []struct {
		name     string
		min      int
		max      int
		expected []int
	}{
		{
			name:     "is inclusive",
			min:      0,
			max:      1,
			expected: []int{0, 1},
		},
		{
			name:     "returns a slice with a single element when min = max",
			min:      1,
			max:      1,
			expected: []int{1},
		},
		{
			name:     "returns an empty slice when min > max",
			min:      1,
			max:      0,
			expected: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Range(tc.min, tc.max); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestInterpose(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		sep      string
		expected []string
	}{
		{
			name:     "with even slice length",
			slice:    []string{"1", "2", "3", "4"},
			sep:      ",",
			expected: []string{"1", ",", "2", ",", "3", ",", "4"},
		},
		{
			name:     "with odd slice length",
			slice:    []string{"1", "2", "3"},
			sep:      ",",
			expected: []string{"1", ",", "2", ",", "3"},
		},
		{
			name:     "empty slice",
			slice:    []string{},
			sep:      ",",
			expected: []string{},
		},
		{
			name:     "one element",
			slice:    []string{"a"},
			sep:      ",",
			expected: []string{"a"},
		},
		{
			name:     "two elements",
			slice:    []string{"a", "b"},
			sep:      ",",
			expected: []string{"a", ",", "b"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Interpose(tc.slice, tc.sep); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestForPage(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		page     int
		size     int
		expected []int
	}{
		{
			name:     "1 through 9, page 1 with size 3",
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			page:     1,
			size:     3,
			expected: []int{1, 2, 3},
		},
		{
			name:     "1 through 9, page 2 with size 3",
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			page:     2,
			size:     3,
			expected: []int{4, 5, 6},
		},
		{
			name:     "1 through 9, page 3 with size 3",
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			page:     3,
			size:     3,
			expected: []int{7, 8, 9},
		},
		{
			name:     "page < 1 is equal to page = 1",
			slice:    []int{1, 2, 3},
			page:     0,
			size:     3,
			expected: []int{1, 2, 3},
		},
		{
			name:     "page * size > len(slice)",
			slice:    []int{1, 2, 3},
			page:     2,
			size:     3,
			expected: []int{},
		},
		{
			name:     "size = len(slice)",
			slice:    []int{1, 2, 3},
			page:     1,
			size:     3,
			expected: []int{1, 2, 3},
		},
		{
			name:     "size > len(slice)",
			slice:    []int{1, 2, 3},
			page:     1,
			size:     10,
			expected: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ForPage(tc.slice, tc.page, tc.size); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestKeyBy(t *testing.T) {
	type person struct {
		firstName string
		lastName  string
	}

	bobRoss := person{"bob", "ross"}
	alyssaHacker := person{"alyssa", "hacker"}
	bobMartin := person{"bob", "martin"}

	testCases := []struct {
		name     string
		input    []person
		f        func(p person) string
		expected map[string]person
	}{
		{
			name:     "first name",
			input:    []person{bobRoss, alyssaHacker, bobMartin},
			f:        func(p person) string { return p.firstName },
			expected: map[string]person{"bob": bobMartin, "alyssa": alyssaHacker},
		},
		{
			name:  "last name",
			input: []person{bobRoss, alyssaHacker, bobMartin},
			f:     func(p person) string { return p.lastName },
			expected: map[string]person{
				"ross":   bobRoss,
				"hacker": alyssaHacker,
				"martin": bobMartin,
			},
		},
		{
			name:     "fixed value will return the last value",
			input:    []person{bobRoss, alyssaHacker, bobMartin},
			f:        func(p person) string { return "foo" },
			expected: map[string]person{"foo": bobMartin},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := KeyBy(tc.input, tc.f); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestLastBy(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		matcher  Matcher[int, int]
		expected int
	}{
		{
			name:     "last even number",
			slice:    []int{1, 2, 3, 4},
			matcher:  func(_, n int) bool { return n%2 == 0 },
			expected: 4,
		},
		{
			name:     "last odd number",
			slice:    []int{1, 2, 3, 4},
			matcher:  func(_, n int) bool { return n%2 == 1 },
			expected: 3,
		},
		{
			name:     "return zero value when no element matches",
			slice:    []int{1, 2, 3, 4},
			matcher:  func(_, _ int) bool { return false },
			expected: 0,
		},
		{
			name:     "empty slice with true matcher",
			slice:    []int{},
			matcher:  func(_, _ int) bool { return true },
			expected: 0,
		},
		{
			name:     "empty slice with false matcher",
			slice:    []int{},
			matcher:  func(_, _ int) bool { return false },
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := LastBy(tc.slice, tc.matcher); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestLastByE(t *testing.T) {
	testCases := []struct {
		name        string
		slice       []int
		matcher     Matcher[int, int]
		expected    int
		expectedErr error
	}{
		{
			name:        "last even with odd numbers",
			slice:       []int{1, 3},
			matcher:     func(_, n int) bool { return n%2 == 0 },
			expected:    0,
			expectedErr: errors.NewValueNotFoundError(),
		},
		{
			name:        "last odd with even numbers",
			slice:       []int{2, 4},
			matcher:     func(_, n int) bool { return n%2 == 1 },
			expected:    0,
			expectedErr: errors.NewValueNotFoundError(),
		},
		{
			name:        "empty slice with false matcher",
			slice:       []int{},
			matcher:     func(_, n int) bool { return false },
			expected:    0,
			expectedErr: errors.NewValueNotFoundError(),
		},
		{
			name:        "empty slice with true matcher",
			slice:       []int{},
			matcher:     func(_, n int) bool { return true },
			expected:    0,
			expectedErr: errors.NewValueNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LastByE(tc.slice, tc.matcher)

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}

			if tc.expectedErr == nil {
				return
			}

			if err == nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Expected error '%v'. Got error '%v'", tc.expectedErr, err)
			}
		})
	}
}

func TestPad(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		size     int
		pad      int
		expected []int
	}{
		{
			name:     "fill to 5 with 1s",
			slice:    []int{1, 2, 3},
			size:     5,
			pad:      1,
			expected: []int{1, 2, 3, 1, 1},
		},
		{
			name:     "fill empty slice to 5 with 5s",
			slice:    []int{},
			size:     5,
			pad:      5,
			expected: []int{5, 5, 5, 5, 5},
		},
		{
			name:     "when size = len(slice), it returns the slice",
			slice:    []int{1, 2, 3},
			size:     3,
			pad:      4,
			expected: []int{1, 2, 3},
		},
		{
			name:     "when size < len(slice), it returns the slice",
			slice:    []int{1, 2, 3},
			size:     1,
			pad:      4,
			expected: []int{1, 2, 3},
		},

		// negative cases (pad left)
		{
			name:     "fill to 5 with 1s",
			slice:    []int{1, 2, 3},
			size:     -5,
			pad:      3,
			expected: []int{3, 3, 1, 2, 3},
		},
		{
			name:     "fill empty slice to 5 with 5s",
			slice:    []int{},
			size:     -5,
			pad:      5,
			expected: []int{5, 5, 5, 5, 5},
		},
		{
			name:     "when size = len(slice), it returns the slice",
			slice:    []int{1, 2, 3},
			size:     -3,
			pad:      4,
			expected: []int{1, 2, 3},
		},
		{
			name:     "when size < len(slice), it returns the slice",
			slice:    []int{1, 2, 3},
			size:     -1,
			pad:      4,
			expected: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if got := Pad(tc.slice, tc.size, tc.pad); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestPrepend(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		value    int
		expected []int
	}{
		{
			name:     "prepend 0 to [1,2,3]",
			slice:    []int{1, 2, 3},
			value:    0,
			expected: []int{0, 1, 2, 3},
		},
		{
			name:     "prepend to empty slice",
			slice:    []int{},
			value:    1,
			expected: []int{1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if got := Prepend(tc.slice, tc.value); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestRandomResultIsIncludedInSlice(t *testing.T) {
	slice := Range(1, 10)

	for i := 0; i < 10; i++ {
		rand := Random(slice)
		if !Contains(slice, ValueEquals(rand)) {
			t.Errorf("expected slice '%v' to contain '%v'", slice, rand)
		}
	}
}

func TestRandomSampling(t *testing.T) {
	slice := Range(1, 1000)

	sampleA := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		sampleA = append(sampleA, Random(slice))
	}

	sampleB := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		sampleB = append(sampleB, Random(slice))
	}

	if reflect.DeepEqual(sampleA, sampleB) {
		t.Errorf("got equal samples sampleA='%v' sampleB='%v'", sampleA, sampleB)
	}
}

func TestRandomWithEmptySliceReturnsZeroValue(t *testing.T) {
	if got := Random([]int{}); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}

	if got := Random([]string{}); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}

	if got := Random([]bool{}); got != false {
		t.Errorf("expected false, got %t", got)
	}
}

func TestRandomEErrorsWhenSliceIsEmpty(t *testing.T) {
	expected := errors.NewEmptyCollectionError()

	if _, got := RandomE([]int{}); got.Error() != expected.Error() {
		t.Errorf("Expected %q. Got %q", expected, got)
	}
}

func TestShuffleChangesTheSlice(t *testing.T) {
	slice := Range(1, 10)
	copyOfSlice := Copy(slice)

	got := Shuffle(slice)
	if reflect.DeepEqual(got, copyOfSlice) {
		t.Errorf("shuffle did not change the slice")
	}

	if !reflect.DeepEqual(slice, got) {
		t.Errorf("Expected '%v'. Got '%v'", got, slice)
	}
}

func TestSkip(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		skip     int
		expected []int
	}{
		{
			name:     "skipping first element on a slice of length 3",
			slice:    []int{1, 2, 3},
			skip:     1,
			expected: []int{2, 3},
		},
		{
			name:     "skipping first element on a slice of length 1",
			slice:    []int{1},
			skip:     1,
			expected: []int{},
		},
		{
			name:     "when skip = len(slice)",
			slice:    []int{1, 2, 3},
			skip:     3,
			expected: []int{},
		},
		{
			name:     "when skip > len(slice)",
			slice:    []int{1, 2, 3},
			skip:     10,
			expected: []int{},
		},
		{
			name:     "skipping on an empty slice",
			slice:    []int{},
			skip:     10,
			expected: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Skip(tc.slice, tc.skip); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestSkipUntil(t *testing.T) {
	alwaysTrueMatcher := func(_, _ any) bool { return true }
	alwaysFalseMatcher := func(_, _ any) bool { return false }

	testCases := []struct {
		name     string
		slice    []int
		matcher  AnyMatcher
		expected []int
	}{
		{
			name:     "after 3",
			slice:    []int{1, 2, 3, 4, 5},
			matcher:  ValueEquals(3),
			expected: []int{3, 4, 5},
		},
		{
			name:     "returns empty if matcher does not match",
			slice:    []int{1, 2, 3, 4, 5},
			matcher:  alwaysFalseMatcher,
			expected: []int{},
		},
		{
			name:     "returns the same slice if matcher matches the first element",
			slice:    []int{1, 2, 3, 4, 5},
			matcher:  alwaysTrueMatcher,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "not matching on an empty slice returns an empty slice",
			slice:    []int{},
			matcher:  alwaysFalseMatcher,
			expected: []int{},
		},
		{
			name:     "matching on an empty slice returns an empty slice",
			slice:    []int{},
			matcher:  alwaysTrueMatcher,
			expected: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SkipUntil(tc.slice, tc.matcher); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestSkipWhile(t *testing.T) {
	alwaysTrueMatcher := func(_, _ any) bool { return true }
	alwaysFalseMatcher := func(_, _ any) bool { return false }

	testCases := []struct {
		name     string
		slice    []int
		matcher  AnyMatcher
		expected []int
	}{
		{
			name:     "skip while value is less than 3",
			slice:    []int{1, 2, 3, 4, 5},
			matcher:  ValueLT(3),
			expected: []int{3, 4, 5},
		},
		{
			name:     "returns the same slice if matcher does not match",
			slice:    []int{1, 2, 3, 4, 5},
			matcher:  alwaysFalseMatcher,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "returns empty slice if matcher matches the first element",
			slice:    []int{1, 2, 3, 4, 5},
			matcher:  alwaysTrueMatcher,
			expected: []int{},
		},
		{
			name:     "not matching on an empty slice returns an empty slice",
			slice:    []int{},
			matcher:  alwaysFalseMatcher,
			expected: []int{},
		},
		{
			name:     "matching on an empty slice returns an empty slice",
			slice:    []int{},
			matcher:  alwaysTrueMatcher,
			expected: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SkipWhile(tc.slice, tc.matcher); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestNth(t *testing.T) {
	testCases := []struct {
		name     string
		nth      int
		slice    []rune
		expected []rune
	}{
		{
			name:     "slice's len is divisible by nth",
			nth:      2,
			slice:    []rune{'a', 'b', 'c', 'd', 'e', 'f'},
			expected: []rune{'a', 'c', 'e'},
		},
		{
			name:     "slice's len is not divisible by nth",
			nth:      4,
			slice:    []rune{'a', 'b', 'c', 'd', 'e', 'f'},
			expected: []rune{'a', 'e'},
		},
		{
			name:     "nth is odd",
			nth:      3,
			slice:    []rune{'a', 'b', 'c', 'd', 'e', 'f'},
			expected: []rune{'a', 'd'},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Nth(tc.slice, tc.nth)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}

			lenGot := len(got)
			capGot := cap(got)
			lenExpec := len(tc.expected)
			capExpec := cap(tc.expected)

			if lenGot != lenExpec {
				t.Errorf("expected cap '%d' len. Got '%d'", lenExpec, lenGot)
			}

			if capGot != capExpec {
				t.Errorf("expected cap '%d' len. Got '%d'", capExpec, capGot)
			}
		})
	}
}

func TestNthOffset(t *testing.T) {
	testCases := []struct {
		name     string
		nth      int
		offset   int
		slice    []rune
		expected []rune
	}{
		{
			name:     "slice's len is divisible by nth",
			nth:      2,
			offset:   1,
			slice:    []rune{'a', 'b', 'c', 'd', 'e', 'f'},
			expected: []rune{'b', 'd', 'f'},
		},
		{
			name:     "slice's len is not divisible by nth",
			nth:      4,
			offset:   1,
			slice:    []rune{'a', 'b', 'c', 'd', 'e', 'f'},
			expected: []rune{'b', 'f'},
		},
		{
			name:     "nth is odd",
			nth:      3,
			offset:   1,
			slice:    []rune{'a', 'b', 'c', 'd', 'e', 'f'},
			expected: []rune{'b', 'e'},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NthOffset(tc.slice, tc.nth, tc.offset); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestSliding(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		window   int
		expected [][]int
	}{
		{
			name:     "1 through 5 with a window of 0",
			slice:    []int{1, 2, 3, 4, 5},
			window:   0,
			expected: nil,
		},
		{
			name:     "1 through 5 with a window of 1",
			slice:    []int{1, 2, 3, 4, 5},
			window:   1,
			expected: [][]int{{1}, {2}, {3}, {4}, {5}},
		},
		{
			name:     "1 through 5 with a window of 2",
			slice:    []int{1, 2, 3, 4, 5},
			window:   2,
			expected: [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
		},
		{
			name:     "1 through 5 with a window of 3",
			slice:    []int{1, 2, 3, 4, 5},
			window:   3,
			expected: [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}},
		},
		{
			name:     "1 through 5 with a window of 4",
			slice:    []int{1, 2, 3, 4, 5},
			window:   4,
			expected: [][]int{{1, 2, 3, 4}, {2, 3, 4, 5}},
		},
		{
			name:     "1 through 5 with a window of 5",
			slice:    []int{1, 2, 3, 4, 5},
			window:   5,
			expected: [][]int{{1, 2, 3, 4, 5}},
		},
		{
			name:     "window > len(slice)",
			slice:    []int{1, 2, 3, 4, 5},
			window:   10,
			expected: [][]int{{1, 2, 3, 4, 5}},
		},
		{
			name:     "empty slice with a window o 0",
			slice:    []int{},
			window:   0,
			expected: nil,
		},
		{
			name:     "empty slice with a window o 1",
			slice:    []int{},
			window:   1,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Sliding(tc.slice, tc.window); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}

func TestSlidingStep(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		window   int
		step     int
		expected [][]int
	}{
		{
			name:     "1 through 5 with a window of 0 and step 1",
			slice:    []int{1, 2, 3, 4, 5},
			window:   0,
			step:     1,
			expected: nil,
		},
		{
			name:     "1 through 5 with a window of 1 and step 2",
			slice:    []int{1, 2, 3, 4, 5},
			window:   1,
			step:     2,
			expected: [][]int{{1}, {3}, {5}},
		},
		{
			name:     "1 through 5 with a window of 1 and step 3",
			slice:    []int{1, 2, 3, 4, 5},
			window:   1,
			step:     3,
			expected: [][]int{{1}, {4}},
		},
		{
			name:     "1 through 5 with a window of 1 and step 4",
			slice:    []int{1, 2, 3, 4, 5},
			window:   1,
			step:     4,
			expected: [][]int{{1}, {5}},
		},
		{
			name:     "1 through 5 with a window of 3 and step 2",
			slice:    []int{1, 2, 3, 4, 5},
			window:   3,
			step:     2,
			expected: [][]int{{1, 2, 3}, {3, 4, 5}},
		},
		{
			name:     "step > len(slice)",
			slice:    []int{1, 2, 3, 4, 5},
			window:   1,
			step:     6,
			expected: [][]int{{1}},
		},
		{
			name:     "step = 0",
			slice:    []int{1, 2, 3, 4, 5},
			window:   1,
			step:     0,
			expected: nil,
		},
		{
			name:     "step = -1",
			slice:    []int{1, 2, 3, 4, 5},
			window:   1,
			step:     -1,
			expected: nil,
		},
		{
			name:     "window > len(slice)",
			slice:    []int{1, 2, 3, 4, 5},
			window:   10,
			step:     1,
			expected: [][]int{{1, 2, 3, 4, 5}},
		},
		{
			name:     "empty slice with a window o 0",
			slice:    []int{},
			window:   0,
			step:     1,
			expected: nil,
		},
		{
			name:     "empty slice with a window o 1",
			slice:    []int{},
			window:   1,
			step:     1,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SlidingStep(tc.slice, tc.window, tc.step); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected '%v'. Got '%v'", tc.expected, got)
			}
		})
	}
}
