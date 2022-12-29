package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkPartition(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	isEven := func(i int) bool { return i%2 == 0 }

	for n := 0; n < b.N; n++ {
		Partition(slice, isEven)
	}
}
