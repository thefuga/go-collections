package internal

import "math"

func Min[T Relational](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T Relational](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func DivCeil[T Number](a, b T) T {
	return T(math.Ceil(float64(a) / float64(b)))
}

func DivFloor[T Number](a, b T) T {
	return T(math.Floor(float64(a) / float64(b)))
}
