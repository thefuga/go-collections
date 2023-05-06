package slice

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/thefuga/go-collections"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[string]
		i           int
		v           string
		err         error
	}{
		{
			"empty collection",
			Collection[string]{},
			0,
			"",
			fmt.Errorf("value not found: empty collection"),
		},
		{
			"calling Get with out of bounds index",
			Collect("foo"),
			2,
			"",
			fmt.Errorf("value not found: index out of bounds"),
		},
		{
			"calling Get on a collection with values",
			Collect("foo", "bar", "baz"),
			1,
			"bar",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			getV := tc.sut.Get(1)
			getVE, err := tc.sut.GetE(1)

			if !reflect.DeepEqual(getV, getVE) && !reflect.DeepEqual(getV, tc.v) {
				t.Errorf("expected get value to be '%s'. got '%s'", getV, tc.v)

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
	testCases := []struct {
		description string
		input       Collection[string]
		expectation Collection[string]
	}{
		{
			"empty collection",
			Collection[string]{},
			Collection[string]{"foo"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if !reflect.DeepEqual(tc.input.Push("foo"), tc.expectation) {
				t.Error("")
			}
		})
	}
}

func TestPut(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[string]
		i           int
		v           string
		expectation Collection[string]
	}{
		{
			"putting at 0 on empty collection",
			Collection[string]{},
			0,
			"foo",
			Collection[string]{"foo"},
		},
		{
			"putting at 1 on empty collection",
			Collection[string]{},
			1,
			"foo",
			Collection[string]{"", "foo"},
		},
		{
			"putting in the middle of a collection",
			Collection[string]{"foo", "baz", "foo", "bar", "baz"},
			1,
			"bar",
			Collection[string]{"foo", "bar", "baz", "foo", "bar", "baz"},
		},
		{
			"prepending to a collection",
			Collection[string]{"bar", "baz"},
			0,
			"foo",
			Collection[string]{"foo", "bar", "baz"},
		},
		{
			"appending to a collection",
			Collection[string]{"foo", "bar"},
			2,
			"baz",
			Collection[string]{"foo", "bar", "baz"},
		},
		{
			"appending to a high cap collection",
			append(make(Collection[string], 0, 10), "foo", "baz"),
			1,
			"bar",
			append(make(Collection[string], 0, 10), "foo", "bar", "baz"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.sut.Put(tc.i, tc.v)

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
		sut         Collection[string]
		v           string
		count       int
		capacity    int
	}{
		{
			"popping an empty collection",
			Collection[string]{},
			"",
			0,
			0,
		},
		{
			"popping a collection with items",
			Collection[string]{"foo", "bar"},
			"bar",
			1,
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualV := tc.sut.Pop()

			if actualV != tc.v {
				t.Errorf("expected %s. got %s", tc.v, actualV)
			}

			if tc.sut.Count() != tc.count {
				t.Errorf("expected count after popping to be %d. got %d", tc.count, tc.sut.Count())
			}

			if tc.sut.Capacity() != tc.capacity {
				t.Errorf("expected capacity after popping to be %d. got %d", tc.capacity, tc.sut.Capacity())
			}
		})
	}
}

func TestPopE(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[string]
		v           string
		err         error
		count       int
		capacity    int
	}{
		{
			"popping an empty collection",
			Collection[string]{},
			"",
			fmt.Errorf("value not found: empty collection"),
			0,
			0,
		},
		{
			"popping a collection with items",
			Collection[string]{"foo", "bar"},
			"bar",
			nil,
			1,
			2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualV, actualErr := tc.sut.PopE()

			if actualV != tc.v {
				t.Errorf("expected %s. got %s", tc.v, actualV)
			}

			if actualErr != nil && tc.err != nil {
				if actualErr.Error() != tc.err.Error() {
					t.Errorf("expected '%s'. got '%s'", tc.err.Error(), actualErr.Error())
				}
			}

			if tc.sut.Count() != tc.count {
				t.Errorf("expected count after popping to be %d. got %d", tc.count, tc.sut.Count())
			}

			if tc.sut.Capacity() != tc.capacity {
				t.Errorf("expected capacity after popping to be %d. got %d", tc.capacity, tc.sut.Capacity())
			}
		})
	}
}

func TestEach(t *testing.T) {
	var eachResult Collection[any]
	sut := Collection[any]{1, "foo", 1.1}

	sut.Each(func(_ int, v any) {
		eachResult = append(eachResult, v)
	})

	if !reflect.DeepEqual(eachResult, sut) {
		t.Errorf("expected visited values to be %v. got %v", sut, eachResult)
	}
}

