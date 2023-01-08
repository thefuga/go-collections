package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkTakePositive(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	half := len(slice) / 2

	for n := 0; n < b.N; n++ {
		Take(slice, half)
	}
}

func BenchmarkTakeNegative(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	half := len(slice) / -2

	for n := 0; n < b.N; n++ {
		Take(slice, half)
	}
}
