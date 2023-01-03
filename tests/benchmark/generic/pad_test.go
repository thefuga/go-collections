package generic

import (
	"testing"

	"github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkPad(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	half := len(slice) / 2

	for n := 0; n < b.N; n++ {
		collections.Pad(slice[:half], len(slice), 0)
	}
}
