package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wiryls/pkg/config"
)

func TestDefault(t *testing.T) {
	assert := assert.New(t)
	{ // test write then read

		type Data struct {
			Message string `json:"orz"`
		}

		in := Data{
			Message: "OTZ",
		}

		tmp, err := ioutil.TempFile("", "*.json")
		assert.NoError(err)

		dst := tmp.Name()
		if err == nil {
			defer os.Remove(dst)
		}

		err = config.Save(&in, dst)
		assert.NoError(err)

		out := Data{}

		err = config.Load(&out, dst)
		assert.NoError(err)

		assert.Equal(in.Message, out.Message)
	}
}
