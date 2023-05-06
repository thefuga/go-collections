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
	result := make([]int, 0, n)

	for i := 0; i < n; i++ {
		result = append(result, i)
	}

	return result
}

func BenchmarkSliceCollectionCopy(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sliceCollection = sliceCollection.Copy()
	}
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
