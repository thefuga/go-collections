package main

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
		t.Error("Tha string wasn't found!")
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
		t.Error("The inserted key should be the former lenght of the collection")
	}

	if item, _ := collection.Get(key); item != 1.0 {
		t.Error("The keys were messed up :/")
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
