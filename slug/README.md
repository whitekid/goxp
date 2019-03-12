# UUID to slug for url

```go
uid := uuid.NewV4()

slug := Slug(uid)
fmt.Printf("uuid=%s\n", uid.String())
fmt.Printf("slug=%s\n", slug)

uid1 := UUID(slug)
fmt.Printf("decode=%s\n", uid1.String())
```

output

```text
uuid=25e49082-4ba0-4334-b6fc-b5f111bfad9b
slug=JeSQgkugQzS2_LXxEb-tmw
decode=25e49082-4ba0-4334-b6fc-b5f111bfad9b
```
