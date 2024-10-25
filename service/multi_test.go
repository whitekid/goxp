package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMulti(t *testing.T) {
	services := []Interface{newSampleService(), newSampleService()}
	m := NewMulti(services...)
	m.Append(newSampleService())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	require.Nil(t, m.Serve(ctx))
	for _, svc := range services {
		require.True(t, svc.(*simpleService).started)
	}
}
