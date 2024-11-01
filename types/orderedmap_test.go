package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/log"
)

func TestOrderedMap(t *testing.T) {
	m := NewOrderedMap[int, string]()

	m.Set(1, "string1")
	m.Set(2, "string2")
	m.Set(3, "string3")
	m.Set(4, "string4")
	m.Delete(3)
	m.Delete(8)

	require.Equal(t, 3, m.Len())

	for k := range m.Keys() {
		v, _ := m.Get(k)
		log.Infof("%d %s", k, v)
	}

	log.Infof("====")
	m.Each(func(i int, k int, v string) bool {
		log.Infof("%d %d %s", i, k, v)
		return true
	})
}
