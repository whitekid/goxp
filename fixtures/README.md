# Fixtures

## env

```go
import "github.com/whitekid/goxp/fixtures"

func TestFixture(t *testing.T) {
    defer fixtures.Env("HELLO", "WORLD")()
    require.Equal(t, "WORLD", os.Getenv("HELLO"))
}
```
