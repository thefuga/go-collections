package collection

import "testing"

const collectionSize = 1000

var (
	slice              = buildIntSlice(collectionSize)
	sliceResult        int
	genericSliceResult int
)

var (
	mapInt       = buildMapFromIntSlice[int](buildIntSlice(collectionSize))
	mapIntResult int
)

var (
	collection       = CollectSlice(buildIntSlice(collectionSize))
	collectionResult int
)

var (
	collectResult   Collection[int, int]
	makeSliceResult []int
	makeMapResult   map[int]int
)

func buildIntSlice(n int) []int {
	slice := make([]int, 0, n)

	for i := 0; i < n; i++ {
		slice = append(slice, i)
	}

	return slice
}

func buildMapFromIntSlice[K comparable](in []int) map[K]int {
	m := make(map[K]int, len(in))

	for k, v := range in {
		var iface interface{} = k
		mk := iface.(K)
		m[mk] = v
	}

	return m
}

func genericGet[T any](k int, in []T) T {
	return in[k]
}

func BenchmarkCollect(b *testing.B) {
	var r Collection[int, int]

	for n := 0; n < b.N; n++ {
		r = CollectSlice(buildIntSlice(collectionSize))
	}
	collectResult = r
}

func BenchmarkMakeSlice(b *testing.B) {
	var r []int

	for n := 0; n < b.N; n++ {
		r = buildIntSlice(collectionSize)
	}

	makeSliceResult = r
}

func BenchmarkMakeMap(b *testing.B) {
	var r map[int]int

	for n := 0; n < b.N; n++ {
		r = buildMapFromIntSlice[int](buildIntSlice(collectionSize))
	}

	makeMapResult = r
}

func BenchmarkSliceGet(b *testing.B) {
	var r int
	sliceLen := len(slice) - 1

	for n := 0; n < b.N; n++ {
		r = slice[n%sliceLen]
	}

	sliceResult = r
}

func BenchmarkGenericSliceGet(b *testing.B) {
	var r int
	sliceLen := len(slice) - 1

	for n := 0; n < b.N; n++ {
		r = genericGet(n%sliceLen, slice)
	}

	genericSliceResult = r
}

func BenchmarkCollectionGet(b *testing.B) {
	var r int
	collectionLen := collection.Count() - 1

	for n := 0; n < b.N; n++ {
		r, _ = collection.Get(n % collectionLen)
	}

	collectionResult = r
}

func BenchmarkMapGet(b *testing.B) {
	var r int
	mapLen := len(mapInt) - 1

	for n := 0; n < b.N; n++ {
		r = mapInt[n%mapLen]
	}

	mapIntResult = r
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

func BenchmarkSliceAppend(b *testing.B) {
	for n := 0; n < b.N; n++ {
		slice = append(slice, n)
	}
}

func BenchmarkSlicePut(b *testing.B) {
	for n := 0; n < b.N; n++ {
		slice = append(slice[:(len(slice)/2)+1], slice[len(slice)/2:]...)
		slice[(len(slice) / 2)] = n
	}
}

func BenchmarkMapInsert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mapInt[len(mapInt)] = n
	}
}
