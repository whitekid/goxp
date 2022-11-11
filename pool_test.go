package goxp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSyncPool(t *testing.T) {
	type Obj struct {
		v string
	}

	p := NewPool(func() *Obj { return &Obj{} })

	require.Equal(t, &Obj{}, p.Get())

	o := &Obj{v: "str"}
	p.Put(o)

	require.Equal(t, o, p.Get())
	require.Equal(t, &Obj{}, p.Get())
}
