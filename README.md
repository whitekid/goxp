# go utility function collection

[![Go](https://github.com/whitekid/goxp/actions/workflows/go.yml/badge.svg)](https://github.com/whitekid/goxp/actions/workflows/go.yml)

need more detailed usage? please refer test cases.

may be not classfied yet..

- `After()` - run func after some duration
- `AvailablePort()` - return random available tcp ports
- `AvailableUdpPort()` - return random available udp ports
- `ClearBit()` - clear bit postion
- `DoWithWorker()` - run go routine with n works
- `Every()` - run goroutine in every duration
- `FileExists()` - return true if file exists
- `Filename()` - return current source file name
- `IsContextDone()` - return true if context is done
- `JsonRecode()` - redecode as new type
- `NewPool()` - `sync.Pool` with type
- `SetBit()` - set bit position
- `SetNX()` - acts as redis SetNX
- `StrToTime()` - parse standard time format as easy
- `URLToListenAddr()` - parse url and get listenable address, ports

## `Exec()` - simple run command

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

## `IfThen()` - run func as condition

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

## Random string/ byte generator

### `RandomByte()` - generate random byte

```go
b := RandomByte(10)
// hex.EncodeToString(b) = "4d46ef2f87b8191daf58"
```

### `RandomString()` - generate random string

```go
s := RandomString(10)
// s = "$c&I$#LR3Y"
```

### `RandomStringWith()` - generate random string

```go
s := RandomStringWith(10, []rune("abcdefg"))
// s = "bbffedabda"
```

### `RandomStringWithCrypto()` - generate random string with `crpto.rand`

```go
s := RandomStringWithCrypto(10)
// s = "d0tu0r3)oZ"
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

## sub packages

- [cryptox](cryptox) - encrypt/ decrypt functions
- [fixtures](fixtures) - useful fixture functions for test
- [flags](flags) - cobra & viper make easy
- [fx](fx) - experimental: some functional functions
- [httptest](httptest) - test http sever make easy
- [log](log) - simple log powered by zap
- [request](request) - simple http client
- [retry](retry) - retrier with backoff
- [service](service) - simple service framework
- [slug](slug) - uuid to slug
- [types](types) - Some useful types
- [validator](validator) - validator make easy
