package etag

import (
	"sync"
	"sync/atomic"
	"time"
)

var version atomic.Uint64
var mu sync.Mutex
var waiters []chan struct{}
var lastBump atomic.Int64
var bumpScheduled atomic.Bool

func Init() {
	version.Store(1)
	lastBump.Store(time.Now().UnixNano())
}

func Current() uint64 {
	return version.Load()
}

// Bump increments the ETag version and notifies all waiters, rate limited to 1 sec
func Bump() {
	now := time.Now().UnixNano()
	last := lastBump.Load()

	// if last bump was within 1s, schedule a delayed bump if one is not already scheduled
	if time.Duration(now-last) < time.Second {
		if bumpScheduled.CompareAndSwap(false, true) {
			go func() {
				// wait the remaining time in the 1s window
				sleepFor := time.Second - time.Duration(now-last)
				if sleepFor < 0 {
					sleepFor = time.Second
				}
				time.Sleep(sleepFor)
				doBump()
				bumpScheduled.Store(false)
			}()
		}
		return
	}

	// otherwise bump immediately
	doBump()
}

func doBump() {
	now := time.Now().UnixNano()
	lastBump.Store(now)

	version.Add(1)

	mu.Lock()
	defer mu.Unlock()
	for _, ch := range waiters {
		close(ch)
	}
	waiters = nil
}

// block until the ETag changes or timeout expires
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
