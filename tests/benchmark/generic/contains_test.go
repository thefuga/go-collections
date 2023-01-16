package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkContains(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	half := len(slice) / 2

	for n := 0; n < b.N; n++ {
		Contains(slice, ValueEquals[int](half))
	}
}

func BenchmarkContainsWithDeepMatcher(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	half := len(slice) / 2

	for n := 0; n < b.N; n++ {
		Contains(slice, ValueDeepEquals[int](half))
	}
}
