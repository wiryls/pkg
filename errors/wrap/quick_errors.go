package wrap

import (
	"fmt"

	"github.com/wiryls/pkg/errors/detail"
)

// Message wraps an error with a message.
//  - Return nil if `inner` is nil.
//  - It uses `fmt.Errorf(msg+": %w", inner)` to create an error.
func Message(inner error, msg string) error {
	if inner == nil {
		return nil
	}

	return fmt.Errorf(msg+": %w", inner)
}

// MessageStack wraps an error with a message, stack trace.
//  - Return nil if `inner` is nil.
//  - It creates a quick and simple detail.Detail error.
func MessageStack(inner error, msg string, skipCaller uint) error {
	if inner == nil {
		return nil
	}

	return detail.New(
		msg,
		detail.FlagInner(inner),
		detail.FlagStackTrace(skipCaller+1),
	)
}

// MessageAlias wraps an error with a message, alias.
//  - Return nil if `inner` is nil.
//  - It creates a quick and simple detail.Detail error.
func MessageAlias(inner error, msg string, alias error) error {
	if inner == nil {
		return nil
	}

	return detail.New(
		msg,
		detail.FlagInner(inner),
		detail.FlagAlias(alias),
	)
}

// MessageAliasStack wraps an error with a message, alias and stacktrace.
//  - Return nil if `inner` is nil.
//  - It creates a quick and simple detail.Detail error.
func MessageAliasStack(inner error, msg string, alias error, skipCaller uint) error {
	if inner == nil {
		return nil
	}

	return detail.New(
		msg,
		detail.FlagInner(inner),
		detail.FlagAlias(alias),
		detail.FlagStackTrace(skipCaller+1),
	)
}
