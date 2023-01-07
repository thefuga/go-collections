package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSpliceN(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	spliceAt := len(slice) / 2
	splitSize := len(slice) / 2

	for n := 0; n < b.N; n++ {
		SpliceN(slice, spliceAt, splitSize)
	}
}
