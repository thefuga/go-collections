package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSplice(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	spliceAt := len(slice) / 2

	for n := 0; n < b.N; n++ {
		Splice(slice, spliceAt)
	}
}
