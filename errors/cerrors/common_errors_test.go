package cerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wiryls/pkg/errors/cerrors"
)

func TestInternal(t *testing.T) {
	assert := assert.New(t)
	{
		err := cerrors.Internal("", nil)
		assert.NoError(err)
	}
	{
		msg := "internal"
		err := cerrors.Internal(msg, nil)
		assert.Error(err)

		internal := (*cerrors.InternalError)(nil)
		assert.True(errors.As(err, &internal))
		assert.True(errors.Is(err, cerrors.ErrInternal))
		assert.Equal(err, internal)
		assert.Equal(msg, err.Error())
	}
	{
		level0 := cerrors.NilArgument("orz")
		assert.Error(level0)
		assert.Nil(errors.Unwrap(level0))
		assert.Equal("argument `orz`, nil is not allowed", level0.Error())

		level1 := cerrors.Internal("inner", level0)
		assert.Error(level1)
		assert.Equal(level0, errors.Unwrap(level1))
		assert.Equal("inner: argument `orz`, nil is not allowed", level1.Error())

		level2 := cerrors.Internal("outter", level1)
		assert.Error(level2)
		assert.Equal("outter: inner: argument `orz`, nil is not allowed", level2.Error())

		var e0 *cerrors.InvalidArgumentError
		assert.True(errors.As(level2, &e0))
		assert.True(errors.Is(level2, cerrors.ErrInvalidArgument))
		assert.Equal(level0, e0)
		assert.Equal("orz", e0.Argument)

		var e1 *cerrors.InternalError
		assert.True(errors.As(level1, &e1))
		assert.True(errors.Is(level1, cerrors.ErrInternal))
		assert.Equal(level1, e1)

		var e2 *cerrors.InternalError
		assert.True(errors.As(level2, &e2))
		assert.True(errors.Is(level2, cerrors.ErrInternal))
		assert.Equal(level2, e2)
	}
	{
		type IAE = cerrors.InvalidArgumentError
		type IE = cerrors.InternalError
		iae := cerrors.InvalidArgument("bar", "baz")
		err := cerrors.Internal("foo", iae)
		assert.True(errors.Is(err, cerrors.ErrInvalidArgument))
		assert.True(errors.Is(err, cerrors.ErrInternal))

		if out := (*IE)(nil); errors.As(err, &out) {
			assert.Equal(err, out)
		} else {
			assert.Fail("`errors.As` failed")
		}

		if out := (*IAE)(nil); errors.As(err, &out) {
			assert.Equal(iae, out)
		} else {
			assert.Fail("`errors.As` failed")
		}
	}
}
