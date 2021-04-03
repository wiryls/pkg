package flow

// I the input type.
// Waiting for go 2.0 generic
type I interface{}

// O the output type.
type O interface{}

// Converter is used to convert I to O.
type Converter func(in I) (out O, ok bool)

// Forward trys to fetch input from src, map it and send it to dst.
func Forward(input <-chan I, output chan<- O, convert Converter) {
	var (
		todo        []I
		done        []O
		peak        O
		src         = input
		dst         chan<- O
		cvt         chan struct{}
		convertible = make(chan struct{})
	)
	close(convertible)

	for src != nil || cvt != nil || dst != nil {
		select {
		case in, put := <-src:
			if put {
				todo = append(todo, in)

				if cvt == nil {
					cvt = convertible // enable channel cvt
				}

			} else {
				src = nil // block channel src
			}

		case <-cvt:
			if item, ok := convert(todo[0]); ok {
				done = append(done, item)
				if dst == nil {
					dst = output // enable channel dst
					peak = done[0]
				}
			}

			if todo = todo[1:]; len(todo) == 0 {
				cvt = nil // block channel cvt
			}

		case dst <- peak:
			done = done[1:]

			if len(done) != 0 {
				peak = done[0]

			} else {
				dst = nil // block channel dst
			}
		}
	}
}

// ForwardSlice trys to fetch input from src, map it and send it to dst.
func ForwardSlice(input <-chan []I, output chan<- []O, convert Converter) {
	var (
		todo        []I
		done        []O
		src         = input
		dst         chan<- []O
		cvt         chan struct{}
		convertible = make(chan struct{})
	)
	close(convertible)

	for src != nil || cvt != nil || dst != nil {
		select {
		case in, put := <-src:

			if put {
				todo = append(todo, in...)

				if cvt == nil {
					cvt = convertible
				}

			} else {
				src = nil
			}

		case <-cvt:
			if that, ok := convert(todo[0]); ok {
				done = append(done, that)
				if dst == nil {
					dst = output
				}
			}

			if todo = todo[1:]; len(todo) == 0 {
				cvt = nil
			}

		case dst <- done:
			done = nil
			dst = nil
		}
	}
}
