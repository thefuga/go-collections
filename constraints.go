package collections

type (
	UnsignedInteger interface {
		uint8 | uint16 | uint32 | uint64 | uint
	}

	SignedInteger interface {
		int8 | int16 | int32 | int64 | int
	}

	Integer interface {
		UnsignedInteger | SignedInteger
	}

	Float interface {
		float32 | float64
	}

	Number interface {
		Integer | Float
	}

	Relational interface {
		Number | string
	}
)
