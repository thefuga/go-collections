package kv

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/thefuga/go-collections"
)

func TestCollect(t *testing.T) {
	items := []string{"foo", "bar", "baz"}

	collection := Collect(items...)

	for k, v := range items {
		foundV, ok := collection[k]

		if !ok {
			t.Errorf("key %d should've been found", k)
		}

		if v != foundV {
			t.Errorf("found value should be %s. got %s", v, foundV)
		}
	}
}

func TestCollectSlice(t *testing.T) {
	items := []string{"foo", "bar", "baz"}

	collection := CollectSlice(items)

	for k, v := range items {
		foundV, ok := collection[k]

		if !ok {
			t.Errorf("key %d should've been found", k)
		}

		if v != foundV {
			t.Errorf("found value should be %s. got %s", v, foundV)
		}
	}
}

func TestCollectMap(t *testing.T) {
	items := map[string]string{"a": "foo", "b": "bar", "c": "baz"}

	collection := CollectMap(items)

	for k, v := range items {
		foundV, ok := collection[k]

		if !ok {
			t.Errorf("key %s should've been found", k)
		}

		if v != foundV {
			t.Errorf("found value should be %s. got %s", v, foundV)
		}
	}
}

func TestCombineE(t *testing.T) {
	testCases := []struct {
		description string
		keys        []string
		values      []string
		expected    map[string]string
		err         error
	}{
		{
			"mismatching keys and values lenghts",
			[]string{"a", "b"},
			[]string{"foo", "bar", "baz"},
			map[string]string{},
			fmt.Errorf("keys and values don't have the same length"),
		},
		{
			"keys and values can be combined",
			[]string{"a", "b", "c"},
			[]string{"foo", "bar", "baz"},
			map[string]string{"a": "foo", "b": "bar", "c": "baz"},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			combined, err := CombineE(tc.keys, tc.values)

			if tc.err != nil {
				if err == nil {
					t.Errorf("expected error '%s'. got nil", tc.err.Error())
				} else if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}

				return
			}

			combinedCount := combined.Count()

			if combinedCount != len(tc.keys) || combinedCount != len(tc.values) {
				t.Errorf(
					"combined collection should have the same count as keys and values. got %d",
					combinedCount,
				)
			}

			for i := 0; i < combinedCount; i++ {
				foundV, ok := combined[tc.keys[i]]

				if !ok {
					t.Errorf("key %s should've been found", tc.keys[i])
				}

				if foundV != tc.values[i] {
					t.Errorf("found value should be %s. got %s", tc.values[i], foundV)
				}
			}
		})
	}
}

func TestCountBy(t *testing.T) {
	type testCase[T comparable] struct {
		description   string
		collection    Collection[int, int]
		mapper        func(n int) T
		expectedCount map[T]int
	}

	testCases := []testCase[bool]{
		{
			description:   "count evens",
			collection:    Collect(1, 2, 3, 4, 5),
			mapper:        func(n int) bool { return n%2 == 0 },
			expectedCount: map[bool]int{true: 2, false: 3},
		},
		{
			description:   "count odds",
			collection:    Collect(1, 2, 3, 4, 5),
			mapper:        func(n int) bool { return n%2 == 1 },
			expectedCount: map[bool]int{true: 3, false: 2},
		},
		{
			description:   "count ones",
			collection:    Collect(1, 2, 1, 3, 1),
			mapper:        func(n int) bool { return n == 1 },
			expectedCount: map[bool]int{true: 3, false: 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualCount := CountBy(tc.collection, tc.mapper)

			if !reflect.DeepEqual(actualCount, tc.expectedCount) {
				t.Errorf("expected CountBy to equal %v. Got %v", tc.expectedCount, actualCount)
			}
		})
	}
}

func TestEach(t *testing.T) {
	eachResult := make(Collection[string, string])
	collection := CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"})

	collection.Each(func(k string, v string) {
		eachResult.Put(k, v)
	})

	if !reflect.DeepEqual(eachResult, collection) {
		t.Errorf("expected visited values to be %v. got %v", collection, eachResult)
	}
}

