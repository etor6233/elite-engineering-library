package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLoaderStampsOnce(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	l := NewLoader(s)
	var calls int32
	load := func(context.Context) ([]byte, error) {
		atomic.AddInt32(&calls, 1)
		time.Sleep(20 * time.Millisecond)
		return []byte("value"), nil
	}
	const n = 20
	var wg sync.WaitGroup
	values := make([]string, n)
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, err := l.Get(context.Background(), "k1", time.Minute, load)
			values[i] = string(v)
			errs[i] = err
		}(i)
	}
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected exactly 1 load, got %d", got)
	}
	for i := 0; i < n; i++ {
		if errs[i] != nil || values[i] != "value" {
			t.Fatalf("caller %d got %q/%v", i, values[i], errs[i])
		}
	}
}

func TestLoaderPopulatesCache(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	l := NewLoader(s)
	var calls int32
	load := func(context.Context) ([]byte, error) {
		atomic.AddInt32(&calls, 1)
		return []byte("value"), nil
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); err != nil {
		t.Fatal(err)
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected cache hit on second Get, loads=%d", got)
	}
}

func TestLoaderErrorNotCached(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	l := NewLoader(s)
	calls := 0
	boom := errors.New("boom")
	load := func(context.Context) ([]byte, error) {
		calls++
		return nil, boom
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); !errors.Is(err, boom) {
		t.Fatalf("expected boom on retry, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected reload after error, calls=%d", calls)
	}
	if s.Len() != 0 {
		t.Fatalf("error poisoned cache, len=%d", s.Len())
	}
}
