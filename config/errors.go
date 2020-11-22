package config

import (
	"errors"
)

// errors
var (
	ErrEncoding = errors.New("encoding failed")
	ErrDecoding = errors.New("decoding failed")
	ErrReading  = errors.New("reading failed")
	ErrWriting  = errors.New("writing failed")
	ErrParsing  = errors.New("parsing failed")
)
