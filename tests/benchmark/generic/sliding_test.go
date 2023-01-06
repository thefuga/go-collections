package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSliding(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	window := len(slice) / 10

	for n := 0; n < b.N; n++ {
		Sliding(slice, window)
	}
}
