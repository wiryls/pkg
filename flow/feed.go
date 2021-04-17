package flow

import (
	"sync/atomic"
)

// Feed could be used to group some tasks with a cancel function or something like that.
//
// Note: We may implements FeedWithChan \ FeedWithContext like this.
//
// Note: it is goroutine-safe and never copy after first use.
type Feed struct {
	Flow
	Mark uint32
}

func (f *Feed) Push(fun func(u32 *uint32)) {
	if fun != nil {
		f.Flow.Push(func() { fun(&f.Mark) })
	}
}

func (f *Feed) Sync(sync uint32) {
	atomic.StoreUint32(&f.Mark, sync)
}
