package slice

import (
	"testing"

	"github.com/thefuga/go-collections/slice"
	"github.com/thefuga/go-collections/tests/benchmark"
)

var (
	sliceCollection       = slice.Collect(benchmark.BuildIntSlice()...)
	sliceCollectionResult int
)

func buildIntSlice(n int) []int {
	slice := make([]int, 0, n)

	for i := 0; i < n; i++ {
		slice = append(slice, i)
	}

	return slice
}

func BenchmarkSliceCollectionPush(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sliceCollection = sliceCollection.Push(n)
	}
}

func BenchmarkCollectionSlicePut(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sliceCollection = sliceCollection.Put(sliceCollection.Count()/2, n)
	}
}
