package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSlidingStep(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	window := len(slice) / 10
	step := 10

	for n := 0; n < b.N; n++ {
		SlidingStep(slice, window, step)
	}
}
