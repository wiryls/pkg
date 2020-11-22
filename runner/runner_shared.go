package runner

import (
	"sync"
)

type shared struct {
	lock sync.RWMutex
	stat State
	once sync.Once

	lerr sync.Mutex
	cerr error
}

func (s *shared) whilerunning(do func() error) (err error) {
	defer s.lock.RUnlock()
	/*_*/ s.lock.RLock()

	if err == nil && s.stat != StateRunning {
		err = whoops.UnexpectedState(s.stat, StateRunning)
	}

	if err == nil && do != nil {
		err = do()
	}

	return
}

func (s *shared) closeasync(do func()) error {

	// fast check stat
	if s.stat.get() == StateStopped {
		return whoops.UnexpectedState(StateStopped)
	}

	// notify looping to exit.
	s.lock.RLock()
	if s.stat != StateStopped {
		s.once.Do(do)
	}
	s.lock.RUnlock()

	return nil
}

// Close this Determination and wait until stop running.
//  - Use it from its `Callback` may cause deadlock. Please use
//    `CloseAsync()` instead.
func (s *shared) close(do func()) error {

	err := s.closeasync(do)
	if err != nil {
		return err
	}

	// wait unitl `run` ends.
	defer s.lerr.Unlock()
	/*_*/ s.lerr.Lock()
	return s.cerr
}

func (s *shared) run(
	// invoke by order
	prepare func(), // prepare the trigger
	booting func() error,
	running func() error,
	trigger func(), // trigger to exit
	closing func() error,
) (err error) {
	defer s.lerr.Unlock()
	/*_*/ s.lerr.Lock()
	/*_*/ s.cerr = nil

	// defer StateClosing -> StateStopped
	defer s.lock.Unlock()
	defer s.stat.set(StateStopped)
	defer s.lock.Lock()

	// StateStopped -> StateBooting -> StateRunning
	if err == nil {
		err = s.onBooting(prepare, booting)
	}

	// StateRunning
	if err == nil {
		err = s.onRunning(running)
	}

	// StateRunning -> StateClosing
	if s.cerr == nil {
		s.cerr = s.onClosing(trigger, closing)
	}

	if err == nil {
		err = s.cerr
	}

	return
}

// from StateStopped to StateBooting and then StateRunning
func (s *shared) onBooting(prepare func(), booting func() error) (err error) {
	// fast check stat
	if stat := s.stat.get(); stat != StateStopped {
		return whoops.UnexpectedState(s.stat, StateStopped)
	}

	// slow path booting
	defer s.lock.Unlock()
	/*_*/ s.lock.Lock()

	if s.stat == StateStopped {
		// [0] set stat
		defer s.stat.set(StateRunning)
		/*_*/ s.stat.set(StateBooting)

		// [1] init
		s.once = sync.Once{}
		if prepare != nil {
			prepare()
		}

		// [2] callback
		if booting != nil {
			err = booting()
		}
	}

	return
}

// keep running, stat won't be changed.
func (s *shared) onRunning(running func() error) (err error) {
	// fast check stat
	if stat := s.stat.get(); stat != StateRunning {
		return whoops.UnexpectedState(s.stat, StateRunning)
	}

	// slow path running
	defer s.lock.RUnlock()
	/*_*/ s.lock.RLock()

	if s.stat == StateRunning && running != nil {
		err = running()
	}

	return
}

// from StateRunning to StateClosing
func (s *shared) onClosing(trigger func(), closing func() error) (err error) {
	// fast check stat
	if stat := s.stat.get(); stat != StateRunning {
		return whoops.UnexpectedState(s.stat, StateRunning)
	}

	// slow path closing
	defer s.lock.Unlock()
	/*_*/ s.lock.Lock()

	if s.stat == StateRunning || s.stat == StateBooting {
		// [0] stat
		s.stat.set(StateClosing)

		// [1] release
		if trigger != nil {
			s.once.Do(trigger)
		}

		// [2] callback
		if closing != nil {
			err = closing()
		}
	}

	return
}
