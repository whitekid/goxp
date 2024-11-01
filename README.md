# go utility function collection

[![Go](https://github.com/whitekid/goxp/actions/workflows/go.yml/badge.svg)](https://github.com/whitekid/goxp/actions/workflows/go.yml)

need more detailed usage? please refer test cases.

## goroutines

|                      |                                                         |
| -------------------- | ------------------------------------------------------- |
| `DoWithWorker()`     | iterate `chan` and do with workers                      |
| `Every()`            | run func with interval                                  |
| `After()`            | run function with delay                                 |
| `Async()`,`Async2()` | run function with goroutine and get result asynchronous |

## JSON/ XML IO

|                 |                     |
| --------------- | ------------------- |
| `ReadJSON[T]()` | decode json to type |
| `ReadXML[T]()`  | decode xml to type  |

## misc

|                   |                        |
| ----------------- | ---------------------- |
| `SetBits()`       |                        |
| `ClearBits()`     |                        |
| `IsContextDone()` |                        |
| `Must(error)`     | panic if error not nil |

## string parse with defaults

|                  |     |
| ---------------- | --- |
| `AtoiDef()`      |     |
| `BarseBoolDef()` |     |
| `ParseIntDef()`  |     |

## Parse time and format

|                 |                                             |
| --------------- | ------------------------------------------- |
| `StrToTime()`   | parse string to time for well known layouts |
| `TimeWithLayout | `time.Time` with layouts                    |
| `RFC1123ZTime`  | `Mon, 02 Jan 2006 15:04:05 -0700`           |
| `RFC3339Time`   | `2006-01-02T15:04:05Z07:00`                 |

## Shell execution

### `Exec()` - simple run command

run command and output to stdin/stdout

```go
exc := Exec("ls", "-al")
err := exc.Do(context.Background())
require.NoError(t, err)
```

run command and get output

```go
exc := Exec("ls", "-al")
output, err := exc.Output(context.Background())
require.NoError(t, err)
require.Contains(t, string(output), "README.md")
```

### `Pipe()` - run command with pipe

## Conditional execution

### `If()`, `Else()`, `IfThen()` - run func as condition

```go
IfThen(true, func() { fmt.Printf("true\n") })
// true

IfThen(true, func() { fmt.Printf("true\n") }, func() { fmt.Printf("false\n") })
// true

IfThen(false, func() { fmt.Printf("true\n") }, func() { fmt.Printf("false\n") })
// false

IfThen(false, func() { fmt.Printf("true\n") }, func() { fmt.Printf("false\n") }, func() { fmt.Printf("false\n") })
// false
```

[play](https://go.dev/play/p/wNadBmhNYR-)

### `Ternary()`, `TernaryF()`, `TernaryCF()`

## Random string/ byte generator

### `RandomByte()` - generate random byte

```go
b := RandomByte(10)
// hex.EncodeToString(b) = "4d46ef2f87b8191daf58"
```

### `RandomString()`, `RandomStringWith()` - generate random string

```go
s := RandomString(10)
// s = "$c&I$#LR3Y"

s := RandomStringWith(10, []rune("abcdefg"))
// s = "bbffedabda"
```

## `Timer()` - measure execution time

```go
func doSomething() {
    defer goxp.Timer("doSomething()")()
    time.Sleep(500 * time.Millisecond)
}

doSomething()
// time takes 500.505063ms: doSomething()
```

[play](https://go.dev/play/p/Wcj2Hw5CLL6)

## Tuple

|                    |                              |
| ------------------ | ---------------------------- |
| `Tuple2`, `Tuple3` | `Pack()` and `Unpack()`      |
| `T2`,`T3`          | construct tuple with element |

## sub packages

- [chanx](chanx) - `chan` extensions
- [cobrax](cobrax) - `cobra` and `viper` utility functions
- [cryptox](cryptox) - encrypt/ decrypt functions
- [fixtures](fixtures) - useful fixture functions for test
- [flags](flags) - cobra & viper make easy
- [fx](fx) - experimental: some functional functions by `iter.Seq`
- [httptest](httptest) - test http sever make easy
- [iterx](iterx) - `iter.Seq` extensions
- [log](log) - simple log powered by zap
- [mapx](mapx) - map extesions
- [requests](requests) - simple http client
- [retry](retry) - retrier with backoff
- [services](services) - simple service framework
- [sets](sets) - Set type
- [slicex](slicex) - slice extensions
- [slug](slug) - uuid to slug
- [testx](testx) - unit test utility functions
- [types](types) - Some useful types like `OrderedMap`
- [validate](validate) - validator make easy
- [x509x](x509x) - x509 utility functions
