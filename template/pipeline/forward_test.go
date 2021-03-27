package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wiryls/pkg/template/pipeline"
)

func TestForwardOne(t *testing.T) {
	assert := assert.New(t)

	{
		input := make(chan pipeline.I)
		output := make(chan pipeline.O)

		total := 100000
		result := make([]bool, total)

		identity := func(i pipeline.I) pipeline.O { return i }
		generate := func() {
			defer close(input)
			for i := 0; i < total; i++ {
				input <- i
			}
		}
		receive := func() {
			for i := 0; i < total; i++ {
				result[(<-output).(int)] = true
			}
		}

		go generate()
		go pipeline.Forward(input, output, identity)
		receive()

		for i := 0; i < total; i++ {
			assert.True(result[i])
		}
	}
}

func TestForwardSome(t *testing.T) {
	assert := assert.New(t)

	{
		input := make(chan []pipeline.I)
		output := make(chan []pipeline.O)

		total := 100000
		result := make([]bool, total)

		identity := func(i pipeline.I) pipeline.O { return i }
		generate := func() {
			batch := 100
			defer close(input)
			buffer := make([]pipeline.I, 0, batch)
			for i := 0; i < total; i++ {
				buffer = append(buffer, i)
				if i%batch == 0 || i == total-1 {
					input <- buffer
					buffer = make([]pipeline.I, 0, batch)
				}
			}
		}
		receive := func() {
			k := 0
			for o := range output {
				for _, x := range o {
					result[(x).(int)] = true
					if k++; k == total {
						return
					}
				}
			}
		}

		go generate()
		go pipeline.ForwardSlice(input, output, identity)
		receive()

		for i := 0; i < total; i++ {
			assert.True(result[i])
		}
	}
}
