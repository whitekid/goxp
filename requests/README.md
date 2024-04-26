# Easy Request Library for golang

```go
resp, err := Get("https://api.github.com").Do(context.Background())

r := map[string]string{}
resp.JSON(&r)
defer resp.Body.Close()

fmt.Printf("%s\n", r["hub_url"])
```

for more detailed usage please refer [test cases](request_test.go),