func TestTap(t *testing.T) {
	sut := Collection[any]{1, "foo", 1.1}

	c := sut.Tap(func(c Collection[any]) {
		if !reflect.DeepEqual(sut, c) {
			t.Errorf("expected tap collection to equal %v. got %v", sut, c)
		}
	})

	if !reflect.DeepEqual(sut, c) {
		t.Errorf("expected returned collection to equal %v. got %v", sut, c)
	}
}
func TestReversingTwiceYieldsTheSameCollection(t *testing.T) {
	coll := Collect(collections.Range(1, 10)...)

	expected := coll.Copy()
	got := coll.Reverse().Reverse()

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestReverse(t *testing.T) {
	testCases := []struct {
		description string
		input       Collection[int]
		expected    Collection[int]
	}{
		{
			"reversing empty collection",
			Collection[int]{},
			Collection[int]{},
		},
		{
			"reversing collection with a single element",
			Collection[int]{1},
			Collection[int]{1},
		},
		{
			"reversing a collection with 10 elements",
			Collection[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Collection[int]{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := tc.input.Reverse()

			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[any]
		input       any
		i           int
	}{
		{
			"searching on an empty collection",
			Collection[any]{},
			"foo",
			-1,
		},
		{
			"searching a nonexistent element",
			Collection[any]{1, "foo", 1.0},
			"bar",
			-1,
		},
		{
			"searching an existing element",
			Collection[any]{1, "foo", 1.0},
			"foo",
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			i := tc.sut.Search(tc.input)

			if i != tc.i {
				t.Errorf("expected resulting index to be %d. got %d", tc.i, i)
			}
		})
	}
}

func TestSearchE(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[any]
		input       any
		i           int
		err         error
	}{
		{
			"searching on an empty collection",
			Collection[any]{},
			"foo",
			-1,
			fmt.Errorf("value not found"),
		},
		{
			"searching a nonexistent element",
			Collection[any]{1, "foo", 1.0},
			"bar",
			-1,
			fmt.Errorf("value not found"),
		},
		{
			"searching an existing element",
			Collection[any]{1, "foo", 1.0},
			"foo",
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			i, err := tc.sut.SearchE(tc.input)

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
	sut := Collect(4, 1, 3, 2)
	expectedCollection := Collect(1, 2, 3, 4)

	f := func(current, next int) bool {
		return current < next
	}

	if !reflect.DeepEqual(sut.Sort(f), expectedCollection) {
		t.Errorf("expected sorted collection to be %v. got %v", expectedCollection, sut)
	}
}

func TestMap(t *testing.T) {
	testCases := []struct {
		description      string
		sut              Collection[int]
		mappedCollection Collection[int]
	}{
		{
			"mapping an empty collection",
			Collection[int]{},
			Collection[int]{},
		},
		{
			"mapping a collection with values",
			Collect(1, 2, 3),
			Collect(2, 3, 4),
		},
	}

	f := func(_ int, v int) int {
		return v + 1
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mappedCollection := tc.sut.Map(f)

			if !reflect.DeepEqual(mappedCollection, tc.mappedCollection) {
				t.Errorf("expected mapped collection to be %v. got %v", tc.mappedCollection, mappedCollection)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[string]
		v           string
	}{
		{
			"calling First on an empty collection",
			Collection[string]{},
			"",
		},
		{
			"calling first on a collection with values",
			Collect("foo", "bar", "baz"),
			"foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v := tc.sut.First()

			if v != tc.v {
				t.Errorf("expected returned value to be '%s', got '%s'", tc.v, v)
			}
		})
	}
}

func TestFirstE(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[string]
		v           string
		err         error
	}{
		{
			"calling First on an empty collection",
			Collection[string]{},
			"",
			fmt.Errorf("value not found: empty collection"),
		},
		{
			"calling first on a collection with values",
			Collect("foo", "bar", "baz"),
			"foo",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v, err := tc.sut.FirstE()

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
		sut         Collection[string]
		v           string
		err         error
	}{
		{
			"calling Last on an empty collection",
			Collection[string]{},
			"",
			fmt.Errorf("value not found: empty collection"),
		},
		{
			"calling last on a collection with values",
			Collect("foo", "bar", "baz"),
			"baz",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			v, err := tc.sut.LastE()

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

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[string]
		isEmpty     bool
	}{
		{
			"calling IsEmpty on an empty collection",
			Collection[string]{},
			true,
		},
		{
			"calling IsEmpty on a collection with values",
			Collection[string]{"foo"},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if isEmpty := tc.sut.IsEmpty(); isEmpty != tc.isEmpty {
				t.Errorf("expect %v. got %v", tc.isEmpty, isEmpty)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[int]
		expected    Collection[int]
	}{
		{
			"does not change empty collection",
			Collection[int]{},
			Collection[int]{},
		},
		{
			"does not change collection with a single element",
			Collection[int]{1},
			Collection[int]{1},
		},
		{
			"does not change collection with 100 elements",
			Collect(collections.Range(1, 100)...),
			Collect(collections.Range(1, 100)...),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if got := tc.collection.Copy(); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected  %v, got %v", got, got)
			}
		})
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[int]
		matcher     collections.Matcher[int, int]
		contains    bool
	}{
		{
			"collection contains at least one matching value",
			Collect(1, 2, 3, 4),
			collections.ValueEquals[int](3),
			true,
		},
		{
			"collection does not contain matching values",
			Collect(1, 2, 3, 4),
			collections.ValueEquals[int](5),
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if contains := tc.collection.Contains(tc.matcher); contains != tc.contains {
				t.Errorf("Contains result should be  %v. got %v", tc.contains, contains)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	collection := Collect(slice...)

	if !reflect.DeepEqual(collection.ToSlice(), slice) {
		t.Errorf("collection converted to slice should equal %v", slice)
	}
}

func TestForgetE(t *testing.T) {
	testCases := []struct {
		description string
		sut         Collection[string]
		expected    Collection[string]
		i           int
		err         error
	}{
		{
			"deleting a nonexistent key",
			Collect("foo", "bar", "baz"),
			Collect("foo", "bar", "baz"),
			3,
			fmt.Errorf("index out of bounds"),
		},
		{
			"deleting a valid key",
			Collect("foo", "bar", "baz"),
			Collect("foo", "baz"),
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.sut.ForgetE(tc.i)
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
