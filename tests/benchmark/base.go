package benchmark

const collectionSize = 1000

func BuildIntSlice() []int {
	slice := make([]int, 0, collectionSize)

	for i := 0; i < collectionSize; i++ {
		slice = append(slice, i)
	}

	return slice
}
