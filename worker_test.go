package utils

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoWithWorker(t *testing.T) {
	sum := int32(0)

	ch := make(chan int32)

	sumTo := 1000
	DoWithWorker(4,
		func() {
			defer close(ch)
			for i := 0; i < sumTo; i++ {
				ch <- int32(i)
			}
		},
		func(i int) {
			for x := range ch {
				atomic.AddInt32(&sum, x)
			}
		})

	require.Equal(t, int32(499500), sum)
}

func TestDoWithWorkerCancel(t *testing.T) {
	ch := make(chan int)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*500)
	defer cancel()

	t1 := time.Now()
	DoWithWorker(4,
		func() {
			defer close(ch)
			for i := 0; i < 1000; i++ {
				ch <- i
			}
		},
		func(i int) {
		exit:
			for x := range ch {
				time.Sleep(time.Millisecond*100 + time.Duration(x*0))

				select {
				case <-ctx.Done():
					break exit
				default:
				}
			}
		})

	require.Truef(t, time.Now().Before(t1.Add(time.Second)), "work should done in %s, it takes %s", time.Second, time.Since(t1))
}
