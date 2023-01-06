package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSkipWhile(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	halfway := len(slice) / 2
	matcher := func(i, _ any) bool { return i.(int) == halfway }

	for n := 0; n < b.N; n++ {
		SkipWhile(slice, matcher)
	}
}
