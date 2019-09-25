package detail

import (
	"fmt"

	"golang.org/x/xerrors"
)

// New creates a detail error with optional stacktrace, an inner error,
// alias and so on.
func New(what interface{}, flags ...Flag) error {
	detail := Detail{cause: what}
	fill(&detail, flags...)
	return &detail
}

// Make detail error. It is used to create other errors with stack trace and
// an inner error.
func Make(what interface{}, flags ...Flag) Detail {
	detail := Detail{cause: what}
	fill(&detail, flags...)
	return detail
}

// fill flags.
func fill(detail *Detail, flags ...Flag) {
	for _, f := range flags {
		if f != nil {
			f(detail)
		}
	}
}

// Detail error contains a message, stack traces and an inner error.
// Usually, it is used by other errors and should always be created by `Make`
// or `New`.
type Detail struct {
	cause interface{}
	alias error
	inner error
	frame *xerrors.Frame
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
	} else if e.alias != nil {
		p.Print(e.alias.Error())
	}

	if e.frame != nil {
		e.frame.Format(p)
	}
	return e.inner
}

// Is compares this error with another. The alias will also be compared.
func (e *Detail) Is(err error) bool {
	return e.alias == err || e == err
}

// Unwrap this error and get the inner error.
func (e *Detail) Unwrap() error {
	return e.inner
}

// Cause is an alias of Unwrap. May be needed in `pkg/errors`.
func (e *Detail) Cause() error {
	return e.inner
}