func TestGetE(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		key         string
		value       string
		err         error
	}{
		{
			"key not found",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"d",
			"",
			fmt.Errorf("key '%s' not found", "d"),
		},
		{
			"key exists",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"b",
			"bar",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			value, err := tc.collection.GetE(tc.key)

			if tc.err != nil {
				if err == nil {
					t.Errorf("expected error '%s'. got nil", tc.err.Error())
				} else if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}

				return
			}

			if value != tc.value {
				t.Errorf("expected found value to be '%s'. got '%s'", tc.value, value)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		key         string
		value       string
	}{
		{
			"key not found",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"d",
			"",
		},
		{
			"key exists",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"b",
			"bar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			value := tc.collection.Get(tc.key)

			if value != tc.value {
				t.Errorf("expected found value to be '%s'. got '%s'", tc.value, value)
			}
		})
	}
}

func TestSearchE(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		value       string
		key         string
		err         error
	}{
		{
			"value not found",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"qux",
			"",
			fmt.Errorf("value not found"),
		},
		{
			"value is found",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"bar",
			"b",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			key, err := tc.collection.SearchE(tc.value)

			if tc.err != nil {
				if err == nil {
					t.Errorf("expected error '%s'. got nil", tc.err.Error())
				} else if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}

				return
			}

			if key != tc.key {
				t.Errorf("expected found value to be '%s'. got '%s'", tc.key, key)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		value       string
		key         string
	}{
		{
			"value not found",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"qux",
			"",
		},
		{
			"value is found",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"bar",
			"b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			key := tc.collection.Search(tc.value)

			if key != tc.key {
				t.Errorf("expected found value to be '%s'. got '%s'", tc.key, key)
			}
		})
	}
}

func TestMap(t *testing.T) {
	collection := CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"})
	expected := CollectMap(map[string]string{"a": "afoo", "b": "bbar", "c": "cbaz"})

	mapped := collection.Map(func(k, v string) string {
		return k + v
	})

	if !reflect.DeepEqual(expected, mapped) {
		t.Errorf("mapped collection should be %v. got %v", expected, mapped)
	}
}

func TestCount(t *testing.T) {
	collection := CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"})

	if count := collection.Count(); count != 3 {
		t.Errorf("count should be %d. got %d", 3, count)
	}
}

func TestPut(t *testing.T) {
	collection := CollectMap(map[string]string{"a": "foo", "b": "bar"})
	expected := CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"})

	collection.Put("c", "baz")

	if !reflect.DeepEqual(expected, collection) {
		t.Errorf("collection after put should be %v. got %v", expected, collection)
	}
}

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[int, string]
		isEmpty     bool
	}{
		{
			"collection is empty",
			Collect[string](),
			true,
		},
		{
			"collection is not empty",
			Collect("foo", "bar", "baz"),
			false,
		},
	}

	for _, tc := range testCases {
		if isEmpty := tc.collection.IsEmpty(); isEmpty != tc.isEmpty {
			t.Errorf("IsEmpty result shoud be %v. got %v", tc.isEmpty, isEmpty)
		}
	}
}

func TestKeys(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []string{"foo", "bar", "baz"}
	collection := Combine(keys, values)

	actualKeys := collection.Keys().Sort(collections.Asc[string]()).ToSlice()

	if !reflect.DeepEqual(keys, actualKeys) {
		t.Errorf("expected collection keys to be %v. got %v", keys, actualKeys)
	}
}

func TestValues(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []string{"bar", "baz", "foo"}
	collection := Combine(keys, values)

	actualValues := collection.Values().Sort(collections.Asc[string]()).ToSlice()

	if !reflect.DeepEqual(values, actualValues) {
		t.Errorf("expected collection values to be %v. got %v", values, actualValues)
	}
}

func TestToSlice(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []string{"bar", "baz", "foo"}
	collection := Combine(keys, values)

	slice := collection.ToSlice()
	collections.Sort(slice, collections.Asc[string]())

	if !reflect.DeepEqual(values, slice) {
		t.Errorf("expected collection values to be %v. got %v", values, slice)
	}
}

