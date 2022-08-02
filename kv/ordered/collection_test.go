package ordered

import (
	"reflect"
	"sort"
	"testing"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/errors"
	"github.com/thefuga/go-collections/slice"
)

func TestCollectSliceMethod(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	collection := CollectSlice(intSlice)

	for k, v := range intSlice {
		if item := collection.Get(k); item != v {
			t.Error("The keys weren't preserved!")
		}
	}
}

func TestCollectMapMethod(t *testing.T) {
	mapString := map[string]any{"f": "foo", "b": "bar"}

	collection := CollectMap(mapString)

	for k, v := range mapString {
		if item := collection.Get(k); item != v {
			t.Error("The keys weren't preserved!")
		}
	}
}

func TestGetMethod(t *testing.T) {
	collection := Collect(1, 2)

	value, err := collection.GetE(0)

	if err != nil {
		t.Error(err)
	}

	if value != 1 {
		t.Error("Wrong value returned!")
	}

	if _, err = collection.GetE(3); err.Error() != "key '3' not found" {
		t.Error(err)
		t.Error("Getting an unexisting key must return an error!")
	}
}

func TestEachMethod(t *testing.T) {
	collection := Collect[any]("foo", 1, 1.5)

	var (
		foundString bool
		foundInt    bool
		foundFloat  bool
		count       int
	)

	collection.Each(func(k int, v any) {
		switch v {
		case "foo":
			foundString = true
		case 1:
			foundInt = true
		case 1.5:
			foundFloat = true
		default:
			t.Error("A value that was not present on the collection was found!")
		}

		count++
	})

	if !foundString {
		t.Error("The string wasn't found!")
	}

	if !foundInt {
		t.Error("The int value wasn't found!")
	}

	if !foundFloat {
		t.Error("The float value wasn't found!")
	}

	if count < collection.Count() {
		t.Error("The method didn't iterate over collected items!")
	}

	if count > collection.Count() {
		t.Error("The method iterate more times then the items count!")
	}
}

func TestTap(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"})

	collection.Tap(func(c Collection[string, string]) {
		if c.Count() != collection.Count() {
			t.Error("The collections are not equal")
		}
	})
}

func TestSearchEMethod(t *testing.T) {
	items := map[string]any{"foo": "foo", "int": 1, "float": 1.0}
	collection := CollectMap(items)

	for k, v := range items {
		foundKey, err := collection.SearchE(v)

		if foundKey != k {
			t.Error("found key is different than the key corresponding to v")
		}

		if err != nil {
			t.Error(err)
		}
	}

	if _, err := collection.SearchE('a'); err.Error() != "value not found" {
		t.Error("searching an unexisting item must return an error")
	}
}

func TestKeys(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"})
	expectedKeys := slice.Collect("bar", "baz", "foo")

	keys := collection.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("expected %v. Got %v", expectedKeys, keys)
	}
}

func TestSort(t *testing.T) {
	collection := Collect(3, 2, 1)

	collection.Sort(collections.Asc[int]())

	expectedCurrent := 1

	collection.Each(func(_ int, value int) {
		if value != expectedCurrent {
			t.Error("Collection wasn't sorted!")
		}

		expectedCurrent++
	})
}

func TestMap(t *testing.T) {
	collection := Collect(1, 2, 3, 4)

	allEven := collection.Map(func(_ int, v int) int {
		if v%2 == 0 {
			return v
		}

		return v + 1
	})

	allEven.Each(func(_ int, v int) {
		if v%2 != 0 {
			t.Error("Expected all values to be even!")
		}
	})
}

