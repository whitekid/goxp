# Service; simple service framework

make simple service

```go
type simpleService struct {
}

func newSampleService() Interface {
    return &simpleService{}
}

func (s *simpleService) Serve(ctx context.Context) error {
    goex.Every(ctx, time.Second, func() error {
        if goex.IsContextDone(ctx) {
            return nil
        }

        log.Infof("Now: %s", time.Now().UTC().Format(time.RFC3339))
        return nil
    })

    <-ctx.Done()
    return nil
}
```

run the service and wait terminate

```go
    svc := newSampleService()
    svc.Serve(ctx)
```
