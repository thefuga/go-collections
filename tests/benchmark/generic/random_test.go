package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkRandom(b *testing.B) {
	slice := benchmark.BuildIntSlice()

	for n := 0; n < b.N; n++ {
		Random(slice)
	}
}