func TestOnly(t *testing.T) {
	collection := CollectMap(map[string]int{"foo": 123, "bar": 456, "baz": 789})
	expectedNewCollection := CollectMap(map[string]int{"foo": 123, "bar": 456})
	keys := []string{"foo", "bar"}

	newCollection := collection.Only(keys)

	if !reflect.DeepEqual(newCollection.values, expectedNewCollection.values) {
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

func TestFirst(t *testing.T) {
	collection := Collect(1, 2, 3)

	if collection.First() != 1 {
		t.Error("The value returned wasn't the first value on the collection!")
	}
}

func TestFirstEmpty(t *testing.T) {
	collection := Collect[any]()

	if collection.First() != nil {
		t.Error("The collection is empty. No value should've been returned")
	}
}

func TestLast(t *testing.T) {
	collection := Collect(1, 2, 3)

	if collection.Last() != 3 {
		t.Error("The value returned wasn't the first value on the collection!")
	}
}

func TestLastEmpty(t *testing.T) {
	collection := Collect[any]()

	if collection.Last() != nil {
		t.Error("The collection is empty. No value should've been returned")
	}
}

func TestPut(t *testing.T) {
	collection := Collect[any](1, "foo", true)

	// collection.Put("float", 1.0)
	collection.Put(3, 1.0)

	key, err := collection.SearchE(1.0)

	if err != nil {
		t.Error("Element wasn't inserted")
	}

	if key != 3 {
		t.Error("The given key wasn't preserved")
	}

	if item := collection.Get(key); item != 1.0 {
		t.Error("The keys were messed up :/")
	}
}

func TestPushWithIncrementableKeys(t *testing.T) {
	var (
		collection        = Collect[any](1, "foo", true)
		expectedPushedKey = collection.Count()
		pushedValue       = 1.0
	)

	collection.Push(pushedValue)

	key, err := collection.SearchE(pushedValue)

	if err != nil {
		t.Error("Element wasn't inserted")
	}

	if key != expectedPushedKey {
		t.Errorf("Expected last element's key to be %d. Got %d", 3, key)
	}

	if item := collection.Get(key); item != pushedValue {
		t.Error("The keys were messed up :/")
	}
}

func TestPushWithNonIncrementableKeys(t *testing.T) {
	var (
		collection  = CollectMap(map[string]any{"a": 1, "b": " foo", "c": true})
		pushedValue = 1.0
	)

	collection.Push(pushedValue)

	_, err := collection.SearchE(pushedValue)

	if err == nil {
		t.Error("Element shouldn't have been inserted")
	}
}

func TestPop(t *testing.T) {
	firstValue := 1
	secondValue := "foo"
	lastValue := true
	collection := Collect[any](firstValue, secondValue, lastValue)

	if collection.Pop() != lastValue {
		t.Error("Popped element wasn't the last!")
	}

	if collection.Last() == lastValue {
		t.Error("Last element wasn't popped!")
	}

	if count := collection.Count(); count != 2 {
		t.Errorf("Expected 2 elements, got %d", count)
	}

	if actualFirstValue := collection.First(); actualFirstValue != firstValue {
		t.Errorf("Expected first value to equal %v, got: %v", firstValue, actualFirstValue)
	}

	if newLastValue := collection.Last(); newLastValue != secondValue {
		t.Errorf("Expected last value to equal %v, got: %v", secondValue, newLastValue)
	}
}

func TestGenericGet(t *testing.T) {
	collection := Collect[any](1, "foo", Collect[any]("bar"))

	intValue, err := Get[int, int](collection, 0)

	if err != nil {
		t.Error(err)
	}

	if intValue != 1 {
		t.Error("Wrong value returned!")
	}

	stringValue, err := Get[int, string](collection, 1)

	if err != nil {
		t.Error(err)
	}

	if stringValue != "foo" {
		t.Error("Wrong value returned!")
	}

	collectionValue, err := Get[int, Collection[int, any]](collection, 2)

	if err != nil {
		t.Error("Getting a generic existing value must always work!")
	}

	if collectionValue.Count() != 1 {
		t.Error("Wrong value returned!")
	}

	if _, err := Get[int, int](collection, 1); err.Error() != "interface conversion: interface {} is string, not int" {
		t.Error("Trying to get a value with the wrong type parameter must return a type error!")
	}
}

func TestToSlice(t *testing.T) {
	values := []int{1, 2, 3, 4}
	collection := CollectSlice(values)

	slice := collection.ToSlice()

	if !reflect.DeepEqual(slice, values) {
		t.Error("ToSlice method didn't return the correct underlying values")
	}

	valuesLen := len(values)
	sliceCap := cap(slice)
	sliceLen := len(slice)

	if sliceCap != valuesLen || sliceLen != valuesLen {
		t.Errorf("Expected sliceLen and sliceCap to equal valuesLen\n"+
			"sliceCap: %d\n"+"sliceLen: %d\n"+"valuesLen: %d\n",
			sliceCap, sliceLen, valuesLen)
	}
}

func TestCombine(t *testing.T) {
	keys := makeCollection[string, string](2)
	keys.Put("0", "first_name")
	keys.Put("1", "last_name")

	values := makeCollection[string, string](2)
	values.Put("0", "Jon")
	values.Put("1", "Doe")

	combined := keys.Combine(values)

	if actualFirstName := combined.Get("first_name"); actualFirstName != "Jon" {
		t.Errorf("Expected first name to be %s, got %s", "Jon", actualFirstName)
	}

	if actualLastName := combined.Get("last_name"); actualLastName != "Doe" {
		t.Errorf("Expected first name to be %s, got %s", "Doe", actualLastName)
	}

	if len(combined.keys) != len(combined.values) {
		t.Error("combined.keys should have the same length as combined.values")
	}
}

func TestCombineDiffKeyLengths(t *testing.T) {
	keys := makeCollection[string, string](2)
	keys.Put("0", "first_name")

	values := makeCollection[string, string](2)
	values.Put("0", "Jon")
	values.Put("1", "Doe")

	combined := keys.Combine(values)

	if actualFirstName := combined.Get("first_name"); actualFirstName != "Jon" {
		t.Errorf("Expected first name to be %s, got %s", "Jon", actualFirstName)
	}

	if _, err := combined.GetE("last_name"); err == nil {
		t.Error("last_name key shouldn't have been combined")
	}
}

func TestConcat(t *testing.T) {
	collectionA := CollectMap(map[string]string{"foo": "a", "bar": "b"})
	collectionB := CollectMap(map[string]string{"baz": "c"})
	expectedCollection := CollectMap(map[string]string{"foo": "a", "bar": "b", "baz": "c"}).
		Sort(collections.Asc[string]())

	concat := collectionA.Concat(collectionB).Sort(collections.Asc[string]())

	if !reflect.DeepEqual(expectedCollection.values, concat.values) {
		t.Errorf(
			"expected concatenated collection values to be %v. Got %v",
			expectedCollection.values,
			concat.values,
		)
	}

	if !reflect.DeepEqual(expectedCollection.keys, concat.keys) {
		t.Errorf(
			"expected concatenated keys collection to be %v. Got %v",
			expectedCollection.keys,
			concat.keys,
		)
	}
}

func TestConcatWithDuplicatedKeys(t *testing.T) {
	collectionA := CollectMap(map[string]string{"foo": "a", "bar": "b"})
	collectionB := CollectMap(map[string]string{"foo": "c"})
	expectedCollection := CollectMap(map[string]string{"foo": "a", "bar": "b"}).
		Sort(collections.Asc[string]())

	concat := collectionA.Concat(collectionB).Sort(collections.Asc[string]())

	if !reflect.DeepEqual(expectedCollection.values, concat.values) {
		t.Errorf(
			"expected concatenated collection values to be %v. Got %v",
			expectedCollection.values,
			concat.values,
		)
	}

	if !reflect.DeepEqual(expectedCollection.keys, concat.keys) {
		t.Errorf(
			"expected concatenated keys collection to be %v. Got %v",
			expectedCollection.keys,
			concat.keys,
		)
	}
}

func TestConcatWithIntKeys(t *testing.T) {
	collectionA := Collect("foo", "bar")
	collectionB := Collect("baz")
	expectedCollection := Collect("foo", "bar", "baz").Sort(collections.Asc[string]())

	concat := collectionA.Concat(collectionB).Sort(collections.Asc[string]())

	if !reflect.DeepEqual(expectedCollection.values, concat.values) {
		t.Errorf(
			"expected concatenated collection values to be %v. Got %v",
			expectedCollection.values,
			concat.values,
		)
	}

	if !reflect.DeepEqual(expectedCollection.keys, concat.keys) {
		t.Errorf(
			"expected concatenated keys collection to be %v. Got %v",
			expectedCollection.keys,
			concat.keys,
		)
	}
}

func TestContainsKey(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "a", "bar": "b"})

	if !collection.Contains(collections.KeyEquals("foo")) {
		t.Error("collection should contain 'foo' key")
	}
}

