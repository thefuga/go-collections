package slice

import "fmt"

func ExampleCollect() {
	fmt.Printf("%v\n", Collect(1, 2, 3, 4))
	fmt.Printf("%v", Collect("foo"))
	// Output:
	// [1 2 3 4]
	// [foo]
}

func ExampleCollection_Get() {
	c := Collection[int]{1, 2, 3, 4}
	fmt.Printf("%d", c.Get(3))
	// Output:
	// 4
}

func ExampleCollection_GetE() {
	c := Collection[int]{1, 2, 3, 4}
	_, err := c.GetE(6)
	fmt.Printf("%v", err)
	// Output:
	// value not found: index out of bounds
}

func ExampleCollection_Put() {
	c := Collection[int]{1, 2, 3, 4}
	fmt.Printf("%v", c.Put(4, 5))
	// Output:
	// [1 2 3 4 5]
}

func ExampleCollection_Push() {
	c := Collection[int]{1, 2, 3, 4}
	fmt.Printf("%v", c.Push(5))
	// Output:
	// [1 2 3 4 5]
}

func ExampleCollection_Pop() {
	c := Collection[int]{1, 2, 3, 4}
	fmt.Printf("%v", c.Pop())
	// Output:
	// 4
}

func ExampleCollection_First() {
	c := Collection[int]{1, 2, 3, 4}
	fmt.Printf("%v", c.First())
	// Output:
	// 1
}

func ExampleCollection_Last() {
	c := Collection[int]{1, 2, 3, 4}
	last := c.Last()
	fmt.Printf("%v", last)
	// Output:
	// 4
}

func ExampleCollection_Sort() {
	c := Collection[int]{2, 1, 4, 3}
	fmt.Printf("%v", c.Sort(func(current, next int) bool {
		return current < next
	}))
	// Output:
	// [1 2 3 4]
}

func ExampleCollection_Capacity() {
	c := Collection[int]{1, 2, 3, 4}
	fmt.Printf("%v", c.Capacity())
	// Output:
	// 4
}

func ExampleCollection_Each() {
	c := Collection[int]{1, 2, 3, 4}
	sut := Collection[int]{5, 6, 7}

	sut.Each(func(_ int, v int) {
		c = c.Push(v)
	})

	fmt.Printf("%v", c)
	// Output:
	// [1 2 3 4 5 6 7]
}

func ExampleCollection_IsEmpty() {
	c := Collection[int]{}
	fmt.Printf("%v", c.IsEmpty())
	// Output:
	// true
}

func ExampleCollection_Map() {
	c := Collection[int]{1, 2, 3, 4}

	fmt.Printf("%v", c.Map(func(_ int, v int) int {
		return v * 2
	}))
	// Output:
	// [2 4 6 8]
}

func ExampleCollection_Tap() {
	c := Collection[int]{1, 2, 3, 4}

	fmt.Printf("%v", c.Tap(func(c Collection[int]) {
		fmt.Printf("%v ", c.Count())
	}))
	// Output:
	// 4 [1 2 3 4]
}

func ExampleCollection_Search() {
	c := Collection[int]{1, 2, 3, 4}

	fmt.Printf("%v", c.Search(4))
	// Output:
	// 3
}

func ExampleCollection_Count() {
	c := Collect(1, 2, 3, 4)
	fmt.Printf("%d", c.Count())
	// Output:
	// 4
}
