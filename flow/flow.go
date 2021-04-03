package flow

import (
	"runtime"
	"sync"
)

// New create a flow to process some tasks.
func New(limit int) *Flow {
	return &Flow{
		count: 0,
		limit: limit,
		tasks: []func(){},
		mutex: sync.Mutex{},
		group: sync.WaitGroup{},
	}
}

// Flow process something.
// User f := Flow{} or f := &Flow{} or f := flow.New(int) to create it.
type Flow struct {
	count int
	limit int
	tasks []func()
	mutex sync.Mutex
	group sync.WaitGroup
}

// Append a task to the executor.
func (f *Flow) Append(task func()) {
	if f != nil && task != nil {
		defer f.mutex.Unlock()
		/*_*/ f.mutex.Lock()

		f.tasks = append(f.tasks, task)
		if f.limit <= 0 {
			f.limit = runtime.NumCPU()
		}
		if f.count < f.limit {
			f.count++
			f.group.Add(1)
			go f.low()
		}
	}
}

// Wait until all task done.
func (f *Flow) Wait() {
	f.group.Wait()
}

func (f *Flow) low() {
	defer f.group.Done()

	f.mutex.Lock()
	for len(f.tasks) != 0 {
		action := f.tasks[0]
		f.tasks = f.tasks[1:]
		f.mutex.Unlock()

		action()

		f.mutex.Lock()
	}
	f.count--
	f.mutex.Unlock()
}
