package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkTakeWhile(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	matcher := ValueLT(len(slice) / 2)

	for n := 0; n < b.N; n++ {
		TakeWhile(slice, matcher)
	}
}
