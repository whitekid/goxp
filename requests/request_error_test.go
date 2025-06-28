package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
)

// Test error handling and edge cases
func TestRequestErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Run("invalid URL", func(t *testing.T) {
		_, err := Get("://invalid-url").Do(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "fail to create request")
	})

	t.Run("connection refused", func(t *testing.T) {
		_, err := Get("http://localhost:0").Do(ctx) // Port 0 should be refused
		require.Error(t, err)
		require.Contains(t, err.Error(), "request failed")
	})

	t.Run("context timeout", func(t *testing.T) {
		shortCtx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		_, err := Get("https://httpbin.org/delay/1").Do(shortCtx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "request failed")
	})

	t.Run("malformed JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"invalid": json}`)) // Missing quotes around json
		}))
		defer server.Close()

		resp, err := Get(server.URL).Do(ctx)
		require.NoError(t, err)

		var result map[string]interface{}
		err = resp.JSON(&result)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid character")
	})
}

func TestRequestJSONEncoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body) // Echo back the body
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Run("multiple JSON values", func(t *testing.T) {
		resp, err := Post(server.URL).
			JSON(map[string]string{"key1": "value1"}).
			JSON(map[string]string{"key2": "value2"}).
			Do(ctx)
		require.NoError(t, err)

		body := resp.String()
		require.Contains(t, body, `"key1":"value1"`)
		require.Contains(t, body, `"key2":"value2"`)
	})

	t.Run("invalid JSON value", func(t *testing.T) {
		// Test with a value that can't be JSON encoded
		invalidValue := make(chan int) // Channels can't be JSON encoded
		_, err := Post(server.URL).JSON(invalidValue).Do(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to encode JSON value")
	})
}

func TestResponseReuse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := Get(server.URL).Do(ctx)
	require.NoError(t, err)

	// First call to String()
	body1 := resp.String()
	require.Equal(t, "test response", body1)

	// Second call to String() should return empty because body is already read
	body2 := resp.String()
	require.Empty(t, body2)
}

func TestRequestConcurrency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("concurrent response"))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const numRequests = 10
	results := make(chan error, numRequests)

	// Test concurrent requests with shared client
	client := &http.Client{}
	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := Get(server.URL).WithClient(client).Do(ctx)
			if err != nil {
				results <- err
				return
			}
			
			body := resp.String()
			
			if body != "concurrent response" {
				results <- fmt.Errorf("unexpected response body: %s", body)
				return
			}
			
			results <- nil
		}()
	}

	// Check all results
	for i := 0; i < numRequests; i++ {
		select {
		case err := <-results:
			require.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for concurrent requests")
		}
	}
}

func TestCompressionErrorHandling(t *testing.T) {
	t.Run("malformed gzip", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			w.Write([]byte("not gzip data")) // Invalid gzip data
		}))
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Should panic due to Must() call in request.go when gzip decompression fails
		require.Panics(t, func() {
			resp, err := Get(server.URL).Do(ctx)
			if err == nil {
				resp.String()
			}
		})
	})

	t.Run("unsupported encoding", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "unknown")
			w.Write([]byte("test data"))
		}))
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := Get(server.URL).Do(ctx)
		require.NoError(t, err)

		// Should handle unknown encoding gracefully
		body := resp.String()
		require.Equal(t, "test data", body)
	})
}

func TestHeadersClosure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back headers as JSON
		headers := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		w.Header().Set("Content-Type", "application/json")
		goxp.Must(json.NewEncoder(w).Encode(headers))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test that all headers are set correctly (tests closure variable capture fix)
	headers := map[string]string{
		"X-Test-1": "value1",
		"X-Test-2": "value2", 
		"X-Test-3": "value3",
		"X-Test-4": "value4",
		"X-Test-5": "value5",
	}

	resp, err := Get(server.URL).Headers(headers).Do(ctx)
	require.NoError(t, err)

	var receivedHeaders map[string]string
	err = resp.JSON(&receivedHeaders)
	require.NoError(t, err)

	// Verify all headers were set correctly
	for k, v := range headers {
		require.Equal(t, v, receivedHeaders[k], "Header %s should have value %s", k, v)
	}
}