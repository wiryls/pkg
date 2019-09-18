package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wiryls/pkg/errors/detail"
)

// Flag is used to add optional parameter to `errors.New(...)`
// By using those flags, we could custom error while `New` error.
// See: https://stackoverflow.com/a/26326418
type Flag func(*option)
type option struct {
	WithCaller bool
	SkipCaller int
	InnerError error
}

// FlagStackTrace add a stack trace to our error.
//  - `skip` top n stack trace.
func FlagStackTrace(skip uint) Flag {
	return func(o *option) {
		o.WithCaller = true
		o.SkipCaller = int(skip)
	}
}

// FlagInnerError adds an inner error to this error. All error messages will
// be joined with a ": " when outputting.
func FlagInnerError(inner error) Flag {
	return func(o *option) { o.InnerError = inner }
}

// New error form an `interface{}`.
//  - Return nil If `what` is nil.
//  - It is the same as `errors.New("string")` if no flags.
func New(what interface{}, flags ...Flag) error {

	// parse option
	o := &option{}
	for _, f := range flags {
		if f != nil {
			f(o)
		}
	}

	// create error
	var err error
	switch {
	case what == nil:
		err = nil

	case len(flags) == 0:
		fallthrough
	default:
		err = errors.New(fmt.Sprint(what))

	case o.InnerError != nil || o.WithCaller:
		ce := detail.New(what, o.SkipCaller, o.InnerError)
		err = &ce
	}

	// return error
	return err
}

// Multiple errors happend.
type Multiple []error

func (e Multiple) Error() string {
	message := make([]string, 0, len(e))
	for _, err := range e {
		if err != nil {
			message = append(message, err.Error())
		}
	}
	return strings.Join(message, "; ")
}

// TODO: add a `Format` method for `Multiple`.
