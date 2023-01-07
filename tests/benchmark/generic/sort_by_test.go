package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSortBy(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	identity := func(i int) int { return i }

	for n := 0; n < b.N; n++ {
		SortBy(slice, identity)
	}
}
