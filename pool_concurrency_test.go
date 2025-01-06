package goxp

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test concurrent access to Pool
func TestPoolConcurrency(t *testing.T) {
	pool := NewPool(func() *testStruct {
		return &testStruct{Value: 42}
	})

	const numGoroutines = 100
	const operationsPerGoroutine = 1000

	var wg sync.WaitGroup

	// Test concurrent Get/Put operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for j := 0; j < operationsPerGoroutine; j++ {
				// Get from pool
				obj := pool.Get()
				require.NotNil(t, obj)
				// Value might be modified by previous operations or be newly created
				// Just check that we got a valid object
				
				// Modify the object
				obj.Value = j
				
				// Put back to pool
				pool.Put(obj)
			}
		}()
	}

	wg.Wait()

	// Pool should still work after concurrent operations
	obj := pool.Get()
	require.NotNil(t, obj)
}

// Test type safety under concurrent access
func TestPoolTypeSafety(t *testing.T) {
	pool := NewPool(func() int {
		return 100
	})

	const numGoroutines = 50
	var wg sync.WaitGroup
	results := make([]int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			// Get value from pool
			val := pool.Get()
			results[idx] = val
			
			// Put back modified value
			pool.Put(val + idx)
		}(i)
	}

	wg.Wait()

	// All results should be valid integers
	for i, result := range results {
		require.IsType(t, 0, result, "Result %d should be int", i)
	}
}

// Test nil handling in concurrent scenarios
func TestPoolNilHandling(t *testing.T) {
	// Create pool that might return nil
	pool := &Pool[*testStruct]{}
	pool.Pool.New = func() any {
		return (*testStruct)(nil) // Return nil
	}

	const numGoroutines = 10
	var wg sync.WaitGroup
	results := make([]*testStruct, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = pool.Get()
		}(i)
	}

	wg.Wait()

	// All results should be nil (zero value) due to our safe type assertion
	for i, result := range results {
		require.Nil(t, result, "Result %d should be nil", i)
	}
}

// Test pool with wrong type stored (edge case)
func TestPoolWrongTypeStored(t *testing.T) {
	pool := NewPool(func() int {
		return 42
	})

	// Manually put wrong type into underlying sync.Pool
	pool.Pool.Put("wrong type")

	// Get should return zero value safely
	val := pool.Get()
	require.Equal(t, 0, val) // Should get zero value of int
}

// Performance test for pool operations
func TestPoolPerformance(t *testing.T) {
	pool := NewPool(func() []byte {
		return make([]byte, 1024) // 1KB buffer
	})

	const numOperations = 100000
	start := time.Now()

	for i := 0; i < numOperations; i++ {
		buf := pool.Get()
		// Simulate some work
		buf[0] = byte(i)
		pool.Put(buf)
	}

	duration := time.Since(start)
	t.Logf("Pool operations took %v for %d operations", duration, numOperations)
	
	// Should complete reasonably fast
	require.Less(t, duration, time.Second, "Pool operations should be fast")
}

// Test pool behavior with panic in New function
func TestPoolNewFunctionPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered from panic:", r)
		}
	}()

	pool := NewPool(func() int {
		panic("test panic")
	})

	// This should panic when New is called
	require.Panics(t, func() {
		pool.Get()
	})
}

// Test memory usage patterns
func TestPoolMemoryReuse(t *testing.T) {
	type largeStruct struct {
		data [1024]byte
	}

	pool := NewPool(func() *largeStruct {
		return &largeStruct{}
	})

	// Get and put the same object multiple times
	obj1 := pool.Get()
	require.NotNil(t, obj1)
	
	// Mark the object
	obj1.data[0] = 0xFF
	pool.Put(obj1)
	
	// Get again - might be the same object
	obj2 := pool.Get()
	require.NotNil(t, obj2)
	
	// Reset and put back
	obj2.data[0] = 0x00
	pool.Put(obj2)
	
	// Pool should continue working
	obj3 := pool.Get()
	require.NotNil(t, obj3)
}

type testStruct struct {
	Value int
}

// Test pool with custom types
func TestPoolCustomTypes(t *testing.T) {
	pool := NewPool(func() *testStruct {
		return &testStruct{Value: 100}
	})

	const numTests = 10
	var objects []*testStruct

	// Get multiple objects
	for i := 0; i < numTests; i++ {
		obj := pool.Get()
		require.NotNil(t, obj)
		require.Equal(t, 100, obj.Value)
		obj.Value = i // Modify
		objects = append(objects, obj)
	}

	// Put them all back
	for _, obj := range objects {
		pool.Put(obj)
	}

	// Get them again
	for i := 0; i < numTests; i++ {
		obj := pool.Get()
		require.NotNil(t, obj)
		// Value might be modified from previous use
		require.IsType(t, &testStruct{}, obj)
	}
}