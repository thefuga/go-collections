package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkNth(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	nth := 10

	for n := 0; n < b.N; n++ {
		Nth(slice, nth)
	}
}
