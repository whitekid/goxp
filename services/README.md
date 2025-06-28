# services - Simple Service Framework

**Lightweight service framework with lifecycle management, graceful shutdown, and context-aware execution.**

## Key Features

- **Simple Interface**: Easy-to-implement service contract
- **Lifecycle Management**: Built-in service startup and shutdown
- **Signal Handling**: Graceful shutdown on system signals
- **Context-Aware**: Proper context propagation and cancellation
- **Well Tested**: Comprehensive test coverage

## Quick Start

```go
import (
    "context"
    "time"
    "github.com/whitekid/goxp"
    "github.com/whitekid/goxp/services"
    "github.com/whitekid/goxp/log"
)

// Implement the Service interface
type TimerService struct {
    interval time.Duration
}

func NewTimerService(interval time.Duration) services.Interface {
    return &TimerService{interval: interval}
}

func (s *TimerService) Serve(ctx context.Context) error {
    // Use goxp.Every for periodic tasks
    return goxp.Every(ctx, s.interval, func() error {
        if goxp.IsContextDone(ctx) {
            return nil
        }

        log.Infof("Heartbeat: %s", time.Now().UTC().Format(time.RFC3339))
        return nil
    })
}
```

## Running Services

### Single Service
```go
func main() {
    ctx := context.Background()
    
    // Create and run service
    svc := NewTimerService(time.Second)
    if err := svc.Serve(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### With Signal Handling
```go
func main() {
    // Context with signal cancellation
    ctx, cancel := signal.NotifyContext(context.Background(), 
        os.Interrupt, syscall.SIGTERM)
    defer cancel()
    
    svc := NewTimerService(time.Second)
    if err := svc.Serve(ctx); err != nil {
        log.Errorf("Service error: %v", err)
    }
    
    log.Info("Service shutdown gracefully")
}
```

## Service Interface

```go
type Interface interface {
    Serve(ctx context.Context) error
}
```

Services must implement the `Serve` method that:
- Accepts a context for cancellation
- Returns an error if the service fails
- Handles context cancellation gracefully
- Cleans up resources on shutdown

## Advanced Examples

### HTTP Server Service
```go
type HTTPService struct {
    server *http.Server
}

func NewHTTPService(addr string, handler http.Handler) services.Interface {
    return &HTTPService{
        server: &http.Server{
            Addr:    addr,
            Handler: handler,
        },
    }
}

func (s *HTTPService) Serve(ctx context.Context) error {
    // Start server in goroutine
    go func() {
        if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Errorf("HTTP server error: %v", err)
        }
    }()
    
    // Wait for context cancellation
    <-ctx.Done()
    
    // Graceful shutdown
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    return s.server.Shutdown(shutdownCtx)
}
```

### Database Connection Service
```go
type DatabaseService struct {
    db *sql.DB
    config DatabaseConfig
}

func (s *DatabaseService) Serve(ctx context.Context) error {
    // Initialize database connection
    var err error
    s.db, err = sql.Open(s.config.Driver, s.config.DSN)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    s.db.SetMaxOpenConns(s.config.MaxOpenConns)
    s.db.SetMaxIdleConns(s.config.MaxIdleConns)
    
    // Health check loop
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return s.db.Close()
        case <-ticker.C:
            if err := s.db.PingContext(ctx); err != nil {
                log.Warnf("Database ping failed: %v", err)
            }
        }
    }
}
```

### Multiple Services
```go
func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), 
        os.Interrupt, syscall.SIGTERM)
    defer cancel()
    
    // Run multiple services concurrently
    g, ctx := errgroup.WithContext(ctx)
    
    // Web server
    g.Go(func() error {
        return NewHTTPService(":8080", myHandler).Serve(ctx)
    })
    
    // Background tasks
    g.Go(func() error {
        return NewTimerService(time.Minute).Serve(ctx)
    })
    
    // Database service
    g.Go(func() error {
        return NewDatabaseService(dbConfig).Serve(ctx)
    })
    
    // Wait for all services
    if err := g.Wait(); err != nil {
        log.Errorf("Service group error: %v", err)
    }
}
```

## Best Practices

### Graceful Shutdown
```go
func (s *MyService) Serve(ctx context.Context) error {
    // Setup resources
    resource := acquireResource()
    defer resource.Close()  // Always clean up
    
    // Main service loop
    for {
        select {
        case <-ctx.Done():
            log.Info("Shutting down gracefully")
            return ctx.Err()
        case work := <-s.workChan:
            if err := s.processWork(ctx, work); err != nil {
                log.Errorf("Work processing error: %v", err)
            }
        }
    }
}
```

### Error Handling
```go
func (s *MyService) Serve(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return nil  // Normal shutdown, not an error
        default:
            if err := s.doWork(ctx); err != nil {
                if errors.Is(err, context.Canceled) {
                    return nil  // Context cancellation is expected
                }
                return fmt.Errorf("service failed: %w", err)
            }
        }
    }
}
```

### Resource Management
```go
func (s *MyService) Serve(ctx context.Context) error {
    // Initialize with context timeout
    initCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    if err := s.initialize(initCtx); err != nil {
        return fmt.Errorf("initialization failed: %w", err)
    }
    
    // Ensure cleanup on exit
    defer s.cleanup()
    
    // Service logic...
    return s.run(ctx)
}
```

---

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.
