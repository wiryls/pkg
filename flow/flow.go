package flow

import (
	"runtime"
	"sync"
)

// New create a flow to process some tasks.
func New(limit int) *Flow {
	return &Flow{limit: limit}
}

// Flow process something.
//
// Create it via:
// f := Flow{}, or
// f := &Flow{}, or
// f := flow.New(int)
//
// Note: it is goroutine-safe and never copy after first use.
type Flow struct {
	count int
	limit int
	tasks []func()
	mutex sync.Mutex
	inner sync.Mutex
}

// Push a task to the executor.
func (f *Flow) Push(task func()) {
	if f != nil && task != nil {
		defer f.mutex.Unlock()
		/*_*/ f.mutex.Lock()

		f.tasks = append(f.tasks, task)
		if f.limit <= 0 {
			f.limit = runtime.NumCPU()
		}
		if f.count == 0 {
			f.inner.Lock()
		}
		if f.count < f.limit {
			f.count++
			go f.low()
		}
	}
}

// Wait until all task done.
func (f *Flow) Wait() {
	defer f.inner.Unlock()
	/*_*/ f.inner.Lock()
}

func (f *Flow) low() {
	f.mutex.Lock()
	for len(f.tasks) != 0 {
		action := f.tasks[0]
		f.tasks = f.tasks[1:]
		f.mutex.Unlock()

		action()

		f.mutex.Lock()
	}
	f.count--
	if f.count == 0 {
		f.inner.Unlock()
	}
	f.mutex.Unlock()
}