func TestContainsValue(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "a", "bar": "b"})

	if !collection.Contains(collections.ValueEquals("a")) {
		t.Error("collection should contain 'a' value")
	}
}

func TestEvery(t *testing.T) {
	collection := Collect(2, 2, 2, 2)

	if !collection.Every(collections.ValueEquals(2)) {
		t.Error("all elements in the collection are equal")
	}

	collection.Put(4, 1)

	if collection.Every(collections.ValueEquals(2)) {
		t.Error("the collection contains a different element")
	}
}

func TestFirstOrFail(t *testing.T) {
	var (
		foundKey   any
		foundValue any
		foundErr   error

		collection = Collect("foo", "bar")
	)

	foundKey, foundValue, foundErr = collection.FirstOrFail(collections.ValueEquals("bar"))
	if foundKey != 1 || foundValue != "bar" || foundErr != nil {
		t.Errorf(
			"Expected %d, %s, %v, got %d, %s, %v",
			1, "bar", nil,
			foundKey, foundValue, foundErr,
		)
	}

	foundKey, foundValue, foundErr = collection.FirstOrFail(collections.ValueEquals("baz"))
	if foundKey != 0 || foundValue != "" || foundErr == nil {
		t.Errorf(
			"Expected %d, %s, %v, got %d, %s, %v",
			0, "", errors.NewValueNotFoundError(),
			foundKey, foundValue, foundErr,
		)
	}
}

