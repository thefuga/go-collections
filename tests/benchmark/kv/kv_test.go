package kv

import (
	"testing"

	"github.com/thefuga/go-collections/kv"
	"github.com/thefuga/go-collections/tests/benchmark"
)

var (
	collection       = kv.CollectSlice(benchmark.BuildIntSlice())
	collectResult    kv.Collection[int, int]
	collectionResult int
)

func BenchmarkCollect(b *testing.B) {
	var r kv.Collection[int, int]

	for n := 0; n < b.N; n++ {
		r = kv.CollectSlice(benchmark.BuildIntSlice())
	}
	collectResult = r
}

func BenchmarkCollectionGet(b *testing.B) {
	var r int
	collectionLen := collection.Count() - 1

	for n := 0; n < b.N; n++ {
		r, _ = collection.GetE(n % collectionLen)
	}

	collectionResult = r
}

func BenchmarkCollectionPut(b *testing.B) {
	for n := 0; n < b.N; n++ {
		collection.Put(collection.Count(), n)
	}
}

func BenchmarkCollectionPush(b *testing.B) {
	for n := 0; n < b.N; n++ {
		collection.Push(n)
	}
}
