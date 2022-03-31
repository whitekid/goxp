# slug for url

## UUID

```go
uid := uuid.New()
sg := slug.NewUUID()
slug := sg.Encode(uid)

fmt.Printf("UUID: %s => slug=%s\n", uid, slug)

uid1 := sg.Decode(slug)
fmt.Printf("slug: %s => UUID=%s\n", slug, uid1)

```

output

```text
UUID: 0ed24fc4-d599-4c06-adae-80526c63d7a4 => slug=DtJPxNWZTAatroBSbGPXpA
slug: DtJPxNWZTAatroBSbGPXpA => UUID=0ed24fc4-d599-4c06-adae-80526c63d7a4
```

## int shortner

```go
encoding := goxp.Shuffle([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"))
shortner := slug.NewShortner(string(encoding))
fmt.Printf("encoding: %s\n\n", string(encoding))

max, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt16))

for i := max.Int64(); i < max.Int64()+10; i++ {
    fmt.Printf("%d => %s\n", i, shortner.Encode(i))
}
```

output

```text
encoding: nPau4_RMk1JScQD927HFhdB8TbUsLNVIA6Crv0fxZp5-myt3lgqwGoEYijzXOWeK

22775 => B9L
22776 => B9A
22777 => B9v
22778 => B9Z
22779 => B9m
22780 => B9l
22781 => B9G
22782 => B9i
22783 => B9O
22784 => B2n
```
