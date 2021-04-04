package flow

import (
	"sync"
	"sync/atomic"
)

// Feed could be used to group some tasks with a cancel function or something like that.
//
// Note: We may implements FeedWithChan \ FeedWithContext like this.
//
// Note: it is goroutine-safe and never copy after first use.
type Feed struct {
	Flow *Flow
	Sync uint32
	wait sync.WaitGroup
}

func (f *Feed) Push(fun func(u32 *uint32)) {
	if f.Flow != nil {
		f.wait.Add(1)
		f.Flow.Push(func() {
			defer f.wait.Done()
			fun(&f.Sync)
		})
	}
}

func (f *Feed) Send(sync uint32) {
	atomic.StoreUint32(&f.Sync, sync)
}

func (f *Feed) Wait() {
	f.wait.Wait()
}
