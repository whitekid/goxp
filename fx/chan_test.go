package fx

import (
	"context"
	"testing"
	"time"

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

func TestFadeIn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chans := []chan int{}
	for i := 0; i < 5; i++ {
		c := make(chan int)
		chans = append(chans, c)
		go func(c chan int, i int) {
			c <- i
			c <- i
			close(c)
		}(c, i)
	}

	ch := FanIn(ctx, chans...)
	r := []int{}
	for x := range ch {
		r = append(r, x)
	}
	r = Sort(r)
	r = Uniq(r)

	require.Equal(t, []int{0, 1, 2, 3, 4}, r)
}

func TestAsync(t *testing.T) {
	ch := Async(func() int {
		time.Sleep(time.Second)
		return 7
	})
	require.Equal(t, 7, <-ch)
}

func TestAsync2(t *testing.T) {
	ch := Async2(func() (int, time.Time) {
		time.Sleep(time.Second)
		return 7, time.Now()
	})
	v := <-ch
	n, tm := v.Unpack()
	require.Equal(t, 7, n)
	require.True(t, tm.Before(time.Now()))
}
