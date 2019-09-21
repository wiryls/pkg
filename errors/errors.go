package errors

import (
	"github.com/wiryls/pkg/errors/detail"
)

// Wrap an error with message, alias and stacktrace.
//  - Return nil if `inner` is nil.
//  - It creates a quick and simple detail.Detail error.
func Wrap(inner error, msg string, alias error, skipCaller uint) error {
	if inner == nil {
		return nil
	}

	return detail.New(
		msg,
		detail.FlagInner(inner),
		detail.FlagAlias(alias),
		detail.FlagStackTrace(skipCaller),
	)
}
