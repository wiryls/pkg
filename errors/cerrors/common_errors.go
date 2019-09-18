package cerrors

import "github.com/wiryls/pkg/errors/detail"

// This file contains some samples showing how to customize new errors from
// our `detail.Detail`.

// InternalError is of type internal error.
//  - Please use `InternalError` to create it.
type InternalError struct{ detail.Detail }

// Internal creates a wrapped `ErrInternal` with detailed information.
//  - Return nil if both `message` and `inner` are nil.
func Internal(message string, inner error) (err error) {
	switch {
	case message == "" && inner == nil:
		return
	case message == "":
		message = "internal error"
		fallthrough
	default:
		return &InternalError{Detail: detail.New(message, 1, inner)}
	}
}

// InvalidArgumentError is  of type invalid argument.
//  - Please use `InvalidArgument` or `MaybeInvalidArgument` to create it.
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

	err := &InvalidArgumentError{Argument: argument, Reason: reason}
	err.Detail = detail.New(err, 1, nil)
	return err
}

// NilArgument is a simplified version of "InvalidArgument".
//  - Return nil if `name` is not nil.
func NilArgument(argument string) error {
	if argument == "" {
		return nil
	}

	reason := "nil is not allowed"
	err := &InvalidArgumentError{Argument: argument, Reason: reason}
	err.Detail = detail.New(err, 1, nil)
	return err
}

// MaybeInvalidArgument creates an `InvalidArgumentError` if cond is true.
func MaybeInvalidArgument(cond bool, argument string, reason string) error {
	if !cond {
		return nil
	}

	err := &InvalidArgumentError{Argument: argument, Reason: reason}
	err.Detail = detail.New(err, 1, nil)
	return err
}

// MaybeNilArgument creates an `InvalidArgumentError` if argument is not nil.
func MaybeNilArgument(argument interface{}, name string) error {
	if argument != nil {
		return nil
	}

	reason := "nil is not allowed"
	err := &InvalidArgumentError{Argument: name, Reason: reason}
	err.Detail = detail.New(err, 1, nil)
	return err
}
