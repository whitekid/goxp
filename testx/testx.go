package testx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
)

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func Must1[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return v
}

func NoError1[T1 any](t *testing.T, v goxp.Tuple2[T1, error]) T1 {
	require.NoError(t, v.V2)
	return v.V1
}

func NoError2[T1 any, T2 any](t *testing.T, v goxp.Tuple3[T1, T2, error]) (T1, T2) {
	require.NoError(t, v.V3)
	return v.V1, v.V2
}
