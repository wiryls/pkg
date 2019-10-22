package cerrors

import (
	"errors"
	"reflect"

	"github.com/wiryls/pkg/errors/detail"
)

// This file contains some samples showing how to customize new errors from
// our `detail.Detail`. And also try to mix error values with error types.

// error values.
var (
	ErrInternal        = errors.New("internal")
	ErrUnimplement     = errors.New("unimplement")
	ErrInvalidArgument = errors.New("invalid argument")
)

// InternalError is of type internal error.
//  - Please use `InternalError` to create it.
type InternalError struct{ detail.Detail }

// Internal creates a wrapped `ErrInternal` with detailed information.
//  - Return nil if both `message` and `inner` are nil.
func Internal(message string, inner error) (err error) {
	switch {
	case message == "" && isNil(inner):
		return
	case message == "":
		message = "internal error"
		fallthrough
	default:
		return &InternalError{Detail: detail.Make(
			message,
			detail.FlagAlias(ErrInternal),
			detail.FlagInner(inner),
			detail.FlagStackTrace(1),
		)}
	}
}

// InvalidArgumentError is  of type invalid argument.
//  - Please use `InvalidArgument` or `TestInvalidArgument` to create it.
type InvalidArgumentError struct {
	Argument string
	Reason   string
	detail.Detail
}

// Error override the error interface to custom message.
func (e *InvalidArgumentError) Error() string {
	switch {
	case e.Argument != "" && e.Reason != "":
		return "argument `" + e.Argument + "`, " + e.Reason
	case e.Argument != "":
		return "argument `" + e.Argument + "` is invalid"
	case e.Reason != "":
		return "invalid argument, " + e.Reason
	default:
		return "invalid argument"
	}

	// Note:
	// https://gist.github.com/dtjm/c6ebc86abe7515c988ec#gistcomment-2794293
}

// InvalidArgument creates a ErrInvalidArgument.
func InvalidArgument(argument string, reason string) error {
	if argument == "" && reason == "" {
		return nil
	}
	return invalidArgumentError(argument, reason)
}

// TestInvalidArgument creates an `InvalidArgumentError` if cond is true.
func TestInvalidArgument(cond bool, argument string, reason string) error {
	if !cond {
		return nil
	}
	return invalidArgumentError(argument, reason)
}

// NilArgument is a simplified version of "InvalidArgument".
//  - Return nil if `name` is not nil.
func NilArgument(argument string) error {
	if argument == "" {
		return nil
	}
	reason := reasonNilArgument
	return invalidArgumentError(argument, reason)
}

// TestNilArgument creates an `InvalidArgumentError` if argument is not nil.
func TestNilArgument(argument interface{}, name string) error {
	if isNil(argument) {
		reason := reasonNilArgument
		return invalidArgumentError(name, reason)
	}
	return nil
}

// TestNilArgumentIfNoErr is similar to `TestNilArgument` but it tests if err
// is nil before checking the argument.
func TestNilArgumentIfNoErr(err error, argument interface{}, name string) error {

	if err != nil {
		return err
	}

	if isNil(argument) {
		reason := reasonNilArgument
		return invalidArgumentError(name, reason)
	}

	return nil
}

const reasonNilArgument = "nil is not allowed"

func invalidArgumentError(argument string, reason string) error {
	err := &InvalidArgumentError{
		Argument: argument,
		Reason:   reason}
	err.Detail = detail.Make(
		err,
		detail.FlagAlias(ErrInvalidArgument),
		detail.FlagStackTrace(2))
	return err
}

func isNil(i interface{}) bool {
	isNil := false

	if i == nil {
		isNil = true
	}

	if isNil == false {
		value := reflect.ValueOf(i)
		switch value.Kind() {
		case reflect.Ptr,
			reflect.Slice,
			reflect.Interface,
			reflect.Map,
			reflect.Chan,
			reflect.Func:
			isNil = value.IsNil()
		}
	}

	return isNil
}
