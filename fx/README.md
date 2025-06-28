# fx - Functional Programming with Go Iterators

**Experimental functional programming utilities powered by Go 1.24+ `iter.Seq` for modern, lazy evaluation patterns.**

## Key Features

- **Modern Iterators**: Built on Go 1.24+ `iter.Seq` for lazy evaluation
- **Experimental**: Cutting-edge functional programming patterns
- **Performance**: Zero-allocation iteration where possible
- **Composable**: Chain operations efficiently
- **Type Safe**: Full generic support

## Requirements

- **Go 1.24+** - Required for `iter.Seq` support
- Enable experimental features if needed

## Quick Start

```go
import "github.com/whitekid/goxp/fx"

// Lazy evaluation with iterators
numbers := fx.Range(1, 100).
    Filter(func(n int) bool { return n%2 == 0 }).
    Map(func(n int) int { return n * n }).
    Take(5)

// Only computed when consumed
for n := range numbers {
    fmt.Println(n)  // 4, 16, 36, 64, 100
}
```

## Functional Operations

### Lazy Evaluation
All operations are lazy - they don't execute until the iterator is consumed:

```go
// This creates a computation pipeline, but doesn't execute
pipeline := fx.Range(1, 1000000).
    Filter(isPrime).
    Map(square).
    Take(10)

// Only now does computation happen
result := fx.ToSlice(pipeline)
```

### Iterator Composition
```go
// Combine multiple data sources
evens := fx.Range(0, 100).Filter(isEven)
odds := fx.Range(1, 100).Filter(isOdd)
combined := fx.Concat(evens, odds)

// Complex transformations
processed := fx.From(data).
    GroupBy(keyFunc).
    Map(processGroup).
    Flatten()
```

## Available Operations

### Generators
- `Range(start, end)` - Generate number sequences
- `Repeat(value, count)` - Repeat values
- `From(slice)` - Convert slice to iterator

### Transformations
- `Map[T,R](iter, func(T) R)` - Transform elements
- `Filter(iter, func(T) bool)` - Select elements
- `FlatMap(iter, func(T) iter.Seq[R])` - Map and flatten

### Aggregations
- `Reduce(iter, func(T,T) T)` - Fold operation
- `Fold(iter, init, func(Acc,T) Acc)` - Fold with initial value
- `GroupBy(iter, func(T) K)` - Group by key

### Utilities
- `Take(iter, n)` - Take first n elements
- `Skip(iter, n)` - Skip first n elements
- `Concat(iter1, iter2)` - Concatenate iterators
- `ToSlice(iter)` - Materialize to slice

## Advanced Examples

### Data Processing Pipeline
```go
type User struct {
    Name string
    Age  int
    City string
}

users := []User{
    {"Alice", 30, "NYC"},
    {"Bob", 25, "SF"},
    {"Carol", 35, "NYC"},
}

// Process users with lazy evaluation
adultsByCity := fx.From(users).
    Filter(func(u User) bool { return u.Age >= 18 }).
    GroupBy(func(u User) string { return u.City }).
    Map(func(group []User) CityStats {
        return CityStats{
            City: group[0].City,
            Count: len(group),
            AvgAge: avgAge(group),
        }
    })

// Only computed when consumed
for stats := range adultsByCity {
    fmt.Printf("%+v\n", stats)
}
```

### Infinite Sequences
```go
// Generate infinite Fibonacci sequence
fibonacci := fx.Generate(func() iter.Seq[int] {
    return func(yield func(int) bool) {
        a, b := 0, 1
        for yield(a) {
            a, b = b, a+b
        }
    }
})

// Take only what you need
first10 := fx.Take(fibonacci, 10)
result := fx.ToSlice(first10)  // [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]
```

### Parallel Processing (Future)
```go
// Experimental: parallel processing with iterators
results := fx.From(largeDataset).
    Parallel(runtime.NumCPU()).
    Map(expensiveOperation).
    Collect()
```

## Performance Characteristics

### Memory Efficiency
- **Lazy Evaluation**: Only processes elements when needed
- **No Intermediate Collections**: Avoids creating temporary slices
- **Streaming**: Can handle infinite or very large datasets

### Optimization Tips
```go
// ✅ Good: Lazy pipeline
pipeline := fx.Range(1, 1000000).
    Filter(isPrime).
    Take(100)  // Only processes until 100 primes found

// ❌ Avoid: Materializing large intermediate results
allPrimes := fx.ToSlice(fx.Range(1, 1000000).Filter(isPrime))
first100 := allPrimes[:100]  // Wasteful
```

## Integration with Standard Library

### Converting to/from Slices
```go
// From slice to iterator
iter := fx.From([]int{1, 2, 3, 4, 5})

// From iterator to slice
result := fx.ToSlice(iter)

// Integration with slices package
import "slices"
sorted := slices.Values(slices.SortedFunc(iter, cmp.Compare))
```

### Working with Maps
```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

// Process map entries
processed := fx.FromMap(m).
    Filter(func(k string, v int) bool { return v > 1 }).
    Map(func(k string, v int) string { return fmt.Sprintf("%s=%d", k, v) })
```

## Experimental Features

**Note**: This package contains experimental features that may change in future versions of Go as the iterator design evolves.

- **Parallel Processing**: Experimental parallel iterator support
- **Async Operations**: Integration with goroutines and channels  
- **Custom Iterators**: Advanced iterator composition patterns

## Migration and Compatibility

### From Traditional Loops
```go
// Traditional approach
var results []int
for _, item := range data {
    if condition(item) {
        results = append(results, transform(item))
    }
}

// Functional approach
results := fx.ToSlice(
    fx.From(data).
        Filter(condition).
        Map(transform),
)
```

### From Other Functional Libraries
The fx package provides familiar functional programming patterns but with the performance benefits of Go's native iterators.

---

**Experimental Notice**: This package uses experimental Go 1.24+ features. APIs may change as the Go iterator design evolves.

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.
