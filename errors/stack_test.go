package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorf(t *testing.T) {
	baseErr := Errorf(nil, "failed to read file: %s", "file.txt")
	wrappedErr := Errorf(baseErr, "process file failed with error: %v", baseErr)

	tests := [...]struct {
		name         string
		err          error
		format       string
		expectedMsg  string
		containsMsg  []string
		checkIsError error
		checkAsType  any
	}{
		{
			name:        "Basic error message with %v",
			err:         wrappedErr,
			format:      "%v",
			expectedMsg: "process file failed with error: failed to read file: file.txt",
		},
		{
			name:        "Basic error message with %s",
			err:         wrappedErr,
			format:      "%s",
			expectedMsg: "process file failed with error: failed to read file: file.txt",
		},
		{
			name:         "Full error message with stack trace using %+v",
			err:          wrappedErr,
			format:       "%+v",
			containsMsg:  []string{"process file failed with error: failed to read file: file.txt", "errors.TestErrorf"},
			checkIsError: baseErr,
			checkAsType:  &withStack{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf(tt.format, tt.err)
			require.Falsef(t, tt.expectedMsg != "" && got != tt.expectedMsg, "Expected output to be '%s', got '%s'", tt.expectedMsg, got)

			for _, msg := range tt.containsMsg {
				require.Containsf(t, got, msg, "Expected output to contain '%s', got '%s'", msg, got)
			}

			if tt.checkIsError != nil && !errors.Is(tt.err, tt.checkIsError) {
				require.FailNow(t, "Expected errors.Is to return true for wrapped error")
			}

			if tt.checkAsType != nil {
				require.ErrorAs(t, tt.err, &tt.checkAsType, "Expected errors.As to return true for type %T", tt.checkAsType)
			}
		})
	}
}
