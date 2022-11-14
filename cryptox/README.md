# Cryptox

## Encrypt/ Decrypt

```go
key := goxp.RandomString(10)
enc, err := cryptox.Encrypt(key, `동해 물과 백두산이 마르고 닳도록 동해 물과 백두산이 마르고 닳도록`)
if err != nil {
    panic(err)
}

fmt.Printf("%v\n", enc)
dec, err := cryptox.Decrypt(key, enc)
if err != nil {
    panic(err)
}
fmt.Printf("%v\n", dec)
```

[play](https://go.dev/play/p/-Rl8Ci8x0Xp)
