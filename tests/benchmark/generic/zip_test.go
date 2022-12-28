package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkZip(b *testing.B) {
	slice1, slice2 := benchmark.BuildIntSlice(), benchmark.BuildIntSlice()

	for n := 0; n < b.N; n++ {
		Zip(slice1, slice2)
	}
}
