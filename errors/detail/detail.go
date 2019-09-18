package detail

import (
	"fmt"

	"golang.org/x/xerrors"
)

// New detail error. It is used to create other errors with stack trace and
// an inner error.
//  - `cause` is something like message.
//  - stack trace will be added if`skipCaller` >= 0.
//  - `inner` is an optional inner error of this error.
func New(cause interface{}, skipCaller int, inner error) Detail {
	var frame *xerrors.Frame
	if skipCaller >= 0 {
		caller := xerrors.Caller(skipCaller + 1)
		frame = &caller
	}
	return Detail{
		cause: cause,
		frame: frame,
		inner: inner,
	}
}

// Detail error contains a message, stack traces and an inner error.
// Usually, it is used by other errors and should always be created by `New`.
type Detail struct {
	cause interface{}
	frame *xerrors.Frame
	inner error
}

// Error implements the error interface.
func (e *Detail) Error() string {
	return fmt.Sprint(e)
}

// Format implements the Format method used for *Printf.
func (e *Detail) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

// FormatError formats this error with stack trace and inner errors.
func (e *Detail) FormatError(p xerrors.Printer) (next error) {
	if e.cause != nil {
		switch v := e.cause.(type) {
		case interface{ String() string }:
			p.Print(v.String())

		case interface{ Error() string }:
			p.Print(v.Error())

		default:
			p.Print(e.cause)
		}
	}
	if e.frame != nil {
		e.frame.Format(p)
	}
	return e.inner
}

// Unwrap this error and get the inner error.
func (e *Detail) Unwrap() error {
	return e.inner
}

// Cause is an alias of Unwrap. Maybe used in `pkg/errors`.
func (e *Detail) Cause() error {
	return e.inner
}
