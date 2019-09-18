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
		det := detail.New(nil, -1, nil)
		err := error(&det)
		assert.Error(err)
		assert.Empty(err.Error())
	}
	{
		msg := "whoops"
		det := detail.New(msg, -1, nil)
		err := error(&det)
		assert.Error(err)
		assert.Equal(msg, err.Error())
	}
	{
		inn := errors.New("inner")
		det := detail.New("whoops", -1, inn)
		err := error(&det)
		assert.Error(err)
		assert.Equal("whoops: inner", err.Error())
		assert.Equal(inn, errors.Unwrap(err))
	}
	{
		// TODO: test stack trace
	}
}
