package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkNthOffset(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	nth := 10
	offset := 5

	for n := 0; n < b.N; n++ {
		NthOffset(slice, nth, offset)
	}
}
