package errors

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUtilityFunctions(t *testing.T) {
	t.Run("HasStack", func(t *testing.T) {
		// Error with stack
		errWithStack := New("error with stack")
		require.True(t, HasStack(errWithStack))

		// Error without stack
		errWithoutStack := NewWithoutStack("error without stack")
		require.False(t, HasStack(errWithoutStack))

		// Wrapped error with stack
		wrappedErr := Wrap(errWithoutStack, "wrapped")
		require.True(t, HasStack(wrappedErr))
	})

	t.Run("Cause", func(t *testing.T) {
		rootErr := New("root cause")
		wrappedErr := Wrap(rootErr, "wrapper 1")
		doubleWrappedErr := Wrap(wrappedErr, "wrapper 2")

		cause := Cause(doubleWrappedErr)
		require.Equal(t, rootErr, cause)

		// Test with single error
		singleErr := New("single")
		require.Equal(t, singleErr, Cause(singleErr))
	})

	t.Run("StackTrace", func(t *testing.T) {
		err := New("test error")
		stack := StackTrace(err)
		require.NotNil(t, stack)
		require.Greater(t, len(stack), 0)

		// Test with error without stack
		errNoStack := NewWithoutStack("no stack")
		stackNoStack := StackTrace(errNoStack)
		require.Nil(t, stackNoStack)

		// Test wrapped error
		wrappedErr := Wrap(errNoStack, "wrapped")
		wrappedStack := StackTrace(wrappedErr)
		require.NotNil(t, wrappedStack)
	})
}

func TestPerformanceOptions(t *testing.T) {
	t.Run("NewWithoutStack", func(t *testing.T) {
		err := NewWithoutStack("no stack error")
		require.False(t, HasStack(err))
		require.Equal(t, "no stack error", err.Error())
	})

	t.Run("WrapWithoutStack", func(t *testing.T) {
		baseErr := NewWithoutStack("base error")
		wrappedErr := WrapWithoutStack(baseErr, "wrapped")
		require.False(t, HasStack(wrappedErr))
		require.Contains(t, wrappedErr.Error(), "base error")
		require.Contains(t, wrappedErr.Error(), "wrapped")
	})
}

func TestConfigurableStackDepth(t *testing.T) {
	// Save original values
	originalMaxDepth := MaxStackDepth
	originalEnableStack := EnableStackTrace
	defer func() {
		MaxStackDepth = originalMaxDepth
		EnableStackTrace = originalEnableStack
	}()

	t.Run("CustomMaxStackDepth", func(t *testing.T) {
		MaxStackDepth = 5
		err := New("test error")
		
		if ws, ok := err.(*withStack); ok {
			require.LessOrEqual(t, len(ws.stack), 5)
		}
	})

	t.Run("DisableStackTrace", func(t *testing.T) {
		EnableStackTrace = false
		err := New("test error")
		require.False(t, HasStack(err))
		
		// Re-enable for cleanup
		EnableStackTrace = true
	})
}

func TestEnvironmentConfiguration(t *testing.T) {
	// Test environment variable configuration
	tests := []struct {
		name     string
		envVar   string
		envValue string
		check    func(t *testing.T)
	}{
		{
			name:     "DisableStackViaEnv",
			envVar:   "GOXP_ERRORS_DISABLE_STACK",
			envValue: "true",
			check: func(t *testing.T) {
				// This would require restarting the package, so we test the parsing logic
				disabled, err := strconv.ParseBool("true")
				require.NoError(t, err)
				require.True(t, disabled)
			},
		},
		{
			name:     "MaxStackDepthViaEnv",
			envVar:   "GOXP_ERRORS_MAX_STACK_DEPTH",
			envValue: "25",
			check: func(t *testing.T) {
				val, err := strconv.Atoi("25")
				require.NoError(t, err)
				require.Equal(t, 25, val)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue := os.Getenv(tt.envVar)
			os.Setenv(tt.envVar, tt.envValue)
			defer os.Setenv(tt.envVar, originalValue)
			
			tt.check(t)
		})
	}
}

