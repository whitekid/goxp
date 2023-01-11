package fx

// To convert []any to []T
func To[T1 any](items []any) []T1 { return Map(items, func(x any) T1 { return x.(T1) }) }
