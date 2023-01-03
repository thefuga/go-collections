package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkLastBy(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	f := func(_, n int) bool { return n == len(slice)/2 } // will match halfway through `slice`

	for n := 0; n < b.N; n++ {
		LastBy(slice, f)
	}
}
