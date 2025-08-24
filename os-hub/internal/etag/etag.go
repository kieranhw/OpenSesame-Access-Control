package etag

import (
	"sync"
	"sync/atomic"
	"time"
)

var version atomic.Uint64
var mu sync.Mutex
var waiters []chan struct{}

func Init() {
	version.Store(1)
}

func Current() uint64 {
	return version.Load()
}

func Bump() {
	version.Add(1)

	mu.Lock()
	defer mu.Unlock()
	for _, ch := range waiters {
		close(ch) // notify all waiters
	}
	waiters = nil
}

func Wait(current uint64, timeout time.Duration) bool {
	if Current() != current {
		return true // already changed
	}

	ch := make(chan struct{})
	mu.Lock()
	waiters = append(waiters, ch)
	mu.Unlock()

	select {
	case <-ch:
		return true // changed
	case <-time.After(timeout):
		return false // timeout
	}
}