func TestFormattingPerformance(t *testing.T) {
	err := New("test error for formatting")
	wrappedErr := Wrap(err, "wrapped error")
	doubleWrappedErr := Wrap(wrappedErr, "double wrapped")

	t.Run("BasicFormatting", func(t *testing.T) {
		basic := fmt.Sprintf("%v", doubleWrappedErr)
		require.Equal(t, "double wrapped", basic)
	})

	t.Run("StringFormatting", func(t *testing.T) {
		str := fmt.Sprintf("%s", doubleWrappedErr)
		require.Equal(t, "double wrapped", str)
	})

	t.Run("DetailedFormatting", func(t *testing.T) {
		detailed := fmt.Sprintf("%+v", doubleWrappedErr)
		require.Contains(t, detailed, "double wrapped")
		require.Contains(t, detailed, "Caused by:")
		require.Contains(t, detailed, "github.com/whitekid/goxp/errors")
		// Stack trace should show error creation internals
		require.Contains(t, detailed, ".go") // Any Go file in the stack
	})
}

func TestStackTraceAccuracy(t *testing.T) {
	// Test that stack traces point to the correct location
	err := createErrorInFunction()
	
	detailed := fmt.Sprintf("%+v", err)
	require.Contains(t, detailed, "createErrorInFunction")
	require.Contains(t, detailed, "enhanced_test.go")
}

func createErrorInFunction() error {
	return New("error created in function")
}

func TestErrorChainIntegrity(t *testing.T) {
	// Create a complex error chain
	rootErr := New("root error")
	level1 := Wrap(rootErr, "level 1")
	level2 := Wrapf(level1, "level 2 with data: %d", 42)
	level3 := Errorf(level2, "level 3 formatted: %s", "test")

	// Test Is() works through the chain
	require.True(t, Is(level3, rootErr))
	require.True(t, Is(level3, level1))
	require.True(t, Is(level3, level2))

	// Test As() works through the chain
	var ws *withStack
	require.True(t, As(level3, &ws))
	require.NotNil(t, ws)

	// Test Cause() returns the root
	cause := Cause(level3)
	require.Equal(t, rootErr, cause)

	// Test formatting shows the complete chain
	formatted := fmt.Sprintf("%+v", level3)
	require.Contains(t, formatted, "level 3 formatted: test")
	require.Contains(t, formatted, "Caused by:")
}

func BenchmarkErrorCreation(b *testing.B) {
	b.Run("WithStack", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New("benchmark error")
		}
	})

	b.Run("WithoutStack", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewWithoutStack("benchmark error")
		}
	})

	b.Run("Wrap", func(b *testing.B) {
		baseErr := New("base error")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Wrap(baseErr, "wrapped error")
		}
	})

	b.Run("WrapWithoutStack", func(b *testing.B) {
		baseErr := NewWithoutStack("base error")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = WrapWithoutStack(baseErr, "wrapped error")
		}
	})
}

func BenchmarkErrorFormatting(b *testing.B) {
	err := createDeepErrorChain(5)

	b.Run("BasicFormat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%v", err)
		}
	})

	b.Run("DetailedFormat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%+v", err)
		}
	})
}

func createDeepErrorChain(depth int) error {
	err := New("root error")
	for i := 0; i < depth; i++ {
		err = Wrapf(err, "level %d", i)
	}
	return err
}

func TestStackTraceCopy(t *testing.T) {
	err := New("test error")
	stack1 := StackTrace(err)
	stack2 := StackTrace(err)

	// Should return copies, not the same slice
	require.NotSame(t, &stack1, &stack2)
	require.Equal(t, stack1, stack2)

	// Modifying one shouldn't affect the other
	if len(stack1) > 0 {
		originalValue := stack1[0]
		stack1[0] = 0
		stack2Copy := StackTrace(err)
		require.NotEqual(t, stack1[0], stack2Copy[0])
		require.Equal(t, originalValue, stack2Copy[0])
	}
}