package runner

// Determination is a helper for building service like object.
// Just create a `Determination` and bind it to a `Runnable`.
type Determination struct {
	shared

	runn Runnable
	exit chan struct{}
}

// Bind a `Callback` to this runner.
//
// WARNING: invoking it when state is StateRunning may cause blocked.
func (s *Determination) Bind(runnable Runnable) *Determination {
	defer s.lock.Unlock()
	/*_*/ s.lock.Lock()

	s.runn = runnable
	return s
}

// State of this `Runner`.
func (s *Determination) State() State {
	return s.stat.get()
}

// Run this `Runner`.
//  - Caller will be blocked until error happens or `Close` is called.
func (s *Determination) Run() (err error) {
	return s.run(s.prepare, s.booting, s.running, s.trigger, s.closing)
}

// WhileRunning do something if it is running. It provides a readonly
// channel to check if it stops.
func (s *Determination) WhileRunning(
	do func(exit <-chan struct{}) error,
) (err error) {
	return s.whilerunning(func() error { return do(s.exit) })
}

// CloseAsync sends a signal to close this runner asynchronously.
// This is a non-block version of `Close`.
//  - Only an `ErrUnexpectedState` with a `StateStopped` may be returned.
func (s *Determination) CloseAsync() error {
	return s.closeasync(s.trigger)
}

// Close this Determination and wait until stop running.
//  - Using it in `WhileRunning` will cause deadlock. Please use
//    `CloseAsync()` instead.
func (s *Determination) Close() error {
	return s.close(s.trigger)
}

// BeforeRunning is a default do nothing method. If we create our object
// like:
//
//     type service struct{
// 	     runner.Determination
//     }
//
// This function provide a default `BeforeRunning()` for `service`.
func (s *Determination) BeforeRunning(<-chan struct{}) error { return nil }

// AfterRunning is a default do nothing method. If we create our object
// like:
//
//     type service struct{
// 	     runner.Determination
//     }
//
// This function provide a default `AfterRunning()` for `service`.
func (s *Determination) AfterRunning() error { return nil }

func (s *Determination) prepare() {
	s.exit = make(chan struct{})
}

func (s *Determination) booting() (err error) {
	if s.runn != nil {
		err = s.runn.BeforeRunning(s.exit)
	}
	return
}

func (s *Determination) running() (err error) {
	if s.runn != nil {
		err = s.runn.Running(s.exit)
	}
	return
}

func (s *Determination) trigger() {
	close(s.exit) // as an triggering signal for exiting
}

func (s *Determination) closing() (err error) {
	if s.runn != nil {
		err = s.runn.AfterRunning()
	}
	return
}
