package fx

type Int interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
}

type Number interface {
	Int | float64
}

type Ordered interface {
	Number | rune | byte | uintptr
}
