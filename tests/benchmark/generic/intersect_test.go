package generic

import (
	"testing"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkGenericIntersect(b *testing.B) {
	leftSlice := benchmark.BuildIntSlice()
	rightSlice := benchmark.BuildIntSlice()

	for n := 0; n < b.N; n++ {
		collections.Intersect(leftSlice, rightSlice)
	}
}
