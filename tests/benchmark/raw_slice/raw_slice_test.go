package raw_slice

import "testing"

const collectionSize = 1000

var (
	rawSlice           = buildIntSlice(collectionSize)
	rawSliceResult     int
	genericSliceResult int
)

func buildIntSlice(n int) []int {
	slice := make([]int, 0, n)

	for i := 0; i < n; i++ {
		slice = append(slice, i)
	}

	return slice
}

func BenchmarkSliceGet(b *testing.B) {
	var r int
	sliceLen := len(rawSlice) - 1

	for n := 0; n < b.N; n++ {
		r = rawSlice[n%sliceLen]
	}

	rawSliceResult = r
}

func BenchmarkGenericSliceGet(b *testing.B) {
	var r int
	sliceLen := len(rawSlice) - 1

	for n := 0; n < b.N; n++ {
		r = genericGet(n%sliceLen, rawSlice)
	}

	genericSliceResult = r
}

func BenchmarkSliceAppend(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rawSlice = append(rawSlice, n)
	}
}

func BenchmarkSlicePut(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rawSlice = append(rawSlice[:(len(rawSlice)/2)+1], rawSlice[len(rawSlice)/2:]...)
		rawSlice[(len(rawSlice) / 2)] = n
	}
}

func genericGet[T any](k int, in []T) T {
	return in[k]
}