func TestFlip(t *testing.T) {
	collection := CollectMap(map[string]string{
		"A": "1", "B": "2", "C": "3",
	}).Sort(collections.Asc[string]())

	expectedFlippedCollection := CollectMap(map[string]string{
		"1": "A", "2": "B", "3": "C",
	}).Sort(collections.Asc[string]())

	flippedCollection := collection.Flip().Sort(collections.Asc[string]())

	if !reflect.DeepEqual(flippedCollection.keys, expectedFlippedCollection.keys) {
		t.Logf("%v||%v", flippedCollection.keys, expectedFlippedCollection.keys)
		t.Error("collection keys didn't flip")
	}

	if !reflect.DeepEqual(flippedCollection.values, expectedFlippedCollection.values) {
		t.Logf("%v||%v", flippedCollection.values, expectedFlippedCollection.values)
		t.Error("collection values didn't flip")
	}
}

func TestMerge(t *testing.T) {
	newCollection := CollectMap(map[string]string{
		"A": "foo",
		"B": "bar",
	})

	collectionToMerge := CollectMap(map[string]string{
		"B": "foobar",
		"C": "baz",
	})

	mergedCollection := newCollection.Merge(collectionToMerge)

	for _, v := range []string{"A", "B", "C"} {
		_, err := mergedCollection.GetE(v)

		if err != nil {
			t.Errorf("Expected %v key to be in the collection %v but it was not", v, mergedCollection)
		}
	}

	for _, v := range []string{"foo", "foobar", "baz"} {
		_, err := mergedCollection.SearchE(v)

		if err != nil {
			t.Errorf("Expected %v to be in the collection values %v but it was not", v, mergedCollection)
		}
	}

}

