package flow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wiryls/pkg/flow"
)

func TestForwardOne(t *testing.T) {
	assert := assert.New(t)

	{
		input := make(chan flow.I)
		output := make(chan flow.O)

		total := 100000
		result := make([]bool, total)

		identity := func(i flow.I) flow.O { return i }
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
		go flow.Forward(input, output, identity)
		receive()

		for i := 0; i < total; i++ {
			assert.True(result[i])
		}
	}
}

func TestForwardSome(t *testing.T) {
	assert := assert.New(t)

	{
		input := make(chan []flow.I)
		output := make(chan []flow.O)

		total := 100000
		result := make([]bool, total)

		identity := func(i flow.I) flow.O { return i }
		generate := func() {
			batch := 100
			defer close(input)
			buffer := make([]flow.I, 0, batch)
			for i := 0; i < total; i++ {
				buffer = append(buffer, i)
				if i%batch == 0 || i == total-1 {
					input <- buffer
					buffer = make([]flow.I, 0, batch)
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
		go flow.ForwardSlice(input, output, identity)
		receive()

		for i := 0; i < total; i++ {
			assert.True(result[i])
		}
	}
}
