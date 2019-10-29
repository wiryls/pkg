package lifecycle

import (
	"strconv"
	"sync"
	"sync/atomic"
)

// State of lifecycle.
type State uint32

// LifeCycle States.
const (
	StateStopped State = iota
	StateBooting
	StateRunning
	StateClosing
)

// StateToString convert state to string.
func (s State) String() string {
	switch s {
	case StateStopped:
		return "stopped"
	case StateBooting:
		return "booting"
	case StateRunning:
		return "running"
	case StateClosing:
		return "closing"
	default:
		return "unknown (" + strconv.Itoa(int(s)) + ")"
	}
}

func (s *State) atomicSet(x State) {
	atomic.StoreUint32((*uint32)(s), uint32(x))
}

func (s *State) atomicGet() State {
	return State(atomic.LoadUint32((*uint32)(s)))
}

// LifeCycle is a helper for building service like object.
// Just create a `LifeCycle` and bind it to a `Callback`.
type LifeCycle struct {
	runn Runnable

	lock sync.RWMutex
	stat State
	exit chan struct{}
	once sync.Once

	lerr sync.Mutex
	cerr error
}

// Bind a `Callback` to this runner.
//  - Invoke it while running may be blocked.
//  - See `lifecycle.Callback` for more details.
func (s *LifeCycle) Bind(runnable Runnable) *LifeCycle {
	defer s.lock.Unlock()
	/*_*/ s.lock.Lock()

	s.runn = runnable
	return s
}

// State of this `Runner`.
func (s *LifeCycle) State() State {
	return s.stat.atomicGet()
}

// Run this `Runner`.
//  - Caller will be blocked until error happens or `Close` is called.
func (s *LifeCycle) Run() (err error) {
	defer s.lerr.Unlock()
	/*_*/ s.lerr.Lock()
	/*_*/ s.cerr = nil

	// defer StateClosing -> StateStopped
	defer s.lock.Unlock()
	defer s.stat.atomicSet(StateStopped)
	defer s.lock.Lock()

	// StateStopped -> StateBooting -> StateRunning
	if err == nil {
		err = s.onBooting()
	}

	// StateRunning
	if err == nil {
		err = s.onRunning()
	}

	// StateRunning -> StateClosing
	if s.cerr == nil {
		s.cerr = s.onClosing()
	}

	if err == nil {
		err = s.cerr
	}

	return
}

// WhileRunning do something if it is running.
func (s *LifeCycle) WhileRunning(f func() error) (err error) {
	defer s.lock.RUnlock()
	/*_*/ s.lock.RLock()

	if err == nil && s.stat != StateRunning {
		err = whoops.UnexpectedState(s.stat, StateRunning)
	}

	if err == nil && f != nil {
		err = f()
	}

	return
}

// WhileRunningChan do something if it is running. It provides a readonly
// channel to check if it stops.
func (s *LifeCycle) WhileRunningChan(
	f func(cancel <-chan struct{}) error,
) (err error) {
	defer s.lock.RUnlock()
	/*_*/ s.lock.RLock()

	if err == nil && s.stat != StateRunning {
		err = whoops.UnexpectedState(s.stat, StateRunning)
	}

	if err == nil && f != nil {
		err = f(s.exit)
	}

	return
}

// CloseAsync sends a signal to close this runner asynchronously.
// This is a non-block version of `Close`.
//  - Only an `ErrUnexpectedState` with a `StateStopped` may be returned.
func (s *LifeCycle) CloseAsync() error {

	// fast check stat
	if s.stat.atomicGet() == StateStopped {
		return whoops.UnexpectedState(StateStopped)
	}

	// notify looping to exit.
	s.lock.RLock()
	if s.stat != StateStopped {
		s.once.Do(func() { close(s.exit) })
	}
	s.lock.RUnlock()

	return nil
}

// Close this LifeCycle and wait until stop running.
//  - Use it from its `Callback` may cause deadlock. Please use
//    `CloseAsync()` instead.
func (s *LifeCycle) Close() error {

	err := s.CloseAsync()
	if err != nil {
		return err
	}

	// wait unitl exit.
	defer s.lerr.Unlock()
	/*_*/ s.lerr.Lock()
	return s.cerr
}

// BeforeRunning is the default do nothing method. If we create our object
// like:
//
//     type service struct{
// 	     lifecycle.LifeCycle
//     }
//
// This function provide a default `BeforeRunning()` for `service`.
func (s *LifeCycle) BeforeRunning() error { return nil }

// AfterRunning is the default do nothing method. If we create our object
// like:
//
//     type service struct{
// 	     lifecycle.LifeCycle
//     }
//
// This function provide a default `AfterRunning()` for `service`.
func (s *LifeCycle) AfterRunning() error { return nil }

// from StateStopped to StateBooting and then StateRunning
func (s *LifeCycle) onBooting() (err error) {
	// fast check stat
	if stat := s.stat.atomicGet(); stat != StateStopped {
		return whoops.UnexpectedState(s.stat, StateStopped)
	}

	// slow path booting
	defer s.lock.Unlock()
	/*_*/ s.lock.Lock()

	if s.stat == StateStopped {
		// [0] set stat
		defer s.stat.atomicSet(StateRunning)
		/*_*/ s.stat.atomicSet(StateBooting)

		// [1] init
		s.exit = make(chan struct{})
		s.once = sync.Once{}

		// [2] callback
		if s.runn != nil {
			err = s.runn.BeforeRunning()
		}
	}

	return
}

// keep running, stat won't be changed.
func (s *LifeCycle) onRunning() (err error) {
	// fast check stat
	if stat := s.stat.atomicGet(); stat != StateRunning {
		return whoops.UnexpectedState(s.stat, StateRunning)
	}

	// slow path running
	defer s.lock.RUnlock()
	/*_*/ s.lock.RLock()

	if s.stat == StateRunning && s.runn != nil {
		err = s.runn.Running(s.exit)
	}

	return
}

// from StateRunning to StateClosing
func (s *LifeCycle) onClosing() (err error) {
	// fast check stat
	if stat := s.stat.atomicGet(); stat != StateRunning {
		return whoops.UnexpectedState(s.stat, StateRunning)
	}

	// slow path closing
	defer s.lock.Unlock()
	/*_*/ s.lock.Lock()

	if s.stat == StateRunning || s.stat == StateBooting {
		// [0] stat
		s.stat.atomicSet(StateClosing)

		// [1] release
		s.once.Do(func() { close(s.exit) })

		// [2] callback
		if s.runn != nil {
			err = s.runn.AfterRunning()
		}
	}

	return
}
