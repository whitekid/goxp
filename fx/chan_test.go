package fx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterChan(t *testing.T) {
	ch := make(chan int)

	want := []int{}
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
			want = append(want, i)
		}
		close(ch)
	}()

	got := []int{}
	IterChan(context.Background(), ch, func(i int) {
		got = append(got, i)
	})

	require.Equal(t, want, got)
}
