package kv

import (
	"fmt"
	"sort"

	"github.com/thefuga/go-collections"
)

func ExampleCollectMap() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	c.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", c.ToSlice())
	// Output:
	// [123 456]
}

func ExampleCollectSlice() {
	fmt.Printf("%v", CollectSlice([]int{123, 456}))
	// Output:
	// {[0 1] map[0:123 1:456]}
}

func ExampleCollect() {
	fmt.Printf("%v", Collect([]int{123, 456}))
	// Output:
	// {[0] map[0:[123 456]]}
}

func ExampleCollection_Get() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456, "baz": 789})
	fmt.Printf("%d", c.Get("foo"))
	// Output:
	// 123
}

func ExampleCollection_Put() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})
	c = c.Put("baz", 789)
	fmt.Printf("%v", c.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_Push() {
	c := Collect[int](123, 456)
	c.Push(789)
	fmt.Printf("%v", c.Get(2))
	// Output:
	// 789
}

func ExampleCollection_Pop() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	c.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", c.Pop())
	// Output:
	// 456
}

func ExampleCollection_IsEmpty() {
	c := CollectMap(map[string]int{})
	fmt.Printf("%v", c.IsEmpty())
	// Output:
	// true
}

func ExampleCollection_Count() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})
	fmt.Printf("%v", c.Count())
	// Output:
	// 2
}

func ExampleCollection_Each() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})
	sut := CollectMap(map[string]int{"baz": 789, "fuu": 101})

	sut.Each(func(k string, v int) {
		c.Put(k, v)
	})

	fmt.Printf("%v", c.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_Search() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})
	fmt.Printf("%v", c.Search(123))
	// Output:
	// foo
}

func ExampleCollection_Keys() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	keys := c.Keys()

	sort.Strings(keys) // just to assert the output correctly, ignore it

	fmt.Printf("%v", keys)
	// Output:
	// [bar foo]
}

func ExampleCollection_Sort() {
	collection := Collect(3, 2, 1)

	collection.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", collection.First())
	// Output:
	// 1
}

func ExampleCollection_Map() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	c = c.Map(func(_ string, v int) int {
		return v * 2
	})

	fmt.Printf("%v", c.Get("foo"))
	// Output:
	// 246
}

func ExampleCollection_Only() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	fmt.Printf("%v", c.Only([]string{"foo"}))
	// Output:
	// {[foo] map[foo:123]}
}

func ExampleCollection_First() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	c.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", c.First())
	// Output:
	// 123
}

func ExampleCollection_FirstOrFail() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	foundKey, foundValue, foundErr := c.FirstOrFail(collections.KeyEquals("bar"))

	fmt.Printf("%v, %v, %v", foundKey, foundValue, foundErr)
	// Output:
	// bar, 456, <nil>
}

func ExampleCollection_Last() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	c.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", c.Last())
	// Output:
	// 456
}

func ExampleCollection_ToSlice() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	c.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", c.ToSlice())
	// Output:
	// [123 456]
}

func ExampleCollection_Combine() {
	keys := makeCollection[string, string](2)
	keys.Put("0", "first_name")

	values := makeCollection[string, string](2)
	values.Put("0", "Jon")
	values.Put("1", "Doe")

	combined := keys.Combine(values)

	fmt.Printf("%v", combined.Get("first_name"))
	// Output:
	// Jon
}

func ExampleCollection_Concat() {
	collectionA := CollectMap(map[string]int{"foo": 123, "bar": 456})
	collectionB := CollectMap(map[string]int{"baz": 789})

	concat := collectionA.Concat(collectionB)

	fmt.Printf("%v", concat.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_Contains() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	fmt.Printf("%v", c.Contains(collections.KeyEquals("foo")))
	// Output:
	// true
}

func ExampleCollection_Every() {
	c := Collect(2, 2, 2, 2)

	fmt.Printf("%v", c.Every(collections.ValueEquals(2)))
	// Output:
	// true
}

func ExampleCollection_Flip() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456, "baz": 789})

	c = c.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	flippedCollection := c.Flip()

	fmt.Printf("%v", flippedCollection.Get("foo"))
	// Output:
	// 123
}

func ExampleCollection_Merge() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})
	sut := CollectMap(map[string]int{"baz": 789})

	mergedCollection := c.Merge(sut)

	fmt.Printf("%v", mergedCollection.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_Filter() {
	c := Collect(1, 2, 3, 4)

	filteredCollection := c.Filter(func(k int, v int) bool {
		return v > 2
	})

	filteredCollection.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", filteredCollection.ToSlice())
	// Output:
	// [3 4]
}

func ExampleCollection_Reject() {
	c := Collect(1, 2, 3, 4)

	rejectedCollection := c.Reject(func(k int, v int) bool {
		return v > 2
	})

	rejectedCollection.Sort(collections.Asc[int]()) // just to assert the output correctly, ignore it

	fmt.Printf("%v", rejectedCollection.ToSlice())
	// Output:
	// [1 2]
}

func ExampleCollection_When() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	newCollection := c.When(c.Count() == 2, func(c Collection[string, int]) Collection[string, int] {
		return c.Put("baz", 789)
	})

	fmt.Printf("%v", newCollection.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_WhenEmpty() {
	c := CollectMap(map[string]int{})

	newCollection := c.WhenEmpty(func(c Collection[string, int]) Collection[string, int] {
		return c.Put("foo", 123)
	})

	fmt.Printf("%v", newCollection.Get("foo"))
	// Output:
	// 123
}

func ExampleCollection_WhenNotEmpty() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	newCollection := c.WhenNotEmpty(func(c Collection[string, int]) Collection[string, int] {
		return c.Put("baz", 789)
	})

	fmt.Printf("%v", newCollection.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_Unless() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	newCollection := c.Unless(c.Count() != 2, func(c Collection[string, int]) Collection[string, int] {
		return c.Put("baz", 789)
	})

	fmt.Printf("%v", newCollection.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_UnlessEmpty() {
	c := CollectMap(map[string]int{"foo": 123, "bar": 456})

	newCollection := c.UnlessEmpty(func(c Collection[string, int]) Collection[string, int] {
		return c.Put("baz", 789)
	})

	fmt.Printf("%v", newCollection.Get("baz"))
	// Output:
	// 789
}

func ExampleCollection_UnlessNotEmpty() {
	c := CollectMap(map[string]int{})

	newCollection := c.UnlessNotEmpty(func(c Collection[string, int]) Collection[string, int] {
		return c.Put("foo", 123)
	})

	fmt.Printf("%v", newCollection.Get("foo"))
	// Output:
	// 123
}
