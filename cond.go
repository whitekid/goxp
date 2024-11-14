package goxp

// IfThen run trueF() if condition is true else run the first falseF()
func IfThen(condition bool, thenF func(), falseF ...func()) {
	if condition {
		thenF()
		return
	}

	if len(falseF) > 0 {
		falseF[0]()
	}
}

type ifElse[T any] struct {
	cond  func() bool
	value T
}

func If[T any](cond func() bool, value T) *ifElse[T] {
	return &ifElse[T]{
		value: value,
		cond:  cond,
	}
}

func (i *ifElse[T]) Else(elseValue T) T {
	if i.cond() {
		return i.value
	}

	return elseValue
}

func Ternary[T any](cond bool, trueValue T, falseValue T) T {
	if cond {
		return trueValue
	}

	return falseValue
}

func TernaryF[T any](cond func() bool, trueValue T, falseValue T) T {
	return Ternary(cond(), trueValue, falseValue)
}

func TernaryCF[T any](cond bool, trueFn func() T, falseFn func() T) T {
	if cond {
		return trueFn()
	}

	return falseFn()
}
