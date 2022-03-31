# go utility collection

[![Go](https://github.com/whitekid/goxp/actions/workflows/go.yml/badge.svg)](https://github.com/whitekid/goxp/actions/workflows/go.yml)

need more detailed usage? please refer test cases.

## in this package

may be not classfied yet..

- `AvailablePort()` - allocate random available ports
- `URLToListenAddr()` - parse url and get listenable address, ports
- `RandomString()` - generate random string
- `Timer()` - measure execution time
- `SetBit()` - set bit position
- `ClearBit()` - clear bit postion
- `StrToTime()` - parse standard time format as easy
- `DoWithWorker()` - run go routine with n works
- `Every()` - run goroutine in every duration

## sub packages

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
