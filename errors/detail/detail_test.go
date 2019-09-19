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
		det := detail.New(nil, nil, -1, nil)
		err := error(&det)
		assert.Error(err)
		assert.Empty(err.Error())
	}
	{
		msg := "whoops"
		det := detail.New(msg, nil, -1, nil)
		err := error(&det)
		assert.Error(err)
		assert.Equal(msg, err.Error())
	}
	{
		inn := errors.New("inner")
		det := detail.New("whoops", nil, -1, inn)
		err := error(&det)
		assert.Error(err)
		assert.Equal("whoops: inner", err.Error())
		assert.Equal(inn, errors.Unwrap(err))
	}
	{
		inn := errors.New("inner")
		ali := errors.New("alias")
		mid := detail.New("bar", ali, 0, inn)
		out := detail.New("foo", nil, 0, &mid)
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
