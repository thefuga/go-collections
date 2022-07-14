package collection

import (
	"reflect"
	"sort"
	"testing"
)

func TestCollectSliceMethod(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	collection := CollectSlice(intSlice)

	for k, v := range intSlice {
		if item, _ := collection.Get(k); item != v {
			t.Error("The keys weren't preserved!")
		}
	}
}

func TestCollectMapMethod(t *testing.T) {
	mapString := map[any]any{"foo": "bar", "1": 1}

	collection := CollectMap(mapString)

	for k, v := range mapString {
		if item, _ := collection.Get(k); item != v {
			t.Error("The keys weren't preserved!")
		}
	}
}

func TestGetMethod(t *testing.T) {
	collection := Collect(1, 2)

	value, err := collection.Get(0)

	if err != nil {
		t.Error(err)
	}

	if value != 1 {
		t.Error("Wrong value returned!")
	}

	if _, err = collection.Get(3); err.Error() != "Key '3' wasn't found in the collection!" {
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

	collection.Each(func(k any, v any) {
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

func TestSearchMethod(t *testing.T) {
	items := map[any]any{"foo": "foo", "int": 1, "float": 1.0}
	collection := CollectMap(items)

	for k, v := range items {
		foundKey, err := collection.Search(v)

		if foundKey != k {
			t.Error("found key is different than the key corresponding to v")
		}

		if err != nil {
			t.Error(err)
		}
	}

	if _, err := collection.Search('a'); err.Error() != "Value wasn't found in the collection!" {
		t.Error("searching an unexisting item must return an error")
	}

}

func TestKeys(t *testing.T) {
	collection := CollectMap(map[any]string{"foo": "foo", "bar": "bar", "baz": "baz"})

	keys := collection.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].(string) < keys[j].(string)
	})

	if !reflect.DeepEqual(keys, []any{"bar", "baz", "foo"}) {
		t.Error("the returned keys didn't match the collection keys")
	}
}

func TestSort(t *testing.T) {
	collection := Collect(3, 2, 1)

	collection.Sort(func(current, next int) bool {
		return current < next
	})

	expectedCurrent := 1

	collection.Each(func(_ any, value int) {
		if value != expectedCurrent {
			t.Error("Collection wasn't sorted!")
		}

		expectedCurrent++
	})
}

func TestMap(t *testing.T) {
	collection := Collect(1, 2, 3, 4)

	allEven := collection.Map(func(_ any, v int) int {
		if v%2 == 0 {
			return v
		}

		return v + 1
	})

	allEven.Each(func(_ any, v int) {
		if v%2 != 0 {
			t.Error("Expected all values to be even!")
		}
	})
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

	collection.Put("float", 1.0)

	key, err := collection.Search(1.0)

	if err != nil {
		t.Error("Element wasn't inserted")
	}

	if key != "float" {
		t.Error("The given key wasn't preserved")
	}

	if item, _ := collection.Get(key); item != 1.0 {
		t.Error("The keys were messed up :/")
	}
}

func TestPush(t *testing.T) {
	collection := Collect[any](1, "foo", true)

	collection.Push(1.0)

	key, err := collection.Search(1.0)

	if err != nil {
		t.Error("Element wasn't inserted")
	}

	if key != 3 {
		t.Error("The inserted key should be the former length of the collection")
	}

	if item, _ := collection.Get(key); item != 1.0 {
		t.Error("The keys were messed up :/")
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

	intValue, err := Get[int](collection, 0)

	if err != nil {
		t.Error(err)
	}

	if intValue != 1 {
		t.Error("Wrong value returned!")
	}

	stringValue, err := Get[string](collection, 1)

	if err != nil {
		t.Error(err)
	}

	if stringValue != "foo" {
		t.Error("Wrong value returned!")
	}

	collectionValue, err := Get[Collection[any]](collection, 2)

	if err != nil {
		t.Error("Getting a generic existing value must always work!")
	}

	if collectionValue.Count() != 1 {
		t.Error("Wrong value returned!")
	}

	if _, err := Get[int](collection, 1); err.Error() != "interface conversion: interface {} is string, not int" {
		t.Error("Trying to get a value with the wrong type parameter must return a type error!")
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

	concreteType = Assert[string](genericType)

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
	keys := Collect("first_name", "last_name")
	values := Collect("Jon", "Doe")

	combined := keys.Combine(values)

	if actualFirstName, _ := combined.Get("first_name"); actualFirstName != "Jon" {
		t.Errorf("Expected first name to be %s, got %s", "Jon", actualFirstName)
	}

	if actualLastName, _ := combined.Get("last_name"); actualLastName != "Doe" {
		t.Errorf("Expected first name to be %s, got %s", "Doe", actualLastName)
	}

	if len(combined.keys) != len(combined.values) {
		t.Error("combined.keys should have the same lenght as combined.values")
	}
}

func TestConcat(t *testing.T) {
	collectionA := Collect("foo", "bar")
	collectionB := Collect("baz")
	expectedCollection := Collect("foo", "bar", "baz")

	concat := collectionA.Concat(collectionB)

	if !reflect.DeepEqual(expectedCollection.values, concat.values) {
		t.Error("concatenated collection was different than expected")
	}
}

func TestFirstOrFail(t *testing.T) {
	var (
		foundKey   any
		foundValue any
		foundErr   error

		collection = Collect("foo", "bar")
	)

	foundKey, foundValue, foundErr = collection.FirstOrFail(ValueEquals("bar"))
	if foundKey != 1 || foundValue != "bar" || foundErr != nil {
		t.Errorf(
			"Expected %d, %s, %v, got %d, %s, %v",
			1, "bar", nil,
			foundKey, foundValue, foundErr,
		)
	}

	foundKey, foundValue, foundErr = collection.FirstOrFail(ValueEquals("baz"))
	if foundKey != nil || foundValue != "" || foundErr == nil {
		t.Errorf(
			"Expected %d, %s, %v, got %d, %s, %v",
			0, "", errors.NewValueNotFoundError(),
			foundKey, foundValue, foundErr,
		)
	}
}
