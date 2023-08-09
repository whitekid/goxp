package fx

// To convert []any to []T
func To[T any](items []any) []T { return Map(items, func(x any) T { return x.(T) }) }
