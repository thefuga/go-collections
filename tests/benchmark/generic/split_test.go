package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkSplit(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	numberOfGroups := len(slice) / 10

	for n := 0; n < b.N; n++ {
		Split(slice, numberOfGroups)
	}
}
