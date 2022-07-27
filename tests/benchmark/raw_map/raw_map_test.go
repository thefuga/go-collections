package collection

import (
	"testing"

	"github.com/thefuga/go-collections/tests/benchmark"
)

var (
	mapInt       = buildMapFromIntSlice[int](benchmark.BuildIntSlice())
	mapIntResult int
)

var (
	makeSliceResult []int
	makeMapResult   map[int]int
)

func buildMapFromIntSlice[K comparable](in []int) map[K]int {
	m := make(map[K]int, len(in))

	for k, v := range in {
		var iface interface{} = k
		mk := iface.(K)
		m[mk] = v
	}

	return m
}

func BenchmarkMakeSlice(b *testing.B) {
	var r []int

	for n := 0; n < b.N; n++ {
		r = benchmark.BuildIntSlice()
	}

	makeSliceResult = r
}

func BenchmarkMakeMap(b *testing.B) {
	var r map[int]int

	for n := 0; n < b.N; n++ {
		r = buildMapFromIntSlice[int](benchmark.BuildIntSlice())
	}

	makeMapResult = r
}

func BenchmarkMapGet(b *testing.B) {
	var r int
	mapLen := len(mapInt) - 1

	for n := 0; n < b.N; n++ {
		r = mapInt[n%mapLen]
	}

	mapIntResult = r
}

func BenchmarkMapInsert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mapInt[len(mapInt)] = n
	}
}