func TestOnly(t *testing.T) {
	collection := CollectMap(map[string]int{"foo": 123, "bar": 456, "baz": 789})
	expectedNewCollection := CollectMap(map[string]int{"foo": 123, "bar": 456})
	keys := []string{"foo", "bar"}

	newCollection := collection.Only(keys)

	if !reflect.DeepEqual(newCollection, expectedNewCollection) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestOnlyWithInvalidKeys(t *testing.T) {
	collection := CollectMap(map[string]int{"foo": 123, "bar": 456, "baz": 789})
	keys := []string{"fo", "ar"}

	newCollection := collection.Only(keys)

	if newCollection.Count() != 0 {
		t.Error("The collection should be empty")
	}
}

func TestTap(t *testing.T) {
	var tapped Collection[string, string]
	collection := CollectMap(map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"})

	collection.Tap(func(c Collection[string, string]) {
		tapped = c
	})

	if !reflect.DeepEqual(tapped, collection) {
		t.Errorf("tapped collection should be %v. got %v", collection, tapped)
	}
}

func TestConcat(t *testing.T) {
	a := Combine([]string{"a", "b"}, []string{"foo", "bar"})
	b := Combine([]string{"a", "c", "d"}, []string{"quux", "baz", "qux"})
	expected := Combine([]string{"a", "b", "c", "d"}, []string{"foo", "bar", "baz", "qux"})

	concatenated := a.Concat(b)

	if !reflect.DeepEqual(expected, concatenated) {
		t.Errorf("concatenated collection should be %v. got %v", expected, concatenated)
	}
}

func TestKeysValues(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []string{"bar", "baz", "foo"}

	collection := Combine(keys, values)

	extractedKeys, extractedValues := collection.KeysValues()

	{
		extractedKeys = extractedKeys.Sort(collections.Asc[string]())
		if !reflect.DeepEqual(keys, extractedKeys.ToSlice()) {
			t.Errorf("extracted keys should be %v. got %v", keys, extractedKeys)
		}
	}

	{
		extractedValues = extractedValues.Sort(collections.Asc[string]())
		if !reflect.DeepEqual(values, extractedValues.ToSlice()) {
			t.Errorf("extracted values should be %v. got %v", values, extractedValues)
		}
	}

}

func TestCopy(t *testing.T) {
	collection := Collect("foo", "bar", "baz")
	copiedCollection := collection.Copy()

	if !reflect.DeepEqual(collection, copiedCollection) {
		t.Errorf("copied collection should be %v. got %v", collection, copiedCollection)
	}

	collection.Put(4, "qux")
	if reflect.DeepEqual(collection, copiedCollection) {
		t.Errorf("changing the collection should not affect it's copy")
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[int, string]
		value       string
		contains    bool
	}{
		{
			"collection does not contain the value",
			Collect("foo", "bar", "baz"),
			"qux",
			false,
		},
		{
			"collection contains the value",
			Collect("foo", "bar", "baz"),
			"bar",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if contains := tc.collection.Contains(collections.ValueEquals(tc.value)); contains != tc.contains {
				t.Errorf("Contains result should be %v. got %v", tc.contains, contains)
			}
		})
	}
}

func TestEvery(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[int, string]
		value       string
		contains    bool
	}{
		{
			"not every element on the collection match",
			Collect("foo", "bar", "foo"),
			"foo",
			false,
		},
		{
			"every element matches",
			Collect("foo", "foo", "foo"),
			"foo",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if contains := tc.collection.Every(collections.ValueEquals(tc.value)); contains != tc.contains {
				t.Errorf("Contains result should be %v. got %v", tc.contains, contains)
			}
		})
	}
}

func TestFlip(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []string{"bar", "baz", "foo"}

	collection := Combine(keys, values)
	flipped := collection.Flip()

	{
		flippedKeys := flipped.Keys().Sort(collections.Asc[string]()).ToSlice()

		if !reflect.DeepEqual(flippedKeys, values) {
			t.Errorf("flipped keys should be %v. got %v", values, flippedKeys)
		}
	}

	{
		flippedValues := flipped.Values().Sort(collections.Asc[string]()).ToSlice()

		if !reflect.DeepEqual(flippedValues, keys) {
			t.Errorf("flipped keys should be %v. got %v", keys, flippedValues)
		}
	}
}

