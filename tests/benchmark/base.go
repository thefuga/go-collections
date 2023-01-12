package benchmark

const CollectionSize = 1000

func BuildIntSlice() []int {
	slice := make([]int, 0, CollectionSize)

	for i := 0; i < CollectionSize; i++ {
		slice = append(slice, i)
	}

	return slice
}
