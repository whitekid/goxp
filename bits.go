package goxp

// SetBit Sets the bit at pos in the  n.
func SetBit(n byte, pos int) byte {
	n |= (1 << pos)
	return n
}

// ClearBit Clears the bit at pos in n.
func ClearBit(n byte, pos int) byte {
	return n &^ (1 << pos)
}
