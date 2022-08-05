![version](https://img.shields.io/github/go-mod/go-version/thefuga/go-collections) [![Go Reference](https://pkg.go.dev/badge/github.com/thefuga/go-collections/.svg)](https://pkg.go.dev/github.com/thefuga/go-collections/) ![commit stage](https://github.com/thefuga/go-collections/actions/workflows/commit-stage.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/thefuga/go-collections)](https://goreportcard.com/report/github.com/thefuga/go-collections) [![codecov](https://codecov.io/gh/thefuga/go-collections/branch/main/graph/badge.svg?token=E2TNB1LQQJ)](https://codecov.io/gh/thefuga/go-collections) ![issues](https://img.shields.io/github/issues/thefuga/go-collections)
------
# Collections
Blazingly fast Collections package. Inspired by [Laravel](https://laravel.com/docs/9.x/collections), written in Go.

## Overview
[go-collections](https://github.com/thefuga/go-collections) offers a variety of methods and types to be used with slices and maps.
Initially, its interface was based on Laravel Collections, but many more methods and functionalities will be added.

This package is still under development, there might be breaking changes introduced. Use with caution!

Pull requests are welcome. See [Contributing](https://github.com/thefuga/go-collections#Contributing)

### Generic functions
There are just functions with a slice type parameter to help the usage of slices. The main disadvantage here is the lack of piping.

```go
func generciFunction() {
	Each(
		func(_, v int) {
			fmt.Println(v) // 2, 3, 4, 5
		},
		Map(
			func(_, v int) int {
				fmt.Println(v) // 1, 2, 3, 4
				return v + 1
			},
			[]int{1, 2, 3, 4}i
		),
	)
}
```

### Slice collection
Slice collection is a custom slice type with embedded methods. Most of the available methods here are just calls to the base generic functions, with the advantage of piping calls instead of nesting functions.

```go
func sliceCollection() {
	Collect(1, 2, 3, 4).
		Map(func(_, v int) int {
			fmt.Println(v) // 1, 2, 3, 4
			return v + 1
		}).
		Each(func(_, v int) {
			fmt.Println(v) // 2, 3, 4, 5
		})
}
```

### KV (map) collections
The KV collection is a bit more complex. It uses composed of two structs: the map, holding the keys and values, and a slice of keys, used to enable the collection to be ordered in any way needed by the user. This is important due to the lack of ordering on Go maps. The usage is similar to the Slice collection, with the ability to add comparable keys.
```go
func mapCollection() {
	CollectMap(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}).
		Sort(collections.Asc[int]()).
		Map(func(k string, v int) int {
			fmt.Printf("%s: %v\n", k, v) // a: 1, b: 2, c: 3, d: 4
			return v + 1
		}).
		Each(func(k string, v int) {
			fmt.Printf("%s: %v\n", k, v) // a: 2, b: 3, c: 4, d: 5
		})
}
```

### Why?
Before type parameters (aka generics), Go lacked a way to do even trivial collection stuff without cumbersome boilerplate.
After 1.18 was released, a myriad of possibilities was open. There are many available options right now, such as the awesome [samber/lo](https://github.com/samber/lo).
That said, none of the available options were what we (i.e. maintainers of this package) were looking for... So go-collection was born.
The aim is to not only offer generic methods to be used with slices and maps but also offer custom types to manage collections in a pipeable fashion.

### Major goals
The major goals for this package are:

- offer a good, cleaar interface to interact with slices and maps
- use no external dependencies
- good resource usage without micro-optimizations - this is not intended to be used by real-time or perfomance critical applications

### Limitations
Some of the methods available on the Laravel Collections are - by design - impossible to be implemented in Go.
These methods were either adapted to be used in a way that made sense, or ignored completely.

## Available types and methods
The complete reference of all types and their documentation is listed in their respective godoc.
Below is a quick reference for each one of them.
### Generic
**Types**
- [Float](https://pkg.go.dev/github.com/thefuga/go-collections#Float)
- [Integer](https://pkg.go.dev/github.com/thefuga/go-collections#Integer)
- [Matcher](https://pkg.go.dev/github.com/thefuga/go-collections#Matcher)
- [Number](https://pkg.go.dev/github.com/thefuga/go-collections#Number)
- [Relational](https://pkg.go.dev/github.com/thefuga/go-collections#Relational)
- [SignedInteger](https://pkg.go.dev/github.com/thefuga/go-collections#SignedInteger)
- [UnsignedInteger](https://pkg.go.dev/github.com/thefuga/go-collections#UnsignedInteger)

**Functions**
- [Asc](https://pkg.go.dev/github.com/thefuga/go-collections#Asc)
- [Average](https://pkg.go.dev/github.com/thefuga/go-collections#Average)
- [AverageE](https://pkg.go.dev/github.com/thefuga/go-collections#AverageE)
- [Copy](https://pkg.go.dev/github.com/thefuga/go-collections#Copy)
- [Cut](https://pkg.go.dev/github.com/thefuga/go-collections#Cut)
- [CutE](https://pkg.go.dev/github.com/thefuga/go-collections#CutE)
- [Delete](https://pkg.go.dev/github.com/thefuga/go-collections#Delete)
- [Desc](https://pkg.go.dev/github.com/thefuga/go-collections#Desc)
- [Each](https://pkg.go.dev/github.com/thefuga/go-collections#Each)
- [First](https://pkg.go.dev/github.com/thefuga/go-collections#First)
- [FirstE](https://pkg.go.dev/github.com/thefuga/go-collections#FirstE)
- [Get](https://pkg.go.dev/github.com/thefuga/go-collections#Get)
- [GetE](https://pkg.go.dev/github.com/thefuga/go-collections#GetE)
- [Last](https://pkg.go.dev/github.com/thefuga/go-collections#Last)
- [LastE](https://pkg.go.dev/github.com/thefuga/go-collections#LastE)
- [Map](https://pkg.go.dev/github.com/thefuga/go-collections#Map)
- [Max](https://pkg.go.dev/github.com/thefuga/go-collections#Max)
- [MaxE](https://pkg.go.dev/github.com/thefuga/go-collections#MaxE)
- [Median](https://pkg.go.dev/github.com/thefuga/go-collections#Median)
- [Min](https://pkg.go.dev/github.com/thefuga/go-collections#Min)
- [MinE](https://pkg.go.dev/github.com/thefuga/go-collections#MinE)
- [Pop](https://pkg.go.dev/github.com/thefuga/go-collections#Pop)
- [PopE](https://pkg.go.dev/github.com/thefuga/go-collections#PopE)
- [Push](https://pkg.go.dev/github.com/thefuga/go-collections#Push)
- [Put](https://pkg.go.dev/github.com/thefuga/go-collections#Put)
- [Search](https://pkg.go.dev/github.com/thefuga/go-collections#Search)
- [SearchE](https://pkg.go.dev/github.com/thefuga/go-collections#SearchE)
- [Sort](https://pkg.go.dev/github.com/thefuga/go-collections#Sort)
- [Sum](https://pkg.go.dev/github.com/thefuga/go-collections#Sum)
### Slice collection
Most methods of the slice collection are just calls to the generic functions passing the collection as the slice argument and returning the result to allow piping.
In addition, some methods are available:
- [Collect](https://pkg.go.dev/github.com/thefuga/go-collections/slice#Collect)
- [Collection](https://pkg.go.dev/github.com/thefuga/go-collections/slice#Collection)
  - [Capacity](https://pkg.go.dev/github.com/thefuga/go-collections/slice#Collection.Capacity)
  - [Count](https://pkg.go.dev/github.com/thefuga/go-collections/slice#Collection.Count)
  - [IsEmpty](https://pkg.go.dev/github.com/thefuga/go-collections/slice#Collection.IsEmpty)
  - [Tap](https://pkg.go.dev/github.com/thefuga/go-collections/slice#Collection.Tap)

### Key/Value collection
#### Generic
- [Assert](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Assert)
- [AssertE](https://pkg.go.dev/github.com/thefuga/go-collections/kv#AssertE)
- [CountBy](https://pkg.go.dev/github.com/thefuga/go-collections/kv#CountBy)
- [Get](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Get)
- [Collect](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collect)
- [CollectMap](https://pkg.go.dev/github.com/thefuga/go-collections/kv#CollectMap)
- [CollectSlice](https://pkg.go.dev/github.com/thefuga/go-collections/kv#CollectSlice)
- [Collection](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection)
  - [Combine](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Combine)
  - [Concat](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Concat)
  - [Contains](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Contains)
  - [Count](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Count)
  - [Each](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Each)
  - [Every](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Every)
  - [Filter](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Filter)
  - [First](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.First)
  - [FirstE](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.FirstE)
  - [FirstOrFail](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.FirstOrFail)
  - [Flip](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Flip)
  - [Get](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Get)
  - [GetE](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.GetE)
  - [IsEmpty](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.IsEmpty)
  - [Keys](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Keys)
  - [Last](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Last)
  - [LastE](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.LastE)
  - [Map](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Map)
  - [Merge](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Merge)
  - [Only](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Only)
  - [Pop](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Pop)
  - [PopE](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.PopE)
  - [Push](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Push)
  - [Put](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Put)
  - [Reject](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Reject)
  - [Search](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Search)
  - [SearchE](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.SearchE)
  - [Sort](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Sort)
  - [Tap](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Tap)
  - [ToSlice](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.ToSlice)
  - [Unless](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.Unless)
  - [UnlessEmpty](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.UnlessEmpty)
  - [UnlessNotEmpty](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.UnlessNotEmpty)
  - [When](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.When)
  - [WhenEmpty](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.WhenEmpty)
  - [WhenNotEmpty](https://pkg.go.dev/github.com/thefuga/go-collections/kv#Collection.WhenNotEmpty)
#### Numeric
The numeric collection includes all methods from the generic kv collection and the additional methods listed bellow:
- [Collect](https://pkg.go.dev/github.com/thefuga/go-collections/kv/numeric#Collect)
- [Collection](https://pkg.go.dev/github.com/thefuga/go-collections/kv/numeric#Collection)
  - [Average](https://pkg.go.dev/github.com/thefuga/go-collections/kv/numeric#Collection.Average)
  - [Max](https://pkg.go.dev/github.com/thefuga/go-collections/kv/numeric#Collection.Max)
  - [Median](https://pkg.go.dev/github.com/thefuga/go-collections/kv/numeric#Collection.Median)
  - [Min](https://pkg.go.dev/github.com/thefuga/go-collections/kv/numeric#Collection.Min)
  - [Sum](https://pkg.go.dev/github.com/thefuga/go-collections/kv/numeric#Collection.Sum)

## Performance
Despite the main description, this is not supposed to be a blazingly fast repository. Rather, it's intended to offer a good interface without deprecating performance.
Benchmarks were made comparing the main methods to their respective raw versions using only the native data struct (e.g. slice or map). 
The full result will always be available as a build artifact of the respective release, but below is some of the main methods compared.

### Benchmarks
TODO

## Contributing
TODO
