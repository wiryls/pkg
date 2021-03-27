package pipeline

import "github.com/cheekybits/genny/generic"

// I the input type
type I generic.Type

// O the output type
type O generic.Type

// Forward trys to fetch input from src, map it and send it to dst.
func Forward(source <-chan I, destination chan<- O, mapping func(I) O) {
	var (
		todo      []I
		done      []O
		peak      O
		input     = source
		output    chan<- O
		convert   chan struct{}
		available = make(chan struct{})
	)
	close(available)

	for input != nil || convert != nil || output != nil {
		select {
		case in, put := <-input:
			todo = append(todo, in)

			if !put {
				input = nil
			}
			if len(todo) != 0 {
				convert = available
			}

		case <-convert:
			done = append(done, mapping(todo[0]))
			todo = todo[1:]

			if len(todo) == 0 {
				convert = nil
			}
			if len(done) != 0 {
				output = destination
				peak = done[0]
			}

		case output <- peak:
			done = done[1:]

			if len(done) != 0 {
				peak = done[0]
			} else {
				output = nil
			}
		}
	}
}

// ForwardSlice trys to fetch input from src, map it and send it to dst.
func ForwardSlice(source <-chan []I, destination chan<- []O, mapping func(I) O) {
	var (
		todo      []I
		done      []O
		input     = source
		output    chan<- []O
		convert   chan struct{}
		available = make(chan struct{})
	)
	close(available)

	for input != nil || convert != nil || output != nil {
		select {
		case in, put := <-input:
			todo = append(todo, in...)

			if !put {
				input = nil
			}
			if len(todo) != 0 {
				convert = available
			}

		case <-convert:
			done = append(done, mapping(todo[0]))
			todo = todo[1:]

			if len(todo) == 0 {
				convert = nil
			}
			if len(done) != 0 {
				output = destination
			}

		case output <- done:
			done = nil
			output = nil
		}
	}
}
