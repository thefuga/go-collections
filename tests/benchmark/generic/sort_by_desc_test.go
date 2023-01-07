package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSortByDesc(b *testing.B) {
	slice := Shuffle(benchmark.BuildIntSlice())
	identity := func(i int) int { return i }

	for n := 0; n < b.N; n++ {
		SortByDesc(slice, identity)
	}
}
