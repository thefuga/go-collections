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
		if collection[k] != v {
			t.Error("The keys weren't preserved!")
		}
	}
}

func TestCollectMapMethod(t *testing.T) {
	mapString := map[string]any{"foo": "bar", "1": 1}

	collection := CollectMap(mapString)

	for k, v := range mapString {
		if collection[k] != v {
			t.Error("The keys weren't preserved!")
		}
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
	collection := Collect[any]("foo", 1, 1.5)

	for k, v := range collection {
		foundKey, err := collection.Search(v)

		if foundKey != k {
			t.Error("found key is different than the key corresponding to v")
		}

		if err != nil {
			t.Error(err)
		}
	}

	if _, err := collection.Search('a'); err == nil {
		t.Error("searching an unexisting  item should return an error")
	}

}

func TestKeys(t *testing.T) {
	collection := Collection[string, string]{"foo": "foo", "bar": "bar", "baz": "baz"}

	keys := collection.Keys()
	sort.Strings(keys)

	if !reflect.DeepEqual(keys, []string{"bar", "baz", "foo"}) {
		t.Error("the returned keys didn't match the collection keys")
	}
}
