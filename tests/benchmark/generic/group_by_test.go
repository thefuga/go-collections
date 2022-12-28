package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkGroupBy(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	idEven := func(i int) bool { return i%2 == 0 }

	for n := 0; n < b.N; n++ {
		GroupBy(slice, idEven)
	}
}
