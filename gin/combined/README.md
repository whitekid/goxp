# gogin combined logger

log as apache combined format not gin

```
127.0.0.1 -  [2018-10-12T00:00:36.970016+09:00] "GET /" 200 -1 "" "curl/7.54.0"
```

## Usage

```go
import 	"github.com/whitekid/go-utils/gin/combined"

router := gin.New()
router.Use(gin.Recovery())
router.Use(New(nil))
...
```

or

```go
router := combined.NewRouter()
...
```
