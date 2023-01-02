package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
	"github.com/thefuga/go-collections/tests/benchmark"
)

func BenchmarkKeyBy(b *testing.B) {
	slice := benchmark.BuildIntSlice()
	f := func(i int) int { return i / 10 }

	for n := 0; n < b.N; n++ {
		KeyBy(slice, f)
	}
}
