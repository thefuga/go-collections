package generic

import (
	"testing"

	. "github.com/thefuga/go-collections"
)

func BenchmarkRange(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(0, 1000)
	}
}
