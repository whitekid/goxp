package gin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEngine_RunWithContext(t *testing.T) {
	type fields struct {
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// {"default", fields{}, args{"127.0.0.1:9999"}, false},
		{"random port", fields{}, args{"127.0.0.1:0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if err := g.RunWithContext(ctx, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("Engine.RunWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.NotEqual(t, 0, g.Addr.Port)
		})
	}
}
