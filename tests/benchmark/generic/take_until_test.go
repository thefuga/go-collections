package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkTakeUntil(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	matcher := ValueGT(len(slice) / 2)

	for n := 0; n < b.N; n++ {
		TakeUntil(slice, matcher)
	}
}
