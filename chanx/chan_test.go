package chanx

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/log"
	"golang.org/x/sync/errgroup"
)

func TestIterChan(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := make(chan int)

	want := []int{}
	go func() {
		for i := 0; i < 10; i++ {
			i := i
			log.Debugf("@@@@@@@@ <- %v", i)
			ch <- i
			want = append(want, i)
		}
		log.Debugf("@@@@ close")
		close(ch)
	}()

	got := []int{}
	err := Iter(ctx, ch, func(ctx context.Context, i int) error {
		got = append(got, i)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestFadeIn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chans := []<-chan int{}
	for i := 0; i < 5; i++ {
		c := make(chan int)
		chans = append(chans, c)
		go func() {
			c <- i
			close(c)
		}()
	}

	ch := FanIn(ctx, chans...)
	r := []int{}
	for x := range ch {
		r = append(r, x)
	}

	r = slices.Sorted(slices.Values(r))

	require.Equal(t, []int{0, 1, 2, 3, 4}, r)
}

func TestFadeOut(t *testing.T) {
	testFadeOut(t, []int{1, 2, 3, 4, 5, 6, 7, 8})
}

func testFadeOut[T cmp.Ordered](t *testing.T, items []T) {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, e := range items {
			ch <- e
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chans := FadeOut(ctx, ch, 3)

	got := []T{}
	eg, _ := errgroup.WithContext(ctx)
	var mu sync.Mutex
	for _, ch := range chans {
		ch := ch
		eg.Go(func() error {
			for v := range ch {
				mu.Lock()
				got = append(got, v)
				mu.Unlock()
			}
			return nil
		})
	}
	eg.Wait()

	slices.Sort(items)
	slices.Sort(got)
	require.Equal(t, items, got)
}

func FuzzFadeOut(f *testing.F) {
	f.Add(1, 2, 3, 4, 5, 6, 7, 8, 9)
	f.Fuzz(func(t *testing.T, v1, v2, v3, v4, v5, v6, v7, v8, v9 int) {
		testFadeOut(t, []int{v1, v2, v3, v4, v5, v6, v7, v8, v9})
	})
}
