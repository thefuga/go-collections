package generic

import (
	"testing"

	coll "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

var genericResult [][]int

func BenchmarkZip(b *testing.B) {
	var result [][]int
	slice1, slice2 := benchmark.BuildIntSlice(), benchmark.BuildIntSlice()

	for n := 0; n < b.N; n++ {
		result = coll.Zip(slice1, slice2)
	}

	genericResult = result
}
