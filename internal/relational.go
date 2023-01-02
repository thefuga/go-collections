package internal

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
