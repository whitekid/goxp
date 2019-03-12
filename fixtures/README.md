# Fixtures

```go
import "github.com/whitekid/go-utils/fixtures"

func TestFixture(t *testing.T) {
	defer fixtures.Env("HELLO", "WORLD").Teardown()
	assert.Equal(t, "WORLD", os.Getenv("HELLO"))
}
```
