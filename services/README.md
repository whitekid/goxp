# services - Service Framework

Simple service framework with lifecycle management and graceful shutdown.

## Service Interface

```go
type Interface interface {
    Serve(ctx context.Context) error
}
```

## Basic Example

```go
type TimerService struct {
    interval time.Duration
}

func (s *TimerService) Serve(ctx context.Context) error {
    return goxp.Every(ctx, s.interval, true, func(ctx context.Context) {
        log.Infof("Heartbeat: %s", time.Now().Format(time.RFC3339))
    })
}
```

## With Signal Handling

```go
func main() {
    ctx, cancel := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGTERM)
    defer cancel()

    svc := &TimerService{interval: time.Second}
    if err := svc.Serve(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## HTTP Server Example

```go
type HTTPService struct {
    server *http.Server
}

func (s *HTTPService) Serve(ctx context.Context) error {
    go func() {
        if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
            log.Error(err)
        }
    }()

    <-ctx.Done()

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    return s.server.Shutdown(shutdownCtx)
}
```

## Multiple Services

```go
func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    g, ctx := errgroup.WithContext(ctx)

    g.Go(func() error { return NewHTTPService(":8080").Serve(ctx) })
    g.Go(func() error { return NewTimerService(time.Minute).Serve(ctx) })
    g.Go(func() error { return NewDatabaseService(config).Serve(ctx) })

    if err := g.Wait(); err != nil {
        log.Error(err)
    }
}
```
