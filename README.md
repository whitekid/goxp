# go utility collection

[![Go](https://github.com/whitekid/goxp/actions/workflows/go.yml/badge.svg)](https://github.com/whitekid/goxp/actions/workflows/go.yml)

need more detailed usage? please refer test cases.

## in this package

may be not classfied yet..

- `After()` - run func after some duration
- `AvailablePort()` - return random available tcp ports
- `AvailableUdpPort()` - return random available udp ports
- `ClearBit()` - clear bit postion
- `DoWithWorker()` - run go routine with n works
- `Every()` - run goroutine in every duration
- `FileExists()` - return true if file exists
- `Filename()` - return current source file name
- `IfThen()` - run func as condition
- `IsContextDone()` - return true if context is done
- `JsonRedoce()` - redecode as new type
- `NewPool()` - `sync.Pool` with type
- `RandomByte()` - generate random byte
- `RandomString()` - generate random string
- `RandomStringWithCrypto()` - generate random string with `crpto.rand`
- `SetBit()` - set bit position
- `SetNX()` - acts as redis SetNX
- `StrToTime()` - parse standard time format as easy
- `Timer()` - measure execution time
- `URLToListenAddr()` - parse url and get listenable address, ports

## sub packages

### [cryptox](cryptox)

encrypt/ decrypt functions

### [fixtures](fixtures)

useful fixtures for test

### [flags](flags)

cobra & viper make easy

### [fx](fx)

experimental: some functional functions

### [httptest](httptest)

test http sever make easy

### [log](log)

simple log powered by zap

### [request](request)

simple http client

### [retry](retry)

retrier with backoff

### [service](service)

simple service framework

### [slug](slug)

uuid to slug

### [types](types)

Some useful types,

- `Strings` add useful function such as `Index()`, `Remove()`, `Reader()`...