func TestFlipE(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []string{"bar", "baz", "foo"}

	collection := Combine(keys, values)
	flipped, err := collection.FlipE()

	if err != nil {
		t.Errorf("expected no error. got'%s'.", err.Error())
	}

	{
		flippedKeys := flipped.Keys().Sort(collections.Asc[string]()).ToSlice()

		if !reflect.DeepEqual(flippedKeys, values) {
			t.Errorf("flipped keys should be %v. got %v", values, flippedKeys)
		}
	}

	{
		flippedValues := flipped.Values().Sort(collections.Asc[string]()).ToSlice()

		if !reflect.DeepEqual(flippedValues, keys) {
			t.Errorf("flipped keys should be %v. got %v", keys, flippedValues)
		}
	}
}

func TestFlipEMismatchingTypes(t *testing.T) {
	collection := Collect("foo", "bar", "baz")
	_, err := collection.FlipE()

	if err == nil {
		t.Errorf("expected error. got nil")
	}
}

func TestMerge(t *testing.T) {
	a := Combine([]string{"a", "b"}, []string{"foo", "bar"})
	b := Combine([]string{"a", "c", "d"}, []string{"quux", "baz", "qux"})
	expected := Combine([]string{"a", "b", "c", "d"}, []string{"quux", "bar", "baz", "qux"})

	concatenated := a.Merge(b)

	if !reflect.DeepEqual(expected, concatenated) {
		t.Errorf("concatenated collection should be %v. got %v", expected, concatenated)
	}
}

func TestFilter(t *testing.T) {
	collection := Collect(1, 2, 3, 4)
	expectedNewCollection := CollectMap(map[int]int{2: 3, 3: 4})

	newCollection := collection.Filter(func(k int, v int) bool {
		return v > 2
	})

	if !reflect.DeepEqual(expectedNewCollection, newCollection) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestReject(t *testing.T) {
	collection := Collect(0, 1, 2, 3, 4)
	expectedNewCollection := CollectMap(map[int]int{2: 2, 3: 3, 4: 4})

	newCollection := collection.Reject(func(k int, v int) bool {
		return v < 2
	})

	if !reflect.DeepEqual(expectedNewCollection, newCollection) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestForgetE(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		key         string
		err         error
	}{
		{
			"key not present on the collection",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"d",
			fmt.Errorf("key 'd' not found"),
		},
		{
			"key is present on the collection",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"b",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			collection, err := tc.collection.ForgetE(tc.key)

			if tc.err != nil {
				if err == nil {
					t.Errorf("expected error '%s'. got nil", tc.err.Error())
				} else if tc.err.Error() != err.Error() {
					t.Errorf("expected error '%s'. got '%s'", tc.err.Error(), err.Error())
				}

				return
			}

			if _, ok := tc.collection[tc.key]; ok {
				t.Errorf("deleted key should not be on the receiver collection")
			}

			if _, ok := collection[tc.key]; ok {
				t.Errorf("deleted key should not be on the returned collection")
			}
		})
	}
}

func TestForget(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		key         string
	}{
		{
			"key not present on the collection",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"d",
		},
		{
			"key is present on the collection",
			CollectMap(map[string]string{"a": "foo", "b": "bar", "c": "baz"}),
			"b",
		},
	}

	for _, tc := range testCases {
		collection := tc.collection.Forget(tc.key)

		if _, ok := tc.collection[tc.key]; ok {
			t.Errorf("deleted key should not be on the receiver collection")
		}

		if _, ok := collection[tc.key]; ok {
			t.Errorf("deleted key should not be on the returned collection")
		}
	}
}

func TestWhen(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		returned    Collection[string, string]
		when        bool
	}{
		{
			"when true",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"c": "baz", "d": "qux"}),
			true,
		},
		{
			"when false",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			false,
		},
	}

	for _, tc := range testCases {
		var calledWith Collection[string, string]

		returned := tc.collection.When(
			tc.when,
			func(c Collection[string, string]) Collection[string, string] {
				calledWith = c
				return tc.returned
			},
		)

		if tc.when {
			if !reflect.DeepEqual(calledWith, tc.collection) {
				t.Errorf("function should've reveived %v. got %v", tc.collection, calledWith)
			}
		}

		if !reflect.DeepEqual(returned, tc.returned) {
			t.Errorf("When should've returned %v. got %v", tc.returned, returned)
		}
	}
}

