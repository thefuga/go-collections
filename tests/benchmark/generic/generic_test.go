package generic

import (
	"testing"

	"github.com/thefuga/go-collections"
)

func BenchmarkRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = collections.Range[int](-i, i)
	}
}
