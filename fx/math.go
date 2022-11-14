package fx

func Sum[T Ordered](collection []T) T { return Reduce(collection, func(x T, y T) T { return x + y }) }

func Max[T Ordered](col []T) T {
	switch len(col) {
	case 0:
		panic("empty collection")
	case 1:
		return col[0]
	case 2:
		if col[0] > col[1] {
			return col[0]
		}
		return col[1]
	default:
		return Reduce(col, func(x T, y T) T { return Ternary(x > y, x, y) })
	}
}

func Min[T Ordered](col []T) T {
	switch len(col) {
	case 0:
		panic("empty collection")
	case 1:
		return col[0]
	case 2:
		if col[0] > col[1] {
			return col[1]
		}
		return col[1]
	default:
		return Reduce(col, func(x T, y T) T { return Ternary(x > y, y, x) })
	}
}
