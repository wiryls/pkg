package detail_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wiryls/pkg/errors/detail"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	{
		err := detail.New(nil)
		assert.Error(err)
		assert.Empty(err.Error())
	}
	{
		det := detail.Make(nil)
		err := error(&det)
		assert.Error(err)
		assert.Empty(err.Error())
	}
	{
		msg := "whoops"
		det := detail.Make(msg)
		err := error(&det)
		assert.Error(err)
		assert.Equal(msg, err.Error())
	}
	{
		inn := errors.New("inner")
		det := detail.Make("whoops", detail.FlagInner(inn))
		err := error(&det)
		assert.Error(err)
		assert.Equal("whoops: inner", err.Error())
		assert.Equal(inn, errors.Unwrap(err))
	}
	{
		inn := errors.New("inner")
		ali := errors.New("alias")
		mid := detail.Make("bar", detail.FlagAlias(ali), detail.FlagStackTrace(0), detail.FlagInner(inn))
		out := detail.Make("foo", detail.FlagStackTrace(0), detail.FlagInner(&mid))
		err := error(&out)
		assert.Error(err)
		assert.Equal("foo: bar: inner", err.Error())
		assert.Equal(&mid, errors.Unwrap(err))
		assert.True(errors.Is(err, ali))
	}
	{
		// TODO: test stack trace
	}
}
