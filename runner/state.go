package runner

import (
	"strconv"
	"sync/atomic"
)

// State of runner.
type State uint32

// Determination States.
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

func (s *State) set(x State) {
	atomic.StoreUint32((*uint32)(s), uint32(x))
}

func (s *State) get() State {
	return State(atomic.LoadUint32((*uint32)(s)))
}
