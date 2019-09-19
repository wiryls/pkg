package service

import (
	"errors"
	"fmt"

	"github.com/wiryls/pkg/errors/detail"
)

// Error values.
var (
	ErrUnexpectedState = errors.New("unexpected state")
)

// UnexpectedStateError is an error when unexpected states happen.
type UnexpectedStateError struct {
	Actual   State
	Expected []State
	detail.Detail
}

func (e *UnexpectedStateError) Error() string {
	var msg string
	switch len(e.Expected) {
	case 0:
		msg = fmt.Sprintf(
			"unexpected state %s", e.Actual)
	case 1:
		msg = fmt.Sprintf(
			"expect state %s but get %s", e.Expected[0], e.Actual)
	default:
		msg = fmt.Sprintf(
			"expect state %s but get %s", e.Expected, e.Actual)
	}
	return msg
}

// this struct is something like an internal namespace.
type oops struct{}

var whoops = oops{}

// UnexpectedServiceState creates a StateError with detailed information.
func (oops) UnexpectedServiceState(get State, expected ...State) error {
	err := &UnexpectedStateError{Actual: get, Expected: expected}
	err.Detail = detail.New(
		err,
		detail.FlagAlias(ErrUnexpectedState),
		detail.FlagStackTrace(1))
	return err
}
