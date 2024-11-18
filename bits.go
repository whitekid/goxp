package goxp

// SetBit Sets the bit at pos in the  n.
func SetBit[T uint8 | uint16 | uint32 | uint64](n T, pos int) T {
	n |= (1 << pos)
	return n
}

// ClearBit Clears the bit at pos in n.
func ClearBit[T uint8 | uint16 | uint32 | uint64](n T, pos int) T {
	return n &^ (1 << pos)
}