func TestWhenEmpty(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		returned    Collection[string, string]
	}{
		{
			"when collection is empty",
			Collection[string, string]{},
			CollectMap(map[string]string{"c": "baz", "d": "qux"}),
		},
		{
			"when collection is not empty",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
		},
	}

	for _, tc := range testCases {
		var calledWith Collection[string, string]

		returned := tc.collection.WhenEmpty(
			func(c Collection[string, string]) Collection[string, string] {
				calledWith = c
				return tc.returned
			},
		)

		if tc.collection.IsEmpty() {
			if !reflect.DeepEqual(calledWith, tc.collection) {
				t.Errorf("function should've reveived %v. got %v", tc.collection, calledWith)
			}
		}

		if !reflect.DeepEqual(returned, tc.returned) {
			t.Errorf("When should've returned %v. got %v", tc.returned, returned)
		}
	}
}

func TestWhenNotEmpty(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		returned    Collection[string, string]
	}{
		{
			"when collection is empty",
			Collection[string, string]{},
			Collection[string, string]{},
		},
		{
			"when collection is not empty",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"c": "baz", "d": "qux"}),
		},
	}

	for _, tc := range testCases {
		var calledWith Collection[string, string]

		returned := tc.collection.WhenNotEmpty(
			func(c Collection[string, string]) Collection[string, string] {
				calledWith = c
				return tc.returned
			},
		)

		if !tc.collection.IsEmpty() {
			if !reflect.DeepEqual(calledWith, tc.collection) {
				t.Errorf("function should've reveived %v. got %v", tc.collection, calledWith)
			}
		}

		if !reflect.DeepEqual(returned, tc.returned) {
			t.Errorf("When should've returned %v. got %v", tc.returned, returned)
		}
	}
}

func TestUnless(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		returned    Collection[string, string]
		unless      bool
	}{
		{
			"unless true",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			true,
		},
		{
			"unless false",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"c": "baz", "d": "qux"}),
			false,
		},
	}

	for _, tc := range testCases {
		var calledWith Collection[string, string]

		returned := tc.collection.Unless(
			tc.unless,
			func(c Collection[string, string]) Collection[string, string] {
				calledWith = c
				return tc.returned
			},
		)

		if !tc.unless {
			if !reflect.DeepEqual(calledWith, tc.collection) {
				t.Errorf("function should've reveived %v. got %v", tc.collection, calledWith)
			}
		}

		if !reflect.DeepEqual(returned, tc.returned) {
			t.Errorf("When should've returned %v. got %v", tc.returned, returned)
		}
	}
}

func TestUnlessEmpty(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		returned    Collection[string, string]
	}{
		{
			"unless collection is empty",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"c": "baz", "d": "qux"}),
		},
		{
			"unless collection is not empty",
			Collection[string, string]{},
			Collection[string, string]{},
		},
	}

	for _, tc := range testCases {
		var calledWith Collection[string, string]

		returned := tc.collection.UnlessEmpty(
			func(c Collection[string, string]) Collection[string, string] {
				calledWith = c
				return tc.returned
			},
		)

		if !tc.collection.IsEmpty() {
			if !reflect.DeepEqual(calledWith, tc.collection) {
				t.Errorf("function should've reveived %v. got %v", tc.collection, calledWith)
			}
		}

		if !reflect.DeepEqual(returned, tc.returned) {
			t.Errorf("Unless should've returned %v. got %v", tc.returned, returned)
		}
	}
}

func TestUnlessNotEmpty(t *testing.T) {
	testCases := []struct {
		description string
		collection  Collection[string, string]
		returned    Collection[string, string]
	}{
		{
			"unless collection is empty",
			Collection[string, string]{},
			CollectMap(map[string]string{"c": "baz", "d": "qux"}),
		},
		{
			"unless collection is not empty",
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
			CollectMap(map[string]string{"a": "foo", "b": "bar"}),
		},
	}

	for _, tc := range testCases {
		var calledWith Collection[string, string]

		returned := tc.collection.UnlessNotEmpty(
			func(c Collection[string, string]) Collection[string, string] {
				calledWith = c
				return tc.returned
			},
		)

		if tc.collection.IsEmpty() {
			if !reflect.DeepEqual(calledWith, tc.collection) {
				t.Errorf("function should've reveived %v. got %v", tc.collection, calledWith)
			}
		}

		if !reflect.DeepEqual(returned, tc.returned) {
			t.Errorf("Unless should've returned %v. got %v", tc.returned, returned)
		}
	}
}
