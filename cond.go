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
