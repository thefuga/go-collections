package generic

import (
	"testing"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkForPage(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	size := len(slice) / 10

	for n := 0; n < b.N; n++ {
		collections.ForPage(slice, 5, size)
	}
}
