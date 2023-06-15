package testx

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
)

func NotNil[T1 any](t *testing.T, values *goxp.Tuple2[T1, error]) T1 {
	v, err := values.Unpack()
	require.NotNil(t, err)

	return v
}

func NoError[T1 any](t *testing.T, values *goxp.Tuple2[T1, error]) T1 {
	v, err := values.Unpack()
	require.NoError(t, err)
	return v
}

func NoError2[T1 any, T2 any](t *testing.T, values *goxp.Tuple3[T1, T2, error]) (T1, T2) {
	v1, v2, err := values.Unpack()
	require.NoError(t, err)
	return v1, v2
}
