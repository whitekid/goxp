# fx: functional & generics

Experimental functional programmings & generics

## Slice

### ForEach

```go
r := []string{}
ForEach([]string{"a", "b", "c", "d"}, func(i int, v string) {
    r = append(r, fmt.Sprintf("%d:%s", i, v))
})
require.Equal(t, []string{"0:a", "1:b", "2:c", "3:d"}, r)
```

### ForEachE

```go
r := []string{}
ForEachE([]string{"a", "b", "c", "d"}, func(i int, v string) error {
    r = append(r, fmt.Sprintf("%d:%s", i, v))
    return nil
})
require.Equal(t, []string{"0:a", "1:b", "2:c", "3:d"}, r)
```

### Filter

```go
r := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
require.Equal(t, []int{2, 4}, r)
```

### Map

```go
r := Map([]int{1, 2, 3, 4}, func(v int) string { return strconv.FormatInt(int64(v), 10) })
require.Equal(t, []string{"1", "2", "3", "4"}, r)
```

### Reduce

```go
 r := Reduce([]int{1, 2, 3, 4}, func(x, y int) int { return x + y })
 require.Equal(t, 10, r)
```

### Times

```go
r := Times(10, func(v int) int { return v })
require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, r)
```

### Shuffle

```go
r := Shuffle([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
// []int{4, 1, 5, 2, 7, 3, 6, 9, 8}
```

### Distinct

```go
r := Distinct(1, 2, 3, 4, 4)
require.Equal(t, []{int{1, 2, 3, 4}}, r)
```

### Contains

```go
require.True(t, Contains([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 9))
```

### Index

```go
require.Equal(t, 5, Index([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 6))
```

### Find

```go
v, ok := Find([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(v int) bool { return v == 6 })
require.True(t, ok)
require.Equal(t, 6, v)
```

### Every

```go
require.True(t, Every([]int{1, 2, 3, 4, 5}, []int{2, 4}))
```

### Sample

```go
s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
require.Contains(t, s, Sample(s))
```

### Samples

```go
s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
for _, e := range Samples(s, 5) {
    require.Contains(t, s, e)
}
```

### Zip

```go
r := Zip([]int{1, 2, 3}, []string{"a", "b", "c"})
require.Equal(t, map[int]string{
    1: "a",
    2: "b",
    3: "c",
}, r)
```

## Map

### Keys

```go
v := map[int]string{1: "a", 2: "b", 3: "c"}
require.Equal(t, []int{1, 2, 3}, Keys(v))
```

### Values

```go
v := map[int]string{1: "a", 2: "b", 3: "c"}
require.Equal(t, []string{"a", "b", "c"}, Values(v))
```

### FilterMap

```go
r := FilterMap(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) bool { return k%2 == 0 })
require.Equal(t, map[int]string{2: "b"}, r)
```

### ForEachMap

```go
r := map[string]int{}
ForEachMap(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) { r[v] = k })
require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, r)
```

### MapKeys

```go
r := MapKeys(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) string { return v })
require.Equal(t, map[string]string{"a": "a", "b": "b", "c": "c"}, r)
```

### MapValues

```go
r := MapValues(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) int { return k })
require.Equal(t, map[int]int{1: 1, 2: 2, 3: 3}, r)
```

### MergeMap

```go
m1 := map[int]string{1: "a", 2: "b"}
m2 := map[int]string{3: "c", 4: "d"}
require.Equal(t, map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}, MergeMap(m1, m2))
```

### SampleMap

```go
m := map[int]string{1: "a", 2: "b", 3: "c"}
k, v := SampleMap(m)
require.Equal(t, v, m[k])
```

## Set

```go
set := NewSet[int]()
set.Append(1, 2, 3, 3, 4, 5, 6, 1, 1, 2, 3, 2, 3, 5)
require.Equal(t, []int{1, 2, 3, 4, 5, 6}, set.Slice())
require.Equal(t, len(set.Slice()), set.Len())
for _, e := range set.Slice() {
    require.True(t, set.Has(e))
}
```

## Condition

### IfElse

```go
require.Equal(t, "true", If(func() bool { return true }, "true").Else("false"))
```

## Ternary

```go
r := Ternary(10%2 == 0, "even", "odd")
require.Equal(t, "even", r)
```

## Math

### Sum

```go
require.Equal(t, 6, Sum([]int{1, 2, 3}))
```

### Max

```go
require.Equal(t, 3, Min([]int{2, 1, 3}]))
```

### Min

```go
require.Equal(t, 1, Min([]int{2, 1, 3}]))
```

## alternatives

- <https://github.com/thoas/go-funk>
- <https://github.com/samber/lo>
