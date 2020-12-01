package runner

import (
	"context"
)

// DeterminationWithContext is a helper for building service like object.
// Just create a `DeterminationWithContext` and bind it to a
// `RunnableWithContext`.
type DeterminationWithContext struct {
	shared

	// remain immutable
	rctx context.Context
	runn RunnableWithContext

	// always created when booting
	sctx context.Context
	exit func()
}

// Bind a `RunnableWithContext` to this runner.
//
// WARNING: invoking it when state is StateRunning may cause blocked.
func (s *DeterminationWithContext) Bind(ctx context.Context, runnable RunnableWithContext) *DeterminationWithContext {
	defer s.lock.Unlock()
	/*_*/ s.lock.Lock()

	s.rctx = ctx
	s.runn = runnable
	return s
}

// State of this `Runner`.
func (s *DeterminationWithContext) State() State {
	return s.stat.get()
}

// Run this `Runner`.
//  - Caller will be blocked until error happens or `Close` is called.
func (s *DeterminationWithContext) Run() (err error) {
	return s.run(s.prepare, s.booting, s.running, s.trigger, s.closing)
}

// WhileRunning do something if it is running. It provides a context
// o check if it stops.
func (s *DeterminationWithContext) WhileRunning(do func(context.Context) error) error {
	return s.whilerunning(func() error {
		select {
		case <-s.sctx.Done():
			return ErrRunnerIsClosing
		default:
			return do(s.sctx)
		}
	})
}

// CloseAsync sends a signal to close this runner asynchronously.
// This is a non-block version of `Close`.
//  - Only an `ErrUnexpectedState` with a `StateStopped` may be returned.
func (s *DeterminationWithContext) CloseAsync() error {
	return s.closeasync(s.trigger)
}

// Close this DeterminationWithContext and wait until stop running.
//  - Using it in `WhileRunning` will cause deadlock. Please use
//    `CloseAsync()` instead.
func (s *DeterminationWithContext) Close() error {
	return s.close(s.trigger)
}

// BeforeRunning is a default do nothing method. If we create our object
// like:
//
//     type service struct{
// 	     runner.DeterminationWithContext
//     }
//
// This function provide a default `BeforeRunning()` for `service`.
func (s *DeterminationWithContext) BeforeRunning(<-chan struct{}) error { return nil }

// AfterRunning is a default do nothing method. If we create our object
// like:
//
//     type service struct{
// 	     runner.DeterminationWithContext
//     }
//
// This function provide a default `AfterRunning()` for `service`.
func (s *DeterminationWithContext) AfterRunning() error { return nil }

func (s *DeterminationWithContext) prepare() {
	if s.rctx != nil {
		s.sctx, s.exit = context.WithCancel(s.rctx)
	} else {
		s.sctx, s.exit = context.WithCancel(context.Background())
	}
}

func (s *DeterminationWithContext) booting() (err error) {
	if s.runn != nil {
		err = s.runn.BeforeRunning(s.sctx)
	}
	return
}

func (s *DeterminationWithContext) running() (err error) {
	if s.runn != nil {
		err = s.runn.Running(s.sctx)
	}
	return
}

func (s *DeterminationWithContext) trigger() {
	s.exit() // as an triggering signal for exiting
}

func (s *DeterminationWithContext) closing() (err error) {
	if s.runn != nil {
		err = s.runn.AfterRunning()
	}
	return
}
