package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkTimes(b *testing.B) {
	identity := func(i int) int { return i }

	for n := 0; n < b.N; n++ {
		Times(benchmark.CollectionSize, identity)
	}
}
