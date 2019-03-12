# Easy Request Library for golang

```go
req, err := Get("http://www.google.com").Do()
```

## Request

## Response

### JSON

```go
resp, err := Get("https://api.github.com").Do()

r := map[string]string{}
resp.JSON(&r)
defer resp.Body.Close()

fmt.Printf(r["hub_url"])
```