func TestCountBy(t *testing.T) {
	type testCase[T comparable] struct {
		name          string
		collection    Collection[int, int]
		mapper        func(n int) T
		expectedCount map[T]int
	}

	testCases := []testCase[bool]{
		{
			name:          "count evens",
			collection:    Collect(1, 2, 3, 4, 5),
			mapper:        func(n int) bool { return n%2 == 0 },
			expectedCount: map[bool]int{true: 2, false: 3},
		},
		{
			name:          "count odds",
			collection:    Collect(1, 2, 3, 4, 5),
			mapper:        func(n int) bool { return n%2 == 1 },
			expectedCount: map[bool]int{true: 3, false: 2},
		},
		{
			name:          "count ones",
			collection:    Collect(1, 2, 1, 3, 1),
			mapper:        func(n int) bool { return n == 1 },
			expectedCount: map[bool]int{true: 3, false: 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualCount := CountBy(tc.collection, tc.mapper)
			if !reflect.DeepEqual(actualCount, tc.expectedCount) {
				t.Errorf("expected CountBy to equal %v. Got %v", tc.expectedCount, actualCount)
			}
		})
	}

}

func TestFilter(t *testing.T) {
	collection := Collect(1, 2, 3, 4)
	expectedNewCollection := CollectMap(map[int]int{2: 3, 3: 4})

	newCollection := collection.Filter(func(k int, v int) bool {
		return v > 2
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestFilterWithMap(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo", "bar": "bar"})
	expectedNewCollection := CollectMap(map[string]string{"foo": "foo"})

	newCollection := collection.Filter(func(k string, v string) bool {
		return v == "foo"
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestReject(t *testing.T) {
	collection := Collect(0, 1, 2, 3, 4)
	expectedNewCollection := CollectMap(map[int]int{2: 2, 3: 3, 4: 4})

	newCollection := collection.Reject(func(k int, v int) bool {
		return v < 2
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestForget(t *testing.T) {
	collection := CollectMap(map[string]int{"foo": 1, "bar": 2})

	key := "foo"

	newCollection, err := collection.ForgetE(key)

	if err != nil {
		t.Error(err)
	}

	if _, err = newCollection.GetE(key); err == nil {
		t.Errorf("The key %v must does not exist on collection %v", key, newCollection)

	}
}

func TestWhen(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo"})
	expectedNewCollection := CollectMap(map[string]string{"foo": "foo", "bar": "bar"})

	newCollection := collection.When(true, func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestWhenFalse(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo"})

	newCollection := collection.When(false, func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(collection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", collection, newCollection)
	}
}

func TestWhenEmpty(t *testing.T) {
	collection := CollectMap(map[string]string{})
	expectedNewCollection := CollectMap(map[string]string{"bar": "bar"})

	newCollection := collection.WhenEmpty(func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestWhenEmptyWithNotEmptyCollection(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo"})

	newCollection := collection.WhenEmpty(func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(collection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", collection, newCollection)
	}
}

func TestWhenNotEmpty(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo"})
	expectedNewCollection := CollectMap(map[string]string{"foo": "foo", "bar": "bar"})

	newCollection := collection.WhenNotEmpty(func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestWhenNotEmptyWithEmptyCollection(t *testing.T) {
	collection := CollectMap(map[string]string{})

	newCollection := collection.WhenNotEmpty(func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(collection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", collection, newCollection)
	}
}

func TestUnless(t *testing.T) {
	collection := CollectMap(map[string]string{"foo": "foo"})
	expectedNewCollection := CollectMap(map[string]string{"foo": "foo", "bar": "bar"})

	newCollection := collection.Unless(false, func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestUnlessEmpty(t *testing.T) {
	collection := CollectMap(map[string]string{})
	expectedNewCollection := CollectMap(map[string]string{})

	newCollection := collection.UnlessEmpty(func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}

func TestUnlessNotEmpty(t *testing.T) {
	collection := CollectMap(map[string]string{})
	expectedNewCollection := CollectMap(map[string]string{"bar": "bar"})

	newCollection := collection.UnlessNotEmpty(func(c Collection[string, string]) Collection[string, string] {
		return c.Put("bar", "bar")
	})

	if !reflect.DeepEqual(expectedNewCollection.values, newCollection.values) {
		t.Errorf("expected %v. Got %v", expectedNewCollection, newCollection)
	}
}
