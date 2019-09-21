package detail

import "golang.org/x/xerrors"

// Flag is used to add optional parameter to `detail.New(...)`
// By using those flags, we could custom error while `New` error.
// See: https://stackoverflow.com/a/26326418
type Flag func(*Detail)

// FlagStackTrace add a stack trace to our error.
//  - `skip` top n stack trace.
func FlagStackTrace(skip uint) Flag {
	return func(err *Detail) {
		caller := xerrors.Caller(int(skip + 3))
		err.frame = &caller
	}
}

// FlagAlias set an alias (error value) to this error. We could use
// `errors.Is(err, alias)` to test it.
func FlagAlias(err error) Flag {
	return func(detail *Detail) { detail.alias = err }
}

// FlagInner adds an inner error to this error. All error messages will
// be joined with a ": " when outputting.
func FlagInner(err error) Flag {
	return func(detail *Detail) { detail.inner = err }
}
